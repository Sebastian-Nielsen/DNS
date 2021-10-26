package Helper

import (
	"DNO/handin/Account"
	"sync"
)

type SafeArray_string struct {
	values []string
	mu     sync.Mutex
}
func (a *SafeArray_string) Append(value string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.values = append(a.values, value)
}
func (a *SafeArray_string) Values() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	cpy := make([]string, len(a.values))
	copy(cpy, a.values)
	return cpy
}
func (a *SafeArray_string) Contains(valueToSearchFor string) bool {
	for _, val := range a.values {
		if val == valueToSearchFor {
			return true
		}
	}
	return false
}


type SafeArray_Transaction struct {
	values []Account.Transaction
	mu     sync.Mutex
}
func (a *SafeArray_Transaction) Append(value Account.Transaction) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.values = append(a.values, value)
}
func (a *SafeArray_Transaction) Values() []Account.Transaction {
	a.mu.Lock()
	defer a.mu.Unlock()
	cpy := make([]Account.Transaction, len(a.values))
	copy(cpy, a.values)
	return cpy
}
func (a *SafeArray_Transaction) Contains(transactionToSearchFor Account.Transaction) bool {
	for _, val := range a.values {
		if val == transactionToSearchFor {
			return true
		}
	}
	return false
}


type SafeArray_SignedTransaction struct {
	values []Account.SignedTransaction
	mu     sync.Mutex
}
func (a *SafeArray_SignedTransaction) Append(value Account.SignedTransaction) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.values = append(a.values, value)
}
func (a *SafeArray_SignedTransaction) Values() []Account.SignedTransaction {
	a.mu.Lock()
	defer a.mu.Unlock()
	cpy := make([]Account.SignedTransaction, len(a.values))
	copy(cpy, a.values)
	return cpy
}
func (a *SafeArray_SignedTransaction) Contains(transactionToSearchFor Account.SignedTransaction) bool {
	for _, val := range a.values {
		if val == transactionToSearchFor {
			return true
		}
	}
	return false
}
