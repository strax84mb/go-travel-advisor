package users

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.strale.io/go-travel/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type iUserRepository interface {
	FindByID(id int64, loadUserRoles bool) (database.User, error)
	FindByUsername(username string, loadUserRoles bool) (database.User, error)
}

type iSecurityService interface {
	GenerateJWT(ctx context.Context, user database.User) (string, error)
}

type userService struct {
	userRepo iUserRepository
	secSrvc  iSecurityService
}

func NewUserService(userRepo iUserRepository, secSrvc iSecurityService) *userService {
	return &userService{
		userRepo: userRepo,
		secSrvc:  secSrvc,
	}
}

func (us *userService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := loadUser(
		ctx,
		func() (database.User, error) {
			return us.userRepo.FindByUsername(username, true)
		},
		func(e *logrus.Entry) *logrus.Entry {
			return e.WithField("userName", username)
		},
	)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrUnauthorized{message: "incorrect password"}
	}
	jwt, err := us.secSrvc.GenerateJWT(ctx, *user)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("error generating JWT")
		return "", fmt.Errorf("error generating JWT: %w", err)
	}
	return jwt, nil
}

func loadUser(
	ctx context.Context,
	getter func() (database.User, error),
	logFields func(*logrus.Entry) *logrus.Entry,
) (*database.User, error) {
	user, err := getter()
	switch {
	case err == database.ErrNotFound:
		entry := logrus.WithContext(ctx).WithError(err)
		entry = logFields(entry)
		entry.Warn("user not found")
		return nil, err
	case err != nil:
		entry := logrus.WithContext(ctx).WithError(err)
		entry = logFields(entry)
		entry.Error("user not read")
		return nil, fmt.Errorf("user not read: %w", err)
	default:
		return &user, nil
	}
}

func (us *userService) GetByUsername(ctx context.Context, username string) (*database.User, error) {
	return loadUser(
		ctx,
		func() (database.User, error) {
			return us.userRepo.FindByUsername(username, false)
		},
		func(e *logrus.Entry) *logrus.Entry {
			return e.WithField("userName", username)
		},
	)
}

func (us *userService) GetByID(ctx context.Context, id int64) (*database.User, error) {
	return loadUser(
		ctx,
		func() (database.User, error) {
			return us.userRepo.FindByID(id, false)
		},
		func(e *logrus.Entry) *logrus.Entry {
			return e.WithField("userId", id)
		},
	)
}
