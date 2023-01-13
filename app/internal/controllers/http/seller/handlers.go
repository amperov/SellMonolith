package seller

import (
	"Selling/app/pkg/auth"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

type Service interface {
	SignUp(ctx context.Context, m map[string]interface{}) (int, error)
	SignIn(ctx context.Context, m map[string]interface{}) (string, error)
	UpdateData(ctx context.Context, m map[string]interface{}, UserID int) error
	GetInfo(ctx context.Context, UserID int) (map[string]interface{}, error)
}

type SellerHandler struct {
	ware auth.MiddleWare
	s    Service
}

func NewSellerHandler(ware auth.MiddleWare, s Service) *SellerHandler {
	return &SellerHandler{ware: ware, s: s}
}

func (s *SellerHandler) Register(r *httprouter.Router) {
	r.POST("/api/auth/sign-in", s.AuthUser)
	r.POST("/api/auth/sign-up", s.CreateUser)
	r.PATCH("/api/seller/update", s.ware.IsAuth(s.UpdateData))
	r.GET("/api/seller/info", s.ware.IsAuth(s.GetInfo))
}

func (s *SellerHandler) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input SignUpInput

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

	log.Println(input.ToMap())
	id, err := s.s.SignUp(r.Context(), input.ToMap())
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success": "user with ID %d created"}`, id)))
	if err != nil {
		log.Println(err)
		return
	}

}

func (s *SellerHandler) AuthUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input SignInInput

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
	token, err := s.s.SignIn(r.Context(), input.ToMap())
	if err != nil {
		log.Println(err)
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"JWT": "%s"}`, token)))
}

func (s *SellerHandler) UpdateData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var upd UpdateInput
	UserID := r.Context().Value("user_id").(int)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &upd)
	if err != nil {
		return
	}

	log.Println(upd.ToMap())
	err = s.s.UpdateData(r.Context(), upd.ToMap(), UserID)
	if err != nil {
		log.Println(err)
		return
	}

}

func (s *SellerHandler) GetInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	UserID := r.Context().Value("user_id").(int)

	info, err := s.s.GetInfo(r.Context(), UserID)
	if err != nil {
		return
	}

	marshal, err := json.Marshal(info)
	if err != nil {
		return
	}
	w.Write(marshal)
}
