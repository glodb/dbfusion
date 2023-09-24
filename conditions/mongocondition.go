package conditions

import "go.mongodb.org/mongo-driver/bson/primitive"

// MongoData is a concrete implementation of the DBFusionData interface designed for MongoDB.
// It encapsulates the specific query, cache key, and cache values required for MongoDB queries.
type MongoData struct {
	Query       primitive.D // The MongoDB query.
	CacheKey    string      // The cache key associated with the query.
	CacheValues string      // The cache values generated through hooks for the query.
}

// GetQuery retrieves the MongoDB-specific query for the underlying MongoDB database system.
func (md *MongoData) GetQuery() interface{} {
	if md.Query == nil {
		md.Query = primitive.D{}
	}
	return md.Query
}

// SetQuery sets the MongoDB-specific query for the underlying MongoDB database system.
func (md *MongoData) SetQuery(query interface{}) {
	md.Query = query.(primitive.D)
}

// GetValues retrieves any values associated with the MongoDB query.
// For MongoDB, values are not typically passed separately, so this method returns nil.
func (md *MongoData) GetValues() interface{} {
	return nil
}

// SetValues sets any values associated with the MongoDB query.
// For MongoDB, values are not typically passed separately, so this method does nothing.
func (md *MongoData) SetValues(data interface{}) {
}

// GetCacheKey retrieves the unique cache key associated with the MongoDB query.
func (md *MongoData) GetCacheKey() string {
	return md.CacheKey
}

// SetCacheKey sets the unique cache key for the MongoDB query, allowing for efficient caching of query results.
func (md *MongoData) SetCacheKey(cacheKey string) {
	md.CacheKey = cacheKey
}

// GetCacheValues retrieves the cache values that were generated through hooks for the MongoDB query.
func (md *MongoData) GetCacheValues() string {
	return md.CacheValues
}

// SetCacheValues sets the cache values generated through hooks for the MongoDB query.
func (md *MongoData) SetCacheValues(values string) {
	md.CacheValues = values
}
