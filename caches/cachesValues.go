package caches

/*
Cache Configuration Variables

This section documents the configuration variables related to caching in the dbFusion system. These variables control various aspects of cache behavior and usage.

1. MAX_CACHE_SIZE (int)
   - Description: Specifies the maximum size of the cache in units appropriate for your implementation.
   - Default Value: 1024
   - Usage: This variable sets the maximum capacity of the cache. It determines how much data the cache can store before older data is evicted to make room for new entries.

2. USE_CACHE (bool)
   - Description: Indicates whether caching is allowed or not.
   - Default Value: true
   - Usage: Set this variable to control whether dbFusion should use caching. If set to true, caching will be enabled, and queries may utilize cached data if available.

3. CACHE_DEFAULT_CONNECTIONS (int)
   - Description: Defines the default number of connections to the cache.
   - Default Value: 1000
   - Usage: This variable sets the default number of concurrent connections to the cache system. Adjust it based on the capacity and requirements of your cache implementation.

4. CACHE_MAX_IDLE_CONNECTIONS (int)
   - Description: Specifies the maximum number of idle (unused) cache connections.
   - Default Value: 80
   - Usage: Control the maximum number of cache connections that can remain idle without being closed. This optimization helps manage resource usage.

5. CACHE_MAX_CONNECTIONS (int)
   - Description: Sets the maximum number of total cache connections.
   - Default Value: 1000
   - Usage: Define the maximum limit for concurrent cache connections. Adjust this value according to the scalability and capacity needs of your application.

6. CACHE_PARALLEL_PROCESS (int)
   - Description: Determines how many parallel processes can connect to the cache at the same time.
   - Default Value: 1000
   - Usage: Set this variable to manage the number of simultaneous connections that can be established to the cache system. It influences the concurrency of cache-related operations.

These configuration variables allow you to fine-tune the behavior of the caching system within dbFusion, ensuring that it aligns with your application's requirements and resource constraints.
*/
var MAX_CACHE_SIZE = 1024
var USE_CACHE = true
var CACHE_DEFAULT_CONNECTIONS = 1000
var CACHE_MAX_IDLE_CONNECTIONS = 80
var CACHE_MAX_CONNECTIONS = 1000
var CACHE_PARALLEL_PROCESS = 1000
