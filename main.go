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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/ketul1009/stockscreener-backend/config"
	"github.com/ketul1009/stockscreener-backend/db"
	"github.com/ketul1009/stockscreener-backend/handlers"
	"github.com/ketul1009/stockscreener-backend/middleware"
	"github.com/ketul1009/stockscreener-backend/pkg/logger"
	redisconn "github.com/ketul1009/stockscreener-backend/redis"
	"github.com/ketul1009/stockscreener-backend/routes"
	"github.com/ketul1009/stockscreener-backend/service"
	engine "github.com/ketul1009/stockscreener-backend/stock-engine"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

func createDBPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	// Configure the connection pool
	config.MaxConns = 25 // Maximum number of connections in the pool
	config.MinConns = 5  // Minimum number of connections to maintain
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func main() {
	// Load environment variables
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, skipping...")
		}
	}

	// Initialize logger
	logger.InitLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Initialize database connection pool
	ctx := context.Background()
	pool, err := createDBPool(ctx, cfg.DBURL)
	if err != nil {
		logger.Fatal("Failed to create database pool", zap.Error(err))
	}
	defer pool.Close()

	// Create a single DB instance to be shared
	dbInstance := db.WithPool(pool)

	// Initialize API config
	redisClient := redisconn.NewRedisClient()
	apiConfig := handlers.ApiConfig{
		DB:               dbInstance,
		AuthService:      &service.AuthService{DB: dbInstance},
		ScreenerService:  &service.ScreenerService{DB: dbInstance, Pool: pool},
		WatchlistService: &service.WatchlistService{DB: dbInstance, Pool: pool},
		RedisClient:      redisClient,
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

	stocks, err := apiConfig.ScreenerService.GetStockUniverse(ctx)
	if err != nil {
		logger.Fatal("Failed to get stock universe", zap.Error(err))
		stocks = []engine.Stock{}
	}

	go engine.StartWorker(redisClient, stocks, &engine.ApiConfig{DB: dbInstance})

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
