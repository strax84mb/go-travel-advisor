package database

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	ROLE_USER  = "USER"
	ROLE_ADMIN = "ADMIN"
)

type UserRole struct {
	ID     int64
	UserID int64 `gorm:"type:bigint;not null"`
	Role   Role  `gorm:"type:varchar(100);not null"`
}

type User struct {
	ID           int64
	Username     string     `gorm:"type:varchar(30);not null;index:idx_user_name,unique"`
	PasswordHash string     `gorm:"type:varchar(100);not null"`
	Roles        []UserRole `gorm:"foreignkey:UserID;association_foreignkey:ID"`
}

func (u *User) AfterCreate(tx *gorm.DB) {
	for i := range u.Roles {
		u.Roles[i].UserID = u.ID
	}
	tx.Create(&u.Roles)
}

type Airport struct {
	ID     int64
	Name   string `gorm:"type:varchar(100);not null"`
	CityID int64  `gorm:"type:bigint;not null"`
	City   *City  `gorm:"foreignkey:ID;association_foreignkey:CityID"`
}

type Comment struct {
	ID       int64
	CityID   int64      `gorm:"type:bigint;not null"`
	City     City       `gorm:"foreignkey:ID;association_foreignkey:CityID"`
	PosterID int64      `gorm:"type:bigint;not null"`
	Poster   User       `gorm:"foreignkey:ID;association_foreignkey:PosterID"`
	Text     string     `gorm:"type:varchar(255)"`
	Created  *time.Time `gorm:"type:datetime;not null"`
	Modified *time.Time `gorm:"type:datetime;not null"`
}

type City struct {
	ID       int64
	Name     string    `gorm:"type:varchar(100);not null;index:idx_city_name,unique"`
	Airports []Airport `gorm:"foreignkey:CityID;association_foreignkey:ID"`
	Comments []Comment `gorm:"foreignkey:CityID;association_foreignkey:ID"`
}

type Route struct {
	ID            int64
	SourceID      int64    `gorm:"type:bigint;not null"`
	Source        *Airport `gorm:"foreignkey:ID;association_foreignkey:SourceID"`
	DestinationID int64    `gorm:"type:bigint;not null"`
	Destination   *Airport `gorm:"foreignkey:ID;association_foreignkey:DestinationID"`
	Price         float32  `gorm:"type:float;not null"`
}
