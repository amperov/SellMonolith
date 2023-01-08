package client

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type ClientService interface {
	Get(ctx context.Context, UniqueCode string, Username string) (map[string]interface{}, error)
}

type ClientHandlers struct {
	c ClientService
}

func NewClientHandlers(c ClientService) *ClientHandlers {
	return &ClientHandlers{c: c}
}

func (h *ClientHandlers) Register(r *httprouter.Router) {
	r.GET("/client/:username", h.GetProducts)
}

func (h *ClientHandlers) GetProducts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	uniqueCode := r.URL.Query().Get("uniquecode")
	name := params.ByName("username")

	prods, err := h.c.Get(r.Context(), uniqueCode, name)
	if err != nil {
		return
	}

	prodsMarshalled, err := json.Marshal(prods)
	if err != nil {
		return
	}

	_, err = w.Write(prodsMarshalled)
	if err != nil {
		return
	}

}
