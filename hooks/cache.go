package hooks

// CacheHook is an interface designed for user-defined models to facilitate efficient caching of associated data.
// By implementing this interface, developers can specify a list of cache indexes as strings. These cache indexes
// serve as keys for storing and retrieving data in a caching system, providing a way to optimize data retrieval
// and reduce the load on the underlying database.
//
// When a model implements the CacheHook interface, it defines which cache indexes should be associated with the
// data represented by that model. These cache indexes can be used to store serialized data in a cache, such as a
// Redis or Memcached cache, allowing for quick access to frequently accessed data.
//
// Implementing CacheHook is particularly useful in scenarios where minimizing database queries and improving
// application performance are essential goals. By leveraging cache indexes, developers can achieve faster data
// retrieval and reduce the overall load on the database, resulting in a more responsive and scalable application.
//
// In summary, CacheHook is a powerful interface that empowers developers to integrate efficient caching
// strategies into their models, ultimately enhancing the performance and responsiveness of their applications.
type CacheHook interface {
	// GetCacheIndexes returns a slice of cache index names, represented as strings. These indexes will be used to
	// store and retrieve data in the cache, allowing for efficient caching of associated data.
	GetCacheIndexes() []string
}
