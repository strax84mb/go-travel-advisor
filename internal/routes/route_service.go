package routes

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils"
)

type iRouteRepository interface {
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
}

type iAirportRepository interface{}

type routeService struct {
	routeRepo   iRouteRepository
	airportRepo iAirportRepository
}

func NewRouteService(routeRepo iRouteRepository, airportRepo iAirportRepository) *routeService {
	return &routeService{
		routeRepo:   routeRepo,
		airportRepo: airportRepo,
	}
}

type Direction byte

const (
	DIRECTION_ALL Direction = iota
	DIRECTION_INCOMMING
	DIRECTION_OUTGOING
)

func (rs *routeService) ListRoutes(ctx context.Context, pagination utils.Pagination) ([]database.Route, error) {
	list, err := rs.routeRepo.Find(pagination)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("could not list routes")
		return nil, fmt.Errorf("could not list routes: %w", err)
	}
	return list, nil
}

// Returns booleans from Direction
// Format is (incomming bool, outgoing bool)
func boolsFromDirection(direction Direction) (bool, bool) {
	switch direction {
	case DIRECTION_INCOMMING:
		return true, false
	case DIRECTION_OUTGOING:
		return false, true
	default:
		return true, true
	}
}

func (rs *routeService) ListRoutesForCity(
	ctx context.Context,
	cityID int64,
	direction Direction,
	pagination utils.Pagination,
) ([]database.Route, error) {
	incomming, outgoing := boolsFromDirection(direction)
	list, err := rs.routeRepo.FindRoutesForCity(cityID, incomming, outgoing, pagination)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithFields(logrus.Fields{
				"cityId":              cityID,
				"listIncommingRoutes": incomming,
				"listOutgoingRoutes":  outgoing,
			}).
			Error("could not list routes for city")
		return nil, fmt.Errorf("could not list routes for city: %w", err)
	}
	return list, nil
}

func (rs *routeService) ListRoutesForAirport(
	ctx context.Context,
	airportID int64,
	direction Direction,
	pagination utils.Pagination,
) ([]database.Route, error) {
	incomming, outgoing := boolsFromDirection(direction)
	list, err := rs.routeRepo.FindRoutesForAirport(airportID, incomming, outgoing, pagination)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithFields(logrus.Fields{
				"airportId":           airportID,
				"listIncommingRoutes": incomming,
				"listOutgoingRoutes":  outgoing,
			}).
			Error("could not list routes for airport")
		return nil, fmt.Errorf("could not list routes for airport: %w", err)
	}
	return list, nil
}

func (rs *routeService) RouteByID(ctx context.Context, id int64) (*database.Route, error) {
	route, err := rs.routeRepo.FindByID(id, true)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithField("routeId", id).
			Error("failed to read route")
		return nil, fmt.Errorf("failed to read route: %w", err)
	}
	return route, nil
}

func (rs *routeService) SaveNew(ctx context.Context, route database.Route) (*database.Route, error) {
	savedRoute, err := rs.routeRepo.Insert(route)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			Error("could not save route: %w", err)
	}
	return savedRoute, nil
}

func (rs *routeService) UpdateRoutePrice(ctx context.Context, id int64, price float32) error {
	if err := rs.routeRepo.UpdatePrice(id, price); err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithFields(logrus.Fields{
				"routeId":  id,
				"newPrice": price,
			}).
			Error("could not update route price")
		return fmt.Errorf("could not update route price: %w", err)
	}
	return nil
}

func (rs *routeService) DeleteRoute(ctx context.Context, id int64) error {
	err := rs.routeRepo.Delete(id)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithField("routeId", id).
			Error("could not delete route")
		return fmt.Errorf("could not delete route: %w", err)
	}
	return nil
}
