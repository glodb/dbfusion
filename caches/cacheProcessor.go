package caches

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/glodb/dbfusion/codec"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/oklog/ulid/v2"
	"golang.org/x/sync/semaphore"
)

// cacheProcessor is a singleton structure that handles all cache-related operations for composite and single indexes.
type cacheProcessor struct {
	semaphore *semaphore.Weighted    // Semaphore for controlling parallel cache operations.
	entropy   *ulid.MonotonicEntropy // Entropy source for ULID generation.
}

var (
	instance *cacheProcessor
	once     sync.Once
)

// GetInstance returns a singleton instance of the cache processor.
func GetInstance() *cacheProcessor {
	once.Do(func() {
		instance = &cacheProcessor{}
		instance.semaphore = semaphore.NewWeighted(int64(CACHE_PARALLEL_PROCESS))
		instance.entropy = ulid.Monotonic(rand.Reader, 0)
	})
	return instance
}

// ProcessInsertCache processes the insertion of data into the cache for composite and single indexes.
// It generates composite indexes based on unique keys and stores data with ULID keys in the cache.
// Parameters:
//   - cache (Cache): The cache implementation to use.
//   - indexes ([]string): List of composite indexes to create.
//   - data (map[string]interface{}): The data to be cached.
//   - dbName (string): The name of the database.
//   - entityName (string): The name of the entity.
// Returns:
//   - error: An error if the cache operation encounters issues.
func (cp *cacheProcessor) ProcessInsertCache(cache Cache, indexes []string, data map[string]interface{}, dbName string, entityName string) (err error) {
	if len(indexes) > 10 {
		return dbfusionErrors.ErrCacheIndexesIncreased
	}
	if len(indexes) == 0 {
		return nil
	}
	cp.semaphore.Acquire(context.TODO(), 1)
	defer cp.semaphore.Release(1)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		cacheIndexes := make([]string, 0)

		for _, val := range indexes {
			uniqueKeys := strings.Split(val, ",")

			if len(uniqueKeys) > 5 {
				err = dbfusionErrors.ErrCacheUniqueKeysIncreased
			}

			index := dbName + "_" + entityName + "_" // Add dbName and entityName
			for _, key := range uniqueKeys {

				if value, ok := data[key]; ok {
					index += fmt.Sprintf("%v", value)
					index += "_"
				}
			}

			if len(index) > 1 {
				index = index[:len(index)-1]
				cacheIndexes = append(cacheIndexes, index)
			}
		}
		err = cp.processIndexes(cache, cacheIndexes, data)

	}()

	wg.Wait()
	return err
}

// processIndexes is a helper function that takes preprocessed indexes and creates the second index and stores data in the cache.
// Parameters:
//   - cache (Cache): The cache implementation to use.
//   - cacheIndexes ([]string): List of composite indexes to create.
//   - data (map[string]interface{}): The data to be cached.
// Returns:
//   - error: An error if the cache operation encounters issues.
func (cp *cacheProcessor) processIndexes(cache Cache, cacheIndexes []string, data map[string]interface{}) error {
	ulid := ulid.MustNew(ulid.Timestamp(time.Now()), cp.entropy).String()
	for _, index := range cacheIndexes {
		err := cache.SetKey(index, ulid)
		if err != nil {
			return err
		}
	}
	encodedData, err := codec.GetInstance().Encode(data)
	if err != nil {
		return err
	}
	err = cache.SetKey(ulid, encodedData)
	if err != nil {
		return err
	}

	return nil
}

// ProceessGetCache retrieves data from the composite index in the cache and decodes it into the provided data variable.
// Parameters:
//   - cache (Cache): The cache implementation to use.
//   - key (string): The key to retrieve data from the cache.
//   - data (interface{}): A reference to the variable where the retrieved data will be decoded.
// Returns:
//   - bool: `true` if data is found and retrieved successfully; otherwise, `false`.
//   - error: An error if the cache operation encounters issues.
func (cp *cacheProcessor) ProceessGetCache(cache Cache, key string, data interface{}) (bool, error) {
	firstKey, err := cache.GetKey(key)
	if err != nil {
		return false, err
	}

	if firstKey == nil {
		return false, nil
	}
	redisData, err := cache.GetKey(string(firstKey.([]byte)))
	if redisData == nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	codec.GetInstance().Decode(redisData.([]byte), data)

	return true, nil
}

// ProceessGetQueryCache retrieves data from the cache using a single key and decodes it into the provided data variable.
// Parameters:
//   - cache (Cache): The cache implementation to use.
//   - key (string): The key to retrieve data from the cache.
//   - data (interface{}): A reference to the variable where the retrieved data will be decoded.
// Returns:
//   - bool: `true` if data is found and retrieved successfully; otherwise, `false`.
//   - error: An error if the cache operation encounters issues.
func (cp *cacheProcessor) ProceessGetQueryCache(cache Cache, key string, data interface{}) (bool, error) {
	redisData, err := cache.GetKey(key)

	if err != nil {
		return false, err
	}

	if redisData == nil {
		return false, nil
	}

	codec.GetInstance().Decode(redisData.([]byte), data)
	return true, nil
}

// ProceessSetQueryCache stores data in the cache using a single key.
// Parameters:
//   - cache (Cache): The cache implementation to use.
//   - key (string): The key to store data in the cache.
//   - data (interface{}): The data to be stored.
// Returns:
//   - bool: `true` if data is stored successfully; otherwise, `false`.
//   - error: An error if the cache operation encounters issues.
func (cp *cacheProcessor) ProceessSetQueryCache(cache Cache, key string, data interface{}) (bool, error) {
	encodedData, err := codec.GetInstance().Encode(data)
	if err != nil {
		return false, err
	}
	err = cache.SetKey(key, encodedData)

	if err != nil {
		return false, err
	}
	return true, nil
}
