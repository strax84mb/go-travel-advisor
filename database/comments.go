package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.strale.io/go-travel/models"
)

// AddComment - saving a new comment
func AddComment(text string, username string, cityID int64) error {
	user, err := getUserByUsername(username)
	if err != nil {
		return err
	}
	_, err = GetCityByID(cityID, 0)
	if err != nil {
		return err
	}
	now := time.Now()
	comment := models.Comment{
		CityID:   cityID,
		Created:  now,
		Modified: now,
		PosterID: user.ID,
		Text:     text,
	}
	if err = comment.Insert(context.Background(), db, boil.Infer()); err != nil {
		log.Printf("Error while saving comment! Error: %s\n", err.Error())
		return &StatementError{
			Message: "Error while saving comment!",
		}
	}
	return nil
}

func getCommentByID(id int64) (*models.Comment, error) {
	comment, err := models.Comments(models.CommentWhere.ID.EQ(id)).One(context.Background(), db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{
				Message: fmt.Sprintf("Comment with ID %d not found!", id),
			}
		}
		log.Printf("Error while reading comment! Error: %s\n", err.Error())
		return nil, &StatementError{
			Message: "Error while reading comment!",
		}
	}
	return comment, nil
}

// UpdateComment - updating existing comment by user that posted it
func UpdateComment(id int64, text string, username string) error {
	// check if user exists
	user, err := getUserByUsername(username)
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
	comment.Text = text
	comment.Modified = time.Now()
	if _, err = comment.Update(context.Background(), db, boil.Infer()); err != nil {
		log.Printf("Error while updating comment! Error: %s\n", err.Error())
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
	comment, err := getCommentByID(id)
	if err != nil {
		return err
	}
	if user.Role != UserRoleAdmin && user.Username != username {
		return &ForbidenError{
			Message: "Comment may be changed by admin or original poster only!",
		}
	}
	if _, err = comment.Delete(context.Background(), db); err != nil {
		log.Printf("Error while deleting comment! Error: %s\n", err.Error())
		return &StatementError{
			Message: "Error while deleting comment!",
		}
	}
	return nil
}

func getCommentsForCity(cityID int64, maxComments int) ([]CommentDto, error) {
	comments, err := models.Comments(qm.Load(
		models.CommentRels.Poster),
		models.CommentWhere.CityID.EQ(cityID),
		qm.OrderBy(models.CommentColumns.Created+" desc"),
		qm.Limit(maxComments)).All(context.Background(), db)
	if err != nil {
		log.Printf("Error while fetching comments for city! Error: %s\n", err.Error())
		return nil, &StatementError{
			Message: "Error while fetching comments for city!",
		}
	}
	result := make([]CommentDto, len(comments))
	for i, c := range comments {
		result[i] = CommentDto{
			ID:       c.ID,
			Text:     c.Text,
			Username: c.R.Poster.Username,
			Created:  c.Created,
			Modified: c.Modified,
		}
	}
	return result, nil
}

func deleteCommentsForCity(cityID int64, tx *sql.Tx) error {
	_, err := models.Comments(models.CommentWhere.CityID.EQ(cityID)).DeleteAll(context.Background(), tx)
	if err != nil {
		log.Printf("Error while deleting comments for city with ID %d! Error: %s\n", cityID, err.Error())
		return &StatementError{
			Message: "Error while deleting comments for city!",
		}
	}
	return nil
}
