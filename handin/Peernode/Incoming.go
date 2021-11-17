package Peernode

import (
	. "DNO/handin/Account"
	. "DNO/handin/Helper"
	"golang.org/x/crypto/pkcs12"
	"math/big"
	"net"
	"strconv"
	"time"
	. "DNO/handin/Cryptography"
)

const SLOT_LENGTH = 1000


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
		p.handleBlock(packet.SignedBlock)
		packet.Type = PacketType.BROADCAST_KNOWN_BLOCK
		p.Broadcast(packet)
	case PacketType.BROADCAST_KNOWN_BLOCK:
		p.debugPrintln("received signed packet: [PacketType: BROADCAST_KNOWN_BLOCK]", packet.SignedBlock.Block.ToString())
		p.handleBlock(packet.SignedBlock)
	case PacketType.BROADCAST_GENESIS_BLOCK:
		p.debugPrintln("received signed packet: [PacketType: BROADCAST_GENESIS_BLOCK] with seed", packet.GenesisBlock.Seed)
		p.handleGensisBlock(packet.GenesisBlock)
	}
}
func (p *PeerNode) handleGensisBlock(G GenesisBlock) {
	p.Sequencer.Hardness = G.Hardness
	p.Sequencer.Seed = G.Seed
	p.Sequencer.Tree.Root = G
	p.Sequencer.Tree.LengthOfBestPath = 0
	p.Sequencer.Tree.LeafHashOfBestPath = G.Seed // just hardcode the "hash" of the genesis block to be the seed
	for _, pk := range G.InitialAccounts {
		p.LocalLedger.CreateAccount(pk.ToString(), G.InitialAmount)
	}
	go p.startLotteryProtocol()
}
func (p *PeerNode) startLotteryProtocol() {
	for {
		//p.Sequencer.SlotNumber.Lock()
		p.Sequencer.SlotNumber.Value += 1
		//p.Sequencer.SlotNumber.Unlock()
		draw := p.getDraw(p.Sequencer.Seed, p.Sequencer.SlotNumber.Value, p.Keys.Sk)
		lotteryString := p.lotteryString(p.Sequencer.Seed, p.Sequencer.SlotNumber.Value)
		if p.isWinner(draw, p.Keys.Pk, lotteryString) {

			block := Block{
				VerificationKey: p.Keys.Pk,
				SlotNumber:      p.Sequencer.SlotNumber.Value,
				Draw:            draw,
				TransactionIDs:  p.Sequencer.UnsequencedTransactionIDs.PopAll(),
				PrevBlockHash:   p.Sequencer.Tree.LeafHashOfBestPath,
				//LengthToRoot: 	 0 p.Sequencer.Tree.LeafHashOfBestPath.LengthToRoot + 1,
				 //hashOfLeafInBestPath:  U_i
				//LengthOfBestPath: M_i
			}
			p.BroadcastBlock(block)
		}
		time.Sleep(SLOT_LENGTH * time.Millisecond)
	}
}
func (p *PeerNode) isWinner(draw *big.Int, publicKey PublicKey, lotteryString string) bool {
	// publicKey := p.Keys.Pk.ToString()
	numberOfTickets := p.LocalLedger.Accounts[publicKey.ToString()]
	toBeHashed := lotteryString + publicKey.ToString() + string(Hash(draw))
	hash, _ := new(big.Int).SetString(toBeHashed, 10)
	hashOfDrawTimesA := new(big.Int).Mul(hash, big.NewInt(int64(numberOfTickets)))
	return hashOfDrawTimesA.Cmp(p.Sequencer.Hardness) == 0
}
func (p *PeerNode) lotteryString(seed string, slotNumber int) string {
	return "LOTTERY" + seed + strconv.Itoa(slotNumber)
}

func (p *PeerNode) getDraw(seed string, slot int, sk SecretKey) *big.Int {
	 value, _ := new(big.Int).SetString(p.lotteryString(seed, slot), 10)
	 return BigInt_createSignature(value, sk)
}
func (p *PeerNode) verifyDraw(signature string, publicKey PublicKey) bool {
	value := p.lotteryString(p.Sequencer.Seed, p.Sequencer.SlotNumber.Value)
	return Verify(signature, value, publicKey)
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


	//p.handleBlock(packet.SignedBlock)
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
func (p* PeerNode) handleBlock(signedBlock SignedBlock) {
	block := signedBlock.Block
	isSequencerPkReceived := p.Sequencer.PublicKey.N != nil
	if !isSequencerPkReceived {
		p.debugPrintln("The sequencer public key has not been received yet, ignoring block")
		return
	}
	isValidSignature := p.Sequencer.Verify(signedBlock)
	if !isValidSignature {
		p.debugPrintln("Signature on block invalid")
		return
	}

	isValidDraw := p.verifyDraw(signedBlock.Signature, block.VerificationKey)
	if !isValidDraw {
		p.debugPrintln("Draw in block invalid")
		return
	}

	ls := p.lotteryString(p.Sequencer.Seed, p.Sequencer.SlotNumber.Value)
	isWinner := p.isWinner(block.Draw, block.VerificationKey, ls)
	if !isWinner {
		p.debugPrintln("Draw in block is not a winner")
		return
	}


/*
	//fmt.Println("Thread trying to lock with blockNumber", p.Sequencer.SlotNumber.Value)
	p.Sequencer.SlotNumber.Lock()
	isNextBlock := p.Sequencer.SlotNumber.Value + 1 == signedBlock.Block.SlotNumber
	if !isNextBlock {
		p.debugPrintln("Block with number", signedBlock.Block.SlotNumber, "is not the next. Current is number", p.Sequencer.SlotNumber.Value)

		//fmt.Println("Thread unlock with blockNumber", p.Sequencer.SlotNumber.Value)
		p.Sequencer.SlotNumber.Unlock()
		return
	}
	p.Sequencer.SlotNumber.Value += 1
	//fmt.Println("Thread unlock with blockNumber", p.Sequencer.SlotNumber.Value)
	p.Sequencer.SlotNumber.Unlock()
*/




//newBlock
//
//	get newblock in
//
//	check if we cannot insert newBlock into map:
//		map.insert(newBlock)
//
//	check if in map
//		(block.prevHash) -> block
//	(block.prevHash) -> block


	_, parentExistInTree := p.Sequencer.Tree.BlockHashToBlock[block.PrevBlockHash]
	if parentExistInTree {
		p.Sequencer.Tree.Insert(block)

		p.insertRecursivelyBlocksThatAreParentTo(block)
	} else {
		// We cant add the block to the tree since its parent isn't in the tree
		p.Sequencer.Tree.BlocksWhoseParentIsNotInTree[block.Hash()] = block
	}

	// Block is valid, so insert it into tree

	p.debugPrintln("Extending list of unapplied Ids of transactions for block number:", block.SlotNumber)
	for _, id := range block.TransactionIDs {
		p.UnappliedIDs.Append(id)
	}
	go p.ApplyUnappliedIDs()
}
func (p *PeerNode) insertRecursivelyBlocksThatAreParentTo(block Block) {
	childToBlock, doExistChildToBlock := p.Sequencer.Tree.BlocksWhoseParentIsNotInTree[block.Hash()]
	if doExistChildToBlock {
		p.Sequencer.Tree.Insert(childToBlock)
		p.insertRecursivelyBlocksThatAreParentTo(childToBlock)
	}
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
