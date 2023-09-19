package connections

import "github.com/glodb/dbfusion/ftypes"

const (
	MONGO = ftypes.DBTypes(1)
	MYSQL = ftypes.DBTypes(2)
)

type baseConnections interface {
	base
	Connect(uri string) error
	ConnectWithCertificate(uri string, filePath string) error
	DisConnect()
}
