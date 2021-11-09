package Account

import (
	. "DNO/handin/Cryptography"
	"fmt"
	"strconv"
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
			"ID:   " + s.ID + "\n" +
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
	if !validSignature  {
		fmt.Printf("Signature is invalid! Transaction from (%s...) to (%s...).\n", t.From[:8], t.To[:8])
		return
	}
	amountPositive := t.Amount > 0
	if !amountPositive {
		fmt.Printf("The amount %d on the signed transaction is negative\n", t.Amount)
		return
	}
	willAmountInFromAccBecomeNegative := l.Accounts[t.From] - t.Amount < 0
	if willAmountInFromAccBecomeNegative {
		//fmt.Printf("Will become less than 0\n", t.Amount)
		return
	}
	l.Accounts[t.From] -= t.Amount
	l.Accounts[t.To] += t.Amount
	//fmt.Println("Account t.From (", t.From[:4] ,") is now:", l.Accounts[t.From])
	//fmt.Println("Account t.To (", t.To[:4] ,") is now:", l.Accounts[t.To])
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
	pk := ToPublicKey(st.From)
	return Verify(st.Signature, msg, pk)
}
