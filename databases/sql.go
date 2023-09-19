package databases

import "github.com/glodb/dbfusion/dbconnections"

type SQLQuery interface {
	dbconnections.DBConnections
	ExecuteSQL(sql string, args ...interface{}) error
}
