//go:build test

package servicetest

import "gitlab.strale.io/go-travel/internal/database"

type MockCityRepository struct {
	FindByIDFn   func(int64, bool) (database.City, error)
	FindByNameFn func(string, bool) (database.City, error)
}

func (m *MockCityRepository) FindByID(id int64, preload bool) (database.City, error) {
	return m.FindByIDFn(id, preload)
}

func (m *MockCityRepository) FindByName(name string, preload bool) (database.City, error) {
	return m.FindByNameFn(name, preload)
}
