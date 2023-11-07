package orderedmap_test

import (
	"testing"

	"github.com/storybehind/gocontainer/orderedmap"
)

func testOrderedMap(t *testing.T, om *orderedmap.OrderedMap[int, int]) {

	// insert results
	var prevKeyValuePair orderedmap.KeyValuePair[int, int]
	var isReplaced bool
	// search results
	var searchedKeyValuePair orderedmap.KeyValuePair[int, int]
	var has bool
	// delete results
	var deletedKeyValuePair orderedmap.KeyValuePair[int, int]
	var isDeleted bool

	assertLen := func(expLen int64) {
		if om.Len() != expLen {
			t.Errorf("os.Len() = %d, expLen = %d", om.Len(), expLen)
		}
	}
	assertInsert := func() {
		if isReplaced {
			t.Errorf("assertInsert err; prevKeyValuePair = %v, isReplaced = %v", prevKeyValuePair, isReplaced)
		}
	}
	assertReplace := func(expPrevKeyValuePair orderedmap.KeyValuePair[int, int]) {
		if prevKeyValuePair.GetKey() != expPrevKeyValuePair.GetKey() || prevKeyValuePair.GetValue() != expPrevKeyValuePair.GetValue() || !isReplaced {
			t.Errorf("assertReplace err; prevKeyValuePair = %v, expPrevKeyValuePair = %v, isReplaced = %v", prevKeyValuePair, expPrevKeyValuePair, isReplaced)
		}
	}
	assertSearch := func(expSearchedKeyValuePair orderedmap.KeyValuePair[int, int], expHas bool) {
		if expSearchedKeyValuePair.GetKey() != searchedKeyValuePair.GetKey() || expSearchedKeyValuePair.GetValue() != searchedKeyValuePair.GetValue() || has != expHas {
			t.Errorf("assertSearch err; searchedKeyValuePair = %v, expSearchedKeyValuePair = %v, has = %v, expHas = %v", searchedKeyValuePair, expSearchedKeyValuePair, has, expHas)
		}
	}
	assertDelete := func(expDeletedKeyValuePair orderedmap.KeyValuePair[int, int], expIsDeleted bool) {
		if deletedKeyValuePair.GetKey() != expDeletedKeyValuePair.GetKey() || deletedKeyValuePair.GetValue() != expDeletedKeyValuePair.GetValue() || isDeleted != expIsDeleted {
			t.Errorf("assertDelete err; deletedKeyValuePair = %v, expDeletedKeyValuePair = %v,  isDeleted = %v, expDeletedKeyValuePair = %v", deletedKeyValuePair, expDeletedKeyValuePair, isDeleted, expIsDeleted)
		}
	}
	var zeroKeyValuePair orderedmap.KeyValuePair[int, int]
	assertEmptyMap := func() {
		assertLen(0)
		deletedKeyValuePair, isDeleted = om.Delete(0)
		assertDelete(zeroKeyValuePair, false)
		deletedKeyValuePair, isDeleted = om.DeleteMin()
		assertDelete(zeroKeyValuePair, false)
		deletedKeyValuePair, isDeleted = om.DeleteMax()
		assertDelete(zeroKeyValuePair, false)

		searchedKeyValuePair, has = om.Get(0)
		assertSearch(zeroKeyValuePair, false)
		searchedKeyValuePair, has = om.GetGreater(0)
		assertSearch(zeroKeyValuePair, false)
		searchedKeyValuePair, has = om.GetGreaterThanOrEqual(0)
		assertSearch(zeroKeyValuePair, false)
		searchedKeyValuePair, has = om.GetLower(0)
		assertSearch(zeroKeyValuePair, false)
		searchedKeyValuePair, has = om.GetLowerThanOrEqual(0)
		assertSearch(zeroKeyValuePair, false)
	}

	// insert keys and its corresponding value to the map. Every key must be unique and not present already in the map
	insertElements := func(keys, values []int) {
		sz := om.Len()
		for i := range keys {
			prevKeyValuePair, isReplaced = om.ReplaceOrInsert(keys[i], values[i])
			assertInsert()
			sz++
			assertLen(sz)
		}
	}

	// assert search conditions in the original map i.e {2,20},{3,30},{5,50},{7,70},{11,110}
	assertOriginalMapSearch := func() {
		searchedKeyValuePair, has = om.Get(1)
		assertSearch(zeroKeyValuePair, false)
		searchedKeyValuePair, has = om.Get(2)
		assertSearch(orderedmap.NewKeyValuePair[int, int](2, 20), true)
		searchedKeyValuePair, has = om.GetGreater(0)
		assertSearch(orderedmap.NewKeyValuePair[int, int](2, 20), true)
		searchedKeyValuePair, has = om.GetGreater(2)
		assertSearch(orderedmap.NewKeyValuePair[int, int](3, 30), true)
		searchedKeyValuePair, has = om.GetGreater(13)
		assertSearch(zeroKeyValuePair, false)
		searchedKeyValuePair, has = om.GetGreaterThanOrEqual(0)
		assertSearch(orderedmap.NewKeyValuePair[int, int](2, 20), true)
		searchedKeyValuePair, has = om.GetGreaterThanOrEqual(2)
		assertSearch(orderedmap.NewKeyValuePair[int, int](2, 20), true)
		searchedKeyValuePair, has = om.GetGreaterThanOrEqual(13)
		assertSearch(zeroKeyValuePair, false)
		searchedKeyValuePair, has = om.GetLower(4)
		assertSearch(orderedmap.NewKeyValuePair[int, int](3, 30), true)
		searchedKeyValuePair, has = om.GetLower(5)
		assertSearch(orderedmap.NewKeyValuePair[int, int](3, 30), true)
		searchedKeyValuePair, has = om.GetLower(1)
		assertSearch(zeroKeyValuePair, false)
		searchedKeyValuePair, has = om.GetLowerThanOrEqual(4)
		assertSearch(orderedmap.NewKeyValuePair[int, int](3, 30), true)
		searchedKeyValuePair, has = om.GetLowerThanOrEqual(5)
		assertSearch(orderedmap.NewKeyValuePair[int, int](5, 50), true)
		searchedKeyValuePair, has = om.GetLowerThanOrEqual(1)
		assertSearch(zeroKeyValuePair, false)
		searchedKeyValuePair, has = om.Max()
		assertSearch(orderedmap.NewKeyValuePair[int, int](11, 110), true)
		searchedKeyValuePair, has = om.Min()
		assertSearch(orderedmap.NewKeyValuePair[int, int](2, 20), true)
	}

	// delete keys in the tree. All keys should be present
	deleteElements := func(keys []int) {
		initLen := om.Len()
		for _, key := range keys {
			deletedKeyValuePair, isDeleted = om.Delete(key)
			assertDelete(orderedmap.NewKeyValuePair[int, int](key, 10*key), true)
			initLen--
			assertLen(initLen)
		}
	}

	// map must contain exp key
	assertGet := func(exp int) {
		if _, has := om.Get(exp); !has {
			t.Errorf("Get err; %d should be present in map", exp)
		}
	}

	// empty tree assert
	assertEmptyMap()

	// insert pairs in map: {2,20},{3,30},{5,50},{7,70},{11,110}
	insertElements([]int{5, 2, 3, 11, 7}, []int{50, 20, 30, 110, 70})
	// tree should contain 2,3,5,7,11

	// insert pair: {5,50} again
	prevKeyValuePair, isReplaced = om.ReplaceOrInsert(5, 50)
	assertReplace(orderedmap.NewKeyValuePair[int, int](5, 50))

	// original map contains pairs: {2,20},{3,30},{5,50},{7,70},{11,110}
	assertOriginalMapSearch()

	// delete key: 1 which is not present in original map
	deletedKeyValuePair, isDeleted = om.Delete(1)
	assertDelete(zeroKeyValuePair, false)
	assertLen(5)

	//delete keys: 2,5,11
	deleteElements([]int{5, 2, 11})
	assertGet(3)
	assertGet(7)

	// insert pairs: {2,20},{5,50},{11,110} and bring back original map
	insertElements([]int{2, 5, 11}, []int{20, 50, 110})
	assertOriginalMapSearch()

	// delete all keys
	deleteElements([]int{2, 3, 5, 7, 11})
	assertEmptyMap()

	//bring back original tree from empty tree
	insertElements([]int{5, 2, 3, 11, 7}, []int{50, 20, 30, 110, 70})
	assertOriginalMapSearch()

	// delete all keys
	deleteElements([]int{2, 3, 5, 7, 11})
}

func TestOrderedMapRbTreeTag(t *testing.T) {
	om := orderedmap.New[int, int](func(k1, k2 int) bool {
		return k1 < k2
	})
	testOrderedMap(t, om)
}

func TestOrderedMapAvlTreeTag(t *testing.T) {
	om := orderedmap.NewByTag[int, int](func(k1, k2 int) bool {
		return k1 < k2
	}, orderedmap.AvlTreeTag)
	testOrderedMap(t, om)
}
