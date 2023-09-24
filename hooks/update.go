package hooks

// PreUpdate is an interface that user-defined models can implement to define pre-update hooks.
// Pre-update hooks are executed before updating data in the database, allowing custom
// actions or modifications to be applied to the data before the update operation.
type PreUpdate interface {
	// PreUpdate is a method that should be implemented to perform custom actions or modifications
	// on the data before it is updated in the database.
	//
	// Example Usage:
	//   func (model *MyModel) PreUpdate() PreUpdate {
	//       // Implement logic to customize the data before update.
	//       // This can include validation checks or transformations.
	//       return model // Return the instance of the model with modifications.
	//   }
	//
	//   // When data is updated in the database for MyModel, the PreUpdate method
	//   // will be called automatically to customize the data before the update.
	//   // Any modifications made within PreUpdate will affect the data being updated.
	PreUpdate() PreUpdate
}

// PostUpdate is an interface that user-defined models can implement to define post-update hooks.
// Post-update hooks are executed after data has been successfully updated in the database,
// allowing additional actions or processing to be performed on the updated data.
type PostUpdate interface {
	// PostUpdate is a method that should be implemented to perform custom actions or processing
	// on the data that has been successfully updated in the database.
	//
	// Example Usage:
	//   func (model *MyModel) PostUpdate() PostUpdate {
	//       // Implement logic to process the updated data.
	//       // This can include generating audit logs or notifications.
	//       return model // Return the instance of the model after processing.
	//   }
	//
	//   // After data is updated in the database for MyModel, the PostUpdate method
	//   // will be called automatically with the updated data.
	//   // Any processing done within PostUpdate will affect the updated data.
	PostUpdate() PostUpdate
}
