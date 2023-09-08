package dbconnections

import "github.com/glodb/dbfusion/caches"

type PreCreateReturn struct {
	EntityName   string
	mData        map[string]interface{}
	Data         interface{}
	keys         string
	placeholders string
	values       []interface{}
}
type DBBase interface {
	ChangeDatabase(dbName string) error
	SetCache(*caches.Cache)
	PreInsert(data interface{}) (entityName PreCreateReturn, err error)
	PostInsert(cache *caches.Cache, data interface{}, mData map[string]interface{}, dbName string, entityName string) error
}
