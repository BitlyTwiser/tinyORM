package dialects

import "database/sql"

var Databases map[string]DialectHandler

func init() {
	Databases = map[string]DialectHandler{
		"mysql":    &Mysql{},
		"postgres": &Postgres{},
		"sqlite":   &SQLite{},
	}
}

// DialectHandler is the primary interface that all database types must comply too
type DialectHandler interface {
	Create(model any) error
	Update(model any) error
	Delete(model any) error
	Where(model any, stmt string, limit int, args ...any) error
	Find(model any, args ...any) error
	Raw(query string, args ...any) (sql.Result, error)
	SetDB(connDB *sql.DB)
	QueryString(connInfo DBConfig) string
}

type DBConfig struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	Pool     int    `yaml:"pool"`
	Connect  bool   `yaml:"connect"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
	Dialect  string `yaml:"dialect"`
}

type MultiTenantDialectHandler struct {
	Handlers []DialectHandler
}
