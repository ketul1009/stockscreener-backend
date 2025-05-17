package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/ketul1009/stockscreener-backend/handlers"
)

// InitRoutes initializes all v1 routes for the application
func InitRoutes(apiConfig handlers.ApiConfig) *chi.Mux {
	// Initialize v1 routes
	v1Router := chi.NewRouter()

	// Health check routes
	v1Router.Get("/healthz", handlers.ReadinessHandler)
	v1Router.Get("/err", handlers.ErrHandler)

	// User routes
	v1Router.Post("/users", apiConfig.HandlerCreateUser)
	v1Router.Get("/users", apiConfig.HandlerGetUsers)
	v1Router.Post("/login", apiConfig.HandlerLogin)
	v1Router.Post("/register", apiConfig.HandlerRegister)
	v1Router.Get("/me", apiConfig.HandlerGetUserFromToken)
	v1Router.Put("/update-user", apiConfig.HandlerUpdateUser)

	// Screener routes
	v1Router.Post("/screeners", apiConfig.HandlerCreateScreener)
	v1Router.Get("/screeners", apiConfig.HandlerGetScreeners)
	v1Router.Put("/screeners/{id}", apiConfig.HandlerUpdateScreener)
	v1Router.Delete("/screeners", apiConfig.HandlerDeleteScreener)

	// Job Producer and Result routes
	v1Router.Post("/jobs", apiConfig.HandlerCreateJob)
	v1Router.Get("/jobs/result", apiConfig.HandlerGetJobResult)
	v1Router.Get("/jobs", apiConfig.HandlerGetJobId)

	// Watchlist routes
	v1Router.Post("/watchlists", apiConfig.HandlerCreateWatchlist)
	v1Router.Get("/watchlists", apiConfig.HandlerGetWatchlist)
	v1Router.Get("/watchlists/all", apiConfig.HandlerGetAllWatchlists)
	v1Router.Put("/watchlists/{id}", apiConfig.HandlerUpdateWatchlist)
	v1Router.Delete("/watchlists/{id}", apiConfig.HandlerDeleteWatchlist)

	return v1Router
}
