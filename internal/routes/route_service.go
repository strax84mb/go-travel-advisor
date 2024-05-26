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

func (rs *routeService) ListRoutesForCity(ctx context.Context, cityID int64, direction Direction, pagination utils.Pagination) ([]database.Route, error) {
	return nil, nil
}

func (rs *routeService) ListRoutesForAirport(ctx context.Context, airportID int64, direction Direction, pagination utils.Pagination) ([]database.Route, error) {
	return nil, nil
}

func (rs *routeService) RouteByID(ctx context.Context, id int64) (*database.Route, error) {
	return nil, nil
}

func (rs *routeService) SaveNew(ctx context.Context, route database.Route) (*database.Route, error) {
	return nil, nil
}

func (rs *routeService) UpdateRoutePrice(ctx context.Context, id int64, price float32) error {
	return nil
}

func (rs *routeService) DeleteRoute(ctx context.Context, id int64) error {
	return nil
}
