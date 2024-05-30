package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.strale.io/go-travel/internal/airports/repository"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils"
	"gitlab.strale.io/go-travel/internal/utils/testutils"
)

type setupStruct struct {
	airport1 *database.Airport
	airport2 *database.Airport
	city1    *database.City
	city2    *database.City
}

type airportRepository interface {
	List(pagination utils.Pagination) ([]database.Airport, error)
	ListInCity(input repository.ListInCityInput) ([]database.Airport, error)
	FindByID(id int64) (database.Airport, error)
	SaveAirport(airport database.Airport) (database.Airport, error)
	UpdateAirport(airport database.Airport) error
	DeleteByID(id int64) error
	DeleteByCityID(cityID int64) error
}

type repoTestKit struct {
	repo  airportRepository
	setup *setupStruct
	db    *database.Database
}

func setupAirports() (*repoTestKit, error) {
	setup := setupStruct{
		city1: &database.City{
			Name: "Beograd",
		},
		city2: &database.City{
			Name: "Amsterdam",
		},
		airport1: &database.Airport{
			Name: "Nikola Tesla",
		},
		airport2: &database.Airport{
			Name: "Schipol",
		},
	}
	db, err := testutils.ConnectToTestDB(true)
	if err != nil {
		return nil, err
	}

	tx := db.DB.Create(setup.city1)
	if tx.Error != nil {
		return nil, tx.Error
	}
	tx = db.DB.Create(setup.city2)
	if tx.Error != nil {
		return nil, tx.Error
	}
	setup.airport1.CityID = setup.city1.ID
	tx = db.DB.Create(setup.airport1)
	if tx.Error != nil {
		return nil, tx.Error
	}
	setup.airport2.CityID = setup.city2.ID
	tx = db.DB.Create(setup.airport2)
	if tx.Error != nil {
		return nil, tx.Error
	}

	repo := repository.NewAirportRepository(&db)

	return &repoTestKit{
		repo:  repo,
		setup: &setup,
		db:    &db,
	}, nil
}

func TestFindByID(t *testing.T) {
	kit, err := setupAirports()
	assert.NoError(t, err)
	airport, err := kit.repo.FindByID(kit.setup.airport1.ID)
	assert.NoError(t, err)
	assert.Equal(t, kit.setup.airport1.ID, airport.ID)
	assert.Equal(t, kit.setup.airport1.Name, airport.Name)
	assert.Equal(t, kit.setup.city1.ID, airport.CityID)
	assert.NotNil(t, airport.City)
	assert.Equal(t, kit.setup.city1.ID, airport.City.ID)
	assert.Equal(t, kit.setup.city1.Name, airport.City.Name)
}

func airportPresent(airports []database.Airport, id int64, name string) bool {
	for _, airport := range airports {
		if airport.ID == id && airport.Name == name {
			return true
		}
	}
	return false
}

func TestList(t *testing.T) {
	kit, err := setupAirports()
	assert.NoError(t, err)
	airports, err := kit.repo.List(utils.PaginationFrom(0, 10))
	assert.NoError(t, err)
	assert.True(t, airportPresent(airports, kit.setup.airport1.ID, kit.setup.airport1.Name))
	assert.True(t, airportPresent(airports, kit.setup.airport2.ID, kit.setup.airport2.Name))
}

func TestListInCity(t *testing.T) {
	kit, err := setupAirports()
	assert.NoError(t, err)
	airports, err := kit.repo.ListInCity(repository.ListInCityInput{
		Pagination: utils.PaginationFrom(0, 10),
		CityID:     kit.setup.city1.ID,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(airports))
	assert.True(t, airportPresent(airports, kit.setup.airport1.ID, kit.setup.airport1.Name))
}

func TestDelete(t *testing.T) {
	kit, err := setupAirports()
	assert.NoError(t, err)
	id := kit.setup.airport1.ID
	err = kit.repo.DeleteByID(id)
	assert.NoError(t, err)
	_, err = kit.repo.FindByID(id)
	assert.EqualError(t, err, "entity not found")
}

func TestDeleteInCity(t *testing.T) {
	kit, err := setupAirports()
	assert.NoError(t, err)
	id := kit.setup.airport1.ID
	err = kit.repo.DeleteByCityID(kit.setup.city1.ID)
	assert.NoError(t, err)
	_, err = kit.repo.FindByID(id)
	assert.EqualError(t, err, "entity not found")
}

func TestUpdate(t *testing.T) {
	kit, err := setupAirports()
	assert.NoError(t, err)
	id := kit.setup.airport1.ID
	err = kit.repo.UpdateAirport(database.Airport{
		ID:     id,
		Name:   "Updated name",
		CityID: kit.setup.city1.ID,
	})
	assert.NoError(t, err)
	airport, err := kit.repo.FindByID(id)
	assert.NoError(t, err)
	assert.Equal(t, "Updated name", airport.Name)
	assert.Equal(t, kit.setup.city1.ID, airport.CityID)
}
