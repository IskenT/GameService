package http

import (
	"context"
	"fmt"

	"github.com/IskenT/MultiGameServices/pkg/logger"

	"github.com/labstack/echo/v4"
)

type HTTPServer interface {
	Start()
	Stop(ctx context.Context)
}

type Handler struct {
	echo          *echo.Echo
	serverPort    string
	walletService WalletService
	logger        logger.Logger
	errc          chan<- error
}

func NewHandler(
	ServerPort string,
	walletService WalletService,
	logger logger.Logger,
	errc chan<- error,
) *Handler {
	server := &Handler{
		echo:          echo.New(),
		walletService: walletService,
		serverPort:    ServerPort,
		logger:        logger,
		errc:          errc,
	}

	return server
}

func (s *Handler) Start() {
	s.echo.POST("/wallet/deposit", s.handleDepositAmount)
	s.echo.POST("/wallet/withdraw", s.handleWithdrawAmount)
	s.echo.GET("/wallet/balance/:userId", s.handleGetBalance)

	s.logger.Info("running echo server on ", s.serverPort)
	port := fmt.Sprintf(":%v", s.serverPort)
	if err := s.echo.Start(port); err != nil {
		s.logger.ErrorF("echo server failed: %v", err)
		s.errc <- fmt.Errorf("echo server failed: %v", err)
	}
}

func (s *Handler) Stop(ctx context.Context) {
	err := s.echo.Shutdown(ctx)
	if err != nil {
		s.logger.ErrorF("echo server failed: %v", err)
		s.errc <- err
	}
}

func (s *Handler) handleDepositAmount(ctx echo.Context) error {
	return s.walletService.HandleDepositAmount(ctx)
}

func (s *Handler) handleWithdrawAmount(ctx echo.Context) error {
	return s.walletService.HandleWithdrawAmount(ctx)
}

func (s *Handler) handleGetBalance(ctx echo.Context) error {
	return s.walletService.HandleGetBalance(ctx)
}
