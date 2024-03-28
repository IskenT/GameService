package restservice

import (
	"context"
	"fmt"

	"github.com/IskenT/MultiGameServices/internal/api"
	repository "github.com/IskenT/MultiGameServices/internal/infrastructure/repository/restrepo"
	"github.com/IskenT/MultiGameServices/pkg/cache"
	"github.com/IskenT/MultiGameServices/pkg/logger"
	"github.com/google/uuid"
)

type WalletService interface {
	DepositAmount(ctx context.Context, balance *api.Balance) (*cache.Balance, error)
	WithdrawAmount(ctx context.Context, balance *api.Balance) (*cache.Balance, error)
	GetAmount(ctx context.Context, userId uuid.UUID) (*cache.Balance, error)
}

type walletService struct {
	walletRepository repository.WalletRepository
	logger           logger.Logger
}

func NewWalletService(ctx context.Context, walletRepository repository.WalletRepository, logger logger.Logger) WalletService {
	return &walletService{
		walletRepository: walletRepository,
		logger:           logger,
	}
}

func (i *walletService) DepositAmount(ctx context.Context, good *api.Balance) (*cache.Balance, error) {
	err := i.walletRepository.Debit(ctx, good.UserId, good.Amount)
	if err != nil {
		return nil, fmt.Errorf("error on deposit: %w", err)
	}

	return i.walletRepository.Get(ctx, good.UserId)
}

func (i *walletService) WithdrawAmount(ctx context.Context, good *api.Balance) (*cache.Balance, error) {
	err := i.walletRepository.Withdraw(ctx, good.UserId, good.Amount)
	if err != nil {
		return nil, fmt.Errorf("error on deposit: %w", err)
	}

	return i.walletRepository.Get(ctx, good.UserId)
}

func (i *walletService) GetAmount(ctx context.Context, userId uuid.UUID) (*cache.Balance, error) {
	return i.walletRepository.Get(ctx, userId)
}
