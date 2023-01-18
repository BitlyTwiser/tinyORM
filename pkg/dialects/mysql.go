package dialects

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

const DIALECT_TYPE_MYSQL = "mysql"

type Mysql struct {
	db *sql.DB
	mu sync.Mutex
}

var _ DialectHandler = (*Mysql)(nil)

func (m *Mysql) Create(model any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return Create(m.db, model, DIALECT_TYPE_MYSQL)
}

func (m *Mysql) Update(model any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return Update(m.db, model, DIALECT_TYPE_MYSQL)
}

func (m *Mysql) Delete(model any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return Delete(m.db, model, DIALECT_TYPE_MYSQL)
}

func (m *Mysql) BulkDelete(model any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return BulkDelete(m.db, model, DIALECT_TYPE_MYSQL)
}

func (m *Mysql) Find(model any, args ...any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return Find(m.db, model, DIALECT_TYPE_MYSQL, args...)
}

func (m *Mysql) Where(model any, stmt string, limit int, args ...any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return Where(m.db, model, stmt, limit, DIALECT_TYPE_MYSQL, args...)
}

func (m *Mysql) Raw(query string, args ...any) (*RawQuery, error) {
	return Raw(m.db, query, args...)
}

// Alters the database that queries are for.
func (m *Mysql) SetDB(connDB *sql.DB) {
	m.db = connDB
}

func (m *Mysql) QueryString(connInfo DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", connInfo.User, connInfo.Password, connInfo.Host, connInfo.Port, connInfo.Database)
}
