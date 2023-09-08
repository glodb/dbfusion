package databases

import "github.com/glodb/dbfusion/query"

type SQLQuery interface {
	query.Query
	// New method to create a table.
	Join(joinOperator string, tablename string, condition string, args ...interface{})
	ExecuteSQL(sql string, args ...interface{}) error
}
