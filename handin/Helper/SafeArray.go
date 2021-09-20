package Helper

import (
	"sync"
)

type SafeArray_string struct {
	values []string
	mu     sync.Mutex
}
func (a *SafeArray_string) Append(value string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.values = append(a.values, value)
}
func (a *SafeArray_string) Values() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	cpy := make([]string, len(a.values))
	copy(cpy, a.values)
	return cpy
}
func (a *SafeArray_string) Contains(valueToSearchFor string) bool {
	for _, val := range a.values {
		if val == valueToSearchFor {
			return true
		}
	}
	return false
}