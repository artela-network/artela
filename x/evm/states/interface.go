package states

import (
	cosmos "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artela-network/artela-evm/vm"
)

// ExtStateDB defines an extension to the interface provided by the go-ethereum
// codebase to support additional states transition functionalities. In particular
// it supports appending a new entry to the states journal through
// AppendJournalEntry so that the states can be reverted after running
// stateful precompiled contracts.
type ExtStateDB interface {
	vm.StateDB
	AppendJournalEntry(JournalEntry)
}

// Keeper provide underlying storage of StateDB
type Keeper interface {
	// Read methods
	GetAccount(ctx cosmos.Context, addr common.Address) *StateAccount
	GetState(ctx cosmos.Context, addr common.Address, key common.Hash) common.Hash
	GetCode(ctx cosmos.Context, codeHash common.Hash) []byte
	// the callback returns false to break early
	ForEachStorage(ctx cosmos.Context, addr common.Address, cb func(key, value common.Hash) bool)

	// Write methods, only called by `StateDB.Commit()`
	SetAccount(ctx cosmos.Context, addr common.Address, account StateAccount) error
	SetState(ctx cosmos.Context, addr common.Address, key common.Hash, value []byte)
	SetCode(ctx cosmos.Context, codeHash []byte, code []byte)
	DeleteAccount(ctx cosmos.Context, addr common.Address) error
}
