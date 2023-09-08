package caches

/*
main interface for the caches all the caches need to implement this interface to be used with dbFusion
*/

type Cache interface {
	ConnectCache(connectionUri string, password ...string) error
	IsConnected() bool
	DisconnectCache()
	GetKey()
	SetKey(key string, value interface{})
	DeleteKey()
	UpdateKey()
}
