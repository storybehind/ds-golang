package segmenttree

type SegmentTreeNode[V ValueType] struct {
	id, leftEnd, rightEnd int64
	value                 V
}

func (stn SegmentTreeNode[V]) GetLeftEnd() int64 {
	return stn.leftEnd
}

func (stn SegmentTreeNode[V]) GetRightEnd() int64 {
	return stn.rightEnd
}

func (stn SegmentTreeNode[V]) GetValue() V {
	return stn.value
}

type ValueType interface {
	pushLazyValue(left, right ValueType)
	update(updateValue any)
}

type MergeFunc[V ValueType] func(left, right *SegmentTreeNode[V]) V

type SegmentTree[V ValueType] struct {
	n          	int64
	nodes      	[]*SegmentTreeNode[V]
	mergeFunc 	MergeFunc[V]
}

func New[V ValueType](n int64, initValues []V, mergeFunc MergeFunc[V]) *SegmentTree[V] {
	segTree := &SegmentTree[V]{
		n:          n,
		nodes:      make([]*SegmentTreeNode[V], n << 2),
		mergeFunc: 	mergeFunc,
	}
	segTree.build(1, 0, n-1, initValues)
	return segTree
}

type MergeQueryFunc[V ValueType] func(v1, v2 V) V

func (segTree *SegmentTree[V]) QueryInterval(leftEnd, rightEnd int64, mergeQueryFunc MergeQueryFunc[V], zeroValue V) V {
	return segTree.queryInterval(1, leftEnd, rightEnd, mergeQueryFunc, zeroValue)
}

func (segTree *SegmentTree[V]) QueryPoint(point int64) V {
	return segTree.queryPoint(1, point)
}

func (segTree *SegmentTree[V]) UpdateInterval(leftEnd, rightEnd int64, updateValue any) {
	segTree.updateInterval(1, leftEnd, rightEnd, updateValue)
}

func (segTree *SegmentTree[V]) build(id, leftEnd, rightEnd int64, initValues []V) {
	segTree.nodes[id] = &SegmentTreeNode[V]{
		id:       id,
		leftEnd:  leftEnd,
		rightEnd: rightEnd,
	}
	if leftEnd == rightEnd {
		segTree.nodes[id].value = initValues[leftEnd]
		return
	}
	mid := leftEnd + (rightEnd - leftEnd) / 2
	segTree.build(id << 1, leftEnd, mid, initValues)
	segTree.build(id << 1 | 1, mid+1, rightEnd, initValues)
	segTree.nodes[id].value = segTree.mergeFunc(segTree.nodes[id << 1], segTree.nodes[id << 1 | 1])
}

func (segTree *SegmentTree[V]) queryInterval(id, leftEnd, rightEnd int64, mergeQueryFunc MergeQueryFunc[V], zeroValue V) V {
	segTreeNode := segTree.nodes[id]
	if segTreeNode.leftEnd > rightEnd || segTreeNode.rightEnd < leftEnd {
		return zeroValue
	}
	if segTreeNode.leftEnd >= leftEnd && segTreeNode.rightEnd <= rightEnd {
		return segTreeNode.value
	}
	segTreeNode.value.pushLazyValue(segTree.nodes[id << 1].value, segTree.nodes[id << 1 | 1].value)
	return mergeQueryFunc(segTree.queryInterval(id << 1, leftEnd, rightEnd, mergeQueryFunc, zeroValue), 
							segTree.queryInterval(id << 1 | 1, leftEnd, rightEnd, mergeQueryFunc, zeroValue))
}

func (segTree *SegmentTree[V]) queryPoint(id, point int64) V {
	segTreeNode := segTree.nodes[id]
	if segTreeNode.leftEnd == segTreeNode.rightEnd {
		return segTreeNode.value
	} 
	segTreeNode.value.pushLazyValue(segTree.nodes[id << 1].value, segTree.nodes[id << 1 | 1].value)
	mid := segTreeNode.leftEnd + (segTreeNode.rightEnd - segTreeNode.leftEnd) / 2
	if point <= mid {
		return segTree.queryPoint(id << 1, point)
	}
	return segTree.queryPoint(id << 1 | 1, point)
}

func (segTree *SegmentTree[V]) updateInterval(id, leftEnd, rightEnd int64, updateValue any) {
	segTreeNode := segTree.nodes[id]
	if segTreeNode.leftEnd > rightEnd || segTreeNode.rightEnd < leftEnd {
		return
	}
	if segTreeNode.leftEnd >= leftEnd && segTreeNode.rightEnd <= rightEnd {
		segTreeNode.value.update(updateValue)
		return
	}
	segTreeNode.value.pushLazyValue(segTree.nodes[id << 1].value, segTree.nodes[id << 1 | 1].value)
	segTree.updateInterval(id << 1, leftEnd, rightEnd, updateValue)
	segTree.updateInterval(id << 1 | 1, leftEnd, rightEnd, updateValue)
	segTreeNode.value = segTree.mergeFunc(segTree.nodes[id << 1], segTree.nodes[id << 1 | 1])
}

