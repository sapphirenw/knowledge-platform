package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sapphirenw/ai-content-creation-api/src/beta"
	"github.com/sapphirenw/ai-content-creation-api/src/customer"
)

// Function to define all routes in the api
func addRoutes(
	mux *chi.Mux,
) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Successfully hit index"))
	})

	mux.Route("/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Successfully hit index"))
		})

		r.Route("/beta", beta.Handler)
		r.Route("/customers/{customerId}", customer.Handler)
	})
}
