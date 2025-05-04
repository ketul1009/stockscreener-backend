package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/ketul1009/stockscreener-backend/db"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not set")
	}

	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	apiConfig := apiConfig{
		DB: db.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", readinessHandler)
	v1Router.Get("/err", errHandler)
	v1Router.Post("/users", apiConfig.handlerCreateUser)

	router.Mount("/v1", v1Router)

	log.Printf("Server is running on port %s", port)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
