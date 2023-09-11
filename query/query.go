package query

// Define an interface for query operations.
type Query interface {
	CRUD
	// AddCompositeCacheIndex()
	// AddCompositeIndex()
	Filter(QMap)
	Sort(order interface{}, args ...interface{})
	Paginate(QMap)
	Distinct(field string)

	RegisterSchema()

	// New method for specifying query conditions.
	Where(condition string, args ...interface{})

	// New methods for grouping and ordering.
	GroupBy(keys string)
	OrderBy(order interface{}, args ...interface{})

	// New methods for bulk operations.
	CreateMany([]interface{})
	UpdateMany([]interface{})
	DeleteMany(QMap)

	Skip(skip int64)
	Limit(limit int64)
}
