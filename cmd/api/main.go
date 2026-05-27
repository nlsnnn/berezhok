package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql"
	"github.com/nlsnnn/berezhok/internal/adapters/rabbitmq"
	"github.com/nlsnnn/berezhok/internal/adapters/redis"
	"github.com/nlsnnn/berezhok/internal/adapters/s3/yandex"
	"github.com/nlsnnn/berezhok/internal/lib/closer"
	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/shared/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	closer.SetLogger(log)

	if cfg.JWTSecret == "" {
		log.Error("JWT_SECRET environment variable is required")
		os.Exit(1)
	}

	log.Info("start app", slog.String("env", cfg.Env))

	// PostgreSQL database
	db, err := postgresql.New(ctx, cfg.Db)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	closer.Add("database", func(_ context.Context) error {
		db.Close()
		return nil
	})

	// S3 Storage
	s3Storage, err := yandex.NewStorage(cfg.S3)
	if err != nil {
		log.Error("failed to initialize S3 storage", "error", err)
		os.Exit(1)
	}

	// Redis
	redisClient, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		log.Error("failed to initialize Redis", "error", err)
		os.Exit(1)
	}
	closer.Add("redis", func(_ context.Context) error {
		return redisClient.Close()
	})

	// RabbitMQ
	rabbitClient, err := rabbitmq.New(cfg.URL, log)
	if err != nil {
		log.Error("failed to initialize RabbitMQ", "error", err)
		os.Exit(1)
	}
	closer.Add("rabbitmq", func(_ context.Context) error {
		return rabbitClient.Close()
	})

	api := application{
		cfg:      cfg,
		pool:     db,
		log:      log,
		s3:       s3Storage,
		redis:    redisClient,
		rabbitmq: rabbitClient,
	}

	serverErrCh := make(chan error, 1)

	go func() {
		if err := api.run(api.mount()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down gracefully")
	case err := <-serverErrCh:
		if err != nil {
			log.Error("failed to start server", sl.Err(err))
			os.Exit(1)
		}
		log.Info("server stopped")
		return
	}

	stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := api.shutdown(shutdownCtx); err != nil {
		log.Error("failed to shutdown server gracefully", sl.Err(err))
	}

	// Закрытие всех ресурсов через глобальный closer
	closerCtx, closerCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer closerCancel()

	if err := closer.CloseAll(closerCtx); err != nil {
		slog.Error("failed to close all resources", "err", err)
	}

	log.Info("app stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
