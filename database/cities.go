package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"gitlab.strale.io/go-travel/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetAllCities - list all cities
func GetAllCities(maxComments int) ([]CityDto, error) {
	cities, err := models.Cities().All(context.Background(), db)
	if err != nil {
		log.Printf("Error while reading all cities! Error: %s\n", err.Error())
		return nil, &StatementError{
			Message: "Error while reading all cities!",
		}
	}
	result := make([]CityDto, len(cities))
	for i, city := range cities {
		comments, err := getCommentsForCity(city.ID, maxComments)
		if err != nil {
			return nil, err
		}
		result[i] = CityDto{
			ID:       city.ID,
			Name:     city.Name,
			Country:  city.Country,
			Comments: comments,
		}
	}
	return result, nil
}

// GetCityByID - get city by ID
func GetCityByID(id int64, maxComments int) (CityDto, error) {
	city, err := models.Cities(models.CityWhere.ID.EQ(id)).One(context.Background(), db)
	if err != nil {
		if err == sql.ErrNoRows {
			return CityDto{}, &NotFoundError{
				Message: fmt.Sprintf("City with id %d not found!", id),
			}
		}
		log.Printf("Error while reading city! Error: %s\n", err.Error())
		return CityDto{}, &StatementError{
			Message: "Error while reading city!",
		}
	}
	comments, err := getCommentsForCity(id, maxComments)
	if err != nil {
		return CityDto{}, err
	}
	result := CityDto{
		ID:       city.ID,
		Name:     city.Name,
		Country:  city.Country,
		Comments: comments,
	}
	return result, nil
}

func getCityByNameAndCountry(name string, country string) (CityDto, error) {
	city, err := models.Cities(qm.Where("LOWER(name) = LOWER(?)", name), qm.And("LOWER(country) = LOWER(?)", country)).One(context.Background(), db)
	if err != nil {
		if err == sql.ErrNoRows {
			return CityDto{}, &NotFoundError{
				Message: fmt.Sprintf("City with name %s in country %s not found!", name, country),
			}
		}
		log.Printf("Error while reading city! Error: %s\n", err.Error())
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
	var nfe *NotFoundError
	if err == nil {
		return &ForbidenError{
			Message: fmt.Sprintf("City with name %s in country %s is already saved!", name, country),
		}
	} else if !errors.As(err, &nfe) {
		return err
	}
	city := models.City{
		Name:    name,
		Country: country,
	}
	if err = city.Insert(context.Background(), db, boil.Infer()); err != nil {
		log.Printf("Error while saving city! Error: %s\n", err.Error())
		return &StatementError{
			Message: "Error while saving city!",
		}
	}
	return nil
}

// UpdateCity - updates a city
func UpdateCity(id int64, name string, country string) error {
	city, err := models.Cities(models.CityWhere.ID.EQ(id)).One(context.Background(), db)
	if err != nil {
		if err == sql.ErrNoRows {
			return &NotFoundError{
				Message: fmt.Sprintf("City with id %d not found!", id),
			}
		}
		log.Printf("Error while reading city! Error: %s\n", err.Error())
		return &StatementError{
			Message: "Error while reading city!",
		}
	}
	city.Name = name
	city.Country = country
	if _, err = city.Update(context.Background(), db, boil.Infer()); err != nil {
		log.Printf("Error while updating city! Error: %s\n", err.Error())
		return &StatementError{
			Message: "Error while updating city!",
		}
	}
	return nil
}

// DeleteCity - delete a city with given ID
func DeleteCity(id int64) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return &StatementError{
			Message: "Error while starting transaction",
		}
	}
	if err = deleteCityInTransaction(id, tx); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return &StatementError{
			Message: "Error while commiting transaction",
		}
	}
	return nil
}

func deleteCityInTransaction(id int64, tx *sql.Tx) error {
	err := deleteCommentsForCity(id, tx)
	if err != nil {
		return err
	}
	rowsAff, err := models.Cities(models.CityWhere.ID.EQ(id)).DeleteAll(context.Background(), tx)
	if err != nil {
		log.Printf("Error while deleting city! Error: %s\n", err.Error())
		return &StatementError{
			Message: "Error while deleting city!",
		}
	} else if rowsAff == 0 {
		return &NotFoundError{
			Message: fmt.Sprintf("City with ID %d not found!", id),
		}
	}
	return nil
}
