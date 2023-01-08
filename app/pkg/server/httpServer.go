package server

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(router *httprouter.Router) error {
	s.httpServer = &http.Server{
		Addr:              ":8080",
		Handler:           router,
		TLSConfig:         nil,
		ReadTimeout:       20 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       20 * time.Second,
		MaxHeaderBytes:    20,
	}

	return s.httpServer.ListenAndServe()
}
