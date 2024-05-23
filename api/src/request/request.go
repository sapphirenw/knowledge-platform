package request

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// Encodes the given type into json and writes it to the request. If the encoding fails, then
// an `http.StatusInternalServerError` is sent instead.
func Encode[T any](w http.ResponseWriter, r *http.Request, logger *slog.Logger, status int, v T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.ErrorContext(r.Context(), "There was an issue encoding the data", "error", err)
		http.Error(w, "There was an issue encoding the body", http.StatusInternalServerError)
	}
}

// Decodes the data as the given type and ensures the data is valid. If the data is not valid
// or there is an issue decoding the request, an `http.StatusBadRequest` is written and (nil, false) is returned
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
