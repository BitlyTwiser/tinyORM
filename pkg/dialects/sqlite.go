package dialects

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
	mu sync.RWMutex
}

var _ DialectHandler = (*SQLite)(nil)

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

// Will accept arbitrary arguments, though only 1 is used, which should be the ID of the object to find.
// If an ID is not passed, ALL objects of the model will be returned
// If an ID IS passed, only a single object should ever be found.
// Uses a Read lock to avoid multiple reads from conflicting
func (s *SQLite) Find(model any, args ...any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return nil
}

// Find an object from database basd off of query params
// Uses a Read lock for concurrent safety
func (s *SQLite) Where(model any, stmt string, limit int, args ...any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return nil
}

// Excutes a raw SQL query. Should be used with caution
// No locks or safety mechanisms in place.
func (s *SQLite) Raw(query string, args ...any) (sql.Result, error) {
	return nil, nil
}

func (s *SQLite) SetDB(connDB *sql.DB) {
	s.db = connDB
}

func (s *SQLite) QueryString(connInfo DBConfig) string {
	return ""
}
