package Helper

import "fmt"

type Socket struct {
	Ip string
	Port string
}
func (s* Socket) ToString() string {
	return fmt.Sprintf("%s:%s", s.Ip, s.Port)
}
