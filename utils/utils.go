package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/ftypes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// utils is a singleton utility class.
type utils struct{}

var (
	instance *utils
	once     sync.Once
)

// GetInstance returns a single instance of the singleton.
func GetInstance() *utils {
	once.Do(func() {
		instance = &utils{}
	})
	return instance
}

// AssignData assigns a value to a variable based on its type using reflection.
// It supports various data types such as string, int, float, bool, etc.,
// and sets the value to the provided reflectValue.
func (u *utils) AssignData(val interface{}, reflectValue reflect.Value) {
	switch reflectValue.Type().String() {
	case "string":
		if reflectValue.CanSet() {
			reflectValue.SetString(fmt.Sprintf("%s", *val.(*interface{})))
		}
	case "int", "int8", "int16", "int32", "int64":
		if reflectValue.CanSet() {
			intVal, _ := strconv.ParseInt(fmt.Sprintf("%s", *val.(*interface{})), 10, 64)
			reflectValue.SetInt(intVal)
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		if reflectValue.CanSet() {
			intVal, _ := strconv.ParseUint(fmt.Sprintf("%s", *val.(*interface{})), 10, 64)
			reflectValue.SetUint(intVal)
		}
	case "float32", "float64":
		if reflectValue.CanSet() {
			intVal, _ := strconv.ParseFloat(fmt.Sprintf("%s", *val.(*interface{})), 10)
			reflectValue.SetFloat(intVal)
		}
	case "complex32", "complex64":
		if reflectValue.CanSet() {
			intVal, _ := strconv.ParseComplex(fmt.Sprintf("%s", *val.(*interface{})), 10)
			reflectValue.SetComplex(intVal)
		}
	case "bool":
		if reflectValue.CanSet() {
			intVal, _ := strconv.ParseBool(fmt.Sprintf("%s", *val.(*interface{})))
			reflectValue.SetBool(intVal)
		}
	case "[]byte":
		if reflectValue.CanSet() {
			intVal, _ := (*val.(*interface{})).([]byte)
			reflectValue.SetBytes(intVal)
		}
	}
}

// buildSqlData constructs SQL data for a key-value pair and appends it to the provided query and values.
// It handles cases where the key contains "IN" to build SQL IN clauses.
func (u *utils) buildSqlData(key string, val interface{}, cacheKey *string, values *string, query *string, valuesInterface *[]interface{}) {
	tempKey := key

	if strings.Contains(strings.ToLower(key), " in ") || strings.Contains(strings.ToLower(key), " in") {
		tempKey = "IN"
	}

	switch tempKey {
	case "IN":
		inValues := val.([]interface{})
		*valuesInterface = append(*valuesInterface, inValues...)
		inquery := "("
		for idx := range inValues {
			if idx == 0 {
				inquery += "?"
			} else {
				inquery += ",?"
			}
		}
		if *cacheKey == "" {
			*cacheKey += fmt.Sprintf("%s_%v", key, val)
		} else {
			*cacheKey += fmt.Sprintf("_%s_%v", key, val)
		}
		inquery += ")"
		if *values == "" {
			*values += fmt.Sprintf("%v", val)
		} else {
			*values += fmt.Sprintf("_%v", val)
		}
		*query += fmt.Sprintf("%s %s ", key, inquery)
	default:
		if val != nil {
			*valuesInterface = append(*valuesInterface, val)
			if *values == "" {
				*values += fmt.Sprintf("%v", val)
			} else {
				*values += fmt.Sprintf("_%v", val)
			}
			*query += fmt.Sprintf("%s ? ", key)
		} else {
			*query += fmt.Sprintf("%s", key)
		}
		if *cacheKey == "" {
			*cacheKey += fmt.Sprintf("%s_%v", key, val)
		} else {
			*cacheKey += fmt.Sprintf("_%s_%v", key, val)
		}
	}
}

// GetSqlFusionData constructs a SQL DBFusionData object based on the provided data.
// It converts data of various types, such as QMap, DMap, or a map[string]interface{},
// into SQL-compatible DBFusionData.
func (u *utils) GetSqlFusionData(data interface{}) (conditions.DBFusionData, error) {
	dbFusionData := &conditions.SqlData{}
	valuesInterface := make([]interface{}, 0)
	values := ""
	cacheKey := ""
	query := ""
	if value, ok := data.(conditions.DBFusionData); ok {
		dbFusionData = value.(*conditions.SqlData)
		return dbFusionData, nil
	} else if value, ok := data.(ftypes.QMap); ok {
		for key, val := range value {
			u.buildSqlData(key, val, &cacheKey, &values, &query, &valuesInterface)
		}
	} else if value, ok := data.(ftypes.DMap); ok {
		for _, val := range value {
			u.buildSqlData(val.Key, val.Value, &cacheKey, &values, &query, &valuesInterface)
		}
	} else if value, ok := data.(map[string]interface{}); ok {
		for key, val := range value {
			u.buildSqlData(key, val, &cacheKey, &values, &query, &valuesInterface)
		}
	} else {
		return dbFusionData, dbfusionErrors.ErrInvalidType
	}
	dbFusionData.SetCacheKey(cacheKey)
	dbFusionData.SetCacheValues(values)
	dbFusionData.SetValues(valuesInterface)
	dbFusionData.SetQuery(query)
	return dbFusionData, nil
}

// buildMongoData constructs MongoDB data for a key-value pair and appends it to the provided query and values.
// It also updates cacheKey and values for MongoDB data.
func (u *utils) buildMongoData(key string, val interface{}, cacheKey *string, values *string) primitive.E {
	if *cacheKey == "" {
		*cacheKey += fmt.Sprintf("%s_%v", key, val)
	} else {
		*cacheKey += fmt.Sprintf("_%s_%v", key, val)
	}
	if *values == "" {
		*values += fmt.Sprintf("%v", val)
	} else {
		*values += fmt.Sprintf("_%v", val)
	}

	return primitive.E{Key: key, Value: val}
}

// GetMongoFusionData constructs a MongoDB DBFusionData object based on the provided data.
// It converts data of various types, such as QMap, DMap, or a map[string]interface{},
// into MongoDB-compatible DBFusionData.
func (u *utils) GetMongoFusionData(data interface{}) (conditions.DBFusionData, error) {
	dbFusionData := &conditions.MongoData{}
	values := ""
	cacheKey := ""
	query := primitive.D{}
	if value, ok := data.(conditions.DBFusionData); ok {
		dbFusionData = value.(*conditions.MongoData)
		return dbFusionData, nil
	} else if value, ok := data.(ftypes.QMap); ok {
		for key, val := range value {
			singleData := u.buildMongoData(key, val, &cacheKey, &values)
			query = append(query, singleData)
		}
	} else if value, ok := data.(ftypes.DMap); ok {
		for _, val := range value {
			singleData := u.buildMongoData(val.Key, val.Value, &cacheKey, &values)
			query = append(query, singleData)
		}
	} else if value, ok := data.(map[string]interface{}); ok {
		for key, val := range value {
			singleData := u.buildMongoData(key, val, &cacheKey, &values)
			query = append(query, singleData)
		}
	} else {
		return dbFusionData, dbfusionErrors.ErrInvalidType
	}
	dbFusionData.SetCacheKey(cacheKey)
	dbFusionData.SetCacheValues(values)
	dbFusionData.SetQuery(query)
	return dbFusionData, nil
}
