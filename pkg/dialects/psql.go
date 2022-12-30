package dialects

import (
	"context"
	"database/sql"
	"fmt"
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

	result, err := pd.db.Exec(query.Query, query.Args...)

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
	return nil
}

func (pd *Postgres) Delete(model any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	return nil
}

// Will accept arbitrary arguments, though only 1 is used, which should be the ID of the object to find.
// If an ID is not passed, ALL objects of the model will be returned
// If an ID IS passed, only a single object should ever be found.
func (pd *Postgres) Find(model any, args ...any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	return nil
}

// Will return all rows found unless <= 1 rows are present in result of query
// Will accept a limit, limit of <= 0 will return all rows found matching the query
// Where is an all in 1 method with no chaining. Pass in the model, statement, desired limit (if there is one, else pass in <= 0), and any arguments to satiate the query
func (pd *Postgres) Where(model any, stmt string, limit int, args ...any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	err := sqlbuilder.QueryAndUpdate("where", model, stmt, limit, args)

	if err != nil {
		return err
	}

	return nil
}

// Just straight up performs a raw query. All work is done by the user, this is just an interface for the ExecContext function.
func (pd *Postgres) Raw(query string, args ...any) (sql.Result, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result, err := pd.db.ExecContext(ctx, query, args)

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
