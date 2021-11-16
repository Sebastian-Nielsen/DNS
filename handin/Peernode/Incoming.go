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
			Type: 						PacketType.PULL_REPLY,
			MessagesSent: 				p.MessagesSent.Values,
			PeersInArrivalOrderValues:  p.PeersInArrivalOrder.Values(),
			TransactionsSeen: 			p.TransactionsSeen.Values(),
			SequencerPublicKey: 		p.Sequencer.KeyPair.Pk,
		}
		p.debugPrintln("Received packet: [Type: PULL] ... Sending back: \n\tpacket{peersInArrivalOrder:", packet.PeersInArrivalOrderValues, "}")
		p.debugPrintln("Sending public key:", p.Sequencer.KeyPair.Pk.ToString()[:10], "...")
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

		_,  signedTransactionIsSeen := p.SignedTransactionsSeen.Get(signedTransaction.ID)
		if !signedTransactionIsSeen {
			p.SignedTransactionsSeen.Put(signedTransaction.ID, signedTransaction)
			p.Sequencer.UnsequencedTransactionIDs.Append(signedTransaction.ID)
			go p.ApplyUnappliedIDs()
			packet.Type = PacketType.BROADCAST_KNOWN_SIGNED_TRANSACTION
			p.Broadcast(packet)
		}
	case PacketType.BROADCAST_KNOWN_SIGNED_TRANSACTION:
		p.debugPrintln("received signed packet: [PacketType: BROADCASTED_KNOWN_SIGNED_TRANSACTION]", packet.SignedTransaction.ToString())
		signedTransaction := packet.SignedTransaction

		_,  signedTransactionIsSeen := p.SignedTransactionsSeen.Get(signedTransaction.ID)
		if !signedTransactionIsSeen {
			p.SignedTransactionsSeen.Put(signedTransaction.ID, signedTransaction)
			p.Sequencer.UnsequencedTransactionIDs.Append(signedTransaction.ID)
			// p.Broadcast(packet)
			go p.ApplyUnappliedIDs()
		}
	case PacketType.BROADCAST_BLOCK:
		//p.debugPrintln("received signed packet: [PacketType: BROADCAST_BLOCK]", packet.SignedBlock.Block.ToString())
		p.ExtendUnappliedIDsIfValidBlock(packet.SignedBlock)
		packet.Type = PacketType.BROADCAST_KNOWN_BLOCK
		p.Broadcast(packet)
	case PacketType.BROADCAST_KNOWN_BLOCK:
		p.debugPrintln("received signed packet: [PacketType: BROADCAST_KNOWN_BLOCK]", packet.SignedBlock.Block.ToString())
		p.ExtendUnappliedIDsIfValidBlock(packet.SignedBlock)
	case PacketType.BROADCAST_GENESIS_BLOCK:
		p.debugPrintln("received signed packet: [PacketType: BROADCAST_GENESIS_BLOCK] with seed", packet.GenesisBlock.Seed)
		p.handleGensisBlock(packet.GenesisBlock)
	}
}
func (p *PeerNode) handleGensisBlock(G GenesisBlock) {
	p.Sequencer.Hardness = G.Hardness
	p.Sequencer.Seed = G.Seed
	for _, pk := range G.InitialAccounts {
		p.LocalLedger.CreateAccount(pk.ToString(), G.InitialAmount)
	}
}
func (p *PeerNode) HandlePullReplyPacket(packet Packet) {
	// Put the sequencer public to the received pk
	p.debugPrintln("Setting sequencer public key to:", packet.SequencerPublicKey.ToString()[:10])
	p.Sequencer.PublicKey = packet.SequencerPublicKey

	// Add all messages contained in the PULL-REPLY packet to our set of messages
	p.extendMessagesSentSet(packet.MessagesSent)

	// Add all peer ports contained in the PULL-REPLY packet to our list of ports
	p.extendPeersInArrivalOrder(packet.PeersInArrivalOrderValues)

	// Apply all transactions contained in the PULL-REPLY packet on our local ledger
	//p.applyAllTransactions(packet.TransactionsSeen)

	// Apply all signed transactions contained in the PULL-REPLY packet on our local ledger
	//p.applyAllSignedTransactions(packet.SignedTransactionsSeen)


	//p.ExtendUnappliedIDsIfValidBlock(packet.SignedBlock)
}
func (p *PeerNode) ApplyUnappliedIDs() {
	if p.unappliedIDSMutexIsLocked { return }
	//fmt.Println("A thread entered 'ApplyUnappliedIDS")
	p.unappliedIDsMutex.Lock()
	p.unappliedIDSMutexIsLocked = true
	defer p.unappliedIDsMutex.Unlock()

	for !(p.UnappliedIDs.IsEmpty()) {

		id := p.UnappliedIDs.Get(0)
		transaction, doWeHaveNextTransactionToApply := p.SignedTransactionsSeen.Get(id)
		if !doWeHaveNextTransactionToApply {
			p.unappliedIDSMutexIsLocked = false
			return
		}

		p.UnappliedIDs.PopHead()
		p.LocalLedger.ApplySignedTransaction(transaction)

	}
	p.unappliedIDSMutexIsLocked = false
}
func (p* PeerNode) ExtendUnappliedIDsIfValidBlock(signedBlock SignedBlock) {
	SequencerPkIsReceived := p.Sequencer.PublicKey.N != nil
	if !SequencerPkIsReceived {
		p.debugPrintln("The sequencer public key has not been received yet, ignoring block")
		return
	}
	isValidSignature := p.Sequencer.Verify(signedBlock)
	if !isValidSignature {
		p.debugPrintln("Signature on block invalid")
		return
	}

	//fmt.Println("Thread trying to lock with blockNumber", p.Sequencer.BlockNumber.Value)
	p.Sequencer.BlockNumber.Lock()
	isNextBlock := p.Sequencer.BlockNumber.Value + 1 == signedBlock.Block.BlockNumber
	if !isNextBlock {
		p.debugPrintln("Block with number", signedBlock.Block.BlockNumber, "is not the next. Current is number", p.Sequencer.BlockNumber.Value)

		//fmt.Println("Thread unlock with blockNumber", p.Sequencer.BlockNumber.Value)
		p.Sequencer.BlockNumber.Unlock()
		return
	}
	p.Sequencer.BlockNumber.Value += 1
	//fmt.Println("Thread unlock with blockNumber", p.Sequencer.BlockNumber.Value)
	p.Sequencer.BlockNumber.Unlock()

	p.debugPrintln("Extending list of unapplied Ids of transactions for block number:", signedBlock.Block.BlockNumber)
	for _, id := range signedBlock.Block.TransactionIDs {
		p.UnappliedIDs.Append(id)
	}
	go p.ApplyUnappliedIDs()
}
func (p *PeerNode) applyAllSignedTransactions(signedTransactions []SignedTransaction) {
	p.debugPrintln("Applying ", len(signedTransactions), " signed transactions")
	for _, signedTransaction := range signedTransactions {
		p.debugPrintln("Applying signed transaction: " + signedTransaction.ToString())
		p.LocalLedger.ApplySignedTransaction(signedTransaction)
		p.SignedTransactionsSeen.Put(signedTransaction.ID, signedTransaction)
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
