package connections

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BitlyTwiser/tinyORM/pkg/databases"
	"github.com/BitlyTwiser/tinyORM/pkg/logger"
	"gopkg.in/yaml.v2"
)

const (
	FILE_RECUR_DEPTH = 5
	databaseFileName = "database.yml"
)

var Connections = make(map[string]databases.DatabaseHandler)

// Initialize database connection via loading the database.yml for the given connection.
// will set the database handlers to the appropriate *sql.DB
func InitDatabaseConnection(dbConnType string) error {
	var db *sql.DB
	var handle databases.DatabaseHandler
	var found bool

	config, err := loadDatabaseConfig(dbConnType)

	if err != nil {
		return err
	}

	connConfig, found := config[dbConnType]

	if !found || connConfig == nil {
		return fmt.Errorf("database connection %s was not found in database.yml. Please check the file", dbConnType)
	}

	if handle, found = databases.Databases[connConfig.Dialect]; !found {
		return fmt.Errorf("please check provided dialect in database.yml. Provided dialect: %v", connConfig.Dialect)
	}

	db, err = sql.Open(connConfig.Dialect, handle.QueryString(*connConfig))
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	// May not need this anymore
	handle.SetDB(map[string]*sql.DB{dbConnType: db})

	Connections[dbConnType] = handle

	return nil
}

func loadDatabaseConfig(dbConnType string) (map[string]*databases.DBConfig, error) {
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

	// Can be an array of connTypes
	config := map[string]*databases.DBConfig{dbConnType: new(databases.DBConfig)}
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

func readDatabaseFile(f io.Reader, config map[string]*databases.DBConfig) error {
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
