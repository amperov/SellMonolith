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
	r.GET("/seller/category", h.ware.IsAuth(h.GetAllCategory))
	r.GET("/seller/category/:cat_id", h.ware.IsAuth(h.GetCategory))
	r.POST("/seller/category", h.ware.IsAuth(h.CreateCategory))
	r.PATCH("/seller/category/:cat_id", h.ware.IsAuth(h.UpdateCategory))
	r.DELETE("/seller/category/:cat_id", h.ware.IsAuth(h.DeleteCategory))
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input CreateCategoryInput
	UserID := r.Context().Value("user_id").(int)

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

	input.UserID = UserID

	id, err := h.cat.Create(r.Context(), input.ToMap())

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "category with ID %d created"}`, id)))
	if err != nil {
		log.Println(err)
		return
	}
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var upd UpdateCategoryInput
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &upd)
	if err != nil {
		return
	}

	id, err := h.cat.Update(r.Context(), upd.ToMap(), UserID, CatID)

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "category with ID %d updated"}`, id)))
	if err != nil {
		return
	}

}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
	err = h.cat.Delete(r.Context(), UserID, CatID)
	if err != nil {
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "category with deleted"}`)))
	if err != nil {
		return
	}
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

	cat, err := h.cat.GetOne(r.Context(), UserID, CatID)
	if err != nil {
		return
	}
	subcategories, err := h.sc.GetAll(r.Context(), UserID, CatID)
	if err != nil {
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
		return
	}

	_, err = w.Write(catMarshalled)
	if err != nil {
		return
	}
}

func (h *CategoryHandler) GetAllCategory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := fmt.Sprintf("%v", r.Context().Value("user_id"))
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return
	}

	cats, err := h.cat.GetAll(r.Context(), UserID)
	if err != nil {
		return
	}

	log.Printf("%v", cats)
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
		log.Println(err)
		return
	}
}
