package conditions

// import (
// 	"fmt"
// 	"log"
// 	"strings"

// 	"github.com/glodb/dbfusion/ftypes"
// )

type SqlData struct {
	Query             string
	Values            []interface{}
	CacheKey          string
	CacheValues       string
	QueryDefaultCache bool
}

func (sd *SqlData) GetQuery() interface{} {
	return sd.Query
}
func (sd *SqlData) GetValues() interface{} {
	return sd.Values

}
func (sd *SqlData) SetQuery(query interface{}) {
	sd.Query = query.(string)
}
func (sd *SqlData) SetValues(values interface{}) {
	sd.Values = values.([]interface{})
}

func (sd *SqlData) GetCacheKey() string {
	return sd.CacheKey
}
func (sd *SqlData) SetCacheKey(cacheKey string) {
	sd.CacheKey = cacheKey
}

func (sd *SqlData) GetCacheValues() string {
	return sd.CacheValues
}

func (sd *SqlData) SetCacheValues(cachecValues string) {
	sd.CacheValues = cachecValues
}
