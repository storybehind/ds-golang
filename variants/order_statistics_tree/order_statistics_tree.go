package orderstatisticstree

import "ds/bbst"

type OrderStatisticsTree[K any] struct {
	bbst.BBST[K]
	less bbst.Less[K]
}

type subtreeSize int64

func getSubtreeSize[K any](node bbst.BBSTNode[K]) subtreeSize {
	if node != nil {
		return node.GetAugmentedData().(subtreeSize)
	}
	return subtreeSize(0)
}

func New[K any](less bbst.Less[K], tag bbst.ConcreteTag) *OrderStatisticsTree[K] {
	updateAugmentedData := func(node bbst.BBSTNode[K]) {
		node.SetAugmentedData(1 + getSubtreeSize(node.GetLeft()) + getSubtreeSize(node.GetRight()))
	}
	switch tag {
		case bbst.AvlTreeTag:
			avlTree := bbst.NewAvlTreeByLessAndUpdateAugmentedData(less, (*bbst.UpdateAugmentedData[K])(&updateAugmentedData))
			t := &OrderStatisticsTree[K]{
				less: less,
			}
			t.BBST = avlTree
			return t
		default:
			panic("unsupported bbst concrete tag")
	}
}

// rank of key stating from zero. 
// Ex: rank of minimum element will be zero. 
// Returns -1 if key is not found in the tree
func (ost *OrderStatisticsTree[K]) Rank(key K) int64 {
	rank := int64(0)
	node := ost.GetRoot()
	for node != nil {
		cmp := ost.compare(key, node.GetKey())
		switch cmp {
			case 0:
				rank += int64(getSubtreeSize(node.GetLeft()))
				return rank
			case -1:
				node = node.GetLeft()
			case 1:
				rank += 1 + int64(getSubtreeSize(node.GetLeft()))
				node = node.GetRight()
		}
	}
	return -1
}

//return key element whose rank(key) = r. 
// Ex : for r == 0 , return minimum key. If r >= ost.Len(), return zeroValue, false
func (ost *OrderStatisticsTree[K]) Select(r int64) (_ K, _ bool) {
	node := ost.GetRoot()
	for node != nil {
		rank := int64(getSubtreeSize(node.GetLeft()))
		if rank == r {
			return node.GetKey(), true
		}
		if r < rank {
			node = node.GetLeft()
			continue
		}
		r -= 1 + rank
		node = node.GetRight()
	}
	return
}

func (ost *OrderStatisticsTree[K]) compare(k1, k2 K) int {
	if ost.less(k1, k2) {
		return -1
	}
	if ost.less(k2, k1) {
		return 1
	}
	return 0
}

