package database

import (
	"database/sql"

	"gitlab.strale.io/go-travel/common"
)

// GetAllCities - list all cities
func GetAllCities(maxCommants int) ([]CityDto, error) {
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
			comments, err := getCommentsForCity(city.ID, maxCommants)
			if err != nil {
				return err
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

func countAllCities() (int, error) {
	query := `SELECT COUNT(city.id) FROM city`
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return 0, &common.GeneralError{
			Message:  "Error while reading city by ID!",
			Location: "database.cities.GetAllCities",
			Cause:    err,
		}
	}
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, &common.GeneralError{
				Message:  "Error while reading columns",
				Location: "database.cities.countAllCities",
				Cause:    err,
			}
		}
	}
	return count, nil
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
			comments, err := getCommentsForCity(id, maxComments)
			if err != nil {
				return nil, err
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

func getCityByNameAndCountry(name string, country string) (CityDto, error) {
	query := `SELECT id FROM city WHERE name = $1 AND country = $2`
	rows, err := db.Query(query, name, country)
	if err != nil {
		return CityDto{}, &common.GeneralError{
			Message:  "Error while querying for city by name and country",
			Location: "database.cities.getCityByNameAndCountry",
			Cause:    err,
		}
	}
	defer rows.Close()
	city := CityDto{
		Name:    name,
		Country: country,
	}
	for rows.Next() {
		err = rows.Scan(&city.ID)
		if err != nil {
			return CityDto{}, &common.GeneralError{
				Message:  "Error while reading columns",
				Location: "database.cities.getCityByNameAndCountry",
				Cause:    err,
			}
		}
	}
	return city, nil
}

// AddNewCity - save new city
func AddNewCity(name string, country string) error {
	city, err := getCityByNameAndCountry(name, country)
	if err != nil {
		return &common.GeneralError{
			Message:  "Error while checking if such city already exists",
			Location: "database.cities.AddNewCity",
			Cause:    err,
		}
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
func UpdateCity(id int64, name string, country string) error {
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
func DeleteCity(id int64) error {
	err := deleteCommentsForCity(id)
	if err != nil {
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
