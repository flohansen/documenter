package app

import "context"

type Server struct {
	Repository DocumentationRepository
}

func NewServer(repo DocumentationRepository) *Server {
	return &Server{
		Repository: repo,
	}
}

func (s *Server) Run(ctx context.Context) error {
	return nil
}
