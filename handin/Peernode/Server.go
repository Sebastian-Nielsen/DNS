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
		//p.debugPrintln("About to go into receive")
		packet, ok := p.Ipc.Receive(conn)
		//p.debugPrintln("[Listen] Received packet: ", packet, " - ok:", ok)
		if !ok { 
			p.debugPrintln("Listen: Receive error")
			return
		}
		go p.handleIncomming(packet, conn)
	}
}
func (p *PeerNode) startServer(port string) {
	p.Listener, _ = net.Listen("tcp", "localhost:" + port)
	p.printPort(p.Listener)

	// go p.ListenForNewConns()
}
