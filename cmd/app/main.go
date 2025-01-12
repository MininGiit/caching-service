package main

import (
	"cachingService/cmd/config"
	"cachingService/internal/infrastructure/cache"
	"cachingService/internal/infrastructure/controller/http"
	"cachingService/internal/usecase"
	"fmt"
)

func main() {
	cfg := config.Init()
	fmt.Println(cfg)
	cache := cache.New(cfg.Cache.MaxSize, cfg.Cache.DefaultTtl)
	uc := usecase.New(cache)
	server := http.NewServer(cfg.Server.PortHost, uc)
	server.StartServer()
}