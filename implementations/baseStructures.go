package implementations

import "reflect"

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
	dataType      reflect.Type
	dataValue     reflect.Value
}

type preUpdateReturn struct {
	entityName string
	queryData  interface{}
}
