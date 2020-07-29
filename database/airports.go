package database

import (
	"errors"
	"fmt"
	"log"
)

// SaveAirport - save new airport to DB
func SaveAirport(airportID int64, name string, cityID int64) (AirportDto, error) {
	city, _, err := GetCityByID(cityID, 0)
	if err != nil {
		return AirportDto{}, err
	}
	_, err = loadAirportByAirportID(airportID)
	if err == nil {
		return AirportDto{}, &ForbidenError{
			Message: fmt.Sprintf("Airport with AirportID %d is already saved", airportID),
		}
	}
	if !errors.As(err, &NotFoundError{}) {
		return AirportDto{}, &StatementError{
			Message: fmt.Sprintf("Error while checking if AirportID %d is already taken", airportID),
			Cause:   err,
		}
	}
	airport := Airport{
		CityID:    cityID,
		Name:      name,
		AirportID: airportID,
	}
	if err = gdb.Create(&airport).Error; err != nil {
		return AirportDto{}, &StatementError{
			Message: "Error while saving airport",
			Cause:   err,
		}
	}
	return makeAirportDto(airport, city), nil
}

func airportNotFoundError(id int64) error {
	return &NotFoundError{fmt.Sprintf("Airport with ID %d not found!", id)}
}

// UpdateAirport - change airport data in DB
func UpdateAirport(id int64, airportID int64, name string, cityID int64) (AirportDto, error) {
	airport, err := loadAirport(id)
	if err != nil {
		return AirportDto{}, nil
	}
	city, _, err := GetCityByID(cityID, 0)
	if err != nil {
		return AirportDto{}, err
	}
	airportWithAirportID, err := loadAirportByAirportID(airportID)
	if err != nil {
		if !errors.As(err, &NotFoundError{}) {
			return AirportDto{}, &StatementError{
				Message: fmt.Sprintf("Error while checking if AirportID %d is already taken by another airport", airportID),
				Cause:   err,
			}
		}
	} else if airportWithAirportID.ID != airport.ID {
		return AirportDto{}, &ForbidenError{
			Message: fmt.Sprintf("AirportID %d is already taken by another airport", airportID),
		}
	}
	airport.AirportID = airportID
	airport.Name = name
	airport.CityID = cityID
	airport.City = City{}
	if err = gdb.Save(&airport).Error; err != nil {
		return AirportDto{}, &StatementError{
			Message: "Error while updating airport data",
			Cause:   err,
		}
	}
	return makeAirportDto(airport, city), nil
}

// DeleteAirport - delete airport from DB
func DeleteAirport(id int64) error {
	curDB := gdb.Delete(&Airport{ID: id})
	if curDB.RecordNotFound() {
		return airportNotFoundError(id)
	} else if curDB.Error != nil {
		log.Printf("Error while deleting airport with ID %d! Error: %s", id, curDB.Error.Error())
		return &StatementError{
			Message: fmt.Sprintf("Error while deleting airport with ID %d", id),
			Cause:   gdb.Error,
		}
	}
	return nil
}

func loadAirportByAirportID(airportID int64) (Airport, error) {
	airport := Airport{}
	curDB := gdb.Where(&Airport{AirportID: airportID}).First(&airport)
	if curDB.RecordNotFound() {
		return Airport{}, &NotFoundError{Message: fmt.Sprintf("Airport with AirportID %d not found", airportID)}
	} else if curDB.Error != nil {
		log.Printf("Error while loading airport with AirportID %d! Error: %s", airportID, curDB.Error.Error())
		return Airport{}, &StatementError{
			Message: fmt.Sprintf("Error while loading airport with AirportID %d", airportID),
			Cause:   curDB.Error,
		}
	}
	return airport, nil
}

func loadAirport(id int64) (Airport, error) {
	airport := Airport{}
	curDB := gdb.First(&airport, id)
	if curDB.RecordNotFound() {
		return Airport{}, airportNotFoundError(id)
	} else if curDB.Error != nil {
		log.Printf("Error while loading airport with ID %d! Error: %s", id, curDB.Error.Error())
		return Airport{}, &StatementError{
			Message: fmt.Sprintf("Error while loading airport with ID %d", id),
			Cause:   curDB.Error,
		}
	}
	return airport, nil
}

func makeAirportDto(airport Airport, city CityDto) AirportDto {
	return AirportDto{
		ID:        airport.ID,
		AirportID: airport.AirportID,
		Name:      airport.Name,
		City:      city,
	}
}

// GetAirport - read airport from DB
func GetAirport(id int64) (AirportDto, error) {
	airport, err := loadAirport(id)
	if err != nil {
		return AirportDto{}, err
	}
	city, _, err := GetCityByID(airport.CityID, 0)
	if err != nil {
		return AirportDto{}, err
	}
	return makeAirportDto(airport, city), nil
}

func countAllAirports() (int, error) {
	var count int
	if err := gdb.Model(&Airport{}).Count(&count).Error; err != nil {
		log.Printf("Error while counting all airports! Error: %s", err.Error())
		return 0, &StatementError{
			Message: "Error while counting all airports",
			Cause:   err,
		}
	}
	return count, nil
}

// ListAirports - load all airports from DB
func ListAirports() ([]AirportDto, error) {
	count, err := countAllAirports()
	if err != nil {
		return nil, err
	}
	airports := make([]Airport, count)
	if err := gdb.Preload("City").Find(&airports).Error; err != nil {
		return nil, &StatementError{
			Message: "Error while loading all airports",
			Cause:   err,
		}
	}
	result := make([]AirportDto, count)
	for i, a := range airports {
		result[i] = AirportDto{
			ID:        a.ID,
			AirportID: a.AirportID,
			Name:      a.Name,
			City: CityDto{
				ID:      a.CityID,
				Name:    a.City.Name,
				Country: a.City.Country,
			},
		}
	}
	return result, nil
}
