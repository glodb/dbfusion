package implementations

import (
	"fmt"

	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/ftypes"
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
