package orderedset

// Balanced Binary Search Node interface with support for augmentation
type BBSTNodeAugmented[K, A any] interface {
	BBSTNode[K]
	// Returns left augmented node
	GetLeftAugmented() BBSTNodeAugmented[K, A]
	// Returns right augmented node
	GetRightAugmented() BBSTNodeAugmented[K, A]
	// Returns parent augmented node
	GetParentAugmented() BBSTNodeAugmented[K, A]
	// Returns node's augmented value
	GetAugmentedValue() A
}

type rbTreeNodeAugmented[K, A any] struct {
	left, right, parent *rbTreeNodeAugmented[K, A]
	color color
	key K
	augmentedValue A
}

func (rbTreeNodeAugmented *rbTreeNodeAugmented[K, A]) GetLeft() BBSTNode[K] {
	return rbTreeNodeAugmented.left
}

func (rbTreeNodeAugmented *rbTreeNodeAugmented[K, A]) GetRight() BBSTNode[K] {
	return rbTreeNodeAugmented.right
}

func (rbTreeNodeAugmented *rbTreeNodeAugmented[K, A]) GetParent() BBSTNode[K] {
	return rbTreeNodeAugmented.parent
}

func (rbTreeNodeAugmented *rbTreeNodeAugmented[K, A]) GetLeftAugmented() BBSTNodeAugmented[K, A] {
	return rbTreeNodeAugmented.left
}

func (rbTreeNodeAugmented *rbTreeNodeAugmented[K, A]) GetRightAugmented() BBSTNodeAugmented[K, A] {
	return rbTreeNodeAugmented.right
}

func (rbTreeNodeAugmented *rbTreeNodeAugmented[K, A]) GetParentAugmented() BBSTNodeAugmented[K, A] {
	return rbTreeNodeAugmented.parent
}

func (rbTreeNodeAugmented *rbTreeNodeAugmented[K, A]) GetKey() K {
	return rbTreeNodeAugmented.key
}

func (rbTreeNodeAugmented *rbTreeNodeAugmented[K, A]) GetAugmentedValue() A {
	return rbTreeNodeAugmented.augmentedValue
}

// Maintains unique set of keys and invariant of node's augmented value. 
// Supports insertion, deletion of keys in O(t * log n) time where n is number of keys in the set and t is time required to maintain node's invariant.
// Search operation takes O(log n) time.
// Can be embedded to support additional functionalities
type RbTreeAugmented[K, A any] struct {
	root     *rbTreeNodeAugmented[K, A]
	sentinel *rbTreeNodeAugmented[K, A]
	updateAugmentValue func(BBSTNodeAugmented[K, A], BBSTNodeAugmented[K, A]) A
	less     func(k1, k2 K) bool
	cmp      compare[K]
	len      int64
}

// Returns instance of Red-Black Tree.
// Less method determines the order of key.
// k1 precedes k2 if and only if Less(k1, k2) return true.
// k1 equals k2 if and only if !Less(k1, k2) && !Less(k2, k1) holds true.

// updateAugmentValue maintains invariants of node's augmented value.
// First argument gives node pointer whose invariant has to be maintained.
// Second argument gives sentinel node pointer (can be thought of nil leaf nodes or root's parent)
func NewRbTreeAugmented[K, A any](less func(k1, k2 K) bool, updateAugmentValue func(BBSTNodeAugmented[K, A], BBSTNodeAugmented[K, A]) A) *RbTreeAugmented[K, A] {
	sentinel := &rbTreeNodeAugmented[K, A] {
		color: BLACK,
	}
	return &RbTreeAugmented[K, A] {
		root:     sentinel,
		sentinel: sentinel,
		less:     less,
		updateAugmentValue: updateAugmentValue,
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
func (rbTreeAugmented *RbTreeAugmented[K, A]) Get(key K) (_ K, _ bool) {
	var node BBSTNode[K] = searchNode[K](rbTreeAugmented.root, key, rbTreeAugmented.cmp, rbTreeAugmented.sentinel)
	if node != nil {
		return node.GetKey(), true
	}
	return
}

// GetGreater looks for smallest key that is strictly greater than key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (rbTreeAugmented *RbTreeAugmented[K, A]) GetGreater(key K) (_ K, _ bool) {
	var greaterNode BBSTNode[K] = searchGreaterNode[K](rbTreeAugmented.root, key, rbTreeAugmented.cmp, rbTreeAugmented.sentinel)
	if greaterNode != nil {
		return greaterNode.GetKey(), true
	}
	return
}

// GetGreaterThanOrEqual looks for smallest key that is greater than or equal to key in the tree, returning it. It returns (zeroValue, zeroValue, false) if unable to find that key
func (rbTreeAugmented *RbTreeAugmented[K, A]) GetGreaterThanOrEqual(key K) (_ K, _ bool) {
	var greaterThanOrEqualNode BBSTNode[K] = searchGreaterThanOrEqualNode[K](rbTreeAugmented.root, key, rbTreeAugmented.cmp, rbTreeAugmented.sentinel)
	if greaterThanOrEqualNode != nil {
		return greaterThanOrEqualNode.GetKey(), true
	}
	return
}

// GetLower looks for greatest key that is strictly lower than key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (rbTreeAugmented *RbTreeAugmented[K, A]) GetLower(key K) (_ K, _ bool) {
	var lowerNode BBSTNode[K] = searchLowerNode[K](rbTreeAugmented.root, key, rbTreeAugmented.cmp, rbTreeAugmented.sentinel)
	if lowerNode != nil {
		return lowerNode.GetKey(), true
	}
	return
}

// GetLowerThanOrEqual looks for greatest key that is lower than or equal to key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (rbTreeAugmented *RbTreeAugmented[K, A]) GetLowerThanOrEqual(key K) (_ K, _ bool) {
	var lowerThanOrEqualNode BBSTNode[K] = searchLowerThanOrEqualNode[K](rbTreeAugmented.root, key, rbTreeAugmented.cmp, rbTreeAugmented.sentinel)
	if lowerThanOrEqualNode != nil {
		return lowerThanOrEqualNode.GetKey(), true
	}
	return
}

// Max returns the largest key in the tree, or (zeroValue, false) if the tree is empty
func (rbTreeAugmented *RbTreeAugmented[K, A]) Max() (_ K, _ bool) {
	if rbTreeAugmented.root == rbTreeAugmented.sentinel {
		return
	}
	var maxNode BBSTNode[K] = getMaxNode[K](rbTreeAugmented.root, rbTreeAugmented.sentinel)
	return maxNode.GetKey(), true
}

// Min returns the smallest key in the tree, or (zeroValue, false) if the tree is empty
func (rbTreeAugmented *RbTreeAugmented[K, A]) Min() (_ K, _ bool) {
	if rbTreeAugmented.root == rbTreeAugmented.sentinel {
		return
	}
	var minNode BBSTNode[K] = getMinNode[K](rbTreeAugmented.root, rbTreeAugmented.sentinel)
	return minNode.GetKey(), true
}

// Len returns the number of keys currently in the tree.
func (rbTreeAugmented *RbTreeAugmented[K, A]) Len() int64 {
	return rbTreeAugmented.len
}

// Returns root node of the tree. 
func (rbTreeAugmented *RbTreeAugmented[K, A]) GetRoot() BBSTNodeAugmented[K, A] {
	return rbTreeAugmented.root
}

// Returns sentinel node of the tree (can be thought of nil leaf nodes or root's parent)
func (rbTreeAugmented *RbTreeAugmented[K, A]) GetSentinel() BBSTNodeAugmented[K, A] {
	return rbTreeAugmented.sentinel
}

// ReplaceOrInsert adds the given key to the tree.
// If a key in the tree already equals the given one, it is removed from the tree and returned, and the second return value is true.
// Otherwise, (zeroValue, false)
// panics if nil is inserted
func (rbTreeAugmented *RbTreeAugmented[K, A]) ReplaceOrInsert(key K) (_ K, _ bool) {
	var y *rbTreeNodeAugmented[K, A] = rbTreeAugmented.sentinel
	x := rbTreeAugmented.root
	for x != rbTreeAugmented.sentinel {
		y = x
		if rbTreeAugmented.less(key, x.key) {
			x = x.left
		} else if rbTreeAugmented.less(x.key, key) {
			x = x.right
		} else {
			var prevKey K = x.key
			x.key = key
			return prevKey, true
		}
	}
	var z *rbTreeNodeAugmented[K, A] = new(rbTreeNodeAugmented[K, A])
	z.parent = y
	if y == rbTreeAugmented.sentinel {
		rbTreeAugmented.root = z
	} else if rbTreeAugmented.less(key, y.key) {
		y.left = z
	} else {
		y.right = z
	}
	z.left = rbTreeAugmented.sentinel
	z.right = rbTreeAugmented.sentinel
	z.color = RED
	z.key = key
	rbTreeAugmented.len++
	rbTreeAugmented.replaceOrInsertFixup(z)
	return
}

func (rbTreeAugmented *RbTreeAugmented[K, A]) replaceOrInsertFixup(z *rbTreeNodeAugmented[K, A]) {
	for z.parent.color == RED {
		if z.parent == z.parent.parent.left {
			var y *rbTreeNodeAugmented[K, A] = z.parent.parent.right
			if y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
				continue
			} else if z == z.parent.right {
				z = z.parent
				rbTreeAugmented.leftRotate(z)
			}
			z.parent.color = BLACK
			z.parent.parent.color = RED
			rbTreeAugmented.rightRotate(z.parent.parent)
		} else {
			var y *rbTreeNodeAugmented[K, A] = z.parent.parent.left
			if y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
				continue
			} else if z == z.parent.left {
				z = z.parent
				rbTreeAugmented.rightRotate(z)
			}
			z.parent.color = BLACK
			z.parent.parent.color = RED
			rbTreeAugmented.leftRotate(z.parent.parent)
		}
	}
	rbTreeAugmented.root.color = BLACK
}

func (rbTreeAugmented *RbTreeAugmented[K, A]) leftRotate(x *rbTreeNodeAugmented[K, A]) {
	var y *rbTreeNodeAugmented[K, A] = x.right
	x.right = y.left
	if y.left != rbTreeAugmented.sentinel {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == rbTreeAugmented.sentinel {
		rbTreeAugmented.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
	x.augmentedValue = rbTreeAugmented.updateAugmentValue(x, rbTreeAugmented.sentinel)
	y.augmentedValue = rbTreeAugmented.updateAugmentValue(y, rbTreeAugmented.sentinel)
}

func (rbTreeAugmented *RbTreeAugmented[K, A]) rightRotate(y *rbTreeNodeAugmented[K, A]) {
	var x *rbTreeNodeAugmented[K, A] = y.left
	y.left = x.right
	if x.right != rbTreeAugmented.sentinel {
		x.right.parent = y
	}
	x.parent = y.parent
	if y.parent == rbTreeAugmented.sentinel {
		rbTreeAugmented.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}
	x.right = y
	y.parent = x
	y.augmentedValue = rbTreeAugmented.updateAugmentValue(y, rbTreeAugmented.sentinel)
	x.augmentedValue = rbTreeAugmented.updateAugmentValue(x, rbTreeAugmented.sentinel)
}

// Delete the key in the tree and return its value.
// If key is not found in the tree, returns (zeroValue, false)
func (rbTreeAugmented *RbTreeAugmented[K, A]) Delete(key K) (_ K, _ bool) {
	var z BBSTNode[K] = searchNode[K](rbTreeAugmented.root, key, rbTreeAugmented.cmp, rbTreeAugmented.sentinel)
	if z == nil {
		return
	}
	var deletedKey K = z.GetKey()
	rbTreeAugmented.delete(z.(*rbTreeNodeAugmented[K, A]))
	rbTreeAugmented.len--
	return deletedKey, true
}

// Delete the maximum key in the tree and return its value.
// On calling empty tree, returns (zeroValue, false)
func (rbTreeAugmented *RbTreeAugmented[K, A]) DeleteMax() (_ K, _ bool) {
	if rbTreeAugmented.root == rbTreeAugmented.sentinel {
		return
	}
	var z *rbTreeNodeAugmented[K, A] = getMaxNode[K](rbTreeAugmented.root, rbTreeAugmented.sentinel).(*rbTreeNodeAugmented[K, A])
	var deletedKey K = z.key
	rbTreeAugmented.delete(z)
	rbTreeAugmented.len--
	return deletedKey, true
}

// Delete the minimum key in the tree and return its value.
// On calling empty tree, returns (zeroValue, false)
func (rbTreeAugmented *RbTreeAugmented[K, A]) DeleteMin() (_ K, _ bool) {
	if rbTreeAugmented.root == rbTreeAugmented.sentinel {
		return
	}
	var z *rbTreeNodeAugmented[K, A] = getMinNode[K](rbTreeAugmented.root, rbTreeAugmented.sentinel).(*rbTreeNodeAugmented[K, A])
	var deletedKey K = z.key
	rbTreeAugmented.delete(z)
	rbTreeAugmented.len--
	return deletedKey, true
}

func (rbTreeAugmented *RbTreeAugmented[K, A]) transplant(u, v *rbTreeNodeAugmented[K, A]) {
	if u.parent == rbTreeAugmented.sentinel {
		rbTreeAugmented.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	v.parent = u.parent
}

func (rbTreeAugmented *RbTreeAugmented[K, A]) delete(z *rbTreeNodeAugmented[K, A]) {
	var x *rbTreeNodeAugmented[K, A]
	var y *rbTreeNodeAugmented[K, A] = z
	var yOriginalColor color = y.color
	if z.left == rbTreeAugmented.sentinel {
		x = z.right
		rbTreeAugmented.transplant(z, z.right)
	} else if z.right == rbTreeAugmented.sentinel {
		x = z.left
		rbTreeAugmented.transplant(z, z.left)
	} else {
		y = getMinNode[K](z.right, rbTreeAugmented.sentinel).(*rbTreeNodeAugmented[K, A])
		yOriginalColor = y.color
		x = y.right
		if y.parent == z {
			x.parent = y
		} else {
			rbTreeAugmented.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		rbTreeAugmented.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}
	z.left = nil
	z.right = nil
	z.parent = nil
	if yOriginalColor == BLACK {
		rbTreeAugmented.deleteFixup(x)
	}
}

func (rbTreeAugmented *RbTreeAugmented[K, A]) deleteFixup(x *rbTreeNodeAugmented[K, A]) {
	for x != rbTreeAugmented.root && x.color == BLACK {
		if x == x.parent.left {
			var w *rbTreeNodeAugmented[K, A] = x.parent.right
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				rbTreeAugmented.leftRotate(x.parent)
				w = x.parent.right
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
				continue
			} else if w.right.color == BLACK {
				w.left.color = BLACK
				w.color = RED
				rbTreeAugmented.rightRotate(w)
				w = x.parent.right
			}
			w.color = x.parent.color
			x.parent.color = BLACK
			w.right.color = BLACK
			rbTreeAugmented.leftRotate(x.parent)
			x = rbTreeAugmented.root
		} else {
			var w *rbTreeNodeAugmented[K, A] = x.parent.left
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				rbTreeAugmented.rightRotate(x.parent)
				w = x.parent.left
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
				continue
			} else if w.left.color == BLACK {
				w.right.color = BLACK
				w.color = RED
				rbTreeAugmented.leftRotate(w)
				w = x.parent.left
			}
			w.color = x.parent.color
			x.parent.color = BLACK
			w.left.color = BLACK
			rbTreeAugmented.rightRotate(x.parent)
			x = rbTreeAugmented.root
		}
	}
	x.color = BLACK
}

type RbAugmentedIterator[K, A any] struct {
	next   			*rbTreeNodeAugmented[K, A]
	rbTreeAugmented *RbTreeAugmented[K, A]
}

// Returns an iterator pointing to least key node in the tree or to sentinel node if tree is empty.
// Used to iterate keys in the ascending order.
func (rbTreeAugmented *RbTreeAugmented[K, A]) Begin() OrderedSetForwardIterator[K] {
	var next *rbTreeNodeAugmented[K, A] = rbTreeAugmented.root
	if next != rbTreeAugmented.sentinel {
		next = getMinNode[K](rbTreeAugmented.root, rbTreeAugmented.sentinel).(*rbTreeNodeAugmented[K, A])
	}
	return &RbAugmentedIterator[K, A]{
		next:   next,
		rbTreeAugmented: rbTreeAugmented,
	}
}

// Calling Next() moves the iterator to the next greater node and returns its key.
// If Next() is called on last key(or greatest key), it returns (zeroValue, false)
func (rbAugmentedIterator *RbAugmentedIterator[K, A]) Next() (_ K, _ bool) {
	if rbAugmentedIterator.next == rbAugmentedIterator.rbTreeAugmented.sentinel {
		return
	}
	rbAugmentedIterator.next = Next[K](rbAugmentedIterator.next, rbAugmentedIterator.rbTreeAugmented.sentinel).(*rbTreeNodeAugmented[K, A])
	if rbAugmentedIterator.next == rbAugmentedIterator.rbTreeAugmented.sentinel {
		return
	}
	return rbAugmentedIterator.next.key, true
}

// Returns the key pointed by iterator. Returns (zeroValue, false) if this is called on empty tree or an iterator has completed traversing all the keys
func (rbAugmentedIterator *RbAugmentedIterator[K, A]) Key() (_ K, _ bool) {
	if rbAugmentedIterator.next != rbAugmentedIterator.rbTreeAugmented.sentinel {
		return rbAugmentedIterator.next.key, true
	}
	return
}

// Deletes the key the pointed by iterator, moves the iterator to next greater key.
// Returns the next greater key if it's present. Otherwise, returns (zeroValue, false).
// panics on calling Remove() in empty tree or an iterator has completed traversing all the keys
func (rbAugmentedIterator *RbAugmentedIterator[K, A]) Remove() (_ K, _ bool) {
	var todelete *rbTreeNodeAugmented[K, A] = rbAugmentedIterator.next
	nextKey, hasNext := rbAugmentedIterator.Next()
	rbAugmentedIterator.rbTreeAugmented.delete(todelete)
	rbAugmentedIterator.rbTreeAugmented.len--
	return nextKey, hasNext
}

type ReverseRbAugmentedIterator[K, A any] struct {
	prev   *rbTreeNodeAugmented[K, A]
	rbTreeAugmented *RbTreeAugmented[K, A]
}

// Returns an reverse iterator pointing to greatest key node in the tree or to sentinel node if tree is empty.
// Used to iterate keys in the descending order
func (rbTreeAugmented *RbTreeAugmented[K, A]) Rbegin() OrderedSetReverseIterator[K] {
	var prev *rbTreeNodeAugmented[K, A] = rbTreeAugmented.root
	if prev != rbTreeAugmented.sentinel {
		prev = getMaxNode[K](rbTreeAugmented.root, rbTreeAugmented.sentinel).(*rbTreeNodeAugmented[K, A])
	}
	return &ReverseRbAugmentedIterator[K, A]{
		prev:   prev,
		rbTreeAugmented: rbTreeAugmented,
	}
}

// Calling Prev() moves the reverse iterator to the next smaller node and returns its key.
// If Prev() is called on last key (or smallest key), it returns (zeroValue, false)
func (reverseRbAugmentedIterator *ReverseRbAugmentedIterator[K, A]) Prev() (_ K, _ bool) {
	var rbTreeAugmented *RbTreeAugmented[K, A] = reverseRbAugmentedIterator.rbTreeAugmented
	var prev *rbTreeNodeAugmented[K, A] = reverseRbAugmentedIterator.prev
	if prev == rbTreeAugmented.sentinel {
		return
	}
	reverseRbAugmentedIterator.prev = Prev[K](reverseRbAugmentedIterator.prev, reverseRbAugmentedIterator.rbTreeAugmented.sentinel).(*rbTreeNodeAugmented[K, A])
	if reverseRbAugmentedIterator.prev == rbTreeAugmented.sentinel {
		return
	}
	return reverseRbAugmentedIterator.prev.key, true
}

// Returns the key pointed by reverse iterator. Returns (zeroValue, false) if this is called on empty tree or an iterator has completed traversing all the keys
func (reverseRbAugmentedIterator *ReverseRbAugmentedIterator[K, A]) Key() (_ K, _ bool) {
	if reverseRbAugmentedIterator.prev != reverseRbAugmentedIterator.rbTreeAugmented.sentinel {
		return reverseRbAugmentedIterator.prev.key, true
	}
	return
}

// Deletes the key the pointed by reverse iterator, moves the reverse iterator to next smaller key.
// Returns the next smaller key if it's present. Otherwise, returns (zeroValue, false).
// panics on calling Remove() in empty tree or an iterator has completed traversing all the keys
func (reverseRbAugmentedIterator *ReverseRbAugmentedIterator[K, A]) Remove() (_ K, _ bool) {
	var todelete *rbTreeNodeAugmented[K, A] = reverseRbAugmentedIterator.prev
	key, hasPrev := reverseRbAugmentedIterator.Prev()
	reverseRbAugmentedIterator.rbTreeAugmented.delete(todelete)
	reverseRbAugmentedIterator.rbTreeAugmented.len--
	return key, hasPrev
}



