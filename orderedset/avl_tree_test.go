package orderedset_test

import (
	"testing"

	orderedset "github.com/storybehind/gocontainer/orderedset"
)

func TestAvlTree(t *testing.T) {
	avlTree := orderedset.NewAvlTree[int](func(k1, k2 int) bool { return k1 < k2 })
	testOrderedSet(t, avlTree)
}

func TestAvlTreeIterator(t *testing.T) {
	avlTree := orderedset.NewAvlTree[int](func(k1, k2 int) bool { return k1 < k2 })
	testOrderedSetForwardIterator(t, avlTree)
}

func TestAvlTreeReverseIterator(t *testing.T) {
	avlTree := orderedset.NewAvlTree[int](func(k1, k2 int) bool { return k1 < k2 })
	testOrderedSetReverseIterator(t, avlTree)
}
