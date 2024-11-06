package bmap

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"sync"
)

type Bmap[K comparable, V any] struct {
	keys       []K
	values     map[K]V
	keyIndices map[K]int
	mutex      sync.RWMutex
	sort       func(V, V) bool
	sortKeys   func(K, K) bool
	sticky     bool
	stickyKeys bool
}

type bmap[K comparable, V any] interface {
	Set(K, V)
	Get(K) (V, bool)
	Delete(K)
	Swap(K, K) error
	Sort(func(V, V) bool)
	SortAdvanced(func(V, V) bool, bool, bool)
	SortKeys(func(K, K) bool, bool, bool)
	Range() func(func(K, V) bool)
	Keys() func(func(K) bool)
	Values() func(func(V) bool)
	Len() int
	String() string
}

func (bmap *Bmap[K, V]) Set(key K, value V) {
	bmap.mutex.Lock()
	go func() {
		defer bmap.mutex.Unlock()
		if bmap.values == nil {
			bmap.values = make(map[K]V)
		}
		if bmap.keyIndices == nil {
			bmap.keyIndices = make(map[K]int)
		}
		_, ok := bmap.values[key]
		bmap.values[key] = value
		if !ok {
			if bmap.sortKeys != nil {
				i, _ := slices.BinarySearchFunc(bmap.keys, key, func(a, b K) int {
					if bmap.sortKeys(a, b) {
						return 1
					}
					return -1
				})
				bmap.keyIndices[key] = i
				bmap.keys = slices.Insert(bmap.keys, i, key)
			} else {
				bmap.keyIndices[key] = len(bmap.keys)
				bmap.keys = append(bmap.keys, key)
			}
		}
	}()
}

func (bmap *Bmap[K, V]) Get(key K) (V, bool) {
	var nilVal V
	bmap.mutex.RLock()
	if bmap.values == nil {
		bmap.mutex.RUnlock()
		return nilVal, false
	}
	value, ok := bmap.values[key]
	bmap.mutex.RUnlock()
	if !ok {
		return nilVal, false
	}
	return value, true
}

func (bmap *Bmap[K, V]) Delete(key K) {
	bmap.mutex.Lock()
	go func() {
		defer bmap.mutex.Unlock()
		if bmap.values == nil {
			return
		}
		_, ok := bmap.values[key]
		if !ok {
			return
		}
		delete(bmap.values, key)
		keyIndex := bmap.keyIndices[key]
		if keyIndex == len(bmap.keyIndices) {
			bmap.keys = bmap.keys[:keyIndex]
		} else {
			bmap.keys = append(bmap.keys[:keyIndex], bmap.keys[keyIndex+1:]...)
			for _, k := range bmap.keys[keyIndex:] {
				bmap.keyIndices[k]--
			}
		}
		delete(bmap.keyIndices, key)
	}()
}

func (bmap *Bmap[K, V]) Swap(key1, key2 K) error {
	bmap.mutex.RLock()
	if bmap.values == nil {
		bmap.mutex.RUnlock()
		return errors.New("bmap is empty")
	}
	index1, ok1 := bmap.keyIndices[key1]
	index2, ok2 := bmap.keyIndices[key2]
	bmap.mutex.RUnlock()
	if !ok1 {
		return errors.New("key 1 not found in bmap")
	}
	if !ok2 {
		return errors.New("key 2 not found in bmap")
	}
	bmap.mutex.Lock()
	go func() {
		bmap.values[key1], bmap.values[key2] = bmap.values[key2], bmap.values[key1]
		bmap.keyIndices[key1], bmap.keyIndices[key2] = index2, index1
		bmap.keys[index1], bmap.keys[index2] = bmap.keys[index2], bmap.keys[index1]
		bmap.mutex.Unlock()
	}()
	return nil
}

func (bmap *Bmap[K, V]) Sort(s func(V, V) bool) {
	bmap.mutex.Lock()
	if bmap.keys == nil {
		bmap.mutex.Unlock()
		return
	}
	go func() {
		sort.Slice(bmap.keys, func(i, j int) bool {
			return s(bmap.values[bmap.keys[i]], bmap.values[bmap.keys[j]])
		})
		for i, k := range bmap.keys {
			bmap.keyIndices[k] = i
		}
		bmap.mutex.Unlock()
	}()
}

func (bmap *Bmap[K, V]) SortAdvanced(s func(V, V) bool, stable bool, sticky bool) {
	bmap.mutex.Lock()
	if bmap.keys == nil {
		bmap.mutex.Unlock()
		return
	}
	if sticky {
		bmap.sort = s
		bmap.sortKeys = nil
	}
	go func() {
		if stable {
			sort.SliceStable(bmap.keys, func(i, j int) bool {
				return s(bmap.values[bmap.keys[i]], bmap.values[bmap.keys[j]])
			})
		} else {
			sort.Slice(bmap.keys, func(i, j int) bool {
				return s(bmap.values[bmap.keys[i]], bmap.values[bmap.keys[j]])
			})
		}
		for i, k := range bmap.keys {
			bmap.keyIndices[k] = i
		}
		bmap.mutex.Unlock()
	}()
}

func (bmap *Bmap[K, V]) SortKeys(s func(K, K) bool, stable bool, sticky bool) {
	bmap.mutex.Lock()
	if bmap.keys == nil {
		bmap.mutex.Unlock()
		return
	}
	if sticky {
		bmap.sortKeys = s
		bmap.sort = nil
	}
	go func() {
		if stable {
			sort.SliceStable(bmap.keys, func(i, j int) bool {
				return s(bmap.keys[i], bmap.keys[j])
			})
		} else {
			sort.Slice(bmap.keys, func(i, j int) bool {
				return s(bmap.keys[i], bmap.keys[j])
			})
		}
		for i, k := range bmap.keys {
			bmap.keyIndices[k] = i
		}
		bmap.mutex.Unlock()
	}()
}

func (bmap *Bmap[K, V]) Len() int {
	return len(bmap.keys)
}

func (bmap *Bmap[K, V]) Range() func(yield func(K, V) bool) {
	return func(yield func(K, V) bool) {
		bmap.mutex.RLock()
		if bmap.keys == nil {
			bmap.mutex.RUnlock()
			return
		}
		for _, key := range bmap.keys {
			val, ok := bmap.values[key]
			bmap.mutex.RUnlock()
			if ok && !yield(key, val) {
				return
			}
			bmap.mutex.RLock()
		}
		bmap.mutex.RUnlock()
	}
}

func (bmap *Bmap[K, V]) Values() func(yield func(V) bool) {
	return func(yield func(V) bool) {
		bmap.mutex.RLock()
		if bmap.keys == nil || bmap.values == nil {
			bmap.mutex.RUnlock()
			return
		}
		for _, key := range bmap.keys {
			val, ok := bmap.values[key]
			bmap.mutex.RUnlock()
			if ok && !yield(val) {
				return
			}
			bmap.mutex.RLock()
		}
		bmap.mutex.RUnlock()
	}
}

func (bmap *Bmap[K, V]) Keys() func(yield func(K) bool) {
	return func(yield func(K) bool) {
		bmap.mutex.RLock()
		if bmap.keys == nil {
			bmap.mutex.RUnlock()
			return
		}
		for _, key := range bmap.keys {
			bmap.mutex.RUnlock()
			if !yield(key) {
				return
			}
			bmap.mutex.RLock()
		}
		bmap.mutex.RUnlock()
	}
}

func (bmap *Bmap[K, V]) String() string {
	bmap.mutex.RLock()
	if bmap.keys == nil {
		bmap.mutex.RUnlock()
		return ""
	}

	var str string
	for _, key := range bmap.keys {
		val, ok := bmap.values[key]
		if ok {
			str = fmt.Sprintf("%s%s: %s\n", str, fmt.Sprint(key), fmt.Sprint(val))
		}
	}
	bmap.mutex.RUnlock()
	return str
}

func (bmap *Bmap[K, V]) Map() map[K]V {
	return bmap.values
}
