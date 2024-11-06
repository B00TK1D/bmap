package bmap

import (
	"sync"
)

type Bsem struct {
	mutex sync.RWMutex
	sem   sync.WaitGroup
}

func (bsem *Bsem) Add(delta int) {
	bsem.mutex.Lock()
	bsem.sem.Add(delta)
	bsem.mutex.Unlock()
}

func (bsem *Bsem) Done() {
	bsem.sem.Done()
}

func (bsem *Bsem) Wait() {
	bsem.mutex.RLock()
	bsem.sem.Wait()
	bsem.mutex.RUnlock()
}
