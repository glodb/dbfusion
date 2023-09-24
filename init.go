// Package DBfusion serves as an Object-Relational Mapping (ORM) layer designed
// to work with MySQL and MongoDB databases. It provides a unified interface for
// managing database connections and executing queries, offering intelligent
// caching options to improve performance.
//
// Connection Management:
// To establish a database connection, users should call 'getMongoConnection' or
// 'getSqlConnection' depending on the target database system. These functions
// return a connection object ready for database operations.
//
// Caching Support:
// DBfusion offers caching support, allowing users to pass a cache instance as
// an option. This cache can automatically store query results and intelligently
// decide whether to retrieve data from the cache or query the database.
//
// Query Execution:
// DBfusion intelligently decides whether to query the database or use cached
// results, providing flexibility to users. Users can choose to force a database
// query instead of relying on cached data when necessary.
//
// Purpose:
// DBfusion simplifies database interactions by providing a consistent API for
// both MySQL and MongoDB. It abstracts the complexities of database connections
// and caching, allowing developers to focus on their application's logic.
//
// Usage:
// Developers can import and use this package to obtain database connections,
// execute queries, and leverage caching for performance optimization.
// Note:
// DBfusion simplifies database access and caching, offering developers an
// efficient way to interact with MySQL and MongoDB databases.
package dbfusion
