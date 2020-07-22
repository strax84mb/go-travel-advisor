package database

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"

	"gitlab.strale.io/go-travel/common"
	cmn "gitlab.strale.io/go-travel/common"
)

func saveUser(user User) error {
	// Check if username is taken
	query := `SELECT COUNT(user.id) FROM user WHERE user.username = '$1'`
	rows, err := db.Query(query, user.username)
	defer rows.Close()
	if err != nil {
		return &cmn.GeneralError{
			Message:  "Error while checking if username is taken",
			Location: "database.users.saveUser",
		}
	}
	var columns []string
	var usercount int
	for rows.Next() {
		columns, err = rows.Columns()
		err = rows.Scan(&usercount)
	}
	if err != nil {
		return &cmn.GeneralError{
			Message:  "Error while reading columns",
			Location: "database.users.saveUser",
		}
	}
	if len(columns) > 0 {
		usercount, _ = strconv.Atoi(columns[0])
	}
	if usercount > 0 {
		return &usernameTakenError{}
	}
	// Generate salt and encode password
	salt := generateSalt()
	hexSalt := hex.EncodeToString(salt)
	hashedPassword := encodePassword(user.password, salt)
	user.salt = hexSalt
	user.password = hashedPassword
	// Save user
	res, err := db.Exec("INSERT INTO user (username, password, salt, role) VALUES ($1, $2, $3, $4)",
		user.username, user.password, user.salt, user.role)
	if rowsAffected, _ := res.RowsAffected(); rowsAffected < 1 {
		return &cmn.GeneralError{
			Message:  "Insert into user failed!",
			Location: "database.users.saveUser",
		}
	}
	return nil
}

type usernameTakenError struct{}

func (u *usernameTakenError) Error() string {
	return "Username is taken!"
}

func generateSalt() []byte {
	salt := make([]byte, 16)
	rand.Read(salt)
	return salt
}

func encodePassword(password string, salt []byte) string {
	h := sha512.New()
	h.Write(salt)
	h.Write([]byte(password))
	hashedPassword := h.Sum(salt)
	return hex.EncodeToString(hashedPassword)
}

// SaveNewUser - save new user in db
func SaveNewUser(username string, password string) error {
	user := User{
		username: username,
		password: password,
		role:     UserRoleUser,
	}
	return saveUser(user)
}

// GetUserByUsernameAndPassword - Get user from DB by username and verify password
// returns (UserDto, salt)
func GetUserByUsernameAndPassword(username string, password string) (*UserDto, string, error) {
	query := `SELECT user.id, user.username, user.password, user.role, user.salt FROM user WHERE user.username = $1`
	rows, err := db.Query(query, username)
	defer rows.Close()
	if err != nil {
		return &UserDto{}, "", &cmn.GeneralError{
			Message:  "Error while checking if username is taken",
			Location: "database.users.GetUserByUsernameAndPassword",
		}
	}
	var id int
	var user string
	var pass string
	var userrole string
	var usersalt string
	for rows.Next() {
		err = rows.Scan(&id, &user, &pass, &userrole, &usersalt)
		if err != nil {
			return &UserDto{}, "", &cmn.GeneralError{
				Message:  "Error while reading columns",
				Location: "database.users.GetUserByUsernameAndPassword",
			}
		}
		salt, _ := hex.DecodeString(usersalt)
		encodedPassword := encodePassword(password, salt)
		if encodedPassword != pass {
			return &UserDto{}, "", &cmn.GeneralError{
				Message:  "Incorrect password!",
				Location: "database.users.GetUserByUsernameAndPassword",
			}
		}
	}
	if user == "" {
		return &UserDto{}, "", &cmn.GeneralError{
			Message:  "Incorrect username!",
			Location: "database.users.GetUserByUsernameAndPassword",
		}
	}
	userDto := &UserDto{
		ID:       int64(id),
		Username: user,
		Role:     userrole,
	}
	return userDto, usersalt, nil
}

// GetUserSaltByUsername - Get user salt for
func GetUserSaltByUsername(username string) (string, error) {
	query := `SELECT user.salt FROM user WHERE user.username = $1`
	rows, err := db.Query(query, username)
	defer rows.Close()
	if err != nil {
		return "", &cmn.GeneralError{
			Message:  "Error while checking if username is taken",
			Location: "database.users.GetUserSaltByUsername",
		}
	}
	var salt string
	for rows.Next() {
		err = rows.Scan(&salt)
		if err != nil {
			return "", &cmn.GeneralError{
				Message:  "Error while reading columns",
				Location: "database.users.GetUserSaltByUsername",
			}
		}
	}
	if salt == "" {
		return "", &cmn.GeneralError{
			Message:  "Username does not exist!",
			Location: "database.users.GetUserSaltByUsername",
		}
	}
	return salt, nil
}

func getUserByUsername(username string) (UserDto, *common.GeneralError) {
	value, found, err := performSingleSelection(
		"database.users.getUserByUsername",
		func(_ []interface{}) (*sql.Rows, error) {
			query := `SELECT id, role FROM user WHERE username = $1`
			return db.Query(query, username)
		},
		func(rows *sql.Rows) (interface{}, error) {
			user := UserDto{
				Username: username,
			}
			err := rows.Scan(&user.ID, &user.Role)
			return user, err
		})
	if err != nil {
		return UserDto{}, err
	}
	if !found {
		return UserDto{}, &common.GeneralError{
			Message:   fmt.Sprintf("Username %s not found!", username),
			Location:  "database.users.getUserByUsername",
			ErrorType: common.UserNotFound,
		}
	}
	return value.(UserDto), nil
}
