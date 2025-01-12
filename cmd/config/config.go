package config

import (
 	"flag"
 	"time"
 	"github.com/caarlos0/env/v6"
)

type Config struct {
 	Server ServerConfig
 	Cache  CacheConfig
 	Log    LogConfig
}

type ServerConfig struct {
 	PortHost string `env:"SERVER_HOST_PORT"`
}

type CacheConfig struct {
 	MaxSize     int           `env:"CACHE_SIZE"`
 	DefaultTtl  time.Duration `env:"DEFAULT_CACHE_TTL"`
}

type LogConfig struct {
 	Level string `env:"LOG_LEVEL"`
}

func Init() *Config {
    cfg := Config{
        Server: ServerConfig{PortHost: "localhost:8080"},
        Cache: CacheConfig{
            MaxSize:     10,
            DefaultTtl:  time.Minute,
        },
        Log: LogConfig{Level: "WARN"},
    }
    env.Parse(&cfg)
    var portHost string
    var maxSize int
    var defaultTtlInt int 
    var logLevel string

    flag.StringVar(&portHost, "server-host-port", "", "Server host and port")
    flag.IntVar(&maxSize, "cache-size", 0, "Maximum size of cache")
    flag.IntVar(&defaultTtlInt, "default-cache-ttl", 0, "Default TTL for cache")
    flag.StringVar(&logLevel, "log-level", "", "Log level (debug, info, warning, error)")

    flag.Parse()
    if portHost != "" {
        cfg.Server.PortHost = portHost
    }
    if maxSize != 0 {
        cfg.Cache.MaxSize = maxSize
    }
    if defaultTtlInt != 0 {
  		defaultTtl := time.Second * time.Duration(defaultTtlInt)
        cfg.Cache.DefaultTtl = defaultTtl
    }
    if logLevel != "" {
        cfg.Log.Level = logLevel
    }
 	return &cfg
}