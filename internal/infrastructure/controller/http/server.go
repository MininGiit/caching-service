/*
Пакет http содержит реализацию http сервера для взамиодейсвтия с сервисом по протоколу http
*/
package http

import (
	"cachingService/internal/logger"
	"cachingService/internal/usecase"
	"context"
	"net/http"
)

// Server структура http сервера
type Server struct {
	httpServer *http.Server
	logger     logger.Logger
}

// NewServer создание нового экземпляра сервера
func NewServer(ctx context.Context, hostPort string, uc usecase.IUseCase, logger logger.Logger) *Server {
	handler := newHandler(ctx, uc, logger)
	router := handler.initRouter()
	return &Server{
		httpServer: &http.Server{
			Handler: router,
			Addr:    hostPort,
		},
		logger: logger,
	}
}

// StartServer запуск сервера
func (s *Server) StartServer() {
	s.logger.Info("Start server", "addr", s.httpServer.Addr)
	s.httpServer.ListenAndServe()
}

// Shutdown корректное завершение работы сервера
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
