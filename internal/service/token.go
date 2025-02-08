package service

import (
	"errors"

	"github.com/dusk-chancellor/dc-sso/internal/models"
	tkn "github.com/dusk-chancellor/dc-sso/pkg/token"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// validates token by parsing it
func (s *Service) ValidateToken(token string) (bool, error) {
	childLogger := s.log.With(
		zap.String("operation", "ValidateToken()"),
	)

	if _, err := tkn.ParseToken(token, s.jwt.Secret); err != nil {
		if errors.Is(err, tkn.ErrInvalidToken) {
			return false, err
		}

		childLogger.Warn("Failed to parse token:", zap.Error(err))
		childLogger.Debug("unparsed token:", zap.String("token", token))
		return false, err
	}

	return true, nil
}

// refreshes both access & refresh tokens at once
func (s *Service) RefreshToken(token string) (string, string, error) {
	childLogger := s.log.With(
		zap.String("operation", "RefreshToken()"),
	)

	claims, err := tkn.ParseToken(token, s.jwt.Secret)
	if err != nil {
		if errors.Is(err, tkn.ErrInvalidToken) {
			return "", "", err
		}
		
		childLogger.Warn("Failed to parse token:", zap.Error(err))
		childLogger.Debug("unparsed token:", zap.String("token", token))
		return "", "", err
	}
	//
	userID := uuid.MustParse(claims.UserID)
	// 
	userData := &models.User{
		ID: userID,
		Username: claims.Subject,
		Email: claims.Email,
	}

	accessToken, err := tkn.GenerateToken(userData, s.jwt.Secret, s.jwt.AccessTokenDuration)
	refreshToken, err := tkn.GenerateToken(userData, s.jwt.Secret, s.jwt.RefreshTokenDuration)
	if err != nil {
		childLogger.Error("Failed to generate new tokens:", zap.Error(err))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
