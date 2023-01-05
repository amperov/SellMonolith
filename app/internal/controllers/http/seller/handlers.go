package seller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"time"
)

type Service interface {
	SignUp(ctx context.Context, m map[string]interface{}) (int, error)
	SignIn(ctx context.Context, m map[string]interface{}) (string, error)
}

type sellerHandler struct {
	s Service
}

func (s *sellerHandler) Register(r *httprouter.Router) {
	r.POST("", s.AuthUser)
	r.POST("", s.CreateUser)
}

func (s *sellerHandler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input SignUpInput

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &input)
	if err != nil {
		return
	}

	id, err := s.s.SignUp(r.Context(), input.ToMap())
	if err != nil {
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("user with ID %d created", id)))
	if err != nil {
		return
	}

}

func (s *sellerHandler) AuthUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input SignInInput

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		return
	}
	token, err := s.s.SignIn(r.Context(), input.ToMap())
	if err != nil {
		return
	}

	cookie := http.Cookie{
		Name:       "JWT",
		Value:      token,
		Expires:    time.Time{}.Add(90 * 24 * 60 * 60 * time.Second),
		RawExpires: "",
		HttpOnly:   true,
	}

	http.SetCookie(w, &cookie)
}
