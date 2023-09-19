package dbconnections

import (
	"reflect"
	"strings"

	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/hooks"
	"github.com/glodb/dbfusion/queryoptions"
	"github.com/glodb/dbfusion/set"
)

type DBCommon struct {
	cache      *caches.Cache
	currentDB  string
	tableName  string
	whereQuery interface{}
	skip       int64
	limit      int64
	projection interface{}
	sort       interface{}
	joins      string
}

func (dbc *DBCommon) SetCache(cache *caches.Cache) {
	dbc.cache = cache
}

func (dbc *DBCommon) ChangeDatabase(dbName string) error {
	dbc.currentDB = dbName
	return nil
}

func (dbc *DBCommon) setTable(tableName string) {
	dbc.tableName = tableName
}

func (dbc *DBCommon) isFieldSet(val reflect.Value) bool {
	if val.IsValid() {
		return !reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
	}
	return false
}

func (dbc *DBCommon) checkPtr(data interface{}) (reflect.Value, reflect.Type) {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()

	switch dataType.Kind() {
	case reflect.Ptr:
		ptrValue := reflect.ValueOf(data)
		dataValue = ptrValue.Elem()
		dataType = dataValue.Type()
	}
	return dataValue, dataType
}

func (dbc *DBCommon) getEntityName(data interface{}) (entityData entityData, err error) {
	dataValue, dataType := dbc.checkPtr(data)

	entityData.dataValue = dataValue
	entityData.dataType = dataType

	structType := 0
	entitySet := false

	switch dataType.Kind() {
	case reflect.Struct:
		structType = 1
	case reflect.Map:
		// Additional check to ensure it's a map[string]interface{}
		if dataType.Key().Kind() == reflect.String && dataType.Elem().Kind() == reflect.Interface {
			if dbc.tableName == "" {
				err = dbfusionErrors.ErrEntityNameRequired
				return
			}
			entityData.entityName = dbc.tableName
			entitySet = true
			structType = 2
		} else {
			err = dbfusionErrors.ErrStringMapRequired
			return
		}
	default:
		err = dbfusionErrors.ErrStringMapRequired
		return
	}

	if !entitySet {
		if value, ok := interface{}(data).(hooks.Entity); ok {
			entityData.entityName = value.GetEntityName()
		} else {
			entityData.entityName = dataType.Name()
		}
	}

	entityData.structType = structType
	return
}

func (dbc *DBCommon) preInsert(data interface{}) (preCreateData preCreateReturn, err error) {

	nameData, nameErr := dbc.getEntityName(data)

	if nameErr != nil {
		err = nameErr
		return
	}

	dataValue := nameData.dataValue
	dataType := nameData.dataType

	structType := nameData.structType
	preCreateData.entityName = nameData.entityName

	if value, ok := interface{}(data).(hooks.PreInsert); ok {
		data = value.PreInsert()
		dataValue = reflect.ValueOf(data)
		dataType = dataValue.Type()
	}

	keys := ""
	placeholders := ""
	values := make([]interface{}, 0)
	if structType == 1 { //This is typ of struct and can have cache keys implemented
		mData := make(map[string]interface{})
		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)

			if structType == 1 {
			}
			rawtags := strings.Split(field.Tag.Get("dbfusion"), ",")
			tagName := rawtags[0]

			tags := set.ConvertArray[string](rawtags)

			if tagName == "" {
				continue
			}
			value := dataValue.Field(i).Interface()

			if tags.Contains("omitempty") {
				if !dbc.isFieldSet(dataValue.Field(i)) {
					continue
				}
			}

			mData[tagName] = value
			keys += tagName + ","
			placeholders += "?,"
			values = append(values, value)
		}
		preCreateData.mData = mData
		if len(keys) > 0 {
			preCreateData.keys = keys[:len(keys)-1]
			preCreateData.placeholders = placeholders[:len(placeholders)-1]
			preCreateData.values = values
		}
	} else {
		//This has been checked in this function for type
		preCreateData.mData = data.(map[string]interface{})
	}

	preCreateData.Data = data
	return
}

func (dbc *DBCommon) postInsert(cache *caches.Cache, data interface{}, mData map[string]interface{}, dbName string, entityName string) error {
	if val, ok := interface{}(data).(hooks.CacheHook); ok {
		if cache != nil {
			err := caches.GetInstance().ProcessInsertCache(*cache, val.GetCacheIndexes(), mData, dbName, entityName)
			if err != nil {
				return err
			}
		} else {
			return dbfusionErrors.ErrNoValidCacheFound
		}
	}

	if value, ok := interface{}(data).(hooks.PostInsert); ok {
		value = value.PostInsert()
	}
	return nil
}

//----Implement a same interface to use common class
//----Check for pre find hooks
//----Instead of forceDB use cache options to check
//----1--ForceDB
//----2--QueryResultsSaveCache
//----Check if the passed query is type dbFusionData
//----if forceDB is checked read from db without any checking
//----***If force db is not checked
//----Check if passed query is already saved in the cache
//----Check only on equal operator
//----If anyother operator included skip checking part
//----if data is not saved in cache or query is not eligible
//----Check if this query result is already in the cache

//TODO: Test mysql for find
//Make skip, sort, limit, and projection to work with find
//Get the values from this class and pass to front class
func (dbc *DBCommon) preFind(cache *caches.Cache, result interface{}, dbFusionOptions ...queryoptions.FindOptions) (prefindReturn preFindReturn, err error) {
	var nameData entityData
	var options queryoptions.FindOptions
	if dbFusionOptions != nil {
		options = dbFusionOptions[0]
	}
	if value, ok := interface{}(result).(hooks.PreFind); ok {
		result = value.PreFind()
	}

	nameData, err = dbc.getEntityName(result)

	if err != nil {
		return
	}

	prefindReturn.entityName = nameData.entityName

	var dbFusionData conditions.DBFusionData
	if value, ok := dbc.whereQuery.(conditions.DBFusionData); !ok {
		prefindReturn.query = dbc.whereQuery
		prefindReturn.whereQuery = dbc.whereQuery
		prefindReturn.queryDatabase = true
		return
	} else {
		dbFusionData = value
	}

	if options.ForceDB { //No need to check cache as it is forced to query db only
		prefindReturn.query = dbFusionData.GetQuery()
		prefindReturn.whereQuery = dbc.whereQuery
		prefindReturn.queryDatabase = true
	} else {
		ok := false
		redisKey := dbc.currentDB + "_" + prefindReturn.entityName + dbFusionData.GetCacheValues()
		if dbFusionData.ShouldQueryDefaultCache() {
			ok, err = caches.GetInstance().ProceessGetCache(*cache, redisKey, result)

			if err != nil {
				return
			}
		}
		skipDB := false

		if !ok { //Data is not found in the redis composite index try to find if this result exists
			redisQueryKey := dbc.currentDB + "_" + prefindReturn.entityName + dbFusionData.GetCacheKey()
			skipDB, err = caches.GetInstance().ProceessGetQueryCache(*cache, redisQueryKey, result)
			if err != nil {
				return
			}

			if !skipDB { //Looked everywhere in cache or not valid to look cache for query now have to query database
				prefindReturn.query = dbFusionData.GetQuery()
				prefindReturn.whereQuery = dbc.whereQuery
				prefindReturn.queryDatabase = true
				return
			}

		}
		prefindReturn.queryDatabase = false
	}
	return
}

func (dbc *DBCommon) postFind(cache *caches.Cache, result interface{}, entityName string, dbFusionOptions ...queryoptions.FindOptions) error {

	var options queryoptions.FindOptions
	if dbFusionOptions != nil {
		options = dbFusionOptions[0]
	}

	if options.CacheResult { //User Asked us to cache the results of this query
		if value, ok := dbc.whereQuery.(conditions.DBFusionData); ok { //Caching is only possible if it is of type DBFusionData
			redisQueryKey := dbc.currentDB + "_" + entityName + value.GetCacheKey()
			caches.GetInstance().ProceessSetQueryCache(*cache, redisQueryKey, result)
		}
	}

	if value, ok := interface{}(result).(hooks.PostFind); ok {
		dbc.whereQuery = value.PostFind()
	}
	return nil
}
