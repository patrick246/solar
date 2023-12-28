package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/patrick246/solar/statistics/internal/config"
	"github.com/patrick246/solar/statistics/internal/database"
	"github.com/patrick246/solar/statistics/internal/httpserver"
	"github.com/patrick246/solar/statistics/internal/listener"
	"github.com/patrick246/solar/statistics/internal/shelly3em"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fatal: %v", err)

		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))

	db, err := database.Connect(cfg.Database)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}

	err = database.Migrate(context.Background(), db, logger)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	srv := httpserver.New(cfg.MetricsAddr)

	lis, err := listener.NewListener(cfg.Broker, logger)
	if err != nil {
		return fmt.Errorf("listener setup: %v", err)
	}

	repo := shelly3em.NewTimescaleRepository(db)
	shellyHandler := shelly3em.NewHandler(repo)

	lis.Handle(cfg.Broker.Topic, shellyHandler.HandleMessage)

	listenCtx, stopListening := context.WithCancel(context.Background())
	defer stopListening()

	stoppedListening := make(chan struct{})

	go func() {
		defer close(stoppedListening)

		err := lis.Listen(listenCtx)
		if err != nil {
			logger.Error("error listening on mqtt", "error", err)
			return
		}
	}()

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error listening on mqtt", "error", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	logger.Info("shutting down")

	shutdownCh, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = srv.Shutdown(shutdownCh)
	if err != nil {
		return err
	}

	stopListening()
	<-stoppedListening

	return nil
}
