package logger

import (
	"log/slog"
	"time"

	"github.com/go-chi/httplog/v2"
)

func NewLogger() *httplog.Logger {
	l := slog.New(NewHandler(nil))
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
