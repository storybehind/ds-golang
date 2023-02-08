package bbst_test

import (
	"ds/bbst"
	"testing"
)

func TestAvlTree(t *testing.T) {
	avlTree := bbst.NewAvlTreeByLess(func(k1, k2 int32) bool {
		return k1 < k2
	})

	// insert results
	var prevKey int32
	var isReplaced bool
	// search results
	var searchedKey int32
	var has bool
	// delete results
	var deletedKey int32
	var isDeleted bool

	assertLen := func (expLen int64) {
		if avlTree.Len() != expLen {
			t.Errorf("avlTree.Len() = %d, expLen = %d", avlTree.Len(), expLen)
		}
	}
	assertInsert := func ()  {
		if prevKey != 0 || isReplaced {
			t.Errorf("assertInsert err; prevKey = %d, isReplaced = %v", prevKey, isReplaced)
		}
	}
	assertReplace := func (expPrevKey int32)  {
		if prevKey != expPrevKey || !isReplaced {
			t.Errorf("assertReplace err; prevKey = %d, expPrevKey = %d, isReplaced = %v", prevKey, expPrevKey, isReplaced)
		}
	}
	assertSearch := func (expSearchedKey int32, expHas bool)  {
		if expSearchedKey != searchedKey || has != expHas {
			t.Errorf("assertSearch err; searchedKey = %d, expSearchedKey = %d, has = %v, expHas = %v", searchedKey, expSearchedKey, has, expHas)
		}
	}
	assertDelete := func (expDeletedKey int32, expIsDeleted bool)  {
		if deletedKey != expDeletedKey || isDeleted != expIsDeleted {
			t.Errorf("assertDelete err; deletedKey = %d, expDeletedKey = %d,  isDeleted = %v, expIsDeleted = %v", deletedKey, expDeletedKey, isDeleted, expIsDeleted)
		}
	}
	assertEmptyTree := func ()  {
		assertLen(0)
		deletedKey, isDeleted = avlTree.Delete(0)
		assertDelete(0, false)
		deletedKey, isDeleted = avlTree.DeleteMin()
		assertDelete(0, false)
		deletedKey, isDeleted = avlTree.DeleteMax()
		assertDelete(0, false)
		
		searchedKey, has = avlTree.Get(0)
		assertSearch(0, false)
		searchedKey, has = avlTree.GetGreater(0)
		assertSearch(0, false)
		searchedKey, has = avlTree.GetGreaterThanOrEqual(0)
		assertSearch(0, false)
		searchedKey, has = avlTree.GetLower(0)
		assertSearch(0, false)
		searchedKey, has = avlTree.GetLowerThanOrEqual(0)
		assertSearch(0, false)
		if avlTree.Has(0) {
			t.Errorf("Has search succeeds in empty tree err")
		}
	}

	// insert keys to the tree. Every key must be unique and not present already in the tree
	insertElements := func (keys []int32)  {
		sz := avlTree.Len()
		for _ , key := range keys {
			prevKey, isReplaced  =avlTree.ReplaceOrInsert(key)
			assertInsert()
			sz++
			assertLen(sz)
		}
	}

	// assert search conditions in the original tree i.e 2,3,5,7,11
	assertOriginalTreeSearch := func () {
		searchedKey, has = avlTree.Get(1)
		assertSearch(0, false)
		searchedKey, has = avlTree.Get(2)
		assertSearch(2, true)
		searchedKey, has = avlTree.GetGreater(0)
		assertSearch(2, true)
		searchedKey, has = avlTree.GetGreater(2)
		assertSearch(3, true)
		searchedKey, has = avlTree.GetGreater(13)
		assertSearch(0, false)
		searchedKey, has = avlTree.GetGreaterThanOrEqual(0)
		assertSearch(2, true)
		searchedKey, has = avlTree.GetGreaterThanOrEqual(2)
		assertSearch(2, true)
		searchedKey, has = avlTree.GetGreaterThanOrEqual(13)
		assertSearch(0, false)
		searchedKey, has = avlTree.GetLower(4)
		assertSearch(3, true)
		searchedKey, has = avlTree.GetLower(5)
		assertSearch(3, true)
		searchedKey, has = avlTree.GetLower(1)
		assertSearch(0, false)
		searchedKey, has = avlTree.GetLowerThanOrEqual(4)
		assertSearch(3, true)
		searchedKey, has = avlTree.GetLowerThanOrEqual(5)
		assertSearch(5, true)
		searchedKey, has = avlTree.GetLowerThanOrEqual(1)
		assertSearch(0, false)
		searchedKey, has = avlTree.Max()
		assertSearch(11, true)
		searchedKey, has = avlTree.Min()
		assertSearch(2, true)
	}

	// delete keys in the tree. All keys should be present
	deleteElements := func (keys []int32)  {
		initLen := avlTree.Len()
		for _, key := range keys {
			deletedKey, isDeleted = avlTree.Delete(key)
			assertDelete(deletedKey, isDeleted)
			initLen--
			assertLen(initLen)
		}
	}

	// tree must contain exp
	assertHas := func (exp int32)  {
		if !avlTree.Has(exp) {
			t.Errorf("Has err; %d should be present in tree", exp)
		}
	}

	assertEmptyTree()

	insertElements([]int32{5,2,3,11,7})
	prevKey, isReplaced = avlTree.ReplaceOrInsert(5)
	assertReplace(5)
	assertOriginalTreeSearch()

	deletedKey, isDeleted = avlTree.Delete(1)
	assertDelete(0, false)
	assertLen(5)
	deleteElements([]int32{5, 2, 11})
	assertHas(3)
	assertHas(7)

	// bring back original keys and search
	insertElements([]int32{2,5,11})
	assertOriginalTreeSearch()

	// delete all keys and bring up original keys
	deleteElements([]int32{2,3,5,7,11})
	assertEmptyTree()
	insertElements([]int32{5,2,3,11,7})
	assertOriginalTreeSearch()
}
