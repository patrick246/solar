package listener

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"

	"github.com/patrick246/solar/statistics/internal/config"
)

type Listener struct {
	cfg    config.Broker
	logger *slog.Logger

	url      *url.URL
	handlers map[string]paho.MessageHandler
}

func NewListener(cfg config.Broker, logger *slog.Logger) (*Listener, error) {
	parsedURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	return &Listener{
		cfg:      cfg,
		logger:   logger,
		url:      parsedURL,
		handlers: map[string]paho.MessageHandler{},
	}, nil
}

func (l *Listener) Handle(topic string, handler paho.MessageHandler) {
	l.handlers[topic] = handler
}

func (l *Listener) Listen(ctx context.Context) error {
	router := paho.NewStandardRouter()

	for topic, handler := range l.handlers {
		router.RegisterHandler(topic, handler)
	}

	clientConfig := autopaho.ClientConfig{
		BrokerUrls:     []*url.URL{l.url},
		Debug:          NewLogger(l.logger, slog.LevelDebug),
		PahoDebug:      NewLogger(l.logger, slog.LevelDebug),
		PahoErrors:     NewLogger(l.logger, slog.LevelError),
		KeepAlive:      20,
		OnConnectionUp: l.onConnect,
		ClientConfig: paho.ClientConfig{
			Router: router,
		},
	}

	if l.url.User != nil {
		if password, ok := l.url.User.Password(); ok {
			clientConfig.SetUsernamePassword(l.url.User.Username(), []byte(password))
		}
	}

	conn, err := autopaho.NewConnection(ctx, clientConfig)
	if err != nil {
		return err
	}

	connectionCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	err = conn.AwaitConnection(connectionCtx)
	if err != nil {
		return err
	}

	<-ctx.Done()

	_, err = conn.Unsubscribe(ctx, &paho.Unsubscribe{
		Topics: []string{l.cfg.Topic},
	})
	if err != nil {
		return err
	}

	err = conn.Disconnect(ctx)
	if err != nil {
		return err
	}

	select {
	case <-conn.Done():
	case <-time.After(10 * time.Second):
	}

	return nil
}

func (l *Listener) onConnect(connectionManager *autopaho.ConnectionManager, _ *paho.Connack) {
	_, err := connectionManager.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{{
			Topic: l.cfg.Topic,
		}},
	})
	if err != nil {
		l.logger.Error("subscribe error: %v", err)
	}
}
