package orderedset

type Iterator[K Less[K]] interface {
	Next() K
	Prev() K
	HasNext() bool
	HasPrev() bool
}