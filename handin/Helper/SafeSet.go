package Helper

import (
	"net"
	"sync"
)

/*
	SafeSet
	A set that avoids race-conditions
*/
type SafeSet_Conn struct {
	mu   sync.Mutex
	Vals map[net.Conn]bool
}

func (s *SafeSet_Conn) Values() map[net.Conn]bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	copyMap := make(map[net.Conn]bool)
	for key, value := range s.Vals {
		copyMap[key] = value
	}
	return copyMap
}
func (s *SafeSet_Conn) Add(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Vals[conn] = true
}
func (s *SafeSet_Conn) Delete(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Vals, conn)
}
func (s *SafeSet_Conn) Contains(conn net.Conn) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Vals[conn]
}
func (s *SafeSet_Conn) ContainsAConnWith(remotePort string) bool {
	for conn := range s.Values() {
		if PortOf(conn.RemoteAddr()) == remotePort {
			return true
		}
	}
	return false
}
func (s *SafeSet_Conn) ToString() string {
	returnVal := "["
	for conn := range s.Values() {
		returnVal += PortOf(conn.RemoteAddr()) + ","
	}
	return returnVal + "]"
}

/*
	SafeSet
	A set that avoids race-conditions
*/
type SafeSet_string struct {
	mu     sync.Mutex
	Values map[string]bool
}

func (s *SafeSet_string) Add(str string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Values[str] = true
}
func (s *SafeSet_string) delete(str string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Values, str)
}
func (s *SafeSet_string) Contains(str string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Values[str]
}

///*
//	SafeSet
//	A set that avoids race-conditions
//*/
//type SafeSet_Transaction struct {
//	mu     sync.Mutex
//	Vals map[ApplyTransaction]bool
//}
//
//func (s *SafeSet_Transaction) Add(transaction ApplyTransaction) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	s.Vals[transaction] = true
//}
//func (s *SafeSet_Transaction) delete(transaction ApplyTransaction) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	delete(s.Vals, transaction)
//}
//func (s *SafeSet_Transaction) Contains(transaction ApplyTransaction) bool {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	return s.Vals[transaction]
//}


///*
//	SafeSet
//	A set that avoids race-conditions
//*/
//type SafeSet_int struct {
//	mu     sync.Mutex
//	Vals map[int]bool
//}
//
//func (s *SafeSet_int) Add(val int) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	s.Vals[val] = true
//}
//func (s *SafeSet_int) delete(val int) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	delete(s.Vals, val)
//}
//func (s *SafeSet_int) Contains(val int) bool {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	return s.Vals[val]
//}
