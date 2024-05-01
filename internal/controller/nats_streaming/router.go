package nstream

import (
	"order-stream/internal/usecase"
	"order-stream/pkg/nats_streaming/server"
)

func NewRouter(t usecase.Order) map[string]server.MsgHandler {
	routes := make(map[string]server.MsgHandler)
	{
		newOrderRoutes(routes, t)
	}

	return routes
}
