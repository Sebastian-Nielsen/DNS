package Helper

import (
	"DNO/handin/Account"
	"strings"
	"sync"
)

type SafeArray_string struct {
	values []string
	mu     sync.Mutex
}
func (a *SafeArray_string) Get(index int) string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.values[index]
}
func (a *SafeArray_string) PopHead() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	 head, tail := a.values[0], a.values[1:]
	 a.values = tail
	return head
}
func (a *SafeArray_string) Set(array []string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.values = array
}
func (a *SafeArray_string) ToString() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return strings.Join(a.values[:], ",")
}
func (a *SafeArray_string) Append(value string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.values = append(a.values, value)
}
func (a *SafeArray_string) PopAll() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	cpy := make([]string, len(a.values))
	copy(cpy, a.values)
	a.values = []string{}
	return cpy
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
func (a *SafeArray_string) IsEmpty() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.values) == 0
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



type SafeArray_Block struct {
	values []Block
	mu     sync.Mutex
}
func (a *SafeArray_Block) Append(value Block) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.values = append(a.values, value)
}
func (a *SafeArray_Block) Values() []Block {
	a.mu.Lock()
	defer a.mu.Unlock()
	cpy := make([]Block, len(a.values))
	copy(cpy, a.values)
	return cpy
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
func (a *SafeArray_SignedTransaction) Set(array []Account.SignedTransaction) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.values = array
}
func (a *SafeArray_SignedTransaction) Contains(transactionToSearchFor Account.SignedTransaction) bool {
	for _, val := range a.values {
		if val == transactionToSearchFor {
			return true
		}
	}
	return false
}
