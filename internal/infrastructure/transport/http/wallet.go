package http

import (
	"context"
	"net/http"

	"github.com/IskenT/MultiGameServices/internal/api"
	"github.com/IskenT/MultiGameServices/internal/infrastructure/usecase/restservice"
	"github.com/IskenT/MultiGameServices/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type WalletService interface {
	HandleDepositAmount(ctx echo.Context) error
	HandleWithdrawAmount(ctx echo.Context) error
	HandleGetBalance(ctx echo.Context) error
}

type WalletHandler struct {
	walletInteractor restservice.WalletService
	logger           logger.Logger
}

func NewWalletService(goodsInteractor restservice.WalletService, logger logger.Logger) WalletService {
	return &WalletHandler{
		walletInteractor: goodsInteractor,
		logger:           logger,
	}
}

func (c *WalletHandler) HandleDepositAmount(ctx echo.Context) error {
	ct, cancel := context.WithCancel(context.Background())
	defer cancel()

	var balance api.Balance

	if err := ctx.Bind(&balance); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid body")
	}

	walletDTO, err := c.walletInteractor.DepositAmount(ct, &balance)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "internal error")
	}

	response := api.GetCacheBalance(walletDTO)
	return ctx.JSON(http.StatusCreated, response)
}

func (c *WalletHandler) HandleWithdrawAmount(ctx echo.Context) error {
	ct, cancel := context.WithCancel(context.Background())
	defer cancel()

	var balance api.Balance

	if err := ctx.Bind(&balance); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid body")
	}

	walletDTO, err := c.walletInteractor.WithdrawAmount(ct, &balance)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "internal error")
	}

	response := api.GetCacheBalance(walletDTO)
	return ctx.JSON(http.StatusCreated, response)
}

func (c *WalletHandler) HandleGetBalance(ctx echo.Context) error {
	ct, cancel := context.WithCancel(context.Background())
	defer cancel()

	userIDStr := ctx.Param("userId")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid userId format")
	}

	walletDTO, err := c.walletInteractor.GetAmount(ct, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "internal error")
	}

	response := api.GetCacheBalance(walletDTO)
	return ctx.JSON(http.StatusCreated, response)
}
