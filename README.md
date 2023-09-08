# DBFusion: Centralized Database Support for Golang


DBFusion is a powerful and versatile database library for Golang that aims to provide centralized database support for both SQL and NoSQL databases. With its extensive set of features and user-friendly APIs, DBFusion simplifies database management tasks, allowing developers to focus on building robust and scalable applications.

## Purpose

The primary purpose of DBFusion is to offer a unified and centralized approach to database interactions in Golang applications. Whether you're working with SQL databases like MySQL, PostgreSQL, or NoSQL databases like MongoDB, DBFusion provides a consistent interface, reducing the complexity of dealing with different database systems.

## Key Features

- **In-Memory Caching with Flexibility**: DBFusion is designed with robust in-memory caching capabilities that leverage indexes and cache indexes for efficient data retrieval. It supports pluggable cache providers such as Memcache and Redis, allowing you to choose the caching solution that best suits your application's needs.

- **Caching Strategies with Customization**: The library offers a variety of caching strategies, including JSON and binary encoding for memory-efficient storage, MD5 hashing, and gob encoding. You can prioritize cache items and set time expiration, with additional support for Memcache and Redis.

- **Comprehensive ORM Functionality**: DBFusion includes a full-featured Object-Relational Mapping (ORM) system that supports a wide range of database interactions. This includes associations (Has One, Has Many, Belongs To, Many To Many, Polymorphism, Single-table inheritance), lifecycle hooks (Before/After Create/Save/Update/Delete/Find), and more.

- **Efficient Query Building and Execution**: With DBFusion's SQL builder, you can construct complex queries with ease. It supports upserts, locking, optimizer/index/comment hints, named arguments, and SQL expressions for advanced query manipulation.

- **Auto Migrations and Schema Management**: Automate database schema migrations as your application evolves. DBFusion simplifies the process of managing schema changes, reducing the risk of errors and ensuring consistency.

- **Chainable APIs for Concise Code**: DBFusion provides chainable APIs that streamline the creation of database queries and operations. This results in concise, readable code that is easier to maintain.

- **Transaction Support for Data Integrity**: The library offers comprehensive transaction management, including transactions, nested transactions, save points, and rollback points. This ensures data integrity and consistency, even in complex database interactions.

- **Sync Database Schema and Query Cache Optimization**: Keep your database schema synchronized with your application's models using the sync database schema support. Additionally, query cache optimization enhances performance by speeding up query responses.

## Installation

To integrate DBFusion into your Golang project, you can use the following import statement:

```go
import "github.com/your-username/dbfusion"


For start should support 
LRU Cache
Redis
MemCached
LevelDB
MemoryStore

Oracle Database
MySQL
PostgreSQL
Microsoft SQL Server
MongoDB
Amazon Relational Database Service (RDS)
MariaDB
IBM Db2
Redis
Elasticsearch
Cassandra
Neo4j
OrientDB
SQLite