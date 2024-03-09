package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/sapphirenw/ai-content-creation-api/src/customer"
)

// Function to define all routes in the api
func addRoutes(
	mux *chi.Mux,
) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		oplog.Info("info here")
		w.Write([]byte("hello world"))
	})
	mux.Route("/customers/{customerId}", customer.CustomerHandler)
}
