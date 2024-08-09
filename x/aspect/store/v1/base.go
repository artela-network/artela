package v1

import (
	"encoding/hex"
	"github.com/artela-network/artela/x/aspect/store"
	v0 "github.com/artela-network/artela/x/aspect/store/v0"
	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math"
)

const (
	maxAspectBoundLimit = math.MaxUint8
	protocolVersion     = store.ProtocolVersion(1)
)

// BaseStore defines a shared base store which can be implemented by all other stores
type BaseStore interface {
	Load(key []byte) ([]byte, error)
	Store(key, value []byte) error
	Version() store.ProtocolVersion

	store.GasMeteredStore
}

type baseStore struct {
	gasMeter v0.GasMeter
	kvStore  sdk.KVStore
	logger   log.Logger
}

func NewBaseStore(logger log.Logger, gasMeter v0.GasMeter, kvStore sdk.KVStore) BaseStore {
	return &baseStore{
		gasMeter: gasMeter,
		kvStore:  kvStore,
		logger:   logger,
	}
}

// Version returns the protocol version of the Store.
func (s *baseStore) Version() store.ProtocolVersion {
	return protocolVersion
}

// Load loads the value from the given Store and do gas metering for the given operation
func (s *baseStore) Load(key []byte) ([]byte, error) {
	if key == nil {
		return nil, store.ErrInvalidStorageKey
	}

	value := s.kvStore.Get(key)

	s.logger.Debug("load from aspect store", "key", abbreviateHex(key), "data", abbreviateHex(value))
	// gas metering after Load, since we are not like EVM, the data length is not known before Load
	if err := s.gasMeter.MeasureStorageLoad(len(key) + len(value)); err != nil {
		return nil, err
	}

	// copy values to avoid unexpected modification
	result := make([]byte, len(value))
	copy(result, value)

	return result, nil
}

// Store stores the value to the given store and do gas metering for the given operation
func (s *baseStore) Store(key, value []byte) error {
	if key == nil {
		return store.ErrInvalidStorageKey
	}

	if len(value) == 0 {
		// if value is nil, we just delete the key, this will not charge gas
		s.logger.Debug("deleting from aspect store", "key", abbreviateHex(key))
		s.kvStore.Delete(key)
	} else {
		if err := s.gasMeter.MeasureStorageStore(len(key) + len(value)); err != nil {
			return err
		}

		s.logger.Debug("saving to aspect store", "key", abbreviateHex(key), "data", abbreviateHex(value))
		s.kvStore.Set(key, value)
	}
	return nil
}

func (s *baseStore) TransferGasFrom(store store.GasMeteredStore) {
	s.gasMeter.UpdateGas(store.Gas())
}

func (s *baseStore) Gas() uint64 {
	return s.gasMeter.RemainingGas()
}

func abbreviateHex(data []byte) string {
	if len(data) > 100 {
		return hex.EncodeToString(data[:100]) + "..."
	}
	return hex.EncodeToString(data)
}
