package connections

import (
	"github.com/glodb/dbfusion/caches"
)

type base interface {
	ChangeDatabase(dbName string) error
	SetCache(*caches.Cache)
}
