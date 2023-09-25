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
con.Table("users").Where(ftypes.DMap{{Key: "firstname =", Value: "Aafaq"}}).SetPageSize(15).Paginate(&users, 2)
```
In this example, we retrieve a page of records with a page size of 15 and page number 2. The results are stored in the users array.

The SetPageSize function allows you to control the number of records per page, making it easy to manage and display large datasets.

This pagination feature is useful when dealing with extensive data and ensures efficient retrieval and display of records.

Limitation for this pagination is that it always need to provide Table function.

Pagination Returns

```go
type PaginationResults struct {
	TotalDocuments int64 // Total number of documents in the query result.
	TotalPages     int64 // Total number of pages based on the pagination criteria.
	CurrentPage    int64 // The current page number being viewed.
	Limit          int64 // The limit of documents displayed per page.
}
```
While other data is in the provided address of the object.

## Mongodb

## MongoDB Functionality in DBFusion

To establish a connection with MongoDB using DBFusion, you can use the following code snippet as a reference:

```go
validDBName := "testDBFusion"
validUri := "mongodb://localhost:27017"
cache := caches.RedisCache{}
err := cache.ConnectCache("localhost:6379")
if err != nil {
  t.Errorf("Error in redis connection, occurred %v", err)
}
options := dbfusion.Options{
  DbName: &validDBName,
  Uri:    &validUri,
  Cache:  &cache,
}
con, err := dbfusion.GetInstance().GetMongoConnection(options)
```
This code establishes a connection to MongoDB with the specified database name and connection URI. It also demonstrates the integration of a Redis cache for improved performance.

These functions allow you to construct complex queries and aggregations in a flexible and expressive manner, making MongoDB integration seamless and powerful within your Go applications.

### MongoDB Conditions in DBFusion

```go
user := models.User{}

// Example 1: Select all records from the "users" collection.
con.FindOne(&user)

// Example 2: Select records with a WHERE condition.
con.Where(ftypes.DMap{{Key: "firstname", Value: "Aafaq"}}).FindOne(&user)

// Example 3: Select specific keys.
con.Where(ftypes.DMap{{Key: "firstname", Value: ftypes.QMap{"$ne":"Aafaq"}}}).
  Select(map[string]bool{"firstname": true}).
  FindOne(&user)

// Example 4: Select with multiple conditions. in this case show firstname but not lastname
con.Where(ftypes.DMap{{"("}, {"firstname =", "Aafaq"}, {" AND lastname =", "Zahid"}, {")"}, {"OR email IN "}, []interface{}{"aafaqzahid9@gmail.com", "aafaq.zahid9@gmail.com"}}).
  Select(map[string]bool{"firstname": true, "lastname":false}).
  Skip(20).
  Limit(10).
  FindOne(&user)
```
These MongoDB conditions provide fine-grained control over your queries, making it easy to retrieve the specific data you need from your MongoDB collections.

### Inserting Data

To insert data into the database, use the `InsertOne` function:

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
con.Table("users").Where(ftypes.DMap{{Key: "firstname =", Value: "Aafaq"}}).SetPageSize(15).Paginate(&users, 2)
```

In this example, we retrieve a page of records with a page size of 15 and page number 2. The results are stored in the users array.

The SetPageSize function allows you to control the number of records per page, making it easy to manage and display large datasets.

This pagination feature is useful when dealing with extensive data and ensures efficient retrieval and display of records.

Limitation for this pagination is that it always need to provide Table function.

Pagination Returns

```go
type PaginationResults struct {
	TotalDocuments int64 // Total number of documents in the query result.
	TotalPages     int64 // Total number of pages based on the pagination criteria.
	CurrentPage    int64 // The current page number being viewed.
	Limit          int64 // The limit of documents displayed per page.
}
```
While other data is in the provided address of the object.
### MongoDB Aggregation in DBFusion

One of the primary functionalities of MongoDB is data retrieval through aggregation pipelines. DBFusion offers extensive support for MongoDB functionality, allowing you to build powerful and flexible queries using an aggregation pipeline. Below are some of the supported MongoDB aggregation functions and their explanations:

- **Match**: Specifies a `$match` stage in the aggregation pipeline for filtering documents based on criteria.
- **Bucket**: Specifies a `$bucket` stage to categorize documents into buckets based on criteria.
- **BucketsAuto**: Specifies a `$bucketAuto` stage to automatically create buckets for documents based on criteria.
- **AddFields**: Specifies a `$addFields` stage to add new fields to documents in the pipeline.
- **GeoNear**: Specifies a `$geoNear` stage for geospatial queries, finding documents near a specified point.
- **Group**: Specifies a `$group` stage to group documents by specified criteria.
- **LimitAggregate**: Specifies a `$limit` stage in the aggregation pipeline to limit the number of documents returned.
- **SkipAggregate**: Specifies a `$skip` stage to skip a specified number of documents in the result set.
- **SortAggregate**: Specifies a `$sort` stage to sort documents in the aggregation pipeline.
- **SortByCount**: Specifies a `$sortByCount` stage in the aggregation pipeline.
- **Project**: Specifies a `$project` stage to specify the projected fields.
- **Unset**: Specifies a `$unset` stage to remove fields from documents.
- **ReplaceWith**: Specifies a `$replaceWith` stage to replace documents with a new one.
- **Merge**: Specifies a `$merge` stage to merge documents.
- **Out**: Specifies a `$out` stage to specify the output location.
- **Facet**: Specifies a `$facet` stage to use aggregation pipelines for multiple categories.
- **CollStats**: Specifies a `$collStats` stage to collect statistics about a collection.
- **IndexStats**: Specifies a `$indexStats` stage to collect statistics about indexes.
- **PlanCacheStats**: Specifies a `$planCacheStats` stage to collect statistics about the query plan cache.
- **Redact**: Specifies a `$redact` stage to control access to data.
- **ReplaceRoot**: Specifies a `$replaceRoot` stage to replace the root document.
- **ReplaceCount**: Specifies a `$replaceCount` stage to replace the count of matching documents.
- **Sample**: Specifies a `$sample` stage to return a random sample of documents.
- **Set**: Specifies a `$set` stage to set values in documents.
- **Unwind**: Specifies a `$unwind` stage to unwind arrays.
- **Lookup**: Specifies a `$lookup` stage to perform a left outer join on documents.
- **GraphLookup**: Specifies a `$graphLookup` stage to perform recursive graph lookup.


### Chaining Aggregation Functions

You can chain these aggregation functions in any order to create powerful aggregation pipelines. Here's an example of how to use these functions in combination:

```go
user := models.User{}
con.Match(ftypes.DMap{{"firstname": "Aafaq"}}).Limit(2).Aggregate(&user)
```
In this example, we use the Match function to filter documents where the "firstname" field is "Aafaq," then we limit the result to 2 documents. Finally, we perform the aggregation and populate the user object with the results.

These aggregation functions provide the flexibility and expressiveness needed to construct complex queries and aggregations in MongoDB, seamlessly integrated with DBFusion.

### Mongodb Aggregate Pagination
Aggregate Pagination works similarly to the Aggregate function, with the added capability of returning paginated results. This feature is particularly useful when dealing with large datasets and allows you to efficiently retrieve and display data.

#### Using Aggregate Pagination

You can use Aggregate Pagination by chaining the aggregation functions as needed and specifying the page number to retrieve the records for that page. Here's an example:

```go
user := models.User{}
con.Match(ftypes.DMap{{"firstname": "Aafaq"}}).Limit(2).Aggregate(&user)
```
In this example, we use the Match function to filter documents where the "firstname" field is "Aafaq" and limit the result to 2 documents.

Aggregate Pagination Results
When using Aggregate Pagination, you can expect to receive the following information:

```go
type PaginationResults struct {
	TotalDocuments int64 // Total number of documents in the query result.
	TotalPages     int64 // Total number of pages based on the pagination criteria.
	CurrentPage    int64 // The current page number being viewed.
	Limit          int64 // The limit of documents displayed per page.
}
```
The other data, such as the query results, will be populated in the provided address of the object.

Aggregate Pagination is a powerful feature that helps you manage and display large datasets effectively, providing control over the number of records per page and improving the user experience.

### Mongodb create indexes

DBFusion for MongoDB supports various types of indexes that can be defined and managed using hooks and the CreateIndexes function. These indexes help optimize query performance and enforce constraints on your data. This documentation provides an overview of supported index interfaces and how to use them.

Supported Index Interfaces
DBFusion supports the following index interfaces that you can implement in your user-defined models to specify which indexes to create:

#### NormalIndexes
Implement the NormalIndexes interface to specify normal indexes that optimize queries on specified fields.

Example Usage:

```go

func (model *MyModel) GetNormalIndexes() []string {
    return []string{"fieldName1:1", "fieldName2:-1"}
}
```
#### UniqueIndexes
Implement the UniqueIndexes interface to specify unique indexes that enforce uniqueness constraints on specified fields.

Example Usage:

```go

func (model *MyModel) GetUniqueIndexes() []string {
    return []string{"uniqueField1:1", "uniqueField2:-1"}
}
```
#### TextIndexes
Implement the TextIndexes interface to specify a text index for full-text search on the specified field.

Example Usage:

```go

func (model *MyModel) GetTextIndex() string {
    return "textField"
}
```
#### TwoDimensionalIndexes
Implement the TwoDimensionalIndexes interface to specify two-dimensional indexes for geospatial data.

Example Usage:

```go

func (model *MyModel) Get2DIndexes() []string {
    return []string{"location"}
}
```
#### TwoDimensionalSpatialIndexes
Implement the TwoDimensionalSpatialIndexes interface to specify two-dimensional spatial indexes for geospatial data with spatial coordinates.

Example Usage:

```go

func (model *MyModel) Get2DSpatialIndexes() []string {
    return []string{"location"}
}
```
#### HashedIndexes
Implement the HashedIndexes interface to specify hashed indexes for efficiently querying hashed values.

Example Usage:

```go

func (model *MyModel) GetHashedIndexes() []string {
    return []string{"hashedField1:1", "hashedField2:-1"}
}
```
#### SparseIndexes
Implement the SparseIndexes interface to specify sparse indexes that only index documents containing the indexed field.

Example Usage:

```go

func (model *MyModel) GetSparseIndexes() []string {
    return []string{"sparseField1:1", "sparseField2:-1"}
}
```
#### Creating Indexes
After implementing the necessary index interfaces in your models, you can call the CreateIndexes function on your DBFusion connection (con). This function is intelligent enough to create indexes only if they haven't been created before.

Example Usage:

```go

err := con.CreateIndexes()
if err != nil {
    // Handle the error
}
```

## Cache Support and Hooks

DBFusion provides seamless cache support for all database operations. One of the primary motivations behind building this package was to enable efficient caching over any database. In this library, we have achieved just that by implementing cache hooks.

### Cache Integration

To enable caching for a Go structure, you need to implement the `GetCacheIndexes()` function. Consider the following example using a `User` structure:

```go
type User struct {
	FirstName string `dbfusion:"firstname"`
	Email     string `dbfusion:"email"`
	Username  string `dbfusion:"username"`
	Password  string `dbfusion:"password"`
	CreatedAt int64  `dbfusion:"createdAt"`
	UpdatedAt int64  `dbfusion:"updatedAt"`
}

func (ne User) GetCacheIndexes() []string {
	return []string{"email", "email,password", "email,username"}
}
```

The GetCacheIndexes() function returns an array of cache indexes. You can specify multiple indexes separated by commas. The library will create cache indexes based on the provided index combinations. It's important to note that the library separates indexes for databases and tables, but it doesn't handle uniqueness in the indexes; this responsibility falls on the implementation.

These cache indexes are created during the insertion of data and updated when data is modified. They are also deleted when you provide a query in the structure object for deletion.

### Cache Usage in Find Operations
Another valuable use of cache is in the Find operation. If you have a query that needs to be executed frequently, you can instruct the library to store the results of the query. Subsequent executions of the same query will then retrieve the data from the cache, improving performance.

Cache support in DBFusion ensures that your database operations are not only efficient but also optimized for speed and responsiveness.

## Hooks Support

In addition to cache hooks, the DBFusion library offers a wide range of hooks to customize and enhance the behavior of database operations. These hooks provide developers with the flexibility to execute code before or after critical database actions.

The supported hooks include:

- **preInsert**: Execute custom logic before inserting data.
- **postInsert**: Execute custom logic after inserting data.
- **preFind**: Execute custom logic before performing a find operation.
- **postFind**: Execute custom logic after a find operation.
- **preUpdate**: Execute custom logic before updating data.
- **postUpdate**: Execute custom logic after updating data.
- **preDelete**: Execute custom logic before deleting data.
- **postDelete**: Execute custom logic after deleting data.

### Example: Pre-Insert Hook

Here's an example of implementing a pre-insert hook:

```go
func (u User) PreInsert() hooks.PreInsert {

	// Sample password hashing to demonstrate the effect of the pre-insert hook
	hasher := md5.New()
	io.WriteString(hasher, u.Password)
	u.Password = fmt.Sprintf("%x", hasher.Sum(nil))
	u.CreatedAt = 0
	return u
}
```
In the code above, the User object implements a pre-insert hook. This hook is used to hash the user's password before inserting it into the database, providing an extra layer of security. Additionally, it sets the CreatedAt field to a specific value during the insertion process.

These hooks offer a powerful way to customize the behavior of DBFusion and integrate your business logic seamlessly with database operations.

## Supported Struct Tags in DBFusion

DBFusion supports a variety of struct tags to customize the behavior of your Go structures when working with databases. These tags are specified within the DBFusion tag and follow the format of `dbfusion:"<tag>..."`. Here are the supported struct tags and their explanations:

- **Field Name**: To specify the name of the field in the database, use the field name after the colon. For example, `dbfusion:"fieldname"`.
  
- **omitempty**: Use the `omitempty` tag to indicate that the field should be omitted in the database if it has a zero or empty value.

- **AUTO_INCREMENT**: If a field should be an auto-incrementing primary key, you can specify it with the `AUTO_INCREMENT` tag.

- **PRIMARY KEY**: Designate a field as the primary key by adding the `PRIMARY KEY` tag.

- **VARCHAR(*)**: To specify the field's data type, such as VARCHAR with a specific length, use the `VARCHAR(*)` tag. Replace `*` with the desired length.

- **NOT NULL**: Use the `NOT NULL` tag to indicate that a field should not contain NULL values.

- **UNIQUE**: The `UNIQUE` tag enforces uniqueness for the field's values in the database.

- **DEFAULT**: Specify a default value for a field in the database using the `DEFAULT` tag.

- **CHECK**: You can add a `CHECK` tag to define custom check constraints for the field.

- **FOREIGN**: The `FOREIGN` tag is used in foreign key relationships to reference another table.

These struct tags allow you to define the database schema and behavior directly within your Go structures, making it convenient to work with databases and tailor your data models to your application's needs.

## Data Types

In DBFusion, we introduce two convenient shorthand types to simplify working with maps and BSON primitive.D objects: `QMap` and `DMap`.

### QMap: String Keys and Interface{} Values

`QMap` is a shorthand for a map with string keys and interface{} (empty interface) values. It provides a flexible way to work with key-value pairs where the key is always a string, and the value can be of any type. For example:

```go
type QMap map[string]interface{}
```
With QMap, you can easily create and manipulate maps with heterogeneous values, which is particularly useful for dynamic data structures.

### DMap: Ordered Map for Queries
DMap is a shorthand for a BSON primitive.D object, primarily used for ordered maps across all queries in DBFusion. 

DMap ensures that the order of fields in your document is preserved, which can be important when constructing complex queries.

These types help streamline your database-related code by providing clear and concise ways to work with maps and objects, simplifying the process of constructing and managing data for your queries.


## Contributing
We welcome contributions to improve the DBFusion library. To contribute, please follow these steps:

Test the library and post issues.
Fork the repository on GitHub.
Clone your forked repository to your local machine.
Make your changes and test them thoroughly.
Commit your changes with clear commit messages.
Push your changes to your forked repository.
Create a pull request on the main repository.
We appreciate your contributions!

## Contact
If you have questions, suggestions, or need assistance with DBFusion, you can contact us:

Email: aafaqzahid9@gmail.com
GitHub Issues: https://github.com/glodb/dbfusion/issues
Feel free to reach out to us with any inquiries or feedback. We're here to help!


