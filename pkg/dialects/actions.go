// Actions is all helper functions that are called within the dialects
package dialects

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/BitlyTwiser/tinyORM/pkg/logger"
	"github.com/BitlyTwiser/tinyORM/pkg/sqlbuilder"
)

func Create(db *sql.DB, model any, dialectType string) error {
	query := sqlbuilder.QueryBuilder("create", model, dialectType)

	if query.Err != nil {
		return query.Err
	}

	stmt, err := db.PrepareContext(context.Background(), query.Query)

	if err != nil {
		return fmt.Errorf("error creating database record. error: %s", err.Error())
	}

	result, err := stmt.Exec(query.Args...)

	if err != nil {
		return fmt.Errorf("error creating database record. Error: %v", err.Error())
	}

	if c, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("error creating records. Error: %s Rows Affected: %d", err.Error(), c)
	}

	return nil
}

func Update(db *sql.DB, model any, dialectType string) error {
	query := sqlbuilder.QueryBuilder("update", model, dialectType)

	if query.Err != nil {
		return query.Err
	}

	stmt, err := db.PrepareContext(context.Background(), query.Query)

	if err != nil {
		return err
	}

	id := query.GetModelID()

	// This should have errored earlier in execution, but just in case
	if id == nil {
		return fmt.Errorf("model ID cannot be nil when calling update. Attempt to use a Raw query to update this model")
	}

	query.Args = append(query.Args[1:], id)

	result, err := stmt.Exec(query.Args...)

	if err != nil {
		return fmt.Errorf("error deleting database record. Error: %v", err.Error())
	}

	if c, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("error deleting records. Error: %s Rows Affected: %d", err.Error(), c)
	}

	return nil
}

// No id is present within the model and no args are passed, the FIRST recored (using limit 1) from the given model will be deleted
// To delete results in bulk, pass in a slice. This will batch delete records, use with caution.
// To delete records with values other than ID, you can insert a model without an ID, but with other fields present.
// i.e. to delete a user by name: Delete(&User{name: "carl"})
// Without an ID field, but with name present, only "carl" will be deleted
// Multiple attributes will be treated as &'s
func Delete(db *sql.DB, model any, dialectType string) error {
	data := sqlbuilder.QueryBuilder("delete", model, dialectType)

	if data.Err != nil {
		return data.Err
	}

	m := reflect.Indirect(reflect.ValueOf(model))
	if m.Kind() == reflect.Slice {
		stmt, err := db.PrepareContext(context.Background(), fmt.Sprintf("DELETE FROM %s", data.TableName))
		if err != nil {
			return err
		}
		result, err := stmt.Exec()
		if err != nil {
			return err
		}

		rowsDeleted, err := result.RowsAffected()
		if err != nil {
			return err
		}

		logger.Log.LogEvent("info", "Deleted rows", "rows deleted", rowsDeleted)

		return nil
	}

	// Means that no attributes were found on model and model is not a slice.
	if data.Query == "" {
		logger.Log.LogEvent("info", "no records were found for the delete query")
		return nil
	}

	stmt, err := db.PrepareContext(context.Background(), data.Query)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(data.Args...)
	if err != nil {
		return fmt.Errorf("error deleting database record. Error: %v", err.Error())
	}

	if c, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("error deleting records. Error: %s Rows Affected: %d", err.Error(), c)
	}

	return nil
}

// Will accept arbitrary arguments, though only 1 is used, which should be the ID of the object to find.
// If an ID is not passed, ALL objects of the model will be returned
// If there is no id and the passed model is not a slice, the first row is returned for the given model
// If an ID IS passed, only a single object should ever be found.
// If an ID is passed, the the model is converted into a slice of model type
func Find(db *sql.DB, model any, dialectType string, args ...any) error {
	var querySymbol string = "?"
	data := sqlbuilder.QueryBuilder("find", model, dialectType)

	if data.Err != nil {
		return data.Err
	}

	// If value is not slice kind and args == 0
	value := reflect.Indirect(reflect.ValueOf(model))
	if len(args) == 0 && value.Kind() == reflect.Slice {
		// Make sure its a slice.

		stmt, err := db.PrepareContext(context.Background(), fmt.Sprintf("SELECT %s FROM %s", sqlbuilder.CoalesceQueryBuilder(value.Type().Elem()), data.TableName))

		if err != nil {
			return err
		}

		rows, err := stmt.Query(args...)

		if err != nil {
			return err
		}

		defer func() {
			if err := rows.Close(); err != nil {
				logger.Log.LogError("error closing datase rows in Find call.", err)
			}
		}()

		//Make new slice to feed into the incoming model slice
		newS := reflect.MakeSlice(reflect.SliceOf(value.Type().Elem()), 0, 0)

		for rows.Next() {
			// Create new pointer to inner struct type
			newVal := reflect.New(value.Type().Elem())

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

		// Ensure rows are closed
		return rows.Close()
	}

	// If no args passed and no slice passed, return first value
	if len(args) == 0 && value.Kind() != reflect.Slice {
		s := fmt.Sprintf("SELECT %s FROM %s LIMIT 1", sqlbuilder.CoalesceQueryBuilder(value.Type()), data.TableName)
		stmt, err := db.PrepareContext(context.Background(), s)

		if err != nil {
			return err
		}

		row := stmt.QueryRow()

		if err := row.Scan(data.ModelAttributes()...); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("no rows found for %s", data.TableName)
			}

			return fmt.Errorf("error records for table: %s. Error: %v", data.TableName, err.Error())
		}

		return nil
	}

	// If not a slice and args are passed, the find operation is much more simple. We expect args[0] to be the ID we are looking for.
	if dialectType == "psql" {
		querySymbol = "$1"
	}
	s := fmt.Sprintf("SELECT %s FROM %s WHERE id = %s", sqlbuilder.CoalesceQueryBuilder(value.Type()), data.TableName, querySymbol)
	stmt, err := db.PrepareContext(context.Background(), s)
	if err != nil {
		return err
	}

	row := stmt.QueryRow(args[0])
	if err := row.Scan(data.ModelAttributes()...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no rows found for id: %v", args[0])
		}

		return fmt.Errorf("error scanning rows for id: %v. Error: %v", args[0], err.Error())
	}

	return nil
}

// Will return all rows found unless <= 1 rows are present in result of query
// Will accept a limit, limit of <= 0 will return all rows found matching the query
// Where is an all in 1 method with no chaining. Pass in the model, statement, desired limit (if there is one, else pass in <= 0), and any arguments to satiate the query
func Where(db *sql.DB, model any, stmt string, limit int, dialectType string, args ...any) error {
	var parsedStmt strings.Builder

	if stmt == "" {
		return errors.New("you cannot pass an empty statement")
	}

	if stmt != "" && len(args) == 0 {
		return errors.New("you must provide attributes for the sql query")
	}

	data := sqlbuilder.QueryBuilder("where", model, dialectType)

	if data.Err != nil {
		return data.Err
	}

	i := 1
	for _, v := range stmt {
		// Handle PSQL conversion
		if v == '?' && dialectType == "psql" {
			parsedStmt.WriteString("$" + strconv.Itoa(i))
			i++
			continue
		}

		parsedStmt.WriteRune(v)
	}

	value := reflect.Indirect(reflect.ValueOf(model))
	// If slice we will scan rows and insert data based off of incoming stmt
	if value.Kind() == reflect.Slice {
		query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", sqlbuilder.CoalesceQueryBuilder(value.Type().Elem()), data.TableName, parsedStmt.String())
		if limit > 0 {
			query = query + fmt.Sprintf(" LIMIT %d", limit)
		}
		s, err := db.PrepareContext(context.Background(), query)

		if err != nil {
			return err
		}

		rows, err := s.Query(args...)

		if err != nil {
			return nil
		}

		defer func() {
			if err := rows.Close(); err != nil {
				logger.Log.LogError("error closing database rows.", err)
			}
		}()

		newS := reflect.MakeSlice(reflect.SliceOf(value.Type().Elem()), 0, 0)
		for rows.Next() {
			newVal := reflect.New(value.Type().Elem())

			if err := rows.Scan(sqlbuilder.PointerAttributes(newVal)...); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("there are no rows for for this query")
				}
				return err
			}

			newS = reflect.Append(newS, newVal.Elem())
		}

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

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", sqlbuilder.CoalesceQueryBuilder(value.Type()), data.TableName, parsedStmt.String())
	if limit > 0 {
		query = query + fmt.Sprintf(" LIMIT %d", limit)
	}
	s, err := db.PrepareContext(context.Background(), query)

	if err != nil {
		return err
	}

	// If not slice, scan row
	row := s.QueryRow(args...)

	if err := row.Scan(data.ModelAttributes()...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no rows found for table name %s", data.TableName)
		}
	}

	return nil
}

// Raw builds a raw query, allowing for a user to either call Exec or All functions to perform execution
func Raw(db *sql.DB, query string, args ...any) (*RawQuery, error) {
	stmt, err := db.PrepareContext(context.Background(), query)

	if err != nil {
		return nil, err
	}

	if len(args) > 0 {
		return &RawQuery{
			stmt:  stmt,
			query: query,
			args:  args,
		}, nil
	}

	return &RawQuery{
		stmt:  stmt,
		query: query,
	}, nil
}
