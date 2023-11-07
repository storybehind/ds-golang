package orderedset

type rbTreeNode[K any] struct {
	left, right, parent *rbTreeNode[K]
	key                 K
	color               color
}

func (rbTreeNode *rbTreeNode[K]) GetLeft() BBSTNode[K] {
	return rbTreeNode.left
}

func (rbTreeNode *rbTreeNode[K]) GetRight() BBSTNode[K] {
	return rbTreeNode.right
}

func (rbTreeNode *rbTreeNode[K]) GetParent() BBSTNode[K] {
	return rbTreeNode.parent
}

func (rbTreeNode *rbTreeNode[K]) GetKey() K {
	return rbTreeNode.key
}

// color type details color of a node
type color byte

const (
	RED color = iota
	BLACK
)

// Maintains unique set of keys.
// Supports insertion, deletion and search operation in O(log n) time where n is number of keys in the set
type RbTree[K any] struct {
	root     *rbTreeNode[K]
	sentinel *rbTreeNode[K]
	less     func(k1, k2 K) bool
	cmp      compare[K]
	len      int64
}

// Returns instance of Red-Black Tree.
// Less method determines the order of key.
// k1 precedes k2 if and only if Less(k1, k2) return true.
// k1 equals k2 if and only if !Less(k1, k2) && !Less(k2, k1) holds true.
func NewRbTree[K any](less func(k1, k2 K) bool) *RbTree[K] {
	sentinel := &rbTreeNode[K]{
		color: BLACK,
	}
	return &RbTree[K]{
		root:     sentinel,
		sentinel: sentinel,
		less:     less,
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
func (rbTree *RbTree[K]) Get(key K) (_ K, _ bool) {
	var node BBSTNode[K] = searchNode[K](rbTree.root, key, rbTree.cmp, rbTree.sentinel)
	if node != nil {
		return node.GetKey(), true
	}
	return
}

// GetGreater looks for smallest key that is strictly greater than key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (rbTree *RbTree[K]) GetGreater(key K) (_ K, _ bool) {
	var greaterNode BBSTNode[K] = searchGreaterNode[K](rbTree.root, key, rbTree.cmp, rbTree.sentinel)
	if greaterNode != nil {
		return greaterNode.GetKey(), true
	}
	return
}

// GetGreaterThanOrEqual looks for smallest key that is greater than or equal to key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (rbTree *RbTree[K]) GetGreaterThanOrEqual(key K) (_ K, _ bool) {
	var greaterThanOrEqualNode BBSTNode[K] = searchGreaterThanOrEqualNode[K](rbTree.root, key, rbTree.cmp, rbTree.sentinel)
	if greaterThanOrEqualNode != nil {
		return greaterThanOrEqualNode.GetKey(), true
	}
	return
}

// GetLower looks for greatest key that is strictly lower than key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (rbTree *RbTree[K]) GetLower(key K) (_ K, _ bool) {
	var lowerNode BBSTNode[K] = searchLowerNode[K](rbTree.root, key, rbTree.cmp, rbTree.sentinel)
	if lowerNode != nil {
		return lowerNode.GetKey(), true
	}
	return
}

// GetLowerThanOrEqual looks for greatest key that is lower than or equal to key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (rbTree *RbTree[K]) GetLowerThanOrEqual(key K) (_ K, _ bool) {
	var lowerThanOrEqualNode BBSTNode[K] = searchLowerThanOrEqualNode[K](rbTree.root, key, rbTree.cmp, rbTree.sentinel)
	if lowerThanOrEqualNode != nil {
		return lowerThanOrEqualNode.GetKey(), true
	}
	return
}

// Max returns the largest key in the tree, or (zeroValue, false) if the tree is empty
func (rbTree *RbTree[K]) Max() (_ K, _ bool) {
	if rbTree.root == rbTree.sentinel {
		return
	}
	var maxNode BBSTNode[K] = getMaxNode[K](rbTree.root, rbTree.sentinel)
	return maxNode.GetKey(), true
}

// Min returns the smallest key in the tree, or (zeroValue, false) if the tree is empty
func (rbTree *RbTree[K]) Min() (_ K, _ bool) {
	if rbTree.root == rbTree.sentinel {
		return
	}
	var minNode BBSTNode[K] = getMinNode[K](rbTree.root, rbTree.sentinel)
	return minNode.GetKey(), true
}

// Len returns the number of keys currently in the tree.
func (rbtree *RbTree[K]) Len() int64 {
	return rbtree.len
}

// ReplaceOrInsert adds the given key to the tree.
// If a key in the tree already equals the given one, it is removed from the tree and returned, and the second return value is true.
// Otherwise, (zeroValue, false)
// panics if nil is inserted
func (rbTree *RbTree[K]) ReplaceOrInsert(key K) (_ K, _ bool) {
	var y *rbTreeNode[K] = rbTree.sentinel
	x := rbTree.root
	for x != rbTree.sentinel {
		y = x
		if rbTree.less(key, x.key) {
			x = x.left
		} else if rbTree.less(x.key, key) {
			x = x.right
		} else {
			var prevKey K = x.key
			x.key = key
			return prevKey, true
		}
	}
	var z *rbTreeNode[K] = new(rbTreeNode[K])
	z.parent = y
	if y == rbTree.sentinel {
		rbTree.root = z
	} else if rbTree.less(key, y.key) {
		y.left = z
	} else {
		y.right = z
	}
	z.left = rbTree.sentinel
	z.right = rbTree.sentinel
	z.color = RED
	z.key = key
	rbTree.len++
	rbTree.replaceOrInsertFixup(z)
	return
}

func (rbTree *RbTree[K]) replaceOrInsertFixup(z *rbTreeNode[K]) {
	for z.parent.color == RED {
		if z.parent == z.parent.parent.left {
			var y *rbTreeNode[K] = z.parent.parent.right
			if y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
				continue
			} else if z == z.parent.right {
				z = z.parent
				rbTree.leftRotate(z)
			}
			z.parent.color = BLACK
			z.parent.parent.color = RED
			rbTree.rightRotate(z.parent.parent)
		} else {
			var y *rbTreeNode[K] = z.parent.parent.left
			if y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
				continue
			} else if z == z.parent.left {
				z = z.parent
				rbTree.rightRotate(z)
			}
			z.parent.color = BLACK
			z.parent.parent.color = RED
			rbTree.leftRotate(z.parent.parent)
		}
	}
	rbTree.root.color = BLACK
}

func (rbTree *RbTree[K]) leftRotate(x *rbTreeNode[K]) {
	var y *rbTreeNode[K] = x.right
	x.right = y.left
	if y.left != rbTree.sentinel {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == rbTree.sentinel {
		rbTree.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (rbTree *RbTree[K]) rightRotate(y *rbTreeNode[K]) {
	var x *rbTreeNode[K] = y.left
	y.left = x.right
	if x.right != rbTree.sentinel {
		x.right.parent = y
	}
	x.parent = y.parent
	if y.parent == rbTree.sentinel {
		rbTree.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}
	x.right = y
	y.parent = x
}

// Delete the key in the tree and return its value.
// If key is not found in the tree, returns (zeroValue, false)
func (rbTree *RbTree[K]) Delete(key K) (_ K, _ bool) {
	var z BBSTNode[K] = searchNode[K](rbTree.root, key, rbTree.cmp, rbTree.sentinel)
	if z == nil {
		return
	}
	var deletedKey K = z.GetKey()
	rbTree.delete(z.(*rbTreeNode[K]))
	rbTree.len--
	return deletedKey, true
}

// Delete the maximum key in the tree and return its value.
// On calling empty tree, returns (zeroValue, false)
func (rbTree *RbTree[K]) DeleteMax() (_ K, _ bool) {
	if rbTree.root == rbTree.sentinel {
		return
	}
	var z *rbTreeNode[K] = getMaxNode[K](rbTree.root, rbTree.sentinel).(*rbTreeNode[K])
	var deletedKey K = z.key
	rbTree.delete(z)
	rbTree.len--
	return deletedKey, true
}

// Delete the minimum key in the tree and return its value.
// On calling empty tree, returns (zeroValue, false)
func (rbTree *RbTree[K]) DeleteMin() (_ K, _ bool) {
	if rbTree.root == rbTree.sentinel {
		return
	}
	var z *rbTreeNode[K] = getMinNode[K](rbTree.root, rbTree.sentinel).(*rbTreeNode[K])
	var deletedKey K = z.key
	rbTree.delete(z)
	rbTree.len--
	return deletedKey, true
}

func (rbTree *RbTree[K]) transplant(u, v *rbTreeNode[K]) {
	if u.parent == rbTree.sentinel {
		rbTree.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	v.parent = u.parent
}

func (rbTree *RbTree[K]) delete(z *rbTreeNode[K]) {
	var x *rbTreeNode[K]
	var y *rbTreeNode[K] = z
	var yOriginalColor color = y.color
	if z.left == rbTree.sentinel {
		x = z.right
		rbTree.transplant(z, z.right)
	} else if z.right == rbTree.sentinel {
		x = z.left
		rbTree.transplant(z, z.left)
	} else {
		y = getMinNode[K](z.right, rbTree.sentinel).(*rbTreeNode[K])
		yOriginalColor = y.color
		x = y.right
		if y.parent == z {
			x.parent = y
		} else {
			rbTree.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		rbTree.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}
	z.left = nil
	z.right = nil
	z.parent = nil
	if yOriginalColor == BLACK {
		rbTree.deleteFixup(x)
	}
}

func (rbTree *RbTree[K]) deleteFixup(x *rbTreeNode[K]) {
	for x != rbTree.root && x.color == BLACK {
		if x == x.parent.left {
			var w *rbTreeNode[K] = x.parent.right
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				rbTree.leftRotate(x.parent)
				w = x.parent.right
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
				continue
			} else if w.right.color == BLACK {
				w.left.color = BLACK
				w.color = RED
				rbTree.rightRotate(w)
				w = x.parent.right
			}
			w.color = x.parent.color
			x.parent.color = BLACK
			w.right.color = BLACK
			rbTree.leftRotate(x.parent)
			x = rbTree.root
		} else {
			var w *rbTreeNode[K] = x.parent.left
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				rbTree.rightRotate(x.parent)
				w = x.parent.left
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
				continue
			} else if w.left.color == BLACK {
				w.right.color = BLACK
				w.color = RED
				rbTree.leftRotate(w)
				w = x.parent.left
			}
			w.color = x.parent.color
			x.parent.color = BLACK
			w.left.color = BLACK
			rbTree.rightRotate(x.parent)
			x = rbTree.root
		}
	}
	x.color = BLACK
}

type RbIterator[K any] struct {
	next   *rbTreeNode[K]
	rbTree *RbTree[K]
}

// Returns an iterator pointing to least key node in the tree.
// Used to iterate keys in the ascending order.
func (rbTree *RbTree[K]) Begin() OrderedSetForwardIterator[K] {
	var next *rbTreeNode[K] = rbTree.root
	if next != rbTree.sentinel {
		next = getMinNode[K](rbTree.root, rbTree.sentinel).(*rbTreeNode[K])
	}
	return &RbIterator[K]{
		next:   next,
		rbTree: rbTree,
	}
}

// Calling Next() moves the iterator to the next greater node and returns its key
// If Next() is called on last key(or greatest key), it returns (zeroValue, false)
func (rbIterator *RbIterator[K]) Next() (_ K, _ bool) {
	if rbIterator.next == rbIterator.rbTree.sentinel {
		return
	}
	rbIterator.next = Next[K](rbIterator.next, rbIterator.rbTree.sentinel).(*rbTreeNode[K])
	if rbIterator.next == rbIterator.rbTree.sentinel {
		return
	}
	return rbIterator.next.key, true
}

// Returns the key pointed by iterator. Returns (zeroValue, false) if this is called on empty tree or an iterator has completed traversing all the keys
func (rbIterator *RbIterator[K]) Key() (_ K, _ bool) {
	if rbIterator.next != rbIterator.rbTree.sentinel {
		return rbIterator.next.key, true
	}
	return
}

// Deletes the key the pointed by iterator, moves the iterator to next greater key.
// Returns the next greater key if it's present. Otherwise, returns (zeroValue, false)
// panics on calling Remove() in empty tree or an iterator has completed traversing all the keys
func (rbIterator *RbIterator[K]) Remove() (_ K, _ bool) {
	var todelete *rbTreeNode[K] = rbIterator.next
	nextKey, hasNext := rbIterator.Next()
	rbIterator.rbTree.delete(todelete)
	rbIterator.rbTree.len--
	return nextKey, hasNext
}

type ReverseRbIterator[K any] struct {
	prev   *rbTreeNode[K]
	rbTree *RbTree[K]
}

// Returns an reverse iterator pointing to greatest key node in the tree
// Used to iterate keys in the descending order
func (rbTree *RbTree[K]) Rbegin() OrderedSetReverseIterator[K] {
	var prev *rbTreeNode[K] = rbTree.root
	if prev != rbTree.sentinel {
		prev = getMaxNode[K](rbTree.root, rbTree.sentinel).(*rbTreeNode[K])
	}
	return &ReverseRbIterator[K]{
		prev:   prev,
		rbTree: rbTree,
	}
}

// Calling Prev() moves the reverse iterator to the next smaller node and returns its key
// If Prev() is called on last key (or smallest key), it returns (zeroValue, false)
func (reverseRbIterator *ReverseRbIterator[K]) Prev() (_ K, _ bool) {
	var rbTree *RbTree[K] = reverseRbIterator.rbTree
	var prev *rbTreeNode[K] = reverseRbIterator.prev
	if prev == rbTree.sentinel {
		return
	}
	reverseRbIterator.prev = Prev[K](reverseRbIterator.prev, reverseRbIterator.rbTree.sentinel).(*rbTreeNode[K])
	if reverseRbIterator.prev == rbTree.sentinel {
		return
	}
	return reverseRbIterator.prev.key, true
}

// Returns the key pointed by reverse iterator. Returns (zeroValue, false) if this is called on empty tree or an iterator has completed traversing all the keys
func (reverseRbIterator *ReverseRbIterator[K]) Key() (_ K, _ bool) {
	if reverseRbIterator.prev != reverseRbIterator.rbTree.sentinel {
		return reverseRbIterator.prev.key, true
	}
	return
}

// Deletes the key the pointed by reverse iterator, moves the reverse iterator to next smaller key.
// Returns the next smaller key if it's present. Otherwise, returns (zeroValue, false)
// panics on calling Remove() in empty tree or an iterator has completed traversing all the keys
func (reverseRbIterator *ReverseRbIterator[K]) Remove() (_ K, _ bool) {
	var todelete *rbTreeNode[K] = reverseRbIterator.prev
	key, hasPrev := reverseRbIterator.Prev()
	reverseRbIterator.rbTree.delete(todelete)
	reverseRbIterator.rbTree.len--
	return key, hasPrev
}
