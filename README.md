DBFusion is ambitious project aimed at created centralized platform in golang that integrates with SQL and NoSQL databases.
The main feature of DBFusion is that it adds a cache on any database.
Currently library inculdes only redis cache internally. But developers can add cache of their own choice by implementing Cache Interface inside caches module. As cache is passed from outside this developer independence is supported.
This library aims to be an ORM currently tested databases are mongo and mysql. This library doesn't aims to be competition to any sql or no sql driver as it inspires all the drivers just tries to be developer friendly.

You can import library by running 
go get github.com/glodb/dbfusion

To integrate DBFusion into your Golang project, you can use the following import statement:

```go
import "github.com/globdb/dbfusion"
```


The library supports
Built in enchanced cache support with cachehooks
Pre and Post operations hooks
Table name hooks which can be defined or extracted from structure name
Built in Pagination support
Aggregate Pagination for mongo db
Chaining API

The common functionality for all databases is cache support and hook

1- Cache Support
One of the main reason for building this package was to support cache over any database seamlessly. In this library we tried to achieve that.
Cache can be enable in any of the golang structure by implementing the hook. 

Consider the following structure

```go
type UserTest struct {
	FirstName string `dbfusion:"firstname"`
	Email     string `dbfusion:"email"`
	Username  string `dbfusion:"username"`
	Password  string `dbfusion:"password"`
	CreatedAt int64  `dbfusion:"createdAt"`
	UpdatedAt int64  `dbfusion:"updatedAt"`
}
```

To integrate the cache implement the function GetCacheIndexes() as follows

```go
func (ne UserTest) GetCacheIndexes() []string {
	return []string{"email", "email,password", "email,username"}
}
```

The function returns cache indexes in form of array.
Multiple indexes are separated by comma's. The library will create the cache indexes on the basis of provided indexes.
The things to consider in creating cache indexes is library separates the index on database and tables but it doesn't support uniquenes in the indexes that needs to be handled by implementations.

2- Hooks
Apart from cache hooks library supports a lot of other hooks
 - preInsert
  - postInsert
  - preFind
  - postFind
  - preUpdate
  - postUpdate
  - preDelete
  - postDelete

Note this 
# DBFusion: Centralized Database Support for Golang

# DBFusion

## Objective

DBFusion is an ambitious project aimed at creating a centralized platform that seamlessly integrates both SQL and NoSQL databases. Currently in the development stage, DBFusion's mission is to provide a comprehensive solution for managing diverse database systems with the potential for further expansion.

## Features

- **Full-Featured ORM**: DBFusion boasts a powerful Object-Relational Mapping (ORM) system for MongoDB and MySQL, with plans to extend support to other databases in the future.

- **Flexible Caching System**: A versatile caching system is at the core of DBFusion. It offers two distinct caching options:
  - Cache storage for entire objects, enabling efficient cache updates.
  - Storage of query results based on user-defined requirements.

- **Extensive API**: DBFusion provides an extensive API, giving programmers the freedom to choose between accessing data stored in the cache or querying the database directly.

- **Database Hooks**: DBFusion includes hooks for various database operations, allowing developers to intervene at different stages of the process. The supported hooks are:
  - preInsert
  - postInsert
  - preFind
  - postFind
  - preUpdate
  - postUpdate
  - preDelete
  - postDelete


- **Pagination Support**: DBFusion offers pagination support for both SQL and NoSQL databases, enhancing data retrieval efficiency.

- **Specialized Pagination**: MongoDB users benefit from specialized aggregatePagination for optimized query performance.

- **Driver Compatibility**: DBFusion is designed to work seamlessly with original database drivers' conditions.

- **Integrated Condition Builder**: The framework includes its own condition builder, making the process of creating cache keys straightforward.

## Development Status

DBFusion is currently in the development stage, with ongoing efforts to expand its capabilities and provide a robust solution for developers working with SQL and NoSQL databases.


DBFusion is a powerful and versatile database library for Golang that aims to provide centralized database support for both SQL and NoSQL databases. With its extensive set of features and user-friendly APIs, DBFusion simplifies database management tasks, allowing developers to focus on building robust and scalable applications.

## Purpose

The primary purpose of DBFusion is to offer a unified and centralized approach to database interactions in Golang applications. Whether you're working with SQL databases like MySQL, PostgreSQL, or NoSQL databases like MongoDB, DBFusion provides a consistent interface, reducing the complexity of dealing with different database systems.

## Key Features

- **In-Memory Caching with Flexibility**: DBFusion is designed with robust in-memory caching capabilities that leverage indexes and cache indexes for efficient data retrieval. It supports pluggable cache providers such as Memcache and Redis, allowing you to choose the caching solution that best suits your application's needs.

- **Caching Strategies with Customization**: The library offers a variety of caching strategies, including JSON and binary encoding for memory-efficient storage, MD5 hashing, and gob encoding. You can prioritize cache items and set time expiration, with additional support for Memcache and Redis.

- **Comprehensive ORM Functionality**: DBFusion includes a full-featured Object-Relational Mapping (ORM) system that supports a wide range of database interactions. This includes associations (Has One, Has Many, Belongs To, Many To Many, Polymorphism, Single-table inheritance), lifecycle hooks (Before/After Create/Save/Update/Delete/Find), and more.

- **Efficient Query Building and Execution**: With DBFusion's SQL builder, you can construct complex queries with ease. It supports upserts, locking, optimizer/index/comment hints, named arguments, and SQL expressions for advanced query manipulation.

- **Chainable APIs for Concise Code**: DBFusion provides chainable APIs that streamline the creation of database queries and operations. This results in concise, readable code that is easier to maintain.

- **Sync Database Schema and Query Cache Optimization**: Keep your database schema synchronized with your application's models using the sync database schema support. Additionally, query cache optimization enhances performance by speeding up query responses.

## Installation

To integrate DBFusion into your Golang project, you can use the following import statement:

```go
import "github.com/your-username/dbfusion"
```

For start should support 
Redis
MemCached

MySQL
PostgreSQL
MongoDB
/////
MariaDB
Oracle Database
Microsoft SQL Server
Amazon Relational Database Service (RDS)
IBM Db2
Elasticsearch
Cassandra
Neo4j
OrientDB
SQLite