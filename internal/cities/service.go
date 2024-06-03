package cities

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.strale.io/go-travel/internal/cities/repository"
	"gitlab.strale.io/go-travel/internal/database"
)

type cityRepository interface {
	Find(input repository.FindInput) ([]database.City, error)
	FindByID(id int64, preload bool) (database.City, error)
	SaveNew(city database.City) (database.City, error)
	Update(city database.City) error
	Delete(id int64) error
}

type airportRepository interface {
	DeleteByCityID(cityID int64) error
}

type cityService struct {
	cityRepo    cityRepository
	airportRepo airportRepository
}

func NewCityService(cityRepo cityRepository, airportRepo airportRepository) *cityService {
	return &cityService{
		cityRepo:    cityRepo,
		airportRepo: airportRepo,
	}
}

func (cs *cityService) ListCities(ctx context.Context, offset int, limit int) ([]database.City, error) {
	cities, err := cs.cityRepo.Find(repository.FindInput{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		logrus.WithError(err).WithContext(ctx).Error("error listing cities")
		return nil, fmt.Errorf("error listing cities: %w", err)
	}
	return cities, nil
}

func (cs *cityService) FindByID(ctx context.Context, id int64) (database.City, error) {
	city, err := cs.cityRepo.FindByID(id, true)
	if err != nil && err != database.ErrNotFound {
		logrus.WithError(err).WithContext(ctx).Error("error loading city")
		return database.City{}, fmt.Errorf("error loading city: %w", err)
	}
	return city, err
}

func (cs *cityService) SaveNewCity(ctx context.Context, name string) (database.City, error) {
	city, err := cs.cityRepo.SaveNew(database.City{Name: name})
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("error while saving city")
		return database.City{}, fmt.Errorf("error while saving city: %w", err)
	}
	return city, nil
}

func (cs *cityService) UpdateCity(ctx context.Context, id int64, name string) error {
	err := cs.cityRepo.Update(database.City{
		ID:   id,
		Name: name,
	})
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("error updating city")
		return fmt.Errorf("error updating city: %w", err)
	}
	return nil
}

func (cs *cityService) DeleteCity(ctx context.Context, id int64) error {
	err := cs.airportRepo.DeleteByCityID(id)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithField("cityId", id).
			Error("error deleting airports in city")
		return fmt.Errorf("error deleting airports in city: %w", err)
	}
	err = cs.cityRepo.Delete(id)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithField("cityId", id).
			Error("error deleting city")
		return fmt.Errorf("error deleting city: %w", err)
	}
	return nil
}

func (cs *cityService) ImportCities(ctx context.Context, content []byte) error {
	reader := bytes.NewReader(content)
	csvReader := csv.NewReader(reader)
	var (
		fields []string
		err    error
	)
	for {
		fields, err = csvReader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			logrus.WithContext(ctx).WithError(err).
				Error("failed to read city row")
		}
		_, err = cs.cityRepo.SaveNew(database.City{
			Name: fields[0],
		})
		if err != nil {
			logrus.WithContext(ctx).WithError(err).
				WithField("csv_row", strings.Join(fields, ",")).
				Error("failed to import city")
		}
	}
}
