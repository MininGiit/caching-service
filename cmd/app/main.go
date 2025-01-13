package main

import (
	"cachingService/cmd/config"
	"cachingService/internal/infrastructure/cache"
	"cachingService/internal/infrastructure/controller/http"
	"cachingService/internal/infrastructure/logger"
	"cachingService/internal/usecase"
	"context"
	"fmt"
)

func main() {
	cfg := config.Init()
	fmt.Println(cfg)
	cache := cache.New(cfg.Cache.MaxSize, cfg.Cache.DefaultTtl)
	uc := usecase.New(cache)
	logger := logger.New(cfg.Log.Level)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server := http.NewServer(ctx, cfg.Server.PortHost, uc, logger)
	server.StartServer()
}