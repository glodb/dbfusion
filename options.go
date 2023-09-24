package dbfusion

import (
	"github.com/glodb/dbfusion/caches"
)

// Options is a structure that holds configuration options for connecting to a database.
type Options struct {
	// DbName is a pointer to a string representing the name of the database to connect to.
	DbName *string

	// Uri is a pointer to a string representing the connection URI for the database.
	Uri *string

	// CertificatePath is a pointer to a string representing the file path to SSL/TLS certificate,
	// if required for secure connections. It can be nil if no certificate is needed.
	CertificatePath *string

	// Cache is an instance of a cache that can be associated with the database connection.
	// It allows for caching data to improve query performance.
	Cache caches.Cache
}
