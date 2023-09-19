package query

import (
	"github.com/glodb/dbfusion/queryoptions"
)

type CRUD interface {
	InsertOne(interface{}) error
	FindOne(interface{}, ...queryoptions.FindOptions) error
	UpdateOne(interface{}) error
	DeleteOne(interface{}) error
}
