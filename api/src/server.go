package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/httprate"
)

func NewServer(
	logger *httplog.Logger,
) http.Handler {
	// create the chi router
	mux := chi.NewRouter()

	// add all middleware
	mux.Use(httplog.RequestLogger(logger))
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RedirectSlashes)
	mux.Use(middleware.ThrottleBacklog(50, 300, time.Second*10)) // adjust
	mux.Use(httprate.LimitByIP(100, 1*time.Minute))

	mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// add the routes
	addRoutes(mux)

	return mux
}
