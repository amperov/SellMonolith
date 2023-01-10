package subcategory

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
)

type SubcatService interface {
	Create(ctx context.Context, m map[string]interface{}, UserID, CatID int) (int, error)
	Update(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID int) (int, error)
	Delete(ctx context.Context, UserID, CatID, SubCatID int) error
	Get(ctx context.Context, UserID, CatID, SubCatID int) (map[string]interface{}, error)
}

type ProductService interface {
	GetAll(ctx context.Context, UserID, CatID, SubCatID int) ([]map[string]interface{}, error)
}

type SubcategoryHandler struct {
	ware auth.MiddleWare
	sc   SubcatService
	ps   ProductService
}

func NewSubcategoryHandler(ware auth.MiddleWare, sc SubcatService, ps ProductService) *SubcategoryHandler {
	return &SubcategoryHandler{ware: ware, sc: sc, ps: ps}
}

func (h *SubcategoryHandler) Register(r *httprouter.Router) {
	r.GET("/seller/category/:cat_id/subcategory/:subcat_id", h.ware.IsAuth(h.GetSubcategory))
	r.POST("/seller/category/:cat_id", h.ware.IsAuth(h.CreateSubcategory))
	r.PATCH("/seller/category/:cat_id/subcategory/:subcat_id", h.ware.IsAuth(h.UpdateSubcategory))
	r.DELETE("/seller/category/:cat_id/subcategory/:subcat_id", h.ware.IsAuth(h.DeleteSubcategory))
}

func (h *SubcategoryHandler) CreateSubcategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input CreateSubcatInput

	UserID := r.Context().Value("user_id").(int)

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		log.Println(err)
		return
	}

	input.CatID = CatID

	id, err := h.sc.Create(r.Context(), input.ToMap(), UserID, CatID)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "subcategory with ID %d created"}`, id)))
	if err != nil {
		log.Println(err)
		return
	}

}

func (h *SubcategoryHandler) UpdateSubcategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input UpdateSubcatInput

	UserID := r.Context().Value("user_id").(int)
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		return
	}
	input.CatID = CatID
	id, err := h.sc.Update(r.Context(), input.ToMap(), UserID, CatID, SubCatID)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "subcategory with ID %d updated"}`, id)))
	if err != nil {
		return
	}
}

func (h *SubcategoryHandler) DeleteSubcategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	UserID := r.Context().Value("user_id").(int)

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

	err = h.sc.Delete(r.Context(), UserID, CatID, SubCatID)
	if err != nil {
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "subcategory deleted"}`)))
	if err != nil {
		return
	}
}

func (h *SubcategoryHandler) GetSubcategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Get Subcategory")
	UserID := r.Context().Value("user_id").(int)

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		log.Println(catID)
		return
	}

	subcatID := params.ByName("subcat_id")
	SubCatID, err := strconv.Atoi(subcatID)
	if err != nil {
		log.Println(err)
		return
	}
	cat, err := h.sc.Get(r.Context(), UserID, CatID, SubCatID)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(cat)
	products, err := h.ps.GetAll(r.Context(), UserID, CatID, SubCatID)
	if err != nil {
		log.Println(err)
	}

	cat["products"] = products

	catMarshalled, err := json.Marshal(cat)
	if err != nil {
		return
	}
	_, err = w.Write(catMarshalled)
	if err != nil {
		return
	}
}
