package connections

import "github.com/glodb/dbfusion/queryoptions"

// Connection is an interface that combines various functionalities for interacting with a database.
// It extends the following interfaces:
// - crud: Provides methods for common CRUD (Create, Read, Update, Delete) operations on records.
// - baseConnections: Defines methods for managing database connections and cache settings.
type Connection interface {
	crud
	baseConnections
	// Paginate performs pagination on a database query result.
	// It takes an interface representing the query criteria and an integer representing the page number.
	// The method returns PaginationResults, providing information about the total number of documents,
	// total pages, current page number, and the limit of documents per page.
	Paginate(interface{}, int) (PaginationResults, error)

	// FindMany retrieves multiple records from the database based on the provided query criteria.
	// It takes an interface representing the query criteria and optional FindOptions.
	// An error is returned if the operation encounters any issues.
	FindMany(interface{}, ...queryoptions.FindOptions) error

	// InsertMany inserts multiple records into the database.
	// It takes an interface representing a collection of records to be inserted.
	// An error is returned if the insertion fails.
	InsertMany(interface{}) error

	// UpdateMany updates multiple records in the database based on the provided query criteria.
	// It takes two interfaces, one representing the query criteria and the other the update data.
	// The third boolean parameter specifies whether to return the updated records.
	// An error is returned if the operation encounters any issues.
	UpdateMany(interface{}, interface{}, bool) error

	// DeleteMany removes multiple records from the database based on the provided query criteria.
	// It takes one or more interfaces representing the query criteria.
	// An error is returned if the deletion operation fails.
	DeleteMany(...interface{}) error
}
