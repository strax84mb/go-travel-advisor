package airports

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.strale.io/go-travel/internal/airports/repository"
	"gitlab.strale.io/go-travel/internal/database"
)

type airportRepository interface {
	List(input repository.ListInput) ([]database.Airport, error)
	ListInCity(input repository.ListInCityInput) ([]database.Airport, error)
	FindByID(id int64) (database.Airport, error)
	SaveAirport(airport database.Airport) (database.Airport, error)
	UpdateAirport(airport database.Airport) error
	DeleteByID(id int64) error
	DeleteByCityID(cityID int64) error
}

type cityRepository interface {
	FindByID(id int64, preload bool) (database.City, error)
	FindByName(name string, preload bool) (database.City, error)
}

type AirportService struct {
	airportRepo airportRepository
	cityRepo    cityRepository
}

func NewAirportService(airportRepo airportRepository, cityRepo cityRepository) *AirportService {
	return &AirportService{
		airportRepo: airportRepo,
		cityRepo:    cityRepo,
	}
}

func (as *AirportService) ListAirports(ctx context.Context, limit int, offset int) ([]database.Airport, error) {
	airports, err := as.airportRepo.List(repository.ListInput{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("error while listing airports")
		return nil, fmt.Errorf("error while listing airports: %w", err)
	}
	return airports, err
}

func (as *AirportService) ListAirportsInCity(ctx context.Context, cityID int64, limit, offset int) ([]database.Airport, error) {
	airports, err := as.airportRepo.ListInCity(repository.ListInCityInput{
		CityID: cityID,
		Pagination: repository.ListInput{
			Offset: offset,
			Limit:  limit,
		},
	})
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithField("cityId", cityID).
			Error("error while listing airports in city")
		return nil, fmt.Errorf("error while listing airports in city: %w", err)
	}
	return airports, err
}

func (as *AirportService) FindByID(ctx context.Context, id int64) (database.Airport, error) {
	airport, err := as.airportRepo.FindByID(id)
	if err != nil && err != database.ErrNotFound {
		logrus.WithContext(ctx).WithError(err).
			WithField("airportId", id).
			Error("error loading airport")
		return database.Airport{}, fmt.Errorf("error loading airport: %w", err)
	}
	return airport, nil
}

func (as *AirportService) SaveNewAirport(ctx context.Context, airport database.Airport) (database.Airport, error) {
	_, err := as.cityRepo.FindByID(airport.CityID, false)
	switch {
	case err == database.ErrNotFound:
		logrus.WithContext(ctx).WithError(err).
			WithField("cityId", airport.CityID).
			Error("city not found")
		return database.Airport{}, fmt.Errorf("city not found: %w", err)
	case err != nil:
		logrus.WithContext(ctx).WithError(err).
			WithField("cityId", airport.CityID).
			Error("failed to check if city exists")
		return database.Airport{}, fmt.Errorf("failed to check if city exists: %w", err)
	}
	savedAirport, err := as.airportRepo.SaveAirport(airport)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			Error("failed to save airport")
		return database.Airport{}, fmt.Errorf("failed to save airport: %w", err)
	}
	return savedAirport, nil
}

func (as *AirportService) UpdateAirport(ctx context.Context, airport database.Airport) error {
	_, err := as.cityRepo.FindByID(airport.CityID, false)
	switch {
	case err == database.ErrNotFound:
		logrus.WithContext(ctx).WithError(err).
			WithField("cityId", airport.CityID).
			Error("city not found")
		return fmt.Errorf("city not found: %w", err)
	case err != nil:
		logrus.WithContext(ctx).WithError(err).
			WithField("cityId", airport.CityID).
			Error("failed to check if city exists")
		return fmt.Errorf("failed to check if city exists: %w", err)
	}
	if err = as.airportRepo.UpdateAirport(airport); err != nil {
		logrus.WithContext(ctx).WithError(err).
			Error("failed to update airport")
		return fmt.Errorf("failed to update airport: %w", err)
	}
	return nil
}

func (as *AirportService) DeleteAirport(ctx context.Context, id int64) error {
	if err := as.airportRepo.DeleteByID(id); err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithField("id", id).
			Error("failed to delete airport")
		return fmt.Errorf("failed to delete airport: %w", err)
	}
	return nil
}

func (as *AirportService) DeleteAirportsInCity(ctx context.Context, cityID int64) error {
	if err := as.airportRepo.DeleteByCityID(cityID); err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithField("cityId", cityID).
			Error("failed to delete all airports in the city")
		return fmt.Errorf("failed to delete all airports in the city: %w", err)
	}
	return nil
}

func (as *AirportService) ImportAirports(ctx context.Context, content []byte) error {
	reader := bytes.NewReader(content)
	csvReader := csv.NewReader(reader)
	var (
		fields []string
		city   database.City
		err    error
	)
	for {
		fields, err = csvReader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			logrus.WithContext(ctx).WithError(err).
				Error("failed to read city row")
			continue
		}
		city, err = as.cityRepo.FindByName(fields[1], false)
		if err != nil {
			logrus.WithContext(ctx).WithError(err).
				WithField("cityName", fields[1]).
				Error("error while reading city by name")
		}
		_, err = as.airportRepo.SaveAirport(database.Airport{
			Name:   fields[0],
			CityID: city.ID,
		})
		if err != nil {
			logrus.WithContext(ctx).WithError(err).
				WithField("csv_row", strings.Join(fields, ",")).
				Error("failed to import city")
		}
	}
}
