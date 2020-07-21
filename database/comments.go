package database

import (
	"time"

	"gitlab.strale.io/go-travel/common"
)

// AddComment - saving a new comment
func AddComment(text string, username string, cityID int64) error {
	user, err := getUserByUsername(username)
	if err != nil {
		return &common.GeneralError{
			Message:  "Error while finding user with given username",
			Location: "database.comments.AddComment",
			Cause:    err,
		}
	}
	_, found, err := GetCityByID(cityID, 0)
	if err != nil || !found {
		return &common.GeneralError{
			Message:  "Could not find city with given ID",
			Location: "database.comments.AddComment",
			Cause:    err,
		}
	}
	statement := `INSERT INTO comment (city_id, user_id, content, created, modified) VALUES ($1, $2, $3, $4, $5)`
	result, err := db.Exec(statement, cityID, user.ID, text, time.Now(), time.Now())
	if err != nil {
		return &common.GeneralError{
			Message:  "Error while saving comment",
			Location: "database.comments.AddComment",
			Cause:    err,
		}
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return &common.GeneralError{
			Message:  "Error while checking insert result",
			Location: "database.comments.AddComment",
			Cause:    err,
		}
	}
	if affected < 1 {
		return &common.GeneralError{
			Message:  "No comment saved",
			Location: "database.comments.AddComment",
			Cause:    err,
		}
	}
	return nil
}
