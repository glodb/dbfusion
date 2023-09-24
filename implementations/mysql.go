package implementations

import (
	"database/sql"
	"fmt"
	"math"

	_ "github.com/go-sql-driver/mysql"

	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/hooks"
	"github.com/glodb/dbfusion/joins"
	"github.com/glodb/dbfusion/queryoptions"
	"github.com/glodb/dbfusion/utils"
)

// MySql represents a type for interacting with a MySQL database. It embeds the SqlBase type to reuse its methods and fields.
type MySql struct {
	SqlBase         // Embedding SqlBase for code reuse.
	db      *sql.DB // db is a reference to a MySQL database connection.
}

func (ms *MySql) ConnectWithCertificate(uri string, filePath string) error {
	return nil
}

// Connect establishes a connection to a MySQL database using the provided URI and sets it in the MySql instance.
// It returns an error if the connection cannot be established.
//
// Parameters:
// - uri (string): The URI to connect to the MySQL database.
//
// Returns:
// - error: An error indicating any issues with the database connection, or nil if the connection is successful.
func (ms *MySql) Connect(uri string) error {
	// Attempt to establish a connection to the MySQL database using the provided URI.
	db, err := sql.Open("mysql", uri)
	if err != nil {
		// If an error occurs during the connection attempt, return the error.
		return err
	}

	// Set the established database connection in the MySql instance.
	ms.db = db

	// Return nil, indicating a successful connection.
	return nil
}

// Table sets the table name for the SQL operation and returns the updated MySql instance.
//
// Parameters:
// - tablename (string): The name of the table to set for SQL operations.
//
// Returns:
// - connections.SQLConnection: The updated MySql instance with the specified table name.
func (ms *MySql) Table(tablename string) connections.SQLConnection {
	// Set the provided table name for SQL operations.
	ms.setTable(tablename)

	// Return the updated MySql instance.
	return ms
}

// InsertOne inserts a single record into the MySQL database table.
//
// Parameters:
// - data (interface{}): The data to insert into the table.
//
// Returns:
// - error: An error if the insertion operation fails, or nil if successful.
func (ms *MySql) InsertOne(data interface{}) error {
	// Ensure that values are reset after the operation.
	defer ms.refreshValues()

	// Create the SQL insert query, values, and preCreateData.
	query, values, preCreateData, err := ms.createSqlInsert(data)

	if err != nil {
		return err
	}

	// Execute the SQL insert query with the provided values.
	_, err = ms.db.Exec(query, values...)

	// Perform post-insert operations if the insert was successful.
	if err == nil {
		err = ms.postInsert(ms.cache, preCreateData.Data, preCreateData.mData, ms.currentDB, preCreateData.entityName)
	}

	// Return any errors encountered during the operation.
	return err
}

// FindOne retrieves a single record from the MySQL database table based on the provided conditions.
//
// Parameters:
// - result (interface{}): A pointer to the result structure where the retrieved data will be stored.
// - dbFusionOptions (...queryoptions.FindOptions): Optional FindOptions to customize the query.
//
// Returns:
// - error: An error if the retrieval operation fails, or nil if successful.
func (ms *MySql) FindOne(result interface{}, dbFusionOptions ...queryoptions.FindOptions) error {
	// Ensure that values are reset after the operation.
	defer ms.refreshValues()

	// Initialize valuesInterface to store query values.
	valuesInterface := make([]interface{}, 0)

	// Check if a WHERE condition is specified.
	if ms.whereQuery != nil {
		// Get the SQL fusion data and update the WHERE condition.
		query, err := utils.GetInstance().GetSqlFusionData(ms.whereQuery)
		if err != nil {
			return err
		}
		ms.whereQuery = query
		valuesInterface = append(valuesInterface, query.GetValues().([]interface{})...)
	} else {
		// If no WHERE condition is provided, create an empty one.
		ms.whereQuery = &conditions.SqlData{}
	}

	// Prepare for the preFind operation.
	prefindReturn, err := ms.preFind(ms.cache, result, dbFusionOptions...)

	if err != nil {
		return err
	}

	// If the value is found in cache, no need to query the database.
	if prefindReturn.queryDatabase {

		// Append any HAVING values to valuesInterface.
		if len(ms.havingValues) != 0 {
			valuesInterface = append(valuesInterface, ms.havingValues...)
		}

		// Create the SQL SELECT query for retrieving one record.
		query := ms.createFindQuery(prefindReturn.entityName, true)

		// Execute the query and retrieve the data.
		rows, err := ms.db.Query(query, valuesInterface...)
		if err != nil {
			return err
		}
		_, err = ms.readSqlDataFromRows(rows, prefindReturn.dataType, prefindReturn.dataValue)
		if err != nil {
			return err
		}
	}

	// Perform post-find operations.
	err = ms.postFind(ms.cache, result, prefindReturn.entityName, dbFusionOptions...)
	return err
}

// UpdateAndFindOne updates a record in the MySQL database table based on the provided conditions and retrieves the updated record.
//
// Parameters:
// - data (interface{}): The data to update the record with.
// - result (interface{}): A pointer to the result structure where the retrieved data will be stored.
// - upsert (bool): Indicates whether to perform an upsert (insert if record not found).
//
// Returns:
// - error: An error if the update operation fails, or nil if successful.
func (ms *MySql) UpdateAndFindOne(data interface{}, result interface{}, upsert bool) error {
	// Ensure that values are reset after the operation.
	defer ms.refreshValues()

	// Initialize valuesInterface to store query values.
	valuesInterface := make([]interface{}, 0)

	// Check if a WHERE condition is specified.
	if ms.whereQuery != nil {
		// Get the SQL fusion data and update the WHERE condition.
		query, err := utils.GetInstance().GetSqlFusionData(ms.whereQuery)
		if err != nil {
			return err
		}
		ms.whereQuery = query
		valuesInterface = append(valuesInterface, query.GetValues().([]interface{})...)
	} else {
		// If no WHERE condition is provided, create an empty one.
		ms.whereQuery = &conditions.SqlData{}
	}

	// Prepare for the preUpdate operation.
	preUpdateReturn, err := ms.preUpdate(result, connections.MYSQL)

	// Create a SQL SELECT query to retrieve the record.
	query := ms.createFindQuery(preUpdateReturn.entityName, true)

	// Execute the query to retrieve the record.
	rows, err := ms.db.Query(query, valuesInterface...)
	if err != nil {
		return err
	}

	// Read the retrieved data into the result structure.
	rowsCount, err := ms.readSqlDataFromRows(rows, preUpdateReturn.dataType, preUpdateReturn.dataValue)
	if err != nil {
		return err
	}

	oldValues := make([]string, 0)
	newValues := make([]string, 0)
	updateCache := false
	var cacheHook hooks.CacheHook

	// Check if the result implements the CacheHook interface.
	if value, ok := interface{}(result).(hooks.CacheHook); ok {
		// Create a map of tag values from the result structure.
		tagMapValue, err := ms.createTagValueMap(result)
		if err == nil {
			// Get the old cache values and prepare to update the cache.
			oldValues = ms.getAllCacheValues(value, tagMapValue, preUpdateReturn.entityName)
			updateCache = true
			cacheHook = value
		}
	}

	// Check if the record is not found, and upsert is enabled.
	if rowsCount == 0 && upsert {
		// Insert the record into the database.
		query, values, _, err := ms.createSqlInsert(data)
		if err != nil {
			return err
		}
		_, err = ms.db.Exec(query, values...)
		if err != nil {
			return err
		}
	} else {
		// Update the record in the database.
		commands, setValues, err := ms.buildMySqlUpdate(data,
			entityData{
				entityName: preUpdateReturn.entityName,
				dataType:   preUpdateReturn.dataType,
				dataValue:  preUpdateReturn.dataValue,
				structType: preUpdateReturn.structType})

		if err != nil {
			return err
		}
		setValues = append(setValues, valuesInterface...)
		query := ms.createUpdateQuery(preUpdateReturn.entityName, commands, true)
		_, err = ms.db.Exec(query, setValues...)

		if err != nil {
			return err
		}
	}

	// Merge the results of the select query and data provided to update the cache values.
	merged := ms.merge(data, result)

	if updateCache {
		// Create a map of tag values from the merged result.
		tagMapValue, _ := ms.createTagValueMap(merged)
		// Get the new cache values and update the cache.
		newValues = ms.getAllCacheValues(cacheHook, tagMapValue, preUpdateReturn.entityName)
	}
	err = ms.postUpdate(ms.cache, result, preUpdateReturn.entityName, oldValues, newValues)

	return nil
}

// DeleteOne deletes a record from the MySQL database table based on the provided conditions.
//
// Parameters:
// - sliceData (interface{}...): Variable number of data slices.
//
// Returns:
// - error: An error if the delete operation fails, or nil if successful.
func (ms *MySql) DeleteOne(sliceData ...interface{}) error {
	// Ensure that values are reset after the operation.
	defer ms.refreshValues()

	var data interface{}
	if len(sliceData) != 0 {
		data = sliceData[0]
	}

	// Check if a WHERE condition is specified.
	if ms.whereQuery != nil {
		// Get the SQL fusion data and update the WHERE condition.
		query, err := utils.GetInstance().GetSqlFusionData(ms.whereQuery)
		if err != nil {
			return err
		}
		ms.whereQuery = query
	} else {
		// If no WHERE condition is provided, create an empty one.
		ms.whereQuery = &conditions.SqlData{}
	}

	// Prepare for the preDelete operation.
	preDeleteData, err := ms.preDelete(data)

	if err != nil {
		return err
	}

	if data != nil { // Need to delete from a struct
		whereConditions, dataInterface, err := ms.buildMySqlDeleteData(preDeleteData.dataType, preDeleteData.dataValue)
		selectQuery := fmt.Sprintf("SELECT * from %s LIMIT 1", preDeleteData.entityName)
		if whereConditions != "" {
			selectQuery = fmt.Sprintf("SELECT * from %s WHERE %s LIMIT 1", preDeleteData.entityName, whereConditions)
		}

		// Execute the SELECT query to check if the record exists.
		rows, err := ms.db.Query(selectQuery, dataInterface...)
		if err != nil {
			return err
		}
		rowsCount, err := ms.readSqlDataFromRows(rows, preDeleteData.dataType, preDeleteData.dataValue)

		if rowsCount == 1 {
			// If the record exists, create a DELETE query and execute it.
			deleteQuery := ms.createDeleteQuery(preDeleteData.entityName, whereConditions, true)
			_, err := ms.db.Query(deleteQuery, dataInterface...)
			if err != nil {
				return err
			}
		}
	} else { // Need to delete based on WHERE conditions
		// Create a DELETE query and execute it.
		deleteQuery := ms.createDeleteQuery(preDeleteData.entityName, "", true)
		_, err := ms.db.Query(deleteQuery, ms.whereQuery.(conditions.DBFusionData).GetValues().([]interface{})...)
		if err != nil {
			return err
		}
	}
	return nil
}

// DisConnect disconnects from the MySQL database.
//
// Returns:
// - error: An error if the disconnection fails, or nil if successful.
func (ms *MySql) DisConnect() error {
	return ms.db.Close()
}

// Paginate fetches a page of results from the MySQL database and populates the provided results interface.
//
// Parameters:
// - results: A pointer to the interface where the query results will be stored.
// - pageNumber: The page number to retrieve (starting from 0).
//
// Returns:
// - paginationResults: A struct containing pagination information.
// - error: An error if the pagination process fails, or nil if successful.
func (ms *MySql) Paginate(results interface{}, pageNumber int) (connections.PaginationResults, error) {
	defer ms.refreshValues()

	// Check if a whereQuery exists and convert it to SQL format if necessary
	if ms.whereQuery != nil {
		query, err := utils.GetInstance().GetSqlFusionData(ms.whereQuery)
		if err != nil {
			return connections.PaginationResults{}, err
		}
		ms.whereQuery = query
	} else {
		ms.whereQuery = &conditions.SqlData{}
	}

	var paginationResults connections.PaginationResults
	countQuery := ms.createCountQuery(ms.tableName)

	var count int64
	row, err := ms.db.Query(countQuery)
	if err != nil {
		return paginationResults, err
	}

	countQueryRows := 0
	for row.Next() {
		countQueryRows++
		err = row.Scan(&count)
	}

	if err != nil {
		return paginationResults, err
	}

	if countQueryRows == 0 {
		return paginationResults, dbfusionErrors.ErrNoRecordFound
	}

	// Calculate pagination information
	paginationResults.TotalDocuments = count
	paginationResults.TotalPages = int64(math.Ceil((float64(count) / float64(ms.pageSize))))
	paginationResults.Limit = int64(ms.pageSize)
	paginationResults.CurrentPage = int64(pageNumber)

	// Set limit and skip for the SQL query
	ms.limit = int64(ms.pageSize)
	ms.skip = int64(pageNumber * ms.pageSize)

	findQuery := ms.createFindQuery(ms.tableName, false)
	rows, err := ms.db.Query(findQuery)

	if err != nil {
		return paginationResults, err
	}

	// Read SQL rows into the provided results interface
	ms.readSqlRowsToArray(rows, results)

	return paginationResults, nil
}

// CreateTable creates a database table based on the provided data structure.
//
// Parameters:
// - data: The data structure that represents the table schema.
// - ifNotExist: A boolean flag indicating whether to create the table only if it doesn't already exist.
//
// Returns:
// - error: An error if the table creation process fails, or nil if successful.
func (ms *MySql) CreateTable(data interface{}, ifNotExist bool) error {
	// Generate the SQL query for creating the table
	query, err := ms.createTableQuery(data, ifNotExist)
	if err != nil {
		return err
	}

	// Execute the SQL query to create the table
	_, err = ms.db.Exec(query)
	return err
}

// New methods for bulk operations.
func (ms *MySql) InsertMany(interface{}) error {
	return nil
}
func (mc *MySql) FindMany(interface{}, ...queryoptions.FindOptions) error {
	return nil
}
func (ms *MySql) UpdateMany(interface{}, interface{}, bool) error {
	return nil
}
func (ms *MySql) DeleteMany(...interface{}) error {
	return nil
}

// Skip sets the number of records to skip when performing a query.
//
// Parameters:
// - skip: The number of records to skip.
//
// Returns:
// - connections.SQLConnection: The MySQL connection instance for method chaining.
func (ms *MySql) Skip(skip int64) connections.SQLConnection {
	ms.skip = skip
	return ms
}

// Limit sets the maximum number of records to retrieve when performing a query.
//
// Parameters:
// - limit: The maximum number of records to retrieve.
//
// Returns:
// - connections.SQLConnection: The MySQL connection instance for method chaining.
func (ms *MySql) Limit(limit int64) connections.SQLConnection {
	ms.limit = limit
	return ms
}

// Select specifies the columns to include in the query result.
//
// Parameters:
// - keys: A map where the keys represent column names, and the values indicate whether to include the column.
//
// Returns:
// - connections.SQLConnection: The MySQL connection instance for method chaining.
func (ms *MySql) Select(keys map[string]bool) connections.SQLConnection {
	selectionKeys := make([]string, 0)

	for key, val := range keys {
		if val {
			selectionKeys = append(selectionKeys, key)
		}
	}
	ms.projection = selectionKeys
	return ms
}

// Sort sets the sorting order for query results based on a specified column.
//
// Parameters:
// - sortKey: The name of the column to sort by.
// - sortdesc: An optional boolean parameter to indicate descending sorting (default is ascending).
//
// Returns:
// - connections.SQLConnection: The MySQL connection instance for method chaining.
func (ms *MySql) Sort(sortKey string, sortdesc ...bool) connections.SQLConnection {
	sortString := sortKey
	sortVal := " ASC"
	if len(sortdesc) > 0 {
		if !sortdesc[0] {
			sortVal = " DESC"
		}
	}
	sortString += sortVal
	if ms.sort != nil {
		sortedValues := ms.sort.(string)
		if sortedValues != "" {
			sortedValues += "," + sortString
		}
		ms.sort = sortedValues
	} else {
		ms.sort = sortString
	}
	return ms
}

// Where specifies the WHERE clause for query filtering.
//
// Parameters:
// - query: The query condition for filtering.
//
// Returns:
// - connections.SQLConnection: The MySQL connection instance for method chaining.
func (ms *MySql) Where(query interface{}) connections.SQLConnection {
	ms.whereQuery = query
	return ms
}

// Join specifies a join operation to be performed in the query.
//
// Parameters:
// - join: The join configuration.
//
// Returns:
// - connections.SQLConnection: The MySQL connection instance for method chaining.
func (ms *MySql) Join(join joins.Join) connections.SQLConnection {
	query := ""
	switch join.Operator {
	case joins.CROSS_JOIN:
		query = "CROSS JOIN "
	case joins.INNER_JOIN:
		query = "INNER JOIN "
	case joins.LEFT_JOIN:
		query = "LEFT JOIN "
	case joins.RIGHT_JOIN:
		query = "RIGHT JOIN "
	}

	query += join.TableName
	query += " ON " + join.Condition

	if ms.joins != "" {
		query += " " + ms.joins
	}
	ms.joins = query
	return ms
}

// GroupBy specifies the GROUP BY clause for grouping query results.
//
// Parameters:
// - fieldName: The name of the column to group by.
//
// Returns:
// - connections.SQLConnection: The MySQL connection instance for method chaining.
func (ms *MySql) GroupBy(fieldName string) connections.SQLConnection {
	if ms.groupBy != "" {
		ms.groupBy += ","
	}
	ms.groupBy += fieldName
	return ms
}

// Having specifies the HAVING clause for filtering grouped query results.
//
// Parameters:
// - data: The HAVING condition for filtering.
//
// Returns:
// - connections.SQLConnection: The MySQL connection instance for method chaining.
func (ms *MySql) Having(data interface{}) connections.SQLConnection {
	dbfusionData, _ := utils.GetInstance().GetSqlFusionData(data)
	if ms.havingString != "" {
		ms.havingString += "," + dbfusionData.GetQuery().(string)
	} else {
		ms.havingString += dbfusionData.GetQuery().(string)
	}
	ms.havingValues = append(ms.havingValues, dbfusionData.GetValues().([]interface{})...)
	return ms
}

// ExecuteSQL allows executing custom SQL queries with optional parameters.
//
// Parameters:
// - sql: The custom SQL query string.
// - args: Optional arguments for query parameters.
//
// Returns:
// - error: An error, if any, encountered during query execution.
func (ms *MySql) ExecuteSQL(sql string, args ...interface{}) error { return nil }

// SetPageSize sets the page size for paginated query results.
//
// Parameters:
// - limit: The page size (number of records per page).
func (ms *MySql) SetPageSize(limit int) {
	ms.pageSize = limit
}
