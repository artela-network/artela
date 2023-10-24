package states

// Derived from https://github.com/ethereum/go-ethereum/blob/v1.12.0/core/types/state_account.go

import (
	"bytes"
	"math/big"
)

// StateAccount is the Ethereum consensus representation of accounts.
// These objects are stored in the storage of auth module.
type StateAccount struct {
	Nonce    uint64
	Balance  *big.Int
	CodeHash []byte
}

// NewEmptyAccount returns an empty account.
func NewEmptyAccount() *StateAccount {
	return &StateAccount{
		Balance:  new(big.Int),
		CodeHash: emptyCodeHash,
	}
}

// IsContract returns if the account contains contract code.
func (acct StateAccount) IsContract() bool {
	return !bytes.Equal(acct.CodeHash, emptyCodeHash)
}
