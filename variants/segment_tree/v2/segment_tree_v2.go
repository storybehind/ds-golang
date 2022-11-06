package segmenttreev2


type SegmentTreeNode[V any] struct {
	id, leftEnd, rightEnd int64
	Value                 V
}

func (stn SegmentTreeNode[V]) GetLeftEnd() int64 {
	return stn.leftEnd
}

func (stn SegmentTreeNode[V]) GetRightEnd() int64 {
	return stn.rightEnd
}

type MergeFunc[V any] func(cur, left, right *SegmentTreeNode[V])

type PushFunc[V any] func(cur, left, right *SegmentTreeNode[V])

type SegmentTree[V any] struct {
	n          	int64
	nodes      	[]*SegmentTreeNode[V]
	mergeFunc 	MergeFunc[V]
	pushFunc 	*PushFunc[V]
}

func New[V any](n int64, initValues []V, mergeFunc MergeFunc[V], pushFunc *PushFunc[V]) *SegmentTree[V] {
	segTree := &SegmentTree[V]{
		n:          n,
		nodes:      make([]*SegmentTreeNode[V], n << 2),
		mergeFunc: 	mergeFunc,
		pushFunc: 	pushFunc,
	}
	segTree.build(1, 0, n-1, initValues)
	return segTree
}

type MergeQueryFunc[V any] func(v1, v2 V) V

func (segTree *SegmentTree[V]) QueryInterval(leftEnd, rightEnd int64, mergeQueryFunc MergeQueryFunc[V], zeroValue V) V {
	return segTree.queryInterval(1, leftEnd, rightEnd, mergeQueryFunc, zeroValue)
}

func (segTree *SegmentTree[V]) QueryPoint(point int64) V {
	return segTree.queryPoint(1, point)
}

type UpdateQueryFunc[V any] func(cur *SegmentTreeNode[V], updateValue any)

func (segTree *SegmentTree[V]) UpdateInterval(leftEnd, rightEnd int64, updateQueryFunc UpdateQueryFunc[V], updateValue any) {
	segTree.updateInterval(1, leftEnd, rightEnd, updateQueryFunc, updateValue)
}

func (segTree *SegmentTree[V]) build(id, leftEnd, rightEnd int64, initValues []V) {
	segTree.nodes[id] = &SegmentTreeNode[V]{
		id:       id,
		leftEnd:  leftEnd,
		rightEnd: rightEnd,
	}
	if leftEnd == rightEnd {
		segTree.nodes[id].Value = initValues[leftEnd]
		return
	}
	mid := leftEnd + (rightEnd - leftEnd) / 2
	segTree.build(id << 1, leftEnd, mid, initValues)
	segTree.build(id << 1 | 1, mid+1, rightEnd, initValues)
	segTree.mergeFunc(segTree.nodes[id], segTree.nodes[id << 1], segTree.nodes[id << 1 | 1])
}

func (segTree *SegmentTree[V]) queryInterval(id, leftEnd, rightEnd int64, mergeQueryFunc MergeQueryFunc[V], zeroValue V) V {
	segTreeNode := segTree.nodes[id]
	if segTreeNode.leftEnd > rightEnd || segTreeNode.rightEnd < leftEnd {
		return zeroValue
	}
	if segTreeNode.leftEnd >= leftEnd && segTreeNode.rightEnd <= rightEnd {
		return segTreeNode.Value
	}
	if segTree.pushFunc != nil {
		(*segTree.pushFunc)(segTreeNode, segTree.nodes[id << 1], segTree.nodes[id << 1 | 1])
	}
	return mergeQueryFunc(segTree.queryInterval(id << 1, leftEnd, rightEnd, mergeQueryFunc, zeroValue), 
							segTree.queryInterval(id << 1 | 1, leftEnd, rightEnd, mergeQueryFunc, zeroValue))
}

func (segTree *SegmentTree[V]) queryPoint(id, point int64) V {
	segTreeNode := segTree.nodes[id]
	if segTreeNode.leftEnd == segTreeNode.rightEnd {
		return segTreeNode.Value
	}
	if segTree.pushFunc != nil {
		(*segTree.pushFunc)(segTreeNode, segTree.nodes[id << 1], segTree.nodes[id << 1 | 1])
	} 
	mid := segTreeNode.leftEnd + (segTreeNode.rightEnd - segTreeNode.leftEnd) / 2
	if point <= mid {
		return segTree.queryPoint(id << 1, point)
	}
	return segTree.queryPoint(id << 1 | 1, point)
}

func (segTree *SegmentTree[V]) updateInterval(id, leftEnd, rightEnd int64, updateQueryFunc UpdateQueryFunc[V], updateValue any) {
	segTreeNode := segTree.nodes[id]
	if segTreeNode.leftEnd > rightEnd || segTreeNode.rightEnd < leftEnd {
		return
	}
	if segTreeNode.leftEnd >= leftEnd && segTreeNode.rightEnd <= rightEnd {
		updateQueryFunc(segTreeNode, updateValue)
		return
	}
	if segTree.pushFunc != nil {
		(*segTree.pushFunc)(segTreeNode, segTree.nodes[id << 1], segTree.nodes[id << 1 | 1])
	}
	segTree.updateInterval(id << 1, leftEnd, rightEnd, updateQueryFunc, updateValue)
	segTree.updateInterval(id << 1 | 1, leftEnd, rightEnd, updateQueryFunc, updateValue)
	segTree.mergeFunc(segTreeNode, segTree.nodes[id << 1], segTree.nodes[id << 1 | 1])
}

