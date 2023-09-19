package query

import (
	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/joins"
)

// Define an interface for query operations.
type Query interface {
	CRUD
	conditions.ConditionBuilder
	// AddCompositeCacheIndex()
	// AddCompositeIndex()
	Paginate(QMap)
	Distinct(field string)

	RegisterSchema()

	Where(interface{}) Query

	// New methods for grouping and ordering.
	GroupBy(keys string)
	OrderBy(order interface{}, args ...interface{})
	// New methods for bulk operations.
	CreateMany([]interface{})
	UpdateMany([]interface{})
	DeleteMany(QMap)

	Skip(skip int64) Query
	Limit(limit int64) Query
	Sort(map[string]bool) Query
	Project(keys map[string]bool) Query
	Join(join joins.Join) Query
}
