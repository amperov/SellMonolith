package auth

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type MiddleWare struct {
	tm TokenManager
}

func NewMiddleWare(tm TokenManager) MiddleWare {
	return MiddleWare{tm: tm}
}

func (w *MiddleWare) IsAuth(handle httprouter.Handle) httprouter.Handle {

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		header := request.Header.Get("Authorization")
		headerArray := strings.Split(header, " ")
		if len(headerArray) != 2 {

		}
		logrus.Println(headerArray[1])
		id, err := w.tm.ValidateToken(headerArray[1])
		if err != nil {
			logrus.Println(err)
			return
		}
		ctx := context.WithValue(request.Context(), "userID", id)
		request.WithContext(ctx)

		handle(writer, request.WithContext(ctx), params)
	}
}
