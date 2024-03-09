package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

func NewServer(
	logger *httplog.Logger,
) http.Handler {
	// create the chi router
	mux := chi.NewRouter()

	// add all middleware
	mux.Use(httplog.RequestLogger(logger))

	// add the routes
	addRoutes(mux)

	return mux
}
