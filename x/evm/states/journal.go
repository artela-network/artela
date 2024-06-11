package states

// Derived from https://github.com/ethereum/go-ethereum/blob/v1.12.0/core/state/journal.go

import (
	"bytes"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
)

// JournalEntry is a modification entry in the states change journal that can be
// Reverted on demand.
type JournalEntry interface {
	// Revert undoes the changes introduced by this journal entry.
	Revert(*StateDB)

	// Dirtied returns the Ethereum address modified by this journal entry.
	Dirtied() *common.Address
}

// ----------------------------------------------------------------------------
// 								   journal
// ----------------------------------------------------------------------------

// journal contains the list of states modifications applied since the last states
// commit. These are tracked to be able to be reverted in the case of an execution
// exception or request for reversal.
type journal struct {
	entries []JournalEntry         // Current changes tracked by the journal
	dirties map[common.Address]int // Dirty accounts and the number of changes
}

// newJournal creates a new initialized journal.
func newJournal() *journal {
	return &journal{
		dirties: make(map[common.Address]int),
	}
}

// append inserts a new modification entry to the end of the change journal.
func (j *journal) append(entry JournalEntry) {
	j.entries = append(j.entries, entry)
	if addr := entry.Dirtied(); addr != nil {
		j.dirties[*addr]++
	}
}

// Revert undoes a batch of journalled modifications along with any Reverted
// dirty handling too.
func (j *journal) Revert(statedb *StateDB, snapshot int) {
	for i := len(j.entries) - 1; i >= snapshot; i-- {
		// Undo the changes made by the operation
		j.entries[i].Revert(statedb)

		// Drop any dirty tracking induced by the change
		if addr := j.entries[i].Dirtied(); addr != nil {
			if j.dirties[*addr]--; j.dirties[*addr] == 0 {
				delete(j.dirties, *addr)
			}
		}
	}
	j.entries = j.entries[:snapshot]
}

// length returns the current number of entries in the journal.
func (j *journal) length() int {
	return len(j.entries)
}

// sortedDirties sort the dirty addresses for deterministic iteration
func (j *journal) sortedDirties() []common.Address {
	keys := make([]common.Address, 0, len(j.dirties))
	for k := range j.dirties {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i].Bytes(), keys[j].Bytes()) < 0
	})
	return keys
}

type (
	// Changes to the account trie.
	createObjectChange struct {
		account *common.Address
	}
	resetObjectChange struct {
		prev *stateObject
	}
	suicideChange struct {
		account     *common.Address
		prev        bool // whether account had already suicided
		prevbalance *big.Int
	}

	// Changes to individual accounts.
	balanceChange struct {
		account *common.Address
		prev    *big.Int
	}
	nonceChange struct {
		account *common.Address
		prev    uint64
	}
	storageChange struct {
		account       *common.Address
		key, prevalue common.Hash
	}
	codeChange struct {
		account            *common.Address
		prevcode, prevhash []byte
	}

	// Changes to other states values.
	refundChange struct {
		prev uint64
	}
	addLogChange struct{}

	// Changes to the access list
	accessListAddAccountChange struct {
		address *common.Address
	}
	accessListAddSlotChange struct {
		address *common.Address
		slot    *common.Hash
	}

	transientStorageChange struct {
		account       *common.Address
		key, prevalue common.Hash
	}
)

// ----------------------------------------------------------------------------
// 								createObjectChange
// ----------------------------------------------------------------------------

func (ch createObjectChange) Revert(s *StateDB) {
	delete(s.stateObjects, *ch.account)
}

func (ch createObjectChange) Dirtied() *common.Address {
	return ch.account
}

// ----------------------------------------------------------------------------
// 								resetObjectChange
// ----------------------------------------------------------------------------

func (ch resetObjectChange) Revert(s *StateDB) {
	s.setStateObject(ch.prev)
}

func (ch resetObjectChange) Dirtied() *common.Address {
	return nil
}

// ----------------------------------------------------------------------------
// 								 suicideChange
// ----------------------------------------------------------------------------

func (ch suicideChange) Revert(s *StateDB) {
	obj := s.getStateObject(*ch.account)
	if obj != nil {
		obj.suicided = ch.prev
		obj.setBalance(ch.prevbalance)
	}
}

func (ch suicideChange) Dirtied() *common.Address {
	return ch.account
}

// ----------------------------------------------------------------------------
// 								 balanceChange
// ----------------------------------------------------------------------------

func (ch balanceChange) Revert(s *StateDB) {
	s.getStateObject(*ch.account).setBalance(ch.prev)
}

func (ch balanceChange) Dirtied() *common.Address {
	return ch.account
}

// ----------------------------------------------------------------------------
// 								 nonceChange
// ----------------------------------------------------------------------------

func (ch nonceChange) Revert(s *StateDB) {
	s.getStateObject(*ch.account).setNonce(ch.prev)
}

func (ch nonceChange) Dirtied() *common.Address {
	return ch.account
}

// ----------------------------------------------------------------------------
// 								 codeChange
// ----------------------------------------------------------------------------

func (ch codeChange) Revert(s *StateDB) {
	s.getStateObject(*ch.account).setCode(common.BytesToHash(ch.prevhash), ch.prevcode)
}

func (ch codeChange) Dirtied() *common.Address {
	return ch.account
}

// ----------------------------------------------------------------------------
// 								storageChange
// ----------------------------------------------------------------------------

func (ch storageChange) Revert(s *StateDB) {
	s.getStateObject(*ch.account).setState(ch.key, ch.prevalue)
}

func (ch storageChange) Dirtied() *common.Address {
	return ch.account
}

// ----------------------------------------------------------------------------
// 								refundChange
// ----------------------------------------------------------------------------

func (ch refundChange) Revert(s *StateDB) {
	s.refund = ch.prev
}

func (ch refundChange) Dirtied() *common.Address {
	return nil
}

// ----------------------------------------------------------------------------
// 								addLogChange
// ----------------------------------------------------------------------------

func (ch addLogChange) Revert(s *StateDB) {
	s.logs = s.logs[:len(s.logs)-1]
}

func (ch addLogChange) Dirtied() *common.Address {
	return nil
}

// ----------------------------------------------------------------------------
// 						  accessListAddAccountChange
// ----------------------------------------------------------------------------

func (ch accessListAddAccountChange) Revert(s *StateDB) {
	/*
		One important invariant here, is that whenever a (addr, slot) is added, if the
		addr is not already present, the add causes two journal entries:
		- one for the address,
		- one for the (address,slot)
		Therefore, when unrolling the change, we can always blindly delete the
		(addr) at this point, since no storage adds can remain when come upon
		a single (addr) change.
	*/
	s.accessList.DeleteAddress(*ch.address)
}

func (ch accessListAddAccountChange) Dirtied() *common.Address {
	return nil
}

// ----------------------------------------------------------------------------
// 					       accessListAddSlotChange
// ----------------------------------------------------------------------------

func (ch accessListAddSlotChange) Revert(s *StateDB) {
	s.accessList.DeleteSlot(*ch.address, *ch.slot)
}

func (ch accessListAddSlotChange) Dirtied() *common.Address {
	return nil
}

// ----------------------------------------------------------------------------
// 					       transientStorageChange
// ----------------------------------------------------------------------------

func (ch transientStorageChange) Revert(s *StateDB) {
	s.SetTransientState(*ch.account, ch.key, ch.prevalue)
}

func (ch transientStorageChange) Dirtied() *common.Address {
	return nil
}
