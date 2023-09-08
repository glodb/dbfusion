package hooks

type CacheHook interface {
	//Comma seprated cache indexes are required
	GetCacheIndexes() []string
}
