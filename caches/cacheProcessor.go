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

type CacheProcessor struct {
	semaphore *semaphore.Weighted
	entropy   *ulid.MonotonicEntropy
}

var (
	instance *CacheProcessor
	once     sync.Once
)

// GetInstance returns a singleton instance of the Factory.
func GetInstance() *CacheProcessor {
	once.Do(func() {
		instance = &CacheProcessor{}
		instance.semaphore = semaphore.NewWeighted(int64(CACHE_PARALLEL_PROCESS))
		instance.entropy = ulid.Monotonic(rand.Reader, 0)
	})
	return instance
}

func (cp *CacheProcessor) ProcessInsertCache(cache Cache, indexes []string, data map[string]interface{}, dbName string, entityName string) (err error) {
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

			index := dbName + "_" + entityName + "_"
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

func (cp *CacheProcessor) processIndexes(cache Cache, cacheIndexes []string, data map[string]interface{}) error {
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

func (cp *CacheProcessor) ProceessGetCache(cache Cache, key string, data interface{}) (bool, error) {
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

func (cp *CacheProcessor) ProceessGetQueryCache(cache Cache, key string, data interface{}) (bool, error) {
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

func (cp *CacheProcessor) ProceessSetQueryCache(cache Cache, key string, data interface{}) (bool, error) {
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
