package http

import (
	"cachingService/internal/usecase"
	"context"
	"net/http"
)

type Server struct {
	httpServer   *http.Server
}

func NewServer(ctx context.Context, hostPort string, uc usecase.IUseCase) *Server{
	handler := NewHandler(ctx, uc)
	router := handler.InitRouter()
	return &Server{
		httpServer: &http.Server{
			Handler:      router,
			Addr:         hostPort,
		},
	}
}

func (s *Server) StartServer() {
	s.httpServer.ListenAndServe()
}
