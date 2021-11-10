package Helper

import (
	. "DNO/handin/Cryptography"
	"strconv"
	"strings"
)

type Sequencer struct {
	UnsequencedTransactionIDs SafeArray_string
	PublicKey                 PublicKey
	KeyPair                   KeyPair
	BlockNumber               SafeCounter
}

type Block struct {
	BlockNumber             int
	TransactionIDs          []string
}



type SignedBlock struct {
	Block     Block
	Signature string
}

func (b *Block) ToString() string {
	return strconv.Itoa(b.BlockNumber) + ":" + strings.Join(b.TransactionIDs, ",")
}

func (s *Sequencer) Sign(block Block) SignedBlock {
	blockString := block.ToString()
	signature := CreateSignature(blockString, s.KeyPair.Sk)
	return SignedBlock {Block: block, Signature: signature}
}

func (s *Sequencer) Verify(signedBlock SignedBlock) bool {
	return Verify(signedBlock.Signature, signedBlock.Block.ToString(), s.PublicKey)
}
