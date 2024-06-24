package slogger

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
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
	logger.ErrorContext(ctx, fmt.Sprintf("%s: %s", message, err), args...)
	return fmt.Errorf(message)
}

func ServerError(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	statusCode int,
	message string,
	args ...any,
) {
	logger.Error(message, args...)
	http.Error(w, fmt.Sprintf("There was an error: %s", message), statusCode)
}
