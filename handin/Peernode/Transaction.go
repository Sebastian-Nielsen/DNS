package Peernode

import (
	. "DNO/handin/Account"
	. "DNO/handin/Helper"
)

/*
	ApplyTransaction related methods
*/
func (p *PeerNode) CreateAccountInLedger(id string, initialAmount int) {
	//p.LocalLedger.CreateAccount(id, initialAmount)
	p.LocalLedger.Accounts[id] = initialAmount
}

func (p *PeerNode) MakeAndBroadcastSignedTransaction(amount int, id string, from string, to string) {
	transaction := Transaction{ID: id, From: from, To: to, Amount: amount}

	p.debugPrintln("Applying transaction:", transaction, ". Accounts before:", p.LocalLedger.Accounts)
	p.LocalLedger.ApplyTransaction(&transaction)
	p.TransactionsSeen.Append(transaction)
	p.debugPrintln("Applying transaction:", transaction, ". Accounts after:", p.LocalLedger.Accounts)

	p.Broadcast(
		Packet{
			Type: PacketType.BROADCAST_TRANSACTION,
			Transaction: transaction,
		},
	)
}

func (p *PeerNode) MakeAndBroadcastTransaction(amount int, id string, from string, to string) {
	transaction := Transaction{ID: id, From: from, To: to, Amount: amount}

	p.debugPrintln("Applying transaction:", transaction, ". Accounts before:", p.LocalLedger.Accounts)
	p.LocalLedger.ApplyTransaction(&transaction)
	p.TransactionsSeen.Append(transaction)
	p.debugPrintln("Applying transaction:", transaction, ". Accounts after:", p.LocalLedger.Accounts)

	p.Broadcast(
		Packet{
			Type: PacketType.BROADCAST_TRANSACTION,
			Transaction: transaction,
		},
	)
}

