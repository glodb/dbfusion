package caches

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"golang.org/x/sync/semaphore"
)

type RedisCache struct {
	pool      *redis.Pool
	semaphore *semaphore.Weighted
}

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

	// Test the connection here if needed
	conn := pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}

	return pool, nil
}

func (rc *RedisCache) IsConnected() bool {
	if rc.pool == nil {
		return false
	}

	conn := rc.pool.Get()
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return false
	}

	return true
}

func (rc *RedisCache) DisconnectCache() {
	rc.pool.Close()
}

func (rc *RedisCache) GetKey(key string) (interface{}, error) {
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)
	data, err := rc.pool.Get().Do("GET", key)
	return data, err
}

func (rc *RedisCache) SetKey(key string, value interface{}) error {
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)
	_, err := rc.pool.Get().Do("SET", key, value)
	return err
}

func (rc *RedisCache) FlushAll() {
	rc.semaphore.Acquire(context.TODO(), 1)
	defer rc.semaphore.Release(1)
	rc.pool.Get().Do("FLUSHALL")
}

func (rc *RedisCache) UpdateKey() {

}

func (rc *RedisCache) DeleteKey() {

}
