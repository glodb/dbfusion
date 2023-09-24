package implementations

import "reflect"

// entityData represents information about an entity, including its name, data type, and value.
type entityData struct {
	entityName string        // The name of the entity.
	structType int           // The type of structure (e.g., struct, map, etc.).
	dataType   reflect.Type  // The data type of the entity.
	dataValue  reflect.Value // The value of the entity.
}

// preCreateReturn contains information returned from a pre-create hook, including the entity name,
// data in map[string]interface{} format, the original data, keys, placeholders, and values for a database operation.
type preCreateReturn struct {
	entityName   string                 // The name of the entity.
	mData        map[string]interface{} // Data in map[string]interface{} format.
	Data         interface{}            // The original data.
	keys         string                 // Keys for database operation.
	placeholders string                 // Placeholders for database operation.
	values       []interface{}          // Values for database operation.
}

// preFindReturn contains information returned from a pre-find hook, including the entity name,
// query, whereQuery, queryDatabase flag, data type, and data value.
type preFindReturn struct {
	entityName    string        // The name of the entity.
	query         interface{}   // The query for database operation.
	whereQuery    interface{}   // The WHERE query for database operation.
	queryDatabase bool          // Flag indicating whether to query the database.
	dataType      reflect.Type  // The data type of the entity.
	dataValue     reflect.Value // The data value of the entity.
}

// preUpdateReturn contains information returned from a pre-update hook, including the entity name,
// queryData, data type, data value, and struct type.
type preUpdateReturn struct {
	entityName string        // The name of the entity.
	queryData  interface{}   // The query data for database operation.
	dataType   reflect.Type  // The data type of the entity.
	dataValue  reflect.Value // The data value of the entity.
	structType int           // The type of structure (e.g., struct, map, etc.).
}

// preDeleteReturn contains information returned from a pre-delete hook, including the entity name,
// data type, and data value.
type preDeleteReturn struct {
	entityName string        // The name of the entity.
	dataType   reflect.Type  // The data type of the entity.
	dataValue  reflect.Value // The data value of the entity.
}
