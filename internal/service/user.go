package service

import (
	"context"
	"errors"

	"github.com/dusk-chancellor/dc-sso/internal/models"
	"github.com/dusk-chancellor/dc-sso/internal/repo"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// user service logic

// retrieves user data by id
func (s *Service) GetUser(ctx context.Context, id string) (models.User, error) {
	childLogger := s.log.With(
		zap.String("operation", "GetUser()"), 
		zap.String("user ID", id),
	)

	user, err := s.getter.GetUserBy(ctx, repo.ID, id)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) { // user does not exist
			return models.User{}, err
		}

		childLogger.Error("failed to get user:", zap.Error(err))
		return models.User{}, err
	}

	return user, nil
}

// updates user info, returning successfulness of op
func (s *Service) UpdateUser(ctx context.Context, id, username, email string) (bool, error) {
	childLogger := s.log.With(
		zap.String("operation", "UpdateUser()"),
	)

	var field repo.Field
	var value string

	switch {
	case username != "": // updating username
		field = repo.Username
		value = username
	case email != "": // updating email
		field = repo.Email
		value = email
	default:
		s.log.Warn("username & email empty")
	}

	childLogger.With(
		zap.Any("field", field),
		zap.String("value", value),
	)
	
	if err := s.updater.UpdateByID(ctx, id, field, value); err != nil {
		if errors.Is(err, repo.ErrUserNotFound) { // user does not exist
			return false, err
		}

		childLogger.Error("failed to update user:", zap.Error(err))
		return false, err
	}

	return true, nil
}

// compares old password with password saved & hashed in db,
// then updates user password
func (s *Service) ChangePassword(ctx context.Context, id, oldPassword, newPassword string) (bool, error) {
	childLogger := s.log.With(
		zap.String("operation", "ChangePassword"),
		zap.String("user id", id),
	)

	user, err := s.getter.GetUserBy(ctx, repo.ID, id)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) { // user does not exist
			return false, err
		}

		childLogger.Error("failed to get user by id:", zap.Error(err))
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(oldPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) { // passwords don't match
			return false, ErrWrongPassword
		}

		childLogger.Error("failed to compare hash and password:", zap.Error(err))
		childLogger.Debug("passwords:", zap.String("hashed", string(user.Password)), zap.String("oldPassword", oldPassword))
		return false, err
	}

	if err := s.updater.UpdateByID(ctx, id, repo.Password, newPassword); err != nil {
		childLogger.Error("failed to change password:", zap.Error(err))
		return false, err
	}

	return true, nil
}

// updates user role based on `email`
func (s *Service) SetRole(ctx context.Context, email, wantsRole string) (bool, error) {
	childLogger := s.log.With(
		zap.String("operation", "SetRole()"),
		zap.String("wants role", wantsRole),
	)

	if err := s.updater.UpdateRole(ctx, email, wantsRole); err != nil {
		if errors.Is(err, repo.ErrUserNotFound) { // user does not exist
			return false, err
		}

		childLogger.Error("failed to set role", zap.Error(err))
		return false, err
	}

	return true, nil
}
