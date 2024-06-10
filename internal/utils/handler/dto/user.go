package dto

import "gitlab.strale.io/go-travel/internal/database"

//go:generate ffjson -nodecoder $GOFILE

type UserDto struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type LoginTokenDto struct {
	Token string `json:"token"`
}

func UserToDto(user *database.User) *UserDto {
	return &UserDto{
		ID:       user.ID,
		Username: user.Username,
	}
}
