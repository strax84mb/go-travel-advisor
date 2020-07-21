package database

import (
	"database/sql"
	"log"

	// Used to initiate DB
	_ "github.com/proullon/ramsql/driver"
)

var db *sql.DB

// InitDb - Used to initialize DB
func InitDb() {
	var err error
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
	_, err = db.Exec("CREATE TABLE city (id BIGSERIAL PRIMARY KEY NOT NULL, name TEXT NOT NULL, country TEXT NOT NULL)")
	handleInitError("sql.Exec : Error : %s\n", err)
	_, err = db.Exec("CREATE TABLE comment (id BIGSERIAL PRIMARY KEY NOT NULL, city_id BIGSERIAL NOT NULL, user_id BIGSERIAL NOT NULL, content TEXT NOT NULL, created DATETIME NOT NULL, modified DATETIME NOT NULL)")
	handleInitError("sql.Exec : Error : %s\n", err)
}

func handleInitError(text string, err error) {
	if err != nil {
		log.Fatalf(text, err.Error())
	}
}
