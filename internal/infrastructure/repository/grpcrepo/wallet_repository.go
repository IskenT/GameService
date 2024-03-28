package grpcrepository

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
	Get(ctx context.Context, userId uuid.UUID) (*cache.Balance, error)
}

type walletRepository struct {
	cache  *cache.Cache
	logger logger.Logger
}

func NewWalletRepository(ctx context.Context, cache *cache.Cache, logger logger.Logger) WalletRepository {
	return &walletRepository{
		cache:  cache,
		logger: logger,
	}
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
