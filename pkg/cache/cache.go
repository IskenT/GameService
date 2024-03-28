package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type (
	Balance struct {
		UserId uuid.UUID
		Amount int
	}

	Cache struct {
		Items map[uuid.UUID]CacheItem
		mu    sync.RWMutex
		ttl   time.Duration
	}

	CacheItem struct {
		Balance   Balance
		ExpiresAt time.Time
	}
)

func (c *Cache) Get(key uuid.UUID) (*Balance, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.Items[key]
	if item.IsExpired() {
		return &Balance{}, false
	} else {
		return &item.Balance, ok
	}
}

func (c *Cache) Set(key uuid.UUID, balance Balance) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	item := CacheItem{
		Balance:   balance,
		ExpiresAt: time.Now().Add(c.ttl),
	}
	c.Items[key] = item
	return nil
}

func (c *Cache) Deposit(key uuid.UUID, amount int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.Items[key]
	if !exists {
		balance := Balance{
			UserId: key,
			Amount: amount,
		}
		entry = CacheItem{
			Balance:   balance,
			ExpiresAt: time.Now().Add(c.ttl),
		}
	} else {
		entry.Balance.Amount += amount
		entry.ExpiresAt = time.Now().Add(c.ttl)
	}

	c.Items[key] = entry

	return nil
}

func (c *Cache) Withdraw(key uuid.UUID, amount int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.Items[key]
	if !exists {
		return fmt.Errorf("user not found")
	}

	if entry.Balance.Amount < amount {
		return fmt.Errorf("user not found")
	}

	entry.Balance.Amount -= amount
	entry.ExpiresAt = time.Now().Add(c.ttl)
	c.Items[key] = entry

	return nil
}

func NewCache(ttl time.Duration) (*Cache, error) {
	c := &Cache{
		Items: make(map[uuid.UUID]CacheItem),
		mu:    sync.RWMutex{},
		ttl:   ttl,
	}
	go c.CleanUp()
	return c, nil
}

func (c *Cache) CleanUp() {
	for range time.Tick(c.ttl) {
		c.mu.RLock()
		tempSlice := make([]uuid.UUID, 0, len(c.Items))
		for key, item := range c.Items {
			if item.IsExpired() {
				tempSlice = append(tempSlice, key)
			}
		}
		c.mu.RUnlock()
		c.mu.Lock()
		for _, value := range tempSlice {
			delete(c.Items, value)
		}
		c.mu.Unlock()
	}
}

func (i *CacheItem) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}
