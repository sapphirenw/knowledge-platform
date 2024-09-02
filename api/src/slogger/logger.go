package slogger

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/gorilla/websocket"
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

// Function to write an error message with varying levels of exposition.
// Include `logString` for what you want to log to the server, and
// `visibleString` for what error you want the user to see
func Error(
	ctx context.Context,
	logger *slog.Logger,
	message string,
	err error,
	args ...any,
) error {
	if logger == nil {
		logger = slog.Default()
	}
	logger.ErrorContext(ctx, fmt.Sprintf("%s: %s", message, err), args...)
	return fmt.Errorf("%s: %w", message, err)
}

// Http writer error
func ServerError(
	w http.ResponseWriter,
	logger *slog.Logger,
	statusCode int,
	message string,
	err error,
	args ...any,
) {
	if logger == nil {
		logger = slog.Default()
	}
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", message, err), args...)
	} else {
		logger.Error(message, args...)
	}
	http.Error(w, message, statusCode)
}

// Websocket error
func WsError(
	conn *websocket.Conn,
	logger *slog.Logger,
	message string,
	err error,
	args ...any,
) {
	if logger == nil {
		logger = slog.Default()
	}
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", message, err), args...)
	} else {
		logger.Error(message, args...)
	}
	conn.WriteMessage(websocket.TextMessage, []byte(message))
}
