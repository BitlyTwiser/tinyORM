package databases

import "database/sql"

var Databases map[string]DatabaseHandler

func init() {
	Databases = map[string]DatabaseHandler{
		"mysql":    &Mysql{},
		"postgres": &Postgres{},
		"sqlite":   &SQLite{},
	}
}

// DatabaseHandler is the primary interface that all database types must comply too
type DatabaseHandler interface {
	Create(model any) error
	Update(model any) error
	Delete(model any) error
	Where(stmt string, args ...any) error
	Find(model any, id string) error
	Raw(query string) (any, error)
	SetDB(connInfo map[string]*sql.DB)
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

type MultiTenantDatabaseHandler struct {
	Handlers []DatabaseHandler
}
