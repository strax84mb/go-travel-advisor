package database

import "time"

// UserDto - user object to expose
type UserDto struct {
	ID       int64
	Username string
	Role     string
}

// CommentDto - comment object to expose
type CommentDto struct {
	ID       int64
	Text     string
	Username string
	Created  time.Time
	Modified time.Time
}

// CityDto - city object to expose
type CityDto struct {
	ID       int64
	Name     string
	Country  string
	Comments []CommentDto
}

// AirportDto - airport object to expose
type AirportDto struct {
	ID        int64
	AirportID int64
	Name      string
	City      CityDto
}

// RouteDto - route object to expose
type RouteDto struct {
	ID          int64
	Source      AirportDto
	Destination AirportDto
	price       float32
}
