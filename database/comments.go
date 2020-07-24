package database

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.strale.io/go-travel/common"
)

// AddComment - saving a new comment
func AddComment(text string, username string, cityID int64) *common.GeneralError {
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
		return &common.GeneralError{
			Message:  "Error while saving comment!",
			Location: "database.comments.AddComment",
			Cause:    e,
		}
	}
	return nil
}

func getCommentByID(id int64) (Comment, *common.GeneralError) {
	comment := Comment{}
	curDB := gdb.Find(&comment, id)
	if curDB.Error != nil {
		if curDB.RecordNotFound() {
			return Comment{}, &common.GeneralError{
				Message:   "Comment with ID %d not found!",
				Location:  "database.comments.getCommentByID",
				ErrorType: common.CommentNotFount,
			}
		}
		log.Printf("Error while reading comment! Error: %s\n", curDB.Error.Error())
		return Comment{}, &common.GeneralError{
			Message:  "Error while reading comment!",
			Location: "database.comments.getCommentByID",
			Cause:    curDB.Error,
		}
	}
	return comment, nil
}

// UpdateComment - updating existing comment by user that posted it
func UpdateComment(id int64, text string, username string, cityID int64) *common.GeneralError {
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
		return &common.GeneralError{
			Message:   "Comment may be changed by original poster only!",
			Location:  "database.comments.UpdateComment",
			ErrorType: common.UserNotAllowed,
		}
	}
	now := time.Now()
	comment.Text = text
	comment.Modified = &now
	if e := gdb.Save(&comment).Error; e != nil {
		log.Printf("Error while updating comment! Error: %s\n", e.Error())
		return &common.GeneralError{
			Message:  "Error while updating comment!",
			Location: "database.comments.UpdateComment",
			Cause:    e,
		}
	}
	return nil
}

// DeleteComment - delete a comment by admin or original poster
func DeleteComment(id int64, username string) *common.GeneralError {
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
		return &common.GeneralError{
			Message:   "Comment may be changed by admin or original poster only!",
			Location:  "database.comments.DeleteComment",
			ErrorType: common.UserNotAllowed,
		}
	}
	if e := gdb.Delete(&Comment{ID: id}).Error; e != nil {
		log.Printf("Error while deleting comment! Error: %s\n", e.Error())
		return &common.GeneralError{
			Message:  "Error while deleting comment!",
			Location: "database.comments.DeleteComment",
			Cause:    e,
		}
	}
	return nil
}

func countCommentsForCity(cityID int64) (int, *common.GeneralError) {
	var count int
	if e := gdb.Model(&Comment{}).Where(&Comment{CityID: cityID}).Count(&count).Error; e != nil {
		log.Printf("Error while counting comments for city! Error: %s\n", e.Error())
		return 0, &common.GeneralError{
			Message:  "Error while counting comments for city!",
			Location: "database.comments.countCommentsForCity",
			Cause:    e,
		}
	}
	return count, nil
}

func getCommentsForCity(cityID int64, maxComments int) ([]CommentDto, *common.GeneralError) {
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
		return nil, &common.GeneralError{
			Message:  "Error while fetching comments for city!",
			Location: "database.comments.getCommentsForCity",
			Cause:    e,
		}
	}
	return comments, nil
}

func deleteCommentsForCity(cityID int64, tx *gorm.DB) error {
	return tx.Where(&Comment{CityID: cityID}).Delete(&Comment{}).Error
}
