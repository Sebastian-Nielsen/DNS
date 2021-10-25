package Peernode

import "net"

func (p *PeerNode) ListenForNewConns() {
	// Continously listen for new connection requests
	p.debugPrintln("Listening for new connections ...")
	for {
		newConn, _ := p.Listener.Accept()
		p.OpenConnections.Add(newConn)
		p.debugPrintln("Got a new connection ...", newConn.LocalAddr(), "->", newConn.RemoteAddr())
		go p.Listen(newConn)
	}
}
func (p *PeerNode) Listen(conn net.Conn) {
	defer conn.Close()
	for {
		packet, ok := p.Ipc.Receive(conn)
		if !ok { p.debugPrintln("Receive error") }
		go p.handleIncomming(packet, conn)
	}
}
func (p *PeerNode) startServer(port string) {
	p.Listener, _ = net.Listen("tcp", ":" + port)
	p.printPort(p.Listener)

	go p.ListenForNewConns()
}
