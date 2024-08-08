package v0

import (
	"encoding/json"
	"github.com/artela-network/artela/x/aspect/store"
	"github.com/artela-network/artela/x/aspect/types"
	evmtypes "github.com/artela-network/artela/x/evm/artela/types"
	artelasdkType "github.com/artela-network/aspect-core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"slices"
	"sort"
)

var _ store.AccountStore = (*accountStore)(nil)

// accountStore is the version 0 account Store, this is no longer maintained.
// Deprecated.
type accountStore struct {
	BaseStore

	ctx *types.AccountStoreContext
}

// NewAccountStore creates a new instance of account Store.
// Deprecated
func NewAccountStore(ctx *types.AccountStoreContext) store.AccountStore {
	var meter GasMeter
	if ctx.ChargeGas() {
		meter = NewGasMeter(ctx)
	} else {
		meter = NewNoOpGasMeter(ctx)
	}

	return &accountStore{
		BaseStore: NewBaseStore(meter, ctx.CosmosContext().KVStore(ctx.EVMStoreKey())),
		ctx:       ctx,
	}
}

func (s *accountStore) Used() (bool, error) {
	// check all binding keys, if any of them is not empty, then the account Store is used
	bindingKeys := []string{V0VerifierBindingKeyPrefix, V0ContractBindKeyPrefix}
	account := s.ctx.Account
	for _, bindingKey := range bindingKeys {
		prefixStore := s.NewPrefixStore(bindingKey)
		storeKey := AccountKey(account.Bytes())
		rawJSON, err := s.Load(prefixStore, storeKey)
		if err != nil {
			return false, err
		}
		if len(rawJSON) > 0 {
			return true, nil
		}
	}

	return false, nil
}

func (s *accountStore) MigrateFrom(old store.AccountStore) error {
	panic("cannot migrate to v0 Store")
}

func (s *accountStore) Init() error {
	return nil
}

func (s *accountStore) getBindingKeyAndLimit(joinPoint uint64, isCA bool) ([]struct {
	key   string
	limit uint8
}, error) {
	joinPointI64 := int64(joinPoint)
	isTxLevel := artelasdkType.CheckIsTransactionLevel(joinPointI64)
	isVerifier := artelasdkType.CheckIsTxVerifier(joinPointI64)

	// for EoA account we can only bind verifier aspect
	if !isCA && !isVerifier {
		return nil, store.ErrInvalidBinding
	}

	// only allow 1 verifier for each contract
	verifierLimit := maxContractVerifierBoundLimit
	if isCA {
		verifierLimit = 1
	}

	bindingKeysAndLimit := make([]struct {
		key   string
		limit uint8
	}, 0)
	if isTxLevel {
		bindingKeysAndLimit = append(bindingKeysAndLimit, struct {
			key   string
			limit uint8
		}{key: V0ContractBindKeyPrefix, limit: maxAspectBoundLimit})
	}
	if isVerifier {
		bindingKeysAndLimit = append(bindingKeysAndLimit, struct {
			key   string
			limit uint8
		}{key: V0VerifierBindingKeyPrefix, limit: verifierLimit})
	}

	return bindingKeysAndLimit, nil
}

// StoreBinding stores the binding of the aspect with the given ID to the account.
func (s *accountStore) StoreBinding(aspectID common.Address, version uint64, joinPoint uint64, priority int8, isCA bool) error {
	bindingKeysAndLimit, err := s.getBindingKeyAndLimit(joinPoint, isCA)
	if err != nil {
		return err
	}

	account := s.ctx.Account
	for _, bindingKeyAndLimit := range bindingKeysAndLimit {
		prefixStore := s.NewPrefixStore(bindingKeyAndLimit.key)
		storeKey := AccountKey(account.Bytes())
		rawJSON, err := s.Load(prefixStore, storeKey)
		if err != nil {
			return err
		}

		bindings := make([]*evmtypes.AspectMeta, 0)
		if len(rawJSON) > 0 {
			if err := json.Unmarshal(rawJSON, &bindings); err != nil {
				return store.ErrStorageCorrupted
			}
		}

		if len(bindings) >= int(bindingKeyAndLimit.limit) {
			return store.ErrBindingLimitExceeded
		}

		// check duplicates
		for _, binding := range bindings {
			if binding.Id == aspectID {
				return store.ErrAlreadyBound
			}
		}

		newAspect := &evmtypes.AspectMeta{
			Id:       aspectID,
			Version:  uint256.NewInt(version),
			Priority: int64(priority),
		}

		bindings = append(bindings, newAspect)

		// re-sort aspects by priority
		if len(bindings) > 1 {
			sort.Slice(bindings, evmtypes.NewBindingPriorityComparator(bindings))
		}

		bindingJSON, err := json.Marshal(bindings)
		if err != nil {
			return err
		}

		if err := s.Store(prefixStore, storeKey, bindingJSON); err != nil {
			return err
		}

		s.ctx.Logger().Debug("aspect binding added",
			"aspect", aspectID.Hex(),
			"account", account.Hex(),
			"storekey", bindingKeyAndLimit.key)
	}

	return nil
}

// RemoveBinding removes the binding of the aspect with the given ID from the account.
func (s *accountStore) RemoveBinding(aspectID common.Address, joinPoint uint64, isCA bool) error {
	bindingKeysAndLimit, err := s.getBindingKeyAndLimit(joinPoint, isCA)
	if err != nil {
		return err
	}

	account := s.ctx.Account
	for _, bindingKeyAndLimit := range bindingKeysAndLimit {
		prefixStore := s.NewPrefixStore(bindingKeyAndLimit.key)
		storeKey := AccountKey(account.Bytes())
		rawJSON, err := s.Load(prefixStore, storeKey)
		if err != nil {
			return err
		}

		bindings := make([]*evmtypes.AspectMeta, 0)
		if len(rawJSON) > 0 {
			if err := json.Unmarshal(rawJSON, &bindings); err != nil {
				return store.ErrStorageCorrupted
			}
		}

		if len(bindings) >= int(bindingKeyAndLimit.limit) {
			return store.ErrBindingLimitExceeded
		}

		toDelete := slices.IndexFunc(bindings, func(meta *evmtypes.AspectMeta) bool {
			return meta.Id == aspectID
		})

		if toDelete < 0 {
			// not found
			return nil
		}

		bindings = slices.Delete(bindings, toDelete, toDelete+1)
		bindingJSON, err := json.Marshal(bindings)
		if err != nil {
			return err
		}

		if err := s.Store(prefixStore, storeKey, bindingJSON); err != nil {
			return err
		}

		s.ctx.Logger().Debug("aspect binding removed",
			"aspect", aspectID.Hex(),
			"account", account.Hex(),
			"storekey", bindingKeyAndLimit.key)
	}

	return nil
}

// LoadAccountBoundAspects loads all aspects bound to the given account.
func (s *accountStore) LoadAccountBoundAspects(filter types.BindingFilter) ([]types.Binding, error) {
	bindingKeys := make([]string, 0, 1)
	if filter.JoinPoint != nil {
		joinPoint := *filter.JoinPoint
		if string(joinPoint) == artelasdkType.JoinPointRunType_VerifyTx.String() {
			bindingKeys = append(bindingKeys, V0VerifierBindingKeyPrefix)
		} else {
			bindingKeys = append(bindingKeys, V0ContractBindKeyPrefix)
		}
	} else if filter.VerifierOnly {
		bindingKeys = append(bindingKeys, V0VerifierBindingKeyPrefix)
	} else if filter.TxLevelOnly {
		bindingKeys = append(bindingKeys, V0ContractBindKeyPrefix)
	} else {
		bindingKeys = append(bindingKeys, V0VerifierBindingKeyPrefix, V0ContractBindKeyPrefix)
	}

	account := s.ctx.Account
	bindingSet := make(map[common.Address]struct{})
	bindings := make([]types.Binding, 0)
	for _, bindingKey := range bindingKeys {
		prefixStore := s.NewPrefixStore(bindingKey)
		storeKey := AccountKey(account.Bytes())
		rawJSON, err := s.Load(prefixStore, storeKey)
		if err != nil {
			return nil, err
		}

		aspectMetas := make([]*evmtypes.AspectMeta, 0)
		if len(rawJSON) > 0 {
			if err := json.Unmarshal(rawJSON, &aspectMetas); err != nil {
				return nil, store.ErrStorageCorrupted
			}
		}

		for _, aspectMeta := range aspectMetas {
			if _, ok := bindingSet[aspectMeta.Id]; ok {
				continue
			}
			bindings = append(bindings, types.Binding{
				Account:  aspectMeta.Id,
				Version:  aspectMeta.Version.Uint64(),
				Priority: int8(aspectMeta.Priority),
			})
			bindingSet[aspectMeta.Id] = struct{}{}
		}
	}

	return bindings, nil
}
