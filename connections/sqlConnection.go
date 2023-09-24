package connections

import (
	"github.com/glodb/dbfusion/joins"
)

// SQLConnection is an interface that extends the base Connection interface and provides
// methods specific to SQL database interactions. It allows executing SQL queries, defining
// table structures, specifying query criteria, sorting, limiting, and more.
type SQLConnection interface {
	Connection
	// ExecuteSQL executes a SQL query with optional arguments.
	// It takes a SQL query string and a variadic list of arguments and returns an error if the query execution fails.
	ExecuteSQL(sql string, args ...interface{}) error

	// CreateTable creates a database table based on the provided tableType.
	// It takes a tableType interface representing the table structure and a boolean indicating whether to create the table if it doesn't exist.
	// It returns an error if the table creation process encounters any issues.
	CreateTable(tableType interface{}, ifNotExist bool) error

	// Where specifies the criteria for filtering records in the SQL database.
	// It takes an interface representing the filter criteria and returns the modified SQLConnection.
	Where(interface{}) SQLConnection

	// Table specifies the name of the SQL database table to query.
	// It takes the name of the table as a parameter and returns the modified SQLConnection.
	Table(tableName string) SQLConnection

	// GroupBy specifies a field for grouping records in the SQL database.
	// It takes the field name as a parameter and returns the modified SQLConnection.
	GroupBy(fieldname string) SQLConnection

	// Having specifies conditions for filtering grouped records in the SQL database.
	// It takes an interface representing the conditions and returns the modified SQLConnection.
	Having(conditions interface{}) SQLConnection

	// Skip specifies the number of records to skip in the result set.
	// It takes an integer representing the number of records to skip and returns the modified SQLConnection.
	Skip(skip int64) SQLConnection

	// Limit specifies the maximum number of records to return in the result set.
	// It takes an integer representing the limit and returns the modified SQLConnection.
	Limit(limit int64) SQLConnection

	// Sort specifies the sorting order of the result set based on a specified key.
	// It takes the key for sorting and an optional boolean indicating descending order.
	// It returns the modified SQLConnection.
	Sort(sortKey string, sortdesc ...bool) SQLConnection

	// Select specifies the fields to include or exclude in the result set.
	// It takes a map where keys represent field names and values represent inclusion/exclusion flags.
	// It returns the modified SQLConnection.
	Select(keys map[string]bool) SQLConnection

	// Join specifies a join operation to combine records from multiple tables in the SQL database.
	// It takes a joins.Join object representing the join operation and returns the modified SQLConnection.
	Join(join joins.Join) SQLConnection
}
