package dbfusion

import (
	"fmt"
	"sync"

	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/implementations"
)

// Factory represents the singleton factory.
type connectionsFactory struct {
	// Add fields and methods relevant to your factory here.
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

// Supports multiple DB Types
func (c *connectionsFactory) getConnection(option Options, dbType ftypes.DBTypes) (connection connections.Connection, err error) {
	if con, ok := c.connections[dbType]; ok {
		connection = con
		return
	}

	if option.Uri == nil {
		return nil, dbfusionErrors.ErrUriRequiredForConnection
	}
	switch dbType {
	case connections.MONGO:
		connection = &implementations.MongoConnection{}
	case connections.MYSQL:
		connection = &implementations.MySql{}
	}

	if option.CertificatePath != nil {
		err = connection.ConnectWithCertificate(*option.Uri, *option.CertificatePath)
	} else {
		err = connection.Connect(*option.Uri)
	}

	if err != nil {
		return nil, fmt.Errorf("failed connection on driver %w", err)
	}

	if option.Cache != nil {
		if option.Cache.IsConnected() {
			connection.SetCache(&option.Cache)
		} else {
			return nil, dbfusionErrors.ErrCacheNotConnected
		}
	}

	if option.DbName != nil {
		connection.ChangeDatabase(*option.DbName)
	} else {
		connection.ChangeDatabase("")
	}
	connection.SetPageSize(DEFAULT_PAGE_SIZE)

	if connection == nil {
		err = dbfusionErrors.ErrDBTypeNotSupported
		return
	}

	return
}
func (c *connectionsFactory) GetMySqlConnection(option Options) (connection connections.SQLConnection, err error) {
	con, err := c.getConnection(option, connections.MYSQL)

	if err != nil {
		return nil, err
	}
	return con.(connections.SQLConnection), nil
}

func (c *connectionsFactory) GeMongoConnection(option Options) (connection connections.MongoConnection, err error) {
	con, err := c.getConnection(option, connections.MONGO)

	if err != nil {
		return nil, err
	}
	return con.(connections.MongoConnection), nil
}

func (c *connectionsFactory) CloseConnection(dbType ftypes.DBTypes) error {
	if con, ok := c.connections[dbType]; ok {
		con.DisConnect()
		return nil
	}
	return dbfusionErrors.ErrConnectionNotAvailable
}

func (c *connectionsFactory) CloseAllConnections() {
	for _, value := range c.connections {
		value.DisConnect()
	}
}
