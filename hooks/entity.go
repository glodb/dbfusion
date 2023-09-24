package hooks

// Entity is an interface that user-defined models can implement to specify their entity name.
// The entity name is used to uniquely identify and associate models with specific database tables,
// collections, or entities. This can be especially useful when performing database operations
// or defining cache keys based on the entity name.
type Entity interface {
	// GetEntityName is a method that should be implemented to return the unique entity name
	// associated with a user-defined model. The entity name should be a string that uniquely
	// identifies the model's association with a database table, collection, or entity.
	//
	// Example Usage:
	//   func (model *MyModel) GetEntityName() string {
	//       // Implement logic to return the entity name for the model.
	//       // This name will be used for database operations or cache key generation.
	//       return "my_table" // Replace with the actual entity name.
	//   }
	//
	//   // Now, whenever database operations or cache-related hooks are performed on MyModel,
	//   // the GetEntityName method will provide the unique entity name for proper association.
	GetEntityName() string
}
