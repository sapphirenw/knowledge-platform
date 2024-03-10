package request

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func Decode[T Validator](w http.ResponseWriter, r *http.Request, logger *slog.Logger) (T, bool) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		// send an error request
		logger.ErrorContext(r.Context(), "Error parsing the data", "error", err)
		http.Error(w, "There was an issue parsing the request body", http.StatusBadRequest)
		return v, false
	}
	if problems := v.Valid(r.Context()); len(problems) > 0 {
		// send an error request
		logger.ErrorContext(r.Context(), "The request body was incomplete", "problems", problems)
		http.Error(w, fmt.Sprintf("Your request body was incomplete: %v", problems), http.StatusBadRequest)
		return v, false
	}
	return v, true
}
