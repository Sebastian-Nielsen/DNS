package Helper

import "sync"

type SafeCounter struct {
	Value  int
	mu     sync.Mutex
}
func (c *SafeCounter) Lock() {
	c.mu.Lock()
}
func (c *SafeCounter) Unlock() {
	c.mu.Unlock()
}