package v1

import (
	aspect "github.com/artela-network/aspect-core/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artela-network/artela/x/aspect/store"
	v0 "github.com/artela-network/artela/x/aspect/store/v0"
	"github.com/artela-network/artela/x/aspect/types"
)

var _ store.AccountStore = (*accountStore)(nil)

type accountStore struct {
	BaseStore

	ctx *types.AccountStoreContext
}

// NewAccountStore creates a new instance of account store.
func NewAccountStore(ctx *types.AccountStoreContext) store.AccountStore {
	var meter v0.GasMeter
	if ctx.ChargeGas() {
		meter = v0.NewGasMeter(ctx)
	} else {
		meter = v0.NewNoOpGasMeter(ctx)
	}

	return &accountStore{
		BaseStore: NewBaseStore(ctx.CosmosContext().Logger(), meter, ctx.CosmosContext().KVStore(ctx.AspectStoreKey())),
		ctx:       ctx,
	}
}

func (a *accountStore) LoadAccountBoundAspects(filter types.BindingFilter) ([]types.Binding, error) {
	allBindings, err := a.getAllBindings()
	if err != nil {
		return nil, err
	}

	var result []types.Binding
	for _, binding := range allBindings {
		if filter.JoinPoint != nil {
			jp, ok := aspect.JoinPointRunType_value[string(*filter.JoinPoint)]
			if !ok {
				return nil, store.ErrInvalidJoinPoint
			}
			if (binding.JoinPoint | uint16(jp)) == 0 {
				continue
			}

			result = append(result, types.Binding(binding))
		} else if filter.VerifierOnly && aspect.CheckIsTxVerifier(int64(binding.JoinPoint)) {
			result = append(result, types.Binding(binding))
		} else if filter.TxLevelOnly && aspect.CheckIsTransactionLevel(int64(binding.JoinPoint)) {
			result = append(result, types.Binding(binding))
		} else {
			result = append(result, types.Binding(binding))
		}
	}

	return result, nil
}

func (a *accountStore) StoreBinding(aspectID common.Address, version uint64, joinPoint uint64, priority int8, isCA bool) error {
	allBindings, err := a.getAllBindings()
	if err != nil {
		return err
	}

	i64JP := int64(joinPoint)
	isTxVerifier := aspect.CheckIsTxVerifier(i64JP)
	isTxLevel := aspect.CheckIsTransactionLevel(i64JP)

	if !isTxLevel && !isTxVerifier {
		return store.ErrNoJoinPoint
	}

	if !isCA && !isTxVerifier {
		return store.ErrBoundNonVerifierWithEOA
	}

	if len(allBindings) >= maxAspectBoundLimit {
		return store.ErrBindingLimitExceeded
	}

	for _, binding := range allBindings {
		if binding.Account == aspectID {
			return store.ErrAlreadyBound
		}
		if isTxVerifier && isCA && aspect.CheckIsTxVerifier(int64(binding.JoinPoint)) {
			// contract only allowed to bind 1 verifier
			return store.ErrBindingLimitExceeded
		}
	}

	newBinding := Binding{
		Account:   aspectID,
		Version:   version,
		Priority:  priority,
		JoinPoint: uint16(joinPoint),
	}
	allBindings = append(allBindings, newBinding)

	key := store.NewKeyBuilder(V1AccountBindingKeyPrefix).AppendBytes(a.ctx.Account.Bytes()).Build()
	bindingsBytes, err := Bindings(allBindings).MarshalText()
	if err != nil {
		return err
	}

	return a.Store(key, bindingsBytes)
}

func (a *accountStore) getAllBindings() ([]Binding, error) {
	key := store.NewKeyBuilder(V1AccountBindingKeyPrefix).AppendBytes(a.ctx.Account.Bytes()).Build()
	bindings, err := a.Load(key)
	if err != nil {
		return nil, err
	}

	result := make([]Binding, 0, len(bindings)/32)
	for i := 0; i < len(bindings); i += 32 {
		data := bindings[i : i+32]
		if common.Hash(data) == emptyHash {
			// EOF
			return result, nil
		}

		var binding Binding
		if err := binding.UnmarshalText(data); err != nil {
			return nil, err
		}
		result = append(result, binding)
	}

	return result, nil
}

func (a *accountStore) RemoveBinding(aspectID common.Address, _ uint64, _ bool) error {
	key := store.NewKeyBuilder(V1AccountBindingKeyPrefix).AppendBytes(a.ctx.Account.Bytes()).Build()
	allBindings, err := a.getAllBindings()
	if err != nil {
		return err
	}

	toDelete := -1
	for i, binding := range allBindings {
		if binding.Account == aspectID {
			toDelete = i
			break
		}
	}

	if toDelete < 0 {
		// if not bind, return nil
		return nil
	}

	allBindings = append(allBindings[:toDelete], allBindings[toDelete+1:]...)
	bindingsBytes, err := Bindings(allBindings).MarshalText()
	if err != nil {
		return err
	}
	return a.Store(key, bindingsBytes)
}

func (a *accountStore) Used() (bool, error) {
	key := store.NewKeyBuilder(store.AspectProtocolInfoKeyPrefix).AppendBytes(a.ctx.Account.Bytes()).Build()
	protocol, err := a.Load(key)
	if err != nil {
		return false, err
	}

	var protocolVersion store.ProtocolVersion
	if err := protocolVersion.UnmarshalText(protocol); err != nil {
		return false, err
	}

	return protocolVersion == store.ProtocolVersion(1), nil
}

func (a *accountStore) MigrateFrom(old store.AccountStore) error {
	//TODO implement me
	panic("implement me")
}

func (a *accountStore) Init() error {
	versionBytes, err := protocolVersion.MarshalText()
	if err != nil {
		return err
	}

	key := store.NewKeyBuilder(store.AspectProtocolInfoKeyPrefix).AppendBytes(a.ctx.Account.Bytes()).Build()
	return a.Store(key, versionBytes)
}
