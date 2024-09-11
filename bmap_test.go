package bmap

import "testing"

func TestSet(t *testing.T) {
	bmap := Type[int, int]{}.New()
	for i := range 1000 {
		go func() {
			for j := range 100 {
				want := j
				bmap.Set(i, want)
				got, _ := bmap.Get(i)
				if want != got {
					t.Errorf("got %d, wanted %d", got, want)
				}
				bmap.Sort(func(i, j int) bool {
					return i > j
				})
				got, _ = bmap.Get(i)
				if want != got {
					t.Errorf("got %d, wanted %d", got, want)
				}
			}
		}()
	}
}

func TestAsync(t *testing.T) {
	bmap := Type[int, int]{}.NewAsync()
	for i := range 100000 {
		want := i*13
		bmap.Set(i, want)
		got, _ := bmap.Get(i)
		if want != got {
			t.Errorf("got %d, wanted %d", got, want)
		}
	}
}

func TestDelete(t *testing.T) {
	bmap := Type[int, int]{}.New()
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
	bmap := Type[int, int]{}.New()
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
