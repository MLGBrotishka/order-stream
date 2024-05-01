package usecase

import (
	"context"
	"fmt"

	"order-stream/internal/entity"
)

// Order
type OrderUseCase struct {
	repo OrderRepo
}

// New
func NewOrder(r OrderRepo) *OrderUseCase {
	return &OrderUseCase{
		repo: r,
	}
}

// GetOrder - getting order from repo
func (uc *OrderUseCase) GetOrder(ctx context.Context, id entity.OrderUID) (entity.Order, error) {
	order, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return entity.Order{}, fmt.Errorf("OrderUseCase - GetOrder - s.repo.GetById: %w", err)
	}
	return order, nil
}

// SaveOrder - saving to repo
func (uc *OrderUseCase) SaveOrder(ctx context.Context, order entity.Order) error {
	err := uc.repo.Store(ctx, order)
	if err != nil {
		return fmt.Errorf("OrderUseCase - SaveOrder - s.repo.Store: %w", err)
	}
	return nil
}
