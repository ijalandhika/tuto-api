package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ijalandhika/tuto-api/internal/auth"
	authdb "github.com/ijalandhika/tuto-api/internal/auth/db"
	"github.com/ijalandhika/tuto-api/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func registerRoutes(r *chi.Mux, pool *pgxpool.Pool, cfg *config.Config) {
	authQueries := authdb.New(pool)
	authService := auth.NewService(authQueries, cfg.JWT)
	authHandler := auth.NewHandler(authService)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok")) //nolint:errcheck
	})

	r.Post("/auth/signup", authHandler.Signup)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/refresh", authHandler.Refresh)
	r.Post("/auth/logout", authHandler.Logout)
}
