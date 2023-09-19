package caches

/*
Cache Interface

The Cache interface defines a contract for all caches that intend to be used with dbFusion. Any cache implementation must implement this interface to ensure compatibility with dbFusion's caching system.

Methods:

1. ConnectCache(connectionUri string, password ...string) error
   - Description: Establishes a connection to the cache using the provided connection URI and optional password.
   - Parameters:
     - `connectionUri` (string): The URI used to connect to the cache.
     - `password` (string, optional): A password, if required for cache authentication.
   - Returns:
     - `error`: An error if the connection cannot be established.

2. IsConnected() bool
   - Description: Checks if the cache is currently connected.
   - Returns:
     - `bool`: `true` if the cache is connected; otherwise, `false`.

3. DisconnectCache()
   - Description: Closes the connection to the cache gracefully.

4. GetKey(key string) (interface{}, error)
   - Description: Retrieves a value from the cache associated with the provided key.
   - Parameters:
     - `key` (string): The key used to identify the value in the cache.
   - Returns:
     - `interface{}`: The cached value.
     - `error`: An error if the value cannot be retrieved.

5. SetKey(key string, value interface{}) error
   - Description: Stores a value in the cache with the specified key.
   - Parameters:
     - `key` (string): The key to associate with the value in the cache.
     - `value` (interface{}): The value to be stored.
   - Returns:
     - `error`: An error if the value cannot be stored in the cache.

6. FlushAll()
   - Description: Clears all data stored in the cache. Use with caution as it removes all cached items.

7. DeleteKey(key string)
   - Description: Deletes a specific key and its associated value from the cache.
   - Parameters:
     - `key` (string): The key to be deleted from the cache.

8. UpdateKey(key string, value interface{}) error
   - Description: Updates the value associated with a specific key in the cache.
   - Parameters:
     - `key` (string): The key to be updated.
     - `value` (interface{}): The new value to replace the existing one.
   - Returns:
     - `error`: An error if the update cannot be performed.

This interface provides a common set of methods that cache implementations must adhere to, allowing dbFusion to work seamlessly with various caching systems. Implement this interface to create a cache that can be used with dbFusion's caching capabilities.
*/
type Cache interface {
	ConnectCache(connectionUri string, password ...string) error
	IsConnected() bool
	DisconnectCache()
	GetKey(key string) (interface{}, error)
	SetKey(key string, value interface{}) error
	FlushAll()
	DeleteKey()
	UpdateKey()
}
