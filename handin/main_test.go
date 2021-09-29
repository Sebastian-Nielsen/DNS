package main

import (
	. "DNO/handin/Account"
	. "DNO/handin/Cryptography"
	. "DNO/handin/Helper"
	"fmt"
	"math/big"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"
)

/*
	sometimes debug_mode has to be set to false for tests to pass for some mystical reason
*/

func createPeerNode( shouldMockInput bool, shouldPrintDebug bool ) PeerNode {
	return PeerNode{
		OpenConnections:     SafeSet_Conn{   Values: make(map[net.Conn]bool) },
		PeersInArrivalOrder: SafeArray_string{},
		MessagesSent:        SafeSet_string{ Values: make(map[string  ]bool) },
		Ipc:                 IPC{ ConnToEncDecPair: make(map[net.Conn]EncoderDecoderPair) },
		LocalLedger:         MakeLedger(),
		TestMock:            Mock{ ShouldMockInput: shouldMockInput, 
								   ShouldPrintDebug: shouldPrintDebug },
	}
}





func TestUseThisToDebug(t *testing.T) {
	c := CTR{
		SecretKey: "6368616e676520746869732070617373",
	}

	bytesToBeEncrypted := []byte("hello there")

	encryptedBytes := c.Encrypt(bytesToBeEncrypted)
	decryptedBytes := c.Decrypt(encryptedBytes)

	fmt.Println(bytesToBeEncrypted)
	fmt.Println([]byte(string(decryptedBytes)))
}

func TestRSAsigning(t *testing.T) {
	n, d := KeyGen(2000)
	publicKey := PublicKey{N:n, E:big.NewInt(3)}
	secretKey := SecretKey{N:n, D:d}
	msg := big.NewInt( 25632212678324 )

	// Sign the message with the secret key
	signedMsg := Sign(msg, secretKey)
	
	// Hash the message to use when verifying 
	hashedMsg := new(big.Int)
	hashedMsg.SetBytes(GetHash(msg))

	// Verify the signed message against the hash using the public key
	verified := Verify(signedMsg, hashedMsg, publicKey)

	if !verified {
		t.Error("Hash from signed message '" + Encrypt(signedMsg, publicKey).String() + "' was not verified for original hash: " + hashedMsg.String())
	}
}

func TestEncryptionAndDecryptionWithRSAandAES(t *testing.T) {
	// RSA decrypt
	n, d := KeyGen(2000)
	publicKey := PublicKey{N:n, E:big.NewInt(3)}
	secretKey := SecretKey{N:n, D:d}
	msg := big.NewInt( 8795378532487390 )
	RSAmsg := Encrypt(msg, publicKey)

	// AES encrypt the secret key
	cbc := CTR{SecretKey: GenerateNewRndmIV(32)}
	filename := "Cryptography/RSAandAEStest"
	secretKeyString := secretKey.N.String() + ":" + secretKey.D.String()
	cbc.EncryptToFile(filename, secretKeyString)

	// AES decrypt the secret key
	decryptionFromFile := cbc.DecryptFromFile(filename)

	// Create secret key from AES decryption of the file
	splitPos := strings.Index(decryptionFromFile, ":")
	decrN := new(big.Int)
	decrN.SetString(decryptionFromFile[:splitPos], 10)
	decrD := new(big.Int)
	decrD.SetString(decryptionFromFile[splitPos+1:], 10)
	decryptedSecretKey := SecretKey{N:decrN, D:decrD}

	// RSA decrypt with the new secret key
	decryptedMsg := Decrypt(RSAmsg, decryptedSecretKey)

	if decryptedMsg.String() != msg.String() {
		t.Error("Original message '" + msg.String() + "' different from decrypted message'" + decryptedMsg.String() + "'")
	}
}
func TestNewcomerNodeReceivesAllTransactionsAppliedBeforeItEnteredNetwork(t *testing.T) {
	t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()
	var peerNode3_port = AvailablePorts.Next()
	var peerNode4_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true, true)
	peerNode2 := createPeerNode(true, true)
	peerNode3 := createPeerNode(true, true)
	peerNode4 := createPeerNode(true, true)

	goStart(&peerNode1, peerNode1_port, "no_port")
	goStart(&peerNode2, peerNode2_port, peerNode1_port)
	goStart(&peerNode3, peerNode3_port, peerNode2_port)

	IDs 	:= []string{"acc1", "acc2", "acc3"}
	amounts := []int{      100,      0, 40}
	makeAccounts(&peerNode1, IDs, amounts)
	makeAccounts(&peerNode2, IDs, amounts)
	makeAccounts(&peerNode3, IDs, amounts)

	//time.Sleep(4 * time.Second)

	peerNode1.MakeAndBroadcastTransaction(42, "tran1", IDs[0], IDs[1])

	time.Sleep(1 * time.Second)

	peerNode2.MakeAndBroadcastTransaction(5, "tran2", IDs[1], IDs[2])

	time.Sleep(1 * time.Second)

	makeAccounts(&peerNode4, IDs, amounts)
	goStart(&peerNode4, peerNode4_port, peerNode2_port)

	time.Sleep(1 * time.Second)

	if peerNode4.LocalLedger.Accounts[IDs[0]] != amounts[0]-42{
		t.Error("Transaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc1' on latecomer peerNode4.",
			"\npeerNode4 Accounts:", peerNode4.LocalLedger.Accounts, "\n")
	}
	if peerNode4.LocalLedger.Accounts[IDs[1]] != (amounts[1]+42) - 5 {
		t.Error("Transaction 'tran1' and 'tran2' from 'acc1' to 'acc2' and 'acc2' to 'acc3' on peerNode1 didn't apply to 'acc2' on latecomer peerNode4.",
			"\npeerNode4 Accounts:", peerNode4.LocalLedger.Accounts, "\n")
	}
	if peerNode4.LocalLedger.Accounts[IDs[2]] != amounts[2]+5 {
		t.Error("Transaction 'tran3' from 'acc2' to 'acc3' on peerNode1 didn't apply to 'acc3' on latecomer peerNode4.",
			"\npeerNode4 Accounts:", peerNode4.LocalLedger.Accounts, "\n")
	}
	if len(peerNode4.TransactionsSeen.Values()) != 2 {
		t.Error("Not all transactions was seen by latecomer peerNode4",
			"\npeerNode4 Accounts:", peerNode4.TransactionsSeen.Values(), "\n")
	}
}
func TestTransactionsAreBroadcastedAndAppliedOnAllPeerNodes(t *testing.T) {
	t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()
	var peerNode3_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true, true)
	peerNode2 := createPeerNode(true, true)
	peerNode3 := createPeerNode(true, true)

	goStart(&peerNode1, peerNode1_port, "no_port")
	goStart(&peerNode2, peerNode2_port, peerNode1_port)
	goStart(&peerNode3, peerNode3_port, peerNode2_port)

	IDs 	:= []string{"acc1", "acc2", "acc3"}
	amounts := []int{      100,      0,     40}
	makeAccounts(&peerNode1, IDs, amounts)
	makeAccounts(&peerNode2, IDs, amounts)
	makeAccounts(&peerNode3, IDs, amounts)

	time.Sleep(1 * time.Second)

	peerNode1.MakeAndBroadcastTransaction(20, "tran1", IDs[0], IDs[1])

	time.Sleep(1 * time.Second)

	if peerNode1.LocalLedger.Accounts[IDs[0]] != amounts[0]-20 {
		t.Error("Transaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to acc1' on peerNode1.",
			"\nnpeerNode1 Accounts:", peerNode1.LocalLedger.Accounts, "\n")
	}
	if peerNode1.LocalLedger.Accounts[IDs[1]] != amounts[1]+20 {
		t.Error("Transaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc2' on peerNode1.",
			"\nnpeerNode1 Accounts:", peerNode1.LocalLedger.Accounts, "\n")
	}
	if peerNode3.LocalLedger.Accounts[IDs[0]] != amounts[0]-20 {
		t.Error("Transaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc1' on peerNode3.",
			"\nnpeerNode1 Accounts:", peerNode3.LocalLedger.Accounts, "\n")
	}
	if peerNode3.LocalLedger.Accounts[IDs[1]] != amounts[1]+20 {
		t.Error("Transaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc2' on peerNode3.",
			"\nnpeerNode1 Accounts:", peerNode3.LocalLedger.Accounts, "\n")
	}
}
func TestNodeConnectsToThreeOthersWhenEnteringNetwork(t *testing.T) {
	t.Parallel()

	peerNode1 := createPeerNode(true, true)
	peer1Port := AvailablePorts.Next()
	goStart(&peerNode1, peer1Port, "no_port")

	var peerNode2 = createPeerNode(true, true)
	var peerNode3 = createPeerNode(true, true)
	var peerNode4 = createPeerNode(true, true)
	
	goStart(&peerNode2, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode3, AvailablePorts.Next(), peer1Port)
	peer4Port := AvailablePorts.Next()
	goStart(&peerNode4, peer4Port, peer1Port)

	time.Sleep(1 * time.Second)

	if len(peerNode1.OpenConnections.Values) != 3 {
		t.Error("PeerNode1 doesn't have exactly 3 connections")
	}
	if len(peerNode2.OpenConnections.Values) != 3 {
		t.Error("PeerNode2 doesn't have exactly 3 connections.\npeerNode2 openConnections:",
			peerNode2.OpenConnections.ToString())
	}
	if len(peerNode3.OpenConnections.Values) != 3 {
		t.Error("PeerNode3 doesn't have exactly 3 connections")
	}
	if len(peerNode4.OpenConnections.Values) != 3 {
		t.Error("PeerNode4 doesn't have exactly 3 connections")
	}

	if !peerNode4.PeersInArrivalOrder.Contains(peer4Port) {
		t.Error("peerNode4 (the last one connected to the network) isn't the last node in PeersInArrivalOrder:", 
			peerNode4.PeersInArrivalOrder.Values())
	}
}
func TestNodeConnectsToTenOthersWhenEnteringNetwork(t *testing.T) {
	t.Parallel()

	peerNode1 := createPeerNode(true, true)
	peer1Port := AvailablePorts.Next()
	goStart(&peerNode1, peer1Port, "no_port")

	var peerNode2 = createPeerNode(true, true)
	var peerNode3 = createPeerNode(true, true)
	var peerNode4 = createPeerNode(true, true)
	var peerNode5 = createPeerNode(true, true)
	var peerNode6 = createPeerNode(true, true)
	var peerNode7 = createPeerNode(true, true)
	var peerNode8 = createPeerNode(true, true)
	var peerNode9 = createPeerNode(true, true)
	var peerNode10 = createPeerNode(true, true)
	var peerNode11 = createPeerNode(true, true)

	goStart(&peerNode2, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode3, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode4, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode5, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode6, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode7, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode8, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode9, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode10, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode11, AvailablePorts.Next(), peer1Port)

	time.Sleep(1 * time.Second)

	// By now, all nodefs should have 10 openConnections, just test that 4 of them have it
	if len(peerNode1.OpenConnections.Values) != 10 {
		t.Error("PeerNode1 doesn't have exactly 10 connections")
	}
	if len(peerNode2.OpenConnections.Values) != 10 {
		t.Error("PeerNode2 doesn't have exactly 10 connections")
	}
	if len(peerNode10.OpenConnections.Values) != 10 {
		t.Error("PeerNode3 doesn't have exactly 10 connections")
	}
	if len(peerNode11.OpenConnections.Values) != 10 {
		t.Error("PeerNode4 doesn't have exactly 10 connections")
	}

	// A new node (peerNode12) joins the network
	p12Port := AvailablePorts.Next()
	peerNode12 := createPeerNode(true, true)
	goStart(&peerNode12, p12Port, peer1Port)

	time.Sleep(1 * time.Second)

	if !peerNode11.PeersInArrivalOrder.Contains(p12Port) {
		t.Error("peerNode11 doesn't have the port of peerNode12 (" + p12Port + ") after peerNode1 broadcast it.\n" +
			"PeerNode11's portList:", peerNode11.PeersInArrivalOrder.Values())
	}
}
func TestPeerNodeConnectsToAllNodesWhenEnteringNetwork(t *testing.T) {
	t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()
	var peerNode3_port = AvailablePorts.Next()
	var peerNode4_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true, true)
	peerNode2 := createPeerNode(true, true)
	peerNode3 := createPeerNode(true, true)
	peerNode4 := createPeerNode(true, true)

	goStart(&peerNode1, peerNode1_port, "no_port")
	goStart(&peerNode2, peerNode2_port, peerNode1_port)
	goStart(&peerNode3, peerNode3_port, peerNode1_port)
	goStart(&peerNode4, peerNode4_port, peerNode1_port)

	time.Sleep(4 * time.Second)
	p4Conns := peerNode4.OpenConnections.Values
	if len(p4Conns) != 3 {
		t.Error("peerNode4 doesn't have an open connection to each node in the network.\n" +
			"PeerNode4's openConnections:", p4Conns)
	}
}
func TestReceivedPeerListWhenJoining(t *testing.T) {
	t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()
	var peerNode3_port = AvailablePorts.Next()
	var peerNode4_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true, true)
	peerNode2 := createPeerNode(true, true)
	peerNode3 := createPeerNode(true, true)
	peerNode4 := createPeerNode(true, true)

	goStart(&peerNode1, peerNode1_port, "no_port")
	goStart(&peerNode2, peerNode2_port, peerNode1_port)
	goStart(&peerNode3, peerNode3_port, peerNode1_port)
	goStart(&peerNode4, peerNode4_port, peerNode1_port)

	time.Sleep(2 * time.Second)

	expectedList := []string{peerNode1_port, peerNode2_port, peerNode3_port, peerNode4_port}
	if len(peerNode4.PeersInArrivalOrder.Values()) == 0 {
		t.Errorf("peerNode4 didn't receive peerNode1's list of peers")
	}
	if !reflect.DeepEqual(peerNode4.PeersInArrivalOrder.Values(), expectedList) {
		t.Error("peerNode4 received peerNode1's list of peers in the wrong order:\npeerNode4:",
			peerNode4.PeersInArrivalOrder.Values(), "\nvs\npeerNode1 (expected list):", expectedList)
	}


}
func TestMessageIsSentFromNode1To3Via2(t *testing.T) {
	t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()
	var peerNode3_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true, true)
	peerNode2 := createPeerNode(true, true)
	peerNode3 := createPeerNode(true, true)


	goStart(&peerNode1, peerNode1_port, "no_port")
	goStart(&peerNode2, peerNode2_port, peerNode1_port)
	goStart(&peerNode3, peerNode3_port, peerNode2_port)


	simulateInputFor(&peerNode1, "some_msg")

	if !peerNode3.MessagesSent.Contains("some_msg") {
		t.Errorf("peerNode3 didn't receive peerNode1's msgs")
	}
}
func TestLatercomerNodeEventuallyGetsAllMsgs(t *testing.T) {
	t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true, true)
	peerNode2 := createPeerNode(true, true)

	goStart(&peerNode1, peerNode1_port, "no_port")

	simulateInputFor(&peerNode1, "some_msg_1")
	simulateInputFor(&peerNode1, "some_msg_2")
	simulateInputFor(&peerNode1, "some_msg_3")

	goStart(&peerNode2, peerNode2_port, peerNode1_port)

	if  !peerNode2.MessagesSent.Contains("some_msg_1") ||
		!peerNode2.MessagesSent.Contains("some_msg_2") ||
		!peerNode2.MessagesSent.Contains("some_msg_3") {
		t.Errorf("peerNode2 didn't receive peerNode1's msgs")
	}
}
func TestPeer1ReceivesMsgFromPeer2(t *testing.T) {
	t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true, true)
	peerNode2 := createPeerNode(true, true)

	go peerNode1.Start(peerNode1_port, "no_port")
	go peerNode2.Start(peerNode2_port, peerNode1_port)

	time.Sleep(4 * time.Second)

	simulateInputFor(&peerNode1, "some_msg")

	if !peerNode2.MessagesSent.Contains("some_msg") {
		t.Errorf("peerNode2 didn't receive peerNode1's msg")
	}

}
func TestPeer1CanConnectToPeer2(t *testing.T) {
	t.Parallel()


	peerNode1 := createPeerNode(true, true)
	peerNode2 := createPeerNode(true, true)

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()

	goStart(&peerNode1, peerNode1_port, "no_port")
	goStart(&peerNode2, peerNode2_port, peerNode1_port)

	if len(peerNode1.OpenConnections.Values) != 1 {
		t.Error("peerNode1 doesn't have peerNode2 in its openConnections set:", peerNode1.OpenConnections.Values)
	}
	if len(peerNode2.OpenConnections.Values) != 1 {
		t.Error("peerNode2 doesn't have peerNode1 in its openConnections set:", peerNode2.OpenConnections.Values)
	}
}



func simulateInputFor(peerNode *PeerNode, text string) {
	peerNode.TestMock.SimulatedInputString = text
	time.Sleep(1 * time.Second)
}
func goStart(peerNode *PeerNode, atPort string, remotePort string) {
	go peerNode.Start(atPort, remotePort)
	time.Sleep(250 * time.Millisecond)
}

func makeAccounts(peerNode *PeerNode, IDs []string, startAmounts []int){
	for i, id := range IDs {
		peerNode.CreateAccountInLedger(id, startAmounts[i])
	}
}

