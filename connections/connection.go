package connections

import (
	"github.com/glodb/dbfusion/ftypes"
)

// Define an interface for query operations.
type Connection interface {
	crud
	baseConnections
	Paginate(interface{}, int) (PaginationResults, error)
	SetPageSize(int)
	// New methods for bulk operations.
	FindMany(interface{})
	CreateMany([]interface{})
	UpdateMany([]interface{})
	DeleteMany(ftypes.QMap)
}
