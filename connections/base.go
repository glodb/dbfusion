package connections

import (
	"github.com/glodb/dbfusion/caches"
)

// base is the foundational interface in the connections package, providing essential methods
// for changing the active database and setting the cache for database connections.
type base interface {
	// ChangeDatabase allows changing the active database to the one specified by its name.
	// It takes the name of the database as a parameter and returns an error if the database change fails.
	ChangeDatabase(dbName string) error

	// SetCache sets the cache to be used by the database connection.
	// It takes a pointer to a cache object as a parameter and associates it with the database connection.
	SetCache(*caches.Cache)

	// SetPageSize sets the page size to be used by pagination queries.
	// It takes a int as a parameter and associates it with the pagination results.
	SetPageSize(int)
}
