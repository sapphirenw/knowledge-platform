package slogger

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

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
