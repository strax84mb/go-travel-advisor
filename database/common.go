package database

import (
	"database/sql"

	"gitlab.strale.io/go-travel/common"
)

func performStatement(location string, execution func([]interface{}) (sql.Result, error), params ...interface{}) *common.GeneralError {
	result, err := execution(params)
	if err != nil {
		return &common.GeneralError{
			Message:  "Error while executing statement",
			Location: location,
			Cause:    err,
		}
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return &common.GeneralError{
			Message:  "Error while checking number of affected rows",
			Location: location,
			Cause:    err,
		}
	}
	if affected < 1 {
		return &common.GeneralError{
			Message:   "No rows affected by this statement",
			Location:  location,
			ErrorType: common.NoRowsAffected,
		}
	}
	return nil
}

func performListSelection(
	location string,
	resultSize int,
	array interface{},
	selection func([]interface{}) (*sql.Rows, error),
	conversion func(*sql.Rows, interface{}, int) error,
	params ...interface{}) *common.GeneralError {
	rows, err := selection(params)
	if err != nil {
		return &common.GeneralError{
			Message:  "Error while executing query",
			Location: location,
			Cause:    err,
		}
	}
	defer rows.Close()
	for i := 0; rows.Next() && i < resultSize; i++ {
		err = conversion(rows, array, i)
		if err != nil {
			return &common.GeneralError{
				Message:  "Error while reading columns",
				Location: location,
				Cause:    err,
			}
		}
	}
	return nil
}

func performSingleSelection(
	location string,
	selection func([]interface{}) (*sql.Rows, error),
	conversion func(*sql.Rows) (interface{}, error),
	params ...interface{}) (interface{}, bool, *common.GeneralError) {
	rows, err := selection(params)
	if err != nil {
		return nil, false, &common.GeneralError{
			Message:  "Error executing selection query",
			Location: location,
			Cause:    err,
		}
	}
	defer rows.Close()
	for rows.Next() {
		returnValue, err := conversion(rows)
		if err != nil {
			return nil, false, &common.GeneralError{
				Message:  "Error while reading columns",
				Location: location,
				Cause:    err,
			}
		}
		return returnValue, true, nil
	}
	return nil, false, nil
}
