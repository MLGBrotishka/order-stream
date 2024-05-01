package usecase

import (
	"context"

	"order-stream/internal/entity"
)

type (
	Order interface {
		SaveOrder(context.Context, entity.Order) error
		GetOrder(context.Context, entity.OrderUID) (entity.Order, error)
	}

	OrderRepo interface {
		Store(context.Context, entity.Order) error
		GetAll(context.Context) ([]entity.Order, error)
		GetById(context.Context, entity.OrderUID) (entity.Order, error)
	}
)
