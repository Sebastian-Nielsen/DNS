package Peernode

import (
	. "DNO/handin/Account"
	. "DNO/handin/Cryptography"
	. "DNO/handin/Helper"
	"bufio"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const SEED = "SebastianAndAndreasBlockchain"

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

	if p.IsInitialSequencer() {
		p.waitUntilTenOpenConnections()
		go p.SendInitialGenisisBlock()
	}
	//go p.PullFromNeighbors() // Occassionally send pull-requests to neighbors, asking for their messagesSent set

	p.send() // Continously prompt the user msgs for the peerNode to broadcast
}


func (p *PeerNode) IsInitialSequencer() bool {
	//return p.Sequencer.KeyPair != KeyPair{}
	return p.Sequencer.IsSequencer
}
func (p *PeerNode) waitUntilTenOpenConnections() {
	for len(p.OpenConnections.Values()) < 9 {
		fmt.Println("waiting. Currently: " + strconv.Itoa(len(p.OpenConnections.Values())))
		time.Sleep(500 * time.Millisecond)
	}
}
func (p *PeerNode) SendInitialGenisisBlock() {
	//time.Sleep(5 * time.Second)
	hardness, _ := new(big.Int).SetString("1129860864649195017936278738467622451102563362381435984249705038832445929304130000000", 10)
	//hardness, _ := new(big.Int).SetString("999", 10)
	fmt.Println("Hardness is: " + hardness.String())
	genesisBlock := GenesisBlock{
		Seed: SEED,
		InitialAccounts: GetHardcodedAccPublicKeys(),
		Hardness: hardness,
		InitialAmount: 10e6,
	}
	p.debugPrintln("Broadcasting Genesis block to ", p.OpenConnections.ToString())
	p.Broadcast(
		//Packet { Type: PacketType.BROADCAST_MSG, GenesisBlock: genesisBlock },
		Packet { Type: PacketType.BROADCAST_GENESIS_BLOCK, GenesisBlock: genesisBlock },
	)

	p.handleGensisBlock(genesisBlock)
}

//func (p *PeerNode) BroadcastBlock() {
//	block := Block {
//		SlotNumber:     p.Sequencer.SlotNumber.Value + 1,
//		TransactionIDs: p.Sequencer.UnsequencedTransactionIDs.PopAll(),
//	}
//
//	p.debugPrintln("Broadcasting block with", strconv.Itoa(len(block.TransactionIDs)), "elements")
//	p.Broadcast(
//		Packet { Type: PacketType.BROADCAST_BLOCK, SignedBlock: p.Sequencer.Sign(block) },
//	)
//	p.handleBlock(SignedBlock{})
//}

func (p *PeerNode) BroadcastBlock(block Block) {
	fmt.Println("Broadcasting block with", strconv.Itoa(len(block.TransactionIDs)), "elements")
	signedBlock := p.Sequencer.Sign(block, p.Keys.Sk)
	p.Broadcast(
		Packet { Type: PacketType.BROADCAST_BLOCK, SignedBlock: signedBlock },
	)
	p.handleBlock(signedBlock)
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
		//p.Sequencer.KeyPair = GenKeyPair()
		//p.Sequencer.PublicKey = p.Sequencer.KeyPair.Pk
		p.Sequencer.IsSequencer = true
		//p.debugPrintln("Sequencer public key is:", p.Sequencer.PublicKey.ToString()[:10], "...")
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


