package databases

import (
	"database/sql"
	"sync"
)

type Mysql struct {
	connections map[string]*sql.DB
	mu          sync.Mutex
	database    string
}

var _ DatabaseHandler = (*Mysql)(nil)

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

func (m *Mysql) Find(model any, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

func (m *Mysql) Where(stmt string, args ...any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

func (m *Mysql) Raw(query string) (any, error) {
	return nil, nil
}

// Alters the database that queries are for.
func (m *Mysql) SetDB(connInfo map[string]*sql.DB) {
	m.connections = connInfo
}

func (m *Mysql) QueryString(connInfo DBConfig) string {
	return ""
}
