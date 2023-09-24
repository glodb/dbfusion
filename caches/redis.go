package caches

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"golang.org/x/sync/semaphore"
)

// RedisCache struct implements the Cache interface for Redis.
type RedisCache struct {
	pool      *redis.Pool         // The Redis connection pool.
	semaphore *semaphore.Weighted // Semaphore to control concurrent access to the cache.
}

// ConnectCache is a method of the RedisCache type used to establish a connection to a Redis cache server.
// It creates a connection pool to the Redis server specified by the 'connectionUri'. If a password is provided,
// it authenticates the connection with the password. It also initializes a semaphore for managing concurrent
// cache operations.
//
// Receiver:
//   - rc: A RedisCache instance responsible for connecting to the Redis cache.
//
// Parameters:
//   - connectionUri: The URI or address of the Redis server to connect to.
//   - password: (Optional) A variadic parameter for providing an authentication password for the Redis server.
//
// Returns:
//   - err: An error indicating the success or failure of the cache connection establishment.
func (rc *RedisCache) ConnectCache(connectionUri string, password ...string) error {
	var err error

	// Initialize a semaphore for managing concurrent cache operations.
	rc.semaphore = semaphore.NewWeighted(int64(CACHE_DEFAULT_CONNECTIONS))

	// If the Redis connection pool is not yet created, create it.
	if rc.pool == nil {
		rc.pool, err = rc.newPool(connectionUri, password...)
	}

	// If an error occurs during pool creation, return the error.
	if err != nil {
		return err
	}

	// If there is no error but the pool's connection has an error, attempt to recreate the pool.
	if rc.pool.Get().Err() != nil {
		rc.pool, err = rc.newPool(connectionUri, password...)
	}

	// If an error occurs during pool recreation, return the error.
	if err != nil {
		return err
	}

	// Reinitialize the semaphore to manage concurrent cache operations.
	rc.semaphore = semaphore.NewWeighted(int64(CACHE_DEFAULT_CONNECTIONS))

	// Return nil to indicate successful cache connection establishment.
	return nil
}

// newPool is a private method of the RedisCache type used to create and configure a Redis connection pool.
// It sets the maximum number of idle and active connections in the pool and establishes connections to the
// Redis server specified by the 'uri'. If a password is provided, it authenticates the connection with the password.
//
// Receiver:
//   - rc: A RedisCache instance responsible for managing Redis cache connections.
//
// Parameters:
//   - uri: The URI or address of the Redis server to connect to.
//   - password: (Optional) A variadic parameter for providing an authentication password for the Redis server.
//
// Returns:
//   - pool: A pointer to a configured Redis connection pool.
//   - err: An error indicating the success or failure of the pool creation.
func (rc *RedisCache) newPool(uri string, password ...string) (*redis.Pool, error) {
	// Configure the Redis connection pool with maximum idle and active connections.
	pool := &redis.Pool{
		MaxIdle:   CACHE_MAX_IDLE_CONNECTIONS,
		MaxActive: CACHE_MAX_CONNECTIONS, // Maximum number of connections.
		Dial: func() (redis.Conn, error) {
			var c redis.Conn
			var err error

			// Check if a password is provided, and create a connection accordingly.
			if len(password) != 0 {
				c, err = redis.Dial("tcp", uri, redis.DialPassword(password[0]))
			} else {
				c, err = redis.Dial("tcp", uri)
			}

			// If an error occurs during connection establishment, return the error.
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	// Test the connection by attempting to get a connection from the pool.
	conn := pool.Get()

	// If there is an error with the connection, return an error indicating connection failure.
	if err := conn.Err(); err != nil {
		return nil, err
	}

	// Return the configured Redis connection pool and nil error to indicate success.
	return pool, nil
}

// IsConnected is a method of the RedisCache type used to check whether a connection to the Redis cache server is established.
// It verifies the presence of a connection pool and attempts to acquire a semaphore to ensure exclusive access to this operation.
// It then retrieves a connection from the pool and checks for any errors. If no errors are encountered, it returns true to
// indicate a successful connection; otherwise, it returns false.
//
// Receiver:
//   - rc: A RedisCache instance responsible for managing Redis cache connections.
//
// Returns:
//   - connected: A boolean indicating whether a connection to the Redis cache server is established (true) or not (false).
func (rc *RedisCache) IsConnected() bool {
	// Check if the Redis connection pool is not initialized, indicating no connection.
	if rc.pool == nil {
		return false
	}

	// Acquire a semaphore to ensure exclusive access to this operation.
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)

	// Retrieve a connection from the pool and defer its closure.
	conn := rc.pool.Get()
	defer conn.Close()

	// Check if there are any errors associated with the connection.
	if err := conn.Err(); err != nil {
		return false
	}

	// Return true to indicate that a connection to the Redis cache server is established.
	return true
}

// DisconnectCache is a method of the RedisCache type used to gracefully close the connection to the Redis cache.
// It calls the Close() method on the connection pool, allowing any pending operations to complete before closing
// the connections and releasing associated resources.
//
// Receiver:
//   - rc: A RedisCache instance responsible for managing Redis cache connections.
func (rc *RedisCache) DisconnectCache() {
	// Close the connection pool gracefully.
	rc.pool.Close()
}

// GetKey is a method of the RedisCache type used to retrieve a value associated with a specific key from the Redis cache.
// It acquires a semaphore to ensure exclusive access to this operation, retrieves the value using the "GET" command on
// the Redis connection, and returns the retrieved data along with any potential errors.
//
// Receiver:
//   - rc: A RedisCache instance responsible for managing Redis cache connections.
//
// Parameters:
//   - key: The key for which the associated value is to be retrieved from the Redis cache.
//
// Returns:
//   - data: An interface{} representing the retrieved data associated with the specified key.
//   - err: An error indicating the success or failure of the cache retrieval operation.
func (rc *RedisCache) GetKey(key string) (interface{}, error) {
	// Acquire a semaphore to ensure exclusive access to this operation.
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)

	// Retrieve the value associated with the specified 'key' using the "GET" command on the Redis connection.
	data, err := rc.pool.Get().Do("GET", key)

	// Return the retrieved data and any potential errors.
	return data, err
}

// SetKey is a method of the RedisCache type used to set a key-value pair in the Redis cache.
// It acquires a semaphore to ensure exclusive access to this operation, uses the "SET" command
// on the Redis connection to set the specified 'key' with the provided 'value', and returns any
// potential errors that may occur during the cache update.
//
// Receiver:
//   - rc: A RedisCache instance responsible for managing Redis cache connections.
//
// Parameters:
//   - key: The key under which the 'value' is to be stored in the Redis cache.
//   - value: The value to be associated with the specified 'key' in the Redis cache.
//
// Returns:
//   - err: An error indicating the success or failure of the cache update operation.
func (rc *RedisCache) SetKey(key string, value interface{}) error {
	// Acquire a semaphore to ensure exclusive access to this operation.
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)

	// Use the "SET" command on the Redis connection to set the 'key' with the provided 'value'.
	_, err := rc.pool.Get().Do("SET", key, value)

	// Return any potential errors that may occur during the cache update.
	return err
}

// FlushAll is a method of the RedisCache type used to flush all data from the Redis cache, effectively clearing the cache.
// It acquires a semaphore to ensure exclusive access to this operation, and then uses the "FLUSHALL" command on the Redis
// connection to remove all data stored in the cache.
//
// Receiver:
//   - rc: A RedisCache instance responsible for managing Redis cache connections.
func (rc *RedisCache) FlushAll() {
	// Acquire a semaphore to ensure exclusive access to this operation.
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)

	// Use the "FLUSHALL" command on the Redis connection to remove all data from the cache.
	rc.pool.Get().Do("FLUSHALL")
}

// DelKey is a method of the RedisCache type used to delete a specific key from the Redis cache.
// It acquires a semaphore to ensure exclusive access to this operation, uses the "DEL" command on
// the Redis connection to delete the specified 'key', and returns any potential errors that may occur
// during the key deletion.
//
// Receiver:
//   - rc: A RedisCache instance responsible for managing Redis cache connections.
//
// Parameters:
//   - key: The key to be deleted from the Redis cache.
//
// Returns:
//   - err: An error indicating the success or failure of the key deletion operation.
func (rc *RedisCache) DelKey(key string) error {
	// Acquire a semaphore to ensure exclusive access to this operation.
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)

	// Use the Redis connection pool to send the "DEL" command and delete the specified 'key'.
	_, err := rc.pool.Get().Do("DEL", key)

	// Return any potential errors that may occur during the key deletion.
	return err
}
