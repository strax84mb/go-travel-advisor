package database

import (
	"time"
)

// UserRoleAdmin - Admin role
var UserRoleAdmin = "ADMIN"

// UserRoleUser - User role
var UserRoleUser = "USER"

// UserRoleAny - Any role
var UserRoleAny = "any"

// User - User entity
type User struct {
	ID       int64
	Username string `gorm:"type:varchar(30);unique_index"`
	Password string `gorm:"type:varchar(100)"`
	Salt     string `gorm:"type:varchar(100)"`
	Role     string `gorm:"type:varchar(15)"`
}

// City - City entity
type City struct {
	id      int64
	name    string
	country string
}

// Comment - Comment entity
type Comment struct {
	id       int64
	cityID   int64
	userID   int64
	text     string
	created  time.Time
	modified time.Time
}

// Airport - Airport entity
type Airport struct {
	id        int64
	airportID int64
	name      string
	cityID    int64
}

// Route - Route entity
type Route struct {
	id          int64
	sourceID    int64
	destination int64
	price       float32
}
