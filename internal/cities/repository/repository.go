package repository

import (
	"fmt"

	"gitlab.strale.io/go-travel/internal/database"
	"gorm.io/gorm"
)

type CityRepository struct {
	db *database.Database
}

func NewCityRepository(db *database.Database) *CityRepository {
	return &CityRepository{
		db: db,
	}
}

type FindInput struct {
	Offset int
	Limit  int
}

func (cr *CityRepository) Find(input FindInput) ([]database.City, error) {
	var result []database.City
	tx := cr.db.DB.
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&result)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to read cities: %w", tx.Error)
	}
	return result, nil
}

func (cr *CityRepository) FindByIDs(ids []int64) ([]database.City, error) {
	var cities []database.City
	tx := cr.db.DB.Where("id IN (?)", ids).Find(&cities)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to get cities for IDs: %w", tx.Error)
	}
	return cities, nil
}

func (cr *CityRepository) findOne(where func(*gorm.DB) *gorm.DB, preload bool) (database.City, error) {
	city := database.City{}
	tx := cr.db.DB
	if preload {
		tx = tx.
			Preload("Airports").
			Preload("Comments")
	}
	tx = where(tx).Take(&city)
	if tx.Error == gorm.ErrRecordNotFound {
		return database.City{}, database.ErrNotFound
	} else if tx.Error != nil {
		return database.City{}, fmt.Errorf("failed to read city: %w", tx.Error)
	}
	return city, nil
}

func (cr *CityRepository) FindByID(id int64, preload bool) (database.City, error) {
	return cr.findOne(
		func(tx *gorm.DB) *gorm.DB {
			return tx.Where("id = ?", id)
		},
		false,
	)
}

func (cr *CityRepository) FindByName(name string, preload bool) (database.City, error) {
	return cr.findOne(
		func(tx *gorm.DB) *gorm.DB {
			return tx.Where("name = ?", name)
		},
		false,
	)
}

func (cr *CityRepository) SaveNew(city database.City) (database.City, error) {
	tx := cr.db.DB.Begin()
	if tx.Error != nil {
		return database.City{}, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()
	city.Airports = nil
	city.Comments = nil
	tx.Create(&city)
	if tx.Error != nil {
		return database.City{}, fmt.Errorf("failed to insert city: %w", tx.Error)
	}
	var lastRowID int64
	err := tx.Raw("SELECT last_insert_rowid()").Row().Scan(&lastRowID)
	if err != nil {
		return database.City{}, fmt.Errorf("failed to get last inserted ID: %w", err)
	}
	city.ID = lastRowID
	if tx.Commit(); tx.Error != nil {
		return database.City{}, fmt.Errorf("failed to commit transaction: %w", tx.Error)
	}
	return city, nil
}

func (cr *CityRepository) Update(city database.City) error {
	tx := cr.db.DB.Where("id = ?", city.ID).
		Update("name", city.Name)
	switch {
	case tx.Error == gorm.ErrRecordNotFound:
		return database.ErrNotFound
	case tx.RowsAffected == 0:
		return database.ErrNotFound
	case tx.Error != nil:
		return fmt.Errorf("failed to update city with ID=%d: %w", city.ID, tx.Error)
	default:
		return nil
	}
}

func (cr *CityRepository) Delete(id int64) error {
	tx := cr.db.DB.Delete(&database.City{}, id)
	if tx.Error == gorm.ErrRecordNotFound {
		return database.ErrNotFound
	} else if tx.Error != nil {
		return fmt.Errorf("failed to delete city with ID=%d: %w", id, tx.Error)
	}
	return nil
}
