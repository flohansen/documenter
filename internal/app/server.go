package app

import (
	"context"
	"net/http"
)

type Server struct {
	Handler http.Handler
}

func NewServer(handler http.Handler) *Server {
	return &Server{
		Handler: handler,
	}
}

func (s *Server) Run(ctx context.Context) error {
	return http.ListenAndServe(":3000", s.Handler)
}
