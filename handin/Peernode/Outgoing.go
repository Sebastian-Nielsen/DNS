package Peernode

import (
	. "DNO/handin/Helper"
	"fmt"
	"strings"
)


func (p *PeerNode) send() {
	/* Continously prompt the user for messages to send */
	//p.debugPrintln("[PeerNode:send] Awaiting input to send ... ")
	//p.println("[PeerNode:send] > Type 'm' to view MessagesSent ")
	//p.println("[PeerNode:send] > Type 'c' to view OpenConnections ")
	for {
		msg := strings.TrimSpace(input(p))
		if msg == "m" {
			println(p.MessagesSent.Values)
			continue
		}
		if msg == "c" {
			println(p.OpenConnections.Values)
			continue
		}
		if msg == "a" {
			println(p.PeersInArrivalOrder.Values())
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
		p.debugPrintf("Received msg we already have: %s", packet.Msg)
		return // Ignore the packet
	}
	p.println("Adding and broadcasting msg: '" + packet.Msg + "'")
	p.MessagesSent.Add(packet.Msg)
	p.Broadcast(packet)
}

func (p *PeerNode) Broadcast(packet Packet) {
	p.debugPrintf("[There's now %d openConnections:\n\t%s\n", len(p.OpenConnections.Values), p.OpenConnections.ToString())
	for openConn := range p.OpenConnections.Values {
		p.println("Sending msg to openConn:", openConn.RemoteAddr())
		p.Ipc.Send(packet, openConn)
	}
}
