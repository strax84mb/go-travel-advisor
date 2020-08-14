package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.strale.io/go-travel/models"
)

// SaveAirport - save new airport to DB
func SaveAirport(airportID int64, name string, cityID int64) (AirportDto, error) {
	city, err := GetCityByID(cityID, 0)
	if err != nil {
		return AirportDto{}, err
	}
	_, err = loadAirportByAirportID(airportID)
	if err == nil {
		return AirportDto{}, &ForbidenError{
			Message: fmt.Sprintf("Airport with AirportID %d is already saved", airportID),
		}
	}
	var nfe *NotFoundError
	if !errors.As(err, &nfe) {
		return AirportDto{}, &StatementError{
			Message: fmt.Sprintf("Error while checking if AirportID %d is already taken", airportID),
			Cause:   err,
		}
	}
	airport := &models.Airport{
		CityID:    cityID,
		Name:      name,
		AirportID: airportID,
	}
	if err = airport.Insert(context.Background(), db, boil.Infer()); err != nil {
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
	city, err := GetCityByID(cityID, 0)
	if err != nil {
		return AirportDto{}, err
	}
	airportWithAirportID, err := loadAirportByAirportID(airportID)
	if err != nil {
		var nfe *NotFoundError
		if !errors.As(err, &nfe) {
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
	if _, err = airport.Update(context.Background(), db, boil.Infer()); err != nil {
		return AirportDto{}, &StatementError{
			Message: "Error while updating airport data",
			Cause:   err,
		}
	}
	return makeAirportDto(airport, city), nil
}

// DeleteAirport - delete airport from DB
func DeleteAirport(id int64) error {
	rowsAff, err := models.Airports(models.AirportWhere.ID.EQ(id)).DeleteAll(context.Background(), db)
	if err != nil {
		log.Printf("Error while deleting airport with ID %d! Error: %s", id, err.Error())
		return &StatementError{
			Message: fmt.Sprintf("Error while deleting airport with ID %d", id),
			Cause:   err,
		}
	} else if rowsAff == 0 {
		return airportNotFoundError(id)
	}
	return nil
}

func loadAirportByAirportID(airportID int64) (*models.Airport, error) {
	airport, err := models.Airports(models.AirportWhere.AirportID.EQ(airportID)).One(context.Background(), db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{
				Message: fmt.Sprintf("Airport with AirportID %d not found", airportID),
			}
		}
		log.Printf("Error while loading airport with AirportID %d! Error: %s", airportID, err.Error())
		return nil, &StatementError{
			Message: fmt.Sprintf("Error while loading airport with AirportID %d", airportID),
			Cause:   err,
		}
	}
	return airport, nil
}

func loadAirport(id int64) (*models.Airport, error) {
	airport, err := models.Airports(models.AirportWhere.ID.EQ(id)).One(context.Background(), db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, airportNotFoundError(id)
		}
		log.Printf("Error while loading airport with ID %d! Error: %s", id, err.Error())
		return nil, &StatementError{
			Message: fmt.Sprintf("Error while loading airport with ID %d", id),
			Cause:   err,
		}
	}
	return airport, nil
}

func makeAirportDto(airport *models.Airport, city CityDto) AirportDto {
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
	city, err := GetCityByID(airport.CityID, 0)
	if err != nil {
		return AirportDto{}, err
	}
	return makeAirportDto(airport, city), nil
}

// ListAirports - load all airports from DB
func ListAirports() ([]AirportDto, error) {
	airports, err := models.Airports(qm.Load(models.AirportRels.City)).All(context.Background(), db)
	if err != nil {
		return nil, &StatementError{
			Message: "Error while loading all airports",
			Cause:   err,
		}
	}
	result := make([]AirportDto, len(airports))
	for i, a := range airports {
		result[i] = AirportDto{
			ID:        a.ID,
			AirportID: a.AirportID,
			Name:      a.Name,
			City: CityDto{
				ID:      a.CityID,
				Name:    a.R.City.Name,
				Country: a.R.City.Country,
			},
		}
	}
	return result, nil
}

// ImportSingleAirport - import single airport data
func ImportSingleAirport(airportID int64, name string, cityName string, cityCountry string) error {
	city, err := getCityByNameAndCountry(cityName, cityCountry)
	if err != nil {
		return err
	}
	_, err = SaveAirport(airportID, name, city.ID)
	return err
}
