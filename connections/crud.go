package connections

import (
	"github.com/glodb/dbfusion/queryoptions"
)

type crud interface {
	InsertOne(interface{}) error
	FindOne(interface{}, ...queryoptions.FindOptions) error
	UpdateAndFindOne(interface{}, interface{}, bool) error
	DeleteOne(interface{}) error
}
