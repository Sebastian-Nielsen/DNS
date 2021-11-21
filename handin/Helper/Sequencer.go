package Helper

import (
	. "DNO/handin/Cryptography"
	"math/big"
	"strconv"
	"strings"
)

type Sequencer struct {
	UnsequencedTransactionIDs SafeArray_string
	//PublicKey                 PublicKey
	IsSequencer				  bool
	//KeyPair    				  KeyPair
	SlotNumber 				  SafeCounter
	Seed       				  string
	Hardness				  *big.Int
	Tree 				      Tree
}

type Block struct {
	VerificationKey 		PublicKey
	SlotNumber      		int
	Draw            		*big.Int
	TransactionIDs          []string
	PrevBlockHash			string
	LengthToRoot            int
	HasBeenApplied			bool
}



type GenesisBlock struct {
	Seed				   string
	InitialAccounts        []PublicKey
	InitialAmount		   int
	Hardness			   *big.Int
}

type SignedBlock struct {
	Block     Block
	Signature string
}

func (b *Block) ToString() string {
	return strconv.Itoa(b.SlotNumber) + ":" +
	       b.VerificationKey.ToString() + ":" +
		   b.Draw.String() + ":" +
		   strings.Join(b.TransactionIDs, ",") + ":" +
		   b.PrevBlockHash + ":" +
		   strconv.Itoa(b.LengthToRoot)
}

func (s *Sequencer) Sign(block Block, sk SecretKey) SignedBlock {
	blockString := block.ToString()
	signature := CreateSignature(blockString, sk)
	return SignedBlock {Block: block, Signature: signature}
}

func (b *Block) Hash() string {
	n := new(big.Int)
	n.SetBytes([]byte(b.ToString()))
	return string(Hash(n))
}

func (s *Sequencer) Verify(signedBlock SignedBlock) bool {
	return Verify(signedBlock.Signature, signedBlock.Block.ToString(), signedBlock.Block.VerificationKey)
}
