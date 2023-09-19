package dbconnections

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/joins"
	"github.com/glodb/dbfusion/query"
	"github.com/glodb/dbfusion/queryoptions"
)

type MySql struct {
	conditions.SqlCondition
	SqlBase
	db *sql.DB
}

//TODO: add the communication with certificate
//TODO: add the options for cache enabling and disabling
//TODO: add thhe
func (ms *MySql) ConnectWithCertificate(uri string, filePath string) error {
	return nil
}

func (ms *MySql) Connect(uri string) error {
	db, err := sql.Open("mysql", uri)
	if err != nil {
		return err
	}
	ms.db = db
	return nil
}

func (ms *MySql) Table(tablename string) DBConnections {
	ms.setTable(tablename)
	return ms
}

func (ms *MySql) InsertOne(data interface{}) error {

	query, values, preCreateData, err := ms.createSqlInsert(data)
	_, err = ms.db.Exec(query, values...)

	if err == nil {
		err = ms.postInsert(ms.cache, preCreateData.Data, preCreateData.mData, ms.currentDB, preCreateData.entityName)
	}
	return err
}

func (ms *MySql) FindOne(result interface{}, dbFusionOptions ...queryoptions.FindOptions) error {
	prefindReturn, err := ms.preFind(ms.cache, result, dbFusionOptions...)
	if err != nil {
		return err
	}
	stringQuery := ""
	valuesInterface := make([]interface{}, 0)
	if prefindReturn.queryDatabase {
		if value, ok := interface{}(prefindReturn.whereQuery).(conditions.DBFusionData); ok { //it is DBFusion Data so query and values are already broken
			stringQuery = value.GetQuery().(string)
			valuesInterface = value.GetValues().([]interface{})
		} else if value, ok := interface{}(prefindReturn.whereQuery).(query.DMap); ok { //it is DMap build values
			stringQuery, valuesInterface = ms.createQuery(value)
		} else {
			return dbfusionErrors.ErrSQLQueryTypeNotSupported
		}

		log.Println(stringQuery)
		log.Println(valuesInterface...)

		ms.createFindQuery(prefindReturn.entityName, stringQuery)

		if err != nil {
			return err
		}
	}

	err = ms.postFind(ms.cache, result, prefindReturn.entityName, dbFusionOptions...)
	return err
}

func (ms *MySql) UpdateOne(interface{}) error {
	return nil
}

func (ms *MySql) DeleteOne(interface{}) error {
	return nil
}

func (ms *MySql) DisConnect() {
}

func (ms *MySql) Paginate(qmap query.QMap) {

}
func (ms *MySql) Distinct(field string) {
	ms.db.Close()
}

// New method to create a table.
func (ms *MySql) CreateTable(ifNotExist bool) {

}

func (ms *MySql) RegisterSchema() {}

// New methods for grouping and ordering.
func (ms *MySql) GroupBy(keys string)                            {}
func (ms *MySql) OrderBy(order interface{}, args ...interface{}) {}

// New methods for bulk operations.
func (ms *MySql) CreateMany([]interface{}) {

}
func (ms *MySql) UpdateMany([]interface{}) {

}
func (ms *MySql) DeleteMany(qmap query.QMap) {

}

func (ms *MySql) Skip(skip int64) query.Query {
	ms.skip = skip
	return ms
}
func (ms *MySql) Limit(limit int64) query.Query {
	ms.limit = limit
	return ms
}
func (ms *MySql) Project(keys map[string]bool) query.Query {
	selectionKeys := make([]string, 0)

	for key, val := range keys {
		if val {
			selectionKeys = append(selectionKeys, key)
		}
	}
	ms.projection = selectionKeys
	return ms
}

func (ms *MySql) Sort(sortKeys map[string]bool) query.Query {
	ms.sort = sortKeys
	return ms
}

func (ms *MySql) Where(query interface{}) query.Query {
	ms.whereQuery = query
	return ms
}

func (ms *MySql) Join(join joins.Join) query.Query {
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

	if ms.joins == "" {
		query = ms.joins
	} else {
		query += " " + ms.joins
	}
	log.Println(query)

	return ms
}
func (ms *MySql) ExecuteSQL(sql string, args ...interface{}) error { return nil }
