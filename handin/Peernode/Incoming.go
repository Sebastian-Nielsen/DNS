package Peernode

import (
	. "DNO/handin/Account"
	. "DNO/handin/Helper"
	"net"
)


func (p *PeerNode) handleIncomming(packet Packet, connPacketWasReceivedOn net.Conn) {

	switch packet.Type {
	case PacketType.BROADCAST_MSG:
		p.debugPrintln("Received packet: [Type: BROADCAST_MSG][Msg: " + packet.Msg + "] ... Broadcasting it")
		p.BroadcastMessage(packet)
	case PacketType.PULL:
		packet = Packet{
			Type: PacketType.PULL_REPLY,
			MessagesSent: p.MessagesSent.Values,
			PeersInArrivalOrderValues: p.PeersInArrivalOrder.Values(),
			TransactionsSeen: p.TransactionsSeen.Values(),
		}
		p.debugPrintln("Received packet: [Type: PULL] ... Sending back: \n\tpacket{peersInArrivalOrder:", packet.PeersInArrivalOrderValues, "}")
		p.Ipc.Send(packet, connPacketWasReceivedOn)
	case PacketType.PULL_REPLY:
		p.debugPrintln("Received packet: [Type: PULL-REPLY] ... Extending our messagesSent set with the messages in the packet")
		p.HandlePullReplyPacket(packet)
	case PacketType.LISTENER_PORT:
		listenerPort := packet.Msg
		p.PeersInArrivalOrder.Append(listenerPort)
		p.debugPrintln("Received packet: [Type: LISTENER_PORT] ... Adding its port:", listenerPort, "peersInArrivalOrder is now:\n\t", p.PeersInArrivalOrder.Values())
	case PacketType.BROADCAST_LISTENER_PORT:
		listenerPort := packet.Msg
		p.debugPrintln("received packet: [Type: BROADCAST_LISTENER_PORT] ... ", listenerPort)
		if !p.PeersInArrivalOrder.Contains(listenerPort) {
			p.PeersInArrivalOrder.Append(listenerPort)
			p.debugPrintln("[PacketType: BROADCASTED_LISTENER_PORT] Adding port:", listenerPort, "peersInArrivalOrder is now:\n\t", p.PeersInArrivalOrder.Values())
		} else {
			packet.Type = PacketType.BROADCASTED_KNOWN_LISTENER_PORT
		}
		p.Broadcast(packet)
	case PacketType.BROADCASTED_KNOWN_LISTENER_PORT:
		listenerPort := packet.Msg
		p.debugPrintln("Received packet: [Type: BROADCASTED_KNOWN_LISTENER_PORT] ... ", listenerPort)

		if p.PeersInArrivalOrder.Contains(listenerPort) {
			p.debugPrintf("Received listenerPort we already have: %s\n", packet.Msg)
			return // Ignore the packet
		}
		p.PeersInArrivalOrder.Append(listenerPort)
		p.debugPrintln("[PacketType: BROADCASTED_KNOWN_LISTENER_PORT] Adding port", listenerPort, "peersInArrivalOrder is now:\n\t", p.PeersInArrivalOrder.Values())
		p.Broadcast(packet)
	case PacketType.BROADCAST_TRANSACTION:
		p.debugPrintln("received packet: [PacketType: BROADCASTED_TRANSACTION]", packet.Transaction)
		transaction := packet.Transaction

		transactionIsNotSeen := !p.TransactionsSeen.Contains(transaction)
		if transactionIsNotSeen {
			p.TransactionsSeen.Append(transaction)
			p.LocalLedger.ApplyTransaction(&transaction)
		}

		packet.Type = PacketType.BROADCAST_KNOWN_TRANSACTION
		p.Broadcast(packet)
	case PacketType.BROADCAST_KNOWN_TRANSACTION:
		p.debugPrintln("received packet: [PacketType: BROADCASTED__KNOWN_TRANSACTION]", packet.Transaction)
		transaction := packet.Transaction

		transactionIsNotSeen := !p.TransactionsSeen.Contains(transaction)
		if transactionIsNotSeen {
			p.TransactionsSeen.Append(transaction)
			p.LocalLedger.ApplyTransaction(&transaction)
			p.Broadcast(packet)
		}
	case PacketType.BROADCAST_SIGNED_TRANSACTION:
		p.debugPrintln("received signed packet: [PacketType: BROADCAST_SIGNED_TRANSACTION]", packet.SignedTransaction.ToString())
		signedTransaction := packet.SignedTransaction

		signedTransactionIsNotSeen := !p.SignedTransactionsSeen.Contains(signedTransaction)
		if signedTransactionIsNotSeen {
			p.SignedTransactionsSeen.Append(signedTransaction)
			p.LocalLedger.ApplySignedTransaction(signedTransaction)
		}

		packet.Type = PacketType.BROADCAST_KNOWN_SIGNED_TRANSACTION
		p.Broadcast(packet)
	case PacketType.BROADCAST_KNOWN_SIGNED_TRANSACTION:
		p.debugPrintln("received signed packet: [PacketType: BROADCASTED_KNOWN_SIGNED_TRANSACTION]", packet.SignedTransaction.ToString())
		signedTransaction := packet.SignedTransaction

		signedTransactionIsNotSeen := !p.SignedTransactionsSeen.Contains(signedTransaction)
		if signedTransactionIsNotSeen {
			p.SignedTransactionsSeen.Append(signedTransaction)
			p.LocalLedger.ApplySignedTransaction(signedTransaction)
			p.Broadcast(packet)
		}
	}
}
func (p *PeerNode) HandlePullReplyPacket(packet Packet) {
	// Add all messages contained in the PULL-REPLY packet to our set of messages
	p.extendMessagesSentSet(packet.MessagesSent)

	// Add all peer ports contained in the PULL-REPLY packet to our list of ports
	p.extendPeersInArrivalOrder(packet.PeersInArrivalOrderValues)

	// Apply all transactions contained in the PULL-REPLY packet on our local ledger
	p.applyAllTransactions(packet.TransactionsSeen)

	// Apply all signed transactions contained in the PULL-REPLY packet on our local ledger
	p.applyAllSignedTransactions(packet.SignedTransactionsSeen)
}
func (p *PeerNode) applyAllSignedTransactions(signedTransactions []SignedTransaction) {
	p.debugPrintln("Applying ", len(signedTransactions), " signed transactions")
	for _, signedTransaction := range signedTransactions {
		p.debugPrintln("Applying signed transaction: " + signedTransaction.ToString())
		p.LocalLedger.ApplySignedTransaction(signedTransaction)
		p.SignedTransactionsSeen.Append(signedTransaction)
	}
}
func (p *PeerNode) applyAllTransactions(transactions []Transaction) {
	p.debugPrintln("Applying all transactions:", transactions)
	for _, transaction := range transactions {
		p.LocalLedger.ApplyTransaction(&transaction)
		p.TransactionsSeen.Append(transaction)
	}
}
func (p *PeerNode) extendPeersInArrivalOrder(peers []string) {
	// no contains method i golang, so simple O(n*m) solution
	p.debugPrintln("Extending PeersInArrivalOrder")
	for _, receivedPeer := range peers {
		if !p.PeersInArrivalOrder.Contains(receivedPeer) {
			p.debugPrint("\tAdding:", receivedPeer, "to PeersInArrivalOrder")
			p.PeersInArrivalOrder.Append(receivedPeer)
		}
	}
	p.debugPrintln("PeersInArrivalOrder is now:", p.PeersInArrivalOrder.Values())
}
func (p *PeerNode) extendMessagesSentSet(messages map[string]bool) {
	if len(messages) == 0 { return }

	p.debugPrintln("Adding messages:")
	for msg := range messages {
		p.println("\t", msg)
		p.MessagesSent.Add(msg)
	}
}
