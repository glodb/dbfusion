package conditions

// SqlData is a concrete implementation of the DBFusionData interface designed for SQL databases.
// It encapsulates the specific SQL query, values, cache key, and cache values required for SQL queries.
type SqlData struct {
	Query       string        // The SQL query.
	Values      []interface{} // The values associated with the SQL query.
	CacheKey    string        // The cache key associated with the query.
	CacheValues string        // The cache values generated through hooks for the query.
}

// GetQuery retrieves the SQL-specific query for the underlying SQL database system.
func (sd *SqlData) GetQuery() interface{} {
	return sd.Query
}

// SetQuery sets the SQL-specific query for the underlying SQL database system.
func (sd *SqlData) SetQuery(query interface{}) {
	sd.Query = query.(string)
}

// GetValues retrieves the values associated with the SQL query.
func (sd *SqlData) GetValues() interface{} {
	return sd.Values
}

// SetValues sets the values associated with the SQL query.
func (sd *SqlData) SetValues(values interface{}) {
	sd.Values = values.([]interface{})
}

// GetCacheKey retrieves the unique cache key associated with the SQL query.
func (sd *SqlData) GetCacheKey() string {
	return sd.CacheKey
}

// SetCacheKey sets the unique cache key for the SQL query, allowing for efficient caching of query results.
func (sd *SqlData) SetCacheKey(cacheKey string) {
	sd.CacheKey = cacheKey
}

// GetCacheValues retrieves the cache values that were generated through hooks for the SQL query.
func (sd *SqlData) GetCacheValues() string {
	return sd.CacheValues
}

// SetCacheValues sets the cache values generated through hooks for the SQL query.
func (sd *SqlData) SetCacheValues(cacheValues string) {
	sd.CacheValues = cacheValues
}
