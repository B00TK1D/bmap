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

func TestOrdered(t *testing.T) {
	bmap := Type[int, int]{}.NewAsync()
	for i := range 100000 {
		want := i
		bmap.Set(i, want)
		got, _ := bmap.Get(i)
		if want != got {
			t.Errorf("got %d, wanted %d", got, want)
		}
	}
}
