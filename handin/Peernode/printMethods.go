package Peernode

import (
	. "DNO/handin/Helper"
	"fmt"
	"net"
)

/*
	Printer helper methods
*/
func (p *PeerNode) debugPrintf(text string, args ...interface{}) {
	if p.TestMock.ShouldPrintDebug {
		fmt.Printf("<" + PortOf(p.Listener.Addr()) + "> " + text, args...)
	}
}
func (p *PeerNode) debugPrint(args ...interface{}) {
	if p.TestMock.ShouldPrintDebug {
		fmt.Print( "\t", args, "\n")
	}
}
func (p *PeerNode) debugPrintln(args ...interface{}) {
	if p.TestMock.ShouldPrintDebug {
		fmt.Println("<" + PortOf(p.Listener.Addr()) + ">", args)
	}
}

func (p *PeerNode) println(args ...interface{}) {
	fmt.Println("<" + PortOf(p.Listener.Addr()) + ">", args)
}
func (p *PeerNode) printPort(listener net.Listener) {
	_, port, _ := net.SplitHostPort(listener.Addr().String())
	p.println("Running on port: " + port)
}


