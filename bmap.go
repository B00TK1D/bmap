package bmap

import (
	"errors"
	"fmt"
	"iter"
	"sort"
	"sync"
)

type bmap[K comparable, V any] struct {
	keys       []K
	values     map[K]V
	keyIndices map[K]int
	mutex      *sync.RWMutex
}

type Type[K comparable, V any] struct{}

func (t Type[K, V]) New() bmap[K, V] {
	return bmap[K, V]{
		mutex:      &sync.RWMutex{},
		values:     map[K]V{},
		keyIndices: map[K]int{},
	}
}

func (bmap *bmap[K, V]) Set(key K, value V) {
	bmap.mutex.RLock()
	_, ok := bmap.values[key]
	bmap.mutex.RUnlock()
	bmap.mutex.Lock()
	defer bmap.mutex.Unlock()
	bmap.values[key] = value
	if !ok {
		bmap.keys = append(bmap.keys, key)
		bmap.keyIndices[key] = len(bmap.keys)
	}
}

func (bmap *bmap[K, V]) Get(key K) (V, bool) {
	var nilVal V
	bmap.mutex.RLock()
	value, ok := bmap.values[key]
	bmap.mutex.RUnlock()
	if !ok {
		fmt.Println("Key not found")
		return nilVal, false
	}
	return value, true
}

func (bmap *bmap[K, V]) Delete(key K) error {
	bmap.mutex.Lock()
	defer bmap.mutex.Unlock()
	_, ok := bmap.values[key]
	if !ok {
		return errors.New("Key not found in bmap")
	}
	delete(bmap.values, key)
	keyIndex := bmap.keyIndices[key]
	bmap.keys = append(append(make([]K, len(bmap.keys)-1), bmap.keys[:keyIndex]...), bmap.keys[keyIndex+1:]...)
	delete(bmap.keyIndices, key)
	return nil
}

func (bmap *bmap[K, V]) Swap(key1, key2 K) error {
	bmap.mutex.RLock()
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
	bmap.keyIndices[key1], bmap.keyIndices[key2] = index2, index1
	bmap.keys[index1], bmap.keys[index2] = bmap.keys[index2], bmap.keys[index1]
	bmap.mutex.Unlock()
	return nil
}

func (bmap *bmap[K, V]) Sort(s func(V, V) bool) {
	bmap.mutex.Lock()
	sort.Slice(bmap.keys, func(i, j int) bool {
		return s(bmap.values[bmap.keys[i]], bmap.values[bmap.keys[j]])
	})
	for i, k := range bmap.keys {
		bmap.keyIndices[k] = i
	}
	bmap.mutex.Unlock()
}

func (bmap *bmap[K, V]) SortStable(s func(V, V) bool) {
	bmap.mutex.Lock()
	sort.SliceStable(bmap.keys, func(i, j int) bool {
		return s(bmap.values[bmap.keys[i]], bmap.values[bmap.keys[j]])
	})
	for i, k := range bmap.keys {
		bmap.keyIndices[k] = i
	}
	bmap.mutex.Unlock()
}

func (bmap *bmap[K, V]) SortKeys(s func(K, K) bool) {
	bmap.mutex.Lock()
	sort.Slice(bmap.keys, func(i, j int) bool {
		return s(bmap.keys[i], bmap.keys[j])
	})
	for i, k := range bmap.keys {
		bmap.keyIndices[k] = i
	}
	bmap.mutex.Unlock()
}

func (bmap *bmap[K, V]) SortKeysStable(s func(K, K) bool) {
	bmap.mutex.Lock()
	sort.SliceStable(bmap.keys, func(i, j int) bool {
		return s(bmap.keys[i], bmap.keys[j])
	})
	for i, k := range bmap.keys {
		bmap.keyIndices[k] = i
	}
	bmap.mutex.Unlock()
}

func (bmap *bmap[K, V]) Len() int {
	return len(bmap.keys)
}

func (bmap *bmap[K, V]) Range() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		bmap.mutex.RLock()
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

func (bmap *bmap[K, V]) ValueRange() iter.Seq[V] {
	return func(yield func(V) bool) {
		bmap.mutex.RLock()
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
