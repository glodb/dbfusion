package connections

import (
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/queryoptions"
)

// Define an interface for query operations.
type Connection interface {
	crud
	baseConnections
	// AddCompositeCacheIndex()
	// AddCompositeIndex()
	Paginate(interface{}, ...queryoptions.FindOptions)
	Distinct(field string)

	RegisterSchema()
	// New methods for bulk operations.
	CreateMany([]interface{})
	UpdateMany([]interface{})
	DeleteMany(ftypes.QMap)
}
