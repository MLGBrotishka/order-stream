package nstream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"order-stream/internal/entity"
	"order-stream/internal/usecase"
	"order-stream/pkg/nats_streaming/server"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
)

type orderRoutes struct {
	orderUseCase usecase.Order
}

func newOrderRoutes(routes map[string]server.MsgHandler, t usecase.Order) {
	r := &orderRoutes{t}
	{
		routes["orders"] = r.getOrders()
	}
}

func (r *orderRoutes) getOrders() server.MsgHandler {
	return func(msg *stan.Msg) error {
		order := entity.Order{}
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&order)
		if err != nil {
			return fmt.Errorf("nstream - newOrderRoutes - getOrders - msg.Unmarshal: %w", err)
		}
		validate := validator.New()
		err = validate.Struct(order)
		if err != nil {
			return fmt.Errorf("nstream - newOrderRoutes - getOrders - validation error: %w", err)
		}
		err = r.orderUseCase.SaveOrder(context.Background(), order)
		if err != nil {
			return fmt.Errorf("nstream - newOrderRoutes - getOrders -  r.orderUseCase.SaveOrder: %w", err)
		}
		return nil
	}
}
