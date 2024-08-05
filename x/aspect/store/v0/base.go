package v0

import (
	"github.com/artela-network/artela/x/aspect/store"
	evmtypes "github.com/artela-network/artela/x/evm/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"math"
)

var emptyAddress common.Address

const (
	maxAspectBoundLimit           = math.MaxUint8
	maxContractVerifierBoundLimit = uint8(1)
	protocolVersion               = store.ProtocolVersion(1)
)

type BaseStore interface {
	newPrefixStore(prefixKey string) prefix.Store
	load(prefixStore prefix.Store, key []byte) ([]byte, error)
	store(prefixStore prefix.Store, key, value []byte) error
	Version() store.ProtocolVersion
}

type baseStore struct {
	gasMeter GasMeter
	kvStore  sdk.KVStore
}

func NewBaseStore(gasMeter GasMeter, kvStore sdk.KVStore) BaseStore {
	return &baseStore{
		gasMeter: gasMeter,
		kvStore:  kvStore,
	}
}

// Version returns the protocol version of the store.
func (s *baseStore) Version() store.ProtocolVersion {
	return protocolVersion
}

// newPrefixStore creates an instance of prefix store,
func (s *baseStore) newPrefixStore(prefixKey string) prefix.Store {
	return prefix.NewStore(s.kvStore, evmtypes.KeyPrefix(prefixKey))
}

// load loads the value from the given store and do gas metering for the given operation
func (s *baseStore) load(prefixStore prefix.Store, key []byte) ([]byte, error) {
	if key == nil {
		return nil, store.ErrInvalidStorageKey
	}

	value := prefixStore.Get(key)

	// gas metering after load, since we are not like EVM, the data length is not known before load
	if err := s.gasMeter.measureStorageLoad(len(key) + len(value)); err != nil {
		return nil, err
	}

	return value, nil
}

// store stores the value to the given store and do gas metering for the given operation
func (s *baseStore) store(prefixStore prefix.Store, key, value []byte) error {
	if key == nil {
		return store.ErrInvalidStorageKey
	}

	if err := s.gasMeter.measureStorageStore(len(key) + len(value)); err != nil {
		return err
	}

	prefixStore.Set(key, value)
	return nil
}
