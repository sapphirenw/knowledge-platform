package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sapphirenw/ai-content-creation-api/src/customer"
	"github.com/sapphirenw/ai-content-creation-api/src/tests"
)

// Function to define all routes in the api
func addRoutes(
	mux *chi.Mux,
) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Successfully hit index"))
	})

	mux.Route("/customers/{customerId}", customer.Handler)
	mux.Route("/tests", tests.Handler)
}
