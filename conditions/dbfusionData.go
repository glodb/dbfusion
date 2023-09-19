package conditions

type DBFusionData interface {
	GetQuery() interface{}
	SetQuery(interface{})
	GetValues() interface{}
	SetValues(interface{})
	SetCacheKey(string)
	GetCacheKey() string
	SetCacheValues(string)
	GetCacheValues() string
}
