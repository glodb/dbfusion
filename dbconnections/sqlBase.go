package dbconnections

import (
	"fmt"
	"log"

	"github.com/glodb/dbfusion/query"
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

func (sb *SqlBase) createQuery(query query.DMap) (string, []interface{}) {
	stringQuery := ""
	data := make([]interface{}, 0)
	for _, val := range query {
		stringQuery += val.Key
		data = append(data, val.Value)
	}
	return stringQuery, data
}

func (sb *SqlBase) createFindQuery(entityName string, conditions string) {
	selectionKeys := "*"
	projections := sb.projection.([]string)
	if len(projections) != 0 {
		selectionKeys = ""
		for id, key := range projections {
			if id != 0 {
				selectionKeys += ", "
			}
			selectionKeys += key
		}
	}

	query := fmt.Sprintf("SELECT %s FROM %s", selectionKeys, entityName)

	if sb.joins != "" {
		query = fmt.Sprintf(query+" %s", sb.joins)
	}

	if conditions != "" {
		query = fmt.Sprintf(query+" WHERE %s", conditions)
	}
	log.Println(query)
}
