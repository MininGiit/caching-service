package main

import (
	"cachingService/internal/infrastructure/cache"
	"cachingService/internal/infrastructure/controller/http"
	"cachingService/internal/usecase"
	"time"
)

func main() {
	cache := cache.New(10, time.Second * 5)
	uc := usecase.New(cache)
	server := http.NewServer(uc)
	server.StartServer()
}