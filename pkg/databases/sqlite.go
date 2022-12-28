package databases

import (
	"database/sql"
	"sync"
)

type SQLite struct {
	connections map[string]*sql.DB
	mu          sync.RWMutex
}

var _ DatabaseHandler = (*SQLite)(nil)

// Create a record in the database
// Uses write lock for concurrent safety
func (s *SQLite) Create(model any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

// Update a value in the database
// Uses write lock for concurrent safety
func (s *SQLite) Update(model any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

// Deletes a value from the database.
// Uses a write lock for concurrent safety
func (s *SQLite) Delete(model any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

// Lookup data via a models ID and the model struct itself.
// Uses a Read lock to avoid multiple reads from conflicting
func (s *SQLite) Find(model any, id string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return nil
}

// Find an object from database basd off of query params
// Uses a Read lock for concurrent safety
func (s *SQLite) Where(stmt string, args ...any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return nil
}

// Excutes a raw SQL query. Should be used with caution
// No locks or safety mechanisms in place.
func (s *SQLite) Raw(query string) (any, error) {
	return nil, nil
}

func (s *SQLite) SetDB(connInfo map[string]*sql.DB) {
	s.connections = connInfo
}

func (s *SQLite) QueryString(connInfo DBConfig) string {
	return ""
}
