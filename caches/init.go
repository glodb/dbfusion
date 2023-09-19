/*
Caches Package

The `caches` package defines the use of caching within the dbFusion system. It has been developed to provide extensible and intelligent caching implementations.

**Cache Interface:**
- As caching is a pluggable component within dbFusion, any cache implementation must adhere to the `CacheInterface`.

**Package Purpose:**
- The primary purpose of this package is to provide the core architecture for caching and to offer an example implementation using Redis.

**Cache Processor:**
- This package includes a cache processor that handles the intricacies of caching.
- The cache processor performs the following tasks:
  - It processes keys and constructs composite indexes.
  - It searches for composite indexes and provides updates.
  - Deletion is straightforward, currently deleting the first index, while composite indexes invalidate on expiration or are moved away for LRU (Least Recently Used) eviction.


Developers can use this package as a foundation for creating their own custom cache solutions within the dbFusion framework. It provides the necessary architecture and an example Redis implementation for reference.

Please note that while this package offers a Redis-based cache example, developers are encouraged to tailor their cache implementations to meet specific use cases and requirements.
*/
package caches
