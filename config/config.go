package config

import (
	"os"
	"time"
)

type Config struct {
	Port           string
	DBURL          string
	AllowedOrigins []string
	JWTSecret      string
	JWTExpiration  time.Duration
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, ErrMissingDBURL
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, ErrMissingJWTSecret
	}

	// Default to 24 hours if not specified
	jwtExpiration := 24 * time.Hour
	if exp := os.Getenv("JWT_EXPIRATION"); exp != "" {
		if duration, err := time.ParseDuration(exp); err == nil {
			jwtExpiration = duration
		}
	}

	// Default allowed origins
	allowedOrigins := []string{"https://*", "http://*"}
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		allowedOrigins = []string{origins}
	}

	return &Config{
		Port:           port,
		DBURL:          dbURL,
		AllowedOrigins: allowedOrigins,
		JWTSecret:      jwtSecret,
		JWTExpiration:  jwtExpiration,
	}, nil
}
