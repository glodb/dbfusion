package implementations

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/hooks"
	"github.com/glodb/dbfusion/queryoptions"
	"github.com/glodb/dbfusion/set"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DBCommon struct {
	cache        *caches.Cache
	currentDB    string
	tableName    string
	whereQuery   interface{}
	skip         int64
	limit        int64
	projection   interface{}
	sort         interface{}
	joins        string
	groupBy      string
	havingString string
	havingValues []interface{}
	orderBy      string
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

func (dbc *DBCommon) createTagValueMap(data interface{}) (tagMapValue map[string]interface{}, err error) {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()

	tagMapValue = make(map[string]interface{})
	switch dataType.Kind() {
	case reflect.Ptr:
		ptrValue := reflect.ValueOf(data)
		dataValue = ptrValue.Elem()
		dataType = dataValue.Type()
	}

	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)

		rawtags := strings.Split(field.Tag.Get("dbfusion"), ",")
		tagName := rawtags[0]

		if tagName == "" {
			continue
		}
		value := dataValue.Field(i).Interface()

		tagMapValue[tagName] = value
	}

	return
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
	case reflect.Slice:
		entityData.entityName = dbc.tableName
		entitySet = true
		structType = 2
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
	prefindReturn.dataValue = nameData.dataValue
	prefindReturn.dataType = nameData.dataType

	var dbFusionData conditions.DBFusionData

	if value, ok := dbc.whereQuery.(conditions.DBFusionData); !ok {
		return prefindReturn, dbfusionErrors.ErrInvalidType
	} else {
		dbFusionData = value
	}

	if options.ForceDB { //No need to check cache as it is forced to query db only
		prefindReturn.query = dbFusionData.GetQuery()
		prefindReturn.whereQuery = dbc.whereQuery
		prefindReturn.queryDatabase = true
	} else {
		ok := false
		redisKey := dbc.currentDB + "_" + prefindReturn.entityName + "_" + dbFusionData.GetCacheValues()
		ok, err = caches.GetInstance().ProceessGetCache(*cache, redisKey, result)

		if err != nil {
			return
		}
		skipDB := false

		if !ok { //Data is not found in the redis composite index try to find if this result exists
			redisQueryKey := dbc.currentDB + "_" + prefindReturn.entityName + "_" + dbFusionData.GetCacheKey()
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
		if value, ok := dbc.whereQuery.(conditions.DBFusionData); ok { //Caching is only possible if it is of type dbFusionData
			redisQueryKey := dbc.currentDB + "_" + entityName + "_" + value.GetCacheKey()
			caches.GetInstance().ProceessSetQueryCache(*cache, redisQueryKey, result)
		}
	}

	if value, ok := interface{}(result).(hooks.PostFind); ok {
		dbc.whereQuery = value.PostFind()
	}
	return nil
}

func (dbc *DBCommon) buildMongoUpdate(data interface{}, nameData entityData) (interface{}, error) {
	dataValue := nameData.dataValue
	dataType := nameData.dataType

	structType := nameData.structType
	tagValueMap := make(map[string]interface{})

	var topMap interface{}
	if structType == 1 { //Its a structure
		queryMap := primitive.D{}
		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)

			if structType == 1 {
			}
			rawtags := strings.Split(field.Tag.Get("dbfusion"), ",")
			tagName := rawtags[0]

			if tagName == "" {
				continue
			}
			value := dataValue.Field(i).Interface()

			if !dbc.isFieldSet(dataValue.Field(i)) {
				continue
			}

			singlePoint := primitive.E{Key: tagName, Value: value}
			queryMap = append(queryMap, singlePoint)
			tagValueMap[tagName] = value
		}
		topMap = primitive.D{{Key: "$set", Value: queryMap}}

	} else {
		if value, ok := data.(ftypes.QMap); ok {
			queryMap := primitive.D{}
			for key, val := range value {
				singlePoint := primitive.E{Key: key, Value: val}
				queryMap = append(queryMap, singlePoint)
			}
			topMap = primitive.D{{Key: "$set", Value: queryMap}}
		} else if value, ok := data.(ftypes.DMap); ok {
			queryMap := primitive.D(value)
			topMap = primitive.D{{Key: "$set", Value: queryMap}}
		} else if value, ok := data.(map[string]interface{}); ok {
			queryMap := primitive.D{}
			for key, val := range value {
				singlePoint := primitive.E{Key: key, Value: val}
				queryMap = append(queryMap, singlePoint)
			}
			topMap = primitive.D{{Key: "$set", Value: queryMap}}
		} else {
			return topMap, dbfusionErrors.ErrInvalidType
		}
	}
	return topMap, nil
}

func (dbc *DBCommon) preUpdate(data interface{}, dbType ftypes.DBTypes) (preUpdateData preUpdateReturn, err error) {

	if value, ok := interface{}(data).(hooks.PreUpdate); ok {
		data = value.PreUpdate()
	}

	nameData, nameErr := dbc.getEntityName(data)

	if nameErr != nil {
		err = nameErr
		return
	}
	preUpdateData.entityName = nameData.entityName

	if dbType == connections.MONGO {
		preUpdateData.queryData, err = dbc.buildMongoUpdate(data, nameData)
	} else if dbType == connections.MYSQL {

	}
	return
}

func (dbc *DBCommon) postUpdate(cache *caches.Cache, result interface{}, entityName string, oldValues []string, newValues []string) error {

	caches.GetInstance().ProceessUpdateCache(*cache, oldValues, newValues, result)

	if value, ok := interface{}(result).(hooks.PostUpdate); ok {
		result = value.PostUpdate()
	}

	return nil
}

//Return false if doesn't implement hooks.Cachehook
//If where query is null ask to query the db and update whole cache keys
//Check if the key is found in the db
func (dbc *DBCommon) getAllCacheValues(data hooks.CacheHook, tagValueMap map[string]interface{}, entityName string) []string {

	cacheKeys := make([]string, 0)
	for _, cacheKey := range data.GetCacheIndexes() {
		internalKeys := strings.Split(cacheKey, ",")

		cacheKey := ""
		for _, internalKey := range internalKeys {
			if value, ok := tagValueMap[internalKey]; ok {
				if cacheKey == "" {
					cacheKey += fmt.Sprintf("%v", value)
				} else {
					cacheKey += fmt.Sprintf("_%v", value)
				}
			}
		}
		cacheKey = dbc.currentDB + "_" + entityName + "_" + cacheKey
		cacheKeys = append(cacheKeys, cacheKey)
	}
	return cacheKeys
}

func (dbc *DBCommon) refreshValues() {
	dbc.tableName = ""
	dbc.whereQuery = nil
	dbc.skip = 0
	dbc.limit = 0
	dbc.projection = nil
	dbc.sort = nil
	dbc.joins = ""
	dbc.groupBy = ""
	dbc.havingString = ""
	dbc.havingValues = make([]interface{}, 0)
	dbc.orderBy = ""
}
