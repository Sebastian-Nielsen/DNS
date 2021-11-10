package Peernode

import (
	. "DNO/handin/Account"
	. "DNO/handin/Cryptography"
	. "DNO/handin/Helper"
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)


/*
	PeerNode
	- Contionusly listen for new connection requests
		- add new connections to OpenConnections

	- process incomming messages
	- process outgoing messages

	- Occassionally send a pull request to a neighbor (taking turns)
*/
type PeerNode struct {
	OpenConnections     	SafeSet_Conn
	PeersInArrivalOrder 	SafeArray_string
	MessagesSent        	SafeSet_string
	Listener            	net.Listener
	Ipc                 	IPC
	TestMock            	Mock
	LocalLedger         	*Ledger
	TransactionsSeen    	SafeArray_Transaction
	SignedTransactionsSeen  SafeMap_string_to_SignedTransaction
	Keys				    KeyPair
	Sequencer               Sequencer
	UnappliedIDs            SafeArray_string
	unappliedIDsMutex       sync.Mutex
	unappliedIDSMutexIsLocked bool
}


/*
	atPort: port at which the PeerNode should listen at.
	remotePort: port of remote PeerNode to connect to initially.
				if remotePort is "" (empty string) then prompt the
				user for a port.
*/
func (p *PeerNode) Start(atPort, remotePort string) {
	socket := p.createSocket(remotePort)
	
	p.startServer(atPort)
	
	p.DialNetwork(socket)

	if p.IsSequencer() {
		go p.PeriodicallySendUnsequencedTransactions()
	}
	//go p.PullFromNeighbors() // Occassionally send pull-requests to neighbors, asking for their messagesSent set

	p.send() // Continously prompt the user msgs for the peerNode to broadcast
}


func (p *PeerNode) IsSequencer() bool {
	return p.Sequencer.KeyPair != KeyPair{}
}
func (p *PeerNode) PeriodicallySendUnsequencedTransactions() {
	for {
		time.Sleep(30 * time.Second)
		p.BroadcastBlock()
	}
}
func (p *PeerNode) BroadcastBlock() {
	block := Block {
		BlockNumber: p.Sequencer.BlockNumber.Value + 1,
		TransactionIDs: p.Sequencer.UnsequencedTransactionIDs.PopAll(),
	}

	p.debugPrintln("Broadcasting block with", strconv.Itoa(len(block.TransactionIDs)), "elements")
	p.Broadcast(
		Packet { Type: PacketType.BROADCAST_BLOCK, SignedBlock: p.Sequencer.Sign(block) },
	)
	p.ExtendUnappliedIDsIfValidBlock(SignedBlock{})
}
func PromptForRemoteSocket(p *PeerNode) Socket {

	ip := "127.0.0.1"  // The handin asks us to also prompt the user for an ip, but really no need for it ...

	fmt.Print("Connect to port: ")
	port := strings.TrimSpace(input(p))
	fmt.Printf("Connecting to: '%s'\n", port)

	return Socket{Ip: ip, Port: port}
}
func input(p *PeerNode) string {
	if !p.TestMock.ShouldMockInput {
		// var returnString string
		// fmt.Scanln(&returnString)
		// return returnString
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		return scanner.Text()
	}

	for p.TestMock.SimulatedInputString == "" {
		time.Sleep(250 * time.Millisecond)
	} // waiting for input
	returnString := p.TestMock.SimulatedInputString
	p.TestMock.SimulatedInputString = ""
	return returnString

}
func Assert(condition bool) {
	if !condition {
		panic("assert failed")
	}
}
func (p *PeerNode) createSocket(remotePort string) Socket {
	if remotePort != "" {
		return Socket{Ip: "127.0.0.1", Port: remotePort}
	}
	return PromptForRemoteSocket(p)
}
func (p *PeerNode) DialNetwork(remoteSocket Socket) {
	conn, err := net.Dial("tcp", remoteSocket.ToString())

	dialIsSuccess := err == nil
	if !dialIsSuccess {
		p.debugPrintln("--------------")
		p.println("DialNetwork failed")
		p.debugPrintln("--------------")
		p.debugPrintln("Adding local port to list")
		p.PeersInArrivalOrder.Append(PortOf(p.Listener.Addr()))
		p.debugPrintln("Setting ourself as the sequencer")
		p.Sequencer.KeyPair = GenKeyPair()
		p.Sequencer.PublicKey = p.Sequencer.KeyPair.Pk
		p.debugPrintln("Sequencer public key is:", p.Sequencer.PublicKey.ToString()[:10], "...")
	} else {
		p.debugPrintln("--------------")
		p.debugPrintln("Local  Addr:", conn.LocalAddr())
		p.debugPrintln("Remote Addr:", conn.RemoteAddr())
		p.debugPrintln("--------------")
		p.OpenConnections.Add(conn)

		p.debugPrintln("[dial] Sending Listener port to dialed connection (" + PortOf(conn.RemoteAddr()) + ")")

		//establishConnectionToNetwork()
		//setup

		p.Ipc.Send(
			Packet{Type: PacketType.LISTENER_PORT, Msg: PortOf(p.Listener.Addr()) },
			conn,
		)

		p.debugPrintln("Sending initial pull-request to neighbor:", conn.RemoteAddr().String())

		p.Ipc.Send(Packet{Type: PacketType.PULL}, conn)
		var packet Packet
		for {
			packet, _ = p.Ipc.Receive(conn)
			if packet.Type == PacketType.PULL_REPLY {break}
			p.debugPrintf("[Dial] Received wrong packet type:", packet.Type)
		}
		p.debugPrintln("[Dial] Received Pull-Reply packet:", packet.PeersInArrivalOrderValues, packet.MessagesSent)

		p.HandlePullReplyPacket(packet)

		peers := packet.PeersInArrivalOrderValues

		p.connectToPeers(peers)

		p.Broadcast(
			Packet{ Type: PacketType.BROADCAST_LISTENER_PORT , Msg: PortOf(p.Listener.Addr()) },
		)

		go p.Listen(conn)
	}
	go p.ListenForNewConns()
}
func (p *PeerNode) connectToPeers(portsOfPeers []string) {
	Assert(p.Listener != nil)
	//upTo := int(math.Min(11, float64(len(portsOfPeers)))-1)

	p.debugPrintln("Connecting to up to 10 peers in received list:\n\t", portsOfPeers)

	if len(portsOfPeers) == 0 {
		p.println("Error: Received empty list from connection")
		return
	}
	newConnections := 0
	for _, port := range portsOfPeers {       // portsOfPeers = ["50001", "50002", ...]
		alreadyHaveAConnectionToThisPort := p.OpenConnections.ContainsAConnWith(port)
		if alreadyHaveAConnectionToThisPort {
			p.debugPrintln("Don't connect to this, we already have a connection to port:", port)
			continue
		}
		isPortOurRemotePort := port == PortOf(p.Listener.Addr())
		if isPortOurRemotePort {
			p.debugPrintln("Don't connect to this, it is ourselves *insert spiderman pointing meme*:", port)
			continue
		}

		conn := p.dial(Socket{Ip: "127.0.0.1", Port: port})
		p.OpenConnections.Add(conn)
		newConnections += 1
		p.debugPrint("\tConnecting to peer:", port)
		if newConnections == 10 {
			break
		}
	}
}
func (p *PeerNode) dial(socket Socket) net.Conn {
	conn, err := net.Dial("tcp", socket.ToString())
	if err != nil {
		p.println(err)
	}
	return conn
}
func (p *PeerNode) PullFromNeighbors() {
	// Requests the full messagesSent set from each neighbor in turn, wait inbetween each request
	time.Sleep(300 * time.Second)
	packet := Packet{Type: PacketType.PULL}
	for {
		for openConn := range p.OpenConnections.Values() {
			p.println("Sending pull-request to neighbor:", openConn.RemoteAddr().String())
			p.Ipc.Send(packet, openConn)
			time.Sleep(30 * time.Second)
		}
	}
}
func GenKeyPair() KeyPair {
	pk, sk := GenKeys(2000)
	return KeyPair{ Pk: pk, Sk: sk }
}


