package database

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gitlab.strale.io/go-travel/common"
	cmn "gitlab.strale.io/go-travel/common"
)

func saveUser(user User) *cmn.GeneralError {
	if !gdb.First(&User{}, "username = ?", user.Username).RecordNotFound() {
		return &cmn.GeneralError{
			Message:   "Username taken",
			Location:  "database.users.saveUser",
			ErrorType: cmn.UserNotAllowed,
		}
	}
	salt := generateSalt()
	hexSalt := hex.EncodeToString(salt)
	hashedPassword := encodePassword(user.Password, salt)
	user.Salt = hexSalt
	user.Password = hashedPassword
	err := gdb.Create(&user).Error
	if err != nil {
		log.Printf("Error while saving user! Error: %s\n", err.Error())
		return &cmn.GeneralError{
			Message:  "Error while saving user",
			Location: "database.users.saveUser",
			Cause:    err,
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
	rand.Seed(time.Now().UnixNano())
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
func SaveNewUser(username string, password string) *cmn.GeneralError {
	user := User{
		Username: username,
		Password: password,
		Role:     UserRoleUser,
	}
	return saveUser(user)
}

// GetUserByUsernameAndPassword - Get user from DB by username and verify password
// returns (UserDto, salt)
func GetUserByUsernameAndPassword(username string, password string) (*UserDto, string, *common.GeneralError) {
	user := User{}
	if gdb.Where("username = ?", username).First(&user).RecordNotFound() {
		return nil, "", &cmn.GeneralError{
			Message:   "Username not found",
			Location:  "database.users.GetUserByUsernameAndPassword",
			ErrorType: cmn.UserNotFound,
		}
	}
	salt, _ := hex.DecodeString(user.Salt)
	encodedPassword := encodePassword(password, salt)
	if encodedPassword != user.Password {
		return &UserDto{}, "", &cmn.GeneralError{
			Message:   "Incorrect password!",
			Location:  "database.users.GetUserByUsernameAndPassword",
			ErrorType: cmn.IncorrectPassword,
		}
	}
	userDto := &UserDto{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	return userDto, user.Salt, nil
}

// GetUserSaltByUsername - Get user salt for
func GetUserSaltByUsername(username string) (string, *common.GeneralError) {
	user := User{}
	if gdb.Select("salt").Where("username = ?", username).First(&user).RecordNotFound() {
		return "", &cmn.GeneralError{
			Message:  "Username does not exist!",
			Location: "database.users.GetUserSaltByUsername",
		}
	}
	return user.Salt, nil
}

func getUserByUsername(username string) (UserDto, *common.GeneralError) {
	user := User{}
	if gdb.Where("username = ?", username).First(&user).RecordNotFound() {
		return UserDto{}, &common.GeneralError{
			Message:   fmt.Sprintf("Username %s not found!", username),
			Location:  "database.users.getUserByUsername",
			ErrorType: common.UserNotFound,
		}
	}
	return UserDto{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}
