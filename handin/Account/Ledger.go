package Account

import ( "sync" )

type Ledger struct {
	Accounts map[string]int
	lock sync.Mutex
}
func MakeLedger() *Ledger {
	ledger := new(Ledger)
	ledger.Accounts = make(map[string]int)
	return ledger
}
func (l *Ledger) CreateAccount(id string, initialAmount int) {
	l.lock.Lock() ; defer l.lock.Unlock()
	l.Accounts[id] = initialAmount
}
