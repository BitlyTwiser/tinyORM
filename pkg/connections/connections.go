package connections

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BitlyTwiser/tinyORM/pkg/dialects"
	"github.com/BitlyTwiser/tinyORM/pkg/logger"
	"gopkg.in/yaml.v2"
)

const (
	FILE_RECUR_DEPTH = 5
	databaseFileName = "database.yml"
	maxIdleTime      = "maxIdleTime"
	maxLifetime      = "maxLifetime"
	maxIdleConn      = "maxIdleConn"
	maxOpenConn      = "maxOpenConn"
)

var Connections = make(map[string]dialects.DialectHandler)

// Initialize database connection via loading the database.yml for the given connection.
// will set the database handlers to the appropriate *sql.DB
func InitDatabaseConnection(dbConnType string) error {
	var db *sql.DB
	var handle dialects.DialectHandler
	var found bool

	config, err := loadDatabaseConfig(dbConnType)

	if err != nil {
		return err
	}

	connConfig, found := config[dbConnType]

	if !found || connConfig == nil {
		return fmt.Errorf("database connection %s was not found in database.yml. Please check the file", dbConnType)
	}

	if handle, found = dialects.Databases[connConfig.Dialect]; !found {
		return fmt.Errorf("please check provided dialect in database.yml. Provided dialect: %v", connConfig.Dialect)
	}

	handle.SetConfig(*connConfig)

	db, err = sql.Open(connConfig.Dialect, handle.QueryString())
	if err != nil {
		return err
	}

	// Set all of the maximums for idle/open connections
	setConnectionDefaults(db, connConfig)

	err = db.Ping()
	if err != nil {
		return err
	}

	// Set handler data
	handle.SetDB(db)

	// Store handler in connections in case of switching handlers
	Connections[dbConnType] = handle

	return nil
}

func setConnectionDefaults(db *sql.DB, config *dialects.DBConfig) {
	if config.MaxIdleConn != 0 {
		db.SetMaxIdleConns(config.MaxIdleConn)
	}

	if config.MaxIdleTime != 0 {
		db.SetConnMaxIdleTime(config.MaxIdleTime)
	}

	if config.MaxLifetime != 0 {
		db.SetConnMaxLifetime(config.MaxLifetime)
	}

	if config.MaxOpenConn != 0 {
		db.SetMaxOpenConns(config.MaxOpenConn)
	}
}

func loadDatabaseConfig(dbConnType string) (map[string]*dialects.DBConfig, error) {
	path := findDatabaseFilePath(databaseFileName, 0)

	if path == "" {
		return nil, logger.Log.LogError("database file not found", fmt.Errorf("could not find database file within project.. please create the database.yml"))
	}

	f, err := os.Open(path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file was not found, please create a database.yml file")
		}
	}

	config := map[string]*dialects.DBConfig{dbConnType: new(dialects.DBConfig)}
	err = readDatabaseFile(f, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

// Iterates several levels up to attempt to locate database.yml file.
func findDatabaseFilePath(path string, level int) string {
	if level == FILE_RECUR_DEPTH {
		return ""
	}
	if _, err := os.Stat(path); err != nil {
		level++

		return findDatabaseFilePath(filepath.Join("../", path), level)
	}

	return path
}

func readDatabaseFile(f io.Reader, config map[string]*dialects.DBConfig) error {
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
