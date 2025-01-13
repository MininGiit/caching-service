package main

import (
	"cachingService/cmd/config"
	"cachingService/internal/infrastructure/cache"
	"cachingService/internal/infrastructure/controller/http"
	"cachingService/internal/infrastructure/logger"
	"cachingService/internal/usecase"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Init()
	cache := cache.New(cfg.Cache.MaxSize, cfg.Cache.DefaultTtl)
	cache.StartCollector()
	defer cache.StartCollector()
	uc := usecase.New(cache)
	logger := logger.New(cfg.Log.Level)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server := http.NewServer(ctx, cfg.Server.PortHost, uc, logger)
	go server.StartServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	logger.Info("Received signal:", "sig", sig)
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Shutdown:", "err", err)
	}
	logger.Info("Gracefull Shutdown")
}
