package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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

	router.Mount("/v1", v1Router)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}

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
