package bmap

import (
	"sort"
	"sync"
	"testing"
)

func TestSet(t *testing.T) {
	bmap := Bmap[int, int]{}
	for i := range 10 {
		go func() {
			for j := range 10 {
				want := j
				bmap.Set(i, want)
				got, _ := bmap.Get(i)
				if want != got {
					t.Errorf("got %d, wanted %d", got, want)
				}
				got, _ = bmap.Get(i)
				if want != got {
					t.Errorf("got %d, wanted %d", got, want)
				}
			}
		}()
	}
}

func TestAsync(t *testing.T) {
	bmap := Bmap[int, int]{}
	for i := range 100000 {
		want := i * 13
		bmap.Set(i, want)
		got, _ := bmap.Get(i)
		if want != got {
			t.Errorf("got %d, wanted %d", got, want)
		}
	}
}

func TestDelete(t *testing.T) {
	bmap := Bmap[int, int]{}
	want := 847392
	bmap.Set(1001, want)
	for i := range 1000 {
		go bmap.Set(i, i)
		go bmap.Set(i*2, i*3)
		go bmap.Delete(i)
		go bmap.Delete(i * 2)
	}
	got, _ := bmap.Get(1001)
	if want != got {
		t.Errorf("got %d, wanted %d", got, want)
	}
}

func TestSwap(t *testing.T) {
	bmap := Bmap[int, int]{}
	want := 107834
	bmap.Set(-1, want)
	for i := range 1000 {
		bmap.Set(i, i)
		bmap.Swap(i-1, i)
	}
	got, _ := bmap.Get(999)
	if want != got {
		t.Errorf("got %d, wanted %d", got, want)
	}
}

func BenchmarkSetBmap(b *testing.B) {
	bmap := Bmap[int, int]{}
	for i := range b.N {
		bmap.Set(i, i)
	}
}

func BenchmarkInitSyncmap(b *testing.B) {
	smap := sync.Map{}
	for i := range b.N {
		smap.Store(i, i)
	}
}

func BenchmarkGetBmap(b *testing.B) {
	bmap := Bmap[int, int]{}
	bmap.Set(0, 0)
	for range b.N {
		bmap.Get(0)
	}
}

func BenchmarkGetSyncmap(b *testing.B) {
	smap := sync.Map{}
	smap.Store(0, 0)
	for range b.N {
		smap.Load(0)
	}
}

func BenchmarkMultiBmap(b *testing.B) {
	bmaps := Bmap[int, *Bmap[int, int]]{}

	for i := range b.N {
		for j := range b.N {
			bmap := Bmap[int, int]{}
			bmap.Set(i, j)
			bmaps.Set(i, &bmap)
		}
		for j := range b.N {
			bmap, _ := bmaps.Get(i)
			bmap.Get(j)
		}
		for j := range b.N {
			bmap, _ := bmaps.Get(i)
			bmap.Delete(j)
		}
	}
}

func BenchmarkMultiSyncmap(b *testing.B) {
	smaps := sync.Map{}

	for i := range b.N {
		for j := range b.N {
			smap := sync.Map{}
			smap.Store(j, j)
			smaps.Store(i, &smap)
		}
		for j := range b.N {
			smap, _ := smaps.Load(i)
			smap.(*sync.Map).Load(j)
		}
		for j := range b.N {
			smap, _ := smaps.Load(i)
			smap.(*sync.Map).Delete(j)
		}
	}
}

func BenchmarkSortMultiBmap(b *testing.B) {
	bmap1 := Bmap[int, int]{}
	bmap2 := Bmap[int, int]{}
	bmap3 := Bmap[int, int]{}

	for i := range b.N {
		want := i
		bmap1.Set(i, want)
		bmap2.Set(i, want)
		bmap3.Set(i, want)
		got, _ := bmap1.Get(i)
		if want != got {
			b.Errorf("got %d, wanted %d", got, want)
		}
		got, _ = bmap2.Get(i)
		if want != got {
			b.Errorf("got %d, wanted %d", got, want)
		}
		got, _ = bmap3.Get(i)
		if want != got {
			b.Errorf("got %d, wanted %d", got, want)
		}
		bmap1.Sort(func(i, j int) bool {
			return i > j
		})
		bmap2.Sort(func(i, j int) bool {
			return i > j
		})
		bmap3.Sort(func(i, j int) bool {
			return i > j
		})
		got, _ = bmap1.Get(i)
		if want != got {
			b.Errorf("got %d, wanted %d", got, want)
		}
		got, _ = bmap2.Get(i)
		if want != got {
			b.Errorf("got %d, wanted %d", got, want)
		}
		got, _ = bmap3.Get(i)
		if want != got {
			b.Errorf("got %d, wanted %d", got, want)
		}
	}
	for i := range 1000 {
		bmap1.Delete(i)
		bmap2.Delete(i)
		bmap3.Delete(i)
	}
}

func BenchmarkSortMultiSyncmap(b *testing.B) {
	smap1 := sync.Map{}
	smap2 := sync.Map{}
	smap3 := sync.Map{}

	for i := range b.N {
		want := i
		smap1.Store(i, want)
		smap2.Store(i, want)
		smap3.Store(i, want)
		got, _ := smap1.Load(i)
		if got != nil && want != got {
			b.Errorf("got %d, wanted %d", got, want)
		}
		got, _ = smap2.Load(i)
		if got != nil && want != got {
			b.Errorf("got %d, wanted %d", got, want)
		}
		got, _ = smap3.Load(i)
		if got != nil && want != got {
			b.Errorf("got %d, wanted %d", got, want)
		}

		// Manually sort the sync maps
		keys := []int{}
		smap1.Range(func(k, v interface{}) bool {
			keys = append(keys, k.(int))
			return true
		})
		sort.Ints(keys)
		for _, k := range keys {
			got, _ = smap1.Load(k)
		}
		keys = []int{}
		smap2.Range(func(k, v interface{}) bool {
			keys = append(keys, k.(int))
			return true
		})
		sort.Ints(keys)
		for _, k := range keys {
			got, _ = smap2.Load(k)
		}
		keys = []int{}
		smap3.Range(func(k, v interface{}) bool {
			keys = append(keys, k.(int))
			return true
		})
		sort.Ints(keys)
		for _, k := range keys {
			got, _ = smap3.Load(k)
		}
	}
	for i := range 1000 {
		smap1.Delete(i)
		smap2.Delete(i)
		smap3.Delete(i)
	}
}
