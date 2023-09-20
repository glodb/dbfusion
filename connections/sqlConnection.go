package connections

import (
	"github.com/glodb/dbfusion/joins"
)

type SQLConnection interface {
	Connection
	ExecuteSQL(sql string, args ...interface{}) error

	Where(interface{}) SQLConnection
	Table(tableName string) SQLConnection
	GroupBy(fieldname string) SQLConnection
	Having(conditions interface{}) SQLConnection
	Skip(skip int64) SQLConnection
	Limit(limit int64) SQLConnection
	Sort(sortKey string, sortdesc ...bool) SQLConnection
	Project(keys map[string]bool) SQLConnection
	Join(join joins.Join) SQLConnection
}