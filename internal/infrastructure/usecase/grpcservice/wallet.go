package grpcservice

import (
	"context"
	"fmt"

	grpcrepository "github.com/IskenT/MultiGameServices/internal/infrastructure/repository/grpcrepo"
	"github.com/IskenT/MultiGameServices/pkg/logger"
	v1 "github.com/IskenT/MultiGameServices/proto/v1"
	"github.com/google/uuid"
)

type WalletService struct {
	v1.WalletServiceServer
	walletRepository grpcrepository.WalletRepository
	logger           logger.Logger
}

func NewWalletService(ctx context.Context, walletRepository grpcrepository.WalletRepository, logger logger.Logger) *WalletService {
	return &WalletService{
		walletRepository: walletRepository,
		logger:           logger,
	}
}

func (h *WalletService) GetBalanceByUserId(ctx context.Context, req *v1.GetWalletBalanceRequest) (*v1.GetWalletBalanceResponse, error) {
	id, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse userId %w", err)
	}
	res, err := h.walletRepository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance %w", err)
	}
	return &v1.GetWalletBalanceResponse{Balance: int32(res.Amount)}, nil
}
