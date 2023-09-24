package connections

import "github.com/glodb/dbfusion/ftypes"

//Supported Dbtypes
const (
	MONGO = ftypes.DBTypes(1)
	MYSQL = ftypes.DBTypes(2)
)

// PaginationResults represents the result of a paginated query, providing information about the total number of documents,
// total pages, current page number, and the limit of documents per page.
type PaginationResults struct {
	TotalDocuments int64 // Total number of documents in the query result.
	TotalPages     int64 // Total number of pages based on the pagination criteria.
	CurrentPage    int64 // The current page number being viewed.
	Limit          int64 // The limit of documents displayed per page.
}

// baseConnections is an interface used by various database connection classes to define common methods
// for managing database connections. It extends the base interface, allowing for changing the active database,
// setting the cache, connecting to a database, disconnecting, and connecting with certificate-based authentication.
type baseConnections interface {
	base

	// Connect establishes a connection to the database using the provided URI.
	// It returns an error if the connection cannot be established.
	Connect(uri string) error

	// ConnectWithCertificate establishes a secure connection to the database using the provided URI and certificate file.
	// It takes the URI and the file path to the certificate as parameters and returns an error if the connection cannot be established.
	ConnectWithCertificate(uri string, filePath string) error

	// DisConnect closes the active connection to the database.
	// It returns an error if the disconnection process encounters any issues.
	DisConnect() error
}
