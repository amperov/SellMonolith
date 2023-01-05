package auth

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type MiddleWare struct {
	tm TokenManager
}

func (w *MiddleWare) IsAuth(handle httprouter.Handle) httprouter.Handle {

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		header := request.Header.Get("Authorization")
		headerArray := strings.Split(header, " ")
		if len(headerArray) != 2 {

		}

		id, err := w.tm.ValidateToken(headerArray[1])
		if err != nil {
			return
		}
		ctx := context.WithValue(request.Context(), "user_id", id)
		request.WithContext(ctx)

		handle(writer, request, params)
	}
}
