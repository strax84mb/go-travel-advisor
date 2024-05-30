package repository

import (
	"fmt"
	"time"

	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type commentRepository struct {
	db *database.Database
}

func NewCommentRepository(db *database.Database) *commentRepository {
	return &commentRepository{db: db}
}

func (cr *commentRepository) doList(pagination utils.Pagination, where func(*gorm.DB) *gorm.DB) ([]database.Comment, error) {
	tx := cr.db.DB
	if where != nil {
		tx = where(tx)
	}
	var comments []database.Comment
	tx.Limit(pagination.Limit).Offset(pagination.Offset).
		Order(clause.OrderByColumn{
			Column: clause.Column{
				Name: "created",
			},
			Desc: true,
		}).
		Find(&comments)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to read comments: %w", tx.Error)
	}
	return comments, nil
}

func (cr *commentRepository) ListComments(pagination utils.Pagination) ([]database.Comment, error) {
	return cr.doList(
		pagination,
		func(tx *gorm.DB) *gorm.DB {
			return tx
		},
	)
}

func (cr *commentRepository) ListCommentsForUser(userID int64, pagination utils.Pagination) ([]database.Comment, error) {
	return cr.doList(
		pagination,
		func(tx *gorm.DB) *gorm.DB {
			return tx.Where("poster_id = ?", userID)
		},
	)
}

func (cr *commentRepository) ListCommentsForCity(cityID int64, pagination utils.Pagination) ([]database.Comment, error) {
	return cr.doList(
		pagination,
		func(tx *gorm.DB) *gorm.DB {
			return tx.Where("city_id = ?", cityID)
		},
	)
}

type FindByIDInput struct {
	ID       int64
	LoadUser bool
	LoadCity bool
}

func (cr *commentRepository) FindByID(input FindByIDInput) (*database.Comment, error) {
	tx := cr.db.DB.Where("id = ?", input.ID)
	if input.LoadCity {
		tx = tx.Preload("City")
	}
	if input.LoadUser {
		tx = tx.Preload("Poster")
	}
	var comment database.Comment
	tx = tx.Take(&comment)
	switch {
	case tx.Error == gorm.ErrRecordNotFound:
		return nil, database.ErrNotFound
	case tx.Error != nil:
		return nil, fmt.Errorf("failed to read by ID: %w", tx.Error)
	default:
		return &comment, nil
	}
}

func (cr *commentRepository) Insert(comment database.Comment) (*database.Comment, error) {
	now := time.Now()
	comment.Modified = &now
	comment.Created = &now
	tx := cr.db.DB.Create(&comment)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to save: %w", tx.Error)
	}
	return &comment, nil
}

func (cr *commentRepository) Update(comment database.Comment) error {
	tx := cr.db.DB.Model(&comment).
		Where("id = ?", comment.ID).
		Updates(
			map[string]interface{}{
				"text":     comment.Text,
				"modified": time.Now(),
			},
		)
	switch {
	case tx.Error == gorm.ErrRecordNotFound || tx.RowsAffected == 0:
		return database.ErrNotFound
	case tx.Error != nil:
		return fmt.Errorf("failed to update: %w", tx.Error)
	default:
		return nil
	}
}

func (cr *commentRepository) Delete(id int64) error {
	tx := cr.db.DB.Delete(&database.Comment{}, id)
	switch {
	case tx.Error == gorm.ErrRecordNotFound || tx.RowsAffected == 0:
		return database.ErrNotFound
	case tx.Error != nil:
		return fmt.Errorf("failed to delete: %w", tx.Error)
	default:
		return nil
	}
}
