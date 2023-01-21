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
	log.Printf("Get products of %s\nUnique Code: %s", name, uniqueCode)

	prods, err := h.c.Get(r.Context(), uniqueCode, name)
	if err != nil {
		log.Printf("error: %+v", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	prodsMarshalled, err := json.Marshal(prods)
	if err != nil {
		log.Printf("error: %+v", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	_, err = w.Write(prodsMarshalled)
	if err != nil {
		log.Printf("error: %+v", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

}

type Request struct {
	Product struct {
		ID int `xml:"id"`
	} `xml:"product"`
	Options struct {
		Option []struct {
			ID     int    `xml:"id,attr"`
			Type   string `xml:"type,attr"`
			Option int    `xml:"option"`
		} `xml:"option"`
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

	log.Printf("%+v", input)

	check, err := h.c.Check(request.Context(), input.Options.Option[0].Option)
	if err != nil {
		log.Print(err)
		return
	}
	if check == false {
		/*checkTwo, err := h.c.Check(request.Context(), input.Options.Option[1].Option)
		if err != nil {
			return
		}
		if checkTwo != true {
			w.WriteHeader(400)
			w.Write([]byte(`"error": "we haven't this products"`))
		}*/
	}
	if check == true {
		log.Println("Request for search key: ", input.Product.ID)
		w.WriteHeader(200)
	}
}
