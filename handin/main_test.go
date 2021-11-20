package main

import (
	. "DNO/handin/Account"
	. "DNO/handin/Cryptography"
	. "DNO/handin/Helper"
	. "DNO/handin/Peernode"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

/*
	sometimes debug_mode has to be set to false for tests to pass for some mystical reason
*/

func createPeerNode( shouldMockInput bool, shouldPrintDebug bool ) PeerNode {
	return PeerNode{
		OpenConnections:     SafeSet_Conn{   Vals: make(map[net.Conn]bool) },
		PeersInArrivalOrder: SafeArray_string{},
		MessagesSent:        SafeSet_string{ Values: make(map[string  ]bool) },
		Ipc:                 IPC{ ConnToEncDecPair: make(map[net.Conn]EncoderDecoderPair) },
		LocalLedger:         MakeLedger(),
		TestMock:            Mock{ ShouldMockInput: shouldMockInput,
								   ShouldPrintDebug: shouldPrintDebug },
	    Keys: 				 GenKeyPair(),
		UnappliedIDs:        SafeArray_string{},
		SignedTransactionsSeen:  SafeMap_string_to_SignedTransaction{ Values: make(map[string] SignedTransaction) },
		Sequencer:			  Sequencer{
								 UnsequencedTransactionIDs: SafeArray_string{},
								 //PublicKey:                 PublicKey{},
								 //KeyPair:                   KeyPair{},
								 IsSequencer: false,
								 SlotNumber:                SafeCounter{Value: -1},
								 Tree:                      Tree{
															   BlockHashToBlock: SafeMap_string_to_Block{ Values: make(map[string]Block)},
															   BlocksThatAreWaitingForTheirParent: make(map[string]*SafeArray_Block),
								 						    },
							 },
	}
}


func Test1(t *testing.T) {

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()
	var peerNode3_port = AvailablePorts.Next()
	var peerNode4_port = AvailablePorts.Next()
	var peerNode5_port = AvailablePorts.Next()
	var peerNode6_port = AvailablePorts.Next()
	var peerNode7_port = AvailablePorts.Next()
	var peerNode8_port = AvailablePorts.Next()
	var peerNode9_port = AvailablePorts.Next()
	var peerNode10_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true, false)
	peerNode2 := createPeerNode(true, true)
	peerNode3 := createPeerNode(true, false)
	peerNode4 := createPeerNode(true, false)
	peerNode5 := createPeerNode(true, false)
	peerNode6 := createPeerNode(true, false)
	peerNode7 := createPeerNode(true, false)
	peerNode8 := createPeerNode(true, false)
	peerNode9 := createPeerNode(true, false)
	peerNode10 := createPeerNode(true, false)


	keypairs := GetHardcodedAccKeyPairs()
	peerNode1.Keys = keypairs[0]
	peerNode2.Keys = keypairs[1]
	peerNode3.Keys = keypairs[2]
	peerNode4.Keys = keypairs[3]
	peerNode5.Keys = keypairs[4]
	peerNode6.Keys = keypairs[5]
	peerNode7.Keys = keypairs[6]
	peerNode8.Keys = keypairs[7]
	peerNode9.Keys = keypairs[8]
	peerNode10.Keys = keypairs[9]

	goStart(&peerNode1, peerNode1_port, "no_port")
	goStart(&peerNode2, peerNode2_port, peerNode1_port)
	goStart(&peerNode3, peerNode3_port, peerNode1_port)
	goStart(&peerNode4, peerNode4_port, peerNode1_port)
	goStart(&peerNode5, peerNode5_port, peerNode1_port)
	goStart(&peerNode6, peerNode6_port, peerNode1_port)
	goStart(&peerNode7, peerNode7_port, peerNode1_port)
	goStart(&peerNode8, peerNode8_port, peerNode1_port)
	goStart(&peerNode9, peerNode9_port, peerNode1_port)
	goStart(&peerNode10, peerNode10_port, peerNode1_port)

	time.Sleep(5 * time.Second)
	peerNode1.MakeAndBroadcastSignedTransaction(
		500,
		PortOf(peerNode1.Listener.Addr())+ ":" + strconv.Itoa(peerNode1.Sequencer.SlotNumber.Value),
		peerNode1.Keys.Pk.ToString(),
		peerNode2.Keys.Pk.ToString(),
	)
	time.Sleep(10 * time.Second)
	for _, p := range []*PeerNode{&peerNode1, &peerNode2, &peerNode3, &peerNode4, &peerNode5,
								  &peerNode6, &peerNode7, &peerNode8, &peerNode9, &peerNode10} {
		fmt.Println(PortOf(p.Listener.Addr()))
		for key, value := range p.LocalLedger.Accounts {
			fmt.Println("[" + key[0:11] + " ; " + strconv.Itoa(value) + "]")
		}
	}
	time.Sleep(10 * time.Second)
	for _, p := range []*PeerNode{&peerNode1, &peerNode2, &peerNode3, &peerNode4, &peerNode5,
		&peerNode6, &peerNode7, &peerNode8, &peerNode9, &peerNode10} {
		fmt.Println(PortOf(p.Listener.Addr()))
		for key, value := range p.LocalLedger.Accounts {
			fmt.Println("[" + key[0:11] + " ; " + strconv.Itoa(value) + "]")
		}
	}
	time.Sleep(10 * time.Second)
	for _, p := range []*PeerNode{&peerNode1, &peerNode2, &peerNode3, &peerNode4, &peerNode5,
		&peerNode6, &peerNode7, &peerNode8, &peerNode9, &peerNode10} {
		fmt.Println(PortOf(p.Listener.Addr()))
		for key, value := range p.LocalLedger.Accounts {
			fmt.Println("[" + key[0:11] + " ; " + strconv.Itoa(value) + "]")
		}
	}
	//peerNode.MakeAndBroadcastSignedTransaction(t.Amount, t.ID+strconv.Itoa(i), t.From, t.To)
	time.Sleep(500 * time.Second)
}






func TestGenerateInitalAccountKeys(t *testing.T) {
	for i := 0; i < 10; i++ {
		pk, sk := GenKeys(2000)
		fmt.Println("KeyPair{ Pk: ToPublicKey(\"" + pk.ToString() + "\"),\nSk:ToSecretKey(\"" + sk.ToString() + "\")},\n")
	}


}

/* Legacy test
func TestConcurrentTransactions(t *testing.T) {
	//t.Parallel()

	var peerNodeA_port = AvailablePorts.Next() // The sequencer
	var peerNodeB_port = AvailablePorts.Next()
	var peerNodeC_port = AvailablePorts.Next()
	
	peerNodeA := createPeerNode(true, false)
	peerNodeB := createPeerNode(true, false)
	peerNodeC := createPeerNode(true, false)
	
	goStart(&peerNodeA, peerNodeA_port, "no_port")
	goStart(&peerNodeB, peerNodeB_port, peerNodeA_port)
	goStart(&peerNodeC, peerNodeC_port, peerNodeA_port)

	AccountNames := []string{
		peerNodeA.Keys.Pk.ToString(),
		peerNodeB.Keys.Pk.ToString(),
		peerNodeC.Keys.Pk.ToString(),
	}
	initialAmounts := []int{500, 0, 0}
	makeAccountsAt(&peerNodeA, AccountNames, initialAmounts)
	makeAccountsAt(&peerNodeB, AccountNames, initialAmounts)
	makeAccountsAt(&peerNodeC, AccountNames, initialAmounts)

	peerNodeB.Keys = peerNodeA.Keys
	peerNodeC.Keys = peerNodeA.Keys

	// wait for phase 1 to finish
	for peerNodeA.Sequencer.PublicKey.N == nil ||
		peerNodeB.Sequencer.PublicKey.N == nil ||
		peerNodeC.Sequencer.PublicKey.N == nil {
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Sequencer pk A:", peerNodeA.Sequencer.PublicKey.ToString()[:10])
	fmt.Println("Sequencer pk B:", peerNodeB.Sequencer.PublicKey.ToString()[:10])
	fmt.Println("Sequencer pk C:", peerNodeC.Sequencer.PublicKey.ToString()[:10])

	var wg sync.WaitGroup
	wg.Add(2)

	go MakeManyTransactions(&peerNodeB, Transaction {Amount: 1, ID: "fromAtoB", From: AccountNames[0], To: AccountNames[1]}, &wg)
	go MakeManyTransactions(&peerNodeC, Transaction {Amount: 1, ID: "fromAToC", From: AccountNames[0], To: AccountNames[2]}, &wg)

	wg.Wait() // wait for the goroutines to finish

	peerNodeA.BroadcastBlock()

	time.Sleep(1 * time.Second)

	ledgerB := peerNodeB.LocalLedger.Accounts
	ledgerC := peerNodeC.LocalLedger.Accounts
	allAccountsAreEqual := reflect.DeepEqual(ledgerB, ledgerC)
	if !allAccountsAreEqual {
		t.Error("Accounts in peerNodeB:", ledgerB, "are not the same as accounts in peerNodeC:", ledgerC)
	}
	amountInAccountA := peerNodeA.LocalLedger.Accounts[AccountNames[0]]
	if amountInAccountA != 0 {
		t.Error("Amount in account A is", amountInAccountA, "but should be 0")
	}

	fmt.Println("Account B has:", peerNodeA.LocalLedger.Accounts[AccountNames[1]])
	fmt.Println("Account C has:", peerNodeA.LocalLedger.Accounts[AccountNames[2]])
}
*/

func MakeManyTransactions(peerNode *PeerNode, t Transaction, wg *sync.WaitGroup) {
	for i := 1; i < 501; i++ {
		peerNode.MakeAndBroadcastSignedTransaction(t.Amount, t.ID+strconv.Itoa(i), t.From, t.To)
		if (i % 100 == 0) {
			fmt.Println("Sent " + strconv.Itoa(i) + " '" + t.ID + "' transactions so far")
		}
	}
	wg.Done()
	fmt.Println("Sent all '" + t.ID + "' transactions")
}


// ============================================ crypto tests ========= START
func TestGenerateAndSignature(t *testing.T) {
	validPassword := "AB$45btc"
	filename := "Cryptography/wallet_test"
	msg := "this is the msg"

	pk, _ := Generate(filename, validPassword)

	// Try to get signature with correct validPassword
	signature, _ := Sign(filename, validPassword, []byte(msg))
	isValidSignature := !Verify(signature, msg, pk)
	if !isValidSignature {
		t.Error("Signature was not verified")
	}

	// Try to get signature with invalid validPassword
	_, err := Sign(filename, validPassword, []byte(msg))
	if err != nil {
		t.Error("A signature was generated with the wrong Password")
	}

	// Try to Generate wallet using a weak validPassword
	_, err = Generate("file", "weak password")
	if err == nil {
		t.Error("Accepted a weak password for encryption")
	}
}
func TestRSAsigningTime(t *testing.T) {
	_, sk := GenKeys(200)
	msg := big.NewInt( 25632212678324 )
	fmt.Println("d:", sk.D.BitLen())

	hashedMsg := new(big.Int)
	hashedMsg.SetBytes(Hash(msg))

	// Time the signing of a hash. Since our signing function computes the hash
	// we just manually compute the hash and the time just the signing (which is the
	// same as what our Decrypt function does)
	startTime := time.Now()
	Decrypt(hashedMsg, sk)
	timeElapsed := time.Since(startTime)

	bitsPerSec := float64(hashedMsg.BitLen() + 1) / timeElapsed.Seconds()
	fmt.Printf("It took %s to sign a hashed message of size 256 bits = %f bits per sec\n", timeElapsed.String(), bitsPerSec)
}
func TestHashingTime(t *testing.T) {
	data := make([]byte, 10000)
	rand.Read(data)
	msg := new(big.Int)
	msg.SetBytes(data)

	// Time the hashing of the message
	startTime := time.Now()
	Hash(msg)
	timeElapsed := time.Since(startTime)

	bitsPerSec := float64(msg.BitLen()) / timeElapsed.Seconds()
	fmt.Printf("It took %s to hash a message of size 10.000 bytes = %f bits per sec\n", timeElapsed.String(), bitsPerSec)
}
func TestRSASigning(t *testing.T) {
	pk, sk := GenKeys(2000)
	msg := big.NewInt( 25632212678324 )

	// CreateSignature the message with the secret key
	signature := BigInt_createSignature(msg, sk)

	// BigInt_verify the signed message against msg using the public key
	verified := BigInt_verify(signature, msg, pk)
	if !verified {
		t.Error("Signed message '" + Encrypt(signature, pk).String() +
			"' was not verified for original message: " + msg.String())
	}

	// Modify the msg by adding 1
	modifiedMsg := big.NewInt(0).Add(msg, big.NewInt(1))

	// Since the signature is modified, we shouldn't get verification
	verified = BigInt_verify(signature, modifiedMsg, pk)
	if verified {
		t.Error("The modified message '" + Encrypt(modifiedMsg, pk).String() +
			"' was wrongly verified for original message: " + msg.String())
	}
}
func TestEncryptionAndDecryptionWithRSAandAES(t *testing.T) {
	// RSA decrypt
	pk, sk := GenKeys(2000)
	msg := big.NewInt( 8795378532487390 )
	RSAEncryptedMsg := Encrypt(msg, pk)

	// AES encrypt the secret key
	ctr := CTR{SecretKey: GenerateNewRndmString(32)}
	filename := "Cryptography/RSAandAEStest"
	secretKeyString := sk.N.String() + ":" + sk.D.String()
	ctr.EncryptToFile(filename, secretKeyString)

	// AES decrypt the secret key
	decryptionFromFile := ctr.DecryptFromFile(filename)

	// Create secret key from AES decryption of the file
	splitPos := strings.Index(decryptionFromFile, ":")
	decrN := new(big.Int)
	decrN.SetString(decryptionFromFile[:splitPos], 10)
	decrD := new(big.Int)
	decrD.SetString(decryptionFromFile[splitPos+1:], 10)
	decryptedSecretKey := SecretKey{N:decrN, D:decrD}

	// RSA decrypt with the new secret key
	decryptedMsg := Decrypt(RSAEncryptedMsg, decryptedSecretKey)

	if decryptedMsg.String() != msg.String() {
		t.Error("Original message '" + msg.String() + "' different from decrypted message'" + decryptedMsg.String() + "'")
	}
}
// ============================================ crypto tests ========= END


// Legacy test from handin_8 and onwards
/*
func TestSignedTransactionsAreBroadcastedAndAppliedOnAllPeerNodes(t *testing.T) {
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

	AccountNames := []string{peerNode1.Keys.Pk.ToString(),
							 peerNode2.Keys.Pk.ToString(),
							 peerNode3.Keys.Pk.ToString()}
	initialAmounts := []int{100, 0, 40}
	makeAccountsAt(&peerNode1, AccountNames, initialAmounts)
	makeAccountsAt(&peerNode2, AccountNames, initialAmounts)
	makeAccountsAt(&peerNode3, AccountNames, initialAmounts)

	time.Sleep(1 * time.Second)

	// Test all signed transactions got received and applied
	peerNode1.MakeAndBroadcastSignedTransaction(20, "tran1", AccountNames[0], AccountNames[1])

	time.Sleep(1 * time.Second)

	if peerNode1.LocalLedger.Accounts[AccountNames[0]] != initialAmounts[0]-20 {
		t.Error("ApplyTransaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc1' on peerNode1.",
			"\nnpeerNode1 Accounts:", peerNode1.LocalLedger.Accounts, "\n")
	}
	if peerNode1.LocalLedger.Accounts[AccountNames[1]] != initialAmounts[1]+20 {
		t.Error("ApplyTransaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc2' on peerNode1.",
			"\nnpeerNode1 Accounts:", peerNode1.LocalLedger.Accounts, "\n")
	}
	if peerNode3.LocalLedger.Accounts[AccountNames[0]] != initialAmounts[0]-20 {
		t.Error("ApplyTransaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc1' on peerNode3.",
			"\nnpeerNode1 Accounts:", peerNode3.LocalLedger.Accounts, "\n")
	}
	if peerNode3.LocalLedger.Accounts[AccountNames[1]] != initialAmounts[1]+20 {
		t.Error("ApplyTransaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc2' on peerNode3.",
			"\nnpeerNode1 Accounts:", peerNode3.LocalLedger.Accounts, "\n")
	}

	// Test that transaction with negative amount is denied
	peerNode2.MakeAndBroadcastSignedTransaction(-10, "tran2", AccountNames[1], AccountNames[2])

	time.Sleep(1 * time.Second)

	if peerNode1.LocalLedger.Accounts[AccountNames[2]] != initialAmounts[2] {
		t.Error("ApplyTransaction 'tran2' from 'acc2' to 'acc3' on peerNode2 applied to 'acc3' on peerNode1 " +
			"despite the amount being negative.",
			"\npeerNode1 Accounts:", peerNode1.LocalLedger.Accounts, "\n")
	}

	// Test that a transaction with a wrong key is denied (peernode2 tries to make a transaction
	// from peernode3's account)
	peerNode2.MakeAndBroadcastSignedTransaction(10, "tran3", AccountNames[2], AccountNames[1])

	time.Sleep(1 * time.Second)

	if peerNode1.LocalLedger.Accounts[AccountNames[2]] != initialAmounts[2] {
		t.Error("ApplyTransaction 'tran3' from 'acc3' to 'acc2' on peerNode2 applied to 'acc3' on peerNode1 " +
			"despite the signed transaction having the wrong secret key.",
			"\npeerNode1 Accounts:", peerNode1.LocalLedger.Accounts, "\n")
	}
}
*/
func TestSignedTransactions(t *testing.T) {
	// t.Parallel()

	ledger1 := MakeLedger()
	id_1, pk_1, sk_1 := createAccountAt(ledger1)

	ledger2 := MakeLedger()
	_, pk_2, _ := createAccountAt(ledger2)

	// Account_1 at ledger_1 creates a signed transaction
	transaction := Transaction{ID: id_1, From: pk_1.ToString(), To: pk_2.ToString(), Amount: 2}
	signedTransaction := ledger1.MakeSignedTransaction(transaction, sk_1)

	// Account_2 at ledger_2 verifies the signed transaction
	isVerified := ledger2.Verify(signedTransaction)
	if !isVerified {
		t.Error("\nError: the signedTransaction (" + signedTransaction.ToString() + "\n) should be verified")
	}

}
func TestCreateSignatureAndVerifyFunctions(t *testing.T) {
	pk, sk := GenKeys(2000)
	msg := "test123123"
	signature := CreateSignature(msg, sk)
	isVerified := Verify(signature, msg, pk)
	if !isVerified {
		t.Error("Verify function on signature (" + signature[:10] + "...) was false for msg: " + msg)
	}
}
func createAccountAt(ledger *Ledger) (string, PublicKey, SecretKey) {
	pk, sk := GenKeys(2000)
	id := CreateSignature(pk.ToString(), sk)
	ledger.CreateAccount(id, 10)
	return id, pk, sk
}

// Legacy test from handin_8 and onwards
/*
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
	makeAccountsAt(&peerNode1, IDs, amounts)
	makeAccountsAt(&peerNode2, IDs, amounts)
	makeAccountsAt(&peerNode3, IDs, amounts)

	//time.Sleep(4 * time.Second)

	peerNode1.MakeAndBroadcastTransaction(42, "tran1", IDs[0], IDs[1])

	time.Sleep(1 * time.Second)

	peerNode2.MakeAndBroadcastTransaction(5, "tran2", IDs[1], IDs[2])

	time.Sleep(1 * time.Second)

	makeAccountsAt(&peerNode4, IDs, amounts)
	goStart(&peerNode4, peerNode4_port, peerNode2_port)

	time.Sleep(1 * time.Second)

	if peerNode4.LocalLedger.Accounts[IDs[0]] != amounts[0]-42{
		t.Error("ApplyTransaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc1' on latecomer peerNode4.",
			"\npeerNode4 Accounts:", peerNode4.LocalLedger.Accounts, "\n")
	}
	if peerNode4.LocalLedger.Accounts[IDs[1]] != (amounts[1]+42) - 5 {
		t.Error("ApplyTransaction 'tran1' and 'tran2' from 'acc1' to 'acc2' and 'acc2' to 'acc3' on peerNode1 didn't apply to 'acc2' on latecomer peerNode4.",
			"\npeerNode4 Accounts:", peerNode4.LocalLedger.Accounts, "\n")
	}
	if peerNode4.LocalLedger.Accounts[IDs[2]] != amounts[2]+5 {
		t.Error("ApplyTransaction 'tran3' from 'acc2' to 'acc3' on peerNode1 didn't apply to 'acc3' on latecomer peerNode4.",
			"\npeerNode4 Accounts:", peerNode4.LocalLedger.Accounts, "\n")
	}
	if len(peerNode4.TransactionsSeen.Vals()) != 2 {
		t.Error("Not all transactions was seen by latecomer peerNode4",
			"\npeerNode4 Accounts:", peerNode4.TransactionsSeen.Vals(), "\n")
	}
}
*/

func TestTransactionsAreBroadcastedAndAppliedOnAllPeerNodes(t *testing.T) {
	// t.Parallel()

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
	makeAccountsAt(&peerNode1, IDs, amounts)
	makeAccountsAt(&peerNode2, IDs, amounts)
	makeAccountsAt(&peerNode3, IDs, amounts)

	time.Sleep(1 * time.Second)

	peerNode1.MakeAndBroadcastTransaction(20, "tran1", IDs[0], IDs[1])

	time.Sleep(1 * time.Second)

	if peerNode1.LocalLedger.Accounts[IDs[0]] != amounts[0]-20 {
		t.Error("ApplyTransaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc1' on peerNode1.",
			"\nnpeerNode1 Accounts:", peerNode1.LocalLedger.Accounts, "\n")
	}
	if peerNode1.LocalLedger.Accounts[IDs[1]] != amounts[1]+20 {
		t.Error("ApplyTransaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc2' on peerNode1.",
			"\nnpeerNode1 Accounts:", peerNode1.LocalLedger.Accounts, "\n")
	}
	if peerNode3.LocalLedger.Accounts[IDs[0]] != amounts[0]-20 {
		t.Error("ApplyTransaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc1' on peerNode3.",
			"\nnpeerNode1 Accounts:", peerNode3.LocalLedger.Accounts, "\n")
	}
	if peerNode3.LocalLedger.Accounts[IDs[1]] != amounts[1]+20 {
		t.Error("ApplyTransaction 'tran1' from 'acc1' to 'acc2' on peerNode1 didn't apply to 'acc2' on peerNode3.",
			"\nnpeerNode1 Accounts:", peerNode3.LocalLedger.Accounts, "\n")
	}
}
func TestNodeConnectsToThreeOthersWhenEnteringNetwork(t *testing.T) {
	// t.Parallel()

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

	if len(peerNode1.OpenConnections.Values()) != 3 {
		t.Error("PeerNode1 doesn't have exactly 3 connections")
	}
	if len(peerNode2.OpenConnections.Values()) != 3 {
		t.Error("PeerNode2 doesn't have exactly 3 connections.\npeerNode2 openConnections:",
			peerNode2.OpenConnections.ToString())
	}
	if len(peerNode3.OpenConnections.Values()) != 3 {
		t.Error("PeerNode3 doesn't have exactly 3 connections")
	}
	if len(peerNode4.OpenConnections.Values()) != 3 {
		t.Error("PeerNode4 doesn't have exactly 3 connections")
	}

	if !peerNode4.PeersInArrivalOrder.Contains(peer4Port) {
		t.Error("peerNode4 (the last one connected to the network) isn't the last node in PeersInArrivalOrder:", 
			peerNode4.PeersInArrivalOrder.Values())
	}
}
func TestNodeConnectsToTenOthersWhenEnteringNetwork(t *testing.T) {
	// t.Parallel()

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
	if len(peerNode1.OpenConnections.Values()) != 10 {
		t.Error("PeerNode1 doesn't have exactly 10 connections")
	}
	if len(peerNode2.OpenConnections.Values()) != 10 {
		t.Error("PeerNode2 doesn't have exactly 10 connections")
	}
	if len(peerNode10.OpenConnections.Values()) != 10 {
		t.Error("PeerNode3 doesn't have exactly 10 connections")
	}
	if len(peerNode11.OpenConnections.Values()) != 10 {
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
	// t.Parallel()

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
	p4Conns := peerNode4.OpenConnections.Values()
	if len(p4Conns) != 3 {
		t.Error("peerNode4 doesn't have an open connection to each node in the network.\n" +
			"PeerNode4's openConnections:", p4Conns)
	}
}
func TestReceivedPeerListWhenJoining(t *testing.T) {
	// t.Parallel()

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
	// t.Parallel()

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
	// t.Parallel()

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
	// t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true, true)
	peerNode2 := createPeerNode(true, true)

	goStart(&peerNode1, peerNode1_port, "no_port")
	goStart(&peerNode2, peerNode2_port, peerNode1_port)

	time.Sleep(4 * time.Second)

	simulateInputFor(&peerNode1, "some_msg")

	if !peerNode2.MessagesSent.Contains("some_msg") {
		t.Errorf("peerNode2 didn't receive peerNode1's msg")
	}

}
func TestPeer1CanConnectToPeer2(t *testing.T) {
	// t.Parallel()


	peerNode1 := createPeerNode(true, true)
	peerNode2 := createPeerNode(true, true)

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()

	goStart(&peerNode1, peerNode1_port, "no_port")
	goStart(&peerNode2, peerNode2_port, peerNode1_port)

	if len(peerNode1.OpenConnections.Values()) != 1 {
		t.Error("peerNode1 doesn't have peerNode2 in its openConnections set:", peerNode1.OpenConnections.Values())
	}
	if len(peerNode2.OpenConnections.Values()) != 1 {
		t.Error("peerNode2 doesn't have peerNode1 in its openConnections set:", peerNode2.OpenConnections.Values())
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

func makeAccountsAt(peerNode *PeerNode, IDs []string, startAmounts []int){
	for i, id := range IDs {
		peerNode.CreateAccountInLedger(id, startAmounts[i])
	}
}

