package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/routes/repository"
	"gitlab.strale.io/go-travel/internal/utils"
	"gitlab.strale.io/go-travel/internal/utils/testutils"
)

type routeRepository interface {
	Find(pagination utils.Pagination) ([]database.Route, error)
	FindRoutesForAirport(
		airportID int64,
		incomming, outgoing bool,
		pagination utils.Pagination,
	) ([]database.Route, error)
	FindRoutesForCity(
		cityID int64,
		incomming, outgoing bool,
		pagination utils.Pagination,
	) ([]database.Route, error)
	FindByID(id int64, loadAirports bool) (*database.Route, error)
	Insert(route database.Route) (*database.Route, error)
	UpdatePrice(id int64, price float32) error
	Delete(id int64) error
	FindDestinations(startAirportsIDs []int64, cityIDsToSkip []int64) ([]database.Route, error)
}

var (
	city1 *database.City = &database.City{
		ID:   1,
		Name: "Beograd",
	}
	city2 *database.City = &database.City{
		ID:   2,
		Name: "Amsterdam",
	}
	city3 *database.City = &database.City{
		ID:   3,
		Name: "London",
	}

	airport1 *database.Airport = &database.Airport{
		ID:     1,
		Name:   "Nikola Tesla",
		CityID: 1,
	}
	airport2 *database.Airport = &database.Airport{
		ID:     2,
		Name:   "Schipol",
		CityID: 2,
	}
	airport3 *database.Airport = &database.Airport{
		ID:     3,
		Name:   "Heathrow",
		CityID: 3,
	}

	route1 *database.Route = &database.Route{
		ID:            1,
		SourceID:      1,
		DestinationID: 2,
		Price:         10,
	}
	route2 *database.Route = &database.Route{
		ID:            2,
		SourceID:      1,
		DestinationID: 3,
		Price:         20,
	}
	route3 *database.Route = &database.Route{
		ID:            3,
		SourceID:      2,
		DestinationID: 1,
		Price:         30,
	}
	route4 *database.Route = &database.Route{
		ID:            4,
		SourceID:      2,
		DestinationID: 3,
		Price:         40,
	}
	route5 *database.Route = &database.Route{
		ID:            5,
		SourceID:      3,
		DestinationID: 1,
		Price:         50,
	}
	route6 *database.Route = &database.Route{
		ID:            6,
		SourceID:      3,
		DestinationID: 2,
		Price:         60,
	}
)

func setupRoutes(t *testing.T) (routeRepository, *database.Database) {
	db, err := testutils.ConnectToTestDB(true)
	assert.NoError(t, err)

	tx := db.DB.Create(city1)
	assert.NoError(t, tx.Error)
	tx = db.DB.Create(city2)
	assert.NoError(t, tx.Error)
	tx = db.DB.Create(city3)
	assert.NoError(t, tx.Error)

	tx = db.DB.Create(airport1)
	assert.NoError(t, tx.Error)
	tx = db.DB.Create(airport2)
	assert.NoError(t, tx.Error)
	tx = db.DB.Create(airport3)
	assert.NoError(t, tx.Error)

	tx = db.DB.Create(route1)
	assert.NoError(t, tx.Error)
	tx = db.DB.Create(route2)
	assert.NoError(t, tx.Error)
	tx = db.DB.Create(route3)
	assert.NoError(t, tx.Error)
	tx = db.DB.Create(route4)
	assert.NoError(t, tx.Error)
	tx = db.DB.Create(route5)
	assert.NoError(t, tx.Error)
	tx = db.DB.Create(route6)
	assert.NoError(t, tx.Error)

	return repository.NewRouteRepository(&db), &db
}

func TestFind(t *testing.T) {
	repo, _ := setupRoutes(t)
	routes, err := repo.Find(utils.PaginationFrom(0, 10))
	assert.NoError(t, err)
	assert.Equal(t, 6, len(routes))
}

func getRoute(routeID int64, routes []database.Route) *database.Route {
	for _, route := range routes {
		if route.ID == routeID {
			return &route
		}
	}
	return nil
}

func TestFindRoutesForAirport_Incoming(t *testing.T) {
	repo, _ := setupRoutes(t)
	routes, err := repo.FindRoutesForAirport(airport1.ID, true, false, utils.PaginationFrom(0, 10))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(routes))
	route := getRoute(3, routes)
	assert.EqualValues(t, 30, route.Price)
	route = getRoute(5, routes)
	assert.EqualValues(t, 50, route.Price)
}

func TestFindRoutesForAirport_Outgoing(t *testing.T) {
	repo, _ := setupRoutes(t)
	routes, err := repo.FindRoutesForAirport(airport1.ID, false, true, utils.PaginationFrom(0, 10))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(routes))
	route := getRoute(1, routes)
	assert.EqualValues(t, 10, route.Price)
	route = getRoute(2, routes)
	assert.EqualValues(t, 20, route.Price)
}

func TestFindRoutesForAirport_All(t *testing.T) {
	repo, _ := setupRoutes(t)
	routes, err := repo.FindRoutesForAirport(airport1.ID, true, true, utils.PaginationFrom(0, 10))
	assert.NoError(t, err)
	assert.Equal(t, 4, len(routes))
	route := getRoute(1, routes)
	assert.EqualValues(t, 10, route.Price)
	route = getRoute(2, routes)
	assert.EqualValues(t, 20, route.Price)
	route = getRoute(3, routes)
	assert.EqualValues(t, 30, route.Price)
	route = getRoute(5, routes)
	assert.EqualValues(t, 50, route.Price)
}

func TestFindRoutesForCity_Incomming(t *testing.T) {
	repo, _ := setupRoutes(t)
	routes, err := repo.FindRoutesForCity(city1.ID, true, false, utils.PaginationFrom(0, 10))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(routes))
	route := getRoute(3, routes)
	assert.EqualValues(t, 30, route.Price)
	route = getRoute(5, routes)
	assert.EqualValues(t, 50, route.Price)
}

func TestFindRoutesForCity_Outgoing(t *testing.T) {
	repo, _ := setupRoutes(t)
	routes, err := repo.FindRoutesForCity(city1.ID, false, true, utils.PaginationFrom(0, 10))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(routes))
	route := getRoute(1, routes)
	assert.EqualValues(t, 10, route.Price)
	route = getRoute(2, routes)
	assert.EqualValues(t, 20, route.Price)
}

func TestFindRoutesForCity_All(t *testing.T) {
	repo, _ := setupRoutes(t)
	routes, err := repo.FindRoutesForCity(city1.ID, true, true, utils.PaginationFrom(0, 10))
	assert.NoError(t, err)
	assert.Equal(t, 4, len(routes))
	route := getRoute(1, routes)
	assert.EqualValues(t, 10, route.Price)
	route = getRoute(2, routes)
	assert.EqualValues(t, 20, route.Price)
	route = getRoute(3, routes)
	assert.EqualValues(t, 30, route.Price)
	route = getRoute(5, routes)
	assert.EqualValues(t, 50, route.Price)
}

func TestFindByID_NotFound(t *testing.T) {
	repo, _ := setupRoutes(t)
	_, err := repo.FindByID(-1, false)
	assert.EqualError(t, err, "entity not found")
}

func routesEqual(r1, r2 *database.Route) bool {
	return r1.ID == r2.ID && r1.Price == r2.Price &&
		r1.SourceID == r2.SourceID &&
		r1.DestinationID == r2.DestinationID
}

func TestFindByID_Success(t *testing.T) {
	repo, _ := setupRoutes(t)
	route, err := repo.FindByID(route1.ID, false)
	assert.NoError(t, err)
	assert.True(t, routesEqual(route1, route))
	assert.Nil(t, route.Source)
	assert.Nil(t, route.Destination)
}

func TestFindByID_Success_LoadAirport(t *testing.T) {
	repo, _ := setupRoutes(t)
	route, err := repo.FindByID(route1.ID, true)
	assert.NoError(t, err)
	assert.True(t, routesEqual(route1, route))
	assert.NotNil(t, route.Source)
	assert.EqualValues(t, airport1.ID, route.Source.ID)
	assert.EqualValues(t, airport1.Name, route.Source.Name)
	assert.EqualValues(t, airport1.CityID, route.Source.CityID)
	assert.NotNil(t, route.Destination)
	assert.EqualValues(t, airport2.ID, route.Destination.ID)
	assert.EqualValues(t, airport2.Name, route.Destination.Name)
	assert.EqualValues(t, airport2.CityID, route.Destination.CityID)
}

func TestFindDestinations(t *testing.T) {
	repo, _ := setupRoutes(t)
	routes, err := repo.FindDestinations(
		[]int64{1, 2},
		[]int64{1},
	)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(routes))
}
