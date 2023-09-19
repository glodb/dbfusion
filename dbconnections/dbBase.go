package dbconnections

import (
	"reflect"

	"github.com/glodb/dbfusion/caches"
)

type entityData struct {
	entityName string
	structType int
	dataType   reflect.Type
	dataValue  reflect.Value
}

type preCreateReturn struct {
	entityName   string
	mData        map[string]interface{}
	Data         interface{}
	keys         string
	placeholders string
	values       []interface{}
}

type preFindReturn struct {
	entityName    string
	query         interface{}
	whereQuery    interface{}
	queryDatabase bool
}

type DBBase interface {
	ChangeDatabase(dbName string) error
	SetCache(*caches.Cache)
	CreateTable(ifNotExist bool)
}
