package Helper

import (
	. "DNO/handin/Cryptography"
	"math/big"
	"strconv"
	"strings"
)

type Sequencer struct {
	UnsequencedTransactionIDs SafeArray_string
	PublicKey                 PublicKey
	KeyPair    				  KeyPair
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
	return b.VerificationKey.ToString() + ":" +
		   strconv.Itoa(b.SlotNumber) + ":" +
		   b.Draw.String() + ":" +
		   strings.Join(b.TransactionIDs, ",") + ":" +
		   string(b.PrevBlockHash) + ":" +
		   strconv.Itoa(b.LengthToRoot)
}

func (s *Sequencer) Sign(block Block) SignedBlock {
	blockString := block.ToString()
	signature := CreateSignature(blockString, s.KeyPair.Sk)
	return SignedBlock {Block: block, Signature: signature}
}

func (b *Block) Hash() string {
	n := new(big.Int)
	n.SetString(b.ToString(), 10)
	return string(Hash(n))
}

func (s *Sequencer) Verify(signedBlock SignedBlock) bool {
	return Verify(signedBlock.Signature, signedBlock.Block.ToString(), s.PublicKey)
}
