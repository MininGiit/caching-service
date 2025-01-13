package http

import (
	"cachingService/internal/logger"
	"cachingService/internal/usecase"
	"context"
	"net/http"
)

type Server struct {
	httpServer   *http.Server
	logger		logger.Logger
}

func NewServer(ctx context.Context, hostPort string, uc usecase.IUseCase, logger logger.Logger) *Server{
	handler := NewHandler(ctx, uc, logger)
	router := handler.InitRouter()
	return &Server{
		httpServer: &http.Server{
			Handler:      router,
			Addr:         hostPort,
		},
		logger: logger,
	}
}

func (s *Server) StartServer() {
	s.logger.Info("Start server", "addr", s.httpServer.Addr)
	s.httpServer.ListenAndServe()
}

func (s *Server)Shutdown(ctx context.Context) error{
	return s.httpServer.Shutdown(ctx)
}
