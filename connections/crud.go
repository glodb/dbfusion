package connections

import (
	"github.com/glodb/dbfusion/queryoptions"
)

// crud is an interface that defines common database operations for creating, reading, updating, and deleting records.
// It provides methods for inserting one record, finding one record with optional query options, updating and finding one record,
// and deleting one or more records.
type crud interface {
	// InsertOne inserts a single record into the database.
	// It takes an interface representing the record to be inserted and returns an error if the insertion fails.
	InsertOne(interface{}) error

	// FindOne retrieves a single record from the database based on the provided query criteria.
	// It takes an interface representing the query criteria and optional FindOptions.
	// If a matching record is found, it populates the provided interface with the result.
	// An error is returned if the operation encounters any issues or if no matching record is found.
	FindOne(interface{}, ...queryoptions.FindOptions) error

	// UpdateAndFindOne updates a single record in the database based on the provided query criteria.
	// It takes two interfaces, one representing the query criteria and the other the update data.
	// The third boolean parameter specifies whether to return the updated record.
	// If the update is successful and a matching record is found, the updated record is populated
	// in the first interface parameter (if the third parameter is true).
	// An error is returned if the operation encounters any issues or if no matching record is found.
	UpdateAndFindOne(interface{}, interface{}, bool) error

	// DeleteOne removes a single record from the database based on the provided query criteria.
	// It takes one or more interfaces representing the query criteria.
	// An error is returned if the deletion operation fails or if no matching record is found.
	DeleteOne(...interface{}) error
}
