package connections

import "github.com/glodb/dbfusion/ftypes"

const (
	MONGO = ftypes.DBTypes(1)
	MYSQL = ftypes.DBTypes(2)
)

type PaginationResults struct {
	TotalDocuments int64
	TotalPages     int64
	CurrentPage    int64
	Limit          int64
}
type baseConnections interface {
	base
	Connect(uri string) error
	ConnectWithCertificate(uri string, filePath string) error
	DisConnect() error
}
