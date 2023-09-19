package conditions

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoData struct {
	Query             primitive.M
	CacheKey          string
	CacheValues       string
	QueryDefaultCache bool
}

func (md mongoData) GetQuery() interface{} {
	return md.Query
}
func (md mongoData) GetValues() interface{} {
	return nil

}
func (md mongoData) GetCacheKey() string {
	return md.CacheKey
}

func (md mongoData) GetCacheValues() string {
	return md.CacheValues
}

func (md mongoData) ShouldQueryDefaultCache() bool {
	return md.QueryDefaultCache
}

type MongoCondition struct {
	ConditionsBase
	ConditionType ConditionType
	Key           string
	Value         interface{}
	innerQuery    mongoData
	checkInner    bool
}

func (mc MongoCondition) checkAndInitializeMap() MongoCondition {

	mc.checkInner = true
	mc.innerQuery.Query = make(primitive.M)
	return mc
}

func (mc MongoCondition) buildCondition(condition ConditionBuilder) mongoData {
	query := mongoData{}
	query.Query = make(primitive.M)
	query.QueryDefaultCache = true
	switch condition.getOperator() {
	case EQUAL:
		{
			query.Query[condition.getKey()] = condition.getValue()
			query.CacheKey = fmt.Sprintf("EQUAL_%v", condition.getValue())
			query.CacheValues = fmt.Sprintf("%v", condition.getValue())
		}
	case NOT_EQUAL:
		{
			query.Query[condition.getKey()] = bson.M{"$ne": condition.getValue()}
			query.CacheKey = fmt.Sprintf("NOT_EQUAL_%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case GREATER_THAN:
		{
			query.Query[condition.getKey()] = bson.M{"$gt": condition.getValue()}
			query.CacheKey = fmt.Sprintf("GREATER_THAN_%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case GREATER_THAN_EQUAL:
		{
			query.Query[condition.getKey()] = bson.M{"$gte": condition.getValue()}
			query.CacheKey = fmt.Sprintf("GREATER_THAN_EQUAL_%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case LESSER_THAN:
		{
			query.Query[condition.getKey()] = bson.M{"$lt": condition.getValue()}
			query.CacheKey = fmt.Sprintf("LESSER_THAN_%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case LESSER_THAN_EQUAL:
		{
			query.Query[condition.getKey()] = bson.M{"$lte": condition.getValue()}
			query.CacheKey = fmt.Sprintf("LESSER_THAN_EQUAL_%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case LIKE:
		{
			query.Query[condition.getKey()] = bson.M{"$regex": condition.getValue()}
			query.CacheKey = fmt.Sprintf("LIKE_%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case NOT_LIKE:
		{
			query.Query[condition.getKey()] = bson.M{"$not": bson.M{"$regex": condition.getValue()}}
			query.CacheKey = fmt.Sprintf("NOT_LIKE_%v", condition.getValue())
			query.QueryDefaultCache = false
		}
	case IN:
		{
			if mc.isArrayOrSlice(condition.getValue()) {
				query.Query[condition.getKey()] = bson.M{"$in": condition.getValue()}
				query.CacheKey = fmt.Sprintf("IN_%v", condition.getValue())
				query.QueryDefaultCache = false
			} else {
				log.Println("WARNING: in query requires array or slic, skipping condition")
			}
		}
	case NOT_IN:
		{
			if mc.isArrayOrSlice(condition.getValue()) {
				query.Query[condition.getKey()] = bson.M{"$nin": condition.getValue()}
				query.CacheKey = fmt.Sprintf("NOT_IN_%v", condition.getValue())
				query.QueryDefaultCache = false
			} else {
				log.Println("WARNING: not in query requires array or slic, skipping condition")
			}
		}
	case IS_NULL:
		{
			query.Query[condition.getKey()] = bson.M{"$exists": false}
			query.QueryDefaultCache = false
		}
	case IS_NOT_NULL:
		{
			query.Query[condition.getKey()] = bson.M{"$exists": true}
			query.QueryDefaultCache = false
		}
	}
	return query
}

func (mc MongoCondition) And(conditions ...ConditionBuilder) ConditionBuilder {
	mc = mc.checkAndInitializeMap()
	data := []bson.M{}
	for _, condition := range conditions {
		if condition.isInner() {
			localQuery := condition.getInnerQuery().(mongoData)
			data = append(data, localQuery.Query)
			mc.innerQuery.CacheKey += "_" + localQuery.CacheKey
			mc.innerQuery.CacheValues += localQuery.CacheValues
			mc.innerQuery.QueryDefaultCache = localQuery.QueryDefaultCache
		} else {
			localData := mc.buildCondition(condition)
			mc.innerQuery.CacheKey += "_" + localData.CacheKey
			data = append(data, localData.Query)
			mc.innerQuery.CacheValues += "_" + localData.CacheValues
			mc.innerQuery.QueryDefaultCache = localData.QueryDefaultCache
		}
	}
	mc.innerQuery.Query["$and"] = data
	return mc
}

func (mc MongoCondition) Or(conditions ...ConditionBuilder) ConditionBuilder {
	mc = mc.checkAndInitializeMap()
	data := make([]primitive.M, 0)
	for _, condition := range conditions {
		if condition.isInner() {
			localQuery := condition.getInnerQuery().(mongoData)
			data = append(data, localQuery.Query)
			mc.innerQuery.CacheKey += "_" + localQuery.CacheKey
			mc.innerQuery.CacheValues += localQuery.CacheValues
			mc.innerQuery.QueryDefaultCache = localQuery.QueryDefaultCache
		} else {
			localData := mc.buildCondition(condition)
			mc.innerQuery.CacheKey += "_" + localData.CacheKey
			data = append(data, localData.Query)
			mc.innerQuery.CacheValues += "_" + localData.CacheValues
			mc.innerQuery.QueryDefaultCache = localData.QueryDefaultCache
		}
	}
	mc.innerQuery.Query["$or"] = data
	return mc
}

func (mc MongoCondition) Add(conditions ...ConditionBuilder) ConditionBuilder {
	mc = mc.checkAndInitializeMap()
	for _, condition := range conditions {
		if condition.isInner() {
			localQuery := condition.getInnerQuery().(mongoData)
			if val, ok := localQuery.Query["$and"]; ok {
				mc.innerQuery.Query["$and"] = val
				mc.innerQuery.CacheKey += "_AND_" + localQuery.CacheKey
			} else if val, ok := localQuery.Query["$or"]; ok {
				mc.innerQuery.Query["$or"] = val
				mc.innerQuery.CacheKey += "_OR_" + localQuery.CacheKey
			}
			mc.innerQuery.CacheValues += localQuery.CacheValues
			mc.innerQuery.QueryDefaultCache = localQuery.QueryDefaultCache
		} else {
			query := mc.buildCondition(condition)
			mc.innerQuery.Query[condition.getKey()] = query.Query[condition.getKey()]
			mc.innerQuery.CacheKey += "_" + query.CacheKey
			mc.innerQuery.CacheValues += "_" + query.CacheValues
			mc.innerQuery.QueryDefaultCache = query.QueryDefaultCache
		}
	}
	return mc
}

func (mc MongoCondition) GetQuery(conditions ConditionBuilder) DBFusionData {
	mc = mc.checkAndInitializeMap()
	return conditions.getInnerQuery()
}

func (mc MongoCondition) getKey() string {
	return mc.Key
}

func (mc MongoCondition) getOperator() ConditionType {
	return mc.ConditionType
}
func (mc MongoCondition) getValue() interface{} {
	return mc.Value
}

func (mc MongoCondition) isInner() bool {
	return mc.checkInner
}

func (mc MongoCondition) getInnerQuery() DBFusionData {
	return mc.innerQuery
}

func (mc MongoCondition) GroupConditions(ConditionBuilder) ConditionBuilder {
	log.Println("WARNING: Grouping is not supported in mongodb use And, OR combinations")
	return mc
}
