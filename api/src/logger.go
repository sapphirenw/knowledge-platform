package main

import (
	"log/slog"

	"github.com/go-chi/httplog/v2"
)

func newLogger() *httplog.Logger {
	logger := httplog.NewLogger("logger", httplog.Options{
		// JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		Tags: map[string]string{
			"version": "v0.1",
			"env":     "dev",
		},
		// QuietDownRoutes: []string{
		// 	"/",
		// 	"/ping",
		// },
		// QuietDownPeriod: 10 * time.Second,
		SourceFieldName: "source",
	})

	return logger
}
