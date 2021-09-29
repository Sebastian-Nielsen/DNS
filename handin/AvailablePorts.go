package main

import (
	"sync"
	"strconv"
)

type SafeIterator struct {
	mu sync.Mutex
	index int
	firstPort int
}

func (i *SafeIterator) Next() string {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.index += 1
	return strconv.Itoa(i.firstPort + i.index)
}


var AvailablePorts = SafeIterator{
	firstPort: 50000,
}

