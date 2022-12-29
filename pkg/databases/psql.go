package databases

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

type Postgres struct {
	connections map[string]*sql.DB
	mu          sync.Mutex
}

var _ DatabaseHandler = (*Postgres)(nil)

func (pd *Postgres) Create(model any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()
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
func (pd *Postgres) Find(model any, id string) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	return nil
}

func (pd *Postgres) Where(stmt string, args ...any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	return nil
}

func (pd *Postgres) Raw(query string) (any, error) {
	return nil, nil
}

func (pd *Postgres) SetDB(connInfo map[string]*sql.DB) {
	for k, v := range connInfo {
		pd.connections[k] = v
	}
}

func (pd *Postgres) QueryString(c DBConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password =%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.Database, "disable")
}
