package repository

import (
	"fmt"

	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils"
	"gorm.io/gorm"
)

type routeRepository struct {
	db *database.Database
}

func NewRouteRepository(db *database.Database) *routeRepository {
	return &routeRepository{
		db: db,
	}
}

func (rr *routeRepository) Find(pagination utils.Pagination) ([]database.Route, error) {
	var list []database.Route
	tx := rr.db.DB.Limit(pagination.Limit).Offset(pagination.Offset).Find(&list)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to read routes: %w", tx.Error)
	}
	return list, nil
}

func (rr *routeRepository) listRoutes(
	incomming, outgoing bool,
	pagination utils.Pagination,
	where func(*gorm.DB) *gorm.DB,
) ([]database.Route, error) {
	var list []database.Route
	tx := where(rr.db.DB)
	if incomming {
		tx = tx.Where("source_id = ?")
	}

}
