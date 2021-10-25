package Account

import (
	. "DNO/handin/Cryptography"
	"math/big"
	"strconv"
	"strings"
)

type SignedTransaction struct {
	ID string // Any string
	From string // A verification key coded as a string
	To string // A verification key coded as a string
	Amount int // Amount to transfer
	Signature string // Potential signature coded as string
}

func (s *SignedTransaction) ToString() string {
	return "\n" +
			"ID:   " + s.ID[:8] + "\n" +
			"From: " + s.From[:8] + "\n" +
			"To:   " + s.To[:8] + "\n" +
			"Amount: " + strconv.Itoa(s.Amount) + "\n" +
			"Signature: " + s.Signature[:8]
}

func (l *Ledger) ApplySignedTransaction(t SignedTransaction) {
	l.lock.Lock() ; defer l.lock.Unlock()
	/* We verify that the t.Signature is a valid RSA
	 * signature on the rest of the fields in t under
	 * the public key t.From.
	 */
	validSignature := l.Verify(t)

	if validSignature {
		l.Accounts[t.From] -= t.Amount
		l.Accounts[t.To] += t.Amount
	}
}

func (l *Ledger) MakeSignedTransaction(t Transaction, sk SecretKey) SignedTransaction {
	return SignedTransaction {
		ID: t.ID,           // We assume that neither ID, From, To, nor Amount is allowed to contain a colon :
		From: t.From,
		To: t.To,
		Amount: t.Amount,
		Signature: CreateSignature(t.ID + ":" + t.From + ":" + t.To + ":" + strconv.Itoa(t.Amount), sk),
	}
}


func (l *Ledger) Verify(st SignedTransaction) bool {
	msg := st.ID + ":" + st.From + ":" + st.To + ":" + strconv.Itoa(st.Amount)
	pk := extractPublicKeyFrom(st)
	return Verify(st.Signature, msg, pk)
}

func extractPublicKeyFrom(st SignedTransaction) PublicKey {
	y := strings.Split(st.To, ":")

	i, _ := strconv.Atoi(y[0])
	j, _ := strconv.Atoi(y[0])
	n := big.NewInt(int64(i))
	e := big.NewInt(int64(j))

	return PublicKey{ N: n, E: e}
}