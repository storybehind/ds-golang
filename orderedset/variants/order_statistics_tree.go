package variants

import "github.com/storybehind/gocontainer/orderedset"

// Maintains unique set of keys.
// Supports insertion, deletion, search, rank and select operations in O(log n) time where n is number of keys in the set
type OrderStatisticsTree[K any] struct {
	*orderedset.RbTreeAugmented[K, int64]
	less func(K, K) bool
	cmp  func(K, K) int
}

// Returns instance of OrderStatisticsTree.
// Less method determines the order of key.
// k1 precedes k2 if and only if Less(k1, k2) return true.
// k1 equals k2 if and only if !Less(k1, k2) && !Less(k2, k1) holds true.
func NewOrderStatisticsTree[K any](less func(k1, k2 K) bool) *OrderStatisticsTree[K] {
	return &OrderStatisticsTree[K]{
		RbTreeAugmented: orderedset.NewRbTreeAugmented[K, int64](less, func(node, sentinel orderedset.BBSTNodeAugmented[K, int64]) int64 {
			return 1 + getSubtreeSize[K](node.GetLeftAugmented(), sentinel) + getSubtreeSize[K](node.GetRightAugmented(), sentinel)
		}),
		less: less,
		cmp: func(k1, k2 K) int {
			if less(k1, k2) {
				return -1
			}
			if less(k2, k1) {
				return 1
			}
			return 0
		},
	}
}

func getSubtreeSize[K any](node, sentinel orderedset.BBSTNodeAugmented[K, int64]) int64 {
	if node == sentinel {
		return 0
	}
	return node.GetAugmentedValue()
}

// rank of key stating from zero.
// Ex: rank of minimum key will be zero.
// Returns -1 if key is not found in the tree
func (ost *OrderStatisticsTree[K]) Rank(key K) int64 {
	rank := int64(0)
	node := ost.GetRoot()
	for node != ost.GetSentinel() {
		cmp := ost.cmp(key, node.GetKey())
		switch cmp {
		case 0:
			rank += getSubtreeSize(node.GetLeftAugmented(), ost.GetSentinel())
			return rank
		case -1:
			node = node.GetLeftAugmented()
		case 1:
			rank += 1 + getSubtreeSize(node.GetLeftAugmented(), ost.GetSentinel())
			node = node.GetRightAugmented()
		}
	}
	return -1
}

//return key element whose rank(key) = r.
// Ex : for r == 0 , return minimum key. If r >= Len(), return zeroValue, false
func (ost *OrderStatisticsTree[K]) Select(r int64) (_ K, _ bool) {
	node := ost.GetRoot()
	for node != ost.GetSentinel() {
		rank := getSubtreeSize(node.GetLeftAugmented(), ost.GetSentinel())
		if rank == r {
			return node.GetKey(), true
		}
		if r < rank {
			node = node.GetLeftAugmented()
			continue
		}
		r -= 1 + rank
		node = node.GetRightAugmented()
	}
	return
}
