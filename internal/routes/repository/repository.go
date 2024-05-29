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

func (rr *routeRepository) doFind(
	pagination utils.Pagination,
	adjustQuery func(*gorm.DB) *gorm.DB,
) ([]database.Route, error) {
	var list []database.Route
	tx := rr.db.DB.Limit(pagination.Limit).Offset(pagination.Offset)
	if adjustQuery != nil {
		tx = adjustQuery(tx)
	}
	tx = tx.Find(&list)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to read routes: %w", tx.Error)
	}
	return list, nil
}

func (rr *routeRepository) Find(pagination utils.Pagination) ([]database.Route, error) {
	return rr.doFind(pagination, nil)
}

func (rr *routeRepository) FindRoutesForAirport(
	airportID int64,
	incomming, outgoing bool,
	pagination utils.Pagination,
) ([]database.Route, error) {
	return rr.doFind(
		pagination,
		func(tx *gorm.DB) *gorm.DB {
			switch {
			case incomming && outgoing:
				tx = tx.Where(tx.Where("routes.destination_id = ?", airportID).Or("routes.source_id = ?", airportID))
			case incomming && !outgoing:
				tx = tx.Where("routes.destination_id = ?", airportID)
			case !incomming && outgoing:
				tx = tx.Where("routes.source_id = ?", airportID)
			}
			return tx
		},
	)
}

func (rr *routeRepository) FindRoutesForCity(
	cityID int64,
	incomming, outgoing bool,
	pagination utils.Pagination,
) ([]database.Route, error) {
	return rr.doFind(
		pagination,
		func(tx *gorm.DB) *gorm.DB {
			switch {
			case incomming && outgoing:
				tx = tx.Joins("JOIN airports ON airports.id = routes.destination_id OR airports.id = routes.source_id")
			case incomming && !outgoing:
				tx = tx.Joins("JOIN airports ON airports.id = routes.destination_id")
			case !incomming && outgoing:
				tx = tx.Joins("JOIN airports ON airports.id = routes.source_id")
			}
			return tx.Where("airports.city_id = ?", cityID)
		},
	)
}

func (rr *routeRepository) FindByID(id int64, loadAirports bool) (*database.Route, error) {
	tx := rr.db.DB.Where("id = ?", id)
	if loadAirports {
		tx = tx.Preload("Source").Preload("Destination")
	}
	var route database.Route
	tx = tx.Take(&route)
	switch {
	case tx.Error == gorm.ErrRecordNotFound:
		return nil, database.ErrNotFound
	case tx.Error != nil:
		return nil, fmt.Errorf("failed to read route: %w", tx.Error)
	default:
		return &route, nil
	}
}

func (rr *routeRepository) Insert(route database.Route) (*database.Route, error) {
	tx := rr.db.DB.Create(&route)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to save new route: %w", tx.Error)
	}
	return &route, nil
}

func (rr *routeRepository) UpdatePrice(id int64, price float32) error {
	tx := rr.db.DB.Where("id = ?", id).Update("price", price)
	if tx.Error != nil {
		return fmt.Errorf("failed to update route price: %w", tx.Error)
	}
	return nil
}

func (rr *routeRepository) Delete(id int64) error {
	tx := rr.db.DB.Where("id = ?", id).Delete(&database.Route{})
	if tx.Error != nil {
		return fmt.Errorf("failed to delete route: %w", tx.Error)
	}
	return nil
}

func (rr *routeRepository) FindDestinations(startAirportsIDs []int64, cityIDsToSkip []int64) ([]database.Route, error) {
	var list []database.Route
	tx := rr.db.DB.Where("routes.start_id IN ?", startAirportsIDs).
		Joins("airports ON routes.destination_id = airports.id").
		Where("airports.city_id NOT IN ?", cityIDsToSkip).
		Preload("Source").Preload("Destination").
		Find(&list)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to list routes: %w", tx.Error)
	}
	return list, nil
}
