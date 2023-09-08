package dbconnections

import (
	"github.com/glodb/dbfusion/query"
)

type DBTypes int

const (
	MONGO = DBTypes(1)
	MYSQL = DBTypes(2)
)

type DBConnections interface {
	DBBase
	query.Query
	Connect(uri string) error
	ConnectWithCertificate(uri string, filePath string) error
	DisConnect()
}
