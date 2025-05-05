package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ketul1009/stockscreener-backend/db"
	"github.com/ketul1009/stockscreener-backend/pkg/utils"
)

type AuthService struct {
	DB *db.Queries
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type UserResponse struct {
	ID        pgtype.UUID      `json:"id"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*LoginResponse, error) {
	user, err := s.DB.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	fmt.Println(user)

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: token,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, username string, email string, password string) (string, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	usernameExists, err := s.CheckUsernameExists(ctx, username)
	if err != nil {
		return "", err
	}
	emailExists, err := s.CheckEmailExists(ctx, email)
	if err != nil {
		return "", err
	}

	if usernameExists {
		return "", errors.New("username already exists")
	}

	if emailExists {
		return "", errors.New("email already exists")
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

func (s *AuthService) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	_, err := s.DB.GetUserByUsername(ctx, username)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (s *AuthService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	_, err := s.DB.GetUserByEmail(ctx, email)
	if err != nil {
		return false, nil
	}

	return true, nil
}
