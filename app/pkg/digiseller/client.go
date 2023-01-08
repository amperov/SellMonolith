package digiseller

import (
	"Selling/app/internal/storage/seller"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Body struct {
	SellerID  int
	Timestamp int64
	Sign      string
}

type Resp struct {
	Retval   int    `json:"retval,omitempty"`
	Token    string `json:"token,omitempty"`
	SellerID int    `json:"seller_id,omitempty"`
	Valid    string `json:"valid,omitempty"`
}

type SellerStore interface {
	GetSellerInfo(ctx context.Context, Username string) (int, string, error)
	SetTransaction(ctx context.Context, model map[string]interface{}) error
}

type CategoryStore interface {
	GetID(ctx context.Context, CategoryName string) (int, error)
}

type HistoryStore interface {
}
type SubcategoryStore interface {
	GetID(ctx context.Context, SubcategoryName string, CategoryID int) (int, error)
}

type ProductStore interface {
	GetSomeProducts(ctx context.Context, SubcatID int, Count int) ([]map[string]interface{}, error)
	DeleteOne(ctx context.Context, ProdID int) error
}

type DigiClient struct {
	c  CategoryStore
	ss SubcategoryStore
	s  SellerStore
	p  ProductStore
}

func NewDigiClient(c CategoryStore, ss SubcategoryStore, s SellerStore, p ProductStore) *DigiClient {
	return &DigiClient{c: c, ss: ss, s: s, p: p}
}

func (c *DigiClient) Auth(ctx context.Context, Username string) string {
	var body Body

	var respStr Resp
	//Searching SellerID by Username or integrate Seller Storage
	id, key, err := c.s.GetSellerInfo(ctx, Username)
	if err != nil {
		logrus.Debugf("get seller ID: %v", err)
		return ""
	}

	body.SellerID = id
	body.Timestamp = time.Now().Unix()

	hash := sha256.New()
	hash.Write([]byte(key + strconv.Itoa(int(body.Timestamp))))
	body.Sign = hex.EncodeToString(hash.Sum(nil))

	BodyMarshalled, err := json.Marshal(body)
	if err != nil {
		logrus.Debugf("marshalling error: %v", err)
		return ""
	}
	reader := bytes.NewReader(BodyMarshalled)

	resp, err := http.Post("https://api.digiseller.ru/api/apilogin", "application/json", reader)
	if err != nil {
		logrus.Debugf("http.Post error: %v", err)
		return ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Debugf("ReadAll error: %v", err)
		return ""
	}

	err = json.Unmarshal(respBody, &respStr)
	if err != nil {
		logrus.Debugf("unmarshal error: %v", err)
		return ""
	}
	return respStr.Token
}

func (c *DigiClient) GetProducts(ctx context.Context, UniqueCode, Token string) ([]map[string]interface{}, error) {
	var input DigiInput
	var tran seller.Transaction

	resp, err := http.Get(fmt.Sprintf("https://api.digiseller.ru/api/purchases/unique-code/%s?token=%s", UniqueCode, Token))
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		return nil, err
	}

	CatID, err := c.c.GetID(ctx, input.Options[0].Name)
	if err != nil {
		return nil, err
	}

	SubcatID, err := c.ss.GetID(ctx, input.Options[0].Value, CatID)
	if err != nil {
		return nil, err
	}

	products, err := c.p.GetSomeProducts(ctx, SubcatID, input.CntGoods)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(products); i++ {
		prod := products[i]
		err := c.p.DeleteOne(ctx, prod["id"].(int))
		if err != nil {
			return nil, err
		}
	}

	tran.UniqueInv = input.Inv
	tran.UniqueCode.UniqueCode = UniqueCode
	tran.UniqueCode.DateConfirmed = input.UniqueCodeState.DateConfirmed
	tran.UniqueCode.DateDelivery = input.UniqueCodeState.DateDelivery
	tran.UniqueCode.DateCheck = input.UniqueCodeState.DateCheck
	tran.CountGoods = input.CntGoods
	tran.Amount = int(input.Amount)
	tran.AmountUSD = int(input.AmountUsd)
	tran.Category = input.Options[0].Name
	tran.Subcategory = input.Options[0].Value
	tran.ClientEmail = input.Email
	tran.Profit, err = strconv.Atoi(input.Profit)

	var contents []string

	for i := 0; i < len(products); i++ {
		prod := products[i]
		contents = append(contents, prod["content"].(string))
	}
	tran.Content = strings.Join(contents, " \n")

	err = c.s.SetTransaction(ctx, tran.ToMap())
	if err != nil {
		return nil, err
	}

	return products, nil
}
