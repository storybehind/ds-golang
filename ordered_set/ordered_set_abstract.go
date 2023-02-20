package orderedset

// OrderedSet interface
type OrderedSet[K Less[K]] interface {
	Get(key K) (_ K, _ bool)
	GetGreater(key K) (_ K, _ bool)
	GetGreaterThanOrEqual(key K) (_ K, _ bool)
	GetLower(key K) (_ K, _ bool)
	GetLowerThanOrEqual(key K) (_ K, _ bool)
	// Has(key K) bool

	Max() (_ K, _ bool)
	Min() (_ K, _ bool)
	Len() int64

	ReplaceOrInsert(key K) (_ K, _ bool)

	Delete(key K) (K, bool)
	DeleteMax() (K, bool)
	DeleteMin() (K, bool)
}

type Less[K any] interface {
	Less(K) bool
}