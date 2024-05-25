package repository

import (
	"fmt"

	"gitlab.strale.io/go-travel/internal/database"
	"gorm.io/gorm"
)

type AirportRepository struct {
	db *database.Database
}

func NewAirportRepository(db *database.Database) *AirportRepository {
	return &AirportRepository{
		db: db,
	}
}

type ListInput struct {
	Offset int
	Limit  int
}

func (ar *AirportRepository) List(input ListInput) ([]database.Airport, error) {
	var airports []database.Airport
	tx := ar.db.DB.
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&airports)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to read airports: %w", tx.Error)
	}
	return airports, nil
}

type ListInCityInput struct {
	Pagination ListInput
	CityID     int64
}

func (ar *AirportRepository) ListInCity(input ListInCityInput) ([]database.Airport, error) {
	var airports []database.Airport
	tx := ar.db.DB.
		Limit(input.Pagination.Limit).
		Offset(input.Pagination.Offset).
		Where("city_id = ?", input.CityID).
		Find(&airports)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to read airports in the city: %w", tx.Error)
	}
	return airports, nil
}

func (ar *AirportRepository) FindByID(id int64) (database.Airport, error) {
	var airport database.Airport
	tx := ar.db.DB.
		Where("id = ?", id).
		Take(&airport)
	if tx.Error != nil {
		return database.Airport{}, fmt.Errorf("failed to find airport: %w", tx.Error)
	}
	return airport, nil
}

func (ar *AirportRepository) SaveAirport(airport database.Airport) (database.Airport, error) {
	tx := ar.db.DB.Begin()
	if tx.Error != nil {
		return database.Airport{}, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()
	tx.Create(&airport)
	if tx.Error != nil {
		return database.Airport{}, fmt.Errorf("failed to save airport: %w", tx.Error)
	}
	var lastInsertRowID int64
	if err := tx.Raw("SELECT last_insert_rowid()").Row().Scan(&lastInsertRowID); err != nil {
		return database.Airport{}, fmt.Errorf("failed to get ID of saved airport: %w", err)
	}
	airport.ID = lastInsertRowID
	tx.Commit()
	if tx.Error != nil {
		return database.Airport{}, fmt.Errorf("failed to commit transaction: %w", tx.Error)
	}
	return airport, nil
}

func (ar *AirportRepository) UpdateAirport(airport database.Airport) error {
	tx := ar.db.DB.Where("id = ?", airport.ID).
		Updates(map[string]interface{}{
			"city_id": airport.CityID,
			"name":    airport.Name,
		})
	if tx.Error != nil {
		return fmt.Errorf("failed to update: %w", tx.Error)
	}
	return nil
}

func (ar *AirportRepository) DeleteByID(id int64) error {
	tx := ar.db.DB.Delete(&database.Airport{}, id)
	switch {
	case tx.Error == gorm.ErrRecordNotFound || tx.RowsAffected == 0:
		return database.ErrNotFound
	case tx.Error != nil:
		return fmt.Errorf("failed to delete: %w", tx.Error)
	default:
		return nil
	}
}

func (ar *AirportRepository) DeleteByCityID(cityID int64) error {
	tx := ar.db.DB.Where("city_id = ?", cityID).
		Delete(&database.Airport{})
	if tx.Error != nil {
		return fmt.Errorf("failed to delete for city: %w", tx.Error)
	}
	return nil
}
