package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"gitlab.strale.io/go-travel/common"
)

// GetAllCities - list all cities
func GetAllCities(maxComments int) ([]CityDto, *common.GeneralError) {
	// count all cities to know how big to make city array
	count, err := countAllCities()
	if err != nil {
		return nil, err
	}
	cities := make([]City, count)
	if e := gdb.Find(&cities).Error; e != nil {
		log.Printf("Error while reading all cities! Error: %s\n", e.Error())
		return nil, &common.GeneralError{
			Message:  "Error while reading all cities!",
			Location: "database.cities.GetAllCities",
			Cause:    e,
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

func countAllCities() (int, *common.GeneralError) {
	var count int
	if err := gdb.Model(&City{}).Count(&count).Error; err != nil {
		log.Printf("Error while counting cities! Error: %s\n", err.Error())
		return 0, &common.GeneralError{
			Message:  "Error while counting cities!",
			Location: "database.cities.countAllCities",
			Cause:    err,
		}
	}
	return count, nil
}

// GetCityByID - get city by ID
func GetCityByID(id int64, maxComments int) (CityDto, bool, *common.GeneralError) {
	city := City{}
	currDB := gdb.First(&city, id)
	if currDB.Error != nil {
		if currDB.RecordNotFound() {
			return CityDto{}, false, &common.GeneralError{
				Message:   fmt.Sprintf("City with ID %d not found!", id),
				Location:  "database.cities.GetCityByID",
				ErrorType: common.CityNotFound,
			}
		}
		log.Printf("Error while reading city! Error: %s\n", currDB.Error.Error())
		return CityDto{}, false, &common.GeneralError{
			Message:  "Error while reading city!",
			Location: "database.cities.GetCityByID",
			Cause:    currDB.Error,
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

func getCityByNameAndCountry(name string, country string) (CityDto, *common.GeneralError) {
	city := City{}
	currDB := gdb.Where("LOWER(name) = LOWER(?) AND LOWER(country) = LOWER(?)", name, country).First(&city)
	if currDB.Error != nil {
		if currDB.RecordNotFound() {
			return CityDto{}, &common.GeneralError{
				Message:   fmt.Sprintf("City with name %s in country %s not found!", name, country),
				Location:  "database.cities.getCityByNameAndCountry",
				ErrorType: common.CityNotFound,
			}
		}
		log.Printf("Error while reading city! Error: %s\n", currDB.Error.Error())
		return CityDto{}, &common.GeneralError{
			Message:  "Error while reading city!",
			Location: "database.cities.getCityByNameAndCountry",
			Cause:    currDB.Error,
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
func AddNewCity(name string, country string) *common.GeneralError {
	_, err := getCityByNameAndCountry(name, country)
	if err != nil && err.ErrorType != common.CityNotFound {
		return err
	}
	city := City{
		Name:    name,
		Country: country,
	}
	if e := gdb.Create(&city).Error; e != nil {
		log.Printf("Error while saving city! Error: %s\n", e.Error())
		return &common.GeneralError{
			Message:  "Error while saving city!",
			Location: "database.cities.AddNewCity",
			Cause:    e,
		}
	}
	return nil
}

// UpdateCity - updates a city
func UpdateCity(id int64, name string, country string) *common.GeneralError {
	city := City{
		ID:      id,
		Name:    name,
		Country: country,
	}
	if gdb.NewRecord(&city) {
		return &common.GeneralError{
			Message:   fmt.Sprintf("City with ID %d not found!", id),
			Location:  "database.cities.UpdateCity",
			ErrorType: common.CityNotFound,
		}
	}
	if e := gdb.Save(&city).Error; e != nil {
		log.Printf("Error while updating city! Error: %s\n", e.Error())
		return &common.GeneralError{
			Message:  "Error while updating city!",
			Location: "database.cities.UpdateCity",
			Cause:    e,
		}
	}
	return nil
}

// DeleteCity - delete a city with given ID
func DeleteCity(id int64) *common.GeneralError {
	err := gdb.Transaction(func(tx *gorm.DB) error {
		err := deleteCommentsForCity(id, tx)
		if err != nil {
			log.Printf("Error while deleting comments of city! Error: %s\n", err.Error())
			return err
		}
		city := City{ID: id}
		err = tx.Delete(&city).Error
		if err != nil {
			log.Printf("Error while reading city! Error: %s\n", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("Error while deleting city! Error: %s\n", err.Error())
		return &common.GeneralError{
			Message:  "Error while deleting city!",
			Location: "database.cities.DeleteCity",
			Cause:    err,
		}
	}
	return nil
}
