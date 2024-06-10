package dto

//go:generate ffjson -noencoder $GOFILE

type UserLoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
