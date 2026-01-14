package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/LeonardoSola/wagon-api/internal/config"
)

var Log *slog.Logger

func Init() error {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	switch config.Config.Env {
	case "development":
		opts.Level = slog.LevelDebug
	case "production":
		opts.Level = slog.LevelInfo
	default:
		return fmt.Errorf("invalid environment: %s", config.Config.Env)
	}

	handler = slog.NewJSONHandler(os.Stdout, opts)

	Log = slog.New(handler)

	slog.SetDefault(Log)

	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	return nil
}
