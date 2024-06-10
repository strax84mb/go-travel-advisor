package repository

import (
	"fmt"

	"gitlab.strale.io/go-travel/internal/database"
	"gorm.io/gorm"
)

type userRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) *userRepository {
	return &userRepository{
		db: db,
	}
}
func (ur *userRepository) findUser(where func(*gorm.DB) *gorm.DB, loadUserRoles bool) (database.User, error) {
	var user database.User
	tx := where(ur.db.DB)
	if loadUserRoles {
		tx = tx.Preload("Roles")
	}
	tx = tx.Take(&user)
	switch {
	case tx.Error == gorm.ErrRecordNotFound:
		return database.User{}, database.ErrNotFound
	case tx.Error != nil:
		return database.User{}, fmt.Errorf("failed to read user: %w", tx.Error)
	default:
		return user, nil
	}
}

func (ur *userRepository) FindByID(id int64, loadUserRoles bool) (database.User, error) {
	return ur.findUser(
		func(tx *gorm.DB) *gorm.DB {
			return tx.Where("id = ?", id)
		},
		loadUserRoles,
	)
}

func (ur *userRepository) FindByUsername(username string, loadUserRoles bool) (database.User, error) {
	return ur.findUser(
		func(tx *gorm.DB) *gorm.DB {
			return tx.Where("username = ?", username)
		},
		loadUserRoles,
	)
}

func (ur *userRepository) FindUsernamesByIDs(ids []int64) (map[int64]string, error) {
	var users []database.User
	tx := ur.db.DB.Model(&database.User{}).
		Where("id IN ?", ids).
		Select("id", "username").
		Scan(&users)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to read usernames for IDs: %w", tx.Error)
	}
	usernames := make(map[int64]string)
	for _, user := range users {
		usernames[user.ID] = user.Username
	}
	return usernames, nil
}
