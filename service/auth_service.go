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

func (s *AuthService) Register(ctx context.Context, username string, email string, password string) (string, int, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", 500, err
	}

	usernameExists, err := s.CheckUsernameExists(ctx, username)
	if err != nil {
		return "", 500, err
	}
	emailExists, err := s.CheckEmailExists(ctx, email)
	if err != nil {
		return "", 500, err
	}

	if usernameExists {
		return "", 403, errors.New("username already exists")
	}

	if emailExists {
		return "", 403, errors.New("email already exists")
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
		return "", 500, err
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		return "", 500, err
	}

	return token, 200, nil
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

func (s *AuthService) GetUserFromToken(ctx context.Context, token string) (*UserResponse, error) {
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	user, err := s.DB.GetUserByUsername(ctx, claims.Subject)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *AuthService) UpdateUser(ctx context.Context, id string, username string, email string) (*LoginResponse, error) {
	existingUser, err := s.DB.GetUserByID(ctx, pgtype.UUID{Bytes: uuid.MustParse(id), Valid: true})
	if err != nil {
		return nil, err
	}

	if existingUser.ID.Bytes != uuid.MustParse(id) {
		return nil, errors.New("user not found")
	}

	usernameExists, err := s.CheckUsernameExists(ctx, username)
	if err != nil {
		return nil, err
	}
	emailExists, err := s.CheckEmailExists(ctx, email)
	if err != nil {
		return nil, err
	}

	if usernameExists && existingUser.Username != username {
		return nil, errors.New("username already exists")
	}

	if emailExists && existingUser.Email != email {
		return nil, errors.New("email already exists")
	}

	user, err := s.DB.UpdateUser(ctx, db.UpdateUserParams{
		ID:        existingUser.ID,
		Username:  username,
		Email:     email,
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return nil, err
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
