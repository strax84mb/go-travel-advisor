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
	ID       int64
	Name     string    `gorm:"type:varchar(100)"`
	Country  string    `gorm:"type:varchar(100)"`
	Comments []Comment `gorm:"foreignkey:CityID;association_foreignkey:ID"`
}

// Comment - Comment entity
type Comment struct {
	ID       int64
	CityID   int64
	PosterID int64
	Poster   User   `gorm:"foreignkey:ID;association_foreignkey:PosterID"`
	Text     string `gorm:"type:varchar(255)"`
	Created  *time.Time
	Modified *time.Time
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
