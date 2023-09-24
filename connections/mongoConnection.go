package connections

// MongoConnection is an interface that extends the base Connection interface and provides
// methods specific to MongoDB database interactions. It allows building and executing MongoDB
// aggregation pipelines, specifying query criteria, sorting, limiting, and more.
type MongoConnection interface {
	Connection

	// Match specifies a $match stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the match criteria and returns the modified MongoConnection.
	Match(interface{}) MongoConnection

	// Bucket specifies a $bucket stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the bucket stage criteria and returns the modified MongoConnection.
	Bucket(interface{}) MongoConnection

	// BucketsAuto specifies a $bucketAuto stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the bucketAuto stage criteria and returns the modified MongoConnection.
	BucketsAuto(interface{}) MongoConnection

	// AddFields specifies a $addFields stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the fields to be added and returns the modified MongoConnection.
	AddFields(interface{}) MongoConnection

	// GeoNear specifies a $geoNear stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the geoNear stage criteria and returns the modified MongoConnection.
	GeoNear(interface{}) MongoConnection

	// Group specifies a $group stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the group stage criteria and returns the modified MongoConnection.
	Group(interface{}) MongoConnection

	// LimitAggregate specifies a $limit stage in a MongoDB aggregation pipeline.
	// It takes an integer representing the maximum number of documents to output and returns the modified MongoConnection.
	LimitAggregate(int) MongoConnection

	// SkipAggregate specifies a $skip stage in a MongoDB aggregation pipeline.
	// It takes an integer representing the number of documents to skip and returns the modified MongoConnection.
	SkipAggregate(int) MongoConnection

	// SortAggregate specifies a $sort stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the sorting criteria and returns the modified MongoConnection.
	SortAggregate(interface{}) MongoConnection

	// SortByCount specifies a $sortByCount stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the sortByCount stage criteria and returns the modified MongoConnection.
	SortByCount(interface{}) MongoConnection

	// Project specifies a $project stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the projected fields and returns the modified MongoConnection.
	Project(interface{}) MongoConnection

	// Unset specifies a $unset stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the fields to unset and returns the modified MongoConnection.
	Unset(interface{}) MongoConnection

	// ReplaceWith specifies a $replaceWith stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the replacement document and returns the modified MongoConnection.
	ReplaceWith(interface{}) MongoConnection

	// Merge specifies a $merge stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the merge stage criteria and returns the modified MongoConnection.
	Merge(interface{}) MongoConnection

	// Out specifies a $out stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the output stage criteria and returns the modified MongoConnection.
	Out(interface{}) MongoConnection

	// Facet specifies a $facet stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the facet stage criteria and returns the modified MongoConnection.
	Facet(interface{}) MongoConnection

	// CollStats specifies a $collStats stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the collStats stage criteria and returns the modified MongoConnection.
	CollStats(interface{}) MongoConnection

	// IndexStats specifies a $indexStats stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the indexStats stage criteria and returns the modified MongoConnection.
	IndexStats(interface{}) MongoConnection

	// PlanCacheStats specifies a $planCacheStats stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the planCacheStats stage criteria and returns the modified MongoConnection.
	PlanCacheStats(interface{}) MongoConnection

	// Redact specifies a $redact stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the redact stage criteria and returns the modified MongoConnection.
	Redact(interface{}) MongoConnection

	// ReplaceRoot specifies a $replaceRoot stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the replaceRoot stage criteria and returns the modified MongoConnection.
	ReplaceRoot(interface{}) MongoConnection

	// ReplaceCount specifies a $replaceCount stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the replaceCount stage criteria and returns the modified MongoConnection.
	ReplaceCount(interface{}) MongoConnection

	// Sample specifies a $sample stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the sample stage criteria and returns the modified MongoConnection.
	Sample(interface{}) MongoConnection

	// Set specifies a $set stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the set stage criteria and returns the modified MongoConnection.
	Set(interface{}) MongoConnection

	// Unwind specifies a $unwind stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the unwind stage criteria and returns the modified MongoConnection.
	Unwind(interface{}) MongoConnection

	// Lookup specifies a $lookup stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the lookup stage criteria and returns the modified MongoConnection.
	Lookup(interface{}) MongoConnection

	// GraphLookup specifies a $graphLookup stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the graphLookup stage criteria and returns the modified MongoConnection.
	GraphLookup(interface{}) MongoConnection

	// Count specifies a $count stage in a MongoDB aggregation pipeline.
	// It takes an interface representing the count stage criteria and returns the modified MongoConnection.
	Count(interface{}) MongoConnection

	// Aggregate executes a MongoDB aggregation pipeline.
	// It takes an interface representing the aggregation pipeline and returns an error if the aggregation fails.
	Aggregate(interface{}) error

	// AggregatePaginate executes a MongoDB aggregation pipeline with pagination.
	// It takes an interface representing the aggregation pipeline and an integer representing the page number.
	// The method returns PaginationResults, providing information about the total number of documents,
	// total pages, current page number, and the limit of documents per page.
	AggregatePaginate(interface{}, int) (PaginationResults, error)

	// Table specifies the MongoDB collection (table) to query.
	// It takes the name of the collection as a parameter and returns the modified MongoConnection.
	Table(tableName string) MongoConnection

	// Where specifies the criteria for filtering documents in the MongoDB collection.
	// It takes an interface representing the filter criteria and returns the modified MongoConnection.
	Where(interface{}) MongoConnection

	// Skip specifies the number of documents to skip in the result set.
	// It takes an integer representing the number of documents to skip and returns the modified MongoConnection.
	Skip(skip int64) MongoConnection

	// Limit specifies the maximum number of documents to return in the result set.
	// It takes an integer representing the limit and returns the modified MongoConnection.
	Limit(limit int64) MongoConnection

	// Sort specifies the sorting order of the result set based on a specified key.
	// It takes the key for sorting and an optional boolean indicating descending order.
	// It returns the modified MongoConnection.
	Sort(sortKey string, sortDesc ...bool) MongoConnection

	// Select specifies the fields to include or exclude in the result set.
	// It takes a map where keys represent field names and values represent inclusion/exclusion flags.
	// It returns the modified MongoConnection.
	Select(keys map[string]bool) MongoConnection

	// CreateIndexes creates one or more indexes in the MongoDB collection based on the provided data.
	// It takes an interface representing index creation data and returns an error if the operation fails.
	CreateIndexes(data interface{}) error
}
