package testutils

import (
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils"
)

type MockRegistry struct {
	reg map[string][][]interface{}
}

func (mr *MockRegistry) Record(function string, args ...interface{}) {
	argsList, found := mr.reg[function]
	if !found {
		argsList = [][]interface{}{}
	}
	argsList = append(argsList, args)
	mr.reg[function] = argsList
}

func (mr *MockRegistry) Times(function string) int {
	argsList, found := mr.reg[function]
	if !found {
		return 0
	}
	return len(argsList)
}

func (mr *MockRegistry) WithArgs(function string, time int, args ...interface{}) bool {
	return false
}

type CityRepositoryMock struct {
	FindFn       func(utils.Pagination) ([]database.City, error)
	FindByIDsFn  func([]int64) ([]database.City, error)
	FindByIDFn   func(int64, bool) (database.City, error)
	FindByNameFn func(string, bool) (database.City, error)
	SaveNewFn    func(database.City) (database.City, error)
	UpdateFn     func(database.City) error
	DeleteFn     func(int64) error
}

func (cr *CityRepositoryMock) Find(pagination utils.Pagination) ([]database.City, error) {
	return cr.FindFn(pagination)
}

func (cr *CityRepositoryMock) FindByIDs(ids []int64) ([]database.City, error) {
	return cr.FindByIDsFn(ids)
}

func (cr *CityRepositoryMock) FindByID(id int64, preload bool) (database.City, error) {
	return cr.FindByIDFn(id, preload)
}

func (cr *CityRepositoryMock) FindByName(name string, preload bool) (database.City, error) {
	return cr.FindByNameFn(name, preload)
}

func (cr *CityRepositoryMock) SaveNew(city database.City) (database.City, error) {
	return cr.SaveNewFn(city)
}

func (cr *CityRepositoryMock) Update(city database.City) error {
	return cr.UpdateFn(city)
}

func (cr *CityRepositoryMock) Delete(id int64) error {
	return cr.DeleteFn(id)
}
