package bbst

//Node of avl tree that holds a particular key
//Supports augmentation
//Maintain left, right and parent pointer for tree traversal
//parent of root node is nil
//left and right pointer of leaf nodes are nil
type AvlTreeNode[K any] struct {
	left, right, parent *AvlTreeNode[K]
	key                 K
	augmentedData       any
	height int64
}

//Get left node 
func (node *AvlTreeNode[K]) GetLeft() BBSTNode[K] {
	return node.left
}

//Get right node
func (node *AvlTreeNode[K]) GetRight() BBSTNode[K] {
	return node.right
}

//Get parent node
func (node *AvlTreeNode[K]) GetParent() BBSTNode[K] {
	return node.parent
}

//Get key
func (node *AvlTreeNode[K]) GetKey() K {
	return node.key
}

//Get augmented data
func (node *AvlTreeNode[K]) GetAugmentedData() any {
	return node.augmentedData
}

//Set augmented data
func (node *AvlTreeNode[K]) SetAugmentedData(augmentedData any) {
	node.augmentedData = augmentedData
}

//IsInterfaceNil() return true if node is a nil pointer
func (node *AvlTreeNode[K]) IsInterfaceNil() (bool) {
	return node == nil
}

// Maintains unique set of keys
type AvlTree[K any] struct {
	root                *AvlTreeNode[K]
	less                Less[K]
	cmp                 compare[K]
	updateAugmentedData *UpdateAugmentedData[K]
	len 				int64
}

//Returns instance of AvlTree.
// Less method determines the order of key. 
// k1 precedes k2 in AvlTree if and only if Less(k1, k2) return true.
// k1 equals k2 if and only if !Less(k1, k2) && !Less(k2, k1) holds true.
func NewAvlTreeByLess[K any](less Less[K]) *AvlTree[K] {
	return &AvlTree[K]{
		root: nil,
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
		updateAugmentedData: nil,
		len: 0,
	}
}

// Returns instance of AvlTree.
// Less method determines the order of key. 
// k1 precedes k2 in AvlTree if and only if Less(k1, k2) return true.
// k1 equals k2 if and only if !Less(k1, k2) && !Less(k2, k1) holds true.
// UpdateAugmentedData method maintains invariant of node's augmentedData as tree performs rotation
// Efficient in case of easy tree augmentation which depends on node's immediate children
func NewAvlTreeByLessAndUpdateAugmentedData[K any](less Less[K], updateAugmentedData *UpdateAugmentedData[K]) *AvlTree[K] {
	return &AvlTree[K]{
		root: nil,
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
		updateAugmentedData: updateAugmentedData,
		len: 0,
	}
}

// Get looks for the key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (avlTree *AvlTree[K]) Get(key K) (_ K, _ bool) {
	return search[K](avlTree.root, key, searchKey, avlTree.cmp)
}

// GetGreater looks for smallest key that is strictly greater than key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (avlTree *AvlTree[K]) GetGreater(key K) (_ K, _ bool) {
	return search[K](avlTree.root, key, searchGreater, avlTree.cmp)
}

// GetGreaterThanOrEqual looks for smallest key that is greater than or equal to key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (avlTree *AvlTree[K]) GetGreaterThanOrEqual(key K) (_ K, _ bool) {
	return search[K](avlTree.root, key, searchGreaterThanOrEqual, avlTree.cmp)
}

// GetLower looks for greatest key that is strictly lower than key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (avlTree *AvlTree[K]) GetLower(key K) (_ K, _ bool) {
	return search[K](avlTree.root, key, searchLower, avlTree.cmp)
}

// GetLowerThanOrEqual looks for greatest key that is lower than or equal to key in the tree, returning it. It returns (zeroValue, false) if unable to find that key
func (avlTree *AvlTree[K]) GetLowerThanOrEqual(key K) (_ K, _ bool) {
	return search[K](avlTree.root, key, searchLowerThanOrEqual, avlTree.cmp)
}

// Has returns true if the given key is in the tree
func (avlTree *AvlTree[K]) Has(key K) bool {
	_, has := search[K](avlTree.root, key, searchKey, avlTree.cmp)
	return has
}

// Max returns the largest key in the tree, or (zeroValue, false) if the tree is empty
func (avlTree *AvlTree[K]) Max() (_ K, _ bool) {
	var zero K
	return search[K](avlTree.root, zero, searchMax, avlTree.cmp)
}

// Min returns the smallest key in the tree, or (zeroValue, false) if the tree is empty
func (avlTree *AvlTree[K]) Min() (_ K, _ bool) {
	var zero K
	return search[K](avlTree.root, zero, searchMin, avlTree.cmp)
}

// Len returns the number of keys currently in the tree.
func (avlTree *AvlTree[K]) Len() int64 {
	return avlTree.len
}

// toRemove details what item to remove in a node.remove call.
type toRemove int

const (
	removeKey toRemove = iota // removes the given item
	removeMin                  // removes smallest item in the subtree
	removeMax                  // removes largest item in the subtree
)

// ReplaceOrInsert adds the given key to the tree.
// If a key in the tree already equals the given one, it is removed from the tree and returned, and the second return value is true.
// Otherwise, (zeroValue, false)
func (avlTree *AvlTree[K]) ReplaceOrInsert(key K) (_ K, _ bool) {
	var prevKey K
	var has bool
	avlTree.root, prevKey, has = avlTree.replaceOrInsert(avlTree.root, key)
	avlTree.root.parent = nil
	if !has {
		avlTree.len++
	}
	return prevKey, has
}

// Delete removes an item equal to the passed in item from the tree, returning it. If no such item exists, returns (zeroValue, false)
func (avlTree *AvlTree[K]) Delete(key K) (K, bool) {
	var deletedKey K
	var deleted bool
	avlTree.root, deletedKey, deleted = avlTree.delete(avlTree.root, key, removeKey)
	if avlTree.root != nil {
		avlTree.root.parent = nil
	}
	if deleted {
		avlTree.len--
	}
	return deletedKey, deleted
}

// DeleteMax removes the largest item in the tree and returns it. If no such item exists, returns (zeroValue, false)
func (avlTree *AvlTree[K]) DeleteMax() (K, bool) {
	var deletedKey K
	var deleted bool
	var zero K
	avlTree.root, deletedKey, deleted = avlTree.delete(avlTree.root, zero, removeMax)
	if avlTree.root != nil {
		avlTree.root.parent = nil
	}
	if deleted {
		avlTree.len--
	}
	return deletedKey, deleted
}

// DeleteMin removes the smallest item in the tree and returns it. If no such item exists, returns (zeroValue, false)
func (avlTree *AvlTree[K]) DeleteMin() (K, bool) {
	var deletedKey K
	var deleted bool
	var zero K
	avlTree.root, deletedKey, deleted = avlTree.delete(avlTree.root, zero, removeMin)
	if avlTree.root != nil {
		avlTree.root.parent = nil
	}
	if deleted {
		avlTree.len--
	}
	return deletedKey, deleted
}

// Get root node of the tree
func (avlTree *AvlTree[K]) GetRoot() BBSTNode[K] {
	return avlTree.root
}

func (avlTree *AvlTree[K]) replaceOrInsert(node *AvlTreeNode[K], key K) (_ *AvlTreeNode[K], _ K, _ bool) {
	if node == nil {
		var zero K
		newNode := &AvlTreeNode[K] {
			left:   nil,
			right:  nil,
			parent: nil,
			key:    key,
			height: 1,
		}
		if avlTree.updateAugmentedData != nil {
			(*avlTree.updateAugmentedData)(newNode)
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
	node.height = 1 + max(getHeight(node.left), getHeight(node.right))
	if avlTree.updateAugmentedData != nil {
		(*avlTree.updateAugmentedData)(node)
	}
	return avlTree.balanceNode(node), prevKey, has
}

func (avlTree *AvlTree[K]) balanceNode(node *AvlTreeNode[K]) *AvlTreeNode[K] {
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

func (avlTree AvlTree[K]) delete(node *AvlTreeNode[K], key K, typ toRemove) (_ *AvlTreeNode[K], _ K, _ bool) {
	if node == nil {
		return
	}

	var deletedKey K
	var deleted bool
	switch typ {
	case removeMin:
		if node.left == nil {
			return node.right, node.key, true
		}
		node.left, deletedKey, deleted = avlTree.delete(node.left, key, typ)
		if node.left != nil {
			node.left.parent = node
		}
	case removeMax:
		if node.right == nil {
			return node.left, node.key, true
		}
		node.right, deletedKey, deleted = avlTree.delete(node.right, key, typ)
		if node.right != nil {
			node.right.parent = node
		}
	case removeKey:
		compare := avlTree.cmp(key, node.key)
		switch compare {
		case 0:
			deletedKey = node.key
			deleted = true
			if node.left != nil {
				var zero, maxKey K
				node.left, maxKey, _ = avlTree.delete(node.left, zero, removeMax)
				if node.left != nil {
					node.left.parent = node
				}
				node.key = maxKey
			} else if node.right != nil {
				var zero, minKey K
				node.right, minKey, _ = avlTree.delete(node.right, zero, removeMin)
				if node.right != nil {
					node.right.parent = node
				}
				node.key = minKey
			} else {
				return nil, node.key, true
			}
		case -1:
			node.left, deletedKey, deleted = avlTree.delete(node.left, key, typ)
			if node.left != nil {
				node.left.parent = node
			}
		case 1:
			node.right, deletedKey, deleted = avlTree.delete(node.right, key, typ)
			if node.right != nil {
				node.right.parent = node
			}
		}
	default:
		panic("invalid remove type")
	}
	node.height = 1 + max(getHeight(node.left), getHeight(node.right))
	if avlTree.updateAugmentedData != nil {
		(*avlTree.updateAugmentedData)(node)
	}
	return avlTree.balanceNode(node), deletedKey, deleted
}

func getHeightDiff[K any](node *AvlTreeNode[K]) int64 {
	return getHeight(node.right) - getHeight(node.left)
}

func getHeight[K any](node *AvlTreeNode[K]) int64 {
	if node == nil {
		return 0
	}
	return node.height
}

func max(a, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}

func (avlTree *AvlTree[K]) rightRotate(y *AvlTreeNode[K]) *AvlTreeNode[K] {
	x := y.left
	y.left = x.right
	if y.left != nil {
		y.left.parent = y
	}

	x.right = y
	y.parent = x

	y.height = max(getHeight(y.left), getHeight(y.right)) + 1
	x.height = max(getHeight(x.left), getHeight(x.right)) + 1

	if avlTree.updateAugmentedData != nil {
		(*avlTree.updateAugmentedData)(y)
		(*avlTree.updateAugmentedData)(x)
	}
	return x
}

func (avlTree *AvlTree[K]) leftRotate(x *AvlTreeNode[K]) *AvlTreeNode[K] {
	y := x.right
	x.right = y.left
	if x.right != nil {
		x.right.parent = x
	}

	y.left = x
	x.parent = y

	x.height = max(getHeight(x.left), getHeight(x.right)) + 1
	y.height = max(getHeight(y.left), getHeight(y.right)) + 1

	if avlTree.updateAugmentedData != nil {
		(*avlTree.updateAugmentedData)(x)
		(*avlTree.updateAugmentedData)(y)
	}
	return y
}
