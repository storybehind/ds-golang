package priorityqueue

// Binary heap node 
type BinaryHeapNode[V any] struct {
	index      int64
	binaryHeap *BinaryHeap[V]
	value      V
}

// Returns node's value
func (bhn *BinaryHeapNode[V]) GetValue() V {
	return bhn.value
}

// Implements priorityqueue with Push, Pop, Top, Update and Remove operations
type BinaryHeap[V any] struct {
	nodes        []*BinaryHeapNode[V]
	priorityFunc func(v1, v2 V) bool
	length       int64
}

// Returns instance of BinaryHeap. 
// v1 has higher priority than v2 if priorityFunc(v1, v2) returns true
func NewBinaryHeap[V any](priorityFunc func(v1, v2 V) bool) *BinaryHeap[V] {
	return &BinaryHeap[V]{
		nodes:        make([]*BinaryHeapNode[V], 1),
		priorityFunc: priorityFunc,
		length:       0,
	}
}

// Returns instance of BinaryHeap. 
// v1 has higher priority than v2 if priorityFunc(v1, v2) returns true
// Takes O(n) time where n = len(initValues)
func InitBinaryHeap[V any](priorityFunc func(v1, v2 V) bool, initValues []V) *BinaryHeap[V] {
	n := len(initValues)
	bh := &BinaryHeap[V] {
		nodes: make([]*BinaryHeapNode[V], n+1),
		priorityFunc: priorityFunc,
		length: int64(n),
	}
	for i, v := range initValues {
		bh.nodes[i+1] = &BinaryHeapNode[V] {
			index: int64(i+1),
			binaryHeap: bh,
			value: v,
		}
	}
	for i:=n/2; i>=1; i-- {
		bh.siftDown(int64(i))
	}
	return bh
}

// Insert value into the binary heap and returns node's pointer to pushed value.
// Takes O(log n) time where n is number of values in the queue.
func (bh *BinaryHeap[V]) Push(value V) *BinaryHeapNode[V] {
	newHeapNode := &BinaryHeapNode[V]{
		index:      int64(len(bh.nodes)),
		binaryHeap: bh,
		value:      value,
	}
	bh.nodes = append(bh.nodes, newHeapNode)
	bh.sift(newHeapNode.index)
	bh.length++
	return newHeapNode
}

// Returns node's pointer of highest priority value. Panics if binary heap is empty. Takes O(1) time.
func (bh *BinaryHeap[V]) Top() (*BinaryHeapNode[V]) {
	return bh.nodes[1]
}

// Remove highest priority value and returns it. Panics if binary heap is empty. Takes O(log n) time where n is number of values in the queue.
func (bh *BinaryHeap[V]) Pop() V {
	return bh.Remove(bh.nodes[1])
}

// Deletes node in binary heap. Takes O(log n) time where n is number of values in the queue.
func (bh *BinaryHeap[V]) Remove(node *BinaryHeapNode[V]) V {
	if node.binaryHeap != bh {
		return node.value
	}
	if bh.length == node.index {
		bh.nodes[bh.length] = nil
		bh.nodes = bh.nodes[:bh.length]
		bh.length--

		node.binaryHeap = nil
		node.index = 0
		return node.value
	}
	nodeIndex := node.index
	bh.swapNodes(bh.length, node.index)
	bh.nodes[bh.length] = nil
	bh.nodes = bh.nodes[:bh.length]
	bh.length--

	node.binaryHeap = nil
	node.index = 0
	
	bh.sift(nodeIndex)
	return node.value
}

// Update node's value to newValue. Takes O(log n) time where n is number of values in the queue.
func (bh *BinaryHeap[V]) Update(node *BinaryHeapNode[V], newValue V) {
	if node.binaryHeap != bh {
		return
	}
	node.value = newValue
	bh.sift(node.index)
}


// Returns number of values currently in the queue
func (bh *BinaryHeap[V]) Len() int64 {
	return bh.length
}

func (bh *BinaryHeap[V]) sift(nodeIndex int64) {
	isSiftUp := (nodeIndex >> 1) > 0 && bh.priorityFunc(bh.nodes[nodeIndex].GetValue(), bh.nodes[nodeIndex >> 1].GetValue())
	if isSiftUp {
		bh.siftUp(nodeIndex)
		return
	}
	bh.siftDown(nodeIndex)
}

func (bh *BinaryHeap[V]) siftUp(nodeIndex int64) {
	for nodeIndex > 1 && bh.priorityFunc(bh.nodes[nodeIndex].GetValue(), bh.nodes[nodeIndex >> 1].GetValue()) {
		bh.swapNodes(nodeIndex, nodeIndex >> 1)
		nodeIndex >>= 1
	}
}

func (bh *BinaryHeap[V]) siftDown(nodeIndex int64) {
	for nodeIndex <= bh.length {
		priorIndex := nodeIndex
		leftIndex := nodeIndex << 1
		rightIndex := (nodeIndex << 1) + 1

		if leftIndex <= bh.length && bh.priorityFunc(bh.nodes[leftIndex].GetValue(), bh.nodes[priorIndex].GetValue()) {
			priorIndex = leftIndex
		}
		if rightIndex <= bh.length && bh.priorityFunc(bh.nodes[rightIndex].GetValue(), bh.nodes[priorIndex].GetValue()) {
			priorIndex = rightIndex
		}

		if priorIndex == nodeIndex {
			return
		}
		bh.swapNodes(nodeIndex, priorIndex)
		nodeIndex = priorIndex
	}
}

func (bh *BinaryHeap[V]) swapNodes(id1, id2 int64) {
	bh.nodes[id1], bh.nodes[id2] = bh.nodes[id2], bh.nodes[id1]
	bh.nodes[id1].index = id1
	bh.nodes[id2].index = id2
}
