package grpc

import (
	"context"
	"net"

	grpc_service "github.com/IskenT/MultiGameServices/internal/infrastructure/usecase/grpcservice"
	"github.com/IskenT/MultiGameServices/pkg/logger"
	v1 "github.com/IskenT/MultiGameServices/proto/v1"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GrpcServer interface {
	Start()
	Stop(ctx context.Context)
}

type Handler struct {
	srv           *grpc.Server
	serverPort    string
	walletService *grpc_service.WalletService
	logger        logger.Logger
	errc          chan<- error
}

func NewHandler(
	ServerPort string,
	walletService *grpc_service.WalletService,
	logger logger.Logger,
	errc chan<- error,
) *Handler {
	logrus := logrus.NewEntry(logrus.StandardLogger())
	server := &Handler{

		srv: grpc.NewServer(
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				grpc_logrus.StreamServerInterceptor(logrus),
				grpc_recovery.StreamServerInterceptor(),
			)),
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_logrus.UnaryServerInterceptor(logrus),
				grpc_recovery.UnaryServerInterceptor(),
			)),
		),
		walletService: walletService,
		serverPort:    ServerPort,
		logger:        logger,
		errc:          errc,
	}
	return server
}

func (s *Handler) Start() {
	s.logger.Info("running grpc server on ", s.serverPort)
	lis, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		s.logger.ErrorF("grpc server failed: %v", err)
		s.errc <- err
	}

	v1.RegisterWalletServiceServer(s.srv, s.walletService)

	if err := s.srv.Serve(lis); err != nil {
		s.logger.ErrorF("grpc server failed: %v", err)
		s.errc <- err
	}
}

func (s *Handler) Stop(ctx context.Context) {
	s.srv.GracefulStop()
}
