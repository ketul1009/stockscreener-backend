package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ketul1009/stockscreener-backend/db"
	"github.com/ketul1009/stockscreener-backend/pkg/utils"
)

type AuthService struct {
	DB *db.Queries
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.DB.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}
	return utils.GenerateJWT(user.Username)
}

func (s *AuthService) Register(ctx context.Context, username string, email string, password string) (string, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user, err := s.DB.CreateUser(ctx, db.CreateUserParams{
		ID:           pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return "", err
	}

	return utils.GenerateJWT(user.Username)
}
