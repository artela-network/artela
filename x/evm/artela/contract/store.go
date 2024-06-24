package contract

import (
	"bytes"
	"encoding/json"
	"errors"
	types2 "github.com/artela-network/aspect-runtime/types"
	"math"
	"math/big"
	"sort"
	"strings"

	artelasdkType "github.com/artela-network/aspect-core/types"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"golang.org/x/exp/slices"

	"github.com/artela-network/artela/x/evm/artela/types"
	evmtypes "github.com/artela-network/artela/x/evm/types"
)

const (
	storageLoadCost     = 50
	storageStoreCost    = 20000
	storageSaveCodeCost = 1000
	storageUpdateCost   = 5000
)

type gasMeter struct {
	gas uint64
}

func newGasMeter(gas uint64) *gasMeter {
	return &gasMeter{
		gas: gas,
	}
}

func (m *gasMeter) measureStorageUpdate(dataLen int) error {
	return m.consume(dataLen, storageUpdateCost)
}

func (m *gasMeter) measureStorageCodeSave(dataLen int) error {
	return m.consume(dataLen, storageSaveCodeCost)
}

func (m *gasMeter) measureStorageStore(dataLen int) error {
	return m.consume(dataLen, storageStoreCost)
}

func (m *gasMeter) measureStorageLoad(dataLen int) error {
	return m.consume(dataLen, storageLoadCost)
}

func (m *gasMeter) remainingGas() uint64 {
	return m.gas
}

func (m *gasMeter) consume(dataLen int, gasCostPer32Bytes uint64) error {
	gas := ((uint64(dataLen) + 32) >> 5) * gasCostPer32Bytes
	if m.gas < gas {
		m.gas = 0
		return types2.OutOfGasError
	}
	m.gas -= gas
	return nil
}

type AspectStore struct {
	storeKey storetypes.StoreKey
	logger   log.Logger
}

type bindingQueryFunc func(sdk.Context, common.Address) ([]*types.AspectMeta, error)

func NewAspectStore(storeKey storetypes.StoreKey, logger log.Logger) *AspectStore {
	return &AspectStore{
		storeKey: storeKey,
		logger:   logger,
	}
}

func (k *AspectStore) newPrefixStore(ctx sdk.Context, fixKey string) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefix(fixKey))
}

func (k *AspectStore) BumpAspectVersion(ctx sdk.Context, aspectId common.Address, gas uint64) (*uint256.Int, uint64, error) {
	meter := newGasMeter(gas)
	version := k.getAspectLastVersion(ctx, aspectId)

	newVersion := version.Add(version, uint256.NewInt(1))
	if err := k.storeAspectVersion(ctx, aspectId, newVersion, meter); err != nil {
		return nil, meter.remainingGas(), err
	}

	return newVersion, meter.remainingGas(), nil
}

// StoreAspectCode aspect code
func (k *AspectStore) StoreAspectCode(ctx sdk.Context, aspectId common.Address, code []byte, version *uint256.Int, gas uint64) (uint64, error) {
	meter := newGasMeter(gas)
	if err := meter.measureStorageCodeSave(len(code)); err != nil {
		return meter.remainingGas(), err
	}

	// store code
	codeStore := k.newPrefixStore(ctx, types.AspectCodeKeyPrefix)
	versionKey := types.AspectVersionKey(
		aspectId.Bytes(),
		version.Bytes(),
	)
	codeStore.Set(versionKey, code)

	k.logger.Info("saved aspect code", "id", aspectId.Hex(), "version", version.String())
	return meter.remainingGas(), nil
}

func (k *AspectStore) GetAspectCode(ctx sdk.Context, aspectId common.Address, version *uint256.Int) ([]byte, *uint256.Int) {
	codeStore := k.newPrefixStore(ctx, types.AspectCodeKeyPrefix)
	if version == nil {
		version = k.getAspectLastVersion(ctx, aspectId)
	}

	if version.Cmp(zero) == 0 {
		return nil, zero
	}

	versionKey := types.AspectVersionKey(
		aspectId.Bytes(),
		version.Bytes(),
	)
	code := codeStore.Get(versionKey)
	return code, version
}

// storeAspectVersion version
func (k *AspectStore) storeAspectVersion(ctx sdk.Context, aspectId common.Address, version *uint256.Int, meter *gasMeter) error {
	var err error
	if version.Cmp(one) == 0 {
		err = meter.measureStorageStore(32)
	} else {
		err = meter.measureStorageUpdate(32)
	}
	if err != nil {
		return err
	}

	versionStore := k.newPrefixStore(ctx, types.AspectCodeVersionKeyPrefix)
	versionKey := types.AspectIdKey(aspectId.Bytes())
	versionStore.Set(versionKey, version.Bytes())

	k.logger.Info("saved aspect version info", "id", aspectId.Hex(), "version", version.String())
	return nil
}

func (k *AspectStore) GetAspectLastVersion(ctx sdk.Context, aspectId common.Address) *uint256.Int {
	return k.getAspectLastVersion(ctx, aspectId)
}

func (k *AspectStore) getAspectLastVersion(ctx sdk.Context, aspectId common.Address) *uint256.Int {
	aspectVersionStore := k.newPrefixStore(ctx, types.AspectCodeVersionKeyPrefix)
	versionKey := types.AspectIdKey(aspectId.Bytes())
	version := uint256.NewInt(0)
	if data := aspectVersionStore.Get(versionKey); data != nil || len(data) > 0 {
		version.SetBytes(data)
	}

	return version
}

// StoreAspectProperty
//
//	@Description:  property storage format
//	 1. {aspectid,key}=>{prperty value}
//	 2. {aspectid,"AspectPropertyAllKeyPrefix"}=>"key1,key2,key3..."
//	@receiver k
//	@param ctx
//	@param aspectId
//	@param prop
//	@return error
func (k *AspectStore) StoreAspectProperty(ctx sdk.Context, aspectId common.Address, prop []types.Property, gas uint64) (uint64, error) {
	meter := newGasMeter(gas)
	if len(prop) == 0 {
		return gas, nil
	}

	// get treemap value
	aspectConfigStore := k.newPrefixStore(ctx, types.AspectPropertyKeyPrefix)
	// get all property key
	propertyAllKey, err := k.getAspectPropertyValue(ctx, aspectId, types.AspectPropertyAllKeyPrefix, meter)
	if err != nil {
		return meter.remainingGas(), err
	}

	keySet := treeset.NewWithStringComparator()
	// add propertyAllKey to keySet
	if len(propertyAllKey) > 0 {
		splitData := strings.Split(string(propertyAllKey), types.AspectPropertyAllKeySplit)
		for _, datum := range splitData {
			keySet.Add(datum)
		}
	}
	for i := range prop {
		key := prop[i].Key
		// add key and deduplicate
		keySet.Add(key)
	}
	// check key limit
	if keySet.Size() > types.AspectPropertyLimit {
		return meter.remainingGas(), errors.New("aspect property limit exceeds")
	}

	// store property key
	for i := range prop {
		key := prop[i].Key
		value := prop[i].Value

		if err := meter.measureStorageCodeSave(len(key) + len(value)); err != nil {
			k.logger.Error("unable to save property", "err", err, "key", key, "value", value)
			return meter.remainingGas(), err
		}

		// store
		aspectPropertyKey := types.AspectPropertyKey(
			aspectId.Bytes(),
			[]byte(key),
		)

		aspectConfigStore.Set(aspectPropertyKey, value)

		k.logger.Info("aspect property updated", "aspect", aspectId.Hex(), "key", key, "value", value)
	}

	// store AspectPropertyAllKey
	keyAry := make([]string, keySet.Size())
	for i, key := range keySet.Values() {
		keyAry[i] = key.(string)
	}
	join := strings.Join(keyAry, types.AspectPropertyAllKeySplit)
	allPropertyKeys := types.AspectPropertyKey(
		aspectId.Bytes(),
		[]byte(types.AspectPropertyAllKeyPrefix),
	)
	aspectConfigStore.Set(allPropertyKeys, []byte(join))

	return meter.remainingGas(), nil
}

func (k *AspectStore) GetAspectPropertyValue(ctx sdk.Context, aspectId common.Address, propertyKey string, gas uint64) ([]byte, uint64, error) {
	meter := newGasMeter(gas)
	value, err := k.getAspectPropertyValue(ctx, aspectId, propertyKey, meter)
	return value, meter.remainingGas(), err
}

func (k *AspectStore) getAspectPropertyValue(ctx sdk.Context, aspectId common.Address, propertyKey string, meter *gasMeter) ([]byte, error) {
	codeStore := k.newPrefixStore(ctx, types.AspectPropertyKeyPrefix)
	aspectPropertyKey := types.AspectPropertyKey(
		aspectId.Bytes(),
		[]byte(propertyKey),
	)

	value := codeStore.Get(aspectPropertyKey)
	return value, meter.measureStorageLoad(len(propertyKey) + len(value))
}

func (k *AspectStore) BindTxAspect(ctx sdk.Context, account common.Address, aspectId common.Address, aspectVersion *uint256.Int, priority int8) error {
	return k.saveBindingInfo(ctx, account, aspectId, aspectVersion, priority,
		k.GetTxLevelAspects, types.ContractBindKeyPrefix, math.MaxUint8)
}

func (k *AspectStore) BindVerificationAspect(ctx sdk.Context, account common.Address, aspectId common.Address, aspectVersion *uint256.Int, priority int8, isContractAccount bool) error {
	if isContractAccount {
		// contract can have only 1 verifier
		return k.saveBindingInfo(ctx, account, aspectId, aspectVersion, priority,
			k.GetVerificationAspects, types.VerifierBindingKeyPrefix, 1)
	} else {
		// EoA can have multiple verifiers
		return k.saveBindingInfo(ctx, account, aspectId, aspectVersion, priority,
			k.GetVerificationAspects, types.VerifierBindingKeyPrefix, math.MaxUint8)
	}
}

func (k *AspectStore) saveBindingInfo(ctx sdk.Context, account common.Address, aspectId common.Address,
	aspectVersion *uint256.Int, priority int8, queryBinding bindingQueryFunc, bindingNameSpace string, limit int,
) error {
	// check aspect existence
	code, version := k.GetAspectCode(ctx, aspectId, aspectVersion)
	if code == nil || version == nil {
		return errors.New("aspect not found")
	}

	// get transaction level aspect binding relationships
	bindings, err := queryBinding(ctx, account)
	if err != nil {
		return err
	}

	if len(bindings) >= limit {
		return errors.New("binding limit exceeded")
	}

	// check duplicates
	existing := -1
	for index, binding := range bindings {
		if bytes.Equal(binding.Id.Bytes(), aspectId.Bytes()) {
			// ignore if binding already exists
			if binding.Priority == int64(priority) &&
				binding.Version.Cmp(aspectVersion) == 0 {
				return nil
			}

			// record existing, replace later
			existing = index
			break
		}
	}

	newAspect := &types.AspectMeta{
		Id:       aspectId,
		Version:  version,
		Priority: int64(priority),
	}

	// replace existing binding
	if existing > 0 {
		bindings[existing] = newAspect
	} else {
		bindings = append(bindings, newAspect)
	}

	// re-sort aspects by priority
	if limit != 1 {
		sort.Slice(bindings, types.NewBindingPriorityComparator(bindings))
	}

	jsonBytes, err := json.Marshal(bindings)
	if err != nil {
		return err
	}

	// save bindings
	aspectBindingStore := k.newPrefixStore(ctx, bindingNameSpace)
	aspectPropertyKey := types.AccountKey(
		account.Bytes(),
	)
	aspectBindingStore.Set(aspectPropertyKey, jsonBytes)

	k.logger.Info("binding info saved",
		"aspect", aspectId.Hex(),
		"contract", account.Hex(),
		"bindings", string(jsonBytes),
	)

	return nil
}

func (k *AspectStore) UnBindContractAspects(ctx sdk.Context, contract common.Address, aspectId common.Address) error {
	txAspectBindings, err := k.GetTxLevelAspects(ctx, contract)
	if err != nil {
		return err
	}
	toDelete := slices.IndexFunc(txAspectBindings, func(meta *types.AspectMeta) bool {
		return bytes.Equal(meta.Id.Bytes(), aspectId.Bytes())
	})
	if toDelete < 0 {
		// not found
		return nil
	}
	txAspectBindings = slices.Delete(txAspectBindings, toDelete, toDelete+1)
	jsonBytes, err := json.Marshal(txAspectBindings)
	if err != nil {
		return err
	}
	// store
	contractBindingStore := k.newPrefixStore(ctx, types.ContractBindKeyPrefix)

	aspectPropertyKey := types.AccountKey(
		contract.Bytes(),
	)
	contractBindingStore.Set(aspectPropertyKey, jsonBytes)

	k.logger.Info("aspect unbound", "aspect", aspectId.Hex(), "contract", contract.String())
	return nil
}

func (k *AspectStore) GetTxLevelAspects(ctx sdk.Context, contract common.Address) ([]*types.AspectMeta, error) {
	return k.getAccountBondAspects(ctx, contract, types.ContractBindKeyPrefix)
}

func (k *AspectStore) GetVerificationAspects(ctx sdk.Context, account common.Address) ([]*types.AspectMeta, error) {
	return k.getAccountBondAspects(ctx, account, types.VerifierBindingKeyPrefix)
}

func (k *AspectStore) getAccountBondAspects(ctx sdk.Context, account common.Address, bindingPrefix string) ([]*types.AspectMeta, error) {
	// retrieve raw binding store
	aspectBindingStore := k.newPrefixStore(ctx, bindingPrefix)
	accountKey := types.AccountKey(
		account.Bytes(),
	)
	rawJSON := aspectBindingStore.Get(accountKey)

	var bindings []*types.AspectMeta
	if len(rawJSON) == 0 {
		return bindings, nil
	}
	if err := json.Unmarshal(rawJSON, &bindings); err != nil {
		return nil, errors.New("failed to unmarshal aspect bindings")
	}
	return bindings, nil
}

func (k *AspectStore) ChangeBoundAspectVersion(ctx sdk.Context, account common.Address, aspectId common.Address, version uint64, isContract bool) error {
	bindingStoreKeys := make([]string, 0, 2)
	bindingStoreKeys = append(bindingStoreKeys, types.VerifierBindingKeyPrefix)
	if isContract {
		bindingStoreKeys = append(bindingStoreKeys, types.ContractBindKeyPrefix)
	}

	for _, bindingStoreKey := range bindingStoreKeys {
		metas, err := k.getAccountBondAspects(ctx, account, bindingStoreKey)
		if err != nil {
			return err
		}

		oldver := uint64(0)
		for _, aspect := range metas {
			if bytes.Equal(aspect.Id.Bytes(), aspectId.Bytes()) {
				oldver = aspect.Version.Uint64()
				aspect.Version = uint256.NewInt(version)
				break
			}
		}

		if oldver == 0 {
			k.logger.Debug("aspect not bound", "aspect", aspectId.Hex(), "account", account.String())
			return nil
		}

		jsonBytes, err := json.Marshal(metas)
		if err != nil {
			return err
		}

		bindingStore := k.newPrefixStore(ctx, bindingStoreKey)
		bindingKey := types.AccountKey(
			account.Bytes(),
		)
		bindingStore.Set(bindingKey, jsonBytes)

		k.logger.Info("aspect bound version changed", "aspect", aspectId.Hex(), "account", account.String(), "old", oldver, "new", version)
	}

	return nil
}

func (k *AspectStore) GetAspectRefValue(ctx sdk.Context, aspectId common.Address) (*treeset.Set, error) {
	aspectRefStore := k.newPrefixStore(ctx, types.AspectRefKeyPrefix)
	aspectPropertyKey := types.AspectIdKey(
		aspectId.Bytes(),
	)

	rawTree := aspectRefStore.Get(aspectPropertyKey)
	if rawTree == nil {
		return nil, nil
	}

	set := treeset.NewWithStringComparator()
	if err := set.UnmarshalJSON(rawTree); err != nil {
		return nil, err
	}
	return set, nil
}

func (k *AspectStore) StoreAspectRefValue(ctx sdk.Context, account common.Address, aspectId common.Address) error {
	dataSet, err := k.GetAspectRefValue(ctx, aspectId)
	if err != nil {
		return err
	}
	if dataSet == nil {
		dataSet = treeset.NewWithStringComparator()
	}
	dataSet.Add(account.String())
	jsonBytes, err := dataSet.MarshalJSON()
	if err != nil {
		return err
	}
	// store
	aspectRefStore := k.newPrefixStore(ctx, types.AspectRefKeyPrefix)

	aspectIdKey := types.AspectIdKey(
		aspectId.Bytes(),
	)
	aspectRefStore.Set(aspectIdKey, jsonBytes)

	k.logger.Info("aspect bound", "aspect", aspectId.Hex(), "account", account.Hex())
	return nil
}

func (k *AspectStore) UnbindAspectRefValue(ctx sdk.Context, account common.Address, aspectId common.Address) error {
	dataSet, err := k.GetAspectRefValue(ctx, aspectId)
	if err != nil {
		return err
	}
	if dataSet == nil {
		return nil
	}
	// remove account
	dataSet.Remove(account.String())
	// marshal set and put treemap with new blockHeight
	jsonBytes, err := dataSet.MarshalJSON()
	if err != nil {
		return err
	}
	// store
	aspectRefStore := k.newPrefixStore(ctx, types.AspectRefKeyPrefix)
	aspectPropertyKey := types.AspectIdKey(
		aspectId.Bytes(),
	)
	aspectRefStore.Set(aspectPropertyKey, jsonBytes)

	k.logger.Info("aspect unbound", "aspect", aspectId.Hex(), "account", account.Hex())
	return nil
}

func (k *AspectStore) UnBindVerificationAspect(ctx sdk.Context, account common.Address, aspectId common.Address) error {
	bindings, err := k.GetVerificationAspects(ctx, account)
	if err != nil {
		return err
	}

	toDelete := slices.IndexFunc(bindings, func(meta *types.AspectMeta) bool {
		return bytes.Equal(meta.Id.Bytes(), aspectId.Bytes())
	})

	if toDelete < 0 {
		// not found
		return nil
	}
	// delete existing
	bindings = slices.Delete(bindings, toDelete, toDelete+1)

	sort.Slice(bindings, types.NewBindingPriorityComparator(bindings))
	jsonBytes, err := json.Marshal(bindings)
	if err != nil {
		return err
	}

	// save bindings
	aspectBindingStore := k.newPrefixStore(ctx, types.VerifierBindingKeyPrefix)
	aspectPropertyKey := types.AccountKey(
		account.Bytes(),
	)
	aspectBindingStore.Set(aspectPropertyKey, jsonBytes)

	k.logger.Info("aspect unbound", "aspect", aspectId.Hex(), "account", account.String())
	return nil
}

// StoreAspectJP
//
//	@Description: Stores the execute conditions of the Aspect Join point. {aspectId,version,'AspectRunJoinPointKey'}==>{value}
//	@receiver k
//	@param ctx
//	@param aspectId
//	@param version: aspect version ,optional，Default Aspect last version
//	@param point  JoinPointRunType value, @see join_point_type.go
//	@return bool Execute Result
func (k *AspectStore) StoreAspectJP(ctx sdk.Context, aspectId common.Address, version uint256.Int, point *big.Int) {
	// check point
	if _, ok := artelasdkType.CheckIsJoinPoint(point); !ok {
		// Default store 0
		point = big.NewInt(0)
	}

	aspectPropertyStore := k.newPrefixStore(ctx, types.AspectJoinPointRunKeyPrefix)
	aspectPropertyKey := types.AspectArrayKey(
		aspectId.Bytes(),
		version.Bytes(),
		[]byte(types.AspectRunJoinPointKey),
	)
	aspectPropertyStore.Set(aspectPropertyKey, point.Bytes())
}

// GetAspectJP
//
//	@Description: get Aspect Join point run
//	@receiver k
//	@param ctx
//	@param aspectId
//	@param version
//	@return *big.Int
func (k *AspectStore) GetAspectJP(ctx sdk.Context, aspectId common.Address, version *uint256.Int) *big.Int {
	// Default last Aspect version
	if version == nil {
		version = k.GetAspectLastVersion(ctx, aspectId)
	}
	codeStore := k.newPrefixStore(ctx, types.AspectJoinPointRunKeyPrefix)
	aspectPropertyKey := types.AspectArrayKey(
		aspectId.Bytes(),
		version.Bytes(),
		[]byte(types.AspectRunJoinPointKey),
	)
	get := codeStore.Get(aspectPropertyKey)

	if nil != get && len(get) > 0 {
		return new(big.Int).SetBytes(get)
	} else {
		return new(big.Int)
	}
}
