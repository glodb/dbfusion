package connections

import (
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/queryoptions"
)

// Define an interface for query operations.
type Connection interface {
	crud
	baseConnections
	Paginate(interface{}, ...queryoptions.FindOptions)
	RegisterSchema()
	// New methods for bulk operations.
	FindMany(interface{})
	CreateMany([]interface{})
	UpdateMany([]interface{})
	DeleteMany(ftypes.QMap)
}
