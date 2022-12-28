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

type DB struct {
	Conn *sql.DB
}

// Initialize database connection via loading the database.yml for the given connection.
// will set the database handlers to the appropriate *sql.DB
func InitDatabase(dbConnType string) error {
	var db *sql.DB
	var handle databases.DatabaseHandler
	var found bool

	config, err := loadDatabaseConfig(dbConnType)

	if err != nil {
		return err
	}

	connConfig, found := config[dbConnType]

	if !found || connConfig == nil {
		return logger.Log.LogError("database connectiong not found", fmt.Errorf("database connection %s was not found in database.yml. Please check the file", dbConnType))
	}

	if handle, found = databases.Databases[connConfig.Dialect]; !found {
		return logger.Log.LogError("database dialect not not supported", fmt.Errorf("please check provided dialect in database.yml. Provided dialect: %v", connConfig.Dialect))
	}

	db, err = sql.Open(connConfig.Dialect, handle.QueryString(*connConfig))
	if err != nil {
		return logger.Log.LogError("error connecting", err)
	}

	err = db.Ping()
	if err != nil {
		return logger.Log.LogError("error connecting to database", err)
	}

	handle.SetDB(map[string]*sql.DB{dbConnType: db})

	return nil
}

func loadDatabaseConfig(dbConnType string) (map[string]*databases.DBConfig, error) {
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

	// Can be an array of connTypes
	config := map[string]*databases.DBConfig{dbConnType: new(databases.DBConfig)}
	err = readDatabaseFile(f, config)

	if err != nil {
		return nil, err
	}

	return config, nil
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
