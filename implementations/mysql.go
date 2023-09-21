package implementations

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/hooks"
	"github.com/glodb/dbfusion/joins"
	"github.com/glodb/dbfusion/queryoptions"
	"github.com/glodb/dbfusion/utils"
)

type MySql struct {
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

func (ms *MySql) Table(tablename string) connections.SQLConnection {
	ms.setTable(tablename)
	return ms
}

func (ms *MySql) InsertOne(data interface{}) error {
	defer ms.refreshValues()
	query, values, preCreateData, err := ms.createSqlInsert(data)

	if err != nil {
		return err
	}
	_, err = ms.db.Exec(query, values...)

	if err == nil {
		err = ms.postInsert(ms.cache, preCreateData.Data, preCreateData.mData, ms.currentDB, preCreateData.entityName)
	}

	return err
}

func (ms *MySql) FindOne(result interface{}, dbFusionOptions ...queryoptions.FindOptions) error {
	defer ms.refreshValues()

	valuesInterface := make([]interface{}, 0)
	if ms.whereQuery != nil {
		query, err := utils.GetInstance().GetSqlFusionData(ms.whereQuery)
		if err != nil {
			return err
		}
		ms.whereQuery = query
		valuesInterface = append(valuesInterface, query.GetValues().([]interface{})...)
	} else {
		ms.whereQuery = &conditions.SqlData{}
	}

	prefindReturn, err := ms.preFind(ms.cache, result, dbFusionOptions...)

	if err != nil {
		return err
	}

	//If the value is found in cache don't need to query database
	if prefindReturn.queryDatabase {

		if len(ms.havingValues) != 0 {
			valuesInterface = append(valuesInterface, ms.havingValues...)
		}

		query := ms.createFindQuery(prefindReturn.entityName, true)

		rows, err := ms.db.Query(query, valuesInterface...)
		if err != nil {
			return err
		}
		_, err = ms.readSqlDataFromRows(rows, prefindReturn.dataType, prefindReturn.dataValue)
		if err != nil {
			return err
		}
	}

	err = ms.postFind(ms.cache, result, prefindReturn.entityName, dbFusionOptions...)
	return err
}

func (ms *MySql) UpdateAndFindOne(data interface{}, result interface{}, upsert bool) error {
	defer ms.refreshValues()
	valuesInterface := make([]interface{}, 0)
	if ms.whereQuery != nil {
		query, err := utils.GetInstance().GetSqlFusionData(ms.whereQuery)
		if err != nil {
			return err
		}
		ms.whereQuery = query
		valuesInterface = append(valuesInterface, query.GetValues().([]interface{})...)
	} else {
		ms.whereQuery = &conditions.SqlData{}
	}

	preUpdateReturn, err := ms.preUpdate(result, connections.MYSQL)
	query := ms.createFindQuery(preUpdateReturn.entityName, true)

	rows, err := ms.db.Query(query, valuesInterface...)
	if err != nil {
		return err
	}
	rowsCount, err := ms.readSqlDataFromRows(rows, preUpdateReturn.dataType, preUpdateReturn.dataValue)
	if err != nil {
		return err
	}

	oldValues := make([]string, 0)
	newValues := make([]string, 0)
	updateCache := false
	var cacheHook hooks.CacheHook

	if value, ok := interface{}(result).(hooks.CacheHook); ok {
		tagMapValue, err := ms.createTagValueMap(result)
		if err == nil {
			oldValues = ms.getAllCacheValues(value, tagMapValue, preUpdateReturn.entityName)
			updateCache = true
			cacheHook = value
		}
	}

	if rowsCount == 0 && upsert { //Insert the record in the db
		query, values, _, err := ms.createSqlInsert(data)
		if err != nil {
			return err
		}
		_, err = ms.db.Exec(query, values...)
		if err != nil {
			return err
		}
	} else {
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

	//Merging the results of select query and data provided to update the cache values
	merged := ms.merge(data, result)

	if updateCache {
		tagMapValue, _ := ms.createTagValueMap(merged)
		newValues = ms.getAllCacheValues(cacheHook, tagMapValue, preUpdateReturn.entityName)
	}
	err = ms.postUpdate(ms.cache, result, preUpdateReturn.entityName, oldValues, newValues)

	return nil
}

func (ms *MySql) DeleteOne(sliceData ...interface{}) error {
	defer ms.refreshValues()
	var data interface{}
	if len(sliceData) != 0 {
		data = sliceData[0]
	}

	if data != nil { //Need to delete from a struct

	} else { //need to delete from whereClause

	}
	return nil
}

func (ms *MySql) DisConnect() error {
	return ms.db.Close()
}

func (ms *MySql) Paginate(interface{}, ...queryoptions.FindOptions) {

}

// New method to create a table.
func (ms *MySql) CreateTable(ifNotExist bool) {

}

func (ms *MySql) RegisterSchema() {}

// New methods for bulk operations.
func (ms *MySql) CreateMany([]interface{}) {

}
func (mc *MySql) FindMany(interface{}) {
}
func (ms *MySql) UpdateMany([]interface{}) {

}
func (ms *MySql) DeleteMany(qmap ftypes.QMap) {

}

func (ms *MySql) Skip(skip int64) connections.SQLConnection {
	ms.skip = skip
	return ms
}
func (ms *MySql) Limit(limit int64) connections.SQLConnection {
	ms.limit = limit
	return ms
}
func (ms *MySql) Project(keys map[string]bool) connections.SQLConnection {
	selectionKeys := make([]string, 0)

	for key, val := range keys {
		if val {
			selectionKeys = append(selectionKeys, key)
		}
	}
	ms.projection = selectionKeys
	return ms
}

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

func (ms *MySql) Where(query interface{}) connections.SQLConnection {
	ms.whereQuery = query
	return ms
}

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
func (ms *MySql) GroupBy(fieldName string) connections.SQLConnection {
	if ms.groupBy != "" {
		ms.groupBy += ","
	}
	ms.groupBy += fieldName
	return ms
}

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
func (ms *MySql) ExecuteSQL(sql string, args ...interface{}) error { return nil }
