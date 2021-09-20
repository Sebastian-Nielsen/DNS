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
	mu     sync.Mutex
	Values map[net.Conn]bool
}

func (s *SafeSet_Conn) Add(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Values[conn] = true
}
func (s *SafeSet_Conn) Delete(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Values, conn)
}
func (s *SafeSet_Conn) Contains(conn net.Conn) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Values[conn]
}
func (s *SafeSet_Conn) ContainsAConnWith(remotePort string) bool {
	for conn := range s.Values {
		if PortOf(conn.RemoteAddr()) == remotePort {
			return true
		}
	}
	return false
}
func (s *SafeSet_Conn) ToString() string {
	returnVal := "["
	for conn := range s.Values {
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
