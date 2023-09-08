package hooks

type PreInsert interface {
	PreInsert() PreInsert
}

type PostInsert interface {
	PostInsert()
}
