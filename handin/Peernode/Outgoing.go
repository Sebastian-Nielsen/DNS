package Peernode

import (
	"DNO/handin/Cryptography"
	. "DNO/handin/Helper"
	"fmt"
	"strings"
	"time"
)

func (p *PeerNode) send() {
	/* Continously prompt the user for messages to send */
	for {
		if p.TestMock.ShouldPrintDebug {
			time.Sleep(200 * time.Millisecond)
		} // wait for all debug prints to finish
		if !p.TestMock.ShouldMockInput {
			fmt.Print("\nType 'm' to view MessagesSent\n" +
				"Type 'c' to view OpenConnections\n" +
				"Type 'p' to view PeersInArrivalOrder\n" +
				"Type 'i' to view ListenerPort\n" +
				"Type 'w' to start interacting with Software Wallet\n" +
				"Or type a message: ")
		}
		msg := strings.TrimSpace(input(p)) // maybe removes too much?
		if msg == "m" {
			fmt.Println("\nAll messages:")
			for k := range p.MessagesSent.Values {
				fmt.Println("\t" + k)
			}
			// fmt.Println(printMessagesSent(p.MessagesSent.Vals))
			continue
		}
		if msg == "c" {
			fmt.Println()
			fmt.Println(p.OpenConnections.Values())
			continue
		}
		if msg == "p" {
			fmt.Println()
			fmt.Println(p.PeersInArrivalOrder.Values())
			continue
		}
		if msg == "i" {
			fmt.Println()
			p.printPort(p.Listener)
			continue
		}
		if msg == "w" {
			for {
				fmt.Println("Debug: (AB!12abc) is a valid password")
				fmt.Println("<Software Wallet> Enter a password:")
				password := input(p)
				fmt.Println("<Software Wallet> Enter a filename: ")
				filename := input(p)
				pk, err := Cryptography.Generate(filename, password)
				if err == nil {
					fmt.Println("Your private key is:\n", pk)
					break
				} else {
					fmt.Println("<Software Wallet> ", err)
					fmt.Println("<Software Wallet> Try again...")
					fmt.Println()
				}
			}
			continue
		}
		p.HandleOutgoing(Packet{Type: PacketType.BROADCAST_MSG, Msg: msg})
	}
}

func (p *PeerNode) HandleOutgoing(packet Packet) {
	switch packet.Type {
	case PacketType.BROADCAST_MSG:
		if p.MessagesSent.Contains(packet.Msg) {
			p.debugPrintf("Cancelling the sending of msg: '%s' (reason: already in messagesSent)", packet.Msg)
			return // Ignore packet
		}
		p.MessagesSent.Add(packet.Msg)
		p.Broadcast(packet)
	default:
		fmt.Printf("[HandleOutgoing] Got an unknown packet.Type: %s", packet.Type)
	}
}

func (p *PeerNode) BroadcastMessage(packet Packet) {
	if p.MessagesSent.Contains(packet.Msg) {
		p.debugPrintf("Received msg we already have: %s\n", packet.Msg)
		return // Ignore the packet
	}
	p.debugPrint("Adding and broadcasting msg: '" + packet.Msg + "'")
	p.MessagesSent.Add(packet.Msg)
	p.Broadcast(packet)
}

func (p *PeerNode) Broadcast(packet Packet) {

	vals := p.OpenConnections.Values()
	p.debugPrintf("[There's now %d openConnections:\n\t%s\n", len(vals), p.OpenConnections.ToString())
	for openConn := range vals {
		p.debugPrintln("Sending msg to openConn:", openConn.RemoteAddr())
		p.Ipc.Send(packet, openConn)
	}
}
