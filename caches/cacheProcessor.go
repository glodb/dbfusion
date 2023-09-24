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

// ProceessSetQueryCache is a method of the cacheProcessor struct used to set data in a cache with the provided key.
//
// Method Signature:
//   func (cp *cacheProcessor) ProceessSetQueryCache(cache Cache, key string, data interface{}) (bool, error)
//
// Parameters:
//   - cache: A Cache interface representing the cache storage where the data will be set.
//   - key: A string representing the key under which the data will be stored in the cache.
//   - data: An interface{} containing the data to be stored in the cache. It will be encoded before storage.
//
// Returns:
//   - bool: A boolean indicating whether the data was successfully set in the cache (true) or not (false).
//   - error: An error if any error occurs during data encoding or cache setting, or nil if the operation is successful.
//
// Description:
// This method encodes the provided data using a codec and then sets it in the cache with the given key.
// It returns a boolean value indicating the success of the operation (true for success, false for failure) and
// an error if any error occurs during encoding or cache setting.
// If the encoding of data fails, it returns false and the encountered error.
// If the cache setting fails, it returns false and the encountered error.
// If both encoding and setting are successful, it returns true and nil to indicate success.
// This method is designed for storing data in a cache with error handling.
func (cp *cacheProcessor) ProceessSetQueryCache(cache Cache, key string, data interface{}) (bool, error) {
	// Encode the provided data using a codec instance.
	encodedData, err := codec.GetInstance().Encode(data)
	if err != nil {
		// If encoding fails, return false and the encountered error.
		return false, err
	}

	// Set the encoded data in the cache using the provided key.
	err = cache.SetKey(key, encodedData)
	if err != nil {
		// If cache setting fails, return false and the encountered error.
		return false, err
	}

	// If both encoding and setting are successful, return true to indicate success.
	return true, nil
}

// ProceessUpdateCache updates the Redis cache by deleting old keys and creating new keys.
// Parameters:
//   - cache (Cache): The cache instance used for data storage and retrieval.
//   - oldKeys ([]string): A slice of old keys to be deleted from the cache.
//   - newKeys ([]string): A slice of new keys to be created in the cache.
//   - data (interface{}): The data to be stored in the cache.
//
// Returns:
//   - bool: A boolean indicating success (true) or failure (false).
//   - error: An error, if any, that occurred during the cache update.
func (cp *cacheProcessor) ProceessUpdateCache(cache Cache, oldKeys []string, newKeys []string, data interface{}) (bool, error) {
	// Check if there are no new keys to create, and return success.
	if len(newKeys) == 0 {
		return true, nil
	}

	compKey := ""

	// This loop deletes the old keys but tries to save the composite key.
	for _, key := range oldKeys {
		if compKey == "" {
			// Attempt to retrieve the firstKey associated with the old key.
			firstKey, err := cache.GetKey(key)
			if err != nil {
				return false, err
			}
			if firstKey != nil {
				compKey = string(firstKey.([]byte))
			}
		}
		// Delete the old key from the cache.
		cache.DelKey(key)
	}

	// If composite key not found, create one.
	if compKey == "" {
		compKey = ulid.MustNew(ulid.Timestamp(time.Now()), cp.entropy).String()
	}

	// Set the composite keys on the new keys in the cache.
	for _, key := range newKeys {
		cache.SetKey(key, compKey)
	}

	// Encode the data and store it in the cache with the composite key.
	encodedData, err := codec.GetInstance().Encode(data)
	if err != nil {
		return false, err
	}
	cache.SetKey(compKey, encodedData)

	// Return success and no error.
	return false, nil
}

// ProceessDeleteCache is a method of the cacheProcessor struct used to perform a batch deletion of cache entries.
// It takes a Cache interface and a list of oldKeys as input, where oldKeys represent the keys to be deleted from the cache.
// Additionally, this function identifies a composite key associated with the first old key found in the cache and deletes it.
//
// Method Signature:
//   func (cp *cacheProcessor) ProceessDeleteCache(cache Cache, oldKeys []string) error
//
// Parameters:
//   - cp: A pointer to the cacheProcessor struct that implements this method.
//   - cache: A Cache interface representing the cache storage where entries will be deleted.
//   - oldKeys: A slice of strings containing the keys to be deleted from the cache.
//
// Returns:
//   - error: An error if any error occurs during cache entry deletion, or nil if the operation is successful.
//
// Description:
// This method starts by attempting to retrieve the firstKey associated with the first old key found in the list.
// If the firstKey exists, it is treated as a composite key, and its value is stored in the compKey variable.
// Then, it iterates through each old key and deletes the corresponding cache entry from the cache.
// If any error occurs during cache entry deletion, the function returns the error immediately.
//
// After deleting the individual cache entries, it proceeds to delete the composite key (if it was found earlier).
// If any error occurs during the composite key deletion, the function returns the error.
//
// This method is designed for batch cache management and cleanup, including the removal of associated composite keys.
func (cp *cacheProcessor) ProceessDeleteCache(cache Cache, oldKeys []string) error {
	// Initialize a variable to store the composite key associated with the first old key found.
	compKey := ""

	// Loop through each old key in the provided list.
	for _, key := range oldKeys {
		// If compKey is empty, attempt to retrieve the firstKey associated with the old key.
		if compKey == "" {
			// Attempt to retrieve the firstKey associated with the old key.
			firstKey, err := cache.GetKey(key)
			if err != nil {
				// If an error occurs during retrieval, return the error.
				return err
			}
			// If a firstKey exists, set compKey to its value.
			if firstKey != nil {
				compKey = string(firstKey.([]byte))
			}
		}

		// Delete the cache entry associated with the current old key.
		err := cache.DelKey(key)

		// If an error occurs during deletion, return the error.
		if err != nil {
			return err
		}
	}

	// Delete the composite key associated with the first old key found.
	err := cache.DelKey(compKey)

	// If an error occurs during deletion, return the error.
	if err != nil {
		return err
	}

	// Return nil to indicate successful deletion of cache entries.
	return nil
}
