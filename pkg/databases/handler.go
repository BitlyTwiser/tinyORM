package databases

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
}
