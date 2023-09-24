package hooks

type NormalIndexes interface {
	GetNormalIndexes() []string
}

type UniqueIndexes interface {
	GetUniqueIndexes() []string
}

type TextIndexes interface {
	GetTextIndex() string
}

type TwoDimensialIndexes interface {
	Get2DIndexes() []string
}

type TwoDimensialSpatialIndexes interface {
	Get2DSpatialIndexes() []string
}

type HashedIndexes interface {
	GetHashedIndexes() []string
}

type SparseIndexes interface {
	GetSparseIndexes() []string
}
