package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type ClientService interface {
	Get(ctx context.Context, UniqueCode string, Username string) ([]map[string]interface{}, error)
}

type ClientHandlers struct {
	c ClientService
}

func NewClientHandlers(c ClientService) *ClientHandlers {
	return &ClientHandlers{c: c}
}

func (h *ClientHandlers) Register(r *httprouter.Router) {
	r.GET("/api/client/:username", h.GetProducts)
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
