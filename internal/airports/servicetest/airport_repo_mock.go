//go:build test

package servicetest

import (
	"gitlab.strale.io/go-travel/internal/airports/repository"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils"
)

type MockAirportRepository struct {
	ListFn           func(utils.Pagination) ([]database.Airport, error)
	ListInCityFn     func(repository.ListInCityInput) ([]database.Airport, error)
	FindByIDFn       func(int64) (database.Airport, error)
	SaveAirportFn    func(database.Airport) (database.Airport, error)
	UpdateAirportFn  func(database.Airport) error
	DeleteByIDFn     func(int64) error
	DeleteByCityIDFn func(int64) error
}

func (m *MockAirportRepository) List(input utils.Pagination) ([]database.Airport, error) {
	return m.ListFn(input)
}

func (m *MockAirportRepository) ListInCity(input repository.ListInCityInput) ([]database.Airport, error) {
	return m.ListInCityFn(input)
}

func (m *MockAirportRepository) FindByID(id int64) (database.Airport, error) {
	return m.FindByIDFn(id)
}

func (m *MockAirportRepository) SaveAirport(airport database.Airport) (database.Airport, error) {
	return m.SaveAirportFn(airport)
}

func (m *MockAirportRepository) UpdateAirport(airport database.Airport) error {
	return m.UpdateAirportFn(airport)
}

func (m *MockAirportRepository) DeleteByID(id int64) error {
	return m.DeleteByIDFn(id)
}

func (m *MockAirportRepository) DeleteByCityID(cityID int64) error {
	return m.DeleteByCityIDFn(cityID)
}
