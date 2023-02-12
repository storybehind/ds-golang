package orderstatisticstree_test

import (
	"ds/bbst"
	orderstatisticstree "ds/variants/order_statistics_tree"
	"testing"
)

func TestOrderStatisticsTree(t *testing.T) {
	ost := orderstatisticstree.New(func(k1, k2 int) bool {
		return k1 < k2
	}, bbst.AvlTreeTag)
	
	assertRank := func (key int, exprank int64)  {
		rank := ost.Rank(key)
		if rank != exprank {
			t.Errorf("assertRank err; key = %d, exprank = %d, found = %d", key, exprank, rank)
		}
	}
	assertSelect := func (rank int64, expkey int, expIsPresent bool) {
		key, isPresent := ost.Select(rank)
		if isPresent != expIsPresent {
			t.Errorf("assertSelect err; isPresent = %v but expIsPresent = %v", isPresent, expIsPresent)
		}
		if isPresent && key != expkey {
			t.Errorf("assertSelect Err; key = %d, but expkey = %d", key, expkey)
		}
	}
	ost.ReplaceOrInsert(1)
	ost.ReplaceOrInsert(2)
	ost.ReplaceOrInsert(3)
	ost.ReplaceOrInsert(4)
	ost.ReplaceOrInsert(5)
	
	// ost contains {1, 2, 3, 4, 5}
	assertRank(1, 0)
	assertRank(2, 1)
	assertRank(3, 2)
	assertRank(4, 3)
	assertRank(5, 4)
	assertRank(6, -1)

	assertSelect(0, 1, true)
	assertSelect(1, 2, true)
	assertSelect(2, 3, true)
	assertSelect(3, 4, true)
	assertSelect(4, 5, true)
	assertSelect(5, 0, false)

	ost.Delete(2)
	ost.Delete(4)

	// ost contains {1, 3, 5}
	assertRank(1, 0)
	assertRank(2, -1)
	assertRank(3, 1)
	assertRank(4, -1)
	assertRank(5, 2)
	assertRank(6, -1)

	assertSelect(0, 1, true)
	assertSelect(1, 3, true)
	assertSelect(2, 5, true)
	assertSelect(3, 0, false)

}