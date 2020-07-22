package database

import (
	"database/sql"
	"fmt"

	"gitlab.strale.io/go-travel/common"
)

// GetAllCities - list all cities
func GetAllCities(maxComments int) ([]CityDto, *common.GeneralError) {
	// count all cities to know how big to make city array
	count, err := countAllCities()
	if err != nil {
		return nil, &common.GeneralError{
			Message:  "Error while counting all the cities",
			Location: "database.cities.GetAllCities",
			Cause:    err,
		}
	}
	cities := make([]CityDto, count)
	err = performListSelection(
		"database.cities.GetAllCities", count, cities[:], getAllCitiesSelection,
		func(rows *sql.Rows, array interface{}, index int) error {
			cities := array.([]CityDto)
			city := CityDto{}
			err := rows.Scan(&city.ID, &city.Name, &city.Country)
			if err != nil {
				return err
			}
			comments, generalError := getCommentsForCity(city.ID, maxComments)
			if generalError != nil {
				return generalError
			}
			city.Comments = comments
			cities[index] = city
			return nil

		})
	return cities, err
}

func getAllCitiesSelection(_ []interface{}) (*sql.Rows, error) {
	query := `SELECT city.id, city.name, city.country FROM city`
	return db.Query(query)
}

func getAllCitiesConversion(rows *sql.Rows, array interface{}, index int) error {
	cities := array.([]CityDto)
	city := CityDto{}
	err := rows.Scan(&city.ID, &city.Name, &city.Country)
	if err == nil {
		cities[index] = city
	}
	return err
}

func countAllCities() (int, *common.GeneralError) {
	value, _, err := performSingleSelection(
		"database.cities.countAllCities",
		func(_ []interface{}) (*sql.Rows, error) {
			query := `SELECT COUNT(city.id) FROM city`
			return db.Query(query)
		},
		func(rows *sql.Rows) (interface{}, error) {
			var count int
			err := rows.Scan(&count)
			return count, err
		})
	if err != nil {
		return 0, err
	}
	return value.(int), nil
}

// GetCityByID - get city by ID
func GetCityByID(id int64, maxComments int) (CityDto, bool, *common.GeneralError) {
	result, found, err := performSingleSelection(
		"database.cities.GetCityByID",
		func(params []interface{}) (*sql.Rows, error) {
			id := params[0].(int64)
			query := `SELECT city.name, city.country FROM city WHERE city.id = $1`
			return db.Query(query, id)
		},
		func(rows *sql.Rows) (interface{}, error) {
			city := CityDto{
				ID: id,
			}
			err := rows.Scan(&city.Name, &city.Country)
			if err != nil {
				return nil, err
			}
			comments, generalError := getCommentsForCity(id, maxComments)
			if generalError != nil {
				return nil, generalError
			}
			city.Comments = comments
			return city, nil
		},
		id)
	if err != nil {
		return CityDto{}, false, err
	}
	if !found {
		return CityDto{}, false, nil
	}
	return result.(CityDto), true, nil
}

func getCityByNameAndCountry(name string, country string) (CityDto, *common.GeneralError) {
	value, found, err := performSingleSelection(
		"database.cities.getCityByNameAndCountry",
		func(_ []interface{}) (*sql.Rows, error) {
			query := `SELECT id FROM city WHERE name = $1 AND country = $2`
			return db.Query(query, name, country)
		},
		func(rows *sql.Rows) (interface{}, error) {
			city := CityDto{
				Name:    name,
				Country: country,
			}
			err := rows.Scan(&city.ID)
			return city, err
		})
	if err != nil {
		return CityDto{}, err
	}
	if !found {
		return CityDto{}, &common.GeneralError{
			Message:   fmt.Sprintf("City with name %s and country %s not found!", name, country),
			Location:  "database.cities.getCityByNameAndCountry",
			ErrorType: common.CityNotFound,
		}
	}
	return value.(CityDto), nil
}

// AddNewCity - save new city
func AddNewCity(name string, country string) *common.GeneralError {
	city, err := getCityByNameAndCountry(name, country)
	if err != nil && err.ErrorType != common.CityNotFound {
		return err
	}
	if city.ID != 0 {
		return &common.GeneralError{
			Message:  "City already exists",
			Location: "database.cities.AddNewCity",
		}
	}
	return performStatement("database.cities.AddNewCity", executeAddNewCity, name, country)
}

func executeAddNewCity(params []interface{}) (sql.Result, error) {
	name := params[0].(string)
	country := params[1].(string)
	statement := `INSERT INTO city (name, country) VALUES ($1, $2)`
	return db.Exec(statement, name, country)
}

// UpdateCity - updates a city
func UpdateCity(id int64, name string, country string) *common.GeneralError {
	return performStatement("database.cities.UpdateCity",
		func(params []interface{}) (sql.Result, error) {
			id := params[0].(int64)
			name := params[1].(string)
			country := params[2].(string)
			statement := `UPDATE city SET name = $1, country = $2 WHERE id = $3`
			return db.Exec(statement, name, country, id)
		},
		id, name, country)
}

// DeleteCity - delete a city with given ID
func DeleteCity(id int64) *common.GeneralError {
	err := deleteCommentsForCity(id)
	if err != nil && err.ErrorType != common.NoRowsAffected {
		return err
	}
	return performStatement("database.cities.DeleteCity",
		func(params []interface{}) (sql.Result, error) {
			id := params[0].(int64)
			statement := `DELETE FROM city WHERE id = $1`
			return db.Exec(statement, id)
		},
		id)
}
