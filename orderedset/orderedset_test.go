package orderedset_test

import (
	"testing"

	orderedset "github.com/storybehind/gocontainer/orderedset"
)

func testOrderedSet(t *testing.T, os orderedset.OrderedSet[int]) {
	// insert results
	var prevKey int
	var isReplaced bool
	// search results
	var searchedKey int
	var has bool
	// delete results
	var deletedKey int
	var isDeleted bool

	assertLen := func(expLen int64) {
		if os.Len() != expLen {
			t.Errorf("os.Len() = %d, expLen = %d", os.Len(), expLen)
		}
	}
	assertInsert := func() {
		if prevKey != 0 || isReplaced {
			t.Errorf("assertInsert err; prevKey = %d, isReplaced = %v", prevKey, isReplaced)
		}
	}
	assertReplace := func(expPrevKey int) {
		if prevKey != expPrevKey || !isReplaced {
			t.Errorf("assertReplace err; prevKey = %d, expPrevKey = %d, isReplaced = %v", prevKey, expPrevKey, isReplaced)
		}
	}
	assertSearch := func(expSearchedKey int, expHas bool) {
		if expSearchedKey != searchedKey || has != expHas {
			t.Errorf("assertSearch err; searchedKey = %d, expSearchedKey = %d, has = %v, expHas = %v", searchedKey, expSearchedKey, has, expHas)
		}
	}
	assertDelete := func(expDeletedKey int, expIsDeleted bool) {
		if deletedKey != expDeletedKey || isDeleted != expIsDeleted {
			t.Errorf("assertDelete err; deletedKey = %d, expDeletedKey = %d,  isDeleted = %v, expIsDeleted = %v", deletedKey, expDeletedKey, isDeleted, expIsDeleted)
		}
	}
	assertEmptyTree := func() {
		assertLen(0)
		deletedKey, isDeleted = os.Delete(0)
		assertDelete(0, false)
		deletedKey, isDeleted = os.DeleteMin()
		assertDelete(0, false)
		deletedKey, isDeleted = os.DeleteMax()
		assertDelete(0, false)

		searchedKey, has = os.Get(0)
		assertSearch(0, false)
		searchedKey, has = os.GetGreater(0)
		assertSearch(0, false)
		searchedKey, has = os.GetGreaterThanOrEqual(0)
		assertSearch(0, false)
		searchedKey, has = os.GetLower(0)
		assertSearch(0, false)
		searchedKey, has = os.GetLowerThanOrEqual(0)
		assertSearch(0, false)
	}

	// insert keys to the tree. Every key must be unique and not present already in the tree
	insertElements := func(keys []int) {
		sz := os.Len()
		for _, key := range keys {
			prevKey, isReplaced = os.ReplaceOrInsert(key)
			assertInsert()
			sz++
			assertLen(sz)
		}
	}

	// assert search conditions in the original tree i.e 2,3,5,7,11
	assertOriginalTreeSearch := func() {
		searchedKey, has = os.Get(1)
		assertSearch(0, false)
		searchedKey, has = os.Get(2)
		assertSearch(2, true)
		searchedKey, has = os.GetGreater(0)
		assertSearch(2, true)
		searchedKey, has = os.GetGreater(2)
		assertSearch(3, true)
		searchedKey, has = os.GetGreater(13)
		assertSearch(0, false)
		searchedKey, has = os.GetGreaterThanOrEqual(0)
		assertSearch(2, true)
		searchedKey, has = os.GetGreaterThanOrEqual(2)
		assertSearch(2, true)
		searchedKey, has = os.GetGreaterThanOrEqual(13)
		assertSearch(0, false)
		searchedKey, has = os.GetLower(4)
		assertSearch(3, true)
		searchedKey, has = os.GetLower(5)
		assertSearch(3, true)
		searchedKey, has = os.GetLower(1)
		assertSearch(0, false)
		searchedKey, has = os.GetLowerThanOrEqual(4)
		assertSearch(3, true)
		searchedKey, has = os.GetLowerThanOrEqual(5)
		assertSearch(5, true)
		searchedKey, has = os.GetLowerThanOrEqual(1)
		assertSearch(0, false)
		searchedKey, has = os.Max()
		assertSearch(11, true)
		searchedKey, has = os.Min()
		assertSearch(2, true)
	}

	// delete keys in the tree. All keys should be present
	deleteElements := func(keys []int) {
		initLen := os.Len()
		for _, key := range keys {
			deletedKey, isDeleted = os.Delete(key)
			assertDelete(key, true)
			initLen--
			assertLen(initLen)
		}
	}

	// tree must contain exp
	assertGet := func(exp int) {
		if _, has := os.Get(exp); !has {
			t.Errorf("Get err; %d should be present in tree", exp)
		}
	}

	// empty tree assert
	assertEmptyTree()

	// insert keys: 2,3,5,7,11
	insertElements([]int{5, 2, 3, 11, 7})
	// tree should contain 2,3,5,7,11

	// insert key: 5 again
	prevKey, isReplaced = os.ReplaceOrInsert(5)
	assertReplace(5)

	// original tree contains keys: 2,3,5,7,11
	assertOriginalTreeSearch()

	// delete key: 1 which is not present in original tree
	deletedKey, isDeleted = os.Delete(1)
	assertDelete(0, false)
	assertLen(5)

	//delete keys: 2,5,11
	deleteElements([]int{5, 2, 11})
	assertGet(3)
	assertGet(7)

	// insert keys: 2,5,11 and bring back original keys
	insertElements([]int{2, 5, 11})
	assertOriginalTreeSearch()

	// delete all keys
	deleteElements([]int{2, 3, 5, 7, 11})
	assertEmptyTree()

	//bring back original tree from empty tree
	insertElements([]int{5, 2, 3, 11, 7})
	assertOriginalTreeSearch()

	// delete all keys
	deleteElements([]int{2, 3, 5, 7, 11})
}

func testOrderedSetForwardIterator(t *testing.T, osi orderedset.OrderedSetI[int]) {
	itr := osi.Begin()
	if _, has := itr.Next(); has {
		t.Errorf("on calling next() empty tree; expected: false, but found true")
	}

	assertNext := func(el, expEl int, has, expHas bool) {
		if has != expHas {
			t.Errorf("on calling next(); expected: %v, but found: %v", expHas, has)
		}
		if el != expEl {
			t.Errorf("on calling next(); expected element: %v, but found: %v", expEl, el)
		}
	}
	L := 5
	for i := 1; i <= L; i++ {
		osi.ReplaceOrInsert(i)
	}
	expEl := 1
	itr = osi.Begin()
	for el, has := itr.Key(); has; el, has = itr.Next() {
		assertNext(el, int(expEl), has, true)
		expEl++
	}

	itr = osi.Begin()
	for el, has := itr.Key(); has; {
		if el == 1 || el == 2 || el == 5 {
			el, has = itr.Remove()
			continue
		}
		el, has = itr.Next()
	}
	if osi.Len() != 2 {
		t.Errorf("Exp Len: 2, but found: %v", osi.Len())
	}
	_, has3 := osi.Get(3)
	if !has3 {
		t.Errorf("Expected key: 3 to be present")
	}
	_, has4 := osi.Get(4)
	if !has4 {
		t.Errorf("Expected key: 4 to be present")
	}
}

func testOrderedSetReverseIterator(t *testing.T, osi orderedset.OrderedSetI[int]) {
	ritr := osi.Rbegin()
	if _, has := ritr.Prev(); has {
		t.Errorf("on calling prev() empty tree; expected: false, but found true")
	}

	assertPrev := func(el, expEl int, has, expHas bool) {
		if has != expHas {
			t.Errorf("on calling prev(); expected: %v, but found: %v", expHas, has)
		}
		if el != expEl {
			t.Errorf("on calling prev(); expected element: %v, but found: %v", expEl, el)
		}
	}
	L := 5
	for i := 1; i <= L; i++ {
		osi.ReplaceOrInsert(i)
	}
	expEl := L
	ritr = osi.Rbegin()
	for el, has := ritr.Key(); has; el, has = ritr.Prev() {
		assertPrev(el, int(expEl), has, true)
		expEl--
	}

	ritr = osi.Rbegin()
	for el, has := ritr.Key(); has; {
		if el == 1 || el == 2 || el == 5 {
			el, has = ritr.Remove()
			t.Logf("el: %v, has: %v", el, has)
			continue
		}
		el, has = ritr.Prev()
		t.Logf("el: %v, has: %v", el, has)
	}
	if osi.Len() != 2 {
		t.Errorf("Exp Len: 2, but found: %v", osi.Len())
	}
	_, has3 := osi.Get(3)
	if !has3 {
		t.Errorf("Expected key: 3 to be present")
	}
	_, has4 := osi.Get(4)
	if !has4 {
		t.Errorf("Expected key: 4 to be present")
	}
}
