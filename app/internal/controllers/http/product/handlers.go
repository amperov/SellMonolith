package product

import (
	"Selling/app/pkg/auth"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ProductService interface {
	Create(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID int) (int, error)
	Update(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID, ProdID int) (int, error)
	Delete(ctx context.Context, UserID, CatID, SubCatID, ProdID int) error
}
type ProductHandler struct {
	ware auth.MiddleWare
	ps   ProductService
}

func NewProductHandler(ware auth.MiddleWare, ps ProductService) *ProductHandler {
	return &ProductHandler{ware: ware, ps: ps}
}

func (h *ProductHandler) Register(r *httprouter.Router) {
	r.POST("/api/seller/category/:cat_id/subcategory/:subcat_id/one", h.ware.IsAuth(h.CreateOne))
	r.POST("/api/seller/category/:cat_id/subcategory/:subcat_id/many", h.ware.IsAuth(h.CreateMany))
	r.PATCH("/api/seller/category/:cat_id/subcategory/:subcat_id/products/:product_id", h.ware.IsAuth(h.UpdateProduct))
	r.DELETE("/api/seller/category/:cat_id/subcategory/:subcat_id/products/:product_id", h.ware.IsAuth(h.DeleteProduct))
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input UpdateProductInput

	UserID := r.Context().Value("user_id").(int)

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(400)
		return
	}

	subcatID := params.ByName("subcat_id")
	SubCatID, err := strconv.Atoi(subcatID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(400)
		return
	}

	prodID := params.ByName("product_id")
	ProductID, err := strconv.Atoi(prodID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(400)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(400)
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(400)
		return
	}

	input.SubCatID = SubCatID
	id, err := h.ps.Update(r.Context(), input.ToMap(), UserID, CatID, SubCatID, ProductID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(400)
		log.Println(err)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "product with ID %d updated"}`, id)))
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(400)
		return
	}
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	UserID := r.Context().Value("user_id").(int)

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	subcatID := params.ByName("subcat_id")
	SubCatID, err := strconv.Atoi(subcatID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	prodID := params.ByName("product_id")
	ProductID, err := strconv.Atoi(prodID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}
	err = h.ps.Delete(r.Context(), UserID, CatID, SubCatID, ProductID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(500)
		return
	}

	_, err = w.Write([]byte(`{"success" : "product deleted"}`))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}
}

func (h *ProductHandler) CreateOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input CreateProductInput

	UserID := r.Context().Value("user_id").(int)

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	subcatID := params.ByName("subcat_id")
	SubCatID, err := strconv.Atoi(subcatID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	input.SubCatID = SubCatID
	id, err := h.ps.Create(r.Context(), input.ToMap(), UserID, CatID, SubCatID)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		log.Println(err)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "product with ID %d created"}`, id)))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}
}

func (h *ProductHandler) CreateMany(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input CreateProductInput

	UserID := r.Context().Value("user_id").(int)

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	subcatID := params.ByName("subcat_id")
	SubCatID, err := strconv.Atoi(subcatID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		log.Println(err)
		return
	}
	Products := strings.Split(input.Content, "\n")

	count := 0
	for i := 0; i < len(Products); i++ {
		input.Content = Products[i]
		input.SubCatID = SubCatID
		_, err := h.ps.Create(r.Context(), input.ToMap(), UserID, CatID, SubCatID)
		if err != nil {
			log.Println(err)
			continue
		}
		count++
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "%d products created"}`, count)))
	if err != nil {
		w.WriteHeader(5)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}
}
