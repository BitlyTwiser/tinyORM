package dialects

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BitlyTwiser/tinyORM/pkg/logger"
	_ "github.com/mattn/go-sqlite3"
)

const DIALECT_TYPE_SQLITE = "sqlite3"

type SQLite struct {
	db     *sql.DB
	mu     sync.Mutex
	config DBConfig
}

var _ DialectHandler = (*SQLite)(nil)

func (s *SQLite) Create(model any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return Create(s.db, model, DIALECT_TYPE_SQLITE)
}

func (s *SQLite) Update(model any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return Update(s.db, model, DIALECT_TYPE_SQLITE)
}

func (s *SQLite) Delete(model any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return Delete(s.db, model, DIALECT_TYPE_SQLITE)
}

func (s *SQLite) BulkDelete(model any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return BulkDelete(s.db, model, DIALECT_TYPE_SQLITE)
}

func (s *SQLite) Find(model any, args ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return Find(s.db, model, DIALECT_TYPE_SQLITE, args...)
}

func (s *SQLite) Where(model any, stmt string, limit int, args ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return Where(s.db, model, stmt, limit, DIALECT_TYPE_SQLITE, args...)
}

func (s *SQLite) Raw(query string, args ...any) (*RawQuery, error) {
	return Raw(s.db, query, args...)
}

func (s *SQLite) SetDB(connDB *sql.DB) {
	s.db = connDB
}

func (s *SQLite) SetConfig(config DBConfig) {
	s.config = config
}

func (s *SQLite) GetConfig() DBConfig {
	return s.config
}

func (s *SQLite) QueryString() string {
	var dsn strings.Builder

	if _, err := os.Stat(s.config.Path); err != nil {
		fullPath, err := filepath.Abs(s.config.Path)
		if err != nil {
			logger.Log.LogError("Could not determine file path", fmt.Errorf("could not build filepath for creating sqlite database"))
			return ""
		}
		logger.Log.LogEvent("warn", "no sqlite database found at given path! Attempting to create database now. (Please note, you will need to manually run migrations and create tables.)", map[string]any{"path": fullPath})
		if os.IsNotExist(err) {
			if _, err := os.Create(fullPath); err != nil {
				logger.Log.LogError("database file not found", fmt.Errorf("could not find sqlitedatabase and could not create the database file within project.. please create the SQLITE database"))
				return ""
			}
		}
	}

	dsn.WriteString("file:")
	dsn.WriteString(s.config.Path)
	// Write back early if not auth as the DSN is complete
	if !s.config.Auth {
		return dsn.String()
	}

	// Complie DSN for encrypting sqlite database
	dsn.WriteString("_auth&")
	dsn.WriteString("_auth_user=")
	dsn.WriteString(s.config.User)
	dsn.WriteString("&_auth_pass=")
	dsn.WriteString(s.config.Password)
	dsn.WriteString("&_auth_crypt=SHA512") // Default to SHA512 encoding

	return dsn.String()
}
