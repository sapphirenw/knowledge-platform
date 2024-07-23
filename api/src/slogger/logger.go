package slogger

import (
	"log/slog"
	"os"
	"time"

	"github.com/go-chi/httplog/v2"
)

func NewLogger() *httplog.Logger {
	// l := slog.New(NewHandler(nil))
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	logger := httplog.Logger{
		Logger: l,
		Options: httplog.Options{
			JSON:             true,
			LogLevel:         slog.LevelDebug,
			Concise:          true,
			RequestHeaders:   false,
			ResponseHeaders:  false,
			MessageFieldName: "message",
			TimeFieldFormat:  time.RFC850,
			Tags: map[string]string{
				"version": "v0.1",
				"env":     "dev",
			},
			QuietDownRoutes: []string{
				"/",
				"/ping",
			},
			QuietDownPeriod: 10 * time.Second,
			SourceFieldName: "source",
		},
	}

	return &logger
}
