package entity

import "github.com/sirupsen/logrus"

type TransactionModel struct {
	ID          int        `json:"id,omitempty"`
	Category    string     `json:"category,omitempty"`
	Subcategory string     `json:"subcategory,omitempty"`
	UniqueCode  UniqueCode `json:"unique_code,omitempty"`
	ClientEmail string     `json:"client_email,omitempty"`
	DateCheck   string     `json:"date_check,omitempty"`
	Content     string     `json:"content,omitempty"`
	Profit      int        `json:"profit,omitempty"`
	Amount      int        `json:"amount,omitempty"`
	AmountUSD   int        `json:"amount_usd,omitempty"`
	CountGoods  int        `json:"count_goods,omitempty"`
	UniqueInv   int        `json:"unique_inv"`
	UserID      int        `json:"user_id"`
}

type UniqueCode struct {
	UniqueCode    string `json:"unique_code,omitempty"`
	State         int    `json:"state,omitempty"`
	DateCheck     string `json:"date_check,omitempty"`
	DateDelivery  string `json:"date_delivery,omitempty"`
	DateConfirmed string `json:"date_confirmed,omitempty"`
}

func (t *TransactionModel) ToMap() map[string]interface{} {
	var m map[string]interface{}
	m["id"] = t.ID
	m["category_name"] = t.Category
	m["subcategory"] = t.Subcategory
	m["client_email"] = t.ClientEmail
	m["content_key"] = t.Content
	m["amount"] = t.Amount
	m["profit"] = t.Profit
	m["amount_usd"] = t.AmountUSD
	m["count"] = t.CountGoods
	m["unique_inv"] = t.UniqueInv
	m["user_id"] = t.UserID

	m["unique_code"] = t.UniqueCode.UniqueCode
	m["date_check"] = t.UniqueCode.DateCheck
	m["date_delivery"] = t.UniqueCode.DateDelivery
	m["date_confirmed"] = t.UniqueCode.DateConfirmed
	if t.UniqueCode.State == 1 {
		m["state"] = "unique code не проверен!"
	} else if t.UniqueCode.State == 2 {
		m["state"] = "товар доставлен, доставка не подтверждена и не опровергнута"
	} else if t.UniqueCode.State == 3 {
		m["state"] = "товар доставлен, доставка подтверждена"
	} else if t.UniqueCode.State == 4 {
		m["state"] = "товар доставлен, но отвергнут"
	} else if t.UniqueCode.State == 5 {
		m["state"] = "уникальный код проверен, товар не доставлен"
	}

	logrus.Debugf("Transaction info: %v", m)

	return m
}
