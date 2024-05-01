package repo

import (
	"context"
	"fmt"
	"order-stream/internal/entity"
	"order-stream/pkg/postgres"

	sq "github.com/Masterminds/squirrel"
)

// OrderRepo
type OrderRepo struct {
	*postgres.Postgres
}

// NewOrder
func NewOrder(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{pg}
}

// Store
func (r *OrderRepo) Store(ctx context.Context, order entity.Order) error {
	sql, args, err := r.Builder.
		Insert("orders").
		Columns("order_uid, order_data").
		Values(order.OrderUID, order).
		ToSql()
	if err != nil {
		return fmt.Errorf("OrderRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OrderRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

// GetAll
func (r *OrderRepo) GetAll(ctx context.Context) ([]entity.Order, error) {
	sql, _, err := r.Builder.
		Select("order_data").
		From("orders").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - GetAll - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - GetAll - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var entities []entity.Order

	for rows.Next() {
		e := entity.Order{}

		err = rows.Scan(&e)
		if err != nil {
			return nil, fmt.Errorf("OrderRepo - GetAll - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

// GetById
func (r *OrderRepo) GetById(ctx context.Context, orderUID entity.OrderUID) (entity.Order, error) {
	sql, _, err := r.Builder.
		Select("order_data").
		From("orders").
		Where(sq.Eq{"order_uid": orderUID}).
		ToSql()
	if err != nil {
		return entity.Order{}, fmt.Errorf("OrderRepo - GetById - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql)
	e := entity.Order{}
	err = row.Scan(&e)
	if err != nil {
		return entity.Order{}, fmt.Errorf("OrderRepo - GetById - row.Scan: %w", err)
	}

	return e, nil
}
