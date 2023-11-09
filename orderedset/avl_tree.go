package orderedset

//Node of avl tree that holds a particular key
//Maintain left, right and parent pointer for tree traversal
type avlTreeNode[K any] struct {
	left, right, parent *avlTreeNode[K]
	key                 K
	height              int64
}

//Get left node
func (node *avlTreeNode[K]) GetLeft() BBSTNode[K] {
	return node.left
}

//Get right node
func (node *avlTreeNode[K]) GetRight() BBSTNode[K] {
	return node.right
}

//Get parent node
func (node *avlTreeNode[K]) GetParent() BBSTNode[K] {
	return node.parent
}

//Get key
func (node *avlTreeNode[K]) GetKey() K {
	return node.key
}

// Maintains unique set of keys.
// Supports insertion, deletion and search operation in O(log n) time where n is number of keys in the set
type AvlTree[K any] struct {
	root     *avlTreeNode[K]
	sentinel *avlTreeNode[K]
	less     func(k1, k2 K) bool
	cmp      compare[K]
	len      int64
}

// Returns instance of AvlTree.
// Less method determines the order of key.
// k1 precedes k2 in AvlTree if and only if Less(k1, k2) return true.
// k1 equals k2 if and only if !Less(k1, k2) && !Less(k2, k1) holds true.
func NewAvlTree[K any](less func(k1, k2 K) bool) *AvlTree[K] {
	sentinel := &avlTreeNode[K]{
		height: 0,
	}
	return &AvlTree[K]{
		root:     sentinel,
		less:     less,
		sentinel: sentinel,
		cmp: func(k1, k2 K) int {
			if less(k1, k2) {
				return -1
			}
			if less(k2, k1) {
				return 1
			}
			return 0
		},
		len: 0,
	}
}

// Get looks for the key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (avlTree *AvlTree[K]) Get(key K) (_ K, _ bool) {
	var node BBSTNode[K] = searchNode[K](avlTree.root, key, avlTree.cmp, avlTree.sentinel)
	if node != nil {
		return node.GetKey(), true
	}
	return
}

// GetGreater looks for smallest key that is strictly greater than key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (avlTree *AvlTree[K]) GetGreater(key K) (_ K, _ bool) {
	var greaterNode BBSTNode[K] = searchGreaterNode[K](avlTree.root, key, avlTree.cmp, avlTree.sentinel)
	if greaterNode != nil {
		return greaterNode.GetKey(), true
	}
	return
}

// GetGreaterThanOrEqual looks for smallest key that is greater than or equal to key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (avlTree *AvlTree[K]) GetGreaterThanOrEqual(key K) (_ K, _ bool) {
	var greaterThanOrEqualNode BBSTNode[K] = searchGreaterThanOrEqualNode[K](avlTree.root, key, avlTree.cmp, avlTree.sentinel)
	if greaterThanOrEqualNode != nil {
		return greaterThanOrEqualNode.GetKey(), true
	}
	return
}

// GetLower looks for greatest key that is strictly lower than key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (avlTree *AvlTree[K]) GetLower(key K) (_ K, _ bool) {
	var lowerNode BBSTNode[K] = searchLowerNode[K](avlTree.root, key, avlTree.cmp, avlTree.sentinel)
	if lowerNode != nil {
		return lowerNode.GetKey(), true
	}
	return
}

// GetLowerThanOrEqual looks for greatest key that is lower than or equal to key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (avlTree *AvlTree[K]) GetLowerThanOrEqual(key K) (_ K, _ bool) {
	var lowerThanOrEqualNode BBSTNode[K] = searchLowerThanOrEqualNode[K](avlTree.root, key, avlTree.cmp, avlTree.sentinel)
	if lowerThanOrEqualNode != nil {
		return lowerThanOrEqualNode.GetKey(), true
	}
	return
}

// Max returns the largest key in the tree, or (zeroValue, false) if the tree is empty
func (avlTree *AvlTree[K]) Max() (_ K, _ bool) {
	if avlTree.root == avlTree.sentinel {
		return
	}
	var maxNode BBSTNode[K] = getMaxNode[K](avlTree.root, avlTree.sentinel)
	return maxNode.GetKey(), true
}

// Min returns the smallest key in the tree, or (zeroValue, false) if the tree is empty
func (avlTree *AvlTree[K]) Min() (_ K, _ bool) {
	if avlTree.root == avlTree.sentinel {
		return
	}
	var minNode BBSTNode[K] = getMinNode[K](avlTree.root, avlTree.sentinel)
	return minNode.GetKey(), true
}

// Len returns the number of keys currently in the tree.
func (avlTree *AvlTree[K]) Len() int64 {
	return avlTree.len
}

// toRemove details what item to remove in a node.remove call.
type toRemove int

const (
	removeKey toRemove = iota // removes the given item
	removeMin                 // removes smallest item in the subtree
	removeMax                 // removes largest item in the subtree
)

// ReplaceOrInsert adds the given key to the tree.
// If a key in the tree already equals the given one, it is removed from the tree and returned, and the second return value is true.
// Otherwise, (zeroValue, false)
func (avlTree *AvlTree[K]) ReplaceOrInsert(key K) (_ K, _ bool) {
	var prevKey K
	var has bool
	avlTree.root, prevKey, has = avlTree.replaceOrInsert(avlTree.root, key)
	avlTree.root.parent = avlTree.sentinel
	if !has {
		avlTree.len++
	}
	return prevKey, has
}

// Delete removes a key equal to the passed in key from the tree, returning it. If no such key exists, returns (zeroValue, false)
func (avlTree *AvlTree[K]) Delete(key K) (K, bool) {
	var deletedKey K
	var deleted bool
	avlTree.root, deletedKey, deleted = avlTree.delete(avlTree.root, key, removeKey)
	avlTree.root.parent = avlTree.sentinel
	if deleted {
		avlTree.len--
	}
	return deletedKey, deleted
}

// DeleteMax removes the largest key in the tree and returns it. If no such item exists, returns (zeroValue, false)
func (avlTree *AvlTree[K]) DeleteMax() (K, bool) {
	var deletedKey K
	var deleted bool
	var zero K
	avlTree.root, deletedKey, deleted = avlTree.delete(avlTree.root, zero, removeMax)
	avlTree.root.parent = avlTree.sentinel
	if deleted {
		avlTree.len--
	}
	return deletedKey, deleted
}

// DeleteMin removes the smallest key in the tree and returns it. If no such item exists, returns (zeroValue, false)
func (avlTree *AvlTree[K]) DeleteMin() (K, bool) {
	var deletedKey K
	var deleted bool
	var zero K
	avlTree.root, deletedKey, deleted = avlTree.delete(avlTree.root, zero, removeMin)
	avlTree.root.parent = avlTree.sentinel
	if deleted {
		avlTree.len--
	}
	return deletedKey, deleted
}

func (avlTree *AvlTree[K]) replaceOrInsert(node *avlTreeNode[K], key K) (_ *avlTreeNode[K], _ K, _ bool) {
	if node == avlTree.sentinel {
		var zero K
		newNode := &avlTreeNode[K]{
			left:   avlTree.sentinel,
			right:  avlTree.sentinel,
			key:    key,
			height: 1,
		}
		return newNode, zero, false
	}
	compare := avlTree.cmp(key, node.key)
	var prevKey K
	var has bool
	switch compare {
	case 0:
		prevKey = node.key
		node.key = key
		has = true
		return node, prevKey, has
	case -1:
		node.left, prevKey, has = avlTree.replaceOrInsert(node.left, key)
		node.left.parent = node
	case 1:
		node.right, prevKey, has = avlTree.replaceOrInsert(node.right, key)
		node.right.parent = node
	}
	node.height = 1 + max(node.left.height, node.right.height)
	return avlTree.balanceNode(node), prevKey, has
}

func (avlTree *AvlTree[K]) balanceNode(node *avlTreeNode[K]) *avlTreeNode[K] {
	nodeHeightDiff := getHeightDiff(node)
	if nodeHeightDiff > 1 {
		rightNodeHeightDiff := getHeightDiff(node.right)
		if rightNodeHeightDiff < 0 {
			// node.right is left-heavy
			node.right = avlTree.rightRotate(node.right)
			node.right.parent = node
		}
		return avlTree.leftRotate(node)
	}
	if nodeHeightDiff < -1 {
		leftNodeHeightDiff := getHeightDiff(node.left)
		if leftNodeHeightDiff > 0 {
			node.left = avlTree.leftRotate(node.left)
			node.left.parent = node
		}
		return avlTree.rightRotate(node)
	}
	return node
}

func (avlTree AvlTree[K]) delete(node *avlTreeNode[K], key K, typ toRemove) (_ *avlTreeNode[K], _ K, _ bool) {
	if node == avlTree.sentinel {
		var zero K
		return node, zero, false
	}

	var deletedKey K
	var deleted bool
	switch typ {
	case removeMin:
		if node.left == avlTree.sentinel {
			return node.right, node.key, true
		}
		node.left, deletedKey, deleted = avlTree.delete(node.left, key, typ)
		node.left.parent = node
	case removeMax:
		if node.right == avlTree.sentinel {
			return node.left, node.key, true
		}
		node.right, deletedKey, deleted = avlTree.delete(node.right, key, typ)
		node.right.parent = node
	case removeKey:
		compare := avlTree.cmp(key, node.key)
		switch compare {
		case 0:
			deletedKey = node.key
			deleted = true
			if node.left != avlTree.sentinel {
				var zero K
				leftMaxNode := getMaxNode[K](node.left, avlTree.sentinel).(*avlTreeNode[K])
				leftMaxNode.left, _, _ = avlTree.delete(node.left, zero, removeMax)
				leftMaxNode.right = node.right
				node = leftMaxNode
			} else if node.right != avlTree.sentinel {
				var zero K
				rightMinNode := getMinNode[K](node.right, avlTree.sentinel).(*avlTreeNode[K])
				rightMinNode.right, _, _ = avlTree.delete(node.right, zero, removeMin)
				rightMinNode.left = node.left
				node = rightMinNode
			} else {
				return avlTree.sentinel, node.key, true
			}
		case -1:
			node.left, deletedKey, deleted = avlTree.delete(node.left, key, typ)
			node.left.parent = node
		case 1:
			node.right, deletedKey, deleted = avlTree.delete(node.right, key, typ)
			node.right.parent = node
		}
	default:
		panic("invalid remove type")
	}
	node.height = 1 + max(node.left.height, node.right.height)
	return avlTree.balanceNode(node), deletedKey, deleted
}

func getHeightDiff[K any](node *avlTreeNode[K]) int64 {
	return node.right.height - node.left.height
}

func max(a, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}

func (avlTree *AvlTree[K]) rightRotate(y *avlTreeNode[K]) *avlTreeNode[K] {
	x := y.left
	y.left = x.right
	y.left.parent = y

	x.right = y
	y.parent = x

	y.height = 1 + max(y.left.height, y.right.height)
	x.height = 1 + max(x.left.height, x.right.height)

	return x
}

func (avlTree *AvlTree[K]) leftRotate(x *avlTreeNode[K]) *avlTreeNode[K] {
	y := x.right
	x.right = y.left
	x.right.parent = x

	y.left = x
	x.parent = y

	x.height = 1 + max(x.left.height, x.right.height)
	y.height = 1 + max(y.left.height, y.right.height)

	return y
}

type AvlIterator[K any] struct {
	next    *avlTreeNode[K]
	avlTree *AvlTree[K]
}

// Returns an iterator pointing to smallest key node in the tree or to sentinel node if tree is empty.
// Used to iterate keys in the ascending order.
func (avlTree *AvlTree[K]) Begin() OrderedSetForwardIterator[K] {
	var next *avlTreeNode[K] = avlTree.root
	if next != avlTree.sentinel {
		next = getMinNode[K](avlTree.root, avlTree.sentinel).(*avlTreeNode[K])
	}
	return &AvlIterator[K]{
		next:    next,
		avlTree: avlTree,
	}
}

// Calling Next() moves the iterator to the next greater node and returns its key.
// If Next() is called on last key(or greatest key), it returns (zeroValue, false)
func (avlIterator *AvlIterator[K]) Next() (_ K, _ bool) {
	if avlIterator.next == avlIterator.avlTree.sentinel {
		return
	}
	avlIterator.next = Next[K](avlIterator.next, avlIterator.avlTree.sentinel).(*avlTreeNode[K])
	if avlIterator.next == avlIterator.avlTree.sentinel {
		return
	}
	return avlIterator.next.key, true
}

// Returns the key pointed by iterator. Returns (zeroValue, false) if this is called on empty tree or an iterator has completed traversing all the keys
func (avlIterator *AvlIterator[K]) Key() (_ K, _ bool) {
	if avlIterator.next != avlIterator.avlTree.sentinel {
		return avlIterator.next.key, true
	}
	return
}

// Deletes the key the pointed by iterator, moves the iterator to next greater key.
// Returns the next greater key if it's present. Otherwise, returns (zeroValue, false)
// panics on calling Remove() in empty tree or an iterator has completed traversing all the keys
func (avlIterator *AvlIterator[K]) Remove() (_ K, _ bool) {
	var todelete *avlTreeNode[K] = avlIterator.next
	var avlTree *AvlTree[K] = avlIterator.avlTree
	nextKey, hasNext := avlIterator.Next()
	avlTree.Delete(todelete.key)
	return nextKey, hasNext
}

type ReverseAvlIterator[K any] struct {
	prev    *avlTreeNode[K]
	avlTree *AvlTree[K]
}

// Returns an reverse iterator pointing to greatest key node in the tree or to sentinel node if tree is empty.
// Used to iterate keys in the descending order
func (avlTree *AvlTree[K]) Rbegin() OrderedSetReverseIterator[K] {
	var prev *avlTreeNode[K] = avlTree.root
	if prev != avlTree.sentinel {
		prev = getMaxNode[K](avlTree.root, avlTree.sentinel).(*avlTreeNode[K])
	}
	return &ReverseAvlIterator[K]{
		prev:    prev,
		avlTree: avlTree,
	}
}

// Calling Prev() moves the reverse iterator to the next smaller node and returns its key.
// If Prev() is called on last key (or smallest key), it returns (zeroValue, false)
func (reverseAvlIterator *ReverseAvlIterator[K]) Prev() (_ K, _ bool) {
	if reverseAvlIterator.prev == reverseAvlIterator.avlTree.sentinel {
		return
	}
	reverseAvlIterator.prev = Prev[K](reverseAvlIterator.prev, reverseAvlIterator.avlTree.sentinel).(*avlTreeNode[K])
	if reverseAvlIterator.prev == reverseAvlIterator.avlTree.sentinel {
		return
	}
	return reverseAvlIterator.prev.key, true
}

// Returns the key pointed by reverse iterator. Returns (zeroValue, false) if this is called on empty tree or an iterator has completed traversing all the keys
func (reverseAvlIterator *ReverseAvlIterator[K]) Key() (_ K, _ bool) {
	if reverseAvlIterator.prev != reverseAvlIterator.avlTree.sentinel {
		return reverseAvlIterator.prev.key, true
	}
	return
}

// Deletes the key the pointed by reverse iterator, moves the reverse iterator to next smaller key.
// Returns the next smaller key if it's present. Otherwise, returns (zeroValue, false).
// panics on calling Remove() in empty tree or an iterator has completed traversing all the keys
func (reverseAvlIterator *ReverseAvlIterator[K]) Remove() (_ K, _ bool) {
	var todelete *avlTreeNode[K] = reverseAvlIterator.prev
	key, hasPrev := reverseAvlIterator.Prev()
	reverseAvlIterator.avlTree.Delete(todelete.key)
	return key, hasPrev
}
