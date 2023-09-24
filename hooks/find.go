package hooks

// PreFind is an interface that user-defined models can implement to define pre-find hooks.
// Pre-find hooks are executed before a database query is performed to customize or modify
// the query parameters. This can be useful for applying filters or additional conditions
// before retrieving data from the database.
type PreFind interface {
	// PreFind is a method that should be implemented to perform custom actions or modifications
	// on the database query parameters before the query is executed.
	//
	// Example Usage:
	//   func (model *MyModel) PreFind() PreFind {
	//       // Implement logic to customize the query parameters before execution.
	//       // This can include adding filters or conditions to the query.
	//       return model // Return the instance of the model with modifications.
	//   }
	//
	//   // When a database find operation is performed on MyModel, the PreFind method
	//   // will be called automatically to customize the query parameters.
	//   // Any modifications made within PreFind will affect the final query.
	PreFind() PreFind
}

// PostFind is an interface that user-defined models can implement to define post-find hooks.
// Post-find hooks are executed after a database query is performed to process or manipulate
// the query results. This can be useful for performing additional actions on retrieved data
// or applying custom processing to query results.
type PostFind interface {
	// PostFind is a method that should be implemented to perform custom actions or processing
	// on the query results obtained from a database query. It receives the query results
	// as a slice of pointers to the model type.
	//
	// Example Usage:
	//   func (model *MyModel) PostFind() PostFind {
	//       // Implement logic to process the query results.
	//       // This can include data manipulation or any other custom actions.
	//       return model // Return the instance of the model after processing.
	//   }
	//
	//   // After a database find operation is performed on MyModel, the PostFind method
	//   // will be called automatically with the query results.
	//   // Any processing done within PostFind will affect the retrieved data.
	PostFind() PostFind
}
