package utils

import (
	"log/slog"

	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

func DefaultLogger() *slog.Logger {
	return slog.New(slogger.NewHandler(&slog.HandlerOptions{Level: slog.LevelInfo}))
}
