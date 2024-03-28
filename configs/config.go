package configs

import (
	"os"
	"sync"
)

type Config struct {
	once       sync.Once
	HttpServer struct {
		Port string
	}

	GrpcServer struct {
		Port string
		Host string
	}

	WebSocketServer struct {
		Port string
		Host string
	}
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.once.Do(func() {
		cfg.HttpServer.Port = getEnv("HTTP_PORT", "")

		cfg.GrpcServer.Port = getEnv("GRPC_PORT", "")
		cfg.GrpcServer.Host = getEnv("GRPC_HOST", "")

		cfg.WebSocketServer.Port = getEnv("SOCKET_PORT", "")
		cfg.WebSocketServer.Host = getEnv("SOCKET_HOST", "")
	})
	return cfg, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
