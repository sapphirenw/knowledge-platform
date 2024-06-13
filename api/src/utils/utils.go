package utils

import (
	"log/slog"

	"github.com/sapphirenw/ai-content-creation-api/src/logger"
)

func DefaultLogger() *slog.Logger {
	return slog.New(logger.NewHandler(&slog.HandlerOptions{Level: slog.LevelInfo}))
}
