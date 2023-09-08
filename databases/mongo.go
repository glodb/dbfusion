package databases

import "github.com/glodb/dbfusion/query"

// Define a separate interface for MongoDB-specific queries.
type MongoQuery interface {
	query.Query
	Aggregate(pipeline []interface{}) error
	// Add other MongoDB-specific methods as needed.

	// New method for handling MongoDB aggregation with options.
	AggregateWithOptions(pipeline []interface{}, options interface{}) error

	// New method for counting documents after aggregation.
	CountAfterAggregate() (int64, error)

	// New method for handling MongoDB distinct aggregation.
	DistinctAggregate(field string, query interface{}) ([]interface{}, error)

	// Aggregation pipeline stages and operators.
	Match(filter interface{}) MongoQuery                         // $match
	Project(spec interface{}) MongoQuery                         // $project
	Group(id interface{}, fields interface{}) MongoQuery         // $group
	Unwind(path string) MongoQuery                               // $unwind
	Lookup(from, localField, foreignField, as string) MongoQuery // $lookup
	ReplaceRoot(newRoot interface{}) MongoQuery                  // $replaceRoot
	Sample(size int64) MongoQuery                                // $sample

	// Add other MongoDB aggregation stages and operators as needed.
}
