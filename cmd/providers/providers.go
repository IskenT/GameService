package providers

import (
	"os"
	"time"

	"github.com/IskenT/MultiGameServices/configs"
	"github.com/IskenT/MultiGameServices/internal/infrastructure/transport/grpc"
	"github.com/IskenT/MultiGameServices/internal/infrastructure/transport/http"
	http_service "github.com/IskenT/MultiGameServices/internal/infrastructure/transport/http"
	grpc_service "github.com/IskenT/MultiGameServices/internal/infrastructure/usecase/grpcservice"

	"github.com/IskenT/MultiGameServices/pkg/cache"
	"github.com/IskenT/MultiGameServices/pkg/logger"
	"github.com/IskenT/MultiGameServices/pkg/logger/zerolog"
)

func ProvideHTTPServer(config *configs.Config, walletService http_service.WalletService, logger logger.Logger, errc chan<- error) http.HTTPServer {
	return http.NewHandler(config.HttpServer.Port, walletService, logger, errc)
}

func ProvideGRPCServer(config *configs.Config, walletService *grpc_service.WalletService, logger logger.Logger, errc chan<- error) grpc.GrpcServer {
	return grpc.NewHandler(config.GrpcServer.Port, walletService, logger, errc)
}

func ProvideConsoleLogger(cnf *configs.Config) (logger.Logger, error) {
	return zerolog.NewZeroLog(os.Stdout)
}

func ProvideCache() (cache.Storage, error) {
	return cache.NewCache(60 * time.Second)
}
