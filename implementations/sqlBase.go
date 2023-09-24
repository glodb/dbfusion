package implementations

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/utils"
)

type SqlBase struct {
	DBCommon
}

// createSqlInsert generates an SQL INSERT query and associated data for inserting a new record into a SQL database table.
//
// Parameters:
// - data: The input data representing the record to be inserted. It can be a struct or a map.
//
// Returns:
// - query: The SQL INSERT query string, e.g., "INSERT INTO tablename (col1, col2, ...) VALUES (?, ?, ...)".
// - values: A slice of interface{} containing the values to be inserted corresponding to placeholders in the query.
// - preCreateData: A preCreateReturn struct that contains additional information about the data transformation.
// - err: An error, if any, encountered during the process.
//
// The createSqlInsert method performs the following steps:
//
// 1. It prepares the data for insertion by calling preInsert, which handles data transformation and validation.
//    preInsert returns preCreateData, a struct containing keys, placeholders, and values for the insertion.
//
// 2. It checks if keys are available in preCreateData. If keys are not present (struct has no cache keys), it constructs
//    keys and placeholders based on the mData map in preCreateData.
//
// 3. It constructs the SQL INSERT query using the entityName, keys, and placeholders obtained from preCreateData.
//
// Example Usage:
//
//   data := User{Username: "john_doe", Email: "john@example.com"}
//   query, values, preData, err := createSqlInsert(data)
//   // query: "INSERT INTO users (username, email) VALUES (?, ?)"
//   // values: ["john_doe", "john@example.com"]
//   // preData: preCreateReturn{entityName: "users", mData: {"username": "john_doe", "email": "john@example.com"}, ...}
//   // err: nil (no error)
//
// This method allows seamless insertion of data into a SQL table with automatic construction of the INSERT query.
func (sb *SqlBase) createSqlInsert(data interface{}) (query string, values []interface{}, preCreateData preCreateReturn, err error) {
	// Step 1: Prepare data for insertion by calling preInsert to get preCreateData.
	preCreateData, err = sb.preInsert(data)
	if err != nil {
		return "", nil, preCreateData, err
	}

	keys := preCreateData.keys
	placeholders := preCreateData.placeholders
	values = preCreateData.values

	// Step 2: If keys are not present (struct has no cache keys), construct keys and placeholders.
	if len(keys) <= 0 {
		values = make([]interface{}, 0)
		keys = ""
		placeholders = ""
		for key, value := range preCreateData.mData {
			keys += key + ","
			placeholders += "?,"
			values = append(values, value)
		}
		keys = keys[:len(keys)-1]                         // Remove the trailing comma from keys.
		placeholders = placeholders[:len(placeholders)-1] // Remove the trailing comma from placeholders.
	}

	// Step 3: Construct the SQL INSERT query.
	query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", preCreateData.entityName, keys, placeholders)
	return query, values, preCreateData, nil
}

// createQuery generates an SQL query string and a slice of interface{} values from a query represented as an ftypes.DMap.
//
// Parameters:
// - query: An ftypes.DMap representing a query, where keys are SQL query components (e.g., SELECT, FROM, WHERE) and values are the associated data.
//
// Returns:
// - stringQuery: The SQL query string constructed from the query components.
// - data: A slice of interface{} containing the data values associated with placeholders in the SQL query.
//
// The createQuery method takes an ftypes.DMap, which typically represents various components of an SQL query. It extracts the query components
// and associated data values to construct a complete SQL query string with placeholders for values. The method then returns the constructed query
// string and the data values to be used with placeholders.
//
// Example Usage:
//
//   query := ftypes.DMap{
//     "SELECT": "username, email",
//     "FROM":   "users",
//     "WHERE":  "age > ? AND status = ?",
//   }
//   stringQuery, data := createQuery(query)
//   // stringQuery: "SELECT username, email FROM users WHERE age > ? AND status = ?"
//   // data: [30, "active"]
//
// This method allows the dynamic construction of SQL queries from query components and associated data values, making it flexible for query generation.
func (sb *SqlBase) createQuery(query ftypes.DMap) (stringQuery string, data []interface{}) {
	// Initialize variables to construct the SQL query.
	stringQuery = ""
	data = make([]interface{}, 0)

	// Iterate through the query components in the ftypes.DMap.
	for _, val := range query {
		// Append the query component (e.g., "SELECT", "FROM", "WHERE") to the query string.
		stringQuery += val.Key
		// Append the associated data value to the data slice for placeholders in the query.
		data = append(data, val.Value)
	}

	return stringQuery, data
}

// createCountQuery generates an SQL query string for counting rows in a database table specified by the entityName.
//
// Parameters:
// - entityName: A string representing the name of the database table for which row count is to be calculated.
//
// Returns:
// - query: The SQL query string for counting rows in the specified table.
//
// The createCountQuery method is responsible for generating an SQL query string that calculates the total number of rows in a given database table.
// It takes into account optional query components such as joins, WHERE conditions, GROUP BY, and HAVING clauses, if they are set in the SqlBase instance.
//
// Example Usage:
//
//   entityName := "users"
//   countQuery := createCountQuery(entityName)
//   // countQuery: "SELECT COUNT(*) FROM users"
//
// This method simplifies the creation of count queries, allowing flexibility to add conditions and clauses as needed.
func (sb *SqlBase) createCountQuery(entityName string) string {
	// Initialize the query string with the basic COUNT(*) query.
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", entityName)

	// If there are join clauses specified, append them to the query.
	if sb.joins != "" {
		query = fmt.Sprintf(query+" %s", sb.joins)
	}

	// If there is a WHERE condition specified, append it to the query.
	if sb.whereQuery != nil {
		whereData := sb.whereQuery.(*conditions.SqlData)
		if whereData.GetQuery().(string) != "" {
			query = fmt.Sprintf(query+" WHERE %v", whereData.GetQuery())
		}
	}

	// If there is a GROUP BY clause specified, append it to the query, and optionally, add a HAVING clause if set.
	if sb.groupBy != "" {
		query = fmt.Sprintf(query+" GROUP BY %s", sb.groupBy)
		if sb.havingString != "" {
			query = fmt.Sprintf(query+" HAVING %s", sb.havingString)
		}
	}

	return query
}

// createFindQuery generates an SQL query string for retrieving data from a database table specified by the entityName.
//
// Parameters:
// - entityName: A string representing the name of the database table from which data is to be retrieved.
// - limitOne: A boolean indicating whether to limit the query to retrieve only one result (true) or not (false).
//
// Returns:
// - query: The SQL query string for retrieving data from the specified table with optional query components.
//
// The createFindQuery method is responsible for generating an SQL query string for retrieving data from a given database table.
// It takes into account optional query components such as projections, joins, WHERE conditions, GROUP BY, HAVING, sorting, limits, and offsets
// that may be set in the SqlBase instance.
//
// Example Usage:
//
//   entityName := "users"
//   limitOne := false
//   findQuery := createFindQuery(entityName, limitOne)
//   // findQuery: "SELECT * FROM users"
//
// This method offers flexibility in constructing find queries with various optional query components.
func (sb *SqlBase) createFindQuery(entityName string, limitOne bool) string {
	selectionKeys := "*" // Default selection is all columns.
	projections := []string{}

	// If projections are specified, update the selectionKeys accordingly.
	if sb.projection != nil {
		projections = sb.projection.([]string)
		if len(projections) != 0 {
			selectionKeys = strings.Join(projections, ", ")
		}
	}

	// Initialize the query string with the basic SELECT statement and selected columns.
	query := fmt.Sprintf("SELECT %s FROM %s", selectionKeys, entityName)

	// If there are join clauses specified, append them to the query.
	if sb.joins != "" {
		query = fmt.Sprintf(query+" %s", sb.joins)
	}

	// If there is a WHERE condition specified, append it to the query.
	if sb.whereQuery != nil {
		whereData := sb.whereQuery.(*conditions.SqlData)
		if whereData.GetQuery().(string) != "" {
			query = fmt.Sprintf(query+" WHERE %v", whereData.GetQuery())
		}
	}

	// If there is a GROUP BY clause specified, append it to the query, and optionally, add a HAVING clause if set.
	if sb.groupBy != "" {
		query = fmt.Sprintf(query+" GROUP BY %s", sb.groupBy)
		if sb.havingString != "" {
			query = fmt.Sprintf(query+" HAVING %s", sb.havingString)
		}
	}

	// If sorting is specified, append the ORDER BY clause to the query.
	if sb.sort != nil {
		query = fmt.Sprintf(query+" ORDER BY %s", sb.sort.(string))
	}

	// If not limiting to one result, add LIMIT and OFFSET clauses as needed.
	if !limitOne {
		if sb.limit != 0 {
			query = fmt.Sprintf(query+" LIMIT %d", sb.limit)
		}
	} else {
		query = fmt.Sprintf(query+" LIMIT %d", 1)
	}

	// If there is a skip offset specified, append the OFFSET clause to the query.
	if sb.skip != 0 {
		query = fmt.Sprintf(query+" OFFSET %d", sb.skip)
	}

	return query
}

// createUpdateQuery generates an SQL UPDATE query string for modifying data in a database table specified by entityName.
//
// Parameters:
// - entityName: A string representing the name of the database table to update.
// - setCommands: A string containing the SET clause with update commands to modify data.
// - limitOne: A boolean indicating whether to limit the query to update only one row (true) or not (false).
//
// Returns:
// - query: The SQL UPDATE query string for modifying data in the specified table with optional query components.
//
// The createUpdateQuery method constructs an SQL UPDATE query string for modifying data in a specified database table.
// It takes into account optional query components such as joins, WHERE conditions, and limiting the number of rows to update.
// The SET clause containing update commands is provided as setCommands.
//
// Example Usage:
//
//   entityName := "users"
//   setCommands := "SET name = 'John', age = 30"
//   limitOne := false
//   updateQuery := createUpdateQuery(entityName, setCommands, limitOne)
//   // updateQuery: "UPDATE users SET name = 'John', age = 30"
//
// This method offers flexibility in constructing update queries with various optional query components.
func (sb *SqlBase) createUpdateQuery(entityName string, setCommands string, limitOne bool) string {
	// Initialize the query string with the basic UPDATE statement and setCommands.
	query := fmt.Sprintf("UPDATE %s %s", entityName, setCommands)

	// If there are join clauses specified, append them to the query.
	if sb.joins != "" {
		query = fmt.Sprintf(query+" %s", sb.joins)
	}

	// If there is a WHERE condition specified, append it to the query.
	if sb.whereQuery != nil {
		whereData := sb.whereQuery.(*conditions.SqlData)
		if whereData.GetQuery().(string) != "" {
			query = fmt.Sprintf(query+" WHERE %v", whereData.GetQuery())
		}
	}

	// If not limiting to one row, add the LIMIT clause as needed.
	if !limitOne {
		if sb.limit != 0 {
			query = fmt.Sprintf(query+" LIMIT %d", sb.limit)
		}
	} else {
		query = fmt.Sprintf(query+" LIMIT %d", 1)
	}

	return query
}

// createDeleteQuery generates an SQL DELETE query string for deleting data from a database table specified by entityName.
//
// Parameters:
// - entityName: A string representing the name of the database table from which data will be deleted.
// - whereConditions: A string containing the WHERE conditions for the DELETE query (optional).
// - limitOne: A boolean indicating whether to limit the query to delete only one row (true) or not (false).
//
// Returns:
// - query: The SQL DELETE query string for deleting data from the specified table with optional query components.
//
// The createDeleteQuery method constructs an SQL DELETE query string for removing data from a specified database table.
// It allows specifying optional query components such as joins, WHERE conditions, and limiting the number of rows to delete.
// The WHERE conditions can be provided directly as whereConditions, or they can be derived from the stored whereQuery.
//
// Example Usage:
//
//   entityName := "users"
//   whereConditions := "age < 18"
//   limitOne := false
//   deleteQuery := createDeleteQuery(entityName, whereConditions, limitOne)
//   // deleteQuery: "DELETE FROM users WHERE age < 18"
//
// This method offers flexibility in constructing delete queries with various optional query components.
func (sb *SqlBase) createDeleteQuery(entityName string, whereConditions string, limitOne bool) string {
	// Initialize the query string with the basic DELETE FROM statement for the specified table.
	query := fmt.Sprintf("DELETE FROM %s", entityName)

	// If there are join clauses specified, append them to the query.
	if sb.joins != "" {
		query = fmt.Sprintf(query+" %s", sb.joins)
	}

	// Check if specific WHERE conditions are provided, and append them to the query if available.
	if whereConditions != "" {
		query = fmt.Sprintf(query+" WHERE %s", whereConditions)
	} else if sb.whereQuery != nil {
		whereData := sb.whereQuery.(*conditions.SqlData)
		if whereData.GetQuery().(string) != "" {
			query = fmt.Sprintf(query+" WHERE %v", whereData.GetQuery())
		}
	}

	// If not limiting to one row, add the LIMIT clause as needed.
	if !limitOne {
		if sb.limit != 0 {
			query = fmt.Sprintf(query+" LIMIT %d", sb.limit)
		}
	} else {
		query = fmt.Sprintf(query+" LIMIT %d", 1)
	}

	return query
}

// readSqlDataFromRows reads data from the given SQL rows and populates a struct or slice of structs based on the provided data type.
//
// Parameters:
// - rows: A pointer to an SQL Rows result containing the retrieved data.
// - dataType: The reflect.Type representing the data structure into which the data should be mapped.
// - dataValue: The reflect.Value representing the value to which the retrieved data should be assigned.
//
// Returns:
// - rowsCount: The number of rows successfully processed and mapped to the provided data structure.
// - error: An error, if any, that occurred during the data mapping process.
//
// The readSqlDataFromRows method is used to read data from an SQL query result (rows) and populate a struct or slice of structs
// based on the provided data type. It dynamically assigns values to the fields of the data structure(s) based on the column names
// in the SQL result set. This method is typically used for mapping database query results to Go data structures.
//
// Example Usage:
//
//   rows := // SQL query result rows
//   dataType := reflect.TypeOf(User{})
//   dataValue := reflect.New(dataType).Elem()
//   rowsCount, err := readSqlDataFromRows(rows, dataType, dataValue)
//
// This method iterates over the rows of an SQL result set, scans the row data into columnData, and assigns it to the corresponding
// fields in the provided data structure. It returns the number of rows successfully processed and any error encountered during
// the process.
func (sb *SqlBase) readSqlDataFromRows(rows *sql.Rows, dataType reflect.Type, dataValue reflect.Value) (int, error) {
	// Initialize the rowsCount to track the number of successfully processed rows.
	rowsCount := 0

	// Check if the rows pointer is nil, indicating no records found.
	if rows == nil {
		return 0, dbfusionErrors.ErrSQLQueryNoRecordFound
	}
	defer rows.Close() // Ensure the rows are closed when done.

	// Get the column names from the SQL result set.
	columnNames, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	// Create a slice of interface{} to hold the column data.
	columnData := make([]interface{}, len(columnNames))
	for i := range columnData {
		var v interface{}
		columnData[i] = &v
	}

	// Create a map to associate tag names with struct fields for efficient assignment.
	tagField := make(map[string]reflect.Value)
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)

		rawtags := strings.Split(field.Tag.Get("dbfusion"), ",")
		tagName := rawtags[0]
		tagField[tagName] = dataValue.Field(i)
	}

	// Iterate through the rows of the SQL result set.
	for rows.Next() {
		rowsCount++
		// Scan the row data into columnData.
		err := rows.Scan(columnData...)
		if err != nil {
			return 0, err
		}

		// Map the scanned column data to the corresponding struct fields based on column names.
		for idx, name := range columnNames {
			if field, ok := tagField[name]; ok {
				utils.GetInstance().AssignData(columnData[idx], field)
			}
		}
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return 0, err
	}

	// Return the number of successfully processed rows and any encountered error.
	return rowsCount, nil
}

// readSqlRowsToArray reads data from the given SQL rows and populates a slice of structs based on the provided results type.
//
// Parameters:
// - rows: A pointer to an SQL Rows result containing the retrieved data.
// - results: A pointer to a slice of structs where the retrieved data should be stored.
//
// Returns:
// - error: An error, if any, that occurred during the data retrieval and population process.
//
// The readSqlRowsToArray method is used to read data from an SQL query result (rows) and populate a slice of structs
// based on the provided results type. It dynamically assigns values to the fields of the struct elements in the slice
// based on the column names in the SQL result set. This method is typically used for mapping database query results to
// a slice of Go data structures.
//
// Example Usage:
//
//   rows := // SQL query result rows
//   var users []User
//   err := readSqlRowsToArray(rows, &users)
//
// This method iterates over the rows of an SQL result set, scans the row data into columnData, and assigns it to the
// corresponding fields in the struct elements of the results slice. It sets the populated results slice to the provided
// results pointer.
func (sb *SqlBase) readSqlRowsToArray(rows *sql.Rows, results interface{}) error {
	// Create a new slice of the same type as results (e.g., &[]Users{})
	resultSliceType := reflect.TypeOf(results).Elem()
	newSlice := reflect.New(resultSliceType).Elem()

	// Get the field names from struct tags.
	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}

	// Create a slice of interface{} to hold the column data.
	columnData := make([]interface{}, len(columnNames))
	for i := range columnData {
		var v interface{}
		columnData[i] = &v
	}

	// Iterate through rows and populate the newSlice.
	for rows.Next() {
		// Create a new element of the slice's element type.
		elementType := resultSliceType.Elem()
		newElement := reflect.New(elementType).Elem()

		// Scan the row into the fields of the newElement.
		if err := rows.Scan(columnData...); err != nil {
			return err
		}

		// Get field names from struct tags and assign data to struct fields.
		tagField := sb.getFieldNames(elementType, newElement)
		for idx, name := range columnNames {
			if fieldName, ok := tagField[name]; ok {
				utils.GetInstance().AssignData(columnData[idx], newElement.FieldByName(fieldName))
			}
		}

		// Append the newElement to the newSlice.
		newSlice = reflect.Append(newSlice, newElement)
	}

	// Set the populated newSlice to the results pointer.
	reflect.ValueOf(results).Elem().Set(newSlice)

	return nil
}

// getFieldNames retrieves field names from struct tags for the given struct type and maps them to their corresponding
// tag names. It is used to associate column names from SQL query results with struct field names.
//
// Parameters:
// - structType: The reflect.Type of the struct for which field names should be retrieved.
// - dataValue: A reflect.Value of the struct instance that corresponds to structType.
//
// Returns:
// - tagField: A map where the keys are tag names and the values are corresponding struct field names.
//
// The getFieldNames method is responsible for extracting field names from struct tags and creating a mapping between
// the tag names used in SQL queries and the actual struct field names. This mapping is crucial for correctly assigning
// query results to struct fields during data retrieval.
//
// Example Usage:
//
//   type User struct {
//       ID   int    `dbfusion:"user_id"`
//       Name string `dbfusion:"user_name"`
//   }
//
//   var u User
//   tagField := getFieldNames(reflect.TypeOf(u), reflect.ValueOf(u))
//   // Result: tagField map with {"user_id": "ID", "user_name": "Name"}
//
// In this example, the getFieldNames method takes the reflect.Type of a User struct and a reflect.Value of an instance
// of that struct. It extracts the tag names ("user_id" and "user_name") and associates them with their respective
// struct field names ("ID" and "Name") in the returned tagField map.
func (sb *SqlBase) getFieldNames(structType reflect.Type, dataValue reflect.Value) map[string]string {
	tagField := make(map[string]string)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		rawtags := strings.Split(field.Tag.Get("dbfusion"), ",")
		tagName := rawtags[0]
		if tagName != "" {
			tagField[tagName] = field.Name
		}
	}
	return tagField
}

// createTableQuery generates a SQL query to create a database table based on the structure of a provided data interface.
// Parameters:
// - data: The data interface for which the table schema should be created. The structure of this data is used to determine
//         the table's columns and data types.
// - ifNotExist: A boolean indicating whether the table should only be created if it does not already exist in the database.
// Returns:
// - query: A string representing the SQL query to create the table.
// - error: An error if any issues occur during query generation.
func (sb *SqlBase) createTableQuery(data interface{}, ifNotExist bool) (string, error) {
	query := ""

	// Get the entity name and data type from the provided data interface.
	name, err := sb.getEntityName(data)
	if err != nil {
		return "", err
	}

	// Determine whether to include "IF NOT EXISTS" in the CREATE TABLE query.
	if ifNotExist {
		query = `CREATE TABLE IF NOT EXISTS ` + name.entityName + ` (`
	} else {
		query = `CREATE TABLE ` + name.entityName + ` (`
	}

	dataType := name.dataType

	// Check if the data type is a struct, as we can only generate table schemas from structs.
	if dataType.Kind() != reflect.Struct {
		return "", dbfusionErrors.ErrInvalidType
	}

	columns := ""

	// Iterate over the fields of the struct to generate column definitions.
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		tags := strings.Split(field.Tag.Get("dbfusion"), ",")

		if columns != "" {
			columns += ","
		}

		// Join struct tags to form the column definition (e.g., "column_name PRIMARY KEY,other_column").
		columns += strings.Join(tags, " ")
	}

	query += columns + ");"

	return query, nil
}
