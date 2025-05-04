package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/ketul1009/stockscreener-backend/config"
	"github.com/ketul1009/stockscreener-backend/db"
	"github.com/ketul1009/stockscreener-backend/handlers"
	"github.com/ketul1009/stockscreener-backend/middleware"
	"github.com/ketul1009/stockscreener-backend/pkg/logger"
	"github.com/ketul1009/stockscreener-backend/routes"
	"github.com/ketul1009/stockscreener-backend/service"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize logger
	logger.InitLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Initialize database connection
	conn, err := pgx.Connect(context.Background(), cfg.DBURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer conn.Close(context.Background())

	// Initialize API config
	apiConfig := handlers.ApiConfig{
		DB:          db.New(conn),
		AuthService: &service.AuthService{DB: db.New(conn)},
	}

	// Create base router
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin", "X-Requested-With", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		ExposedHeaders:   []string{"Link", "Content-Type", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		AllowCredentials: true,
		MaxAge:           300,
		Debug:            true,
	}))

	// Initialize routes
	router.Mount("/v1", routes.InitRoutes(apiConfig))

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server is starting", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Server is shutting down...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
