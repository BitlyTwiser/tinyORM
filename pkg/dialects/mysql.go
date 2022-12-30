package dialects

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	db *sql.DB
	mu sync.Mutex
}

var _ DialectHandler = (*Mysql)(nil)

func (m *Mysql) Create(model any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

func (m *Mysql) Update(model any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

func (m *Mysql) Delete(model any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

func (m *Mysql) Find(model any, args ...any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

func (m *Mysql) Where(model any, stmt string, limit int, args ...any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

func (m *Mysql) Raw(query string, args ...any) (sql.Result, error) {
	return nil, nil
}

// Alters the database that queries are for.
func (m *Mysql) SetDB(connDB *sql.DB) {
	m.db = connDB
}

func (m *Mysql) QueryString(connInfo DBConfig) string {
	return ""
}
