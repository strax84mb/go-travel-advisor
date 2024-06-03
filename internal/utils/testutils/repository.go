package testutils

import (
	"fmt"
	"log"
	"os"

	"gitlab.strale.io/go-travel/internal/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectToTestDB(logSQL bool) (database.Database, error) {
	var lggr logger.Interface
	if logSQL {
		lggr = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
			},
		)
	}
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: lggr,
	})
	if err != nil {
		return database.Database{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.AutoMigrate(
		&database.Comment{},
		&database.Airport{},
		&database.City{},
		&database.User{},
		&database.UserRole{},
		&database.Route{},
	)

	return database.Database{
		DB: db,
	}, nil
}
