package nstream

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

// Config
type Config struct {
	URL       string
	ClusterID string
	ClientID  string
	WaitTime  time.Duration
	Attempts  int
}

// Connection
type Connection struct {
	Config
	Connection   *stan.Conn
	Subscription *stan.Subscription
}

// New
func New(cfg Config) *Connection {
	conn := &Connection{
		Config: cfg,
	}
	return conn
}

func (c *Connection) AttemptConnect() error {
	var err error
	for i := c.Attempts; i > 0; i-- {
		if err = c.connect(); err == nil {
			break
		}

		log.Printf("Nats-streaming is trying to connect, attempts left: %d", i)
		time.Sleep(c.WaitTime)
	}

	if err != nil {
		return fmt.Errorf("nats-streaming - AttemptConnect - c.connect: %w", err)
	}

	return nil
}

// Подключение к NATS Streaming
func (c *Connection) connect() error {
	var err error
	// Создание опций подключения
	opts := []stan.Option{
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}),
	}

	// Подключение
	conn, err := stan.Connect(c.ClusterID, c.ClientID, opts...)
	if err != nil {
		return fmt.Errorf("stan.Connect: %w", err)
	}
	c.Connection = &conn
	return nil
}
