package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

// AddComment - saving a new comment
func AddComment(text string, username string, cityID int64) error {
	user, err := getUserByUsername(username)
	if err != nil {
		return err
	}
	_, _, err = GetCityByID(cityID, 0)
	if err != nil {
		return err
	}
	now := time.Now()
	comment := Comment{
		CityID:   cityID,
		Created:  &now,
		Modified: &now,
		PosterID: user.ID,
		Text:     text,
	}
	if e := gdb.Create(&comment).Error; e != nil {
		log.Printf("Error while saving comment! Error: %s\n", e.Error())
		return &StatementError{
			Message: "Error while saving comment!",
		}
	}
	return nil
}

func getCommentByID(id int64) (Comment, error) {
	comment := Comment{}
	curDB := gdb.Find(&comment, id)
	if curDB.Error != nil {
		if curDB.RecordNotFound() {
			return Comment{}, &NotFoundError{
				Message: fmt.Sprintf("Comment with ID %d not found!", id),
			}
		}
		log.Printf("Error while reading comment! Error: %s\n", curDB.Error.Error())
		return Comment{}, &StatementError{
			Message: "Error while reading comment!",
		}
	}
	return comment, nil
}

// UpdateComment - updating existing comment by user that posted it
func UpdateComment(id int64, text string, username string, cityID int64) error {
	// check if user exists
	user, err := getUserByUsername(username)
	if err != nil {
		return err
	}
	// check if city exists
	_, _, err = GetCityByID(cityID, 0)
	if err != nil {
		return err
	}
	// check if comment exists
	comment, err := getCommentByID(id)
	if err != nil {
		return err
	}
	// check if original poster wants change
	if user.ID != comment.PosterID {
		return &ForbidenError{
			Message: "Comment may be changed by original poster only!",
		}
	}
	now := time.Now()
	comment.Text = text
	comment.Modified = &now
	if e := gdb.Save(&comment).Error; e != nil {
		log.Printf("Error while updating comment! Error: %s\n", e.Error())
		return &StatementError{
			Message: "Error while updating comment!",
		}
	}
	return nil
}

// DeleteComment - delete a comment by admin or original poster
func DeleteComment(id int64, username string) error {
	// check if user exists
	user, err := getUserByUsername(username)
	if err != nil {
		return err
	}
	_, err = getCommentByID(id)
	if err != nil {
		return err
	}
	if user.Role != UserRoleAdmin && user.Username != username {
		return &ForbidenError{
			Message: "Comment may be changed by admin or original poster only!",
		}
	}
	if e := gdb.Delete(&Comment{ID: id}).Error; e != nil {
		log.Printf("Error while deleting comment! Error: %s\n", e.Error())
		return &StatementError{
			Message: "Error while deleting comment!",
		}
	}
	return nil
}

func countCommentsForCity(cityID int64) (int, error) {
	var count int
	if e := gdb.Model(&Comment{}).Where(&Comment{CityID: cityID}).Count(&count).Error; e != nil {
		log.Printf("Error while counting comments for city! Error: %s\n", e.Error())
		return 0, &StatementError{
			Message: "Error while counting comments for city!",
		}
	}
	return count, nil
}

func getCommentsForCity(cityID int64, maxComments int) ([]CommentDto, error) {
	count, err := countCommentsForCity(cityID)
	if err != nil {
		return nil, err
	}
	if count == 0 || maxComments == 0 {
		return nil, nil
	}
	if maxComments < count && maxComments != -1 {
		count = maxComments
	}
	comments := make([]CommentDto, count)
	e := gdb.Debug().Table("comments").Select("comments.id, comments.text, users.username, comments.created, comments.modified").
		Joins("JOIN users ON users.id = comments.posert_id").
		Where(&Comment{CityID: cityID}).
		Order("created desc").
		Limit(count).
		Find(&comments).Error
	if e != nil {
		log.Printf("Error while fetching comments for city! Error: %s\n", e.Error())
		return nil, &StatementError{
			Message: "Error while fetching comments for city!",
		}
	}
	return comments, nil
}

func deleteCommentsForCity(cityID int64, tx *gorm.DB) error {
	if err := tx.Where(&Comment{CityID: cityID}).Delete(&Comment{}).Error; err != nil {
		log.Printf("Error while deleting comments for city with ID %d! Error: %s\n", cityID, err.Error())
		return &StatementError{
			Message: "Error while deleting comments for city!",
		}
	}
	return nil
}
