package priorityqueue_test

import (
	"sort"
	"testing"

	"github.com/storybehind/gocontainer/priorityqueue"
)

func TestNewBinaryHeap(t *testing.T) {
	minHeap := priorityqueue.NewBinaryHeap(func(v1, v2 int) bool {
		return v1 < v2
	})
	checkLen(t, minHeap, 0)

	node5 := minHeap.Push(5)
	checkLen(t, minHeap, 1)
	checkTop(t, minHeap, 5)

	node2 := minHeap.Push(2)
	checkLen(t, minHeap, 2)
	checkTop(t, minHeap, 2)

	node7 := minHeap.Push(7)
	checkLen(t, minHeap, 3)
	checkTop(t, minHeap, 2)

	node51 := minHeap.Push(5)
	checkLen(t, minHeap, 4)
	checkTop(t, minHeap, 2)

	minHeap.Update(node5, 1)
	checkLen(t, minHeap, 4)
	checkTop(t, minHeap, 1)

	checkPop(t, minHeap, 1)
	checkLen(t, minHeap, 3)
	checkTop(t, minHeap, 2)

	checkRemove(t, minHeap, node5)
	checkLen(t, minHeap, 3)
	checkTop(t, minHeap, 2)

	checkRemove(t, minHeap, node7)
	checkLen(t, minHeap, 2)
	checkTop(t, minHeap, 2)

	checkPop(t, minHeap, 2)
	checkLen(t, minHeap, 1)
	checkTop(t, minHeap, 5)

	checkRemove(t, minHeap, node51)
	checkLen(t, minHeap, 0)

	checkRemove(t, minHeap, node2)
	checkLen(t, minHeap, 0)
}

func TestBinaryHeapInit(t *testing.T) {
	initValues := []int{4, 1, 3, 2, 5, 6, 5}
	bh := priorityqueue.InitBinaryHeap[int](func(v1, v2 int) bool {return v1 < v2}, initValues)
	sort.Slice(initValues, func (i, j int) bool {
		return initValues[i] < initValues[j]
	})
	bh.Push(0)
	if v := bh.Pop(); v != 0 {
		t.Errorf("expected 0; found: %d", v)
	}
	for i:=0; bh.Len() > 0; i++ {
		if v := bh.Pop(); v != initValues[i] {
			t.Errorf("expected non decreasing order; exp: %d, found: %d", initValues[i], v)
		}
	}
}

func checkLen[V any](t *testing.T, bh *priorityqueue.BinaryHeap[V], len int64) {
	if n := bh.Len(); n != len {
		t.Errorf("mh.Len() = %d, want= %d", n, len)
	}
}

func checkTop[V comparable](t *testing.T, bh *priorityqueue.BinaryHeap[V], topExpected V) {
	topFoundNode := bh.Top()
	if topExpected != topFoundNode.GetValue() {
		t.Errorf("topFound = %v, topExpected = %v", topFoundNode.GetValue(), topExpected)
	}
}

func checkPop[V comparable](t *testing.T, bh *priorityqueue.BinaryHeap[V], expected V) {
	found := bh.Pop()
	if found != expected {
		t.Errorf("found = %v, expected = %v", found, expected)
	}
}

func checkRemove[V comparable](t *testing.T, bh *priorityqueue.BinaryHeap[V], bhn *priorityqueue.BinaryHeapNode[V]) {
	val := bh.Remove(bhn)
	if val != bhn.GetValue() {
		t.Errorf("val = %v, bhn.GetValue() = %v ; val & bhn.GetValue() must match after removal", val, bhn.GetValue())
	}
}
