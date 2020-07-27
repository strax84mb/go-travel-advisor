package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	// Used to initiate DB
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var gdb *gorm.DB

// InitDb - Used to initialize DB
func InitDb() {
	var err error
	os.Remove("test.db")
	gdb, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	gdb.AutoMigrate(&User{}, &City{}, &Comment{}, &Airport{}, &Route{})
	user := User{
		Username: "admin",
		Password: "admin",
		Role:     UserRoleAdmin,
	}
	genErr := saveUser(user)
	if genErr != nil {
		fmt.Println(genErr.Error())
	}
}

// SetDBForTesting - for testing purposes
func SetDBForTesting(db *gorm.DB) {
	gdb = db
}

func handleInitError(text string, err error) {
	if err != nil {
		log.Fatalf(text, err.Error())
	}
}
