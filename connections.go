package dbfusion

import (
	"fmt"
	"sync"

	"github.com/glodb/dbfusion/dbconnections"
	"github.com/glodb/dbfusion/dbfusionErrors"
)

// Factory represents the singleton factory.
type Connections struct {
	// Add fields and methods relevant to your factory here.
	connections map[dbconnections.DBTypes]dbconnections.DBConnections
}

var (
	instance *Connections
	once     sync.Once
)

// GetInstance returns a singleton instance of the Factory.
func GetInstance() *Connections {
	once.Do(func() {
		instance = &Connections{}
	})
	return instance
}

// Supports multiple DB Types
func (c *Connections) GetConnection(option Options) (connection dbconnections.DBConnections, err error) {
	if con, ok := c.connections[option.DbType]; ok {
		connection = con
		return
	}

	if option.Uri == nil {
		return nil, dbfusionErrors.ErrUriRequiredForConnection
	}
	switch option.DbType {
	case dbconnections.MONGO:
		connection = &dbconnections.MongoConnection{}
	case dbconnections.MYSQL:
		connection = &dbconnections.MySql{}
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

	if connection == nil {
		err = dbfusionErrors.ErrDBTypeNotSupported
		return
	}

	return
}

func (c *Connections) CloseConnection(dbType dbconnections.DBTypes) error {
	if con, ok := c.connections[dbType]; ok {
		con.DisConnect()
		return nil
	}
	return dbfusionErrors.ErrConnectionNotAvailable
}

func (c *Connections) CloseAllConnections() {
	for _, value := range c.connections {
		value.DisConnect()
	}
}
