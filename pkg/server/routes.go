package server

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "llamadeck/docs"
	"llamadeck/webui"
)

func SetupRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Add CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   handler.cfg.Server.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   handler.cfg.Server.AllowedHeaders,
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	if handler.cfg.Server.EnableSwagger {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		))
	}

	// OIDC login and whoami must be reachable without a management key or session
	r.Get("/api/v1/auth/whoami", handler.Whoami())
	if handler.oidc != nil {
		r.Route("/api/v1/auth/oidc", func(r chi.Router) {
			r.Get("/login", handler.OIDCLogin())
			r.Get("/callback", handler.OIDCCallback())
			r.Post("/logout", handler.OIDCLogout())
		})
	}

	// Define routes
	r.Route("/api/v1", func(r chi.Router) {

		if handler.authMiddleware != nil && handler.cfg.Auth.RequireManagementAuth {
			r.Use(handler.authMiddleware.ManagementAuthMiddleware())
