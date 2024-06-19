package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.strale.io/go-travel/internal/cities/repository"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils"
	"gitlab.strale.io/go-travel/internal/utils/testutils"
)

var (
	city1 database.City = database.City{
		Name: "Beograd",
	}
	city2 database.City = database.City{
		Name: "Budapest",
	}
)

type cityRepository interface {
	Find(pagination utils.Pagination) ([]database.City, error)
	FindByIDs(ids []int64) ([]database.City, error)
	FindByID(id int64, preload bool) (database.City, error)
	FindByName(name string, preload bool) (database.City, error)
	Update(city database.City) error
	SaveNew(city database.City) (database.City, error)
	Delete(id int64) error
}

func setupCities() (cityRepository, *database.Database, error) {
	db, err := testutils.ConnectToTestDB(false)
	if err != nil {
		return nil, nil, err
	}
	tx := db.DB.Create(&city1)
	if tx.Error != nil {
		return nil, nil, tx.Error
	}
	tx = db.DB.Create(&city2)
	if tx.Error != nil {
		return nil, nil, tx.Error
	}
	return repository.NewCityRepository(&db), &db, nil
}

func cityPresent(cities []database.City, name string) bool {
	for _, city := range cities {
		if city.Name == name {
			return true
		}
	}
	return false
}

func getCity(cities []database.City, name string) *database.City {
	for _, city := range cities {
		if city.Name == name {
			return &city
		}
	}
	return nil
}

func TestGetCities(t *testing.T) {
	repo, _, err := setupCities()
	assert.NoError(t, err)
	cities, err := repo.Find(utils.PaginationFrom(0, 10))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cities))
	assert.True(t, cityPresent(cities, "Beograd"))
	assert.True(t, cityPresent(cities, "Budapest"))
	bgCity := getCity(cities, "Beograd")
	city, err := repo.FindByID(bgCity.ID, false)
	assert.NoError(t, err)
	assert.Equal(t, bgCity.Name, city.Name)
	city, err = repo.FindByName("Budapest", false)
	assert.NoError(t, err)
	assert.Equal(t, "Budapest", city.Name)
}

func TestCreateUpdateDelete(t *testing.T) {
	repo, _, err := setupCities()
	assert.NoError(t, err)
	city := database.City{
		Name: "Subotica",
	}
	assert.Zero(t, city.ID)
	// save new city
	city, err = repo.SaveNew(city)
	assert.NoError(t, err)
	assert.NotZero(t, city.ID)
	id := city.ID
	assert.Equal(t, "Subotica", city.Name)
	// update
	city.Name = "Novi Sad"
	err = repo.Update(city)
	assert.NoError(t, err)
	city, err = repo.FindByID(id, false)
	assert.NoError(t, err)
	assert.Equal(t, "Novi Sad", city.Name)
	// delete
	err = repo.Delete(id)
	assert.NoError(t, err)
	// confirm deleted
	_, err = repo.FindByID(id, false)
	assert.EqualError(t, err, "entity not found")
	err = repo.Delete(id)
	assert.EqualError(t, err, "entity not found")
	err = repo.Update(city)
	assert.EqualError(t, err, "entity not found")
}

func TestFindCityWithPreload(t *testing.T) {
	repo, db, err := setupCities()
	assert.NoError(t, err)
	tx := db.DB.Create(&database.Airport{
		Name:   "Nikola Tesla",
		CityID: city1.ID,
	})
	assert.NoError(t, tx.Error)
	city, err := repo.FindByID(city1.ID, true)
	assert.NoError(t, err)
	assert.Equal(t, "Beograd", city.Name)
	assert.NotNil(t, city.Airports)
	assert.Equal(t, 1, len(city.Airports))
	assert.Equal(t, "Nikola Tesla", city.Airports[0].Name)
}

func TestFindByIDs(t *testing.T) {
	repo, _, err := setupCities()
	assert.NoError(t, err)
	cities, err := repo.FindByIDs([]int64{city1.ID, city2.ID})
	assert.NoError(t, err)
	assert.True(t, cityPresent(cities, "Beograd"))
	assert.True(t, cityPresent(cities, "Budapest"))
}
