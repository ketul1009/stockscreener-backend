package config

import "errors"

var (
	ErrMissingDBURL     = errors.New("DB_URL environment variable is required")
	ErrMissingJWTSecret = errors.New("JWT_SECRET environment variable is required")
)
