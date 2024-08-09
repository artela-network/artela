package v0

import (
	"math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artela-network/artela/x/aspect/store"
	evmtypes "github.com/artela-network/artela/x/evm/types"
)

var emptyAddress common.Address

const (
	maxAspectBoundLimit           = math.MaxUint8
	maxContractVerifierBoundLimit = uint8(1)
	protocolVersion               = store.ProtocolVersion(0)
)

// BaseStore defines a shared base store which can be implemented by all other stores
type BaseStore interface {
	NewPrefixStore(prefixKey string) prefix.Store
	Load(prefixStore prefix.Store, key []byte) ([]byte, error)
	Store(prefixStore prefix.Store, key, value []byte) error
	Version() store.ProtocolVersion

	store.GasMeteredStore
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

// Version returns the protocol version of the Store.
func (s *baseStore) Version() store.ProtocolVersion {
	return protocolVersion
}

// NewPrefixStore creates an instance of prefix Store,
func (s *baseStore) NewPrefixStore(prefixKey string) prefix.Store {
	return prefix.NewStore(s.kvStore, evmtypes.KeyPrefix(prefixKey))
}

// Load loads the value from the given Store and do gas metering for the given operation
func (s *baseStore) Load(prefixStore prefix.Store, key []byte) ([]byte, error) {
	if key == nil {
		return nil, store.ErrInvalidStorageKey
	}

	value := prefixStore.Get(key)

	// gas metering after Load, since we are not like EVM, the data length is not known before Load
	if err := s.gasMeter.MeasureStorageLoad(len(key) + len(value)); err != nil {
		return nil, err
	}

	return value, nil
}

// Store stores the value to the given store and do gas metering for the given operation
func (s *baseStore) Store(prefixStore prefix.Store, key, value []byte) error {
	if key == nil {
		return store.ErrInvalidStorageKey
	}

	if err := s.gasMeter.MeasureStorageStore(len(key) + len(value)); err != nil {
		return err
	}

	prefixStore.Set(key, value)
	return nil
}

func (s *baseStore) TransferGasFrom(store store.GasMeteredStore) {
	s.gasMeter.UpdateGas(store.Gas())
}

func (s *baseStore) Gas() uint64 {
	return s.gasMeter.RemainingGas()
}
