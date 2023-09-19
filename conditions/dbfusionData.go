package conditions

type DBFusionData interface {
	GetQuery() interface{}
	GetValues() interface{}
	GetCacheKey() string
	GetCacheValues() string
	ShouldQueryDefaultCache() bool
}
