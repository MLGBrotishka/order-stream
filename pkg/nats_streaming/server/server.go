package server

import (
	"fmt"
	"order-stream/pkg/logger"
	nstream "order-stream/pkg/nats_streaming"
	"time"

	"github.com/nats-io/stan.go"
)

const (
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

type MsgHandler func(*stan.Msg) error

type Server struct {
	conn    *nstream.Connection
	notify  chan error
	stop    chan struct{}
	router  map[string]MsgHandler
	logger  logger.Interface
	timeout time.Duration
}

func New(url, clusterId, clientID string, router map[string]MsgHandler, l logger.Interface, opts ...Option) (*Server, error) {
	cfg := nstream.Config{
		URL:       url,
		ClusterID: clusterId,
		ClientID:  clientID,
		WaitTime:  _defaultWaitTime,
		Attempts:  _defaultAttempts,
	}

	c := &Server{
		conn:    nstream.New(cfg),
		notify:  make(chan error),
		stop:    make(chan struct{}),
		router:  router,
		logger:  l,
		timeout: _defaultTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(c)
	}

	err := c.conn.AttemptConnect()
	if err != nil {
		return nil, fmt.Errorf("nstream server - NewServer - c.conn.AttemptConnect: %w", err)
	}

	go c.consumer()

	return c, nil
}

func (s *Server) consumer() {
	for channel := range s.router {
		sub, err := (*s.conn.Connection).Subscribe(channel, func(msg *stan.Msg) {
			err := s.router[channel](msg)
			if err != nil {
				s.logger.Error(err, "nstream server - Server - Router")
			}
		}, stan.DeliverAllAvailable())
		if err != nil {
			s.logger.Error(err, "nstream server - Server - Subscribe")
		}
		defer sub.Unsubscribe()
	}

	<-s.stop
}

func (s *Server) reconnect() {
	close(s.stop)

	err := s.conn.AttemptConnect()
	if err != nil {
		s.notify <- err
		close(s.notify)

		return
	}

	s.stop = make(chan struct{})

	go s.consumer()
}

// Notify
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown
func (s *Server) Shutdown() error {
	select {
	case <-s.notify:
		return nil
	default:
	}

	close(s.stop)
	time.Sleep(s.timeout)

	err := (*s.conn.Connection).Close()
	if err != nil {
		return fmt.Errorf("nstream server - Shutdown - s.Connection.Close: %w", err)
	}

	return nil
}
