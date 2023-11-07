package orderedset

// OrderedSet interface
type OrderedSet[K any] interface {
	Get(key K) (_ K, _ bool)
	GetGreater(key K) (_ K, _ bool)
	GetGreaterThanOrEqual(key K) (_ K, _ bool)
	GetLower(key K) (_ K, _ bool)
	GetLowerThanOrEqual(key K) (_ K, _ bool)

	Max() (_ K, _ bool)
	Min() (_ K, _ bool)
	Len() int64

	ReplaceOrInsert(key K) (_ K, _ bool)

	Delete(key K) (_ K, _ bool)
	DeleteMax() (_ K, _ bool)
	DeleteMin() (_ K, _ bool)
}

// Balanced Binary Search Tree Node interface
type BBSTNode[K any] interface {
	GetLeft() BBSTNode[K]
	GetRight() BBSTNode[K]
	GetParent() BBSTNode[K]
	GetKey() K
}

type compare[K any] func(k1, k2 K) int

func searchNode[K any](node BBSTNode[K], key K, cmp compare[K], sentinel BBSTNode[K]) BBSTNode[K] {
	var curNode BBSTNode[K] = node
	for curNode != sentinel {
		compare := cmp(key, curNode.GetKey())
		if compare == 0 {
			return curNode
		}
		if compare == -1 {
			curNode = curNode.GetLeft()
		} else {
			curNode = curNode.GetRight()
		}
	}
	return nil
}

func searchGreaterNode[K any](node BBSTNode[K], key K, cmp compare[K], sentinel BBSTNode[K]) BBSTNode[K] {
	var greaterNode BBSTNode[K]
	var curNode BBSTNode[K] = node
	for curNode != sentinel {
		compare := cmp(key, curNode.GetKey())
		if compare >= 0 {
			curNode = curNode.GetRight()
			continue
		}
		greaterNode = curNode
		curNode = curNode.GetLeft()
	}
	return greaterNode
}

func searchGreaterThanOrEqualNode[K any](node BBSTNode[K], key K, cmp compare[K], sentinel BBSTNode[K]) BBSTNode[K] {
	var greaterThanOrEqualKeyNode BBSTNode[K]
	var curNode BBSTNode[K] = node
	for curNode != sentinel {
		compare := cmp(key, curNode.GetKey())
		if compare == 0 {
			greaterThanOrEqualKeyNode = curNode
			break
		}
		if compare == 1 {
			curNode = curNode.GetRight()
			continue
		}
		greaterThanOrEqualKeyNode = curNode
		curNode = curNode.GetLeft()
	}
	return greaterThanOrEqualKeyNode
}

func searchLowerNode[K any](node BBSTNode[K], key K, cmp compare[K], sentinel BBSTNode[K]) BBSTNode[K] {
	var lowerNode BBSTNode[K]
	var curNode BBSTNode[K] = node
	for curNode != sentinel {
		compare := cmp(key, curNode.GetKey())
		if compare <= 0 {
			curNode = curNode.GetLeft()
			continue
		}
		lowerNode = curNode
		curNode = curNode.GetRight()
	}
	return lowerNode
}

func searchLowerThanOrEqualNode[K any](node BBSTNode[K], key K, cmp compare[K], sentinel BBSTNode[K]) BBSTNode[K] {
	var lowerThanOrEqualKeyNode BBSTNode[K]
	var curNode BBSTNode[K] = node
	for curNode != sentinel {
		compare := cmp(key, curNode.GetKey())
		if compare == 0 {
			lowerThanOrEqualKeyNode = curNode
			break
		}
		if compare == -1 {
			curNode = curNode.GetLeft()
			continue
		}
		lowerThanOrEqualKeyNode = curNode
		curNode = curNode.GetRight()
	}
	return lowerThanOrEqualKeyNode
}

// panics if node is either nil or sentinel
func getMaxNode[K any](node BBSTNode[K], sentinel BBSTNode[K]) BBSTNode[K] {
	var maxNode BBSTNode[K] = node
	for maxNode.GetRight() != sentinel {
		maxNode = maxNode.GetRight()
	}
	return maxNode
}

// panics if node is either nil or sentinel
func getMinNode[K any](node BBSTNode[K], sentinel BBSTNode[K]) BBSTNode[K] {
	var minNode BBSTNode[K] = node
	for minNode.GetLeft() != sentinel {
		minNode = minNode.GetLeft()
	}
	return minNode
}

type OrderedSetForwardIterator[K any] interface {
	Next() (_ K, _ bool)
	Key() (_ K, _ bool)
	Remove() (_ K, _ bool)
}

type OrderedSetReverseIterator[K any] interface {
	Prev() (_ K, _ bool)
	Key() (_ K, _ bool)
	Remove() (_ K, _ bool)
}

type OrderedSetI[K any] interface {
	OrderedSet[K]
	Begin() OrderedSetForwardIterator[K]
	Rbegin() OrderedSetReverseIterator[K]
}

func Next[K any](node, sentinel BBSTNode[K]) BBSTNode[K] {
	if node.GetRight() != sentinel {
		return getMinNode[K](node.GetRight(), sentinel)
	}
	var next BBSTNode[K] = node
	for next.GetParent() != sentinel && next == next.GetParent().GetRight() {
		next = next.GetParent()
	}
	return next.GetParent()
}

func Prev[K any](node, sentinel BBSTNode[K]) BBSTNode[K] {
	if node.GetLeft() != sentinel {
		return getMaxNode[K](node.GetLeft(), sentinel)
	}
	var prev BBSTNode[K] = node
	for prev.GetParent() != sentinel && prev == prev.GetParent().GetLeft() {
		prev = prev.GetParent()
	}
	return prev.GetParent()
}
