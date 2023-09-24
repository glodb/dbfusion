package dbfusion

import (
	"fmt"
	"sync"

	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/implementations"
)

// Factory represents the singleton lazy factory.
type connectionsFactory struct {
	//Connections saves all the db connections called once
	connections map[ftypes.DBTypes]connections.Connection
}

var (
	instance *connectionsFactory
	once     sync.Once
)

// GetInstance returns a singleton instance of the Factory.
func GetInstance() *connectionsFactory {
	once.Do(func() {
		instance = &connectionsFactory{}
	})
	return instance
}

// getConnection creates and returns a database connection based on the provided options and database type.
// It supports multiple DB Types and returns the appropriate connection instance for the specified type.
// If the connection for the given type already exists, it retrieves and returns that connection.
// Parameters:
//   - option: Options containing connection details, including URI, database name, certificate path, and cache.
//   - dbType: The type of the database to connect to, such as MongoDB or MySQL.
// Returns:
//   - connection: A Connection interface representing the established database connection.
//   - err: An error if any occurred during the connection setup.
func (c *connectionsFactory) getConnection(option Options, dbType ftypes.DBTypes) (connection connections.Connection, err error) {
	// Check if a connection of the specified database type already exists in the connections map.
	if con, ok := c.connections[dbType]; ok {
		// If it exists, assign the existing connection to the 'connection' variable and return.
		connection = con
		return
	}

	// Check if the URI is provided in the options. It is required for establishing a database connection.
	if option.Uri == nil {
		return nil, dbfusionErrors.ErrUriRequiredForConnection
	}

	// Depending on the specified database type, create a new connection instance.
	switch dbType {
	case connections.MONGO:
		connection = &implementations.MongoConnection{}
	case connections.MYSQL:
		connection = &implementations.MySql{}
	}

	// If a certificate path is provided, attempt to connect using the certificate.
	if option.CertificatePath != nil {
		err = connection.ConnectWithCertificate(*option.Uri, *option.CertificatePath)
	} else {
		// Otherwise, connect using the URI.
		err = connection.Connect(*option.Uri)
	}

	// Check if any error occurred during the connection attempt.
	if err != nil {
		return nil, fmt.Errorf("failed connection on driver %w", err)
	}

	// If a cache instance is provided in the options and it is connected, set it for the connection.
	if option.Cache != nil {
		if option.Cache.IsConnected() {
			connection.SetCache(&option.Cache)
		} else {
			return nil, dbfusionErrors.ErrCacheNotConnected
		}
	}

	// If a database name is provided in the options, set it for the connection; otherwise, leave it empty.
	if option.DbName != nil {
		connection.ChangeDatabase(*option.DbName)
	} else {
		connection.ChangeDatabase("")
	}

	// Set the default page size for pagination on the connection.
	connection.SetPageSize(DEFAULT_PAGE_SIZE)

	// If the 'connection' variable is still nil, it means the specified DB type is not supported.
	if connection == nil {
		err = dbfusionErrors.ErrDBTypeNotSupported
		return
	}

	// Return the established connection and any error encountered during the process.
	return
}

// GetMySqlConnection creates and returns a MySQL database connection based on the provided options.
// It uses the 'getConnection' function from the factory to establish the connection.
// If successful, it returns a SQLConnection instance specific to MySQL.
// If an error occurs during the connection process, it returns an error.
func (c *connectionsFactory) GetMySqlConnection(option Options) (connection connections.SQLConnection, err error) {
	// Attempt to establish a connection using the 'getConnection' function for the MySQL database type.
	con, err := c.getConnection(option, connections.MYSQL)

	// Check if an error occurred during the connection attempt.
	if err != nil {
		return nil, err
	}

	// If successful, cast the generic 'con' to a SQLConnection specific to MySQL and return it.
	return con.(connections.SQLConnection), nil
}

// GetMongoConnection creates and returns a MongoDB database connection based on the provided options.
// It utilizes the 'getConnection' function from the factory to establish the connection for MongoDB.
// If the connection is successfully established, it returns a MongoConnection instance specific to MongoDB.
// If any error occurs during the connection process, it returns an error.
func (c *connectionsFactory) GetMongoConnection(option Options) (connection connections.MongoConnection, err error) {
	// Attempt to establish a connection using the 'getConnection' function for the MongoDB database type.
	con, err := c.getConnection(option, connections.MONGO)

	// Check if an error occurred during the connection attempt.
	if err != nil {
		return nil, err
	}

	// If the connection is successful, cast the generic 'con' to a MongoConnection specific to MongoDB and return it.
	return con.(connections.MongoConnection), nil
}

// CloseConnection attempts to close an active database connection of the specified DB type.
// It takes the 'dbType' parameter, which represents the type of database to be closed.
// If a valid connection for the specified DB type is found in the factory's connections map,
// it invokes the 'DisConnect' method on that connection to gracefully close it.
// If no connection is found for the specified DB type, it returns an error indicating that
// the connection is not available.
//
// Parameters:
//   - dbType: The type of database (e.g., MongoDB, MySQL) for which the connection should be closed.
//
// Returns:
//   - error: An error, if any, encountered during the connection closure process.
//            Returns 'ErrConnectionNotAvailable' if no connection is found for the specified DB type.
func (c *connectionsFactory) CloseConnection(dbType ftypes.DBTypes) error {
	// Check if a valid connection exists for the specified DB type in the connections map.
	if con, ok := c.connections[dbType]; ok {
		// Invoke the 'DisConnect' method on the found connection to gracefully close it.
		con.DisConnect()
		return nil
	}

	// If no connection is found for the specified DB type, return an error indicating that the connection is not available.
	return dbfusionErrors.ErrConnectionNotAvailable
}

// CloseAllConnections gracefully closes all active database connections managed by the connections factory.
// It iterates through the connections stored in the connections map and invokes the 'DisConnect' method
// on each connection to ensure that all connections are closed properly.
//
// This function is useful for cleaning up resources and terminating all database connections
// when they are no longer needed, such as when shutting down the application.
//
// Note: It does not return any error because the goal is to close connections without interruption.
func (c *connectionsFactory) CloseAllConnections() {
	// Iterate through the connections stored in the connections map.
	for _, value := range c.connections {
		// Invoke the 'DisConnect' method on each connection to gracefully close it.
		value.DisConnect()
	}
}
