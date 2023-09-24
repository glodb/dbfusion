package implementations

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/hooks"
	"github.com/glodb/dbfusion/queryoptions"
	"github.com/glodb/dbfusion/set"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DBCommon is a struct that represents common configuration and state
// for database operations. It is used as a building block for constructing
// database queries and managing query parameters.
type DBCommon struct {
	cache        *caches.Cache // A cache instance for caching query results.
	currentDB    string        // The name of the current database.
	tableName    string        // The name of the database table being queried.
	whereQuery   interface{}   // The query conditions for filtering results.
	skip         int64         // The number of records to skip in query results.
	limit        int64         // The maximum number of records to return in query results.
	projection   interface{}   // The fields to project in query results.
	sort         interface{}   // The sorting criteria for query results.
	joins        string        // A string representing joins in the query.
	groupBy      string        // The field for grouping query results.
	havingString string        // The HAVING clause for grouped query results.
	havingValues []interface{} // Values for the parameters in the HAVING clause.
	orderBy      string        // The ORDER BY clause for sorting query results.
	pageSize     int           // The number of records per page for paginated queries.
}

// SetCache associates a cache object with the DBCommon instance, enabling caching
// of database query results for improved performance.
//
// This method takes a pointer to a caches.Cache instance as its parameter, allowing
// you to provide a specific cache implementation. Once set, the provided cache object
// will be used by DBCommon to store and retrieve query results, reducing the need for
// repeated database queries when the same data is requested.
//
// Parameters:
// - cache: A pointer to a caches.Cache instance representing the cache object to be set.
//
// Example:
//   cache := caches.NewRedisCache()
//   dbCommon := &DBCommon{}
//   dbCommon.SetCache(cache)
//   // Now, the dbCommon instance is configured to use the specified cache for caching results.
//
// Note:
//   It's essential to initialize and configure the cache object before calling SetCache.
//   The behavior of caching may vary depending on the cache implementation used.
func (dbc *DBCommon) SetCache(cache *caches.Cache) {
	dbc.cache = cache
}

// ChangeDatabase switches the active database to the specified database name.
//
// This method allows you to change the active database connection to the one associated
// with the provided database name. Once called, any subsequent database operations
// will be executed in the context of the newly selected database.
//
// Parameters:
// - dbName: A string representing the name of the target database to switch to.
//
// Returns:
// - error: An error interface that is always nil in this implementation, as changing
//   the database name is considered successful. In case of more complex database
//   systems, this method might return an error if the database name cannot be changed.
//
// Example:
//   dbCommon := &DBCommon{}
//   err := dbCommon.ChangeDatabase("my_database")
//   if err != nil {
//       fmt.Printf("Error changing database: %v", err)
//   } else {
//       fmt.Println("Database changed successfully.")
//   }
//
// Note:
//   This method assumes that the provided database name exists and is accessible
//   with the current database connection. The behavior may vary depending on the
//   database system being used.
func (dbc *DBCommon) ChangeDatabase(dbName string) error {
	dbc.currentDB = dbName
	return nil
}

// setTable sets the name of the database table for subsequent database operations.
//
// This method allows you to specify the name of the database table on which
// subsequent database operations will be performed. It sets the internal
// 'tableName' field of the DBCommon instance to the provided 'tableName' string.
//
// Parameters:
// - tableName: A string representing the name of the database table to set.
//
// Example:
//   dbCommon := &DBCommon{}
//   dbCommon.setTable("users")
//
// Note:
//   It's important to set the table name before executing any database queries or
//   operations to ensure that the correct database table is targeted.
func (dbc *DBCommon) setTable(tableName string) {
	dbc.tableName = tableName
}

// isFieldSet checks if the provided reflect.Value is set or contains a non-zero value.
//
// This method is used to determine whether a field (represented as a reflect.Value)
// is set or contains a non-zero value. It performs the check by validating whether
// the provided reflect.Value is valid and whether its interface value is equal to
// the zero value of its type.
//
// Parameters:
// - val: A reflect.Value representing the field to be checked.
//
// Returns:
// - bool: A boolean indicating whether the field is set or contains a non-zero value.
//
// Example:
//   var str string = "Hello, World!"
//   val := reflect.ValueOf(str)
//   isSet := dbc.isFieldSet(val) // Returns true since 'str' is set and contains a non-zero value.
//
// Note:
//   This method is commonly used for checking if a field in a struct is assigned a value
//   or left as its zero value, which can be useful for conditional logic in database
//   operations or other scenarios.
func (dbc *DBCommon) isFieldSet(val reflect.Value) bool {
	if val.IsValid() {
		return !reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
	}
	return false
}

// checkPtr examines the provided data and checks if it's a pointer. If it is a pointer,
// it returns the value that the pointer points to and its type; otherwise, it returns the
// original data value and its type.
//
// This method is used to handle cases where data may be a pointer to a value or the value
// itself. It ensures that the correct reflect.Value and reflect.Type are obtained for
// further processing.
//
// Parameters:
// - data: An interface{} representing the data to be examined.
//
// Returns:
// - reflect.Value: The reflect.Value of the data, either the original value or the pointed-to
//   value if data is a pointer.
// - reflect.Type: The reflect.Type of the data, representing its data type.
//
// Example:
//   var strValue string = "Hello, World!"
//   strPtr := &strValue
//   val, typ := dbc.checkPtr(strValue) // Returns val as reflect.Value of strValue and typ as reflect.Type of string.
//
// Note:
//   This method is useful for handling scenarios where data may be passed as a pointer
//   or a direct value, ensuring consistency in further reflection operations.
func (dbc *DBCommon) checkPtr(data interface{}) (reflect.Value, reflect.Type) {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()

	// Check if the provided data is a pointer. If so, obtain the value it points to.
	switch dataType.Kind() {
	case reflect.Ptr:
		ptrValue := reflect.ValueOf(data)
		dataValue = ptrValue.Elem()
		dataType = dataValue.Type()
	}

	// Return the obtained value and its type.
	return dataValue, dataType
}

// merge merges two data structures of the same type, a and b, by combining their non-zero values.
// It uses reflection to create a new instance of the same type as a and iterates over the fields of a and b.
// For each field, it selects the non-zero value from either a or b and populates the corresponding field
// in the result with that value.
//
// Parameters:
// - a: The first data structure to be merged.
// - b: The second data structure to be merged.
//
// Returns:
// - interface{}: A new instance with merged values from a and b.
//
// Example:
//   userA := User{Name: "Alice", Age: 30}
//   userB := User{Name: "Bob", Age: 0}
//   mergedUser := ms.merge(userA, userB)
//
// Note:
//   This method is useful for combining data from two instances of the same type while prioritizing
//   non-zero values. It ensures that the resulting instance contains the most relevant data.
func (ms *DBCommon) merge(a interface{}, b interface{}) interface{} {
	// Use reflection to create a new instance of the same type as a
	result := reflect.New(reflect.TypeOf(a)).Interface()

	// Use reflection to get the values of a and b
	valA := reflect.ValueOf(a)
	valB := reflect.ValueOf(b)

	// If a is a pointer, get the value it points to
	if valA.Type().Kind() == reflect.Ptr {
		ptrValue := reflect.ValueOf(a)
		valA = ptrValue.Elem()
	}

	// If b is a pointer, get the value it points to
	if valB.Type().Kind() == reflect.Ptr {
		ptrValue := reflect.ValueOf(b)
		valB = ptrValue.Elem()
	}

	// Iterate over the fields of a and merge values
	for i := 0; i < valA.Type().NumField(); i++ {
		fieldName := valA.Type().Field(i).Name
		fieldA := valA.FieldByName(fieldName)
		fieldB := valB.FieldByName(fieldName)

		if fieldA.IsValid() && !fieldA.IsZero() {
			// If fieldA is valid and not zero, use its value
			reflect.ValueOf(result).Elem().FieldByName(fieldName).Set(fieldA)
		} else if fieldB.IsValid() && !fieldB.IsZero() {
			// If fieldB is valid and not zero, use its value
			reflect.ValueOf(result).Elem().FieldByName(fieldName).Set(fieldB)
		}
	}

	// Return the merged result
	return result
}

// createTagValueMap creates a map containing field tags and their corresponding values from the given data structure.
//
// Parameters:
// - data: The data structure from which field tags and values are extracted.
//
// Returns:
// - tagMapValue: A map where keys are field tags and values are the corresponding field values.
// - err: An error, if any, that occurred during the process.
//
// Example:
//   type User struct {
//       Name     string `dbfusion:"user_name"`
//       Age      int    `dbfusion:"user_age"`
//       IsActive bool   `dbfusion:"active"`
//   }
//
//   user := User{Name: "Alice", Age: 30, IsActive: true}
//   tagMap, err := dbc.createTagValueMap(user)
//
//   // tagMapValue would be {"user_name": "Alice", "user_age": 30, "active": true}
//
// Note:
//   This method is useful for extracting field tags and their corresponding values from a data structure
//   that uses struct tags for mapping to database columns. It returns a map of tag-value pairs.
func (dbc *DBCommon) createTagValueMap(data interface{}) (tagMapValue map[string]interface{}, err error) {
	// Get the reflect.Value and reflect.Type of the data
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()

	// Initialize the tag-value map
	tagMapValue = make(map[string]interface{})

	// If data is a pointer, get the value it points to
	if dataType.Kind() == reflect.Ptr {
		ptrValue := reflect.ValueOf(data)
		dataValue = ptrValue.Elem()
		dataType = dataValue.Type()
	}

	// Iterate over the fields of the data structure
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)

		// Split the raw struct tags into individual tags
		rawTags := strings.Split(field.Tag.Get("dbfusion"), ",")
		tagName := rawTags[0]

		// Skip fields with empty tags
		if tagName == "" {
			continue
		}

		// Get the field's value as an interface
		value := dataValue.Field(i).Interface()

		// Add the tag and its corresponding value to the map
		tagMapValue[tagName] = value
	}

	// Return the tag-value map and any potential error
	return tagMapValue, err
}

// getEntityName determines the entity name, struct type, data type, and data value of the given data.
// It inspects the provided data and determines whether it represents an entity or a basic data type.
// This information is used for database operations.
//
// Parameters:
// - data: The data to inspect and determine its entity-related properties.
//
// Returns:
// - entityData: An entityData structure containing the following information:
//     - entityName: The name of the entity or data type.
//     - structType: An integer representing the struct type (1 for struct, 2 for slice, 3 for string, 0 for others).
//     - dataType: The reflect.Type of the data.
//     - dataValue: The reflect.Value of the data.
// - err: An error, if any, that occurred during the determination process.
//
// Example:
//   user := User{Name: "Alice", Age: 30}
//   entityData, err := dbc.getEntityName(user)
//
//   // entityData would be {entityName: "User", structType: 1, dataType: reflect.TypeOf(user), dataValue: reflect.ValueOf(user)}
//
//   str := "Hello, World!"
//   entityData, err := dbc.getEntityName(str)
//
//   // entityData would be {entityName: "", structType: 3, dataType: reflect.TypeOf(str), dataValue: reflect.ValueOf(str)}
//
// Note:
//   This method helps identify the entity name and its related properties for database operations.
//   It checks whether the data is a struct, slice, map, string, or other data type and extracts relevant information.
func (dbc *DBCommon) getEntityName(data interface{}) (entityData entityData, err error) {
	// Check if data is a pointer and get its dereferenced value and type
	dataValue, dataType := dbc.checkPtr(data)

	// Initialize entityData structure with dataValue and dataType
	entityData.dataValue = dataValue
	entityData.dataType = dataType

	// Initialize structType to 0 (other) and entitySet to false
	structType := 0
	entitySet := false

	// Determine the data type and set structType accordingly
	switch dataType.Kind() {
	case reflect.Struct:
		structType = 1
	case reflect.Slice:
		// For slices, set the entity name to the current table name and mark entitySet as true
		entityData.entityName = dbc.tableName
		entitySet = true
		structType = 2
	case reflect.Map:
		// Additional check to ensure it's a map[string]interface{}
		if dataType.Key().Kind() == reflect.String && dataType.Elem().Kind() == reflect.Interface {
			if dbc.tableName == "" {
				err = dbfusionErrors.ErrEntityNameRequired
				return
			}
			// For map[string]interface{}, set the entity name to the current table name and mark entitySet as true
			entityData.entityName = dbc.tableName
			entitySet = true
			structType = 2
		} else {
			// If it's not a map[string]interface{}, return an error
			err = dbfusionErrors.ErrStringMapRequired
			return
		}
	case reflect.String:
		// For strings, set the entity name to the current table name and mark entitySet as true
		entityData.entityName = dbc.tableName
		entitySet = true
		structType = 3
	default:
		// For other data types, return an error
		err = dbfusionErrors.ErrStringMapRequired
		return
	}

	// If entitySet is still false, check if the data implements the Entity interface or use its data type name
	if !entitySet {
		if value, ok := interface{}(data).(hooks.Entity); ok {
			entityData.entityName = value.GetEntityName()
		} else {
			entityData.entityName = dataType.Name()
		}
	}

	// Set the structType in entityData and return the result
	entityData.structType = structType
	return
}

// preInsert prepares data for insertion into the database and returns pre-insertion details.
//
// This function inspects the provided data, extracts relevant information, and prepares it for insertion
// into the database. It handles struct, slice, map, and custom PreInsert implementations.
//
// Parameters:
// - data: The data to be inserted into the database.
//
// Returns:
// - preCreateData: A preCreateReturn structure containing pre-insertion details, such as entity name,
//   data map, keys, placeholders, values, and the original data.
// - err: An error, if any, that occurred during the preparation process.
//
// Example:
//   user := User{Name: "Alice", Age: 30}
//   preCreateData, err := dbc.preInsert(user)
//
//   // preCreateData would contain details for inserting the "User" entity into the database.
//
//   data := map[string]interface{}{"name": "Bob", "age": 25}
//   preCreateData, err := dbc.preInsert(data)
//
//   // preCreateData would contain details for inserting the map data into the database.
//
//   customData := CustomData{...}
//   preCreateData, err := dbc.preInsert(customData)
//
//   // preCreateData would contain details for inserting customData into the database.
//
// Note:
//   This function prepares data for insertion into the database by inspecting its type, extracting field values,
//   and handling custom PreInsert implementations if available. It returns details that are used in the
//   insertion process.
func (dbc *DBCommon) preInsert(data interface{}) (preCreateData preCreateReturn, err error) {
	// Get entity name and related information using getEntityName function
	nameData, nameErr := dbc.getEntityName(data)

	if nameErr != nil {
		err = nameErr
		return
	}

	// Extract dataValue, dataType, and structType from nameData
	dataValue := nameData.dataValue
	dataType := nameData.dataType
	structType := nameData.structType

	// Set the entity name in preCreateData
	preCreateData.entityName = nameData.entityName

	// Check if the data implements PreInsert interface and apply PreInsert logic
	if value, ok := interface{}(data).(hooks.PreInsert); ok {
		data = value.PreInsert()
		dataValue = reflect.ValueOf(data)
		dataType = dataValue.Type()
	}

	// Initialize variables for keys, placeholders, and values
	keys := ""
	placeholders := ""
	values := make([]interface{}, 0)

	// Handle struct data type (structType == 1)
	if structType == 1 {
		mData := make(map[string]interface{})
		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)

			rawtags := strings.Split(field.Tag.Get("dbfusion"), ",")
			tagName := rawtags[0]

			tags := set.ConvertArray[string](rawtags)

			if tagName == "" {
				continue
			}

			value := dataValue.Field(i).Interface()

			if tags.Contains("omitempty") {
				if !dbc.isFieldSet(dataValue.Field(i)) {
					continue
				}
			}

			mData[tagName] = value
			keys += tagName + ","
			placeholders += "?,"
			values = append(values, value)
		}
		preCreateData.mData = mData

		// Trim the trailing comma from keys and placeholders
		if len(keys) > 0 {
			preCreateData.keys = keys[:len(keys)-1]
			preCreateData.placeholders = placeholders[:len(placeholders)-1]
			preCreateData.values = values
		}
	} else {
		// Handle other data types (e.g., map[string]interface{})
		preCreateData.mData = data.(map[string]interface{})
	}

	// Set the original data in preCreateData
	preCreateData.Data = data
	return
}

// postInsert performs post-insertion operations, including cache handling and post-insert hooks.
//
// After successfully inserting data into the database, this function handles cache updates and invokes
// the PostInsert hook if implemented.
//
// Parameters:
// - cache: A pointer to the Cache interface for cache handling. Can be nil if cache handling is not required.
// - data: The data that was inserted into the database.
// - mData: A map containing data values associated with the inserted entity.
// - dbName: The name of the database where the insertion occurred.
// - entityName: The name of the entity (table or collection) to which the data was inserted.
//
// Returns:
// - err: An error, if any, that occurred during post-insert operations.
//
// Example:
//   user := User{Name: "Alice", Age: 30}
//   err := dbc.postInsert(cache, user, mData, "myDatabase", "users")
//
//   // The function will handle cache updates and invoke PostInsert hook (if implemented).
//
//   data := map[string]interface{}{"name": "Bob", "age": 25}
//   err := dbc.postInsert(cache, data, mData, "myDatabase", "myEntity")
//
//   // The function will handle cache updates and invoke PostInsert hook (if implemented).
//
//   customData := CustomData{...}
//   err := dbc.postInsert(cache, customData, mData, "customDB", "customEntity")
//
//   // The function will handle cache updates and invoke PostInsert hook (if implemented).
func (dbc *DBCommon) postInsert(cache *caches.Cache, data interface{}, mData map[string]interface{}, dbName string, entityName string) error {
	// Check if the data implements the CacheHook interface and cache handling is requested
	if val, ok := interface{}(data).(hooks.CacheHook); ok {
		// Ensure a valid cache instance is available
		if cache != nil {
			// Process cache update based on the CacheHook's cache indexes
			err := caches.GetInstance().ProcessInsertCache(*cache, val.GetCacheIndexes(), mData, dbName, entityName)
			if err != nil {
				return err
			}
		} else {
			// Return an error if no valid cache instance is found
			return dbfusionErrors.ErrNoValidCacheFound
		}
	}

	// Check if the data implements the PostInsert interface and invoke the PostInsert hook if implemented
	if value, ok := interface{}(data).(hooks.PostInsert); ok {
		value = value.PostInsert()
	}

	return nil
}

// preFind prepares the data and options for a find operation.
//
// This function is responsible for setting up the necessary data and conditions for a find operation, including
// determining whether to query the database or retrieve data from the cache. It also handles invoking the PreFind
// hook if implemented on the provided result data.
//
// Parameters:
// - cache: A pointer to the Cache interface for cache handling.
// - result: The result data structure that will be populated with the find results.
// - dbFusionOptions: Optional FindOptions to customize the find operation.
//
// Returns:
// - prefindReturn: A struct containing pre-find information, including entity name, query, where query, and database query flag.
// - err: An error, if any, that occurred during the pre-find process.
//
// Example:
//   var users []User
//   prefindData, err := dbc.preFind(cache, &users, queryoptions.FindOptions{ForceDB: true})
//
//   // The function will prepare data for a find operation, query the database (ForceDB), and return pre-find information.
//
//   var customData []CustomData
//   prefindData, err := dbc.preFind(cache, &customData)
//
//   // The function will prepare data for a find operation, check the cache, and return pre-find information.
func (dbc *DBCommon) preFind(cache *caches.Cache, result interface{}, dbFusionOptions ...queryoptions.FindOptions) (prefindReturn preFindReturn, err error) {
	var nameData entityData
	var options queryoptions.FindOptions

	// Extract FindOptions from the provided arguments, if available
	if len(dbFusionOptions) > 0 {
		options = dbFusionOptions[0]
	}

	// Check if the result data implements the PreFind interface and invoke PreFind hook if implemented
	if value, ok := interface{}(result).(hooks.PreFind); ok {
		result = value.PreFind()
	}

	// Get entity name and related information for the result data
	nameData, err = dbc.getEntityName(result)
	if err != nil {
		return
	}

	// Initialize prefindReturn struct with entity name, data type, and data value
	prefindReturn.entityName = nameData.entityName
	prefindReturn.dataValue = nameData.dataValue
	prefindReturn.dataType = nameData.dataType

	var dbFusionData conditions.DBFusionData

	// Check if dbc.whereQuery is of type conditions.DBFusionData, as it's expected for the whereQuery parameter
	if value, ok := dbc.whereQuery.(conditions.DBFusionData); !ok {
		return prefindReturn, dbfusionErrors.ErrInvalidType
	} else {
		dbFusionData = value
	}

	if options.ForceDB { // Database query is forced, skip cache checks
		prefindReturn.query = dbFusionData.GetQuery()
		prefindReturn.whereQuery = dbc.whereQuery
		prefindReturn.queryDatabase = true
	} else {
		ok := false

		// Construct a cache key for the result data based on database, entity name, and cache values
		redisKey := dbc.currentDB + "_" + prefindReturn.entityName + "_" + dbFusionData.GetCacheValues()

		// Check if the data exists in the cache and retrieve it
		ok, err = caches.GetInstance().ProceessGetCache(*cache, redisKey, result)

		if err != nil {
			return
		}

		skipDB := false

		if !ok { // Data not found in the Redis composite index, check if it exists in the query cache
			// Construct a cache key for the query based on database, entity name, and cache key
			redisQueryKey := dbc.currentDB + "_" + prefindReturn.entityName + "_" + dbFusionData.GetCacheKey()
			// Check if the data exists in the query cache and retrieve it
			skipDB, err = caches.GetInstance().ProceessGetQueryCache(*cache, redisQueryKey, result)
			if err != nil {
				return
			}

			if !skipDB { // Data is not found in the cache or it's not valid to look in the cache for the query, proceed with database query
				prefindReturn.query = dbFusionData.GetQuery()
				prefindReturn.whereQuery = dbc.whereQuery
				prefindReturn.queryDatabase = true
				return
			}
		}

		// Data exists in the cache, set queryDatabase to false (skip database query)
		prefindReturn.queryDatabase = false
	}

	return
}

// postFind handles post-processing tasks after a find operation.
//
// This function is responsible for caching the results of a find operation if requested and invoking the PostFind
// hook if implemented on the result data.
//
// Parameters:
// - cache: A pointer to the Cache interface for cache handling.
// - result: The result data structure obtained from the find operation.
// - entityName: The name of the entity for which the find operation was performed.
// - dbFusionOptions: Optional FindOptions to customize the find operation.
//
// Returns:
// - error: An error, if any, that occurred during post-processing.
//
// Example:
//   var users []User
//   err := dbc.postFind(cache, &users, "users", queryoptions.FindOptions{CacheResult: true})
//
//   // The function will cache the results and invoke the PostFind hook (if implemented) for the "users" entity.
//
//   var customData []CustomData
//   err := dbc.postFind(cache, &customData, "custom_data")
//
//   // The function will cache the results for "custom_data" entity (if CacheResult is true) and not invoke the PostFind hook.
func (dbc *DBCommon) postFind(cache *caches.Cache, result interface{}, entityName string, dbFusionOptions ...queryoptions.FindOptions) error {
	var options queryoptions.FindOptions

	// Extract FindOptions from the provided arguments, if available
	if len(dbFusionOptions) > 0 {
		options = dbFusionOptions[0]
	}

	// Cache the results if CacheResult option is enabled
	if options.CacheResult {
		// Check if the whereQuery is of type conditions.DBFusionData, as caching is only possible for this type
		if value, ok := dbc.whereQuery.(conditions.DBFusionData); ok {
			// Construct a cache key for the query based on database, entity name, and cache key
			redisQueryKey := dbc.currentDB + "_" + entityName + "_" + value.GetCacheKey()
			caches.GetInstance().ProceessSetQueryCache(*cache, redisQueryKey, result)
		}
	}

	// Check if the result data implements the PostFind interface and invoke PostFind hook if implemented
	if value, ok := interface{}(result).(hooks.PostFind); ok {
		dbc.whereQuery = value.PostFind()
	}

	return nil
}

// buildMongoData constructs a MongoDB primitive.D document from the provided data.
//
// This function takes the reflection-based data type and value, extracts field values with valid tags,
// and creates a MongoDB primitive.D document representing the data. Fields with empty "dbfusion" tags or
// zero values are skipped in the document.
//
// Parameters:
// - dataType: The reflection-based data type of the input data.
// - dataValue: The reflection-based data value containing the actual data.
//
// Returns:
// - primitive.D: A MongoDB primitive.D document representing the extracted data.
//
// Example:
//   var user User
//   data := buildMongoData(reflect.TypeOf(user), reflect.ValueOf(user))
//   // 'data' now contains a MongoDB primitive.D document with fields and values from the 'user' struct.
//
//   var customData Custom
//   data := buildMongoData(reflect.TypeOf(customData), reflect.ValueOf(customData))
//   // 'data' now contains a MongoDB primitive.D document with fields and values from the 'customData' struct.
func (dbc *DBCommon) buildMongoData(dataType reflect.Type, dataValue reflect.Value) primitive.D {
	queryMap := primitive.D{}

	// Iterate through the fields of the data structure
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)

		// Extract tags from the field's "dbfusion" tag
		rawtags := strings.Split(field.Tag.Get("dbfusion"), ",")
		tagName := rawtags[0]

		// Skip fields without a valid tag name
		if tagName == "" {
			continue
		}

		// Get the field's value
		value := dataValue.Field(i).Interface()

		// Skip fields with zero values
		if !dbc.isFieldSet(dataValue.Field(i)) {
			continue
		}

		// Create a MongoDB primitive.E element for the field and add it to the document
		singlePoint := primitive.E{Key: tagName, Value: value}
		queryMap = append(queryMap, singlePoint)
	}

	return queryMap
}

// buildMongoUpdate constructs a MongoDB update document based on the provided data and entity information.
//
// This function takes the data and entity information, and depending on the data type and structure type, constructs
// a MongoDB update document suitable for updating records in the database. It supports updating individual fields
// using "$set".
//
// Parameters:
// - data: The data to be used for constructing the update document.
// - nameData: Information about the entity, including its data type and structure type.
//
// Returns:
// - interface{}: A MongoDB update document or nil if the input data is invalid.
// - error: An error if the input data type is not supported.
//
// Example:
//   var userUpdate User
//   updateDoc, err := buildMongoUpdate(userUpdate, nameData)
//   // 'updateDoc' now contains a MongoDB update document suitable for updating 'User' records.
//
//   var customUpdate map[string]interface{}
//   updateDoc, err := buildMongoUpdate(customUpdate, nameData)
//   // 'updateDoc' now contains a MongoDB update document suitable for updating records based on 'customUpdate'.
func (dbc *DBCommon) buildMongoUpdate(data interface{}, nameData entityData) (interface{}, error) {
	dataValue := nameData.dataValue
	dataType := nameData.dataType

	structType := nameData.structType

	var topMap interface{}

	// Check the data type and structure type to determine how to construct the update document
	if structType == 1 { // It's a structure
		queryMap := dbc.buildMongoData(dataType, dataValue)
		topMap = primitive.D{{Key: "$set", Value: queryMap}}
	} else {
		switch data.(type) {
		case ftypes.QMap:
			// Convert data to a primitive.D document
			queryMap := primitive.D{}
			for key, val := range data.(ftypes.QMap) {
				singlePoint := primitive.E{Key: key, Value: val}
				queryMap = append(queryMap, singlePoint)
			}
			topMap = primitive.D{{Key: "$set", Value: queryMap}}
		case ftypes.DMap:
			// Use the data directly as a primitive.D document
			topMap = primitive.D{{Key: "$set", Value: primitive.D(data.(ftypes.DMap))}}
		case map[string]interface{}:
			// Convert data to a primitive.D document
			queryMap := primitive.D{}
			for key, val := range data.(map[string]interface{}) {
				singlePoint := primitive.E{Key: key, Value: val}
				queryMap = append(queryMap, singlePoint)
			}
			topMap = primitive.D{{Key: "$set", Value: queryMap}}
		default:
			return topMap, dbfusionErrors.ErrInvalidType
		}
	}

	return topMap, nil
}

// buildMySqlDeleteData constructs the conditions and values for a MySQL delete query based on the provided data.
//
// This function takes the data and its type, iterates through the fields, and constructs the conditions and values
// required for a MySQL delete query. It creates a WHERE clause with field names and placeholders for values.
//
// Parameters:
// - dataType: The reflect.Type of the data.
// - dataValue: The reflect.Value of the data.
//
// Returns:
// - string: The WHERE clause conditions for the delete query.
// - []interface{}: The values to be used as placeholders in the WHERE clause.
// - error: An error if the input data type is not supported or if there is an issue constructing the conditions.
//
// Example:
//   var userToDelete User
//   conditions, values, err := buildMySqlDeleteData(reflect.TypeOf(userToDelete), reflect.ValueOf(userToDelete))
//   // 'conditions' contains the WHERE clause conditions, and 'values' contains the corresponding values.
func (dbc *DBCommon) buildMySqlDeleteData(dataType reflect.Type, dataValue reflect.Value) (string, []interface{}, error) {
	conditions := ""
	valuesInterface := make([]interface{}, 0)

	// Iterate through the fields of the data type
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)

		rawtags := strings.Split(field.Tag.Get("dbfusion"), ",")
		tagName := rawtags[0]

		if tagName == "" {
			continue
		}
		value := dataValue.Field(i).Interface()

		// Check if the field is set (not zero or nil)
		if !dbc.isFieldSet(dataValue.Field(i)) {
			continue
		}

		// Create the conditions and add the value as a placeholder
		if conditions == "" {
			conditions += fmt.Sprintf("%s = ?", tagName)
		} else {
			conditions += fmt.Sprintf(" AND %s = ?", tagName)
		}
		valuesInterface = append(valuesInterface, value)
	}

	return conditions, valuesInterface, nil
}

// buildMySqlUpdate constructs the SET clause and values for a MySQL update query based on the provided data and entityData.
//
// This function takes the data and its entityData, which contains information about the entity and its type. It iterates
// through the fields of the data and constructs the SET clause for a MySQL update query. It also builds the values to be
// used as placeholders in the query.
//
// Parameters:
// - data: The data for the update query.
// - nameData: The entityData containing information about the entity and its type.
//
// Returns:
// - string: The SET clause for the update query.
// - []interface{}: The values to be used as placeholders in the SET clause.
// - error: An error if the input data type is not supported or if there is an issue constructing the SET clause.
//
// Example:
//   var userToUpdate User
//   entityData, _ := dbc.getEntityName(userToUpdate)
//   setClause, values, err := buildMySqlUpdate(userToUpdate, entityData)
//   // 'setClause' contains the SET clause, and 'values' contains the corresponding values.
func (dbc *DBCommon) buildMySqlUpdate(data interface{}, nameData entityData) (string, []interface{}, error) {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()

	// If the data is a pointer, dereference it
	switch dataType.Kind() {
	case reflect.Ptr:
		ptrValue := reflect.ValueOf(data)
		dataValue = ptrValue.Elem()
		dataType = dataValue.Type()
	}

	structType := 0
	switch dataType.Kind() {
	case reflect.Struct:
		structType = 1
	}
	setString := ""
	valuesInterface := make([]interface{}, 0)

	if structType == 1 { // It's a structure
		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)

			if structType == 1 {
			}
			rawtags := strings.Split(field.Tag.Get("dbfusion"), ",")
			tagName := rawtags[0]

			if tagName == "" {
				continue
			}
			value := dataValue.Field(i).Interface()

			// Check if the field is set (not zero or nil)
			if !dbc.isFieldSet(dataValue.Field(i)) {
				continue
			}
			if setString == "" {
				setString += fmt.Sprintf("%s = ?", tagName)
			} else {
				setString += fmt.Sprintf(",%s = ?", tagName)
			}
			valuesInterface = append(valuesInterface, value)
		}
		setString = "SET " + setString

	} else {
		if value, ok := data.(ftypes.QMap); ok {
			for key, val := range value {
				if setString == "" {
					setString += fmt.Sprintf("%s = ?", key)
				} else {
					setString += fmt.Sprintf(",%s = ?", key)
				}
				valuesInterface = append(valuesInterface, val)
			}
			setString = "SET " + setString
		} else if value, ok := data.(ftypes.DMap); ok {
			for _, val := range value {
				if setString == "" {
					setString += fmt.Sprintf("%s = ?", val.Key)
				} else {
					setString += fmt.Sprintf(",%s = ?", val.Key)
				}
				valuesInterface = append(valuesInterface, val.Value)
			}
			setString = "SET " + setString
		} else if value, ok := data.(map[string]interface{}); ok {
			for key, val := range value {
				if setString == "" {
					setString += fmt.Sprintf("%s = ?", key)
				} else {
					setString += fmt.Sprintf(",%s = ?", key)
				}
				valuesInterface = append(valuesInterface, val)
			}
			setString = "SET " + setString
		} else {
			return "", valuesInterface, dbfusionErrors.ErrInvalidType
		}
	}
	return setString, valuesInterface, nil
}

// preUpdate prepares data for an update operation based on the provided data and database type (dbType).
//
// This function takes the input data and checks if it implements the PreUpdate hook. If so, it calls the PreUpdate
// method to potentially modify the data. It then retrieves entity-related information using the getEntityName function.
// Depending on the database type (dbType), it prepares the data for an update operation, which includes constructing
// the update query or command specific to that database.
//
// Parameters:
// - data: The input data for the update operation.
// - dbType: The type of database (e.g., connections.MONGO) for which the update operation is being prepared.
//
// Returns:
// - preUpdateReturn: A struct containing data prepared for the update operation.
// - error: An error if there is an issue preparing the update data or retrieving entity information.
//
// Example:
//   var userToUpdate User
//   preUpdateData, err := dbc.preUpdate(userToUpdate, connections.MONGO)
//   // 'preUpdateData' contains data prepared for the update operation, specific to MongoDB.
func (dbc *DBCommon) preUpdate(data interface{}, dbType ftypes.DBTypes) (preUpdateData preUpdateReturn, err error) {

	// Check if the data implements the PreUpdate hook and potentially modify it.
	if value, ok := interface{}(data).(hooks.PreUpdate); ok {
		data = value.PreUpdate()
	}

	// Retrieve entity-related information, such as entity name, data type, etc.
	nameData, nameErr := dbc.getEntityName(data)

	if nameErr != nil {
		err = nameErr
		return
	}

	// Initialize the preUpdateData struct with entity-related information.
	preUpdateData.entityName = nameData.entityName
	preUpdateData.dataValue = nameData.dataValue
	preUpdateData.dataType = nameData.dataType
	preUpdateData.structType = nameData.structType

	// Depending on the database type, prepare the data for an update operation.
	if dbType == connections.MONGO {
		preUpdateData.queryData, err = dbc.buildMongoUpdate(data, nameData)
	}

	return
}

// preDelete prepares data for a delete operation based on the provided data.
//
// This function takes the input data and checks if it implements the PreUpdate hook. If so, it calls the PreUpdate
// method to potentially modify the data. It then retrieves entity-related information using the getEntityName function.
// The entity name is required for the delete operation. The prepared data includes entity-related information.
//
// Parameters:
// - data: The input data for the delete operation. If nil, the entity name is obtained from the DBCommon's tableName.
//
// Returns:
// - preDeleteReturn: A struct containing data prepared for the delete operation.
// - error: An error if there is an issue preparing the delete data or retrieving entity information.
//
// Example:
//   var userToDelete User
//   preDeleteData, err := dbc.preDelete(userToDelete)
//   // 'preDeleteData' contains data prepared for the delete operation.
func (dbc *DBCommon) preDelete(data interface{}) (preDeleteData preDeleteReturn, err error) {

	// Check if the data implements the PreUpdate hook and potentially modify it.
	if value, ok := interface{}(data).(hooks.PreUpdate); ok {
		data = value.PreUpdate()
	}

	var nameData entityData

	// If the data is nil, use the entity name from DBCommon's tableName.
	if data == nil {
		nameData.entityName = dbc.tableName
	} else {
		nameData, err = dbc.getEntityName(data)
	}

	if err != nil {
		return
	}

	// Ensure that the entity name is not empty; it's required for the delete operation.
	if nameData.entityName == "" {
		err = dbfusionErrors.ErrEntityNameRequired
		return
	}

	// Initialize the preDeleteData struct with entity-related information.
	preDeleteData.entityName = nameData.entityName
	preDeleteData.dataValue = nameData.dataValue
	preDeleteData.dataType = nameData.dataType

	return
}

// postDelete performs post-delete operations, such as cache management and potential data modification.
//
// This function takes the following actions:
// 1. If the input data implements the CacheHook interface, it retrieves cache-related values for the deleted data and
//    deletes those values from the cache.
// 2. If the input data implements the PostDelete hook, it calls the PostDelete method to potentially modify the data.
//
// Parameters:
// - cache: A pointer to the cache instance used for cache-related operations.
// - data: The input data for the delete operation.
// - entityName: The name of the entity being deleted.
// - results: A map containing the results of the delete operation (e.g., database records that were deleted).
//
// Returns:
// - error: An error if there is an issue performing post-delete operations or if the data modification fails.
//
// Example:
//   var userToDelete User
//   err := dbc.postDelete(cacheInstance, userToDelete, "users", deleteResults)
//   // Perform post-delete operations, such as cache management or data modification.
func (dbc *DBCommon) postDelete(cache *caches.Cache, data interface{}, entityName string, results primitive.M) error {

	// Check if the input data implements the CacheHook interface.
	if value, ok := interface{}(data).(hooks.CacheHook); ok {

		// Build cache-related keys for this data.
		oldValues := dbc.getAllCacheValues(value, results, entityName)

		// Delete all the keys associated with this data from the cache.
		caches.GetInstance().ProceessDeleteCache(*cache, oldValues)
	}

	// Check if the input data implements the PostDelete hook and potentially modify it.
	if value, ok := interface{}(data).(hooks.PostDelete); ok {
		data = value.PostDelete()
	}

	return nil
}

// postUpdate performs post-update operations, including cache management and potential data modification.
//
// This function takes the following actions:
// 1. Updates the cache with new values after an update operation, removing old cache entries.
// 2. If the input data implements the PostUpdate hook, it calls the PostUpdate method to potentially modify the data.
//
// Parameters:
// - cache: A pointer to the cache instance used for cache-related operations.
// - result: The result of the update operation.
// - entityName: The name of the entity being updated.
// - oldValues: A slice of old cache values associated with the updated data.
// - newValues: A slice of new cache values to replace the old ones.
//
// Returns:
// - error: An error if there is an issue performing post-update operations or if the data modification fails.
//
// Example:
//   var updatedData User
//   oldCacheValues := []string{"user:1234", "user:5678"}
//   newCacheValues := []string{"user:7890"}
//   err := dbc.postUpdate(cacheInstance, updatedData, "users", oldCacheValues, newCacheValues)
//   // Perform post-update operations, such as cache management or data modification.
func (dbc *DBCommon) postUpdate(cache *caches.Cache, result interface{}, entityName string, oldValues []string, newValues []string) error {

	// Update the cache with new values, removing old cache entries.
	caches.GetInstance().ProceessUpdateCache(*cache, oldValues, newValues, result)

	// Check if the input data implements the PostUpdate hook and potentially modify it.
	if value, ok := interface{}(result).(hooks.PostUpdate); ok {
		result = value.PostUpdate()
	}

	return nil
}

// getAllCacheValues retrieves all cache keys associated with a data object based on its cache indexes.
//
// This function iterates through the cache indexes defined by the data implementing the CacheHook interface,
// extracts values from the tagValueMap corresponding to each cache index, and constructs cache keys.
// The cache keys are based on the entity name and the values from the tagValueMap, joined by underscores.
//
// Parameters:
// - data: The data object implementing the CacheHook interface, which defines cache indexes.
// - tagValueMap: A map containing tag values extracted from the data object.
// - entityName: The name of the entity for which cache keys are being generated.
//
// Returns:
// - cacheKeys: A slice of cache keys constructed based on cache indexes and tag values.
//
// Example:
//   var userData User
//   tagValues := map[string]interface{}{
//       "id":    1234,
//       "name":  "John",
//   }
//   entity := "users"
//   cacheKeys := dbc.getAllCacheValues(userData, tagValues, entity)
//   // cacheKeys will contain ["testDB_users_1234_John"] based on the cache indexes.
func (dbc *DBCommon) getAllCacheValues(data hooks.CacheHook, tagValueMap map[string]interface{}, entityName string) []string {

	// Initialize an empty slice to store cache keys.
	cacheKeys := make([]string, 0)

	// Iterate through the cache indexes defined by the data object.
	for _, cacheKey := range data.GetCacheIndexes() {
		// Split the cacheKey into internal keys.
		internalKeys := strings.Split(cacheKey, ",")

		// Initialize a cache key string.
		cacheKeyString := ""

		// Iterate through internal keys to construct the cache key.
		for _, internalKey := range internalKeys {
			// Check if the internal key exists in the tagValueMap.
			if value, ok := tagValueMap[internalKey]; ok {
				// Append the value to the cache key string, separated by underscores.
				if cacheKeyString == "" {
					cacheKeyString += fmt.Sprintf("%v", value)
				} else {
					cacheKeyString += fmt.Sprintf("_%v", value)
				}
			}
		}

		// Assemble the final cache key with the current database, entity name, and tag values.
		cacheKeyString = dbc.currentDB + "_" + entityName + "_" + cacheKeyString

		// Append the cache key to the cacheKeys slice.
		cacheKeys = append(cacheKeys, cacheKeyString)
	}

	// Return the slice of constructed cache keys.
	return cacheKeys
}

// refreshValues resets the internal state of the DBCommon instance.
//
// This function sets various properties of the DBCommon instance to their initial or empty values,
// effectively clearing any previously set values or configurations. It's typically used to reset
// the state of the instance for building new database queries.
func (dbc *DBCommon) refreshValues() {
	// Reset various properties to their initial or empty values.
	dbc.tableName = ""
	dbc.whereQuery = nil
	dbc.skip = 0
	dbc.limit = 0
	dbc.projection = nil
	dbc.sort = nil
	dbc.joins = ""
	dbc.groupBy = ""
	dbc.havingString = ""
	dbc.havingValues = make([]interface{}, 0)
	dbc.orderBy = ""
}
