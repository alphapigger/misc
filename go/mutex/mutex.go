package mutex

import (
	"sync"
)

var (
	mu   sync.Mutex
	rwmu sync.RWMutex
	a    uint64
	m    = map[string]string{
		"a": "1",
	}
)

func Mutex() (v string) {
	mu.Lock()
	v = m["a"]
	mu.Unlock()
	return v
}

func RWMutex() (v string) {
	rwmu.RLock()
	v = m["a"]
	rwmu.RUnlock()
	return v
}
