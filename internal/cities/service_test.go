package cities_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.strale.io/go-travel/internal/cities"
	"gitlab.strale.io/go-travel/internal/cities/repository"
	"gitlab.strale.io/go-travel/internal/cities/testutils"
	"gitlab.strale.io/go-travel/internal/database"
)

type cityService interface {
	ListCities(ctx context.Context, offset, limit int) ([]database.City, error)
	FindByID(ctx context.Context, id int64) (database.City, error)
	SaveNewCity(ctx context.Context, name string) (database.City, error)
	UpdateCity(ctx context.Context, id int64, name string) error
	DeleteCity(ctx context.Context, id int64) error
	ImportCities(ctx context.Context, content []byte) error
}

type testKit struct {
	cityRepo *testutils.CityRepositoryMock
	service  cityService
}

func setUp() *testKit {
	logrus.SetOutput(io.Discard)
	cityRepo := &testutils.CityRepositoryMock{}
	service := cities.NewCityService(cityRepo, nil)
	return &testKit{
		service:  service,
		cityRepo: cityRepo,
	}
}

func TestListCities_Success(t *testing.T) {
	kit := setUp()
	ctx := context.Background()
	kit.cityRepo.FindFn = func(fi repository.FindInput) ([]database.City, error) {
		assert.Equal(t, 10, fi.Limit)
		assert.Equal(t, 0, fi.Offset)
		return []database.City{
			{
				ID: 123,
			},
		}, nil
	}

	cities, err := kit.service.ListCities(ctx, 0, 10)

	assert.NoError(t, err)
	assert.NotNil(t, cities)
	assert.Equal(t, 1, len(cities))
	assert.Equal(t, int64(123), cities[0].ID)
}

func TestListCities_Fail(t *testing.T) {
	kit := setUp()
	ctx := context.Background()
	kit.cityRepo.FindFn = func(fi repository.FindInput) ([]database.City, error) {
		return nil, errors.New("some error")
	}

	cities, err := kit.service.ListCities(ctx, 0, 10)

	assert.Nil(t, cities)
	assert.EqualError(t, err, "error listing cities: some error")
}

func TestFindByID_Error(t *testing.T) {
	kit := setUp()
	cityID := 123
	kit.cityRepo.FindByIDFn = func(id int64, preload bool) (database.City, error) {
		assert.True(t, preload)
		assert.EqualValues(t, cityID, id)
		return database.City{}, errors.New("some error")
	}

	_, err := kit.service.FindByID(context.Background(), int64(cityID))

	assert.EqualError(t, err, "error loading city: some error")
}

func TestFindByID_NotFound(t *testing.T) {
	kit := setUp()
	cityID := 123
	kit.cityRepo.FindByIDFn = func(id int64, preload bool) (database.City, error) {
		return database.City{}, database.ErrNotFound
	}

	_, err := kit.service.FindByID(context.Background(), int64(cityID))

	assert.EqualError(t, err, database.ErrNotFound.Error())
}

func TestFindByID_Success(t *testing.T) {
	kit := setUp()
	cityID := 123
	kit.cityRepo.FindByIDFn = func(id int64, preload bool) (database.City, error) {
		assert.True(t, preload)
		assert.EqualValues(t, cityID, id)
		return database.City{
			ID: int64(cityID),
		}, nil
	}

	city, err := kit.service.FindByID(context.Background(), int64(cityID))

	assert.NoError(t, err)
	assert.EqualValues(t, cityID, city.ID)
}
