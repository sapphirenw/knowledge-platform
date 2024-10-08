package beta

import (
	"github.com/go-chi/chi/v5"
)

func Handler(mux chi.Router) {
	mux.Post("/createBetaApiKey", createBetaApiKey)
	mux.Route("/customers", func(r chi.Router) {
		r.Get("/get", getCustomer)
	})
}
