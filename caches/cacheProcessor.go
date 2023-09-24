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

// ProcessInsertCache is a method of the cacheProcessor type used to insert data into a cache with specified indexes.
// It checks for the number of indexes, acquires a semaphore, and processes the data in parallel.
// It constructs cache indexes based on the given data and indexes, and then processes these indexes.
// If the number of unique keys in an index exceeds a certain limit, it returns an error.
//
// Receiver:
//   - cp: A cacheProcessor instance responsible for processing cache operations.
//
// Parameters:
//   - cache: The Cache interface to interact with the cache system.
//   - indexes: A slice of strings representing the indexes to be created.
//   - data: A map containing data to be cached.
//   - dbName: The name of the database.
//   - entityName: The name of the entity or collection in the database.
//
// Returns:
//   - err: An error indicating the success or failure of the cache insertion operation.
func (cp *cacheProcessor) ProcessInsertCache(cache Cache, indexes []string, data map[string]interface{}, dbName string, entityName string) (err error) {
	// Check if the number of indexes exceeds a limit.
	if len(indexes) > 10 {
		return dbfusionErrors.ErrCacheIndexesIncreased
	}
	// If there are no indexes, return without processing.
	if len(indexes) == 0 {
		return nil
	}

	// Acquire a semaphore to control concurrent cache processing.
	cp.semaphore.Acquire(context.TODO(), 1)
	defer cp.semaphore.Release(1)

	// Create a WaitGroup to synchronize parallel processing.
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		// Create a slice to store cache indexes.
		cacheIndexes := make([]string, 0)

		// Iterate through the provided indexes.
		for _, val := range indexes {
			// Split unique keys within each index.
			uniqueKeys := strings.Split(val, ",")

			// Check if the number of unique keys exceeds a limit.
			if len(uniqueKeys) > 5 {
				err = dbfusionErrors.ErrCacheUniqueKeysIncreased
			}

			// Construct a cache index based on dbName, entityName, and unique key values.
			index := dbName + "_" + entityName + "_"
			for _, key := range uniqueKeys {
				if value, ok := data[key]; ok {
					index += fmt.Sprintf("%v", value)
					index += "_"
				}
			}

			// Remove the trailing underscore and add the index to the slice.
			if len(index) > 1 {
				index = index[:len(index)-1]
				cacheIndexes = append(cacheIndexes, index)
			}
		}

		// Process the constructed cache indexes.
		err = cp.processIndexes(cache, cacheIndexes, data)
	}()

	// Wait for parallel processing to complete.
	wg.Wait()
	return err
}

// processIndexes is a method of the cacheProcessor type used to process and set cache indexes and associated data.
// It generates a unique ULID (Universally Unique Lexicographically Sortable Identifier), associates it with each index,
// and stores the ULID in the cache. It also encodes and stores the provided data in the cache using the generated ULID.
//
// Receiver:
//   - cp: A cacheProcessor instance responsible for processing cache operations.
//
// Parameters:
//   - cache: The Cache interface to interact with the cache system.
//   - cacheIndexes: A slice of strings representing cache indexes.
//   - data: A map containing data to be cached.
//
// Returns:
//   - err: An error indicating the success or failure of the cache processing operation.
func (cp *cacheProcessor) processIndexes(cache Cache, cacheIndexes []string, data map[string]interface{}) error {
	// Generate a unique ULID based on the current timestamp and entropy source.
	ulid := ulid.MustNew(ulid.Timestamp(time.Now()), cp.entropy).String()

	// Iterate through the provided cache indexes.
	for _, index := range cacheIndexes {
		// Set each cache index with the generated ULID.
		err := cache.SetKey(index, ulid)
		if err != nil {
			return err
		}
	}

	// Encode the data and store it in the cache using the generated ULID as the key.
	encodedData, err := codec.GetInstance().Encode(data)
	if err != nil {
		return err
	}
	err = cache.SetKey(ulid, encodedData)
	if err != nil {
		return err
	}

	// Return nil to indicate successful cache processing.
	return nil
}

// ProcessGetCache is a method of the cacheProcessor type used to retrieve cached data from the cache system.
// It looks up the cache using the specified key, retrieves the associated data, decodes it, and populates the
// provided 'data' interface with the decoded data. If the key or data is not found in the cache, it returns false.
//
// Receiver:
//   - cp: A cacheProcessor instance responsible for processing cache operations.
//
// Parameters:
//   - cache: The Cache interface to interact with the cache system.
//   - key: The key used to retrieve data from the cache.
//   - data: An interface where the retrieved and decoded data will be populated.
//
// Returns:
//   - found: A boolean indicating whether data was found in the cache (true) or not (false).
//   - err: An error indicating the success or failure of the cache retrieval operation.
func (cp *cacheProcessor) ProceessGetCache(cache Cache, key string, data interface{}) (bool, error) {
	// Retrieve the first key associated with the specified 'key' from the cache.
	firstKey, err := cache.GetKey(key)
	if err != nil {
		return false, err
	}

	// If the firstKey is not found in the cache, return false to indicate data not found.
	if firstKey == nil {
		return false, nil
	}

	// Retrieve the data associated with the firstKey from the cache.
	redisData, err := cache.GetKey(string(firstKey.([]byte)))

	// If the redisData is not found in the cache, return false to indicate data not found.
	if redisData == nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	// Decode the retrieved data and populate the 'data' interface with the decoded data.
	codec.GetInstance().Decode(redisData.([]byte), data)

	// Return true to indicate that data was successfully found and retrieved from the cache.
	return true, nil
}

// ProcessGetQueryCache is a method of the cacheProcessor type used to retrieve cached query results from the cache system.
// It looks up the cache using the specified key, retrieves the associated data, decodes it, and populates the provided
// 'data' interface with the decoded data. If the key or data is not found in the cache, it returns false.
//
// Receiver:
//   - cp: A cacheProcessor instance responsible for processing cache operations.
//
// Parameters:
//   - cache: The Cache interface to interact with the cache system.
//   - key: The key used to retrieve data from the cache.
//   - data: An interface where the retrieved and decoded data will be populated.
//
// Returns:
//   - found: A boolean indicating whether data was found in the cache (true) or not (false).
//   - err: An error indicating the success or failure of the cache retrieval operation.
func (cp *cacheProcessor) ProceessGetQueryCache(cache Cache, key string, data interface{}) (bool, error) {
	// Retrieve the data associated with the specified 'key' from the cache.
	redisData, err := cache.GetKey(key)

	// If an error occurs during cache retrieval, return the error.
	if err != nil {
		return false, err
	}

	// If the 'redisData' is not found in the cache, return false to indicate data not found.
	if redisData == nil {
		return false, nil
	}

	// Decode the retrieved data and populate the 'data' interface with the decoded data.
	codec.GetInstance().Decode(redisData.([]byte), data)

	// Return true to indicate that data was successfully found and retrieved from the cache.
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
