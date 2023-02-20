package orderedset

type binarySearchTree[K Less[K]] interface {
	OrderedSet[K]
	getRoot() binarySearchTreeNode[K]
	setRoot(binarySearchTreeNode[K])
}

type binarySearchTreeNode[K Less[K]] interface {
	getLeft() binarySearchTreeNode[K]
	setLeft(binarySearchTreeNode[K])
	getRight() binarySearchTreeNode[K]
	setRight(binarySearchTreeNode[K])
	getParent() binarySearchTreeNode[K]
	setParent(binarySearchTreeNode[K])
	getKey() K
	setKey(K)
	isNodeNil() bool
}

func cmp[K Less[K]](k1, k2 K) int {
	if k1.Less(k2) {
		return -1
	}
	if k2.Less(k1) {
		return 1
	}
	return 0
}

func getMinNode[K Less[K]](node binarySearchTreeNode[K]) binarySearchTreeNode[K] {
	var minNode binarySearchTreeNode[K] = node
	for !minNode.getLeft().isNodeNil() {
		minNode = minNode.getLeft()
	}
	return minNode
}

func getMaxNode[K Less[K]](node binarySearchTreeNode[K]) binarySearchTreeNode[K] {
	var maxNode binarySearchTreeNode[K] = node
	for !maxNode.getRight().isNodeNil() {
		maxNode = maxNode.getRight()
	}
	return maxNode
}

func searchGreaterNode[K Less[K]](node binarySearchTreeNode[K], key K) binarySearchTreeNode[K] {
	var greaterNode binarySearchTreeNode[K]
	var curNode binarySearchTreeNode[K] = node
	for !curNode.isNodeNil() {
		compare := cmp(key, curNode.getKey())
		if compare >= 0 {
			curNode = curNode.getRight()
			continue
		}
		greaterNode = curNode
		curNode = curNode.getLeft()
	}
	return greaterNode
}

func searchLowerNode[K Less[K]](node binarySearchTreeNode[K], key K) binarySearchTreeNode[K] {
	var lowerNode binarySearchTreeNode[K]
	var curNode binarySearchTreeNode[K] = node
	for !curNode.isNodeNil() {
		compare := cmp(key, curNode.getKey())
		if compare <= 0 {
			curNode = curNode.getLeft()
			continue
		}
		lowerNode = curNode
		curNode = curNode.getRight()
	}
	return lowerNode
}

func searchGreaterThanOrEqualNode[K Less[K]](node binarySearchTreeNode[K], key K) binarySearchTreeNode[K] {
	var greaterThanOrEqualKeyNode binarySearchTreeNode[K]
	var curNode binarySearchTreeNode[K] = node
	for !curNode.isNodeNil() {
		compare := cmp(key, curNode.getKey())
		if compare == 0 {
			greaterThanOrEqualKeyNode = curNode
			break
		}
		if compare == 1 {
			curNode = curNode.getRight()
			continue
		}
		greaterThanOrEqualKeyNode = curNode
		curNode = curNode.getLeft()
	}
	return greaterThanOrEqualKeyNode
}

func searchLowerThanOrEqualNode[K Less[K]](node binarySearchTreeNode[K], key K) binarySearchTreeNode[K] {
	var lowerThanOrEqualKeyNode binarySearchTreeNode[K]
	var curNode binarySearchTreeNode[K] = node
	for !curNode.isNodeNil() {
		compare := cmp(key, curNode.getKey())
		if compare == 0 {
			lowerThanOrEqualKeyNode = curNode
			break
		}
		if compare == -1 {
			curNode = curNode.getLeft()
			continue
		}
		lowerThanOrEqualKeyNode = curNode
		curNode = curNode.getRight()
	}
	return lowerThanOrEqualKeyNode
}

func searchNode[K Less[K]](node binarySearchTreeNode[K], key K) binarySearchTreeNode[K] {
	var curNode binarySearchTreeNode[K] = node
	for !curNode.isNodeNil() {
		compare := cmp(key, curNode.getKey())
		if compare == 0 {
			return curNode
		}
		if compare == -1 {
			curNode = curNode.getLeft()
		} else {
			curNode = curNode.getRight()
		}
	}
	return nil
}

func leftRotate[K Less[K]](bst binarySearchTree[K], x binarySearchTreeNode[K]) {
	var y, p binarySearchTreeNode[K] = x.getRight(), x.getParent()
	var beta binarySearchTreeNode[K] = y.getLeft()
	x.setRight(beta)
	if !beta.isNodeNil() {
		beta.setParent(x)
	}
	y.setParent(p)
	if p.isNodeNil() {
		bst.setRoot(y)
	} else if x == p.getLeft() {
		p.setLeft(y)
	} else {
		p.setRight(y)
	}
	y.setLeft(x)
	x.setParent(y)
}

func rightRotate[K Less[K]](bst binarySearchTree[K], y binarySearchTreeNode[K]) {
	var x, p binarySearchTreeNode[K] = y.getLeft(), y.getParent()
	var beta binarySearchTreeNode[K] = x.getRight()
	y.setLeft(beta)
	if !beta.isNodeNil() {
		beta.setParent(y)
	}
	x.setParent(p)
	if p.isNodeNil() {
		bst.setRoot(x)
	} else if y == p.getLeft() {
		p.setLeft(x)
	} else {
		p.setRight(x)
	}
	x.setRight(y)
	y.setParent(x)
}


