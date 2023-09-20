package conditions

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoData struct {
	Query       primitive.D
	CacheKey    string
	CacheValues string
}

func (md *MongoData) GetQuery() interface{} {
	if md.Query == nil {
		md.Query = primitive.D{}
	}
	return md.Query
}

func (md *MongoData) SetQuery(query interface{}) {
	md.Query = query.(primitive.D)
}
func (md *MongoData) GetValues() interface{} {
	return nil
}
func (md *MongoData) SetValues(data interface{}) {
}

func (md *MongoData) GetCacheKey() string {
	return md.CacheKey
}

func (md *MongoData) SetCacheKey(cacheKey string) {
	md.CacheKey = cacheKey
}
func (md *MongoData) GetCacheValues() string {
	return md.CacheValues
}
func (md *MongoData) SetCacheValues(values string) {
	md.CacheValues = values
}
