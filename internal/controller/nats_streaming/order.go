package nstream

import (
	"context"
	"encoding/json"
	"fmt"
	"order-stream/internal/entity"
	"order-stream/internal/usecase"
	"order-stream/pkg/nats_streaming/server"

	"github.com/nats-io/stan.go"
)

type orderRoutes struct {
	orderUseCase usecase.Order
}

func newOrderRoutes(routes map[string]server.MsgHandler, t usecase.Order) {
	r := &orderRoutes{t}
	{
		routes["get-orders"] = r.getOrders()
	}
}

func (r *orderRoutes) getOrders() server.MsgHandler {
	return func(msg *stan.Msg) error {
		order := entity.Order{}
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			return fmt.Errorf("nstream - newOrderRoutes - getOrders - msg.Unmarshal: %w", err)
		}
		err = r.orderUseCase.SaveOrder(context.Background(), order)
		if err != nil {
			return fmt.Errorf("nstream - newOrderRoutes - getOrders -  r.orderUseCase.SaveOrder: %w", err)
		}
		return nil
	}
}
