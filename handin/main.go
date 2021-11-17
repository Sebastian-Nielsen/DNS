package main

import (
	. "DNO/handin/Account"
	. "DNO/handin/Cryptography"
	. "DNO/handin/Helper"
	. "DNO/handin/Peernode"
	"net"
)


func main() {
	var peerNode = PeerNode{
		OpenConnections:     	SafeSet_Conn{ Vals: make(map[net.Conn]bool) },
		PeersInArrivalOrder: 	SafeArray_string{},
		MessagesSent:        	SafeSet_string{ Values: make(map[string  ]bool) },
		Ipc:                 	IPC{ ConnToEncDecPair: make(map[net.Conn]EncoderDecoderPair) },
		TestMock:            	Mock{ ShouldMockInput: false, ShouldPrintDebug: false },
		Keys:                	GenKeyPair(),
		UnappliedIDs:        	SafeArray_string{},
		SignedTransactionsSeen: SafeMap_string_to_SignedTransaction{ Values: make(map[string] SignedTransaction) },
		Sequencer:			 	Sequencer{
									UnsequencedTransactionIDs: SafeArray_string{},
								 	PublicKey:                 PublicKey{},
								 	KeyPair:                   KeyPair{},
								 	SlotNumber:                SafeCounter{Value: -1},
								},
	}

	peerNode.Start("", "")

	println("---------------------\nTerminating...")
}

