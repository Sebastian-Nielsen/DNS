package Helper

import (
	"fmt"
	"net"
)

type Socket struct {
	Ip string
	Port string
}
func (s* Socket) ToString() string {
	return fmt.Sprintf("%s:%s", s.Ip, s.Port)
}

func PortOf(addr net.Addr) string {
	_, port, _ := net.SplitHostPort(addr.String())
	return port
}
