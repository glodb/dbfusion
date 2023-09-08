package dbconnections

import (
	"log"
	"reflect"
	"strings"

	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/hooks"
)

type DBCommon struct {
	cache     *caches.Cache
	currentDB string
}

func (dbc *DBCommon) SetCache(cache *caches.Cache) {
	dbc.cache = cache
}

func (dbc *DBCommon) ChangeDatabase(dbName string) error {
	dbc.currentDB = dbName
	return nil
}

func (dbc *DBCommon) PreInsert(data interface{}) (precreateData PreCreateReturn, err error) {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()

	structType := 0
	switch dataType.Kind() {
	case reflect.Struct:
		structType = 1
	case reflect.Map:
		// Additional check to ensure it's a map[string]interface{}
		if dataType.Key() == reflect.TypeOf("") && dataType.Elem() == reflect.TypeOf(interface{}(nil)) {
			structType = 2
		} else {
			err = dbfusionErrors.ErrStructOrMapRequired
			return
		}
	default:
		err = dbfusionErrors.ErrStructOrMapRequired
		return
	}

	if value, ok := interface{}(data).(hooks.PreInsert); ok {
		data = value.PreInsert()
		dataValue = reflect.ValueOf(data)
		dataType = dataValue.Type()
	}

	if value, ok := interface{}(data).(hooks.Entity); ok {
		precreateData.EntityName = value.GetEntityName()
	} else {
		precreateData.EntityName = dataType.Name()
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
			tag := strings.Split(field.Tag.Get("dbfusion"), ",")[0]
			if tag == "" {
				continue
			}
			value := dataValue.Field(i).Interface()
			mData[tag] = value
			keys += tag + ","
			placeholders += "?,"
			values = append(values, value)
		}
		precreateData.mData = mData
		if len(keys) > 0 {
			precreateData.keys = keys[:len(keys)-1]
			precreateData.placeholders = placeholders[:len(placeholders)-1]
			precreateData.values = values
		}
	} else {
		//This has been checked in this function for type
		precreateData.mData = data.(map[string]interface{})
	}
	precreateData.Data = data
	return
}

func (dbc *DBCommon) PostInsert(cache *caches.Cache, data interface{}, mData map[string]interface{}, dbName string, entityName string) error {
	if val, ok := interface{}(data).(hooks.CacheHook); ok {
		if cache != nil {
			caches.GetInstance().ProcessInsertCache(*cache, val.GetCacheIndexes(), mData, dbName, entityName)
		} else {
			log.Println("WARNING: No valid cache found to process this hook")
		}
	}

	if value, ok := interface{}(data).(hooks.PostInsert); ok {
		value.PostInsert()
	}
	return nil
}
