package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	// Used to initiate DB
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var gdb dbWrapper

// InitDb - Used to initialize DB
func InitDb() {
	var err error
	os.Remove("test.db")
	initDB, err := gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err.Error())
	}
	gdb = &dbWrapperImpl{initDB}
	gdb.AutoMigrate(&User{}, &City{}, &Comment{}, &Airport{}, &Route{})
	user := User{
		Username: "admin",
		Password: "admin",
		Role:     UserRoleAdmin,
	}
	genErr := saveUser(user)
	if genErr != nil {
		fmt.Println(genErr.Error())
	}
}

// SetDBForTesting - for testing purposes
func SetDBForTesting(db *gorm.DB) {
	//gdb = db
}

func handleInitError(text string, err error) {
	if err != nil {
		log.Fatalf(text, err.Error())
	}
}

type dbWrapper interface {
	AutoMigrate(...interface{}) dbWrapper
	Select(interface{}, ...interface{}) dbWrapper
	Where(interface{}, ...interface{}) dbWrapper
	First(interface{}, ...interface{}) dbWrapper
	Create(interface{}) dbWrapper
	Save(interface{}) dbWrapper
	Delete(interface{}, ...interface{}) dbWrapper
	Model(interface{}) dbWrapper
	Count(interface{}) dbWrapper
	Preload(string, ...interface{}) dbWrapper
	Find(interface{}, ...interface{}) dbWrapper
	Table(string) dbWrapper
	Joins(string, ...interface{}) dbWrapper
	Order(interface{}, ...bool) dbWrapper
	Limit(interface{}) dbWrapper
	NewRecord(interface{}) bool
	RecordNotFound() bool
	Transaction(func(*gorm.DB) error) error
	Error() error
	RowsAffected() int64
}

type dbWrapperImpl struct {
	db *gorm.DB
}

func (db *dbWrapperImpl) AutoMigrate(params ...interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.AutoMigrate(params...)}
}

func (db *dbWrapperImpl) Select(query interface{}, args ...interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.Select(query, args...)}
}

func (db *dbWrapperImpl) Where(query interface{}, args ...interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.Where(query, args...)}
}

func (db *dbWrapperImpl) First(out interface{}, where ...interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.First(out, where...)}
}

func (db *dbWrapperImpl) RecordNotFound() bool {
	return db.db.RecordNotFound()
}

func (db *dbWrapperImpl) Create(value interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.Create(value)}
}

func (db *dbWrapperImpl) Save(value interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.Save(value)}
}

func (db *dbWrapperImpl) Delete(value interface{}, where ...interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.Delete(value, where...)}
}

func (db *dbWrapperImpl) Model(value interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.Model(value)}
}

func (db *dbWrapperImpl) Count(value interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.Count(value)}
}

func (db *dbWrapperImpl) Error() error {
	return db.db.Error
}

func (db *dbWrapperImpl) Preload(column string, conditions ...interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.Preload(column, conditions...)}
}

func (db *dbWrapperImpl) Find(out interface{}, where ...interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.Find(out, where...)}
}

func (db *dbWrapperImpl) NewRecord(value interface{}) bool {
	return db.db.NewRecord(value)
}

func (db *dbWrapperImpl) Table(name string) dbWrapper {
	return &dbWrapperImpl{db.db.Table(name)}
}

func (db *dbWrapperImpl) Joins(query string, args ...interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.Joins(query, args...)}
}

func (db *dbWrapperImpl) Order(value interface{}, reorder ...bool) dbWrapper {
	return &dbWrapperImpl{db.db.Order(value, reorder...)}
}
func (db *dbWrapperImpl) Limit(limit interface{}) dbWrapper {
	return &dbWrapperImpl{db.db.Limit(limit)}
}

func (db *dbWrapperImpl) Transaction(txfunc func(*gorm.DB) error) error {
	return db.Transaction(txfunc)
}

func (db *dbWrapperImpl) RowsAffected() int64 {
	return db.db.RowsAffected
}
