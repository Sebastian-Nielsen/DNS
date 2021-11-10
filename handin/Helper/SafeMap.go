package Helper

import (
	"DNO/handin/Account"
	"sync"
)

/*
	SafeMap
	A map that avoids race-conditions
*/
type SafeMap_string_to_SignedTransaction struct {
	mu     sync.Mutex
	Values map[string] Account.SignedTransaction
}

func (s *SafeMap_string_to_SignedTransaction) Put(key string, value Account.SignedTransaction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Values[key] = value
}
func (s *SafeMap_string_to_SignedTransaction) Get(key string) (Account.SignedTransaction, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.Values[key]
	return value, ok
}
//func (s *SafeMap_Transaction) delete(transaction ApplyTransaction) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	delete(s.Vals, transaction)
//}
//func (s *SafeMap_Transaction) Contains(transaction ApplyTransaction) bool {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	return s.Vals[transaction]
//}
