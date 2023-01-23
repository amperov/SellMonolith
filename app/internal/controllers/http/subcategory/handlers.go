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
	r.GET("/api/seller/category/:cat_id/subcategory/:subcat_id", h.ware.IsAuth(h.GetSubcategory))
	r.POST("/api/seller/category/:cat_id", h.ware.IsAuth(h.CreateSubcategory))
	r.PATCH("/api/seller/category/:cat_id/subcategory/:subcat_id", h.ware.IsAuth(h.UpdateSubcategory))
	r.DELETE("/api/seller/category/:cat_id/subcategory/:subcat_id", h.ware.IsAuth(h.DeleteSubcategory))
}

func (h *SubcategoryHandler) CreateSubcategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input CreateSubcatInput

	UserID := r.Context().Value("user_id").(int)

	if UserID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		w.WriteHeader(400)
		log.Println(err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		log.Println(err)
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		w.WriteHeader(400)
		log.Println(err)
		return
	}

	input.CatID = CatID

	id, err := h.sc.Create(r.Context(), input.ToMap(), UserID, CatID)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "subcategory with ID %d created"}`, id)))
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	w.WriteHeader(201)
}

func (h *SubcategoryHandler) UpdateSubcategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input UpdateSubcatInput

	UserID := r.Context().Value("user_id").(int)

	if UserID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	subcatID := params.ByName("subcat_id")
	SubCatID, err := strconv.Atoi(subcatID)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		return
	}
	input.CatID = CatID
	id, err := h.sc.Update(r.Context(), input.ToMap(), UserID, CatID, SubCatID)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "subcategory with ID %d updated"}`, id)))
	if err != nil {
		return
	}
	w.WriteHeader(200)
}

func (h *SubcategoryHandler) DeleteSubcategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	UserID := r.Context().Value("user_id").(int)

	if UserID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	subcatID := params.ByName("subcat_id")
	SubCatID, err := strconv.Atoi(subcatID)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	err = h.sc.Delete(r.Context(), UserID, CatID, SubCatID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "subcategory deleted"}`)))
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func (h *SubcategoryHandler) GetSubcategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	UserID := r.Context().Value("user_id").(int)

	if UserID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(400)
		log.Println(catID)
		return
	}

	subcatID := params.ByName("subcat_id")
	SubCatID, err := strconv.Atoi(subcatID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(400)
		log.Println(err)
		return
	}
	cat, err := h.sc.Get(r.Context(), UserID, CatID, SubCatID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(500)
		log.Println(err)
		return
	}

	products, err := h.ps.GetAll(r.Context(), UserID, CatID, SubCatID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(500)
		log.Println(err)
	}

	cat["products"] = products

	catMarshalled, err := json.Marshal(cat)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(500)
		return
	}
	_, err = w.Write(catMarshalled)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(500)
		return
	}
}
