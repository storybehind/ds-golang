# gocontainer

gocontainer is a package that provides implementation of some standard containers in Golang. It mainly tries to provide containers that is not supported in standard library and also to provide rich interfaces to these containers.

## Usage

```
go get -u github.com/storybehind/gocontainer
```

## packages

* [orderedset](#orderedset)
  * [Interfaces](#interfaces) 
  * [Red-Black Tree](#RbTree)
  * [AVL Tree](#AvlTree)
  * [OrderStatisticsTree](#OrderStatisticsTree)
  * [Augmentation](#Augmentation)
* [orderedmap](#orderedmap)
* [priorityqueue](#priorityqueue)
  * [BinaryHeap](#BinaryHeap)

### orderedset

orderedset provides containers to support insertion, deletion and search keys while maintaining the order of keys.

#### interfaces

```go OrderedSet
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
```
```go OrderedSetI
// OrderedSet Interface with iterator interfaces
type OrderedSetI[K any] interface {
	OrderedSet[K]
	Begin() OrderedSetForwardIterator[K]
	Rbegin() OrderedSetReverseIterator[K]
}
```
```go OrderedSetForwardIterator
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

```
```go OrderedSetReverseIterator
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
```

Implementations of OrderedSetI : Red-Black Tree, AvlTree, OrderStatisticsTree, RbTreeAugmented

#### RbTree

RbTree (or Red-Black Tree) implements [OrderedSetI](#interfaces). Supports insertion, deletion and search for keys in O(log n) where n is the number of keys in the tree. Keys can be iterated in ascending (or) descending order in O(n) time.

```go
package main

import (
	"fmt"

	"github.com/storybehind/gocontainer/orderedset"
)

func main() {
	// Initialize orderedset
	// Calling New(), return RbTree instance by default
  // Less method determines order of keys. Here k1 < k2 iff Less(k1, k2) returns true	
	os := orderedset.New[int](func(k1, k2 int) bool {return k1 < k2})

	// Insert keys to the set
	for key:=1; key<=5; key++ {
		_, isReplaced := os.ReplaceOrInsert(key)
		fmt.Printf("key: %d, isReplaced: %v\n", key, isReplaced)
	}

	// Search keys
	key, has := os.Get(1)
	fmt.Printf("Get 1; key: %v, has: %v\n", key, has)
	key, has = os.Get(0)
	fmt.Printf("Get 0; key: %v, has: %v\n", key, has)
	key, has = os.GetGreater(0)
	fmt.Printf("GetGreater 0; key: %v, has: %v\n", key, has)
	
	// Delete key
	key, isDeleted := os.Delete(1)
	fmt.Printf("Delete 1; key: %v, isDeleted: %v\n", key, isDeleted)

	// Iterate keys in ascending order
	forwardItr := os.Begin()
	for key, has = forwardItr.Key(); has; key, has = forwardItr.Next() {
		fmt.Printf("forwardItr; key: %v\n", key)
	}

	// Iterate keys in descending order
	reverseItr := os.Rbegin()
	for key, has = reverseItr.Key(); has; key, has = reverseItr.Prev() {
		fmt.Printf("reverseItr; key: %v\n", key)
	}

	// Remove if certain condition satisfies
	forwardItr = os.Begin()
	for key, has = forwardItr.Key(); has; {
		if key == 3 {
			key, has = forwardItr.Remove()
			continue
		}
		key, has = forwardItr.Next()
	}
	// Iterate keys in ascending order
	forwardItr = os.Begin()
	for key, has = forwardItr.Key(); has; key, has = forwardItr.Next() {
		fmt.Printf("forwardItr2; key: %v\n", key)
	}

	// Output
	// key: 1, isReplaced: false
	// key: 2, isReplaced: false
	// key: 3, isReplaced: false
	// key: 4, isReplaced: false
	// key: 5, isReplaced: false
	// Get 1; key: 1, has: true
	// Get 0; key: 0, has: false
	// GetGreater 0; key: 1, has: true
	// Delete 1; key: 1, isDeleted: true
	// forwardItr; key: 2
	// forwardItr; key: 3
	// forwardItr; key: 4
	// forwardItr; key: 5
	// reverseItr; key: 5
	// reverseItr; key: 4
	// reverseItr; key: 3
	// reverseItr; key: 2
	// forwardItr2; key: 2
	// forwardItr2; key: 4
	// forwardItr2; key: 5
}
```

#### AvlTree

AvlTree (or AVL Tree) implements [OrderedSetI](#interfaces). Supports insertion, deletion and search for keys in O(log n) where n is the number of keys in the tree. Keys can be iterated in ascending (or) descending order in O(n) time.

```go
package main

import (
	"fmt"

	"github.com/storybehind/gocontainer/orderedset"
)

func main() {
	// Initialize orderedset
	// Calling NewAvlTree(), return AvlTree instance
	// Less method determines order of keys. Here k1 < k2 iff Less(k1, k2) returns true
	os := orderedset.NewAvlTree[int](func(k1, k2 int) bool {return k1 < k2})

	// Insert keys to the set
	for key:=1; key<=5; key++ {
		_, isReplaced := os.ReplaceOrInsert(key)
		fmt.Printf("key: %d, isReplaced: %v\n", key, isReplaced)
	}

	// Search keys
	key, has := os.Get(1)
	fmt.Printf("Get 1; key: %v, has: %v\n", key, has)
	key, has = os.Get(0)
	fmt.Printf("Get 0; key: %v, has: %v\n", key, has)
	key, has = os.GetGreater(0)
	fmt.Printf("GetGreater 0; key: %v, has: %v\n", key, has)
	
	// Delete key
	key, isDeleted := os.Delete(1)
	fmt.Printf("Delete 1; key: %v, isDeleted: %v\n", key, isDeleted)

	// Iterate keys in ascending order
	forwardItr := os.Begin()
	for key, has = forwardItr.Key(); has; key, has = forwardItr.Next() {
		fmt.Printf("forwardItr; key: %v\n", key)
	}

	// Iterate keys in descending order
	reverseItr := os.Rbegin()
	for key, has = reverseItr.Key(); has; key, has = reverseItr.Prev() {
		fmt.Printf("reverseItr; key: %v\n", key)
	}

	// Remove if certain condition satisfies
	forwardItr = os.Begin()
	for key, has = forwardItr.Key(); has; {
		if key == 3 {
			key, has = forwardItr.Remove()
			continue
		}
		key, has = forwardItr.Next()
	}
	// Iterate keys in ascending order
	forwardItr = os.Begin()
	for key, has = forwardItr.Key(); has; key, has = forwardItr.Next() {
		fmt.Printf("forwardItr2; key: %v\n", key)
	}

	// Output
	// key: 1, isReplaced: false
	// key: 2, isReplaced: false
	// key: 3, isReplaced: false
	// key: 4, isReplaced: false
	// key: 5, isReplaced: false
	// Get 1; key: 1, has: true
	// Get 0; key: 0, has: false
	// GetGreater 0; key: 1, has: true
	// Delete 1; key: 1, isDeleted: true
	// forwardItr; key: 2
	// forwardItr; key: 3
	// forwardItr; key: 4
	// forwardItr; key: 5
	// reverseItr; key: 5
	// reverseItr; key: 4
	// reverseItr; key: 3
	// reverseItr; key: 2
	// forwardItr2; key: 2
	// forwardItr2; key: 4
	// forwardItr2; key: 5
}
```

#### OrderStatisticsTree

OrderStatisticsTree supports insertion, deletion, search, rank and select operations in O(log n) time where n is number of keys in the tree. Keys can be iterated in ascending (or) descending order in O(n) time. It [augments](#Augmentation) node's subtree size.

Rank(key): Determines the index of given key starting from zero. Ex: rank of minimum key will be zero. Returns -1 if key is not found in the tree.

Select(r): Return key element whose rank(key) = r. Ex : for r = 0 , return minimum key. If r >= Len(), return zeroValue, false.

```go
package main

import (
	"fmt"

	"github.com/storybehind/gocontainer/orderedset/variants"
)

func main() {
	// Initialize OrderStatisticsTree
	// Less method determines order of keys. Here k1 < k2 iff Less(k1, k2) returns true
	os := variants.NewOrderStatisticsTree[int](func(k1, k2 int) bool { return k1 < k2 })

	// Insert keys to the set
	for key := 1; key <= 5; key++ {
		_, isReplaced := os.ReplaceOrInsert(key)
		fmt.Printf("key: %d, isReplaced: %v\n", key, isReplaced)
	}

	// Search keys
	key, has := os.Get(1)
	fmt.Printf("Get 1; key: %v, has: %v\n", key, has)
	key, has = os.Get(0)
	fmt.Printf("Get 0; key: %v, has: %v\n", key, has)
	key, has = os.GetGreater(0)
	fmt.Printf("GetGreater 0; key: %v, has: %v\n", key, has)

	// Select key
	key, has = os.Select(0)
	fmt.Printf("Select 0; key: %v, has: %v\n", key, has)
	key, has = os.Select(4)
	fmt.Printf("Select 4; key: %v, has: %v\n", key, has)
	key, has = os.Select(5)
	fmt.Printf("Select 5; key: %v, has: %v\n", key, has)

	// Number of keys between -2(exclusive) and 3(inclusive) i.e (-2, 3]
	count := os.Rank(3) - os.Rank(-2)
	fmt.Printf("count keys: %d\n", count)

	// Delete key
	key, isDeleted := os.Delete(1)
	fmt.Printf("Delete 1; key: %v, isDeleted: %v\n", key, isDeleted)

	// Iterate keys in ascending order
	forwardItr := os.Begin()
	for key, has = forwardItr.Key(); has; key, has = forwardItr.Next() {
		fmt.Printf("forwardItr; key: %v\n", key)
	}

	// Iterate keys in descending order
	reverseItr := os.Rbegin()
	for key, has = reverseItr.Key(); has; key, has = reverseItr.Prev() {
		fmt.Printf("reverseItr; key: %v\n", key)
	}

	// Remove if certain condition satisfies
	forwardItr = os.Begin()
	for key, has = forwardItr.Key(); has; {
		if key == 3 {
			key, has = forwardItr.Remove()
			continue
		}
		key, has = forwardItr.Next()
	}
	// Iterate keys in ascending order
	forwardItr = os.Begin()
	for key, has = forwardItr.Key(); has; key, has = forwardItr.Next() {
		fmt.Printf("forwardItr2; key: %v\n", key)
	}

	//Output
	// key: 1, isReplaced: false
	// key: 2, isReplaced: false
	// key: 3, isReplaced: false
	// key: 4, isReplaced: false
	// key: 5, isReplaced: false
	// Get 1; key: 1, has: true
	// Get 0; key: 0, has: false
	// GetGreater 0; key: 1, has: true
	// Select 0; key: 1, has: true
	// Select 4; key: 5, has: true
	// Select 5; key: 0, has: false
	// count keys: 3
	// Delete 1; key: 1, isDeleted: true
	// forwardItr; key: 2
	// forwardItr; key: 3
	// forwardItr; key: 4
	// forwardItr; key: 5
	// reverseItr; key: 5
	// reverseItr; key: 4
	// reverseItr; key: 3
	// reverseItr; key: 2
	// forwardItr2; key: 2
	// forwardItr2; key: 4
	// forwardItr2; key: 5
}
```

#### Augmentation

RbTreeAugmented maintains unique set of keys and invariant of node's augmented value. Supports insertion, deletion of keys in O(t * log n) time where n is number of keys in the set and t is time required to maintain node's invariant i.e updateAugmentValue time. Search operation takes O(log n) time. Can be embedded to support additional functionalities. [Interval Tree](#IntervalTree) is one such example.

##### IntervalTree:

Interval Tree maintains set of intervals and provide additional functionality (IntervalSearch) whether the given interval overlaps with any intervals in container and returning it.

If each interval interval is represented by start and end values, then Interval Tree augments maximum end value of all intervals in a subtree.

```go
type interval struct {
	start, end uint64
}
```

To implement interval tree:

1) Define a struct and embed RbTreeAugmented

  ```go
  type IntervalTree struct {
	  *orderedset.RbTreeAugmented[interval, uint64]
  }
  ```
  
2) To initialize RbTreeAugmented, provide Less method and updateAugmentValue method.

   2.1) Less method determines the order of intervals.

    ```go
    func (i interval) Less(other interval) bool {
    	if i.start < other.start {
    		return true
    	}
    	if i.start > other.start {
    		return false
    	}
    	if i.end < other.end {
    		return true
    	}
    	if i.end > other.end {
    		return false
    	}
    	return false
    }

    less := func(i1, i2 interval) bool {return i1.Less(i2)}
    ```

   2.2) updateAugmentValue method returns updated augment value. In particular, updated augment value is the maximum of {node's interval end value, node's left augmented value and node's right augmented value}. updateAugmentValue method is also provided with sentinel node to determine its augment value.

   ```go
   func getMaxEndValue(node, sentinel orderedset.BBSTNodeAugmented[interval, uint64]) uint64 {
    	if node == sentinel {
    		return 0 // Given minimum value 
    	}
    	return node.GetAugmentedValue()
    }

   updateAugmentedValue := func(node, sentinel orderedset.BBSTNodeAugmented[interval, uint64]) uint64 {
		  return max(node.GetKey().end, max(getMaxEndValue(node.GetLeftAugmented(), sentinel), getMaxEndValue(node.GetRightAugmented(), sentinel)))
	  }  
   
   ```

3) Provide functionality: IntervalSearch

   ```go
    // Finds an interval that overlaps with interval i and second value is true.
    // If there is no interval that overlaps with interval i, returns (zeroValue, false)
    // Takes O(log n) time where n is the number of intervals in the tree  
    func (it *IntervalTree) IntervalSearch(i interval) (_ interval, _ bool) {
      node := it.GetRoot()
      for node != it.GetSentinel() && !node.GetKey().Overlap(i) {
        if node.GetLeftAugmented() != it.GetSentinel() && node.GetLeftAugmented().GetAugmentedValue() >= i.start {
          node = node.GetLeftAugmented()
        } else {
          node = node.GetRightAugmented()
        }
      }
      if node == it.GetSentinel() {
        return
      }
      return node.GetKey(), true
    }
   ```

Complete Code:

```go
package main

import (
	"fmt"

	"github.com/storybehind/gocontainer/orderedset"
)

type interval struct {
	start, end uint64
}

func (i interval) Less(other interval) bool {
	if i.start < other.start {
		return true
	}
	if i.start > other.start {
		return false
	}
	if i.end < other.end {
		return true
	}
	if i.end > other.end {
		return false
	}
	return false
}

func (i interval) Overlap(other interval) bool {
	return i.start <= other.end && other.start <= i.end
}

// Augments maximum end value of all intervals in a subtree
type IntervalTree struct {
	*orderedset.RbTreeAugmented[interval, uint64]
}

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func getMaxEndValue(node, sentinel orderedset.BBSTNodeAugmented[interval, uint64]) uint64 {
	if node == sentinel {
		return 0 // Given minimum value 
	}
	return node.GetAugmentedValue()
}

func NewIntervalTree() *IntervalTree {
	less := func(i1, i2 interval) bool {return i1.Less(i2)}
	updateAugmentedValue := func(node, sentinel orderedset.BBSTNodeAugmented[interval, uint64]) uint64 {
		return max(node.GetKey().end, max(getMaxEndValue(node.GetLeftAugmented(), sentinel), getMaxEndValue(node.GetRightAugmented(), sentinel)))
	}
	return &IntervalTree {
		RbTreeAugmented: orderedset.NewRbTreeAugmented[interval, uint64](less, updateAugmentedValue),
	}
}

// Finds an interval that overlaps with interval i and second value is true.
// If there is no interval that overlaps with interval i, returns (zeroValue, false)
// Takes O(log n) time where n is the number of intervals in the tree  
func (it *IntervalTree) IntervalSearch(i interval) (_ interval, _ bool) {
	node := it.GetRoot()
	for node != it.GetSentinel() && !node.GetKey().Overlap(i) {
		if node.GetLeftAugmented() != it.GetSentinel() && node.GetLeftAugmented().GetAugmentedValue() >= i.start {
			node = node.GetLeftAugmented()
		} else {
			node = node.GetRightAugmented()
		}
	}
	if node == it.GetSentinel() {
		return
	}
	return node.GetKey(), true
}

func main() {
	// Initialize IntervalTree
	intervalTree := NewIntervalTree()
	
	// Insert intervals
	intervalTree.ReplaceOrInsert(interval{start: 1, end: 3})
	intervalTree.ReplaceOrInsert(interval{start: 5, end: 8})
	intervalTree.ReplaceOrInsert(interval{start: 6, end: 10})
	intervalTree.ReplaceOrInsert(interval{start: 8, end: 9})
	intervalTree.ReplaceOrInsert(interval{start: 15, end: 23})
	intervalTree.ReplaceOrInsert(interval{start: 16, end: 21})
	intervalTree.ReplaceOrInsert(interval{start: 17, end: 19})
	intervalTree.ReplaceOrInsert(interval{start: 19, end: 20})
	intervalTree.ReplaceOrInsert(interval{start: 25, end: 30})
	intervalTree.ReplaceOrInsert(interval{start: 26, end: 26})
	
	i, has := intervalTree.IntervalSearch(interval{start: 22, end: 25})
	fmt.Printf("interval: %v, has: %v\n", i, has)
	
	i, has = intervalTree.IntervalSearch(interval{start: 11, end: 14})
	fmt.Printf("interval: %v, has: %v\n", i, has)

	// Output
	// interval: {15 23}, has: true
	// interval: {0 0}, has: false
}
```

### orderedmap

orderedmap provides OrderedMap container which maintains key value pairs where all keys are unique. Supports insertion, deletion and search operation in O(log n) time where n is number of keys in the map. OrderedMap can be iterated in ascending (or) descending order of keys in O(n) time. Underlying set structure to store keys can be chosen based on Tag. By calling New(), it defaults to [RbTree](#RbTree) tag.

```go
package main

import (
	"fmt"
	"strings"

	"github.com/storybehind/gocontainer/orderedmap"
)

func main() {
	// Initialize orderedmap
	// Calling New(), return orderedmap with RbTree as underlying set structure by default
	om := orderedmap.New[int, string](func(k1, k2 int) bool {return k1 < k2})
	
	// To use AvlTree as underlying set structure
	// om = orderedmap.NewByTag[int, string](func(k1, k2 int) bool {return k1 < k2}, orderedmap.AvlTreeTag)

	// Insery key-value pairs
	_, isReplaced := om.ReplaceOrInsert(1, "1")
	fmt.Printf("key: 1, isReplaced: %v\n", isReplaced)
	_, isReplaced = om.ReplaceOrInsert(2, "2")
	fmt.Printf("key: 2, isReplaced: %v\n", isReplaced)
	_, isReplaced = om.ReplaceOrInsert(3, "3")
	fmt.Printf("key: 3, isReplaced: %v\n", isReplaced)
	_, isReplaced = om.ReplaceOrInsert(4, "4")
	fmt.Printf("key: 4, isReplaced: %v\n", isReplaced)
	_, isReplaced = om.ReplaceOrInsert(5, "5")
	fmt.Printf("key: 5, isReplaced: %v\n", isReplaced)

	// Search keys
	kvpair, has := om.Get(1)
	fmt.Printf("Get 1; key: %v, value: %v, has: %v\n", kvpair.GetKey(), kvpair.GetValue(), has)
	kvpair, has = om.Get(0)
	fmt.Printf("Get 0; key: %v, value: %v, has: %v\n", kvpair.GetKey(), kvpair.GetValue(), has)
	kvpair, has = om.GetGreater(0)
	fmt.Printf("GetGreater 0; key: %v, value: %v, has: %v\n", kvpair.GetKey(), kvpair.GetValue(), has)
	
	// Delete key
	kvpair, isDeleted := om.Delete(1)
	fmt.Printf("Delete 1; key: %v, value: %v, isDeleted: %v\n", kvpair.GetKey(), kvpair.GetValue(), isDeleted)
	
	// Iterate keys in ascending order
	forwardItr := om.Begin()
	for kvpair, has = forwardItr.Key(); has; kvpair, has = forwardItr.Next() {
		fmt.Printf("forwardItr; key: %v, value: %v\n", kvpair.GetKey(), kvpair.GetValue())
	}

	// Iterate keys in descending order
	reverseItr := om.Rbegin()
	for kvpair, has = reverseItr.Key(); has; kvpair, has = reverseItr.Prev() {
		fmt.Printf("reverseItr; key: %v, value: %v\n", kvpair.GetKey(), kvpair.GetValue())
	}

	// Remove if certain condition satisfies
	forwardItr = om.Begin()
	for kvpair, has = forwardItr.Key(); has; {
		if strings.Compare(kvpair.GetValue(), "3") == 0 {
			kvpair, has = forwardItr.Remove()
			continue
		}
		kvpair, has = forwardItr.Next()
	}
	// Iterate keys in ascending order
	forwardItr = om.Begin()
	for kvpair, has = forwardItr.Key(); has; kvpair, has = forwardItr.Next() {
		fmt.Printf("forwardItr2; key: %v, value: %v\n", kvpair.GetKey(), kvpair.GetValue())
	}

	// Output
	// key: 1, isReplaced: false
	// key: 2, isReplaced: false
	// key: 3, isReplaced: false
	// key: 4, isReplaced: false
	// key: 5, isReplaced: false
	// Get 1; key: 1, value: 1, has: true
	// Get 0; key: 0, value: , has: false
	// GetGreater 0; key: 1, value: 1, has: true
	// Delete 1; key: 1, value: 1, isDeleted: true
	// forwardItr; key: 2, value: 2
	// forwardItr; key: 3, value: 3
	// forwardItr; key: 4, value: 4
	// forwardItr; key: 5, value: 5
	// reverseItr; key: 5, value: 5
	// reverseItr; key: 4, value: 4
	// reverseItr; key: 3, value: 3
	// reverseItr; key: 2, value: 2
	// forwardItr2; key: 2, value: 2
	// forwardItr2; key: 4, value: 4
	// forwardItr2; key: 5, value: 5
}
```

### priorityqueue

priorityqueue provides containers in which elements with high priority are served before elements with low priority.

Containers: [Binary Heap](#BinaryHeap)

#### BinaryHeap

BinaryHeap[V] provides Push, Pop, Top, Update and Remove operations. 

Can be initialized with

1) empty queue using NewBinaryHeap. Takes O(1) time.

2) some initial values using InitBinaryHeap. Takes O(n) time where n is the number of initial values.

_Operations:_

Push(V) *BinaryHeapNode[V] - Inserts given value to the container and returns its node pointer. Takes O(log n) time where n is the number of values in the container.

Pop() V - Removes highest priority value from the container and returns it. Takes O(log n) time where n is the number of values in the container.

Top() *BinaryHeapNode[V]- Returns node's pointer to highest priority value in the container. Takes O(1) time.

Remove(*BinaryHeapNode[V]) V - Deletes given node's pointer in binary heap and returns its value. Has no effect if node is already removed. Takes O(log n) time where n is number of values in the queue.

Update(*BinaryHeapNode[V], V) - Update given node's value to new given value. Has no effect if node is already removed. Takes O(log n) time where n is number of values in the queue.

```go
package main

import (
	"fmt"

	"github.com/storybehind/gocontainer/priorityqueue"
)

func main() {
	// Initialize
	// v1 has higher priority than v2 if priorityFunc(v1, v2) return true
	pq := priorityqueue.NewBinaryHeap[int](func(v1, v2 int) bool {return v1 < v2})

	// Initialize with initial values
	// pq := priorityqueue.InitBinaryHeap[int](func(v1, v2 int) bool {return v1 < v2}, []int{1})
	
	// Push value
	ptr1 := pq.Push(1)
	ptr2 := pq.Push(2)
	fmt.Printf("MinValue: %v\n", pq.Top().GetValue())

	// Update value
	pq.Update(ptr1, 3)
	fmt.Printf("MinValue: %v\n", pq.Top().GetValue())

	// Remove
	pq.Remove(ptr2)

	// Empty queue in non decreasing order
	for pq.Len() > 0 {
		fmt.Printf("%d\n", pq.Pop())
	}

	// Output
	// MinValue: 1
	// MinValue: 2
	// 3
}
```










