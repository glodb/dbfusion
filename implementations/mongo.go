package implementations

import (
	"context"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/hooks"
	"github.com/glodb/dbfusion/queryoptions"
	"github.com/glodb/dbfusion/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConnection is a struct representing a connection to a MongoDB database with additional query and aggregation options.
type MongoConnection struct {
	DBCommon
	client         *mongo.Client
	match          interface{}
	count          interface{}
	bucket         interface{}
	bucketAuto     interface{}
	addFields      interface{}
	geoNear        interface{}
	group          interface{}
	limitAggregate int
	skipAggregate  int
	sortAggregate  interface{}
	project        interface{}
	sortCount      interface{}
	unset          interface{}
	replaceWith    interface{}
	merge          interface{}
	out            interface{}
	replaceRoot    interface{}
	facet          interface{}
	collStats      interface{}
	indexStats     interface{}
	planCacheStats interface{}
	redact         interface{}
	replaceCount   interface{}
	sample         interface{}
	set            interface{}
	unwind         interface{}
	lookup         interface{}
	graphLookup    interface{}
}

//TODO: add the communication with certificate
func (mc *MongoConnection) ConnectWithCertificate(uri string, filePath string) error {
	// if mc.certSet {
	// 	connectionURI := fmt.Sprintf(connectionStringMain, connectionConfigMap["dbUserMongo"].(string), connectionConfigMap["dbPasswordMongo"].(string), connectionConfigMap["dbHostMongo"].(string)+":"+strconv.Itoa(int(connectionConfigMap["dbPortMongo"].(float64))), connectionConfigMap["dbNameMongo"].(string), readPreference)
	// 	tlsConfig, err := u.getCustomTLSConfig(connectionConfigMap["dbCAFileMongo"].(string))
	// 	if err != nil {
	// 		return nil, errors.New("Unable to get tls config")
	// 	}
	// 	client, err = mongo.NewClient(options.Client().ApplyURI(connectionURI).SetTLSConfig(tlsConfig))
	// 	if err != nil {
	// 		return nil, errors.New("failed to create client")
	// 	}
	// } else {
	// 	connectionURI := fmt.Sprintf(connectionStringDev, connectionConfigMap["dbHostMongo"].(string)+":"+strconv.Itoa(int(connectionConfigMap["dbPortMongo"].(float64))), connectionConfigMap["dbNameMongo"].(string))
	// 	client, err = mongo.NewClient(options.Client().ApplyURI(connectionURI))

	// 	if err != nil {
	// 		return nil, errors.New("failed to create client")
	// 	}
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	// defer cancel()

	// err = client.Connect()
	// if err != nil {
	// 	log.Fatalf("Failed to connect to cluster: %v", err)
	// }

	// // Force a connection to verify our connection string
	// err = client.Ping(ctx, nil)
	// if err != nil {
	// 	log.Fatalf("Failed to ping cluster: %v", err)
	// }
	// u.client = client
	// u.dbName = dbName
	return nil
}

// Connect establishes a connection to a MongoDB server using the provided URI.
// It takes a URI string as a parameter, specifying the MongoDB server to connect to.
//
// Parameters:
// - uri: A string representing the MongoDB server URI.
//
// Returns:
// - An error if the connection cannot be established successfully, or nil on success.
//
// This function performs the following steps:
// 1. Creates client options based on the provided URI.
// 2. Attempts to connect to the MongoDB server using the created client options.
// 3. Checks the connection's status by sending a Ping command to the server.
// 4. If the connection is successful, it stores the connected client in the MongoConnection struct
//    for future use.
//
// Example Usage:
//   mc := MongoConnection{}
//   err := mc.Connect("mongodb://localhost:27017/mydatabase")
//   if err != nil {
//       log.Fatal("Failed to connect to MongoDB:", err)
//   }
//
// The above example initializes a MongoConnection instance and connects it to a MongoDB server
// running on localhost at port 27017, using the "mydatabase" database.
// If the connection is successful, the mc.client field will hold the connected client.
func (mc *MongoConnection) Connect(uri string) error {
	var err error

	// Create client options with the provided URI.
	clientOptions := options.Client().ApplyURI(uri)

	// Attempt to connect to the MongoDB server.
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err // Return an error if the connection attempt fails.
	}

	// Check the connection to ensure it's alive by sending a Ping command.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err // Return an error if the connection is not successful.
	}

	// Set the connected client in the MongoConnection struct for future use.
	mc.client = client

	// Return nil, indicating a successful connection.
	return nil
}

// Table sets the name of the MongoDB collection to operate on.
//
// Parameters:
// - tablename: A string representing the name of the MongoDB collection.
//
// Returns:
// - A reference to the MongoConnection instance, allowing for method chaining.
//
// This method is used to specify the name of the MongoDB collection that subsequent
// database operations will be performed on. It sets the 'tablename' field of the
// MongoConnection instance to the provided collection name.
//
// Example Usage:
//   mc := MongoConnection{}
//   mc.Table("users").FindOne(&result)
//
// In the above example, the 'Table' method is called to set the collection name to "users,"
// and then a 'FindOne' operation is performed on the "users" collection.
// Method chaining is used to make the code more concise and readable.
func (mc *MongoConnection) Table(tablename string) connections.MongoConnection {
	mc.setTable(tablename)
	return mc
}

// InsertOne inserts a single document into the specified MongoDB collection.
//
// Parameters:
// - data: A pointer to the data structure representing the document to be inserted.
//
// Returns:
// - An error if the insertion operation encounters any issues, otherwise returns nil.
//
// This method inserts a single document into the MongoDB collection specified by
// the 'entityName' field of the MongoConnection instance. The document is provided
// as a pointer to the data structure (usually a struct) to be inserted.
//
// Example Usage:
//   mc := MongoConnection{}
//   user := User{Name: "Alice", Age: 30}
//   err := mc.Table("users").InsertOne(&user)
//
// In the above example, the 'InsertOne' method is called to insert a document into
// the "users" collection. The 'user' variable contains the data to be inserted, and it
// is passed as a pointer to the method. Any errors encountered during insertion are
// returned as an error value.
func (mc *MongoConnection) InsertOne(data interface{}) error {
	// Defer the resetting of connection state to ensure cleanup even if an error occurs.
	defer mc.refreshValues()

	// Prepare the data for insertion and handle any pre-insertion operations.
	preCreateData, err := mc.preInsert(data)

	if err != nil {
		return err
	}

	// Use the MongoDB client to insert the document into the specified collection.
	_, err = mc.client.Database(mc.currentDB).Collection(preCreateData.entityName).InsertOne(context.TODO(), preCreateData.mData)

	if err == nil {
		// Handle any post-insertion operations, such as caching.
		err = mc.postInsert(mc.cache, preCreateData.Data, preCreateData.mData, mc.currentDB, preCreateData.entityName)
	}
	return err
}

// FindOne retrieves a single document from the specified MongoDB collection based on the provided query conditions.
//
// Parameters:
// - result: A pointer to the data structure where the retrieved document will be decoded into.
// - dbFusionOptions: Optional database fusion options that can be applied to the query.
//
// Returns:
// - An error if the retrieval operation encounters any issues, otherwise returns nil.
//
// This method retrieves a single document from the MongoDB collection specified by the 'entityName' field
// of the MongoConnection instance. The retrieved document is decoded into the data structure pointed to
// by 'result'. The query conditions are determined by the 'whereQuery', 'projection', 'skip', and 'sort' fields
// of the MongoConnection instance.
//
// Example Usage:
//   mc := MongoConnection{}
//   var user User
//   err := mc.Table("users").Where(MongoData{"name": "Alice"}).FindOne(&user)
//
// In the above example, the 'FindOne' method is used to retrieve a single document from the "users" collection
// where the "name" field matches "Alice". The retrieved data is decoded into the 'user' variable.
func (mc *MongoConnection) FindOne(result interface{}, dbFusionOptions ...queryoptions.FindOptions) error {
	// Defer the resetting of connection state to ensure cleanup even if an error occurs.
	defer mc.refreshValues()

	// Check if there are specific query conditions provided in 'whereQuery'.
	if mc.whereQuery != nil {
		// Convert the 'whereQuery' into a MongoDB-compatible query.
		query, err := utils.GetInstance().GetMongoFusionData(mc.whereQuery)
		if err != nil {
			return err
		}
		// Update the 'whereQuery' with the MongoDB-compatible query.
		mc.whereQuery = query
	} else {
		// If no query conditions are provided, initialize 'whereQuery' as an empty MongoData.
		mc.whereQuery = &conditions.MongoData{}
	}

	// Prepare for pre-find operations and retrieve pre-find data.
	prefindReturn, err := mc.preFind(mc.cache, result, dbFusionOptions...)
	if err != nil {
		return err
	}

	// Check if the query should be executed against the database.
	if prefindReturn.queryDatabase {
		// Create options for the FindOne operation, including projection, skip, and sort.
		opts := options.FindOneOptions{}
		if mc.projection != nil {
			opts.SetProjection(mc.projection)
		}
		if mc.skip != 0 {
			opts.SetSkip(mc.skip)
		}
		if mc.sort != nil {
			opts.SetSort(mc.sort)
		}

		// Execute the FindOne operation to retrieve a single document.
		err = mc.client.Database(mc.currentDB).Collection(prefindReturn.entityName).FindOne(context.TODO(), prefindReturn.query, &opts).Decode(result)
		if err != nil {
			return err
		}
	}

	// Handle any post-find operations, such as caching.
	err = mc.postFind(mc.cache, result, prefindReturn.entityName, dbFusionOptions...)
	return err
}

// UpdateAndFindOne updates a document in the specified MongoDB collection based on provided data and query conditions,
// and then retrieves the updated document.
//
// Parameters:
// - data: The data used to update the document.
// - result: A pointer to the data structure where the updated document will be decoded into.
// - upsert: If true, insert the document if it does not exist.
//
// Returns:
// - An error if the update operation or retrieval encounters any issues, otherwise returns nil.
//
// This method updates a document in the MongoDB collection specified by the 'entityName' field of the MongoConnection instance.
// The update is performed based on the provided 'data' and query conditions specified in the 'whereQuery'.
// After the update, the method retrieves the updated document and decodes it into the data structure pointed to by 'result'.
//
// Example Usage:
//   mc := MongoConnection{}
//   var updatedUser User
//   err := mc.Table("users").Where(MongoData{"name": "Alice"}).UpdateAndFindOne(updateData, &updatedUser, true)
//
// In the above example, the 'UpdateAndFindOne' method is used to update a document in the "users" collection where
// the "name" field matches "Alice". The updated document is decoded into the 'updatedUser' variable, and if it doesn't exist, it is inserted.
func (mc *MongoConnection) UpdateAndFindOne(data interface{}, result interface{}, upsert bool) error {
	// Defer the resetting of connection state to ensure cleanup even if an error occurs.
	defer mc.refreshValues()

	var fusionQuery conditions.DBFusionData
	// Check if there are specific query conditions provided in 'whereQuery'.
	if mc.whereQuery != nil {
		// Convert the 'whereQuery' into a MongoDB-compatible query.
		query, err := utils.GetInstance().GetMongoFusionData(mc.whereQuery)
		if err != nil {
			return err
		}
		fusionQuery = query
	} else {
		// If no query conditions are provided, initialize 'whereQuery' as an empty MongoData.
		fusionQuery = &conditions.MongoData{}
	}

	// Prepare for pre-update operations and retrieve pre-update data.
	preUpdateReturn, err := mc.preUpdate(data, connections.MONGO)
	if err != nil {
		return err
	}

	// Create options for the FindOneAndUpdate operation, including projection, sort, upsert, and return document.
	opts := options.FindOneAndUpdateOptions{}
	if mc.projection != nil {
		opts.SetProjection(mc.projection)
	}
	if mc.sort != nil {
		opts.SetSort(mc.sort)
	}
	opts.SetUpsert(upsert)
	opts.SetReturnDocument(options.After)

	// Initialize variables for cache update.
	updateCache := false
	oldKeys := []string{}
	newKeys := []string{}
	var cacheHook hooks.CacheHook

	// Check if the 'result' implements the CacheHook interface.
	if value, ok := interface{}(result).(hooks.CacheHook); ok {
		// Attempt to retrieve the existing document before the update.
		err = mc.client.Database(mc.currentDB).Collection(preUpdateReturn.entityName).FindOne(context.TODO(), fusionQuery.GetQuery().(primitive.D)).Decode(result)
		if err != nil {
			return err
		}

		// Create a tag map from the updated document.
		tagMapValue, err := mc.createTagValueMap(result)
		if err == nil {
			// Get old cache keys before the update.
			oldKeys = mc.getAllCacheValues(value, tagMapValue, preUpdateReturn.entityName)
			updateCache = true
			cacheHook = value
		}
	}

	// Perform the FindOneAndUpdate operation to update and retrieve the document.
	err = mc.client.Database(mc.currentDB).Collection(preUpdateReturn.entityName).FindOneAndUpdate(
		context.TODO(),
		fusionQuery.GetQuery().(primitive.D),
		preUpdateReturn.queryData.(primitive.D),
		&opts,
	).Decode(result)
	if err != nil {
		return err
	}

	// If cache update is needed, get new cache keys after the update.
	if updateCache {
		tagMapValue, _ := mc.createTagValueMap(result)
		newKeys = mc.getAllCacheValues(cacheHook, tagMapValue, preUpdateReturn.entityName)
	}

	// Handle any post-update operations, such as caching.
	err = mc.postUpdate(mc.cache, result, preUpdateReturn.entityName, oldKeys, newKeys)
	return err
}

// DeleteOne deletes a document from the specified MongoDB collection based on provided query conditions or data.
//
// Parameters:
// - sliceData: A variadic parameter that accepts optional data to identify the document to be deleted. If provided, the document is identified and deleted based on the data provided.
//
// Returns:
// - An error if the deletion operation encounters any issues, otherwise returns nil.
//
// This method deletes a document from the MongoDB collection specified by the 'entityName' field of the MongoConnection instance.
// The document to be deleted can be identified either by query conditions specified in the 'whereQuery' or by providing specific data.
//
// Example Usage:
//   mc := MongoConnection{}
//   err := mc.Table("users").Where(MongoData{"name": "Alice"}).DeleteOne()
//
// In the above example, the 'DeleteOne' method is used to delete a document in the "users" collection where
// the "name" field matches "Alice".
func (mc *MongoConnection) DeleteOne(sliceData ...interface{}) error {
	// Defer the resetting of connection state to ensure cleanup even if an error occurs.
	defer mc.refreshValues()

	var data interface{}
	// Check if optional data is provided to identify the document to be deleted.
	if len(sliceData) != 0 {
		data = sliceData[0]
	}

	// Check if there are specific query conditions provided in 'whereQuery'.
	if mc.whereQuery != nil {
		// Convert the 'whereQuery' into a MongoDB-compatible query.
		query, err := utils.GetInstance().GetMongoFusionData(mc.whereQuery)
		if err != nil {
			return err
		}
		mc.whereQuery = query
	} else {
		// If no query conditions are provided, initialize 'whereQuery' as an empty MongoData.
		mc.whereQuery = &conditions.MongoData{}
	}

	// Prepare for pre-delete operations and retrieve pre-delete data.
	preDeleteData, err := mc.preDelete(data)
	if err != nil {
		return err
	}

	var deleteQuery primitive.D
	var results primitive.M

	// Check if specific data is provided for document identification (delete by data).
	if data != nil {
		// Build a MongoDB-compatible query to identify the document based on data.
		deleteQuery = mc.buildMongoData(preDeleteData.dataType, preDeleteData.dataValue)

		// Attempt to find and delete the document identified by the query.
		err = mc.client.Database(mc.currentDB).Collection(preDeleteData.entityName).FindOneAndDelete(context.TODO(), deleteQuery).Decode(&results)
	} else {
		// Delete documents based on query conditions (delete by query).

		// Simple delete operation without checking the cache, as cache is not relevant in this case.
		_, err = mc.client.Database(mc.currentDB).Collection(preDeleteData.entityName).DeleteOne(context.TODO(), mc.whereQuery.(conditions.DBFusionData).GetQuery())
	}

	if err != nil {
		return err
	}

	// Handle any post-delete operations, such as cache updates.
	err = mc.postDelete(mc.cache, data, preDeleteData.entityName, results)

	return err
}

// DisConnect closes the connection to the MongoDB server.
//
// Returns:
// - An error if the disconnection encounters any issues, otherwise returns nil.
//
// This method gracefully closes the connection to the MongoDB server that was previously established using the 'Connect' method.
// It is important to call this method when you are finished using the MongoDB connection to release resources and maintain proper cleanup.
//
// Example Usage:
//   mc := MongoConnection{}
//   err := mc.Connect("mongodb://localhost:27017")
//   if err != nil {
//       log.Fatal("Failed to connect to MongoDB:", err)
//   }
//   // Perform MongoDB operations...
//   err = mc.DisConnect() // Close the MongoDB connection when done.
func (mc *MongoConnection) DisConnect() error {
	// Close the MongoDB client connection gracefully.
	return mc.client.Disconnect(context.TODO())
}

// Paginate performs pagination on a MongoDB query and retrieves a specific page of results.
//
// Parameters:
// - results: A pointer to the slice where the query results should be stored.
// - pageNumber: The page number to retrieve (1-based index).
//
// Returns:
// - PaginationResults: A struct containing pagination information.
// - error: An error if the pagination encounters any issues, otherwise returns nil.
//
// This method is used to paginate the results of a MongoDB query, allowing you to retrieve a specific page of data.
// It takes a pointer to a slice where the query results should be stored and the page number to retrieve.
// The method calculates the total number of documents matching the query, total pages, and other pagination details.
// It then retrieves the specified page of results and populates the provided slice with the data.
//
// Example Usage:
//   mc := MongoConnection{}
//   err := mc.Connect("mongodb://localhost:27017")
//   if err != nil {
//       log.Fatal("Failed to connect to MongoDB:", err)
//   }
//   var results []YourDataType
//   pageNumber := 1
//   paginationInfo, err := mc.Paginate(&results, pageNumber)
//   if err != nil {
//       log.Fatal("Failed to paginate query:", err)
//   }
//   // Process the paginated results and use paginationInfo to display pagination controls.
func (mc *MongoConnection) Paginate(results interface{}, pageNumber int) (connections.PaginationResults, error) {
	// Ensure that MongoDB Fusion data is available for the query.
	if mc.whereQuery != nil {
		query, err := utils.GetInstance().GetMongoFusionData(mc.whereQuery)
		if err != nil {
			return connections.PaginationResults{}, err
		}
		mc.whereQuery = query
	} else {
		mc.whereQuery = &conditions.MongoData{}
	}

	// Initialize pagination results.
	var paginationResults connections.PaginationResults

	// Configure options for the MongoDB query.
	opts := options.FindOptions{}
	if mc.projection != nil {
		opts.SetProjection(mc.projection)
	}

	// Count the total number of documents matching the query.
	count, err := mc.client.Database(mc.currentDB).Collection(mc.tableName).CountDocuments(context.TODO(), mc.whereQuery.(conditions.DBFusionData).GetQuery())
	if err != nil {
		return connections.PaginationResults{}, err
	}

	// Populate paginationResults with count and pagination details.
	paginationResults.TotalDocuments = count
	paginationResults.TotalPages = int64(math.Ceil((float64(count) / float64(mc.pageSize))))
	paginationResults.Limit = int64(mc.pageSize)
	paginationResults.CurrentPage = int64(pageNumber)

	// Calculate the skip and limit based on the requested page number.
	mc.limit = int64(mc.pageSize)
	mc.skip = int64((pageNumber - 1) * mc.pageSize)
	opts.SetSkip(mc.skip)
	opts.SetSort(mc.sort)
	opts.SetLimit(mc.limit)

	// Execute the MongoDB query with pagination options.
	cursor, err := mc.client.Database(mc.currentDB).Collection(mc.tableName).Find(context.TODO(), mc.whereQuery.(conditions.DBFusionData).GetQuery(), &opts)
	if err != nil {
		return connections.PaginationResults{}, err
	}

	// Decode and store the results in the provided slice.
	if err = cursor.All(context.TODO(), results); err != nil {
		return connections.PaginationResults{}, err
	}

	return paginationResults, nil
}

func (mc *MongoConnection) InsertMany(interface{}) error {
	return nil
}
func (mc *MongoConnection) FindMany(interface{}, ...queryoptions.FindOptions) error {
	return nil
}
func (mc *MongoConnection) UpdateMany(interface{}, interface{}, bool) error {
	return nil
}
func (mc *MongoConnection) DeleteMany(...interface{}) error {
	return nil
}

// Skip sets the number of documents to skip in a MongoDB query.
//
// Parameters:
// - skip: The number of documents to skip.
//
// Returns:
// - connections.MongoConnection: A reference to the MongoConnection for method chaining.
//
// This method is used to specify how many documents should be skipped in the result set of a MongoDB query.
// It updates the skip value in the MongoConnection object and can be used in method chaining.
func (mc *MongoConnection) Skip(skip int64) connections.MongoConnection {
	mc.skip = skip
	return mc
}

// Limit sets the maximum number of documents to return in a MongoDB query.
//
// Parameters:
// - limit: The maximum number of documents to return.
//
// Returns:
// - connections.MongoConnection: A reference to the MongoConnection for method chaining.
//
// This method is used to specify the maximum number of documents that should be returned in the result set of a MongoDB query.
// It updates the limit value in the MongoConnection object and can be used in method chaining.
func (mc *MongoConnection) Limit(limit int64) connections.MongoConnection {
	mc.limit = limit
	return mc
}

// Select specifies which fields to include or exclude in the query results.
//
// Parameters:
// - keys: A map of field names and a boolean flag indicating whether to include (true) or exclude (false) the field.
//
// Returns:
// - connections.MongoConnection: A reference to the MongoConnection for method chaining.
//
// This method is used to define which fields should be included or excluded in the query results.
// It takes a map where the keys are field names, and the values are boolean flags (true to include, false to exclude).
// The specified fields will be projected into the query results, and others will be omitted.
// It updates the projection value in the MongoConnection object and can be used in method chaining.
func (mc *MongoConnection) Select(keys map[string]bool) connections.MongoConnection {
	selectionKeys := make(map[string]int, 0)

	for key, val := range keys {
		if val {
			selectionKeys[key] = 1
		} else {
			selectionKeys[key] = 0
		}
	}
	mc.projection = selectionKeys
	return mc
}

// Sort specifies the sorting order for the query results based on a field and an optional descending flag.
//
// Parameters:
// - sortKey: The field by which to sort the results.
// - sortdesc: An optional boolean flag indicating descending sorting (true for descending, false for ascending).
//
// Returns:
// - connections.MongoConnection: A reference to the MongoConnection for method chaining.
//
// This method is used to define the sorting order for the query results based on a field.
// You can specify a field to sort by and, optionally, set it to descending order.
// If called multiple times, the sorting fields will be combined.
// It updates the sort value in the MongoConnection object and can be used in method chaining.
func (mc *MongoConnection) Sort(sortKey string, sortdesc ...bool) connections.MongoConnection {
	sortString := sortKey
	sortVal := 1
	if len(sortdesc) > 0 {
		if !sortdesc[0] {
			sortVal = -1
		}
	}

	if mc.sort != nil {
		sortMap := mc.sort.(map[string]interface{})
		sortMap[sortString] = sortVal
		mc.sort = sortMap
	} else {
		sortMap := make(map[string]interface{})
		mc.sort = sortMap
	}
	return mc
}

// Where specifies the query conditions to filter MongoDB query results.
//
// Parameters:
// - query: The query conditions to filter the results.
//
// Returns:
// - connections.MongoConnection: A reference to the MongoConnection for method chaining.
//
// This method is used to set query conditions for filtering MongoDB query results.
// It takes a query object that defines the filtering criteria.
// It updates the whereQuery value in the MongoConnection object and can be used in method chaining.
func (mc *MongoConnection) Where(query interface{}) connections.MongoConnection {
	mc.whereQuery = query
	return mc
}

// SetPageSize sets the maximum number of documents to return per page in paginated queries.
//
// Parameters:
// - limit: The maximum number of documents per page.
//
// This method is used to set the maximum number of documents to return per page in paginated queries.
// It updates the pageSize value in the MongoConnection object.
func (mc *MongoConnection) SetPageSize(limit int) {
	mc.pageSize = limit
}

// parseSortableIndexes parses a list of sortable indexes and converts them into MongoDB-compatible index definitions.
//
// Parameters:
// - indexes: A slice of sortable index definitions in the format "field:order", e.g., "field1:1,field2:-1".
//
// Returns:
// - []bson.D: A slice of MongoDB-compatible index definitions, where each definition is a BSON document (bson.D).
//
// This method is used to convert a list of sortable indexes into MongoDB-compatible index definitions.
// Each sortable index definition is in the format "field:order", where "field" is the index field name,
// and "order" indicates the sorting order (1 for ascending, -1 for descending).
// The method parses these definitions and constructs BSON documents for each index.
// The resulting slice contains MongoDB-compatible index definitions.
func (mc *MongoConnection) parseSortableIndexes(indexes []string) []bson.D {
	mongoIndexes := []bson.D{}
	for _, index := range indexes {
		compounds := strings.Split(index, ",")
		localIndex := bson.D{}
		for _, compoundIndex := range compounds {
			indexValues := strings.Split(compoundIndex, ":")
			if len(indexValues) > 1 {
				intVal, _ := strconv.ParseInt(indexValues[1], 10, 64)
				singleIndex := bson.E{Key: indexValues[0], Value: intVal}
				localIndex = append(localIndex, singleIndex)
			}
		}
		mongoIndexes = append(mongoIndexes, localIndex)
	}
	return mongoIndexes
}

// parseIndexes parses a list of index fields and converts them into MongoDB-compatible index definitions.
//
// Parameters:
// - indexes: A slice of index field names.
// - indexName: The name of the index.
//
// Returns:
// - []bson.D: A slice of MongoDB-compatible index definitions, where each definition is a BSON document (bson.D).
//
// This method is used to convert a list of index fields into MongoDB-compatible index definitions.
// Each index definition consists of the field name and the index name.
// The method constructs BSON documents for each index with the specified field and index name.
// The resulting slice contains MongoDB-compatible index definitions.
func (mc *MongoConnection) parseIndexes(indexes []string, indexName string) []bson.D {
	mongoIndexes := []bson.D{}
	for _, index := range indexes {
		compounds := strings.Split(index, ",")
		localIndex := bson.D{}
		for _, compoundIndex := range compounds {
			singleIndex := bson.E{Key: compoundIndex, Value: indexName}
			localIndex = append(localIndex, singleIndex)
		}
		mongoIndexes = append(mongoIndexes, localIndex)
	}
	return mongoIndexes
}

// CreateIndexes creates indexes on a MongoDB collection based on the specified data structure's index configurations.
//
// Parameters:
// - data: The data structure for which indexes should be created. It should implement the appropriate hooks interfaces
//         (NormalIndexes, UniqueIndexes, TextIndexes, TwoDimensionalIndexes, TwoDimensionalSpatialIndexes,
//         HashedIndexes, SparseIndexes) to specify index configurations.
//
// Returns:
// - error: An error if index creation fails; otherwise, it returns nil.
//
// This method creates indexes on a MongoDB collection based on the index configurations specified in the provided data structure.
// It uses hooks interfaces to determine which indexes to create and their configurations.
func (mc *MongoConnection) CreateIndexes(data interface{}) error {
	// Get the entity name from the provided data structure.
	name, _ := mc.getEntityName(data)

	// Create Normal Indexes if the data structure implements the NormalIndexes hook.
	if val, ok := interface{}(data).(hooks.NormalIndexes); ok {
		indexes := val.GetNormalIndexes()
		compundIndexes := mc.parseSortableIndexes(indexes)
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys: index,
			}
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				return err
			}
		}
	}

	// Create Unique Indexes if the data structure implements the UniqueIndexes hook.
	if val, ok := interface{}(data).(hooks.UniqueIndexes); ok {
		indexes := val.GetUniqueIndexes()
		compundIndexes := mc.parseSortableIndexes(indexes)
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys:    index,
				Options: options.Index().SetUnique(true),
			}
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				return err
			}
		}
	}

	// Create Text Index if the data structure implements the TextIndexes hook.
	if val, ok := interface{}(data).(hooks.TextIndexes); ok {
		indexes := val.GetTextIndex()
		compundIndexes := mc.parseIndexes([]string{indexes}, "text")
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys: index,
			}
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				return err
			}
		}
	}

	// Create 2D Indexes if the data structure implements the TwoDimensionalIndexes hook.
	if val, ok := interface{}(data).(hooks.TwoDimensionalIndexes); ok {
		indexes := val.Get2DIndexes()
		compundIndexes := mc.parseIndexes(indexes, "2d")
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys: index,
			}
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				return err
			}
		}
	}

	// Create 2DSphere Indexes if the data structure implements the TwoDimensionalSpatialIndexes hook.
	if val, ok := interface{}(data).(hooks.TwoDimensionalSpatialIndexes); ok {
		indexes := val.Get2DSpatialIndexes()
		compundIndexes := mc.parseIndexes(indexes, "2dsphere")
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys: index,
			}
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				return err
			}
		}
	}

	// Create Hashed Indexes if the data structure implements the HashedIndexes hook.
	if val, ok := interface{}(data).(hooks.HashedIndexes); ok {
		indexes := val.GetHashedIndexes()
		compundIndexes := mc.parseIndexes(indexes, "hashed")
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys: index,
			}
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				return err
			}
		}
	}

	// Create Sparse Indexes if the data structure implements the SparseIndexes hook.
	if val, ok := interface{}(data).(hooks.SparseIndexes); ok {
		indexes := val.GetSparseIndexes()
		compundIndexes := mc.parseSortableIndexes(indexes)
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys:    index,
				Options: options.Index().SetSparse(true),
			}
			log.Println(mc.parseSortableIndexes(indexes))
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			log.Println(err)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Match sets the $match aggregation stage to filter documents that match the specified criteria.
func (mc *MongoConnection) Match(data interface{}) connections.MongoConnection {
	mc.match = data
	return mc
}

// Bucket sets the $bucket aggregation stage to group documents into buckets based on the specified criteria.
func (mc *MongoConnection) Bucket(data interface{}) connections.MongoConnection {
	mc.bucket = data
	return mc
}

// BucketsAuto sets the $bucketsAuto aggregation stage to group documents into automatic buckets based on the specified criteria.
func (mc *MongoConnection) BucketsAuto(data interface{}) connections.MongoConnection {
	mc.bucketAuto = data
	return mc
}

// AddFields sets the $addFields aggregation stage to add new fields to documents based on the specified expressions.
func (mc *MongoConnection) AddFields(data interface{}) connections.MongoConnection {
	mc.addFields = data
	return mc
}

// GeoNear sets the $geoNear aggregation stage to perform geospatial queries and return documents in proximity to a specified point.
func (mc *MongoConnection) GeoNear(data interface{}) connections.MongoConnection {
	mc.geoNear = data
	return mc
}

// Group sets the $group aggregation stage to group documents based on specified fields and perform aggregation operations.
func (mc *MongoConnection) Group(data interface{}) connections.MongoConnection {
	mc.group = data
	return mc
}

// LimitAggregate sets the $limit aggregation stage to limit the number of documents in the aggregation pipeline.
func (mc *MongoConnection) LimitAggregate(data int) connections.MongoConnection {
	mc.limitAggregate = data
	return mc
}

// SkipAggregate sets the $skip aggregation stage to skip a specified number of documents in the aggregation pipeline.
func (mc *MongoConnection) SkipAggregate(data int) connections.MongoConnection {
	mc.skipAggregate = data
	return mc
}

// SortAggregate sets the $sort aggregation stage to sort documents in the aggregation pipeline based on the specified criteria.
func (mc *MongoConnection) SortAggregate(data interface{}) connections.MongoConnection {
	mc.sortAggregate = data
	return mc
}

// SortByCount sets the $sortByCount aggregation stage to perform a count operation and then sort the result documents.
func (mc *MongoConnection) SortByCount(data interface{}) connections.MongoConnection {
	mc.sortCount = data
	return mc
}

// Project sets the $project aggregation stage to reshape documents and include or exclude fields as specified.
func (mc *MongoConnection) Project(data interface{}) connections.MongoConnection {
	mc.project = data
	return mc
}

// Unset sets the $unset aggregation stage to remove specified fields from documents.
func (mc *MongoConnection) Unset(data interface{}) connections.MongoConnection {
	mc.unset = data
	return mc
}

// ReplaceWith sets the $replaceWith aggregation stage to replace documents with the specified expression.
func (mc *MongoConnection) ReplaceWith(data interface{}) connections.MongoConnection {
	mc.replaceWith = data
	return mc
}

// Merge sets the $merge aggregation stage to merge documents from different collections.
func (mc *MongoConnection) Merge(data interface{}) connections.MongoConnection {
	mc.merge = data
	return mc
}

// Out sets the $out aggregation stage to write the result of the aggregation pipeline to a specified collection.
func (mc *MongoConnection) Out(data interface{}) connections.MongoConnection {
	mc.out = data
	return mc
}

// Facet sets the $facet aggregation stage to apply multiple pipelines to the same input documents.
func (mc *MongoConnection) Facet(data interface{}) connections.MongoConnection {
	mc.facet = data
	return mc
}

// CollStats sets the $collStats aggregation stage to return statistics for a specified collection.
func (mc *MongoConnection) CollStats(data interface{}) connections.MongoConnection {
	mc.collStats = data
	return mc
}

// IndexStats sets the $indexStats aggregation stage to return statistics for a specified collection's indexes.
func (mc *MongoConnection) IndexStats(data interface{}) connections.MongoConnection {
	mc.indexStats = data
	return mc
}

// PlanCacheStats sets the $planCacheStats aggregation stage to return statistics for a specified collection's query plan cache.
func (mc *MongoConnection) PlanCacheStats(data interface{}) connections.MongoConnection {
	mc.planCacheStats = data
	return mc
}

// Redact sets the $redact aggregation stage to control the access to the data in documents.
func (mc *MongoConnection) Redact(data interface{}) connections.MongoConnection {
	mc.redact = data
	return mc
}

// ReplaceRoot sets the $replaceRoot aggregation stage to replace the root document with a specified expression.
func (mc *MongoConnection) ReplaceRoot(data interface{}) connections.MongoConnection {
	mc.replaceRoot = data
	return mc
}

// ReplaceCount sets the $replaceCount aggregation stage to replace documents with a specified count.
func (mc *MongoConnection) ReplaceCount(data interface{}) connections.MongoConnection {
	mc.replaceCount = data
	return mc
}

// Sample sets the $sample aggregation stage to randomly sample documents from a collection.
func (mc *MongoConnection) Sample(data interface{}) connections.MongoConnection {
	mc.sample = data
	return mc
}

// Set sets the $set aggregation stage to add new fields to documents based on the specified expressions.
func (mc *MongoConnection) Set(data interface{}) connections.MongoConnection {
	mc.set = data
	return mc
}

// Unwind sets the $unwind aggregation stage to deconstruct an array field and output one document for each element.
func (mc *MongoConnection) Unwind(data interface{}) connections.MongoConnection {
	mc.unwind = data
	return mc
}

// Lookup sets the $lookup aggregation stage to perform a left outer join between documents from two collections.
func (mc *MongoConnection) Lookup(data interface{}) connections.MongoConnection {
	mc.lookup = data
	return mc
}

// GraphLookup sets the $graphLookup aggregation stage to perform a recursive search on a graph structure.
func (mc *MongoConnection) GraphLookup(data interface{}) connections.MongoConnection {
	mc.graphLookup = data
	return mc
}

// Count sets the $count aggregation stage to return the number of documents in the aggregation pipeline.
func (mc *MongoConnection) Count(data interface{}) connections.MongoConnection {
	mc.count = data
	return mc
}

// refreshAggregation resets all aggregation settings and options in the MongoConnection,
// preparing it for a new aggregation operation with default values.
func (mc *MongoConnection) refreshAggregation() {
	// Reset all aggregation settings to their default values.
	mc.match = nil
	mc.count = nil
	mc.bucket = nil
	mc.bucketAuto = nil
	mc.addFields = nil
	mc.geoNear = nil
	mc.group = nil
	mc.limitAggregate = 0
	mc.skipAggregate = 0
	mc.sortAggregate = nil
	mc.project = nil
	mc.sortCount = nil
	mc.unset = nil
	mc.replaceWith = nil
	mc.merge = nil
	mc.out = nil
	mc.replaceRoot = nil
	mc.facet = nil
	mc.collStats = nil
	mc.indexStats = nil
	mc.planCacheStats = nil
	mc.redact = nil
	mc.replaceCount = nil
	mc.sample = nil
	mc.set = nil
	mc.unwind = nil
	mc.lookup = nil
	mc.graphLookup = nil
}

// createAggregation generates the aggregation pipeline based on the configured aggregation options.
// It constructs an array of BSON documents (pipeline stages) to be used in the aggregation.
func (mc *MongoConnection) createAggregation() bson.A {
	// Initialize an empty array to store the aggregation pipeline stages.
	pipelines := primitive.A{}

	// Add "$match" stage to filter documents based on the specified conditions.
	if mc.match != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$match", Value: mc.match}})
	}

	// Add "$count" stage to count the number of documents in the aggregation result.
	if mc.count != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$count", Value: mc.count}})
	}

	// Add "$bucket" stage for bucketing documents into specified ranges.
	if mc.bucket != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$bucket", Value: mc.bucket}})
	}

	// Add "$bucketAuto" stage for automatically bucketing documents into specified ranges.
	if mc.bucketAuto != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$bucketAuto", Value: mc.bucketAuto}})
	}

	// Add "$addFields" stage to add new fields to documents.
	if mc.addFields != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$addFields", Value: mc.addFields}})
	}

	// Add "$geoNear" stage for geospatial queries.
	if mc.geoNear != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$geoNear", Value: mc.geoNear}})
	}

	// Add "$group" stage to group documents based on specified criteria.
	if mc.group != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$group", Value: mc.group}})
	}

	// Add "$skip" stage to skip a specified number of documents in the result.
	if mc.skipAggregate != 0 {
		pipelines = append(pipelines, primitive.D{{Key: "$skip", Value: mc.skipAggregate}})
	}

	// Add "$limit" stage to limit the number of documents in the result.
	if mc.limitAggregate != 0 {
		pipelines = append(pipelines, primitive.D{{Key: "$limit", Value: mc.limitAggregate}})
	}

	// Add "$sort" stage to sort documents based on specified criteria.
	if mc.sortAggregate != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$sort", Value: mc.sortAggregate}})
	}

	// Add "$project" stage to reshape the documents in the result.
	if mc.project != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$project", Value: mc.project}})
	}

	// Add "$sortCount" stage to sort the result by the count.
	if mc.sortCount != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$sortCount", Value: mc.sortCount}})
	}

	// Add "$unset" stage to remove specified fields from documents.
	if mc.unset != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$unset", Value: mc.unset}})
	}

	// Add "$replaceWith" stage to replace documents with specified values.
	if mc.replaceWith != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$replaceWith", Value: mc.replaceWith}})
	}

	// Add "$merge" stage to merge documents into a target collection.
	if mc.merge != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$merge", Value: mc.merge}})
	}

	// Add "$out" stage to specify the output collection for the aggregation result.
	if mc.out != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$out", Value: mc.out}})
	}

	// Add "$replaceRoot" stage to replace the root of documents with specified values.
	if mc.replaceRoot != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$replaceRoot", Value: mc.replaceRoot}})
	}

	// Add "$facet" stage for multi-faceted aggregations.
	if mc.facet != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$facet", Value: mc.facet}})
	}

	// Add "$collStats" stage to collect statistics about a collection.
	if mc.collStats != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$collStats", Value: mc.collStats}})
	}

	// Add "$indexStats" stage to collect statistics about indexes on a collection.
	if mc.indexStats != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$indexStats", Value: mc.indexStats}})
	}

	// Add "$planCacheStats" stage to collect query plan cache statistics.
	if mc.planCacheStats != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$planCacheStats", Value: mc.planCacheStats}})
	}

	// Add "$redact" stage to restrict data based on specified conditions.
	if mc.redact != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$redact", Value: mc.redact}})
	}

	// Add "$replaceCount" stage to replace the count of documents in the result.
	if mc.replaceCount != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$replaceCount", Value: mc.replaceCount}})
	}

	// Add "$sample" stage to randomly sample documents from the result.
	if mc.sample != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$sample", Value: mc.sample}})
	}

	// Add "$set" stage to set values in documents.
	if mc.set != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$set", Value: mc.set}})
	}

	// Add "$unwind" stage to unwind arrays in documents.
	if mc.unwind != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$unwind", Value: mc.unwind}})
	}

	// Add "$lookup" stage to perform a left outer join with another collection.
	if mc.lookup != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$lookup", Value: mc.lookup}})
	}

	// Add "$graphLookup" stage for recursive graph-like searches.
	if mc.graphLookup != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$graphLookup", Value: mc.graphLookup}})
	}

	// Return the constructed aggregation pipeline.
	return pipelines
}

// Aggregate performs an aggregation query on the MongoDB collection using the specified aggregation stages.
// It constructs and executes an aggregation pipeline based on the configured stages and options.
//
// Parameters:
// - data: An interface{} to which the result of the aggregation will be decoded.
//
// Returns:
// - err: Any error that occurred during the aggregation or decoding process.
//
// Note:
// This method executes the aggregation pipeline defined by the configured stages, such as '$match', '$group', '$sort', etc.,
// using the MongoDB Go driver. It then decodes the result of the aggregation into the provided 'data' interface{}.
// After execution, it cleans up the aggregation and query options to prepare for future operations.
func (mc *MongoConnection) Aggregate(data interface{}) error {
	// Clean up aggregation and query options after execution.
	defer mc.refreshAggregation()
	defer mc.refreshValues()

	// Execute the aggregation query on the MongoDB collection.
	cursor, err := mc.client.Database(mc.currentDB).Collection(mc.tableName).Aggregate(context.TODO(), mc.createAggregation())
	if err != nil {
		return err
	}

	// Decode the result of the aggregation into the provided 'data' interface{}.
	if err = cursor.All(context.TODO(), data); err != nil {
		return err
	}

	return nil
}

// AggregatePaginate performs an aggregation query on the MongoDB collection with pagination support.
// It constructs and executes an aggregation pipeline with the specified stages and pagination parameters.
//
// Parameters:
// - data: An interface{} to which the result of the aggregation will be decoded.
// - pageNumber: The page number for the pagination. Starts from 1.
//
// Returns:
// - paginationResults: An object containing pagination information (total documents, total pages, limit, and current page).
// - err: Any error that occurred during the aggregation or decoding process.
//
// Note:
// This method first adds a '$match' stage if 'mc.match' is specified, then a '$group' stage if 'mc.group' is specified,
// and finally a '$group' stage to count the total documents matching the aggregation criteria.
// After obtaining the total count, it calculates pagination information and applies the appropriate '$skip' and '$limit' stages
// to retrieve the desired page of data. The result is then decoded into the provided 'data' interface{}.
func (mc *MongoConnection) AggregatePaginate(data interface{}, pageNumber int) (paginationResults connections.PaginationResults, err error) {
	// Clean up aggregation and query options after execution.
	defer mc.refreshAggregation()
	defer mc.refreshValues()

	// Create an aggregation stage to count the total documents.
	countGoupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}

	// Initialize the aggregation pipeline with stages.
	pipelines := primitive.A{}

	// Add a '$match' stage if 'mc.match' is specified.
	if mc.match != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$match", Value: mc.match}})
	}

	// Add a '$group' stage if 'mc.group' is specified.
	if mc.group != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$group", Value: mc.group}})
	}

	// Add the count stage to the pipeline.
	pipelines = append(pipelines, countGoupStage)

	// Perform the aggregation to get the total count of documents.
	countData := []ftypes.QMap{}
	cursor, err := mc.client.Database(mc.currentDB).Collection(mc.tableName).Aggregate(context.TODO(), pipelines)
	if err != nil {
		return
	}
	if err = cursor.All(context.TODO(), &countData); err != nil {
		return
	}

	// Check if countData is not empty (i.e., there are documents to paginate).
	if len(countData) > 0 {
		count := int64(countData[0]["count"].(int32))

		// Calculate pagination information.
		paginationResults.TotalDocuments = count
		paginationResults.TotalPages = int64(math.Ceil((float64(count) / float64(mc.pageSize))))
		paginationResults.Limit = int64(mc.pageSize)
		paginationResults.CurrentPage = int64(pageNumber)

		// Configure aggregation parameters for pagination.
		mc.limitAggregate = mc.pageSize
		mc.skipAggregate = int((pageNumber - 1) * mc.pageSize)

		// Execute the aggregation query with pagination parameters.
		cursor, err = mc.client.Database(mc.currentDB).Collection(mc.tableName).Aggregate(context.TODO(), mc.createAggregation())
		if err != nil {
			return
		}
		if err = cursor.All(context.TODO(), data); err != nil {
			return
		}
	}
	return
}
