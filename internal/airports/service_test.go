//go:build test

package airports

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.strale.io/go-travel/internal/airports/servicetest"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils"
)

type serviceKit struct {
	cityRepo    *servicetest.MockCityRepository
	airportRepo *servicetest.MockAirportRepository
	service     *AirportService
}

func setupKit() *serviceKit {
	logrus.SetLevel(logrus.FatalLevel)

	cityRepo := &servicetest.MockCityRepository{}
	airportRepo := &servicetest.MockAirportRepository{}

	service := NewAirportService(airportRepo, cityRepo)

	return &serviceKit{
		cityRepo:    cityRepo,
		airportRepo: airportRepo,
		service:     service,
	}
}

var airport1 = database.Airport{
	ID:     1,
	Name:   "Nikola Tesla",
	CityID: 12,
}

func TestListAirports_Fail(t *testing.T) {
	kit := setupKit()
	someError := errors.New("some error")
	kit.airportRepo.ListFn = func(p utils.Pagination) ([]database.Airport, error) {
		assert.Equal(t, 0, p.Offset)
		assert.Equal(t, 10, p.Limit)
		return nil, someError
	}

	airports, err := kit.service.ListAirports(context.Background(), utils.PaginationFrom(0, 10))
	assert.Nil(t, airports)
	assert.EqualError(t, err, "error while listing airports: some error")
}

func TestListAirports_Success(t *testing.T) {
	kit := setupKit()

	kit.airportRepo.ListFn = func(p utils.Pagination) ([]database.Airport, error) {
		assert.Equal(t, 0, p.Offset)
		assert.Equal(t, 10, p.Limit)
		return []database.Airport{airport1}, nil
	}

	airports, err := kit.service.ListAirports(context.Background(), utils.PaginationFrom(0, 10))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(airports))
	assert.Equal(t, airport1.ID, airports[0].ID)
}
