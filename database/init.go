package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gitlab.strale.io/go-travel/common"

	// Used to initiate DB
	//_ "github.com/proullon/ramsql/driver"

	"github.com/jinzhu/gorm"
	// Used to initiate DB
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *sql.DB
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
	/*
		db, err = sql.Open("ramsql", "MyDataSource")
		handleInitError("sql.Open : Error : %s\n", err)
		_, err = db.Exec("CREATE TABLE user (id BIGSERIAL PRIMARY KEY NOT NULL, username TEXT NOT NULL, password TEXT NOT NULL, salt TEXT NOT NULL, role TEXT NOT NULL)")
		handleInitError("sql.Exec : Error : %s\n", err)
		user := User{
			username: "admin",
			password: "admin",
			role:     UserRoleAdmin,
		}
		err = saveUser(user)
		handleInitError("Error while savin user - Error: %s\n", err)
		err = saveUser(User{
			username: "strahinjamalabosna",
			password: "strale84",
			role:     UserRoleUser,
		})
		handleInitError("Error while savin user - Error: %s\n", err)
		_, err = db.Exec("CREATE TABLE city (id BIGSERIAL PRIMARY KEY NOT NULL, name TEXT NOT NULL, country TEXT NOT NULL)")
		handleInitError("sql.Exec : Error : %s\n", err)
		_, err = db.Exec("CREATE TABLE comment (id BIGSERIAL PRIMARY KEY NOT NULL, city_id BIGSERIAL NOT NULL, user_id BIGSERIAL NOT NULL, content TEXT NOT NULL, created NUMBER NOT NULL, modified NUMBER NOT NULL)")
		handleInitError("sql.Exec : Error : %s\n", err)
		generalError := AddNewCity("Subotica", "Srbija")
		hangleGeneralError(generalError)
		generalError = AddNewCity("Novi Sad", "Srbija")
		hangleGeneralError(generalError)
		generalError = AddComment("Admins comment", "admin", 1)
		hangleGeneralError(generalError)
		count, generalError := countCommentsForCity(1)
		fmt.Println(count)
		generalError = AddComment("Users comment", "strahinjamalabosna", 1)
		hangleGeneralError(generalError)
		count, generalError = countCommentsForCity(1)
		fmt.Println(count)
	*/
}

func handleInitError(text string, err error) {
	if err != nil {
		log.Fatalf(text, err.Error())
	}
}

func hangleGeneralError(err *common.GeneralError) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
