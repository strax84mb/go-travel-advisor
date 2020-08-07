package database

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func saveUser(user User) error {
	if !gdb.First(&User{}, "username = ?", user.Username).RecordNotFound() {
		return &UsernameTakenError{}
	}
	salt := generateSalt()
	hexSalt := hex.EncodeToString(salt)
	hashedPassword := encodePassword(user.Password, salt)
	user.Salt = hexSalt
	user.Password = hashedPassword
	err := gdb.Create(&user).Error()
	if err != nil {
		log.Printf("Error while saving user! Error: %s\n", err.Error())
		return &StatementError{
			Message: "Error while saving user",
		}
	}
	return nil
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
func SaveNewUser(username string, password string) error {
	user := User{
		Username: username,
		Password: password,
		Role:     UserRoleUser,
	}
	return saveUser(user)
}

// GetUserByUsernameAndPassword - Get user from DB by username and verify password
// returns (UserDto, salt)
func GetUserByUsernameAndPassword(username string, password string) (*UserDto, string, error) {
	user := User{}
	if gdb.Where("username = ?", username).First(&user).RecordNotFound() {
		return nil, "", &NotFoundError{
			Message: fmt.Sprintf("Username %s not found", username),
		}
	}
	salt, _ := hex.DecodeString(user.Salt)
	encodedPassword := encodePassword(password, salt)
	if encodedPassword != user.Password {
		return &UserDto{}, "", &UnauthorizedError{
			Message: "Incorrect password!",
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
func GetUserSaltByUsername(username string) (string, error) {
	user := User{}
	if gdb.Select("salt").Where("username = ?", username).First(&user).RecordNotFound() {
		return "", &NotFoundError{
			Message: fmt.Sprintf("Username %s not found", username),
		}
	}
	return user.Salt, nil
}

func getUserByUsername(username string) (UserDto, error) {
	user := User{}
	if gdb.Where("username = ?", username).First(&user).RecordNotFound() {
		return UserDto{}, &NotFoundError{
			Message: fmt.Sprintf("Username %s not found", username),
		}
	}
	return UserDto{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}
