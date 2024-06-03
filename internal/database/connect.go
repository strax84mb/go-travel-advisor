package database

import (
	"fmt"
	"os"

	"gitlab.strale.io/go-travel/internal/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func ConnectToDatabase(dbConfig config.DbConfig) (Database, error) {
	db, err := gorm.Open(sqlite.Open(dbConfig.URL), &gorm.Config{})
	if err != nil {
		return Database{}, fmt.Errorf("failed to connect to database: %w", err)
	}
	db.AutoMigrate(&Comment{}, &Airport{}, &City{}, &User{}, &UserRole{}, &Route{})

	if os.Getenv("INIT_DB") == "true" {
		user := User{}
		tx := db.Where("username = ?", "admin").Take(&user)
		if tx.Error == gorm.ErrRecordNotFound {
			pass, err := bcrypt.GenerateFromPassword([]byte("admin_pass"), bcrypt.DefaultCost)
			if err != nil {
				return Database{}, fmt.Errorf("failed to create admin: %w", err)
			}
			db.Save(&User{
				Username:     "admin",
				PasswordHash: string(pass),
				Roles: []UserRole{
					{
						Role: ROLE_ADMIN,
					},
				},
			})
		}
		user = User{}
		tx = db.Where("username = ?", "user_name").Take(&user)
		if tx.Error == gorm.ErrRecordNotFound {
			pass, err := bcrypt.GenerateFromPassword([]byte("user_pass"), bcrypt.DefaultCost)
			if err != nil {
				return Database{}, fmt.Errorf("failed to create user: %w", err)
			}
			db.Save(&User{
				Username:     "user_name",
				PasswordHash: string(pass),
				Roles: []UserRole{
					{
						Role: ROLE_USER,
					},
				},
			})
		}
	}

	return Database{DB: db}, nil
}
