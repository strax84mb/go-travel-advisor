package database

import (
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type dbWrapperMock struct {
	Value              []interface{}
	ValueRead          *int
	RowsAffectedValues []int64
	RowsAffectedRead   *int
	Errors             []error
	ErrorRead          *int
	NewRecords         []bool
	NewRecordRead      *int
	Counts             []int
	CountRead          *int
	NoRecordFound      []bool
	NoRecordFoundRead  *int
}

func (db *dbWrapperMock) copydb() *dbWrapperMock {
	return &dbWrapperMock{
		Value:              db.Value,
		ValueRead:          db.ValueRead,
		RowsAffectedValues: db.RowsAffectedValues,
		RowsAffectedRead:   db.RowsAffectedRead,
		Errors:             db.Errors,
		ErrorRead:          db.ErrorRead,
		NewRecords:         db.NewRecords,
		NewRecordRead:      db.NewRecordRead,
		Counts:             db.Counts,
		CountRead:          db.CountRead,
		NoRecordFound:      db.NoRecordFound,
		NoRecordFoundRead:  db.NoRecordFoundRead,
	}
}

func (db *dbWrapperMock) setZeros() {
	valueRead := 0
	rowsAffectedRead := 0
	errorRead := 0
	newRecordRead := 0
	countRead := 0
	noRecordFoundRead := 0
	db.ValueRead = &valueRead
	db.RowsAffectedRead = &rowsAffectedRead
	db.ErrorRead = &errorRead
	db.NewRecordRead = &newRecordRead
	db.CountRead = &countRead
	db.NoRecordFoundRead = &noRecordFoundRead
}

func (db *dbWrapperMock) AutoMigrate(args ...interface{}) dbWrapper {
	return db.copydb()
}

func (db *dbWrapperMock) Select(query interface{}, args ...interface{}) dbWrapper {
	return db.copydb()
}

func (db *dbWrapperMock) Where(condition interface{}, values ...interface{}) dbWrapper {
	return db.copydb()
}

func (db *dbWrapperMock) First(entity interface{}, where ...interface{}) dbWrapper {
	q := db.ValueRead
	result := db.Value[*q]
	copier.Copy(entity, result)
	(*q)++
	return db.copydb()
}

func (db *dbWrapperMock) Create(entity interface{}) dbWrapper {
	return db.copydb()
}

func (db *dbWrapperMock) Save(entity interface{}) dbWrapper {
	return db.copydb()
}

func (db *dbWrapperMock) Delete(entity interface{}, where ...interface{}) dbWrapper {
	return db.copydb()
}

func (db *dbWrapperMock) Model(model interface{}) dbWrapper {
	return db.copydb()
}

func (db *dbWrapperMock) Count(count interface{}) dbWrapper {
	q := db.CountRead
	result := db.Counts[*q]
	*(count.(*int)) = result
	(*q)++
	return db.copydb()
}

func (db *dbWrapperMock) Preload(fields string, values ...interface{}) dbWrapper {
	return db.copydb()
}

func (db *dbWrapperMock) Find(entity interface{}, where ...interface{}) dbWrapper {
	q := db.ValueRead
	result := db.Value[*q]
	copier.Copy(entity, result)
	(*q)++
	return db.copydb()
}

func (db *dbWrapperMock) Table(name string) dbWrapper {
	return db.copydb()
}

func (db *dbWrapperMock) Joins(query string, values ...interface{}) dbWrapper {
	return db.copydb()
}

func (db *dbWrapperMock) Order(order interface{}, backs ...bool) dbWrapper {
	return db.copydb()
}

func (db *dbWrapperMock) Limit(number interface{}) dbWrapper {
	return db.copydb()
}

func (db *dbWrapperMock) NewRecord(entity interface{}) bool {
	q := db.NewRecordRead
	err := db.NewRecords[*q]
	(*q)++
	return err
}

func (db *dbWrapperMock) RecordNotFound() bool {
	q := db.NoRecordFoundRead
	nrf := db.NoRecordFound[*q]
	(*q)++
	return nrf
}

func (db *dbWrapperMock) Transaction(txfunc func(*gorm.DB) error) error {
	return nil
}

func (db *dbWrapperMock) Error() error {
	q := db.ErrorRead
	err := db.Errors[*q]
	(*q)++
	return err
}

func (db *dbWrapperMock) RowsAffected() int64 {
	q := db.RowsAffectedRead
	rowsAffected := db.RowsAffectedValues[*q]
	(*q)++
	return rowsAffected
}
