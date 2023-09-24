package hooks

// NormalIndexes is an interface that user-defined models can implement to specify normal indexes
// for a MongoDB collection. Normal indexes are used for efficient querying of fields.
type NormalIndexes interface {
	// GetNormalIndexes should be implemented to return an array of field names on which normal indexes
	// should be created. These indexes are used for optimizing queries on the specified fields.
	//
	// Example Usage:
	//   func (model *MyModel) GetNormalIndexes() []string {
	//       return []string{"fieldName1:1", "fieldName2:-1"}
	//   }
	//
	// In this example, the model specifies that normal indexes should be created on "fieldName1" and "fieldName2".
	GetNormalIndexes() []string
}

// UniqueIndexes is an interface that user-defined models can implement to specify unique indexes
// for a MongoDB collection. Unique indexes ensure that values in the specified fields are unique across
// all documents in the collection.
type UniqueIndexes interface {
	// GetUniqueIndexes should be implemented to return an array of field names on which unique indexes
	// should be created. These indexes enforce uniqueness constraints on the specified fields.
	//
	// Example Usage:
	//   func (model *MyModel) GetUniqueIndexes() []string {
	//       return []string{"uniqueField1:1", "uniqueField2:-1"}
	//   }
	//
	// In this example, the model specifies that unique indexes should be created on "uniqueField1" and "uniqueField2".
	GetUniqueIndexes() []string
}

// TextIndexes is an interface that user-defined models can implement to specify a text index
// for a MongoDB collection. Text indexes enable full-text search on the specified fields.
type TextIndexes interface {
	// GetTextIndex should be implemented to return the name of the field on which a text index
	// should be created. This index allows full-text search on the specified field.
	//
	// Example Usage:
	//   func (model *MyModel) GetTextIndex() string {
	//       return "textField"
	//   }
	//
	// In this example, the model specifies that a text index should be created on the "textField".
	GetTextIndex() string
}

// TwoDimensionalIndexes is an interface that user-defined models can implement to specify two-dimensional
// indexes for a MongoDB collection. Two-dimensional indexes are used for geospatial data.
type TwoDimensionalIndexes interface {
	// Get2DIndexes should be implemented to return an array of field names on which two-dimensional indexes
	// should be created. These indexes are used for geospatial data.
	//
	// Example Usage:
	//   func (model *MyModel) Get2DIndexes() []string {
	//       return []string{"location"}
	//   }
	//
	// In this example, the model specifies that two-dimensional indexes should be created on the "location" field.
	Get2DIndexes() []string
}

// TwoDimensionalSpatialIndexes is an interface that user-defined models can implement to specify two-dimensional
// spatial indexes for a MongoDB collection. These indexes are used for geospatial data with spatial coordinates.
type TwoDimensionalSpatialIndexes interface {
	// Get2DSpatialIndexes should be implemented to return an array of field names on which two-dimensional
	// spatial indexes should be created. These indexes are used for geospatial data with spatial coordinates.
	//
	// Example Usage:
	//   func (model *MyModel) Get2DSpatialIndexes() []string {
	//       return []string{"location"}
	//   }
	//
	// In this example, the model specifies that two-dimensional spatial indexes should be created on the "location" field.
	Get2DSpatialIndexes() []string
}

// HashedIndexes is an interface that user-defined models can implement to specify hashed indexes
// for a MongoDB collection. Hashed indexes are used for efficiently querying hashed values.
type HashedIndexes interface {
	// GetHashedIndexes should be implemented to return an array of field names on which hashed indexes
	// should be created. These indexes are used for optimizing queries on hashed values.
	//
	// Example Usage:
	//   func (model *MyModel) GetHashedIndexes() []string {
	//       return []string{"hashedField1:1", "hashedField2:-1"}
	//   }
	//
	// In this example, the model specifies that hashed indexes should be created on "hashedField1" and "hashedField2".
	GetHashedIndexes() []string
}

// SparseIndexes is an interface that user-defined models can implement to specify sparse indexes
// for a MongoDB collection. Sparse indexes only index documents that contain the indexed field.
type SparseIndexes interface {
	// GetSparseIndexes should be implemented to return an array of field names on which sparse indexes
	// should be created. These indexes only index documents that contain the indexed field, allowing
	// for efficient queries on specific fields.
	//
	// Example Usage:
	//   func (model *MyModel) GetSparseIndexes() []string {
	//       return []string{"sparseField1:1", "sparseField2:-1"}
	//   }
	//
	// In this example, the model specifies that sparse indexes should be created on "sparseField1" and "sparseField2".
	GetSparseIndexes() []string
}
