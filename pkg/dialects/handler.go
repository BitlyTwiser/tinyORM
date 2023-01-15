package dialects

import "database/sql"

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
	Port     int    `yaml:"port,omitempty"`
	Host     string `yaml:"host,omitempty"`
	Pool     int    `yaml:"pool,omitempty"`
	Connect  bool   `yaml:"connect,omitempty"`
	Password string `yaml:"password,omitempty"`
	User     string `yaml:"user,omitempty"`
	Database string `yaml:"database,omitempty"`
	Path     string `yaml:"path,omitempty"`
	Dialect  string `yaml:"dialect"`
	Auth     bool   `yaml:"auth"`
}

type MultiTenantDialectHandler struct {
	Handlers []DialectHandler
}
