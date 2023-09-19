package conditions

import (
	"fmt"
	"log"
	"strings"
)

type SqlData struct {
	Query             string
	Values            []interface{}
	CacheKey          string
	CacheValues       string
	QueryDefaultCache bool
}

func (sd SqlData) GetQuery() interface{} {
	return sd.Query
}
func (sd SqlData) GetValues() interface{} {
	return sd.Values

}
func (sd SqlData) GetCacheKey() string {
	return sd.CacheKey
}

func (sd SqlData) GetCacheValues() string {
	return sd.CacheValues
}

func (sd SqlData) ShouldQueryDefaultCache() bool {
	return sd.QueryDefaultCache
}

type SqlCondition struct {
	ConditionsBase
	ConditionType ConditionType
	Key           string
	Value         interface{}
	AddNot        bool
	innerQuery    SqlData
	checkInner    bool
}

func (sc SqlCondition) initializeQuery() SqlCondition {

	sc.checkInner = true
	sc.innerQuery = SqlData{}
	sc.innerQuery.Values = make([]interface{}, 0)
	sc.innerQuery.CacheValues = ""
	return sc
}

func (sc SqlCondition) buildArraySqlData(condition ConditionBuilder, operator string) SqlData {
	query := SqlData{}
	query.Values = make([]interface{}, 0)

	values := condition.getValue().([]interface{})

	if len(values) <= 0 {
		return query
	}

	query.Query = condition.getKey() + " " + operator + " ("

	if operator == "NOT IN" {
		query.CacheKey = "NOT_IN" + "_"
	} else {
		query.CacheKey = operator + "_"
	}

	for _, val := range values {
		query.Query += "?, "
		query.Values = append(query.Values, val)
		query.CacheKey += fmt.Sprintf("%v", val)
	}

	query.Query = query.Query[:len(query.Query)-2]
	query.Query += ")"

	return query
}

func (sc SqlCondition) buildCondition(condition ConditionBuilder) SqlData {
	query := SqlData{}
	query.QueryDefaultCache = true
	query.Values = make([]interface{}, 0)

	switch condition.getOperator() {
	case EQUAL:
		{
			query.Query = condition.getKey() + " = ?"
			query.Values = append(query.Values, condition.getValue())
			query.CacheKey = fmt.Sprintf("EQUAL_%v", condition.getValue())
			query.CacheValues = fmt.Sprintf("%v", condition.getValue())
		}
	case NOT_EQUAL:
		{
			query.Query = condition.getKey() + " <> ?"
			query.Values = append(query.Values, condition.getValue())
			query.CacheKey = fmt.Sprintf("NOT_EQUAL_%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case GREATER_THAN:
		{
			query.Query = condition.getKey() + " > ?"
			query.Values = append(query.Values, condition.getValue())
			query.CacheKey = fmt.Sprintf("GREATER_THAN_%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case GREATER_THAN_EQUAL:
		{
			query.Query = condition.getKey() + " >= ?"
			query.Values = append(query.Values, condition.getValue())
			query.CacheKey = fmt.Sprintf("GREATER_THAN_EQUAL%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case LESSER_THAN:
		{
			query.Query = condition.getKey() + " < ?"
			query.Values = append(query.Values, condition.getValue())
			query.CacheKey = fmt.Sprintf("LESSER_THAN%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case LESSER_THAN_EQUAL:
		{
			query.Query = condition.getKey() + " <= ?"
			query.Values = append(query.Values, condition.getValue())
			query.CacheKey = fmt.Sprintf("LESSER_THAN_EQUAL%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case LIKE:
		{
			query.Query = condition.getKey() + " LIKE ?"
			query.Values = append(query.Values, condition.getValue())
			query.CacheKey = fmt.Sprintf("LIKE_%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case IN:
		{
			if sc.isArrayOrSlice(condition.getValue()) {
				query = sc.buildArraySqlData(condition, "IN")
			} else {
				log.Println("WARNING: in query requires array or slice, skipping condition")
			}
			query.QueryDefaultCache = false
		}
	case NOT_IN:
		{
			if sc.isArrayOrSlice(condition.getValue()) {
				query = sc.buildArraySqlData(condition, "NOT IN")
			} else {
				log.Println("WARNING: mot in query requires array or slice, skipping condition")
			}
			query.QueryDefaultCache = false
		}
	case IS_NULL:
		{
			query.Query = condition.getKey() + " IS NULL"
			query.CacheKey = fmt.Sprintf("IS_NULL%v", condition.getKey())
			query.QueryDefaultCache = false
		}
	case IS_NOT_NULL:
		{
			query.Query = condition.getKey() + " IS NOT NULL"
			query.CacheKey = fmt.Sprintf("IS_NOT_NULL%v", condition.getKey())
			query.QueryDefaultCache = false
		}
	case NOT_LIKE:
		{
			query.Query = condition.getKey() + " NOT LIKE ?"
			query.CacheKey = fmt.Sprintf("NOT_LIKE_%v", condition.getValue())
			query.Values = append(query.Values, condition.getValue())
			query.QueryDefaultCache = false
		}
	}
	return query
}

func (sc SqlCondition) buildQuery(queryPrefix string, conditions ...ConditionBuilder) ConditionBuilder {
	sc = sc.initializeQuery()
	for _, condition := range conditions {
		if condition.isInner() {
			localQuery := condition.getInnerQuery().(SqlData)
			sc.innerQuery.Query += condition.getInnerQuery().(SqlData).Query
			sc.innerQuery.Values = append(sc.innerQuery.Values, localQuery.Values...)
			sc.innerQuery.CacheKey += "_" + localQuery.CacheKey
			sc.innerQuery.QueryDefaultCache = localQuery.QueryDefaultCache
			sc.innerQuery.CacheValues += localQuery.CacheValues
		} else {
			localQuery := sc.buildCondition(condition)
			sc.innerQuery.Query += queryPrefix + localQuery.Query
			sc.innerQuery.Values = append(sc.innerQuery.Values, localQuery.Values...)
			sc.innerQuery.CacheKey += "_" + localQuery.CacheKey
			sc.innerQuery.QueryDefaultCache = localQuery.QueryDefaultCache
			sc.innerQuery.CacheValues += "_" + localQuery.CacheValues
		}
	}
	return sc
}

func (sc SqlCondition) And(conditions ...ConditionBuilder) ConditionBuilder {
	return sc.buildQuery(" AND ", conditions...)
}

func (sc SqlCondition) Or(conditions ...ConditionBuilder) ConditionBuilder {
	return sc.buildQuery(" OR ", conditions...)
}

func (sc SqlCondition) GroupConditions(conditon ConditionBuilder) ConditionBuilder {
	sc = sc.initializeQuery()

	innerQuery := conditon.getInnerQuery().(SqlData)

	if strings.HasPrefix(innerQuery.Query, " AND ") {
		innerQuery.Query = innerQuery.Query[5:]
		innerQuery.Query = " AND (" + innerQuery.Query
	} else if strings.HasPrefix(innerQuery.Query, " OR ") {
		innerQuery.Query = innerQuery.Query[4:]
		innerQuery.Query = " OR (" + innerQuery.Query
	} else {
		innerQuery.Query = "(" + innerQuery.Query
	}

	sc.innerQuery.Query = innerQuery.Query + ")"
	sc.innerQuery.Values = append(sc.innerQuery.Values, conditon.getInnerQuery().(SqlData).Values...)
	sc.innerQuery.CacheKey = innerQuery.CacheKey
	sc.innerQuery.CacheValues += innerQuery.CacheValues
	sc.innerQuery.QueryDefaultCache = innerQuery.QueryDefaultCache
	return sc
}

func (sc SqlCondition) Add(conditions ...ConditionBuilder) ConditionBuilder {
	for idx, condition := range conditions {
		query := ""
		if condition.isInner() {
			localQuery := condition.getInnerQuery().(SqlData)
			sc.innerQuery.Query += query + localQuery.Query
			sc.innerQuery.Values = append(sc.innerQuery.Values, localQuery.Values...)
			sc.innerQuery.CacheKey += "_" + localQuery.CacheKey
			sc.innerQuery.QueryDefaultCache = localQuery.QueryDefaultCache
			sc.innerQuery.CacheValues += localQuery.CacheValues
		} else {
			localQuery := sc.buildCondition(condition)
			sc.innerQuery.Query += query + localQuery.Query
			sc.innerQuery.Values = append(sc.innerQuery.Values, localQuery.Values...)
			sc.innerQuery.CacheKey += "_" + localQuery.CacheKey
			sc.innerQuery.QueryDefaultCache = localQuery.QueryDefaultCache
			sc.innerQuery.CacheValues += "_" + localQuery.CacheValues
		}

		if idx == 0 {
			if strings.HasPrefix(sc.innerQuery.Query, " AND ") {
				sc.innerQuery.Query = sc.innerQuery.Query[5:]
			}

			if strings.HasPrefix(sc.innerQuery.Query, " OR ") {
				sc.innerQuery.Query = sc.innerQuery.Query[4:]
			}
		}
	}
	return sc
}

func (sc SqlCondition) getKey() string {
	return sc.Key
}

func (sc SqlCondition) getOperator() ConditionType {
	return sc.ConditionType
}
func (sc SqlCondition) getValue() interface{} {
	return sc.Value
}

func (sc SqlCondition) isInner() bool {
	return sc.checkInner
}

func (sc SqlCondition) getInnerQuery() DBFusionData {
	return sc.innerQuery
}

func (sc SqlCondition) GetQuery(condition ConditionBuilder) DBFusionData {
	sc = sc.initializeQuery()
	return condition.getInnerQuery()
}
