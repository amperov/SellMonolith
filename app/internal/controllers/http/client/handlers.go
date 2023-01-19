package client

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

type ClientService interface {
	Get(ctx context.Context, UniqueCode string, Username string) ([]map[string]interface{}, error)
	Check(ctx context.Context, ItemID int) (bool, error)
}

type ClientHandlers struct {
	c ClientService
}

func NewClientHandlers(c ClientService) *ClientHandlers {
	return &ClientHandlers{c: c}
}

func (h *ClientHandlers) Register(r *httprouter.Router) {
	r.GET("/api/client/:username", h.GetProducts)
	r.POST("/api/precheck", h.PreCheck)

}

func (h *ClientHandlers) GetProducts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	uniqueCode := r.URL.Query().Get("uniquecode")
	name := params.ByName("username")

	prods, err := h.c.Get(r.Context(), uniqueCode, name)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	prodsMarshalled, err := json.Marshal(prods)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	_, err = w.Write(prodsMarshalled)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

}

type Request struct {
	Product struct {
		ID int `xml:"id"`
	} `xml:"product"`
	Options []struct {
		ID    int    `xml:"id"`
		Type  string `xml:"type"`
		Value string `xml:"value"`
	} `xml:"options"`
}

func (h *ClientHandlers) PreCheck(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var input Request

	all, err := io.ReadAll(request.Body)
	if err != nil {
		return
	}

	log.Println(string(all))

	err = xml.Unmarshal(all, &input)
	if err != nil {
		return
	}

	log.Println(input)
	check, err := h.c.Check(request.Context(), input.Options[0].ID)
	if err != nil {
		return
	}
	if check == false {
		w.WriteHeader(400)
		w.Write([]byte(`"error": "we haven't this products"`))
		return
	}
	if check == true {
		log.Println("Request for search key: ", input.Product.ID)
		w.WriteHeader(200)
	}
}
