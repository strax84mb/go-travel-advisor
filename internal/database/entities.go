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

func (u *User) AfterCreate(tx *gorm.DB) error {
	for i := range u.Roles {
		u.Roles[i].UserID = u.ID
	}
	tx.Create(&u.Roles)
	return tx.Error
}

type Airport struct {
	ID     int64
	Name   string `gorm:"type:varchar(100);not null"`
	CityID int64  `gorm:"type:bigint;not null"`
	City   *City  `gorm:"foreignkey:CityID;association_foreignkey:ID"`
}

type Comment struct {
	ID       int64
	CityID   int64      `gorm:"type:bigint;not null"`
	City     City       `gorm:"foreignkey:CityID;association_foreignkey:ID"`
	PosterID int64      `gorm:"type:bigint;not null"`
	Poster   User       `gorm:"foreignkey:PosterID;association_foreignkey:ID"`
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
	Source        *Airport `gorm:"foreignkey:SourceID;association_foreignkey:ID"`
	DestinationID int64    `gorm:"type:bigint;not null"`
	Destination   *Airport `gorm:"foreignkey:DestinationID;association_foreignkey:ID"`
	Price         float32  `gorm:"type:float;not null"`
}
