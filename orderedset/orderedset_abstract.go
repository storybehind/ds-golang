package orderedset

// OrderedSet interface
type OrderedSet[K any] interface {
	// Get looks for the key in the set, returning it. It returns (zeroValue, false) if unable to find that key
	Get(key K) (_ K, _ bool)
	// GetGreater looks for smallest key that is strictly greater than key in the set, returning it. It returns (zeroValue, false) if unable to find that key
	GetGreater(key K) (_ K, _ bool)
	// GetGreaterThanOrEqual looks for smallest key that is greater than or equal to key in the set, returning it. It returns (zeroValue, false) if unable to find that key
	GetGreaterThanOrEqual(key K) (_ K, _ bool)
	// GetLower looks for greatest key that is strictly lower than key in the set, returning it. It returns (zeroValue, false) if unable to find that key
	GetLower(key K) (_ K, _ bool)
	// GetLowerThanOrEqual looks for greatest key that is lower than or equal to key in the set, returning it. It returns (zeroValue, false) if unable to find that key
	GetLowerThanOrEqual(key K) (_ K, _ bool)

	// Max returns the largest key in the set, or (zeroValue, false) if the set is empty
	Max() (_ K, _ bool)
	// Min returns the smallest key in the set, or (zeroValue, false) if the set is empty
	Min() (_ K, _ bool)
	// Len returns the number of keys currently in the set.
	Len() int64

	// ReplaceOrInsert adds the given key to the set.
	// If a key in the set already equals the given one, it is removed from the tree and returned, and the second return value is true.
	// Otherwise, (zeroValue, false).
	ReplaceOrInsert(key K) (_ K, _ bool)

	// Delete the key in the set and return its value.
	// If key is not found in the set, returns (zeroValue, false)
	Delete(key K) (_ K, _ bool)
	// Delete the maximum key in the set and return its value.
	// On calling empty set, returns (zeroValue, false)
	DeleteMax() (_ K, _ bool)
	// Delete the minimum key in the set and return its value.
	// On calling empty set, returns (zeroValue, false)
	DeleteMin() (_ K, _ bool)
}

// Balanced Binary Search Tree Node interface
type BBSTNode[K any] interface {
	// Returns left node of Balanced Binary Search Tree Node 
	GetLeft() BBSTNode[K]
	// Returns right node of Balanced Binary Search Tree Node
	GetRight() BBSTNode[K]
	// Returns parent node of Balanced Binary Search Tree Node
	GetParent() BBSTNode[K]
	// Returns key of Balanced Binary Search Tree Node
	GetKey() K
}

// Returns instance of Red-Black Tree.
// Less method determines the order of key.
// k1 precedes k2 if and only if Less(k1, k2) return true.
// k1 equals k2 if and only if !Less(k1, k2) && !Less(k2, k1) holds true.
func New[K any](less func(k1, k2 K) bool) *RbTree[K] {
	return NewRbTree[K](less)
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

// To iterate keys in the ascending order.
type OrderedSetForwardIterator[K any] interface {
	// Calling Next() moves the iterator to the next greater node and returns its key.
	// If Next() is called on last key(or greatest key), it returns (zeroValue, false)
	Next() (_ K, _ bool)
	// Returns the key pointed by iterator. Returns (zeroValue, false) if this is called on empty set or an iterator has completed traversing all the keys.
	Key() (_ K, _ bool)
	// Deletes the key the pointed by iterator, moves the iterator to next greater key.
	// Returns the next greater key if it's present. Otherwise, returns (zeroValue, false).
	// panics on calling Remove() in empty set or an iterator has completed traversing all the keys
	Remove() (_ K, _ bool)
}

// To iterate keys in the descending order
type OrderedSetReverseIterator[K any] interface {
	// Calling Prev() moves the reverse iterator to the next smaller node and returns its key.
	// If Prev() is called on last key (or smallest key), it returns (zeroValue, false)
	Prev() (_ K, _ bool)
	// Returns the key pointed by reverse iterator. Returns (zeroValue, false) if this is called on empty tree or an iterator has completed traversing all the keys
	Key() (_ K, _ bool)
	// Deletes the key the pointed by reverse iterator, moves the reverse iterator to next smaller key.
	// Returns the next smaller key if it's present. Otherwise, returns (zeroValue, false).
	// panics on calling Remove() in empty tree or an iterator has completed traversing all the keys
	Remove() (_ K, _ bool)
}

// OrderedSet Interface with iterator interfaces
type OrderedSetI[K any] interface {
	OrderedSet[K]
	Begin() OrderedSetForwardIterator[K]
	Rbegin() OrderedSetReverseIterator[K]
}

// Returns successor node
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

// Returns predecessor node
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
