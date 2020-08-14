package database

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"gitlab.strale.io/go-travel/models"
)

func saveUser(user models.User) error {
	found, err := models.Users(models.UserWhere.Username.EQ(user.Username)).Exists(context.Background(), db)
	if err != nil {
		return &StatementError{
			Message: "Error while if usernbame is available",
		}
	} else if found {
		return &UsernameTakenError{}
	}
	salt := generateSalt()
	hexSalt := hex.EncodeToString(salt)
	hashedPassword := encodePassword(user.Password, salt)
	user.Salt = hexSalt
	user.Password = hashedPassword
	log.Println(hexSalt)
	log.Println(hashedPassword)
	err = user.Insert(context.Background(), db, boil.Infer())
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
	user := models.User{
		Username: username,
		Password: password,
		Role:     UserRoleUser,
	}
	return saveUser(user)
}

// GetUserByUsernameAndPassword - Get user from DB by username and verify password
// returns (UserDto, salt)
func GetUserByUsernameAndPassword(username string, password string) (*UserDto, string, error) {
	user, err := models.Users(models.UserWhere.Username.EQ(username)).One(context.Background(), db)
	if err != nil && err == sql.ErrNoRows {
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
	user, err := models.Users(models.UserWhere.Username.EQ(username)).One(context.Background(), db)
	if err != nil && err == sql.ErrNoRows {
		return "", &NotFoundError{
			Message: fmt.Sprintf("Username %s not found", username),
		}
	}
	return user.Salt, nil
}

func getUserByUsername(username string) (UserDto, error) {
	user, err := models.Users(models.UserWhere.Username.EQ(username)).One(context.Background(), db)
	if err != nil && err == sql.ErrNoRows {
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
