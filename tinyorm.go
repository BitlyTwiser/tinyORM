package tinyorm

import (
	"fmt"

	"github.com/BitlyTwiser/tinyORM/pkg/connections"
	"github.com/BitlyTwiser/tinyORM/pkg/databases"
	"github.com/BitlyTwiser/tinyORM/pkg/logger"
)

// Connect is the primary entrypoint to tinyorm
// Connect will accept variadic values of connection strings i.e. development, prod etc..
// Each connection string must match a string present within the database.yml.
// If a connection fails whilst connecting to any of the given databases, the failure is tracked and reported.
// If all fail, the application will exit
func Connect(connection string) (databases.DatabaseHandler, error) {
	// Default to development
	if connection == "" {
		connection = "development"
	}

	err := connections.InitDatabaseConnection(connection)

	if err != nil {
		return nil, logger.Log.LogError("error initalizing database connection", err)
	}

	if handle, found := connections.Connections[connection]; found {
		return handle, nil
	}

	return nil, logger.Log.LogError("error connecting to database", fmt.Errorf("no database was found in database.yml"))
}

// Will connect to and handle multiple concurrent database connections
// Will accept variadic set of values each string denoting a database connection within the database.yml file.
// i.e. Development, Prod, RO etc..
func MultiConnect(connections ...string) (databases.MultiTenantDatabaseHandler, error) {
	return databases.MultiTenantDatabaseHandler{}, nil
}
