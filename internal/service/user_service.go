package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Kihouww/mini-ai-compute-platform/internal/model"
	"github.com/Kihouww/mini-ai-compute-platform/internal/repository"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Register(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error) {
	username := strings.TrimSpace(req.Username)
	password := req.Password
	email := strings.TrimSpace(req.Email)

	if username == "" {
		return nil, fmt.Errorf("%w: username is required", ErrInvalidInput)
	}

	if len(password) < 6 {
		return nil, fmt.Errorf("%w: password must be at least 6 characters", ErrInvalidInput)
	}

	existingUser, err := s.userRepo.FindByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(hashBytes),
		Email:        email,
		Status:       "active",
	}

	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &model.RegisterResponse{
		UserID:   userID,
		Username: username,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	username := strings.TrimSpace(req.Username)
	password := req.Password

	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.FindByUsername(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidCredentials
	}

	if err != nil {
		return nil, err
	}

	if user.Status != "active" {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token := fmt.Sprintf("mock-token-user-%d-%d", user.ID, time.Now().Unix())

	return &model.LoginResponse{
		Token:  token,
		UserID: user.ID,
	}, nil
}
