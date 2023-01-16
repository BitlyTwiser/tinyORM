package dialects

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/BitlyTwiser/tinyORM/pkg/logger"
	"github.com/BitlyTwiser/tinyORM/pkg/sqlbuilder"
)

type RawQuery struct {
	stmt  *sql.Stmt
	args  []any
	query string
}

// Executes given query strig to perform which ever action the query denotes
func (rq *RawQuery) Exec() error {
	result, err := rq.stmt.Exec(rq.args...)

	if err != nil {
		return logger.Log.LogError("error occured executing raw query.", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return logger.Log.LogError("error occured calculating affected rows from raw query.", err)
	}
	logger.Log.LogEvent("warn", "executed raw query", "rows affected", rows)

	return nil
}

// All will accept a model, perform the query, and attempt to fill any data values into the given model.
func (rq *RawQuery) All(model any) error {
	var rows *sql.Rows
	var err error
	m := reflect.Indirect(reflect.ValueOf(model))
	if m.Kind() == reflect.Slice {
		if !sqlbuilder.IsPointer(m) {
			return logger.Log.LogError("you must pass a pointer to a struct for all to function", errors.New("in correct value passed. Must be pointer"))
		}

		if len(rq.args) > 0 {
			rows, err = rq.stmt.Query(rq.args...)
			if err != nil {
				return logger.Log.LogError("error occured executing raw query", err)
			}
		} else {
			rows, err = rq.stmt.Query()
			if err != nil {
				return logger.Log.LogError("error occured executing raw query", err)
			}

		}

		defer func() {
			if err := rows.Close(); err != nil {
				logger.Log.LogError("error closing database rows in All query", err)
			}
		}()

		//Make new slice to feed into the incoming model slice
		newS := reflect.MakeSlice(reflect.SliceOf(m.Type().Elem()), 0, 0)

		for rows.Next() {
			// Create new pointer to inner struct type
			newVal := reflect.New(m.Type().Elem())

			// Fill model with data after parsing out attributes for struct
			err := rows.Scan(sqlbuilder.PointerAttributes(newVal)...)
			if err != nil {
				return err
			}

			// Append slice of new val
			newS = reflect.Append(newS, newVal.Elem())
		}

		// Ensure rows did not encounter an error when calling Next()
		if err := rows.Err(); err != nil {
			return err
		}

		// Check if we can set the model, if we can, insert newslice
		v := reflect.ValueOf(model).Elem()
		if v.CanSet() {
			v.Set(newS)
		}

		return rows.Close()
	}

	if !sqlbuilder.IsPointer(m) {
		return logger.Log.LogError("you must pass a pointer to a struct for all to function", errors.New("in correct value passed. Must be pointer"))
	}

	row := rq.stmt.QueryRow(rq.args...)
	if err := row.Scan(sqlbuilder.PointerAttributes(reflect.ValueOf(model))...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return logger.Log.LogError("no rows found for raw query call", err)
		}

		return logger.Log.LogError("error occured scanning raw query results", err)
	}

	return nil
}
