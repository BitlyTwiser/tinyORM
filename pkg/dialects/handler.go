package dialects

import (
	"database/sql"
	"time"
)

var Databases map[string]DialectHandler

func init() {
	Databases = map[string]DialectHandler{
		"mysql":    &Mysql{},
		"postgres": &Postgres{},
		"sqlite3":  &SQLite{},
	}
}

// DialectHandler is the primary interface that all database types must comply too
type DialectHandler interface {
	Create(model any) error
	Update(model any) error
	Delete(model any) error
	Where(model any, stmt string, limit int, args ...any) error
	Find(model any, args ...any) error
	Raw(query string, args ...any) (*RawQuery, error)
	SetDB(connDB *sql.DB)
	QueryString(connInfo DBConfig) string
}

type DBConfig struct {
	Port        int           `yaml:"port,omitempty"`
	Host        string        `yaml:"host,omitempty"`
	Pool        int           `yaml:"pool,omitempty"`
	Connect     bool          `yaml:"connect,omitempty"`
	Password    string        `yaml:"password,omitempty"`
	User        string        `yaml:"user,omitempty"`
	Database    string        `yaml:"database,omitempty"`
	Path        string        `yaml:"path,omitempty"`
	Dialect     string        `yaml:"dialect"`
	Auth        bool          `yaml:"auth"`
	MaxIdleTime time.Duration `yaml:"maxIdleTime,omitempty"`
	MaxLifetime time.Duration `yaml:"maxLifetime,omitempty"`
	MaxIdleConn int           `yaml:"maxIdleConn,omitempty"`
	MaxOpenConn int           `yaml:"maxOpenConn,omitempty"`
}

type MultiTenantDialectHandler struct {
	Handlers map[string]DialectHandler
}

// Append will add database handlers to the Handlers slice
func (mtd *MultiTenantDialectHandler) Set(key string, handler DialectHandler) {
	mtd.Handlers[key] = handler
}

// Empty will determine if there are not database handlers present
func (mtd MultiTenantDialectHandler) Empty() bool {
	return len(mtd.Handlers) == 0
}

// Switch allows the caller to alter to different databases to perform executions again
func (mtd MultiTenantDialectHandler) SwitchDB(database string) DialectHandler {
	if db, found := mtd.Handlers[database]; found {
		return db
	}

	return nil
}
