package conditions

// DBFusionData is an interface designed to abstract database-specific query and data handling in the DBFusion framework.
// Implementations of this interface should provide methods to get and set queries, values, cache keys, and cache values
// according to the requirements of the underlying database system.
type DBFusionData interface {
	// GetQuery retrieves the database-specific query for the underlying database system.
	GetQuery() interface{}

	// SetQuery sets the database-specific query for the underlying database system.
	SetQuery(interface{})

	// GetValues retrieves any values associated with the query. For database systems that support passing values to
	// the driver separately (e.g., SQL), this method provides those values.
	GetValues() interface{}

	// SetValues sets the values associated with the query. For database systems that support passing values to
	// the driver separately (e.g., SQL), this method sets those values.
	SetValues(interface{})

	// SetCacheKey sets a unique cache key for the specific query, allowing for efficient caching of query results.
	SetCacheKey(string)

	// GetCacheKey retrieves the unique cache key associated with the query, which can be used for caching purposes.
	GetCacheKey() string

	// SetCacheValues sets the cache values generated through hooks for the query. These cache values can be used
	// to store precomputed results or data.
	SetCacheValues(string)

	// GetCacheValues retrieves the cache values that were set at the time of generation through hooks. These cache
	// values can be used for efficient data retrieval and storage.
	GetCacheValues() string
}
