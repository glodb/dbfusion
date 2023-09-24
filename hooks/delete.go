package hooks

// PreDelete is an interface that user-defined models can implement to define pre-delete hooks.
// Pre-delete hooks are executed before a data deletion operation, allowing developers to perform
// custom actions or validations before the deletion occurs.
type PreDelete interface {
	// PreDelete is a method that should be implemented to specify the pre-delete logic for a model.
	// This method is called before a data deletion operation is performed.
	//
	// Example Usage:
	//   func (model *MyModel) PreDelete() PreDelete {
	//       // Implement custom pre-delete logic here.
	//       // This method will be executed before data deletion.
	//       return model
	//   }
	//
	//   // Now, whenever a data deletion operation is called on MyModel, the PreDelete method
	//   // will be invoked to execute the defined pre-delete logic.
	PreDelete() PreDelete
}

// PostDelete is an interface that user-defined models can implement to define post-delete hooks.
// Post-delete hooks are executed after a data deletion operation, allowing developers to perform
// additional actions or clean-up tasks following the deletion of data.
type PostDelete interface {
	// PostDelete is a method that should be implemented to specify the post-delete logic for a model.
	// This method is called after a data deletion operation has been successfully executed.
	//
	// Example Usage:
	//   func (model *MyModel) PostDelete() PostDelete {
	//       // Implement custom post-delete logic here.
	//       // This method will be executed after data deletion.
	//       return model
	//   }
	//
	//   // Now, whenever a data deletion operation is completed on MyModel, the PostDelete method
	//   // will be invoked to execute the defined post-delete logic.
	PostDelete() PostDelete
}
