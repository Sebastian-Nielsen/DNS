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
	signedTransaction := p.LocalLedger.MakeSignedTransaction(transaction, p.Keys.Sk)

	// p.debugPrintln("Applying signed-transaction:", signedTransaction.ToString())//, ". Accounts before:", p.LocalLedger.Accounts)
	// p.LocalLedger.ApplySignedTransaction(signedTransaction)
	p.debugPrintln("Adding signed-transaction:", signedTransaction.ToString(), "to signed transactions seen")//, ". Accounts after:", p.LocalLedger.Accounts)
	p.SignedTransactionsSeen.Put(signedTransaction.ID, signedTransaction)

	p.Broadcast(
		Packet{
			Type: PacketType.BROADCAST_SIGNED_TRANSACTION,
			SignedTransaction: signedTransaction,
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

