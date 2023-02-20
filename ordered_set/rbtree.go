package orderedset

type rbTreeNode[K Less[K]] struct {
	left, right, parent *rbTreeNode[K]
	key                 K
	color               color
}

// func (node *rbTreeNode[K]) getLeft() binarySearchTreeNode[K] {
// 	return node.left
// }

// func (node *rbTreeNode[K]) setLeft(value binarySearchTreeNode[K]) {
// 	node.left = value.(*rbTreeNode[K])
// }

// func (node *rbTreeNode[K]) getRight() binarySearchTreeNode[K] {
// 	return node.right
// }

// func (node *rbTreeNode[K]) setRight(value binarySearchTreeNode[K]) {
// 	node.right = value.(*rbTreeNode[K])
// }

// func (node *rbTreeNode[K]) getParent() binarySearchTreeNode[K] {
// 	return node.parent
// }

// func (node *rbTreeNode[K]) setParent(value binarySearchTreeNode[K]) {
// 	node.parent = value.(*rbTreeNode[K])
// }

// func (node *rbTreeNode[K]) getKey() K {
// 	return node.key
// }

// func (node *rbTreeNode[K]) setKey(value K) {
// 	node.key = value
// }

// func (node *rbTreeNode[K]) isNodeNil() bool {
// 	return node == nil
// }

// color type details color of a node
type color byte

const (
	RED color = iota
	BLACK
)

// Maintains unique set of keys
type RbTree[K Less[K]] struct {
	root     *rbTreeNode[K]
	sentinel *rbTreeNode[K]
	// end *RbIterator[K]
	len      int64
}

// Returns instance of Red-Black Tree.
// Less method determines the order of key.
// k1 precedes k2 if and only if Less(k1, k2) return true.
// k1 equals k2 if and only if !Less(k1, k2) && !Less(k2, k1) holds true.
func NewRbTree[K Less[K]]() *RbTree[K] {
	return &RbTree[K]{
		root: nil,
		sentinel: &rbTreeNode[K]{
			color: BLACK,
		},
		// end: &RbIterator[K]{},
		len: 0,
	}
}

// Get looks for the key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (rbTree *RbTree[K]) Get(key K) (_ K, _ bool) {
	if rbTree.root == nil {
		return
	}
	var node *rbTreeNode[K] = rbTree.searchNode(rbTree.root, key)
	if node != nil {
		return node.key, true
	}
	return
}

func (rbTree *RbTree[K]) searchNode(node *rbTreeNode[K], key K) *rbTreeNode[K] {
	var curNode *rbTreeNode[K] = node
	for curNode != rbTree.sentinel {
		compare := cmp(key, curNode.key)
		if compare == 0 {
			return curNode
		}
		if compare == -1 {
			curNode = curNode.left
		} else {
			curNode = curNode.right
		}
	}
	return nil
}

// GetGreater looks for smallest key that is strictly greater than key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (rbTree *RbTree[K]) GetGreater(key K) (_ K, _ bool) {
	if rbTree.root == nil {
		return
	}
	var greaterNode *rbTreeNode[K] = rbTree.searchGreaterNode(rbTree.root, key)
	if greaterNode != nil {
		return greaterNode.key, true
	}
	return
}

func (rbTree *RbTree[K]) searchGreaterNode(node *rbTreeNode[K], key K) *rbTreeNode[K] {
	var greaterNode *rbTreeNode[K]
	var curNode *rbTreeNode[K] = node
	for curNode != rbTree.sentinel {
		compare := cmp(key, curNode.key)
		if compare >= 0 {
			curNode = curNode.right
			continue
		}
		greaterNode = curNode
		curNode = curNode.left
	}
	return greaterNode
}

// GetGreaterThanOrEqual looks for smallest key that is greater than or equal to key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (rbTree *RbTree[K]) GetGreaterThanOrEqual(key K) (_ K, _ bool) {
	if rbTree.root == nil {
		return
	}
	var greaterThanOrEqualNode *rbTreeNode[K] = rbTree.searchGreaterThanOrEqualNode(rbTree.root, key)
	if greaterThanOrEqualNode != nil {
		return greaterThanOrEqualNode.key, true
	}
	return
}

func (rbTree *RbTree[K]) searchGreaterThanOrEqualNode(node *rbTreeNode[K], key K) *rbTreeNode[K] {
	var greaterThanOrEqualKeyNode *rbTreeNode[K]
	var curNode *rbTreeNode[K] = node
	for curNode != rbTree.sentinel {
		compare := cmp(key, curNode.key)
		if compare == 0 {
			greaterThanOrEqualKeyNode = curNode
			break
		}
		if compare == 1 {
			curNode = curNode.right
			continue
		}
		greaterThanOrEqualKeyNode = curNode
		curNode = curNode.left
	}
	return greaterThanOrEqualKeyNode
}

// GetLower looks for greatest key that is strictly lower than key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (rbTree *RbTree[K]) GetLower(key K) (_ K, _ bool) {
	if rbTree.root == nil {
		return
	}
	var lowerNode *rbTreeNode[K] = rbTree.searchLowerNode(rbTree.root, key)
	if lowerNode != nil {
		return lowerNode.key, true
	}
	return
}

func (rbTree *RbTree[K]) searchLowerNode(node *rbTreeNode[K], key K) *rbTreeNode[K] {
	var lowerNode *rbTreeNode[K]
	var curNode *rbTreeNode[K] = node
	for curNode != rbTree.sentinel {
		compare := cmp(key, curNode.key)
		if compare <= 0 {
			curNode = curNode.left
			continue
		}
		lowerNode = curNode
		curNode = curNode.right
	}
	return lowerNode
}

// GetLowerThanOrEqual looks for greatest key that is lower than or equal to key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (rbTree *RbTree[K]) GetLowerThanOrEqual(key K) (_ K, _ bool) {
	var lowerThanOrEqualNode *rbTreeNode[K] = rbTree.searchLowerThanOrEqualNode(rbTree.root, key)
	if lowerThanOrEqualNode != nil {
		return lowerThanOrEqualNode.key, true
	}
	return
}

func (rbTree *RbTree[K]) searchLowerThanOrEqualNode(node *rbTreeNode[K], key K) *rbTreeNode[K] {
	var lowerThanOrEqualKeyNode *rbTreeNode[K]
	var curNode *rbTreeNode[K] = node
	for curNode != rbTree.sentinel {
		compare := cmp(key, curNode.key)
		if compare == 0 {
			lowerThanOrEqualKeyNode = curNode
			break
		}
		if compare == -1 {
			curNode = curNode.left
			continue
		}
		lowerThanOrEqualKeyNode = curNode
		curNode = curNode.right
	}
	return lowerThanOrEqualKeyNode
}


// // Has returns true if the given key is in the tree
// func (rbtree *RbTree[K]) Has(key K) bool {
// 	var node binaryTreeNode[K] = searchNode[K](rbtree.root, key)
// 	if node != nil {
// 		return true
// 	}
// 	return false
// }

// Max returns the largest key in the tree, or (zeroValue, false) if the tree is empty
func (rbTree *RbTree[K]) Max() (_ K, _ bool) {
	if rbTree.root == nil {
		return
	}
	var maxNode *rbTreeNode[K] = rbTree.getMaxNode(rbTree.root)
	return maxNode.key, true
}

// node can't be nil or rbTree.sentinel
func (rbTree *RbTree[K]) getMaxNode(node *rbTreeNode[K]) *rbTreeNode[K] {
	var maxNode *rbTreeNode[K] = node
	for maxNode.right != rbTree.sentinel {
		maxNode = maxNode.right
	}
	return maxNode
}

// Min returns the smallest key in the tree, or (zeroValue, false) if the tree is empty
func (rbTree *RbTree[K]) Min() (_ K, _ bool) {
	if rbTree.root == nil {
		return
	}
	var minNode *rbTreeNode[K] = rbTree.getMinNode(rbTree.root)
	return minNode.key, true
}

// node can't be nil or rbTree.sentinel
func (rbTree *RbTree[K]) getMinNode(node *rbTreeNode[K]) *rbTreeNode[K] {
	var minNode *rbTreeNode[K] = node
	for minNode.left != rbTree.sentinel {
		minNode = minNode.left
	}
	return minNode
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
		if key.Less(x.key) {
			x = x.left
		} else if x.key.Less(key) {
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
	} else if key.Less(y.key) {
		y.left = z
	} else {
		y.right = z
	}
	z.left = rbTree.sentinel
	z.right = rbTree.sentinel
	z.color = RED
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


func (rbTree *RbTree[K]) Delete(key K) (_ K, _ bool) {
	if rbTree.root == nil {
		return
	}
	var z *rbTreeNode[K] = rbTree.searchNode(rbTree.root, key)
	if z == nil {
		return
	}
	var deletedKey K = z.key
	rbTree.delete(z)
	return deletedKey, true
}

func (rbTree *RbTree[K]) DeleteMax() (_ K, _ bool) {
	if rbTree.root == nil {
		return
	}
	var z *rbTreeNode[K] = rbTree.getMaxNode(rbTree.root)
	var deletedKey K = z.key
	rbTree.delete(z)
	return deletedKey, true
}

func (rbTree *RbTree[K]) DeleteMin() (_ K, _ bool) {
	if rbTree.root == nil {
		return
	}
	var z *rbTreeNode[K] = rbTree.getMinNode(rbTree.root)
	var deletedKey K = z.key
	rbTree.delete(z)
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
		y = rbTree.getMinNode(z.right)
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

// type RbIterator[K Less[K]] struct {
// 	cur *rbTreeNode[K]
// 	rbTree *RbTree[K]
// }

// func (rbTree *RbTree[K]) Iterator() *RbIterator[K] {
// 	// var next *rbTreeNode[K] = rbTree.sentinel
// 	// if rbTree.root != nil {
// 	// 	next = rbTree.getMinNode(rbTree.root)
// 	// }
// 	return &RbIterator[K] {
// 		cur: rbTree.sentinel,
// 		rbTree: rbTree,
// 	}
// }

// func (rbTree *RbTree[K]) Begin() *RbIterator[K] {
// 	if rbTree.root == nil {
// 		return rbTree.end
// 	}
// 	return &RbIterator[K] {
// 		cur: rbTree.getMinNode(rbTree.root),
// 		rbTree: rbTree,
// 	}
// }

// func (rbTree *RbTree[K]) End() *RbIterator[K] {
// 	return rbTree.end
// }

// func (rbIterator *RbIterator[K]) Next() (*RbIterator[K]) {
// 	var rbTree *RbTree[K] = rbIterator.rbTree
// 	var cur *rbTreeNode[K] = rbIterator.cur
// 	var key K = cur.key
// 	if cur.right != rbTree.sentinel {
// 		rbIterator.cur = rbTree.getMinNode(cur.right)
// 		return key, true
// 	}
// 	for cur.parent != rbTree.sentinel && cur == cur.parent.right {
// 		cur = cur.parent
// 	}
// 	rbIterator.cur = cur.parent
// 	return key, true
// }
