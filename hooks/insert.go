package hooks

// PreInsert is an interface that user-defined models can implement to define pre-insert hooks.
// Pre-insert hooks are executed before inserting data into the database, allowing custom
// actions or modifications to be applied to the data being inserted.
type PreInsert interface {
	// PreInsert is a method that should be implemented to perform custom actions or modifications
	// on the data before it is inserted into the database.
	//
	// Example Usage:
	//   func (model *MyModel) PreInsert() PreInsert {
	//       // Implement logic to customize the data before insertion.
	//       // This can include setting default values or applying transformations.
	//       return model // Return the instance of the model with modifications.
	//   }
	//
	//   // When data is inserted into the database for MyModel, the PreInsert method
	//   // will be called automatically to customize the data before insertion.
	//   // Any modifications made within PreInsert will affect the data being inserted.
	PreInsert() PreInsert
}

// PostInsert is an interface that user-defined models can implement to define post-insert hooks.
// Post-insert hooks are executed after data has been successfully inserted into the database,
// allowing additional actions or processing to be performed on the inserted data.
type PostInsert interface {
	// PostInsert is a method that should be implemented to perform custom actions or processing
	// on the data that has been successfully inserted into the database.
	//
	// Example Usage:
	//   func (model *MyModel) PostInsert() PostInsert {
	//       // Implement logic to process the inserted data.
	//       // This can include generating additional related records or notifications.
	//       return model // Return the instance of the model after processing.
	//   }
	//
	//   // After data is inserted into the database for MyModel, the PostInsert method
	//   // will be called automatically with the inserted data.
	//   // Any processing done within PostInsert will affect the inserted data.
	PostInsert() PostInsert
}
