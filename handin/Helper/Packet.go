package Helper

import "DNO/handin/Account"

type Packet struct {
	Type string
	Msg  string
	MessagesSent map[string]bool
	PeersInArrivalOrderValues []string

	Transaction Account.Transaction
	TransactionsSeen []Account.Transaction

	SignedTransaction Account.SignedTransaction
	SignedTransactionsSeen []Account.SignedTransaction
}

var PacketType = struct {
	BROADCAST_KNOWN_TRANSACTION string
	BROADCAST_TRANSACTION string
	BROADCAST_KNOWN_SIGNED_TRANSACTION string
	BROADCAST_SIGNED_TRANSACTION string
	BROADCASTED_KNOWN_LISTENER_PORT string
	BROADCAST_LISTENER_PORT string
	LISTENER_PORT string
	BROADCAST_MSG string
	PULL          string
	PULL_REPLY string
} {
	BROADCAST_KNOWN_TRANSACTION: "BROADCAST_KNOWN_TRANSACTION",
	BROADCAST_TRANSACTION: "BROADCAST_TRANSACTION",
	BROADCAST_KNOWN_SIGNED_TRANSACTION: "BROADCAST_KNOWN_SIGNED_TRANSACTION",
	BROADCAST_SIGNED_TRANSACTION: "BROADCAST_SIGNED_TRANSACTION",
	BROADCASTED_KNOWN_LISTENER_PORT: "BROADCASTED_KNOWN_LISTENER_PORT",
	BROADCAST_LISTENER_PORT: "BROADCAST_LISTENER_PORT",
	LISTENER_PORT:           "LISTENER_PORT",
	BROADCAST_MSG:           "BROADCAST_MSG",
	PULL:                    "PULL",
	PULL_REPLY:              "PULL_REPLY",
}
