package Account

type Transaction struct {
	ID string
	From string
	To string
	Amount int
}
func (l *Ledger) ApplyTransaction(t *Transaction) {
	l.lock.Lock() ; defer l.lock.Unlock()
	l.Accounts[t.From] -= t.Amount
	l.Accounts[t.To] += t.Amount
}
