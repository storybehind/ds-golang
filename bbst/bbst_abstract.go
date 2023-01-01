package bbst

type BBST[K any] interface {
	Get(key K) (_ K, _ bool)
	GetGreater(key K) (_ K, _ bool)
	GetGreaterThanOrEqual(key K) (_ K, _ bool)
	GetLower(key K) (_ K, _ bool)
	GetLowerThanOrEqual(key K) (_ K, _ bool)
	Has(key K) bool

	Max() (_ K, _ bool)
	Min() (_ K, _ bool)
	Len() int64

	ReplaceOrInsert(key K) (_ K, _ bool)

	Delete(key K) (K, bool)
	DeleteMax() (K, bool)
	DeleteMin() (K, bool)

	GetRoot() (BBSTNode[K])
}

type BBSTNode[K any] interface {
	GetLeft() BBSTNode[K]
	GetRight() BBSTNode[K]
	GetParent() BBSTNode[K]
	GetKey() K
	GetAugmentedData() any
	SetAugmentedData(augmentedData any)
	IsInterfaceNil() (bool)
}

type ConcreteTag int

const (
	AvlTreeTag ConcreteTag = iota
)

type Less[K any] func(k1, k2 K) bool

type UpdateAugmentedData[K any] func(node BBSTNode[K])

type toSearch int

const (
	searchMin toSearch = iota
	searchMax
	searchGreater
	searchLower
	searchGreaterThanOrEqual
	searchLowerThanOrEqual
	searchKey
)

type compare[K any] func(k1, k2 K) int

func search[K any](node BBSTNode[K], key K, typ toSearch, cmp compare[K]) (_ K, _ bool) {
	if node.IsInterfaceNil() {
		return
	}

	switch typ {
	case searchMin:
		minNode := node
		for !minNode.GetLeft().IsInterfaceNil() {
			minNode = minNode.GetLeft()
		}
		return minNode.GetKey(), true
	case searchMax:
		maxNode := node
		for !maxNode.GetRight().IsInterfaceNil() {
			maxNode = maxNode.GetRight()
		}
		return maxNode.GetKey(), true
	case searchGreater:
		var greaterKey K
		hasGreaterKey := false
		curNode := node
		for !curNode.IsInterfaceNil() {
			compare := cmp(key, curNode.GetKey())
			if compare >= 0 {
				curNode = curNode.GetRight()
				continue
			}
			greaterKey = curNode.GetKey()
			hasGreaterKey = true
			curNode = curNode.GetLeft()
		}
		return greaterKey, hasGreaterKey
	case searchLower:
		var lowerKey K
		hasLowerKey := false
		curNode := node
		for !curNode.IsInterfaceNil() {
			compare := cmp(key, curNode.GetKey())
			if compare <= 0 {
				curNode = curNode.GetLeft()
				continue
			}
			lowerKey = curNode.GetKey()
			hasLowerKey = true
			curNode = curNode.GetRight()
		}
		return lowerKey, hasLowerKey
	case searchGreaterThanOrEqual:
		var greaterOrEqualKey K
		hasGreaterOrEqualKey := false
		curNode := node
		for !curNode.IsInterfaceNil() {
			compare := cmp(key, curNode.GetKey())
			if compare == 0 {
				greaterOrEqualKey = curNode.GetKey()
				hasGreaterOrEqualKey = true
				break
			}
			if compare == 1 {
				curNode = curNode.GetRight()
				continue
			}
			greaterOrEqualKey = curNode.GetKey()
			hasGreaterOrEqualKey = true
			curNode = curNode.GetLeft()
		}
		return greaterOrEqualKey, hasGreaterOrEqualKey
	case searchLowerThanOrEqual:
		var lowerOrEqualKey K
		hasLowerOrEqualKey := false
		curNode := node
		for !curNode.IsInterfaceNil() {
			compare := cmp(key, curNode.GetKey())
			if compare == 0 {
				lowerOrEqualKey = curNode.GetKey()
				hasLowerOrEqualKey = true
				break
			}
			if compare == -1 {
				curNode = curNode.GetLeft()
				continue
			}
			lowerOrEqualKey = curNode.GetKey()
			hasLowerOrEqualKey = true
			curNode = curNode.GetRight()
		}
		return lowerOrEqualKey, hasLowerOrEqualKey
	case searchKey:
		curNode := node
		for !curNode.IsInterfaceNil() {
			compare := cmp(key, curNode.GetKey())
			if compare == 0 {
				return curNode.GetKey(), true
			}
			if compare == -1 {
				curNode = curNode.GetLeft()
			} else {
				curNode = curNode.GetRight()
			}
		}
		return
	default:
		panic("invalid search type")
	}
}
