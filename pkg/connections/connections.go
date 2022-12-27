package connections

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var Connections map[string]*DB

type dbConfig struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	Pool     int    `yaml:"pool"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
	Dialect  string `yaml:"dialect"`
}

type DB struct {
	Conn *sql.DB
}

// This will be a variadic function of multiple dbConnTypes (strings)
// If nothing is passed (len(dbConnType == 0 )) then we will map ALL connections in the database.yml file.
func initDatabase(dbConnType string) (*DB, error) {
	var db *sql.DB

	// Default to development
	if dbConnType == "" {
		dbConnType = "development"
	}

	config, err := loadDatabaseConfig(dbConnType)

	if err != nil {
		return nil, err
	}

	for _, connInfo := range config {
		db, err = sql.Open(connInfo.Dialect, connInfo.String())

		if err != nil {
			return nil, err
		}
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database. error %v", err.Error())
	}

	return &DB{Conn: db}, nil
}

// Depends on a database.yml file to be located next to the database.go
// Will parse multiple connections i.e. Development, Production
// Will hold a connection to EACH unless stated as False in config file.
func loadDatabaseConfig(dbConnType string) (map[string]*dbConfig, error) {
	path, err := filepath.Abs("../../database.yml")
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file was not found, please create a database.yml file")
		}
	}

	config := map[string]*dbConfig{dbConnType: new(dbConfig)}
	err = readDatabaseFile(f, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func readDatabaseFile(f io.Reader, config map[string]*dbConfig) error {
	d, err := io.ReadAll(f)

	if err != nil {
		return fmt.Errorf("error reading data from database.yml file.. make sure the file looks correct. error: %v", err.Error())
	}

	err = yaml.Unmarshal(d, &config)

	if err != nil {
		return fmt.Errorf("error parsing fields from database.yml, check file. error: %v", err.Error())
	}

	return nil
}

func (c dbConfig) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s password =%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.Database, "disable")
}
