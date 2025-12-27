package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/LeonardoSola/wagon-api/internal/config"
)

var Log *slog.Logger

// Init inicializa o logger global da aplicação
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

	// Define como logger global do slog
	slog.SetDefault(Log)

	// Redireciona log.Printf (goose, libs antigas) para o mesmo output
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	return nil
}
