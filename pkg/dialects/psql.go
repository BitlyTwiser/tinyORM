package dialects

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db     *sql.DB
	mu     sync.Mutex
	config DBConfig
}

const DIALECT_TYPE_PSQL = "psql"

var _ DialectHandler = (*Postgres)(nil)

func (pd *Postgres) Create(model any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	return Create(pd.db, model, DIALECT_TYPE_PSQL)
}

func (pd *Postgres) Update(model any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	return Update(pd.db, model, DIALECT_TYPE_PSQL)
}

func (pd *Postgres) Delete(model any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	return Delete(pd.db, model, DIALECT_TYPE_PSQL)
}

func (pd *Postgres) BulkDelete(model any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	return BulkDelete(pd.db, model, DIALECT_TYPE_PSQL)
}

func (pd *Postgres) Find(model any, args ...any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	return Find(pd.db, model, DIALECT_TYPE_PSQL, args...)
}

func (pd *Postgres) Where(model any, stmt string, limit int, args ...any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	return Where(pd.db, model, stmt, limit, DIALECT_TYPE_PSQL, args...)
}

func (pd *Postgres) Raw(query string, args ...any) (*RawQuery, error) {
	return Raw(pd.db, query, args...)
}

func (pd *Postgres) SetDB(connDB *sql.DB) {
	pd.db = connDB
}

func (pd *Postgres) SetConfig(config DBConfig) {
	pd.config = config
}

func (pd *Postgres) GetConfig() DBConfig {
	return pd.config
}

func (pd *Postgres) QueryString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password =%s dbname=%s sslmode=%s", pd.config.Host, pd.config.Port, pd.config.User, pd.config.Password, pd.config.Database, "disable")
}
