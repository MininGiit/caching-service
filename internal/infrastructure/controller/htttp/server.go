package http

import (
	"cachingService/internal/usecase"
	"net/http"

)

type Server struct {
	httpServer   *http.Server
}

func NewServer(uc *usecase.UseCase) *Server{
	handler := NewHandler(uc)
	router := handler.InitRouter()
	return &Server{
		httpServer: &http.Server{
			Handler:      router,
			Addr:         ":8080",
		},
	}
}

func (s *Server) StartServer() {
	s.httpServer.ListenAndServe()
}
