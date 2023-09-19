package connections

// Define a separate interface for MongoDB-specific queries.
type MongoConnection interface {
	Connection
	// Aggregate(pipeline []interface{}) error

	// GetHexFromObjectId()
	// ConvertHexToObjectId()
	// // Add other MongoDB-specific methods as needed.

	// // New method for handling MongoDB aggregation with options.
	// AggregateWithOptions(pipeline []interface{}, options interface{}) error

	// // New method for counting documents after aggregation.
	// CountAfterAggregate() (int64, error)

	// // New method for handling MongoDB distinct aggregation.
	// DistinctAggregate(field string, query interface{}) ([]interface{}, error)

	// // Aggregation pipeline stages and operators.
	// Match(filter interface{}) MongoConnection // $match
	// ProjectAggregate(spec interface{}) MongoQuery                 // $project
	// GroupAggregate(id interface{}, fields interface{}) MongoQuery // $group
	// Unwind(path string) MongoQuery                                // $unwind
	// Lookup(from, localField, foreignField, as string) MongoQuery  // $lookup
	// ReplaceRoot(newRoot interface{}) MongoQuery                   // $replaceRoot
	// Sample(size int64) MongoQuery                                 // $sample

	// Certainly, here are all the top-level aggregation pipeline stages in MongoDB:

	// $addFields:
	// $bucket:
	// $bucketAuto:
	// $collStats:
	// $count:
	// $facet:
	// $geoNear:
	// $graphLookup:
	// $group:
	// $indexStats:
	// $limit:
	// $listLocalSessions:
	// $listSessions:
	// $lookup:
	// $match:
	// $merge:
	// $out:
	// $planCacheStats:
	// $project:
	// $redact:
	// $replaceRoot:
	// $replaceWith:
	// $sample:
	// $set:
	// $skip:
	// $sort:
	// $sortByCount:
	// $unset:
	// $unwind:

	// Deconstructs an array field and generates a separate document for each element in the array.
	Table(tableName string) MongoConnection
	Where(interface{}) MongoConnection
	Skip(skip int64) MongoConnection
	Limit(limit int64) MongoConnection
	Sort(sortKey string, sortdesc ...bool) MongoConnection
	Project(keys map[string]bool) MongoConnection

	// Add other MongoDB aggregation stages and operators as needed.
}
