// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"order-stream/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

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
