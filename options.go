package dbfusion

import (
	"github.com/glodb/dbfusion/caches"
)

type Options struct {
	DbName          *string
	Uri             *string
	CertificatePath *string
	Cache           caches.Cache
}
