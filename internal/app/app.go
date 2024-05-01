// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"order-stream/config"
	v1 "order-stream/internal/controller/http/v1"
	nstream "order-stream/internal/controller/nats_streaming"
	"order-stream/internal/usecase"
	"order-stream/internal/usecase/cache"
	"order-stream/internal/usecase/repo"
	"order-stream/pkg/httpserver"
	"order-stream/pkg/logger"
	nstreamserver "order-stream/pkg/nats_streaming/server"
	"order-stream/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Cache
	cache, err := cache.NewOrderLoad(
		repo.NewOrder(pg),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - cache.NewOrderLoad: %w", err))
	}

	// Use case
	orderUseCase := usecase.NewOrder(
		cache,
	)

	// Nats-streaming client
	nstreamRouter := nstream.NewRouter(orderUseCase)
	nstreamServer, err := nstreamserver.New(cfg.NS.URL, cfg.NS.ClusterID, cfg.NS.ClientID, nstreamRouter, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - nstreamServer - server.New: %w", err))
	}
	// HTTP Server
	handler := gin.New()
	v1.NewOrderRouter(handler, l, orderUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-nstreamServer.Notify():
		l.Error(fmt.Errorf("app - Run - nstreamServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = nstreamServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - nstreamServer.Shutdown: %w", err))
	}
}
