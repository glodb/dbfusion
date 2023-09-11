package dbconnections

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/glodb/dbfusion/query"
)

type MySql struct {
	DBCommon
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

func (ms *MySql) Insert(data interface{}) error {
	preinsertdata, err := ms.preInsert(data)
	if err != nil {
		return err
	}

	keys := preinsertdata.keys
	placeholders := preinsertdata.placeholders
	values := preinsertdata.values
	if len(keys) <= 0 {
		values = make([]interface{}, 0)
		keys = ""
		placeholders = ""
		for key, value := range preinsertdata.mData {
			keys += key + ","
			placeholders += "?,"
			values = append(values, value)
		}
		keys = keys[:len(keys)-1]
		placeholders = placeholders[:len(placeholders)-1]
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", preinsertdata.EntityName, keys, placeholders)
	_, err = ms.db.Exec(query, values...)

	if err == nil {
		err = ms.postInsert(ms.cache, preinsertdata.Data, preinsertdata.mData, ms.currentDB, preinsertdata.EntityName)
	}
	return err
}

func (ms *MySql) Find(interface{}) error {
	return nil
}

func (ms *MySql) Update(interface{}) error {
	return nil
}

func (ms *MySql) Delete(interface{}) error {
	return nil
}

func (ms *MySql) DisConnect() {
}

func (ms *MySql) Filter(qmap query.QMap) {

}

func (ms *MySql) Sort(order interface{}, args ...interface{}) {

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

// New method for specifying query conditions.
func (ms *MySql) Where(condition string, args ...interface{}) {}

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

func (ms *MySql) Skip(skip int64)   {}
func (ms *MySql) Limit(limit int64) {}
