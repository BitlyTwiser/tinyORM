package dialects

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/BitlyTwiser/tinyORM/pkg/sqlbuilder"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
	mu sync.Mutex
}

var _ DialectHandler = (*Postgres)(nil)

func (pd *Postgres) Create(model any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	query := sqlbuilder.QueryBuilder("create", model, "psql")

	if query.Err != nil {
		return query.Err
	}

	stmt, err := pd.db.PrepareContext(context.Background(), query.Query)

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

func (pd *Postgres) Update(model any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	query := sqlbuilder.QueryBuilder("update", model, "psql")

	if query.Err != nil {
		return query.Err
	}

	stmt, err := pd.db.PrepareContext(context.Background(), query.Query)

	if err != nil {
		return err
	}

	id := query.GetModelID()

	// This should have errored earlier in execution, but just in case
	if id == nil {
		return fmt.Errorf("model ID cannot be nil")
	}

	result, err := stmt.Exec(id)

	if err != nil {
		return fmt.Errorf("error deleting database record. Error: %v", err.Error())
	}

	if c, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("error deleting records. Error: %s Rows Affected: %d", err.Error(), c)
	}

	return nil
}

// No id is present within the model and no args are passed, a batch delete will occur.
// To delete records with values other than ID, you can insert a model without an ID, but with other fields present.
// i.e. to delete a user by name: Delete(&User{name: "carl"})
// Without an ID field, but with name present, only "carl" will be deleted
// Multiple attributes will be treated as &'s
func (pd *Postgres) Delete(model any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	data := sqlbuilder.QueryBuilder("delete", model, "psql")

	if data.Err != nil {
		return data.Err
	}

	stmt, err := pd.db.PrepareContext(context.Background(), data.Query)

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
// If an ID IS passed, only a single object should ever be found.
// If an ID is passed, the the model is converted into a slice of model type
func (pd *Postgres) Find(model any, args ...any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	data := sqlbuilder.QueryBuilder("find", model, "psql")

	if data.Err != nil {
		return data.Err
	}

	value := reflect.Indirect(reflect.ValueOf(model))
	if len(args) == 0 {
		// Make sure its a slice.

		if value.Kind() != reflect.Slice {
			return fmt.Errorf("you must pass an slice of model value when not using an ID")
		}
		stmt, err := pd.db.PrepareContext(context.Background(), fmt.Sprintf("SELECT %s FROM %s", sqlbuilder.CoalesceQueryBuilder(value.Type().Elem()), data.TableName))

		if err != nil {
			return err
		}

		rows, err := stmt.Query(args...)

		if err != nil {
			return err
		}

		defer rows.Close()

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

		// Check if we can set the model, if we can, insert newslice
		v := reflect.ValueOf(model).Elem()
		if v.CanSet() {
			v.Set(newS)
		}

		// Ensure rows are closed
		return rows.Close()
	}

	// If not a slice, the find operation is much more simple. We expect args[0] to be the ID we are looking for.
	s := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", sqlbuilder.CoalesceQueryBuilder(value.Type()), data.TableName)
	stmt, err := pd.db.PrepareContext(context.Background(), s)
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
func (pd *Postgres) Where(model any, stmt string, limit int, args ...any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	// Strip attributes from model and build coalesce query
	// stmt and args are passed into coalese query
	// If limit is > 0, limit is also passed
	// if model is a slice, return multiple (as with find)
	// Else we expect to only return the first (pass limit 1 to query, even if limit <= 0)
	return nil
}

// Just straight up performs a raw query. All work is done by the user, this is just an interface for the Exec function.
func (pd *Postgres) Raw(query string, args ...any) (sql.Result, error) {
	stmt, err := pd.db.PrepareContext(context.Background(), query)

	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec(args)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (pd *Postgres) SetDB(connDB *sql.DB) {
	pd.db = connDB
}

func (pd *Postgres) QueryString(c DBConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password =%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.Database, "disable")
}
