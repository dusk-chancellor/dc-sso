package service

import (
	"context"
	"errors"

	"github.com/dusk-chancellor/dc-sso/internal/repo"
	"github.com/dusk-chancellor/dc-sso/pkg/token"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(ctx context.Context, username, email, password string) (string, string, string, error) {
	childLogger := s.log.With(
		zap.String("operation", "Register()"),
		zap.String("username", username),
	)

	// hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		childLogger.Error("failed to generate hashed password:", zap.Error(err))
		return "", "", "", err
	}

	if _, err := s.db.CreateUser(ctx, username, email, hashedPassword); err != nil {
		if errors.Is(err, repo.ErrUserAlreadyExists) {
			return "", "", "", err
		}

		childLogger.Error("failed to create new user:", zap.Error(err))
		return "", "", "", err
	}
	
	// auto-logging after registration
	return s.Login(ctx, username, email, password)
}

func (s *Service) Login(ctx context.Context, username, email, password string) (string, string, string, error) {
	childLogger := s.log.With(
		zap.String("operation", "Login()"),
		zap.String("username", username),
	)

	user, err := s.getter.GetUserBy(ctx, repo.Email, email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return "", "", "", err
		}

		childLogger.Error("failed to retrieve user data:", zap.Error(err))
		return "", "", "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", "", "", ErrWrongPassword
		}

		childLogger.Error("failed to compare passwords", zap.Error(err))
		return "", "", "", err
	}

	accessToken, err := token.GenerateToken(&user, s.jwt.Secret, s.jwt.AccessTokenDuration)
	refreshToken, err := token.GenerateToken(&user, s.jwt.Secret, s.jwt.RefreshTokenDuration)
	if err != nil {
		childLogger.Warn("failed to generate tokens", zap.Error(err))
		return "", "", "", err
	}

	return user.ID.String(), accessToken, refreshToken, nil
}

func (s *Service) Logout(ctx context.Context, token string) (bool, error) {
	return true, nil
}
