package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

// GetAllCities - list all cities
func GetAllCities(maxComments int) ([]CityDto, error) {
	// count all cities to know how big to make city array
	count, err := countAllCities()
	if err != nil {
		return nil, err
	}
	cities := make([]City, count)
	if e := gdb.Find(&cities).Error; e != nil {
		log.Printf("Error while reading all cities! Error: %s\n", e.Error())
		return nil, &StatementError{
			Message: "Error while reading all cities!",
		}
	}
	result := make([]CityDto, count)
	for i := 0; i < count; i++ {
		comments, err := getCommentsForCity(cities[i].ID, maxComments)
		if err != nil {
			return nil, err
		}
		result[i] = CityDto{
			ID:       cities[i].ID,
			Name:     cities[i].Name,
			Country:  cities[i].Country,
			Comments: comments,
		}
	}
	return result, nil
}

func countAllCities() (int, error) {
	var count int
	if err := gdb.Model(&City{}).Count(&count).Error; err != nil {
		log.Printf("Error while counting cities! Error: %s\n", err.Error())
		return 0, &StatementError{
			Message: "Error while counting cities!",
		}
	}
	return count, nil
}

// GetCityByID - get city by ID
func GetCityByID(id int64, maxComments int) (CityDto, bool, error) {
	city := City{}
	currDB := gdb.First(&city, id)
	if currDB.Error != nil {
		if currDB.RecordNotFound() {
			return CityDto{}, false, &NotFoundError{
				Message: fmt.Sprintf("City with ID %d not found!", id),
			}
		}
		log.Printf("Error while reading city! Error: %s\n", currDB.Error.Error())
		return CityDto{}, false, &StatementError{
			Message: "Error while reading city!",
		}
	}
	comments, err := getCommentsForCity(id, maxComments)
	if err != nil {
		return CityDto{}, true, err
	}
	result := CityDto{
		ID:       city.ID,
		Name:     city.Name,
		Country:  city.Country,
		Comments: comments,
	}
	return result, true, nil
}

func getCityByNameAndCountry(name string, country string) (CityDto, error) {
	city := City{}
	currDB := gdb.Where("LOWER(name) = LOWER(?) AND LOWER(country) = LOWER(?)", name, country).First(&city)
	if currDB.Error != nil {
		if currDB.RecordNotFound() {
			return CityDto{}, &NotFoundError{
				Message: fmt.Sprintf("City with name %s in country %s not found!", name, country),
			}
		}
		log.Printf("Error while reading city! Error: %s\n", currDB.Error.Error())
		return CityDto{}, &StatementError{
			Message: "Error while reading city!",
		}
	}
	result := CityDto{
		ID:      city.ID,
		Name:    city.Name,
		Country: city.Country,
	}
	return result, nil
}

// AddNewCity - save new city
func AddNewCity(name string, country string) error {
	_, err := getCityByNameAndCountry(name, country)
	if err == nil {
		return &ForbidenError{
			Message: fmt.Sprintf("City with name %s in country %s is already saved!", name, country),
		}
	} else if !errors.As(err, &NotFoundError{}) {
		return err
	}
	city := City{
		Name:    name,
		Country: country,
	}
	if e := gdb.Create(&city).Error; e != nil {
		log.Printf("Error while saving city! Error: %s\n", e.Error())
		return &StatementError{
			Message: "Error while saving city!",
		}
	}
	return nil
}

// UpdateCity - updates a city
func UpdateCity(id int64, name string, country string) error {
	city := City{
		ID:      id,
		Name:    name,
		Country: country,
	}
	if gdb.NewRecord(&city) {
		return &NotFoundError{
			Message: fmt.Sprintf("City with ID %d not found!", id),
		}
	}
	if e := gdb.Save(&city).Error; e != nil {
		log.Printf("Error while updating city! Error: %s\n", e.Error())
		return &StatementError{
			Message: "Error while updating city!",
		}
	}
	return nil
}

// DeleteCity - delete a city with given ID
func DeleteCity(id int64) error {
	return gdb.Transaction(func(tx *gorm.DB) error {
		err := deleteCommentsForCity(id, tx)
		if err != nil {
			return err
		}
		city := City{ID: id}
		curdb := tx.Delete(&city)
		if curdb.RowsAffected == 0 {
			return &NotFoundError{
				Message: fmt.Sprintf("City with ID %d not found!", id),
			}
		}
		err = curdb.Error
		if err != nil {
			log.Printf("Error while deleting city! Error: %s\n", err.Error())
			return &StatementError{
				Message: "Error while deleting city!",
			}
		}
		return nil
	})
}
