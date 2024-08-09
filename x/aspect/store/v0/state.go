package v0

import (
	"github.com/artela-network/artela/x/aspect/store"
	"github.com/artela-network/artela/x/aspect/types"
)

var _ store.AspectStateStore = (*stateStore)(nil)

// stateStore is the version 0 state Store, this is no longer maintained.
// Deprecated.
type stateStore struct {
	BaseStore

	ctx *types.AspectStoreContext
}

// NewStateStore creates a new instance of account state.
// Deprecated
func NewStateStore(ctx *types.AspectStoreContext) store.AspectStateStore {
	// for state Store, we have already charged gas in host api,
	// so no need to charge it again in the Store
	return &stateStore{
		BaseStore: NewBaseStore(NewNoOpGasMeter(ctx), ctx.CosmosContext().KVStore(ctx.EVMStoreKey())),
		ctx:       ctx,
	}
}

// SetState sets the state of the aspect with the given ID and key.
func (s *stateStore) SetState(key []byte, value []byte) {
	aspectID := s.ctx.AspectID
	prefixStore := s.NewPrefixStore(V0AspectStateKeyPrefix)
	storeKey := AspectArrayKey(aspectID.Bytes(), key)
	if len(value) == 0 {
		prefixStore.Delete(storeKey)
	}
	prefixStore.Set(storeKey, value)
}

// GetState returns the state of the aspect with the given ID and key.
func (s *stateStore) GetState(key []byte) []byte {
	aspectID := s.ctx.AspectID
	prefixStore := s.NewPrefixStore(V0AspectStateKeyPrefix)
	storeKey := AspectArrayKey(aspectID.Bytes(), key)
	return prefixStore.Get(storeKey)
}
