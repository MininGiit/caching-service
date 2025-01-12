package http

import (
	"cachingService/internal/usecase"
	"net/http"

)

type Server struct {
	httpServer   *http.Server
}

func NewServer(hostPort string, uc *usecase.UseCase) *Server{
	handler := NewHandler(uc)
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
