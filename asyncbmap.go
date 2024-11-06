package bmap

import (
	"errors"
	"sort"
	"sync"
)

type asyncBmap[K comparable, V any] struct {
	bmap[K, V]
}

func (t Type[K, V]) NewAsync() asyncBmap[K, V] {
	b := asyncBmap[K, V]{}
	b.mutex = sync.RWMutex{}
	b.values = map[K]V{}
	b.keyIndices = map[K]int{}
	return b
}

func (bmap *asyncBmap[K, V]) Set(key K, value V) {
	bmap.mutex.Lock()
	go func() {
		defer bmap.mutex.Unlock()
		_, ok := bmap.values[key]
		bmap.values[key] = value
		if !ok {
			bmap.keyIndices[key] = len(bmap.keys)
			bmap.keys = append(bmap.keys, key)
		}
	}()
}

func (bmap *asyncBmap[K, V]) Delete(key K) {
	bmap.mutex.Lock()
	go func() {
		defer bmap.mutex.Unlock()
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

func (bmap *asyncBmap[K, V]) Swap(key1, key2 K) error {
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
	go func() {
		bmap.values[key1], bmap.values[key2] = bmap.values[key2], bmap.values[key1]
		bmap.keyIndices[key1], bmap.keyIndices[key2] = index2, index1
		bmap.keys[index1], bmap.keys[index2] = bmap.keys[index2], bmap.keys[index1]
		bmap.mutex.Unlock()
	}()
	return nil
}

func (bmap *asyncBmap[K, V]) Sort(s func(V, V) bool) {
	bmap.mutex.Lock()
	sort.Slice(bmap.keys, func(i, j int) bool {
		return s(bmap.values[bmap.keys[i]], bmap.values[bmap.keys[j]])
	})
	for i, k := range bmap.keys {
		bmap.keyIndices[k] = i
	}
	bmap.mutex.Unlock()
}

func (bmap *asyncBmap[K, V]) SortStable(s func(V, V) bool) {
	bmap.mutex.Lock()
	go func() {
		sort.SliceStable(bmap.keys, func(i, j int) bool {
			return s(bmap.values[bmap.keys[i]], bmap.values[bmap.keys[j]])
		})
		for i, k := range bmap.keys {
			bmap.keyIndices[k] = i
		}
		bmap.mutex.Unlock()
	}()
}

func (bmap *asyncBmap[K, V]) SortKeys(s func(K, K) bool) {
	bmap.mutex.Lock()
	go func() {
		sort.Slice(bmap.keys, func(i, j int) bool {
			return s(bmap.keys[i], bmap.keys[j])
		})
		for i, k := range bmap.keys {
			bmap.keyIndices[k] = i
		}
		bmap.mutex.Unlock()
	}()
}

func (bmap *asyncBmap[K, V]) SortKeysStable(s func(K, K) bool) {
	bmap.mutex.Lock()
	go func() {
		sort.SliceStable(bmap.keys, func(i, j int) bool {
			return s(bmap.keys[i], bmap.keys[j])
		})
		for i, k := range bmap.keys {
			bmap.keyIndices[k] = i
		}
		bmap.mutex.Unlock()
	}()
}
