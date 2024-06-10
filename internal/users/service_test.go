package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/users"
	"gitlab.strale.io/go-travel/internal/users/servicemocks"
	"golang.org/x/crypto/bcrypt"
)

type userService interface {
	Login(ctx context.Context, username, password string) (string, error)
	GetByUsername(ctx context.Context, username string) (*database.User, error)
	GetByID(ctx context.Context, id int64) (*database.User, error)
}

type userServiceTestKit struct {
	ctrl     *gomock.Controller
	service  userService
	userRepo servicemocks.MockiUserRepository
	secSrvc  servicemocks.MockiSecurityService
}

func setup(t *testing.T) *userServiceTestKit {
	logrus.SetLevel(logrus.FatalLevel)
	ctrl := gomock.NewController(t)
	userRepo := servicemocks.NewMockiUserRepository(ctrl)
	secSrvc := servicemocks.NewMockiSecurityService(ctrl)
	service := users.NewUserService(userRepo, secSrvc)
	return &userServiceTestKit{
		ctrl:     ctrl,
		service:  service,
		userRepo: *userRepo,
		secSrvc:  *secSrvc,
	}
}

func TestLogin_Success(t *testing.T) {
	kit := setup(t)
	defer kit.ctrl.Finish()
	ctx := context.Background()
	passBtyes, err := bcrypt.GenerateFromPassword([]byte("my_password"), bcrypt.DefaultCost)
	assert.NoError(t, err)
	user := database.User{
		ID:           1,
		Username:     "strale@example.com",
		PasswordHash: string(passBtyes),
		Roles: []database.UserRole{
			{
				ID:     11,
				UserID: 1,
				Role:   database.ROLE_ADMIN,
			},
		},
	}

	kit.userRepo.EXPECT().
		FindByUsername(gomock.Eq("strale@example.com"), gomock.Eq(true)).
		Times(1).
		Return(user, nil)
	kit.secSrvc.EXPECT().
		GenerateJWT(gomock.Any(), gomock.Eq(user)).
		Times(1).
		Return("success", nil)

	jwt, err := kit.service.Login(ctx, "strale@example.com", "my_password")
	assert.NoError(t, err)
	assert.Equal(t, "success", jwt)
}

func TestLogin_FailedJwtGeneration(t *testing.T) {
	kit := setup(t)
	defer kit.ctrl.Finish()
	ctx := context.Background()
	passBtyes, err := bcrypt.GenerateFromPassword([]byte("my_password"), bcrypt.DefaultCost)
	assert.NoError(t, err)
	user := database.User{
		ID:           1,
		Username:     "strale@example.com",
		PasswordHash: string(passBtyes),
		Roles: []database.UserRole{
			{
				ID:     11,
				UserID: 1,
				Role:   database.ROLE_ADMIN,
			},
		},
	}

	kit.userRepo.EXPECT().
		FindByUsername(gomock.Eq("strale@example.com"), gomock.Eq(true)).
		Times(1).
		Return(user, nil)
	kit.secSrvc.EXPECT().
		GenerateJWT(gomock.Any(), gomock.Eq(user)).
		Times(1).
		Return("", errors.New("some error"))

	_, err = kit.service.Login(ctx, "strale@example.com", "my_password")
	assert.EqualError(t, err, "error generating JWT: some error")
}

func TestLogin_IncorrectPassword(t *testing.T) {
	kit := setup(t)
	defer kit.ctrl.Finish()
	ctx := context.Background()
	user := database.User{
		ID:           1,
		Username:     "strale@example.com",
		PasswordHash: "incorrect password",
		Roles: []database.UserRole{
			{
				ID:     11,
				UserID: 1,
				Role:   database.ROLE_ADMIN,
			},
		},
	}

	kit.userRepo.EXPECT().
		FindByUsername(gomock.Eq("strale@example.com"), gomock.Eq(true)).
		Times(1).
		Return(user, nil)
	kit.secSrvc.EXPECT().
		GenerateJWT(gomock.Any(), gomock.Eq(user)).
		Times(0)

	_, err := kit.service.Login(ctx, "strale@example.com", "my_password")
	assert.EqualError(t, err, "unauthorized: incorrect password")
}

func TestLogin_FailedToLoadUser(t *testing.T) {
	kit := setup(t)
	defer kit.ctrl.Finish()
	ctx := context.Background()

	kit.userRepo.EXPECT().
		FindByUsername(gomock.Eq("strale@example.com"), gomock.Eq(true)).
		Times(1).
		Return(database.User{}, errors.New("some error"))
	kit.secSrvc.EXPECT().
		GenerateJWT(gomock.Any(), gomock.Any()).
		Times(0)

	_, err := kit.service.Login(ctx, "strale@example.com", "my_password")
	assert.EqualError(t, err, "user not read: some error")
}
