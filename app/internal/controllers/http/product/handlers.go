package product

import (
	"Selling/app/pkg/auth"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
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
	r.POST("/seller/category/:cat_id/subcategory/:subcat_id/one", h.ware.IsAuth(h.CreateOne))
	r.POST("/seller/category/:cat_id/subcategory/:subcat_id/many", h.ware.IsAuth(h.CreateMany))
	r.PATCH("/seller/category/:cat_id/subcategory/:subcat_id/products/:product_id", h.ware.IsAuth(h.UpdateProduct))
	r.DELETE("/seller/category/:cat_id/subcategory/:subcat_id/products/:product_id", h.ware.IsAuth(h.DeleteProduct))
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input UpdateProductInput

	userID := fmt.Sprintf("%v", r.Context().Value("user_id"))
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return
	}

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		return
	}

	subcatID := params.ByName("subcat_id")
	SubCatID, err := strconv.Atoi(subcatID)
	if err != nil {
		return
	}

	prodID := params.ByName("product_id")
	ProductID, err := strconv.Atoi(prodID)
	if err != nil {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		return
	}

	id, err := h.ps.Update(r.Context(), input.ToMap(), UserID, CatID, SubCatID, ProductID)
	if err != nil {
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "product with ID %d updated"}`, id)))
	if err != nil {
		return
	}
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := fmt.Sprintf("%v", r.Context().Value("user_id"))
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return
	}

	catID := params.ByName(":cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		return
	}

	subcatID := params.ByName(":subcat_id")
	SubCatID, err := strconv.Atoi(subcatID)
	if err != nil {
		return
	}

	prodID := params.ByName(":product_id")
	ProductID, err := strconv.Atoi(prodID)
	if err != nil {
		return
	}
	err = h.ps.Delete(r.Context(), UserID, CatID, SubCatID, ProductID)
	if err != nil {
		return
	}

	_, err = w.Write([]byte(`{"success" : "product deleted"}`))
	if err != nil {
		return
	}
}

func (h *ProductHandler) CreateOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input CreateProductInput

	userID := fmt.Sprintf("%v", r.Context().Value("user_id"))
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return
	}

	catID := params.ByName(":cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		return
	}

	subcatID := params.ByName(":subcat_id")
	SubCatID, err := strconv.Atoi(subcatID)
	if err != nil {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		return
	}

	id, err := h.ps.Create(r.Context(), input.ToMap(), UserID, CatID, SubCatID)
	if err != nil {
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "product with ID %d created"}`, id)))
	if err != nil {
		return
	}
}

func (h *ProductHandler) CreateMany(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input CreateProductsInput

	userID := fmt.Sprintf("%v", r.Context().Value("user_id"))
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return
	}

	catID := params.ByName(":cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		return
	}

	subcatID := params.ByName(":subcat_id")
	SubCatID, err := strconv.Atoi(subcatID)
	if err != nil {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &input)
	if err != nil {
		return
	}

	//TODO Need Testing
	count := 0
	for i := 0; i < len(input.Products); i++ {

		_, err := h.ps.Create(r.Context(), input.Products[i].ToMap(), UserID, CatID, SubCatID)
		if err != nil {
			continue
		}
		count++
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "%d products created"}`, count)))
	if err != nil {
		return
	}
}
