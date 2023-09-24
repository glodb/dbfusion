// Package hooks provides a framework for defining interfaces that user-defined models can implement to enable
// custom behaviors and interactions with the database. These hooks allow users to define actions that should be
// executed before or after database operations, such as inserts, updates, and deletes.
//
// By implementing the hooks interfaces in their models, users can define pre and post hooks that run before and
// after specific database operations. This enables them to perform custom logic, validations, or data
// transformations tailored to their application's needs.
//
// Additionally, the package supports setting cache values based on these hooks, allowing users to efficiently
// store and retrieve data from a cache to reduce database load and improve performance.
//
// Another feature of this package is support for defining MongoDB indexes through the hooks interface. Users can
// specify the indexes that should be created for their models, optimizing query performance for MongoDB
// databases.
//
// Overall, the hooks package empowers developers to extend and customize the behavior of their models in
// database interactions, cache management, and query optimization.
package hooks
