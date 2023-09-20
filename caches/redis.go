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

// ConnectCache establishes a connection to the Redis cache using the provided URI and optional password.
// It initializes the Redis connection pool if it doesn't exist.
// Parameters:
//   - connectionUri (string): The URI used to connect to the Redis cache.
//   - password (string, optional): An optional password for cache authentication.
// Returns:
//   - error: An error if the connection cannot be established.
func (rc *RedisCache) ConnectCache(connectionUri string, password ...string) error {
	var err error

	rc.semaphore = semaphore.NewWeighted(int64(CACHE_DEFAULT_CONNECTIONS))

	if rc.pool == nil {
		rc.pool, err = rc.newPool(connectionUri, password...)
	}

	if err != nil {
		return err
	}
	if rc.pool.Get().Err() != nil {
		rc.pool, err = rc.newPool(connectionUri, password...)
	}

	if err != nil {
		return err
	}
	rc.semaphore = semaphore.NewWeighted(int64(CACHE_DEFAULT_CONNECTIONS))
	return nil
}

// newPool is a helper function for ConnectCache that creates a new Redis connection pool.
// It configures the pool with the specified maximum idle and maximum active connections, and
// provides a dial function for establishing Redis connections.
//
// Parameters:
//   - uri (string): The URI used to connect to the Redis cache.
//   - password (string, optional): An optional password for cache authentication.
//
// Returns:
//   - *redis.Pool: A new Redis connection pool.
//   - error: An error if the pool cannot be created or a connection cannot be established.
func (rc *RedisCache) newPool(uri string, password ...string) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:   CACHE_MAX_IDLE_CONNECTIONS,
		MaxActive: CACHE_MAX_CONNECTIONS, // max number of connections
		Dial: func() (redis.Conn, error) {
			var c redis.Conn
			var err error
			if len(password) != 0 {
				c, err = redis.Dial("tcp", uri, redis.DialPassword(password[0]))
			} else {
				c, err = redis.Dial("tcp", uri)
			}
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	// Test the connection here if connected
	conn := pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}

	return pool, nil
}

// IsConnected checks if the Redis cache is currently connected.
// Returns:
//   - bool: `true` if the cache is connected; otherwise, `false`.
func (rc *RedisCache) IsConnected() bool {
	if rc.pool == nil {
		return false
	}

	// Acquire a semaphore to ensure exclusive access to this operation.
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)
	conn := rc.pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return false
	}

	return true
}

// DisconnectCache closes the connection to the Redis cache gracefully.
func (rc *RedisCache) DisconnectCache() {
	rc.pool.Close()
}

// GetKey retrieves a value from the Redis cache associated with the provided key.
// Parameters:
//   - key (string): The key used to identify the value in the cache.
// Returns:
//   - interface{}: The cached value.
//   - error: An error if the value cannot be retrieved.
func (rc *RedisCache) GetKey(key string) (interface{}, error) {
	// Acquire a semaphore to ensure exclusive access to this operation.
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)
	data, err := rc.pool.Get().Do("GET", key)
	return data, err
}

// SetKey stores a value in the Redis cache with the specified key.
// Parameters:
//   - key (string): The key to associate with the value in the cache.
//   - value (interface{}): The value to be stored.
// Returns:
//   - error: An error if the value cannot be stored in the cache.
func (rc *RedisCache) SetKey(key string, value interface{}) error {
	// Acquire a semaphore to ensure exclusive access to this operation.
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)
	_, err := rc.pool.Get().Do("SET", key, value)
	return err
}

// FlushAll clears all data stored in the Redis cache. Use with caution as it removes all cached items.
func (rc *RedisCache) FlushAll() {
	// Acquire a semaphore to ensure exclusive access to this operation.
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)
	rc.pool.Get().Do("FLUSHALL")
}

// DelKey deletes a key from the Redis cache.
// It acquires a semaphore to ensure thread safety during the operation.
//
// Parameters:
//   - key (string): The key to be deleted from the cache.
//
// Returns:
//   - error: An error, if any, that occurred during the deletion.
func (rc *RedisCache) DelKey(key string) error {
	// Acquire a semaphore to ensure exclusive access to this operation.
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)

	// Use the Redis connection pool to send the DEL command and delete the key.
	_, err := rc.pool.Get().Do("DEL", key)

	return err
}
