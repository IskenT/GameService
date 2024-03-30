package repository

import (
	"context"
	"errors"

	"github.com/IskenT/MultiGameServices/pkg/cache"
	"github.com/IskenT/MultiGameServices/pkg/logger"
	"github.com/google/uuid"
)

var (
	ErrNotExist = errors.New("good not exist")
)

type WalletRepository interface {
	Debit(ctx context.Context, userId uuid.UUID, amount int) error
	Withdraw(ctx context.Context, userId uuid.UUID, amount int) error
	Get(ctx context.Context, userId uuid.UUID) (*cache.Balance, error)
}

type walletRepository struct {
	cache  cache.Storage
	logger logger.Logger
}

func NewWalletRepository(ctx context.Context, cache cache.Storage, logger logger.Logger) WalletRepository {
	return &walletRepository{
		cache:  cache,
		logger: logger,
	}
}

func (r *walletRepository) Debit(ctx context.Context, userId uuid.UUID, amount int) error {
	r.logger.Info("deposit balance to userId=", userId)
	return r.cache.Deposit(userId, amount)
}

func (r *walletRepository) Withdraw(ctx context.Context, userId uuid.UUID, amount int) error {
	r.logger.Info("withdraw balance from userId=", userId)
	return r.cache.Withdraw(userId, amount)
}

func (r *walletRepository) Get(ctx context.Context, userId uuid.UUID) (*cache.Balance, error) {
	r.logger.Info("get balance by userId=", userId)
	item, ok := r.cache.Get(userId)
	if !ok {
		return nil, ErrNotExist
	} else {
		return item, nil
	}
}
