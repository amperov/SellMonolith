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

type Precheck struct {
	Request struct {
		Product struct {
			ID int `json:"id"`
		} `json:"product"`
	} `json:"request"`
}

func (h *ClientHandlers) PreCheck(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var input Precheck

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
	check, err := h.c.Check(request.Context(), input.Request.Product.ID)
	if err != nil {
		return
	}
	if check == false {
		w.WriteHeader(400)
		w.Write([]byte(`"error": "we haven't this products"`))
		return
	}
	if check == true {
		log.Println("Request for search key: ", input.Request.Product.ID)
		w.WriteHeader(200)
	}
}
