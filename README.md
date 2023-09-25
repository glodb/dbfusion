# DBFusion Documentation

## Overview

DBFusion is a comprehensive GoLang library designed to create a centralized platform that seamlessly integrates with both SQL and NoSQL databases. Its primary feature is the ability to add a cache layer to any database. While the library currently supports Redis as its internal caching mechanism, developers can integrate their preferred cache systems by implementing the Cache Interface within the caches module. DBFusion strives to be a developer-friendly ORM (Object-Relational Mapping) and has been thoroughly tested with MongoDB and MySQL. It is not intended to replace existing SQL or NoSQL drivers but rather to provide an intuitive and user-friendly alternative.

### Features

- **Enhanced Cache Support**: DBFusion offers built-in cache support with cache hooks.
- **Hooks**: Pre and post-operation hooks for customization.
- **Table Name Hooks**: Define or extract table names from the structure.
- **Pagination Support**: Built-in pagination and aggregate pagination for MongoDB.
- **Chaining API**: Fluent API for constructing complex database queries.

## Installation

To incorporate DBFusion into your Go project, use the following import statement:

```go
import "github.com/globdb/dbfusion"
```
## MySql Integration

### Establishing a Connection
To connect to a MySQL database, you need a valid URI and, if caching is required, a working cache connection. Here is how you can establish a connection:

```go
validDBName := "your_db"
validUri := "username:password@tcp(dbhost:port)/dbname"
cache := caches.RedisCache{}
err := cache.ConnectCache("host:port")
if err != nil {
  t.Errorf("Error in Redis connection: %v", err)
}
options :=
  dbfusion.Options{
    DbName: &validDBName,
    Uri:    &validUri,
    Cache:  &cache,
  }
con, err := dbfusion.GetInstance().GetMySqlConnection(options)
```

This code snippet will give you a valid connection to your MySQL database.

### Creating a Table

To create a table in MySQL with indexing, you can define the schema using struct tags. Here's an example schema definition:

```go

type UserCreateTable struct {
	Id        int    `dbfusion:"id,INT,AUTO_INCREMENT,PRIMARY KEY"`
	Email     string `dbfusion:"email,omitempty,VARCHAR(255),NOT NULL,UNIQUE"`
	Phone     string `dbfusion:"phone,VARCHAR(255),NOT NULL"`
	Password  string `dbfusion:"password,VARCHAR(50),NOT NULL"`
	FirstName string `dbfusion:"firstName,VARCHAR(50)"`
	LastName  string `dbfusion:"lastName,VARCHAR(50)"`
	CreatedAt int    `dbfusion:"createdAt,INT"`
	UpdatedAt int    `dbfusion:"updatedAt,INT"`
}
```
After defining the schema, you can create the table using the CreateTable function on the MySQL connection:

```go
CreateTable(dataStructure interface{}, bool IfExists)
```

Simply pass an object of UserCreateTable, and set ifExists to true if you want the query to run only if the table does not exist.

### Chaining Queries
DBFusion supports query chaining for constructing complex database queries. You can use the following functions in a chained manner:

#### Where

The Where function defines the conditions that must be met to retrieve data. It has the following signature:
```go
Where(interface{}) SQLConnection
```
#### Select
Select describes the variables to be selected by the query by default its *.
```go
Select(keys map[string]bool) SQLConnection
```
Select specifies the variables to be selected in the query (default is *). You can use it as follows:
```go
con.Select(map[string]bool{"firstname": true})
```
#### Table
Table sets the name of the table. If the passed variable is not a struct with a name, you can specify the table name using this function:

```go
Table(tableName string) SQLConnection
```
#### Group By
Use GroupBy to specify grouping by a field:

```go
GroupBy(fieldname string) SQLConnection
```
#### Having
Having works similarly to Where but is used for aggregation queries: 
```go
Having(conditions interface{}) SQLConnection
```
#### Sort
Sort is used to specify the sorting order:
```go
Sort(sortKey string, sortdesc ...bool) SQLConnection
```
Example
In order to sort descending on lastname we can use
```go
con.Sort("lastname", true)
```
#### Skip
Skip determines how many rows to skip:

```go
Skip(skip int64) SQLConnection
```
#### Limit 
Limit specifies how many rows to return:

```go
Limit(limit int64) SQLConnection
```
#### Join
Join is used for joining tables. It has the following structure:
```go
Join(join joins.Join) SQLConnection
type Join struct {
	Operator JoinType
	TableName string
	Condition string
}
```

JoinType includes
```go
const (
	INNER_JOIN = JoinType(1)
	LEFT_JOIN = JoinType(2)
	RIGHT_JOIN = JoinType(3)
	CROSS_JOIN = JoinType(4)
)
```
If we want to join user with address we need to use this
```go
con.Table("users").Join(Join{Operator:LEFT_JOIN, TableName:"address", Condition:"address.userId=user.id"})
```

### Inserting Data
To insert data into the database, use the InsertOne function:

```go
InsertOne(interface{}) error
```
The object you pass should have DBFusion struct tags. InsertOne supports all the data types defined and supported by the library. It also creates cache keys defined by cache hooks.

### Retrieving Data with Find
The Find function is essential for retrieving data. It supports caching and options to force fetching data from the database or cache. Here's how it works:
```go
FindOne(interface{}, ...queryoptions.FindOptions) error
ForceDB bool
CacheResult bool
```
You can specify options in FindOptions to skip checking the cache, force results from the database, or cache the results for future queries.

#### Conditions
Conditions are supported by ftypes.DMap. You can use various conditions in your queries. Here are some examples:

```go
// Example 1: Select all records from the "users" table.
user := models.User{}
con.FindOne(&user)

// Example 2: Select records with a WHERE condition.
con.Where(ftypes.DMap{{Key: "firstname =", Value: "Aafaq"}}).FindOne(&user)

// Example 3: Select specific keys.
con.Where(ftypes.DMap{{Key: "firstname =", Value: "Aafaq"}}).Select(map[string]bool{"firstname": true}).FindOne(&user)

// Example 4: Select with DISTINCT.
con.Where(ftypes.DMap{{Key: "Distinct(firstname) =", Value: "Aafaq"}}).Select(map[string]bool{"firstname": true}).FindOne(&user)

// Example 5: Select with multiple conditions.
con.Where(ftypes.DMap{{"("}, {"firstname =", "Aafaq"}, {" AND lastname =", "Zahid"}, {")"}, {"OR email IN "}, []interface{}{"aafaqzahid9@gmail.com", "aafaq.zahid9@gmail.com"}}).
  Select(map[string]bool{"firstname": true}).
  Skip(20).
  Limit(10).
  Join(joins.Join{Operator: joins.INNER_JOIN, TableName: "users b", Condition: "users.firstname=b.firstname"}).
  FindOne(&user)

```
### Update
To update data and find one record, use the UpdateAndFindOne function:
```go
UpdateAndFindOne(interface{}, interface{}, bool) error
```
This function synchronizes the cache after the update.

### Delete
To delete one record from the database, use the DeleteOne function. If you provide a parameter, it checks if caching is implemented and deletes
```go
DeleteOne(...interface{}) error
```
### Pagination

DBFusion supports pagination for retrieving large datasets. By default, the page size is set to 10, but you can customize it using the `SetPageSize` function. Pagination works similarly to the `Find` function but returns results as an array of objects. Here's how you can use it:

```go
users := make([]models.User{}, 0)
pageNumber := 2
con.Where(ftypes.DMap{{Key: "firstname =", Value: "Aafaq"}}).SetPageSize(15).Paginate(&users, 2)
```
In this example, we retrieve a page of records with a page size of 15 and page number 2. The results are stored in the users array.

The SetPageSize function allows you to control the number of records per page, making it easy to manage and display large datasets.

This pagination feature is useful when dealing with extensive data and ensures efficient retrieval and display of records.

## Mongodb

## Cache

## Hooks

## Struct Tags

## Data Types

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