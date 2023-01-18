package category

import (
	"Selling/app/pkg/auth"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"strconv"
)

//Replace interfaces to Structures

type CatService interface {
	Create(ctx context.Context, m map[string]interface{}) (int, error)
	Update(ctx context.Context, m map[string]interface{}, UserID int, CatID int) (int, error)
	Delete(ctx context.Context, UserID int, CatID int) error
	GetOne(ctx context.Context, UserID int, CatID int) (map[string]interface{}, error)
	GetAll(ctx context.Context, UserID int) ([]map[string]interface{}, error)
}
type SubCatService interface {
	GetAll(ctx context.Context, UserID, CatID int) ([]map[string]interface{}, error)
	GetCount(ctx context.Context, UserID, CatID int) (int, error)
}

type ProductService interface {
	GetCount(ctx context.Context, UserID, CatID, SubCatID int) (int, error)
}

type CategoryHandler struct {
	ware auth.MiddleWare
	cat  CatService
	sc   SubCatService
	ps   ProductService
}

func NewCategoryHandler(ware auth.MiddleWare, cat CatService, sc SubCatService, ps ProductService) *CategoryHandler {
	return &CategoryHandler{ware: ware, cat: cat, sc: sc, ps: ps}
}

func (h *CategoryHandler) Register(r *httprouter.Router) {
	r.GET("/api/seller/category", h.ware.IsAuth(h.GetAllCategory))
	r.GET("/api/seller/category/:cat_id", h.ware.IsAuth(h.GetCategory))
	r.POST("/api/seller/category", h.ware.IsAuth(h.CreateCategory))
	r.PATCH("/api/seller/category/:cat_id", h.ware.IsAuth(h.UpdateCategory))
	r.DELETE("/api/seller/category/:cat_id", h.ware.IsAuth(h.DeleteCategory))
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input CreateCategoryInput
	UserID := r.Context().Value("user_id").(int)

	if UserID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		log.Println(err)
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		log.Println(err)
		return
	}

	input.UserID = UserID

	id, err := h.cat.Create(r.Context(), input.ToMap())

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "category with ID %d created"}`, id)))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		log.Println(err)
		return
	}
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var upd UpdateCategoryInput
	UserID := r.Context().Value("user_id").(int)

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
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

	err = json.Unmarshal(body, &upd)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	id, err := h.cat.Update(r.Context(), upd.ToMap(), UserID, CatID)

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "category with ID %d updated"}`, id)))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	UserID := r.Context().Value("user_id").(int)

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	err = h.cat.Delete(r.Context(), UserID, CatID)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "category with deleted"}`)))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	UserID := r.Context().Value("user_id").(int)

	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	cat, err := h.cat.GetOne(r.Context(), UserID, CatID)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	subcategories, err := h.sc.GetAll(r.Context(), UserID, CatID)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		log.Println(err)
	}
	for _, subcategory := range subcategories {
		count, err := h.ps.GetCount(r.Context(), UserID, CatID, subcategory["id"].(int))
		if err != nil {

			logrus.Println(err)
			return
		}
		subcategory["count_products"] = count
	}

	cat["subcategories"] = subcategories

	catMarshalled, err := json.Marshal(cat)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	_, err = w.Write(catMarshalled)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}
}

func (h *CategoryHandler) GetAllCategory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	UserID := r.Context().Value("user_id").(int)

	cats, err := h.cat.GetAll(r.Context(), UserID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	for _, cat := range cats {
		log.Println(cat)
		count, err := h.sc.GetCount(r.Context(), UserID, cat["id"].(int))
		if err != nil {
			log.Println(err)
			count = 0
		}
		cat["count_subcategories"] = count
	}

	catsMarshalled, err := json.Marshal(cats)
	if err != nil {
		return
	}
	_, err = w.Write(catsMarshalled)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		log.Println(err)
		return
	}
}
