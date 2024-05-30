package routes

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.strale.io/go-travel/internal/airports/repository"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils"
)

type pfRouteRepository interface {
	FindDestinations(startAirportsIDs []int64, cityIDsToSkip []int64) ([]database.Route, error)
}

type pfAirportRepository interface {
	ListInCity(input repository.ListInCityInput) ([]database.Airport, error)
}

type pfCityRepository interface {
	FindByIDs(ids []int64) ([]database.City, error)
}

type pathFindingService struct {
	airportRepo pfAirportRepository
	cityRepo    pfCityRepository
	routeRepo   pfRouteRepository
}

func NewPathFindingService(airportRepo pfAirportRepository, cityRepo pfCityRepository, routeRepo pfRouteRepository) *pathFindingService {
	return &pathFindingService{
		airportRepo: airportRepo,
		cityRepo:    cityRepo,
		routeRepo:   routeRepo,
	}
}

func (pfs *pathFindingService) FindCheapestPath(ctx context.Context, startCityID, finishCityID int64) ([]*database.Route, float32, error) {
	walk := newCheapestRouteWalk(startCityID)
	var err error
	for {
		walk.Expanded = false
		err = pfs.exapand(ctx, &walk, &walk.Root)
		if err != nil {
			return nil, 0, err
		}
		if !walk.Expanded {
			break
		}
	}
	if walk.CheapestPathEndNode == nil {
		return nil, 0, nil
	}
	routes, err := pfs.collectRoutesAndInsertCities(ctx, &walk)
	if err != nil {
		return nil, 0, err
	}
	return routes, walk.CheapestPrice, nil
}

func collectCityIDs(routes []*database.Route) []int64 {
	cityIDMap := make(map[int64]bool)
	var cityID int64
	for _, route := range routes {
		cityID = route.Destination.CityID
		if cityID != 0 {
			cityIDMap[cityID] = true
		}
	}
	list := make([]int64, len(cityIDMap))
	index := 0
	for k := range cityIDMap {
		list[index] = k
		index++
	}
	return list
}

func insertCities(routes []*database.Route, cities []database.City) []*database.Route {
	cityMap := make(map[int64]*database.City)
	for _, city := range cities {
		cityMap[city.ID] = &city
	}
	for _, route := range routes {
		route.Source.City = cityMap[route.Source.CityID]
		route.Destination.City = cityMap[route.Destination.CityID]
	}
	for k := range cityMap {
		cityMap[k] = nil
	}
	return routes
}

func (pfs *pathFindingService) collectRoutesAndInsertCities(ctx context.Context, walk *CheapestRouteWalk) ([]*database.Route, error) {
	routes := walk.CollectCheapestPath()
	cityIDs := collectCityIDs(routes)
	cities, err := pfs.cityRepo.FindByIDs(cityIDs)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			Error("could not get passed cities")
		return nil, fmt.Errorf("could not get passed cities: %w", err)
	}
	routes = routes[1:]
	routes = insertCities(routes, cities)
	return routes, nil
}

func (pfs *pathFindingService) getAllAirportIDs(ctx context.Context, cityID int64) ([]int64, error) {
	airportIDs := make(map[int64]bool)
	var (
		airports []database.Airport
		err      error
		page     int = 0
	)
	pageSize := 100
	for {
		airports, err = pfs.airportRepo.ListInCity(repository.ListInCityInput{
			Pagination: utils.PaginationFrom(page, pageSize),
			CityID:     cityID,
		})
		if len(airports) == 0 {
			break
		}
		if err != nil {
			logrus.WithContext(ctx).WithError(err).
				WithFields(logrus.Fields{
					"cityId":   cityID,
					"page":     page,
					"pageSize": pageSize,
				}).
				Error("could not list airports in the city during walk")
			return nil, fmt.Errorf("could not list airports in the city during walk: %w", err)
		}
		for _, airport := range airports {
			airportIDs[airport.ID] = true
		}
		page++
	}
	ids := make([]int64, len(airportIDs))
	index := 0
	for k := range airportIDs {
		ids[index] = k
		index++
	}
	return ids, nil
}

func (pfs *pathFindingService) exapand(ctx context.Context, walk *CheapestRouteWalk, node *CRNode) error {
	if node == nil {
		return nil
	}
	if walk.CheapestPathEndNode != nil && walk.CheapestPathEndNode.AccumulatedPrice <= node.AccumulatedPrice {
		return nil
	}
	var err error
	if node.Destinations == nil {
		walk.Expanded = true
		var airportIDs []int64
		airportIDs, err = pfs.getAllAirportIDs(ctx, node.Route.Destination.CityID)
		if err != nil {
			return err
		}
		routes, err := pfs.routeRepo.FindDestinations(airportIDs, node.ListPassedCities())
		if err != nil {
			logrus.WithContext(ctx).WithError(err).
				Error("failed to list possible destinations")
			return fmt.Errorf("failed to list possible destinations: %w", err)
		}
		// make destinations and recalculate price
		node.Destinations = make([]*CRNode, len(routes))
		for i, route := range routes {
			newNode := &CRNode{
				Parent:           node,
				Route:            route,
				AccumulatedPrice: node.AccumulatedPrice + route.Price,
			}
			if newNode.Route.Destination.CityID == walk.FinalCityID {
				if walk.CheapestPathEndNode == nil {
					walk.CheapestPathEndNode = newNode
					walk.CheapestPrice = newNode.AccumulatedPrice
				} else if walk.CheapestPrice > newNode.AccumulatedPrice {
					walk.CheapestPathEndNode = newNode
					walk.CheapestPrice = newNode.AccumulatedPrice
				}
			}
			node.Destinations[i] = newNode
		}
	}
	for _, destination := range node.Destinations {
		err = pfs.exapand(ctx, walk, destination)
		if err != nil {
			return err
		}
	}
	return nil
}
