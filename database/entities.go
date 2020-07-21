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
	id       int64
	username string
	password string
	salt     string
	role     string
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
