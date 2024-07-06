package request

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

type WsMessage[T any] struct {
	HasError bool `json:"hasError"`
	Content  T    `json:"content"`
}

/*
Writes a message to a websocket. If the passed error is not nil, then the user will be made
aware that there was an error. In either case, a JSON blob will be returned to the user, with
the implementation choosing how to process the passed `message` type.

Currently, the web app assumes that `message` is a string if the err is not nil, and that T is a
valid JSON object in all other cases.
*/
func WriteWs[T any](
	ctx context.Context,
	logger *slog.Logger,
	conn *websocket.Conn,
	message T,
	err error,
	args ...any,
) error {
	logger.Debug("writing to the websocket", "message", message)
	// log the error if applicable
	if err != nil {
		logger.Error(fmt.Sprintf("%v: %s", message, err), args...)
	}

	// create the object
	msg := WsMessage[T]{
		HasError: err != nil,
		Content:  message,
	}
	enc, err := json.Marshal(msg)
	if err != nil {
		return slogger.Error(ctx, logger, "failed to encode the message", err)
	}

	// write to the user
	if err := conn.WriteJSON(enc); err != nil {
		return slogger.Error(ctx, logger, "failed to write the message", err)
	}

	return nil
}
