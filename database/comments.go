package database

import (
	"database/sql"
	"fmt"
	"time"

	"gitlab.strale.io/go-travel/common"
)

// AddComment - saving a new comment
func AddComment(text string, username string, cityID int64) error {
	user, err := getUserByUsername(username)
	if err != nil {
		return err
	}
	_, found, err := GetCityByID(cityID, 0)
	if err != nil {
		return err
	}
	if !found {
		return &common.GeneralError{
			Message:   fmt.Sprintf("City with ID %d not found!", cityID),
			Location:  "database.comments.AddComment",
			ErrorType: common.CityNotFound,
		}
	}
	return performStatement(
		"database.comments.AddComment",
		func(params []interface{}) (sql.Result, error) {
			cityID := params[0].(int64)
			userID := params[1].(int64)
			text := params[2].(string)
			statement := `INSERT INTO comment (city_id, user_id, content, created, modified) VALUES ($1, $2, $3, $4, $5)`
			return db.Exec(statement, cityID, userID, text, time.Now(), time.Now())
		},
		cityID, user.ID, text)
}

func getCommentByID(id int64) (Comment, *common.GeneralError) {
	comment, found, err := performSingleSelection(
		"database.comments.getCommentByID",
		func(params []interface{}) (*sql.Rows, error) {
			commentID := params[0].(int64)
			query := `SELECT city_id, user_id, content, created, modified FROM comment WHERE id = $1`
			return db.Query(query, commentID)
		},
		func(rows *sql.Rows) (interface{}, error) {
			comment := Comment{
				id: id,
			}
			err := rows.Scan(&comment.cityID, &comment.userID, &comment.text, &comment.created, &comment.modified)
			return comment, err
		},
		id)
	if err != nil {
		return Comment{}, err
	}
	if !found {
		return Comment{}, &common.GeneralError{
			Message:   fmt.Sprintf("Comment with ID %d not found!", id),
			Location:  "database.comments.getCommentByID",
			ErrorType: common.CommentNotFount,
		}
	}
	return comment.(Comment), err

}

// UpdateComment - updating existing comment by user that posted it
func UpdateComment(id int64, text string, username string, cityID int64) error {
	// check if user exists
	user, err := getUserByUsername(username)
	if err != nil {
		return &common.GeneralError{
			Message:  "Error while finding user with given username",
			Location: "database.comments.UpdateComment",
			Cause:    err,
		}
	}
	// check if city exists
	_, found, err := GetCityByID(cityID, 0)
	if err != nil {
		return err
	}
	if !found {
		return &common.GeneralError{
			Message:   fmt.Sprintf("City with ID %d not found!", cityID),
			Location:  "database.comments.UpdateComment",
			ErrorType: common.CityNotFound,
		}
	}
	// check if comment exists
	comment, err := getCommentByID(id)
	if err != nil {
		return err
	}
	// check if original poster wants change
	if user.ID != comment.userID {
		return &common.GeneralError{
			Message:   "Comment may be changed by original poster only!",
			Location:  "database.comments.UpdateComment",
			ErrorType: common.UserNotAllowed,
		}
	}
	// do comment update
	return performStatement(
		"database.comments.UpdateComment",
		func(params []interface{}) (sql.Result, error) {
			comment := *(params[0].(*Comment))
			statement := `UPDATE comment SET content = $1, modified = $2 WHERE id = $3`
			return db.Exec(statement, comment.text, time.Now(), comment.id)
		},
		comment)
}

// DeleteComment - delete a comment by admin or original poster
func DeleteComment(id int64, username string) *common.GeneralError {
	// check if user exists
	user, err := getUserByUsername(username)
	if err != nil {
		return &common.GeneralError{
			Message:  "Error while finding user with given username",
			Location: "database.comments.DeleteComment",
			Cause:    err,
		}
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
	return performStatement(
		"database.comments.DeleteComment",
		func(params []interface{}) (sql.Result, error) {
			id := params[0].(int64)
			statement := `DELETE FROM comment WHERE id = $1`
			return db.Exec(statement, id)
		},
		id)
}

func countCommentsForCity(cityID int64) (int, *common.GeneralError) {
	result, _, err := performSingleSelection(
		"database.comments.countCommentsForCity",
		func(_ []interface{}) (*sql.Rows, error) {
			query := `SELECT COUNT(id) FROM comment WHERE city_id = $1`
			return db.Query(query, cityID)
		},
		func(rows *sql.Rows) (interface{}, error) {
			count := 0
			err := rows.Scan(&count)
			return count, err
		})
	if err != nil {
		return 0, err
	}
	return result.(int), nil
}

func getCommentsForCity(cityID int64, maxComments int) ([]CommentDto, *common.GeneralError) {
	count, err := countCommentsForCity(cityID)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}
	if maxComments < count && maxComments != -1 {
		count = maxComments
	}
	list := make([]CommentDto, count)
	err = performListSelection(
		"database.comments.getCommentsForCity",
		count, list[:],
		func(_ []interface{}) (*sql.Rows, error) {
			query := `SELECT comment.id, comment.content, user.username, comment.created, comment.modified
			FROM comment 
			JOIN user ON user.id = comment.user_id
			WHERE comment.city_id = $1
			ORDER BY comment.created DESC LIMIT $2`
			return db.Query(query, cityID, count)
		},
		func(rows *sql.Rows, pointer interface{}, index int) error {
			comment := CommentDto{}
			err := rows.Scan(&comment.ID, &comment.Text, &comment.Username, &comment.Created, &comment.Modified)
			if err != nil {
				return err
			}
			list := pointer.([]CommentDto)
			list[index] = comment
			return nil
		})
	return list, err
}

func deleteCommentsForCity(cityID int64) *common.GeneralError {
	return performStatement(
		"database.comments.deleteCommentsForCity",
		func(_ []interface{}) (sql.Result, error) {
			statement := `DELETE FROM comment WHERE city_id = $1`
			return db.Exec(statement, cityID)
		})
}
