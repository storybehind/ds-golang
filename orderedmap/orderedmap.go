package orderedmap

import "github.com/storybehind/gocontainer/orderedset"

// Container to maintain key value pairs where all keys are unique
// Supports insertion, deletion and search operation in O(log n) time where n is number of keys in the map
type OrderedMap[K, V any] struct {
	os orderedset.OrderedSetI[KeyValuePair[K, V]]
}

// Stores key and value. Use GetKey() and GetValue() to retrieve key and value respectively
type KeyValuePair[K, V any] struct {
	key   K
	value V
}

func NewKeyValuePair[K, V any](key K, value V) KeyValuePair[K, V] {
	return KeyValuePair[K, V]{
		key:   key,
		value: value,
	}
}

func (kvpair KeyValuePair[K, V]) GetKey() K {
	return kvpair.key
}

func (kvpair KeyValuePair[K, V]) GetValue() V {
	return kvpair.value
}

type Tag int

const (
	AvlTreeTag Tag = iota
	RbTreeTag
)

// Returns instance of OrderedMap
// Less method determines the order of key.
// k1 precedes k2 if and only if Less(k1, k2) return true.
// k1 equals k2 if and only if !Less(k1, k2) && !Less(k2, k1) holds true.
// By  default, underlying set structure is RbTree.
func New[K, V any](less func(k1, k2 K) bool) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		os: orderedset.NewRbTree[KeyValuePair[K, V]](func(k1, k2 KeyValuePair[K, V]) bool {
			return less(k1.key, k2.key)
		}),
	}
}

// Returns instance of OrderedMap
// Less method determines the order of key.
// k1 precedes k2 if and only if Less(k1, k2) return true.
// k1 equals k2 if and only if !Less(k1, k2) && !Less(k2, k1) holds true.
// tag specifies underlying set structure. Can be AvlTreeTag or RbTreeTag.
func NewByTag[K, V any](less func(k1, k2 K) bool, tag Tag) *OrderedMap[K, V] {
	switch tag {
	case AvlTreeTag:
		return &OrderedMap[K, V]{
			os: orderedset.NewAvlTree[KeyValuePair[K, V]](func(k1, k2 KeyValuePair[K, V]) bool {
				return less(k1.key, k2.key)
			}),
		}
	case RbTreeTag:
		return &OrderedMap[K, V]{
			os: orderedset.NewRbTree[KeyValuePair[K, V]](func(k1, k2 KeyValuePair[K, V]) bool {
				return less(k1.key, k2.key)
			}),
		}
	default:
		panic("invalid tag type")
	}
}

// Get looks for the key in the tree, returning its KeyValuePair. It returns (zeroValue, false) if unable to find that key
func (om *OrderedMap[K, V]) Get(key K) (_ KeyValuePair[K, V], _ bool) {
	return om.os.Get(KeyValuePair[K, V]{
		key: key,
	})
}

// GetGreater looks for smallest key that is strictly greater than key in the tree, returning its KeyValuePair. It returns (zeroValue, false) if unable to find that key
func (om *OrderedMap[K, V]) GetGreater(key K) (_ KeyValuePair[K, V], _ bool) {
	return om.os.GetGreater(KeyValuePair[K, V]{
		key: key,
	})
}

// GetGreaterThanOrEqual looks for smallest key that is greater than or equal to key in the tree, returning its KeyValuePair. It returns (zeroValue, false) if unable to find that key
func (om *OrderedMap[K, V]) GetGreaterThanOrEqual(key K) (_ KeyValuePair[K, V], _ bool) {
	return om.os.GetGreaterThanOrEqual(KeyValuePair[K, V]{
		key: key,
	})
}

// GetLower looks for greatest key that is strictly lower than key in the tree, returning its KeyValuePair. It returns (zeroValue, false) if unable to find that key
func (om *OrderedMap[K, V]) GetLower(key K) (_ KeyValuePair[K, V], _ bool) {
	return om.os.GetLower(KeyValuePair[K, V]{
		key: key,
	})
}

// GetLowerThanOrEqual looks for greatest key that is lower than or equal to key in the tree, returning its KeyValuePair. It returns (zeroValue, false) if unable to find that key
func (om *OrderedMap[K, V]) GetLowerThanOrEqual(key K) (_ KeyValuePair[K, V], _ bool) {
	return om.os.GetLowerThanOrEqual(KeyValuePair[K, V]{
		key: key,
	})
}

// Max returns KeyValuePair with largest key, or (zeroValue, false) if the map is empty
func (om *OrderedMap[K, V]) Max() (_ KeyValuePair[K, V], _ bool) {
	return om.os.Max()
}

// Min returns KeyValuePair with smallest key, or (zeroValue, false) if the map is empty
func (om *OrderedMap[K, V]) Min() (_ KeyValuePair[K, V], _ bool) {
	return om.os.Min()
}

// Len returns the number of keys currently in the map.
func (om *OrderedMap[K, V]) Len() int64 {
	return om.os.Len()
}

// ReplaceOrInsert adds the given key and value to the map.
// If a key in the map already equals the given one, it is removed from the map and returns its KeyValuePair, and the second return value is true.
// Otherwise, (zeroValue, false)
// panics if nil is inserted
func (om *OrderedMap[K, V]) ReplaceOrInsert(key K, value V) (_ KeyValuePair[K, V], _ bool) {
	return om.os.ReplaceOrInsert(KeyValuePair[K, V]{
		key:   key,
		value: value,
	})
}

// Delete the key in the map and return its KeyValuePair, and the second return value is true.
// If key is not found in the tree, returns (zeroValue, false)
func (om *OrderedMap[K, V]) Delete(key K) (_ KeyValuePair[K, V], _ bool) {
	return om.os.Delete(KeyValuePair[K, V]{
		key: key,
	})
}

// Delete the maximum key in the map and return its KeyValuePair.
// On calling empty tree, returns (zeroValue, false)
func (om *OrderedMap[K, V]) DeleteMax() (_ KeyValuePair[K, V], _ bool) {
	return om.os.DeleteMax()
}

// Delete the minimum key in the map and return its KeyValuePair.
// On calling empty tree, returns (zeroValue, false)
func (om *OrderedMap[K, V]) DeleteMin() (_ KeyValuePair[K, V], _ bool) {
	return om.os.DeleteMin()
}

type OrderedMapIterator[K, V any] struct {
	forwardIterator orderedset.OrderedSetForwardIterator[KeyValuePair[K, V]]
}

// Returns an iterator pointing to least key in the map.
// Used to iterate keys in the ascending order.
func (om *OrderedMap[K, V]) Begin() *OrderedMapIterator[K, V] {
	return &OrderedMapIterator[K, V]{
		forwardIterator: om.os.Begin(),
	}
}

// Calling Next() moves the iterator to the next greater node and returns its KeyValuePair
// If Next() is called on last key(or greatest key), it returns (zeroValue, false)
func (omItr *OrderedMapIterator[K, V]) Next() (_ KeyValuePair[K, V], _ bool) {
	return omItr.forwardIterator.Next()
}

// Returns the KeyValuePair pointed by iterator. Returns (zeroValue, false) if this is called on empty map or an iterator has completed traversing all the keys
func (omItr *OrderedMapIterator[K, V]) Key() (_ KeyValuePair[K, V], _ bool) {
	return omItr.forwardIterator.Key()
}

// Deletes the key the pointed by iterator, moves the iterator to next greater node.
// Returns the next greater KeyValuePair if it's present. Otherwise, returns (zeroValue, false)
// panics on calling Remove() in empty map or an iterator has completed traversing all the keys
func (omItr *OrderedMapIterator[K, V]) Remove() (_ KeyValuePair[K, V], _ bool) {
	return omItr.forwardIterator.Remove()
}

type ReverseOrderedMapIterator[K, V any] struct {
	reverseIterator orderedset.OrderedSetReverseIterator[KeyValuePair[K, V]]
}

// Returns an reverse iterator pointing to greatest key node in the tree
// Used to iterate keys in the descending order
func (om *OrderedMap[K, V]) Rbegin() *ReverseOrderedMapIterator[K, V] {
	return &ReverseOrderedMapIterator[K, V]{
		reverseIterator: om.os.Rbegin(),
	}
}

// Calling Prev() moves the reverse iterator to the next smaller node and returns its KeyValuePair
// If Prev() is called on last key (or smallest key), it returns (zeroValue, false)
func (omRitr *ReverseOrderedMapIterator[K, V]) Prev() (_ KeyValuePair[K, V], _ bool) {
	return omRitr.reverseIterator.Prev()
}

// Returns the key pointed by reverse iterator. Returns (zeroValue, false) if this is called on empty map or an iterator has completed traversing all the keys
func (omRitr *ReverseOrderedMapIterator[K, V]) Key() (_ KeyValuePair[K, V], _ bool) {
	return omRitr.reverseIterator.Key()
}

// Deletes the key the pointed by reverse iterator, moves the reverse iterator to next smaller key.
// Returns the next smaller KeyValuePair if it's present. Otherwise, returns (zeroValue, false)
// panics on calling Remove() in empty map or an iterator has completed traversing all the keys
func (omRitr *ReverseOrderedMapIterator[K, V]) Remove() (_ KeyValuePair[K, V], _ bool) {
	return omRitr.reverseIterator.Remove()
}
