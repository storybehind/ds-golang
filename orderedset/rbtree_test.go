package orderedset_test

import (
	"testing"

	orderedset "github.com/storybehind/gocontainer/orderedset"
)

func TestRbTree(t *testing.T) {
	rbTree := orderedset.NewRbTree[int](func(k1, k2 int) bool { return k1 < k2 })
	testOrderedSet(t, rbTree)
}

func TestRbTreeIterator(t *testing.T) {
	rbTree := orderedset.NewRbTree[int](func(k1, k2 int) bool { return k1 < k2 })
	testOrderedSetForwardIterator(t, rbTree)
}

func TestRbTreeReverseIterator(t *testing.T) {
	rbTree := orderedset.NewRbTree[int](func(k1, k2 int) bool { return k1 < k2 })
	testOrderedSetReverseIterator(t, rbTree)
}
