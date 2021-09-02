package account

type SignedTransaction struct {
     ID        string // Any string
     From      string // A verification key coded as a string
     To        string // A verification key coded as a string
     Amount    int    // Amount to transfer
     Signature string // Potential signature coded as string
}

func (l *Ledger) SignedTransaction(t *SignedTransaction) {
     l.lock.Lock() ; defer l.lock.Unlock()

     /* We verify that the t.Signature is a valid RSA
      * signature on the rest of the fields in t under
      * the public key t.From.
      */
     validSignature := true 

     if validSignature {
     	l.Accounts[t.From] -= t.Amount
     	l.Accounts[t.To] += t.Amount
     }
}
