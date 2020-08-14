package database

// UserRoleAdmin - Admin role
var UserRoleAdmin = "ADMIN"

// UserRoleUser - User role
var UserRoleUser = "USER"

// UserRoleAny - Any role
var UserRoleAny = "any"

/*
// User - User entity
type User struct {
	ID       int64
	Username string `gorm:"type:varchar(30);unique_index;not null"`
	Password string `gorm:"type:varchar(100);not null"`
	Salt     string `gorm:"type:varchar(100);not null"`
	Role     string `gorm:"type:varchar(15);not null"`
}

// City - City entity
type City struct {
	ID       int64
	Name     string    `gorm:"type:varchar(100);not null"`
	Country  string    `gorm:"type:varchar(100);not null"`
	Comments []Comment `gorm:"foreignkey:CityID;association_foreignkey:ID"`
}

// Comment - Comment entity
type Comment struct {
	ID       int64
	CityID   int64      `gorm:"type:bigint;not null"`
	PosterID int64      `gorm:"type:bigint;not null"`
	Poster   User       `gorm:"foreignkey:ID;association_foreignkey:PosterID"`
	Text     string     `gorm:"type:varchar(255)"`
	Created  *time.Time `gorm:"type:datetime;not null"`
	Modified *time.Time `gorm:"type:datetime;not null"`
}

// Airport - Airport entity
type Airport struct {
	ID        int64
	AirportID int64  `gorm:"type:bigint;not null"`
	Name      string `gorm:"type:varchar(100);not null"`
	CityID    int64  `gorm:"type:bigint;not null"`
	City      City   `gorm:"foreignkey:ID;association_foreignkey:CityID"`
}

// Route - Route entity
type Route struct {
	ID            int64
	SourceID      int64   `gorm:"type:bigint;not null"`
	Source        Airport `gorm:"foreignkey:ID;association_foreignkey:SourceID"`
	DestinationID int64   `gorm:"type:bigint;not null"`
	Destination   Airport `gorm:"foreignkey:ID;association_foreignkey:DestinationID"`
	Price         float32 `gorm:"type:float;not null"`
}
*/
