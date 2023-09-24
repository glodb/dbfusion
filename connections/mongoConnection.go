package connections

// Define a separate interface for MongoDB-specific queries.
type MongoConnection interface {
	Connection

	Match(interface{}) MongoConnection
	Bucket(interface{}) MongoConnection
	BucketsAuto(interface{}) MongoConnection
	AddFields(interface{}) MongoConnection
	GeoNear(interface{}) MongoConnection
	Group(interface{}) MongoConnection
	LimitAggregate(int) MongoConnection
	SkipAggregate(int) MongoConnection
	SortAggregate(interface{}) MongoConnection
	SortByCount(interface{}) MongoConnection
	Project(interface{}) MongoConnection
	Unset(interface{}) MongoConnection
	ReplaceWith(interface{}) MongoConnection
	Merge(interface{}) MongoConnection
	Out(interface{}) MongoConnection
	Facet(interface{}) MongoConnection
	CollStats(interface{}) MongoConnection
	IndexStats(interface{}) MongoConnection
	PlanCacheStats(interface{}) MongoConnection
	Redact(interface{}) MongoConnection
	ReplaceRoot(interface{}) MongoConnection
	ReplaceCount(interface{}) MongoConnection
	Sample(interface{}) MongoConnection
	Set(interface{}) MongoConnection
	Unwind(interface{}) MongoConnection
	Lookup(interface{}) MongoConnection
	GraphLookup(interface{}) MongoConnection
	Count(interface{}) MongoConnection
	Aggregate(interface{}) error
	AggregatePaginate(interface{}, int) (PaginationResults, error)

	Table(tableName string) MongoConnection
	Where(interface{}) MongoConnection
	Skip(skip int64) MongoConnection
	Limit(limit int64) MongoConnection
	Sort(sortKey string, sortdesc ...bool) MongoConnection
	Select(keys map[string]bool) MongoConnection
	CreateIndexes(data interface{}) error

	// Add other MongoDB aggregation stages and operators as needed.
}
