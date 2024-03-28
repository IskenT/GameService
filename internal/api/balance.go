package api

import (
	"github.com/IskenT/MultiGameServices/pkg/cache"
	"github.com/google/uuid"
)

type Balance struct {
	UserId uuid.UUID `json:"user_id,omitempty"`
	Amount int       `json:"amount"`
}

func GetCacheBalance(balance *cache.Balance) Balance {
	return Balance{
		UserId: balance.UserId,
		Amount: balance.Amount,
	}
}
