package query

type CRUD interface {
	Insert(interface{}) error
	Find(interface{}) error
	Update(interface{}) error
	Delete(interface{}) error
}
