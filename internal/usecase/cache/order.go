package cache

import (
	"context"
	"fmt"
	"order-stream/internal/entity"
	"order-stream/internal/usecase"
	"sync"
)

type OrderCache struct {
	locker sync.RWMutex
	store  map[entity.OrderUID]entity.Order
	repo   usecase.OrderRepo
}

func NewOrder(r usecase.OrderRepo) *OrderCache {
	cache := OrderCache{
		store: make(map[entity.OrderUID]entity.Order),
		repo:  r,
	}
	return &cache
}

func NewOrderLoad(r usecase.OrderRepo) (*OrderCache, error) {
	orders, err := r.GetAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("NewOrderLoad- s.repo.GetAll: %w", err)
	}
	store := make(map[entity.OrderUID]entity.Order)
	for _, order := range orders {
		store[order.OrderUID] = order
	}
	cache := OrderCache{
		store: store,
		repo:  r,
	}
	return &cache, nil
}

func (c *OrderCache) Store(ctx context.Context, order entity.Order) error {
	err := c.repo.Store(ctx, order)
	if err != nil {
		return fmt.Errorf("OrderCache - Store - s.repo.Store: %w", err)
	}
	c.locker.Lock()
	c.store[order.OrderUID] = order
	c.locker.Unlock()
	return nil
}

func (c *OrderCache) GetAll(context.Context) ([]entity.Order, error) {
	var orders []entity.Order
	c.locker.RLock()
	for key := range c.store {
		orders = append(orders, c.store[key])
	}
	c.locker.RUnlock()
	return orders, nil
}

func (c *OrderCache) GetById(ctx context.Context, order_uid entity.OrderUID) (entity.Order, error) {
	c.locker.RLock()
	order, ok := c.store[order_uid]
	c.locker.RUnlock()
	if ok {
		return order, nil
	}
	order, err := c.repo.GetById(ctx, order_uid)
	if err != nil {
		return entity.Order{}, fmt.Errorf("OrderCache - GetById - s.repo.GetById: %w", err)
	}
	c.locker.Lock()
	c.store[order_uid] = order
	c.locker.Unlock()
	return order, nil
}
