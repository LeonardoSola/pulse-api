package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LeonardoSola/pulse-api/adapters/bootstrap"
	"github.com/LeonardoSola/pulse-api/adapters/cmd/http"
	"github.com/LeonardoSola/pulse-api/adapters/cmd/whatsapp"
	"github.com/LeonardoSola/pulse-api/adapters/database/postgres"
	"github.com/LeonardoSola/pulse-api/adapters/messaging/rabbitmq"
	"github.com/LeonardoSola/pulse-api/internal/config"
	"github.com/LeonardoSola/pulse-api/internal/logger"
)

const shutdownTimeout = 10 * time.Second

func main() {
	if err := run(); err != nil {
		if logger.Log != nil {
			logger.Log.Error("fatal", "error", err)
		} else {
			fmt.Fprintln(os.Stderr, "fatal:", err)
		}
		os.Exit(1)
	}
}

func run() error {
	if err := config.Load(); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	if err := logger.Init(); err != nil {
		return fmt.Errorf("logger: %w", err)
	}
	logger.Log.Info("✅ Logger")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := postgres.NewConnection(); err != nil {
		return fmt.Errorf("database: %w", err)
	}
	defer postgres.Close()
	logger.Log.Info("✅ Database")

	// whatsmeow creates its own tables first; our migrations reference them.
	if err := whatsapp.InitContainer(); err != nil {
		return fmt.Errorf("whatsapp: %w", err)
	}
	defer whatsapp.DisconnectAll()
	logger.Log.Info("✅ WhatsApp")

	current, err := postgres.InitMigrations()
	if err != nil {
		return fmt.Errorf("migrations: %w", err)
	}
	logger.Log.Info("✅ Migrations", "current", current)

	if err := rabbitmq.InitRabbitMQ(); err != nil {
		return fmt.Errorf("messaging: %w", err)
	}
	defer rabbitmq.Close()
	logger.Log.Info("✅ Messaging", "enabled", config.Config.RabbitMQ.Enabled)

	bootstrap.Init()

	if err := bootstrap.Services.DeviceService.StartDevices(ctx); err != nil {
		return fmt.Errorf("devices: %w", err)
	}
	logger.Log.Info("✅ Devices")

	server := &stdhttp.Server{
		Addr:              fmt.Sprintf(":%d", config.Config.Server.Port),
		Handler:           http.Init(),
		ReadHeaderTimeout: 10 * time.Second,
		// No WriteTimeout: GET /devices/{id}/qr is a long-lived SSE stream.
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Log.Info("✅ HTTP", "port", config.Config.Server.Port)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, stdhttp.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}
	case <-ctx.Done():
		slog.Info("shutting down")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// Open SSE streams may outlive the timeout; force-close them.
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Warn("graceful shutdown timed out, closing connections", "error", err)
			_ = server.Close()
		}
	}

	return nil
}
