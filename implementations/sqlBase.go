package implementations

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/utils"
)

type SqlBase struct {
	DBCommon
}

func (sb *SqlBase) createSqlInsert(data interface{}) (string, []interface{}, preCreateReturn, error) {
	preCreateData, err := sb.preInsert(data)
	if err != nil {
		return "", nil, preCreateData, err
	}

	keys := preCreateData.keys
	placeholders := preCreateData.placeholders
	values := preCreateData.values
	if len(keys) <= 0 {
		values = make([]interface{}, 0)
		keys = ""
		placeholders = ""
		for key, value := range preCreateData.mData {
			keys += key + ","
			placeholders += "?,"
			values = append(values, value)
		}
		keys = keys[:len(keys)-1]
		placeholders = placeholders[:len(placeholders)-1]
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", preCreateData.entityName, keys, placeholders)
	return query, values, preCreateData, nil
}

func (sb *SqlBase) createQuery(query ftypes.DMap) (string, []interface{}) {
	stringQuery := ""
	data := make([]interface{}, 0)
	for _, val := range query {
		stringQuery += val.Key
		data = append(data, val.Value)
	}
	return stringQuery, data
}

func (sb *SqlBase) createCountQuery(entityName string) string {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", entityName)

	if sb.joins != "" {
		query = fmt.Sprintf(query+" %s", sb.joins)
	}

	if sb.whereQuery != nil {
		whereData := sb.whereQuery.(*conditions.SqlData)
		if whereData.GetQuery().(string) != "" {
			query = fmt.Sprintf(query+" WHERE %v", whereData.GetQuery())
		}
	}

	if sb.groupBy != "" {
		query = fmt.Sprintf(query+" GROUP BY %s", sb.groupBy)
		if sb.havingString != "" {
			query = fmt.Sprintf(query+" HAVING %s", sb.havingString)
		}
	}

	return query
}
func (sb *SqlBase) createFindQuery(entityName string, limitOne bool) string {
	selectionKeys := "*"
	projections := []string{}

	if sb.projection != nil {
		projections = sb.projection.([]string)
		if len(projections) != 0 {
			selectionKeys = ""
			for id, key := range projections {
				if id != 0 {
					selectionKeys += ", "
				}
				selectionKeys += key
			}
		}
	}

	query := fmt.Sprintf("SELECT %s FROM %s", selectionKeys, entityName)

	if sb.joins != "" {
		query = fmt.Sprintf(query+" %s", sb.joins)
	}

	if sb.whereQuery != nil {
		whereData := sb.whereQuery.(*conditions.SqlData)
		if whereData.GetQuery().(string) != "" {
			query = fmt.Sprintf(query+" WHERE %v", whereData.GetQuery())
		}
	}

	if sb.groupBy != "" {
		query = fmt.Sprintf(query+" GROUP BY %s", sb.groupBy)
		if sb.havingString != "" {
			query = fmt.Sprintf(query+" HAVING %s", sb.havingString)
		}
	}

	if sb.sort != nil {
		query = fmt.Sprintf(query+" ORDER BY %s", sb.sort.(string))
	}

	if !limitOne {
		if sb.limit != 0 {
			query = fmt.Sprintf(query+" LIMIT %d", sb.limit)
		}
	} else {
		query = fmt.Sprintf(query+" LIMIT %d", 1)
	}

	if sb.skip != 0 {
		query = fmt.Sprintf(query+" OFFSET %d", sb.skip)
	}
	return query
}

func (sb *SqlBase) createUpdateQuery(entityName string, setCommands string, limitOne bool) string {
	query := fmt.Sprintf("UPDATE %s %s", entityName, setCommands)

	if sb.joins != "" {
		query = fmt.Sprintf(query+" %s", sb.joins)
	}

	if sb.whereQuery != nil {
		whereData := sb.whereQuery.(*conditions.SqlData)
		if whereData.GetQuery().(string) != "" {
			query = fmt.Sprintf(query+" WHERE %v", whereData.GetQuery())
		}
	}

	if !limitOne {
		if sb.limit != 0 {
			query = fmt.Sprintf(query+" LIMIT %d", sb.limit)
		}
	} else {
		query = fmt.Sprintf(query+" LIMIT %d", 1)
	}
	return query
}

func (sb *SqlBase) createDeleteQuery(entityName string, whereConditions string, limitOne bool) string {
	query := fmt.Sprintf("DELETE FROM %s", entityName)

	if sb.joins != "" {
		query = fmt.Sprintf(query+" %s", sb.joins)
	}

	if whereConditions != "" {
		query = fmt.Sprintf(query+" WHERE %s", whereConditions)
	} else if sb.whereQuery != nil {
		whereData := sb.whereQuery.(*conditions.SqlData)
		if whereData.GetQuery().(string) != "" {
			query = fmt.Sprintf(query+" WHERE %v", whereData.GetQuery())
		}
	}

	if !limitOne {
		if sb.limit != 0 {
			query = fmt.Sprintf(query+" LIMIT %d", sb.limit)
		}
	} else {
		query = fmt.Sprintf(query+" LIMIT %d", 1)
	}
	return query
}

func (sb *SqlBase) readSqlDataFromRows(rows *sql.Rows, dataType reflect.Type, dataValue reflect.Value) (int, error) {

	rowsCount := 0
	if rows == nil {
		return 0, dbfusionErrors.ErrSQLQueryNoRecordFound
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	// Create a slice of interface{} to hold the column data
	columnData := make([]interface{}, len(columnNames))
	for i := range columnData {
		var v interface{}
		columnData[i] = &v
	}

	tagField := make(map[string]reflect.Value)
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)

		rawtags := strings.Split(field.Tag.Get("dbfusion"), ",")
		tagName := rawtags[0]
		tagField[tagName] = dataValue.Field(i)
	}

	// Iterate through the rows
	for rows.Next() {
		rowsCount++
		// Scan the row data into columnData
		err := rows.Scan(columnData...)
		if err != nil {
			return 0, err
		}

		for idx, name := range columnNames {
			if field, ok := tagField[name]; ok {
				utils.GetInstance().AssignData(columnData[idx], field)
			}
		}

	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return rowsCount, nil
}

func (sb *SqlBase) readSqlRowsToArray(rows *sql.Rows, results interface{}) error {
	// Create a new slice of the same type as results (e.g., &[]Users{})
	resultSliceType := reflect.TypeOf(results).Elem()
	newSlice := reflect.New(resultSliceType).Elem()

	// // Get the type of elements within the slice (e.g., Users)
	// elementType := resultSliceType.Elem()
	// elementValue := reflect.New(elementType).Elem()

	// Get the field names from struct tags

	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}

	// Create a slice of interface{} to hold the column data
	columnData := make([]interface{}, len(columnNames))
	for i := range columnData {
		var v interface{}
		columnData[i] = &v
	}

	// log.Println(fieldNames)
	// Iterate through rows and populate the newSlice
	for rows.Next() {
		// Create a new element of the slice's element type
		elementType := resultSliceType.Elem()
		newElement := reflect.New(elementType).Elem()

		// log.Println(newElement.Field(0))
		// Scan the row into the fields of the newElement
		if err := rows.Scan(columnData...); err != nil {
			return err
		}
		tagField := sb.getFieldNames(elementType, newElement)
		for idx, name := range columnNames {
			if fieldName, ok := tagField[name]; ok {
				utils.GetInstance().AssignData(columnData[idx], newElement.FieldByName(fieldName))
			}
		}
		// Append the newElement to the newSlice
		newSlice = reflect.Append(newSlice, newElement)
	}

	// Set the populated newSlice to the results pointer
	reflect.ValueOf(results).Elem().Set(newSlice)

	return nil
}

func (sb *SqlBase) getFieldNames(structType reflect.Type, dataValue reflect.Value) map[string]string {
	tagField := make(map[string]string)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		rawtags := strings.Split(field.Tag.Get("dbfusion"), ",")
		tagName := rawtags[0]
		if tagName != "" {
			tagField[tagName] = field.Name
		}
	}
	return tagField
}

func (sb *SqlBase) createTableQuery(data interface{}, ifNotExist bool) (string, error) {
	query := ""

	name, err := sb.getEntityName(data)

	if err != nil {
		return "", err
	}
	if ifNotExist {
		query = `CREATE TABLE IF NOT EXISTS ` + name.entityName + ` (`
	} else {
		query = `CREATE TABLE ` + name.entityName + ` (`
	}

	dataType := name.dataType

	if dataType.Kind() != reflect.Struct {
		return "", dbfusionErrors.ErrInvalidType
	}

	columns := ""

	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		tags := strings.Split(field.Tag.Get("dbfusion"), ",")

		if columns != "" {
			columns += ","
		}

		columns += strings.Join(tags, " ")
	}

	query += columns + ");"

	return query, nil
}
