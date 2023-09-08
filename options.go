package dbfusion

import (
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/dbconnections"
)

type Options struct {
	DbName          *string
	Uri             *string
	DbType          dbconnections.DBTypes
	CertificatePath *string
	Cache           caches.Cache
}
