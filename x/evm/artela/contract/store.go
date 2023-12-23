package contract

import (
	"bytes"
	"cosmossdk.io/errors"
	"encoding/json"
	"fmt"
	artelasdkType "github.com/artela-network/aspect-core/types"
	"math"
	"math/big"
	"sort"
	"strings"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/status-im/keycard-go/hexutils"
	"golang.org/x/exp/slices"

	"github.com/artela-network/artela/x/evm/artela/types"
	evmtypes "github.com/artela-network/artela/x/evm/types"
)

type AspectStore struct {
	storeKey storetypes.StoreKey

	logger log.Logger
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

func (k *AspectStore) RemoveBlockLevelAspect(ctx sdk.Context, aspectId common.Address) error {
	dataSet, err := k.GetBlockLevelAspects(ctx)
	if err != nil {
		return err
	}
	if dataSet == nil {
		return nil
	}
	delete(dataSet, aspectId.String())
	jsonBytes, err := json.Marshal(dataSet)
	if err != nil {
		return err
	}
	// store
	store := k.newPrefixStore(ctx, types.AspectBlockKeyPrefix)
	aspectBlockKey := types.AspectBlockKey()
	store.Set(aspectBlockKey, jsonBytes)

	k.logger.Debug(
		fmt.Sprintf("setState: RemoveBlockLevelAspect"),
		"key", string(aspectBlockKey),
		"aspect-id", aspectId.Hex(),
		"data-set", string(jsonBytes),
	)
	return nil
}

// StoreBlockLevelAspect key="AspectBlock" value=map[string]int64
func (k *AspectStore) StoreBlockLevelAspect(ctx sdk.Context, aspectId common.Address) error {
	dataSet, err := k.GetBlockLevelAspects(ctx)
	if err != nil {
		return err
	}
	if dataSet == nil {
		// order by
		dataSet = make(map[string]int64)
	}
	// oder by block height
	dataSet[aspectId.String()] = ctx.BlockHeight()
	jsonBytes, err := json.Marshal(dataSet)
	if err != nil {
		return err
	}
	//
	// prefix
	// kv
	store := k.newPrefixStore(ctx, types.AspectBlockKeyPrefix)
	aspectBlockKey := types.AspectBlockKey()
	store.Set(aspectBlockKey, jsonBytes)

	k.logger.Debug(
		fmt.Sprintf("setState: StoreBlockLevelAspect"),
		"key", string(aspectBlockKey),
		"aspect-id", aspectId.Hex(),
		"data-set", string(jsonBytes),
	)
	return nil
}

func (k *AspectStore) GetBlockLevelAspects(ctx sdk.Context) (map[string]int64, error) {
	store := k.newPrefixStore(ctx, types.AspectBlockKeyPrefix)
	blockKey := types.AspectBlockKey()
	get := store.Get(blockKey)
	if get == nil {
		return nil, nil
	}
	blockMap := make(map[string]int64)
	if err := json.Unmarshal(get, &blockMap); err != nil {
		return nil, err
	}
	return blockMap, nil
}

// StoreAspectCode aspect code
func (k *AspectStore) StoreAspectCode(ctx sdk.Context, aspectId common.Address, code []byte) *uint256.Int {
	// get last value
	version := k.GetAspectLastVersion(ctx, aspectId)
	if len(code) == 0 {
		return version
	}

	// store code
	codeStore := k.newPrefixStore(ctx, types.AspectCodeKeyPrefix)
	newVersion := version.Add(version, uint256.NewInt(1))
	versionKey := types.AspectVersionKey(
		aspectId.Bytes(),
		newVersion.Bytes(),
	)
	codeStore.Set(versionKey, code)

	k.logger.Debug(
		fmt.Sprintf("setState: StoreAspectCode"),
		"aspect-id", aspectId.Hex(),
		"aspect-version", fmt.Sprintf("%d", newVersion),
		"aspect-code-hex", hexutils.BytesToHex(code),
	)

	// update last version
	k.StoreAspectVersion(ctx, aspectId, newVersion)
	return newVersion
}

func (k *AspectStore) GetAspectCode(ctx sdk.Context, aspectId common.Address, version *uint256.Int) ([]byte, *uint256.Int) {
	codeStore := k.newPrefixStore(ctx, types.AspectCodeKeyPrefix)
	if version == nil {
		version = k.GetAspectLastVersion(ctx, aspectId)
	}
	versionKey := types.AspectVersionKey(
		aspectId.Bytes(),
		version.Bytes(),
	)
	code := codeStore.Get(versionKey)
	return code, version
}

// StoreAspectVersion version
func (k *AspectStore) StoreAspectVersion(ctx sdk.Context, aspectId common.Address, version *uint256.Int) {
	versionStore := k.newPrefixStore(ctx, types.AspectCodeVersionKeyPrefix)
	versionKey := types.AspectIdKey(
		aspectId.Bytes(),
	)
	versionStore.Set(versionKey, version.Bytes())

	k.logger.Debug(
		fmt.Sprintf("setState: StoreAspectVersion"),
		"aspect-id", aspectId.Hex(),
		"aspect-version", fmt.Sprintf("%d", version),
	)
}

func (k *AspectStore) GetAspectLastVersion(ctx sdk.Context, aspectId common.Address) *uint256.Int {
	aspectVersionStore := k.newPrefixStore(ctx, types.AspectCodeVersionKeyPrefix)
	versionKey := types.AspectIdKey(
		aspectId.Bytes(),
	)
	version := uint256.NewInt(0)
	data := aspectVersionStore.Get(versionKey)
	if data != nil || len(data) > 0 {
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
func (k *AspectStore) StoreAspectProperty(ctx sdk.Context, aspectId common.Address, prop []types.Property) error {

	if len(prop) == 0 {
		return nil
	}

	// get treemap value
	aspectConfigStore := k.newPrefixStore(ctx, types.AspectPropertyKeyPrefix)
	//get all property key
	propertyAllKey := k.GetAspectPropertyValue(ctx, aspectId, types.AspectPropertyAllKeyPrefix)

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
		//add key and deduplicate
		keySet.Add(key)
	}
	//check key limit
	if keySet.Size() > types.AspectPropertyLimit {
		return errors.Wrapf(nil, "The maximum key limit is exceeded, and the maximum allowed is %d now available %d", types.AspectPropertyLimit, keySet.Size())
	}

	// store property key
	for i := range prop {
		key := prop[i].Key
		value := prop[i].Value

		// store
		aspectPropertyKey := types.AspectPropertyKey(
			aspectId.Bytes(),
			[]byte(key),
		)

		aspectConfigStore.Set(aspectPropertyKey, []byte(value))

		k.logger.Debug(
			fmt.Sprintf("setState: StoreAspectProperty"),
			"aspect-id", aspectId.Hex(),
			"aspect-property", fmt.Sprintf("%+v", prop),
		)
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

	return nil
}

func (k *AspectStore) GetAspectPropertyValue(ctx sdk.Context, aspectId common.Address, propertyKey string) []byte {
	if types.AspectProofKey == propertyKey || types.AspectAccountKey == propertyKey {
		// Block query of account and Proof
		return nil
	}
	codeStore := k.newPrefixStore(ctx, types.AspectPropertyKeyPrefix)
	aspectPropertyKey := types.AspectPropertyKey(
		aspectId.Bytes(),
		[]byte(propertyKey),
	)
	return codeStore.Get(aspectPropertyKey)
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
		return errors.Wrap(nil, "aspect not exist")
	}

	// get transaction level aspect binding relationships
	bindings, err := queryBinding(ctx, account)
	if err != nil {
		return err
	}

	if len(bindings) >= limit {
		return errors.Wrap(nil, "aspect binding limit exceeds")
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

	k.logger.Debug(
		fmt.Sprintf("setState: saveBindingInfo"),
		"aspect-id", aspectId.Hex(),
		"countract", account.String(),
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
		return errors.Wrapf(nil, "aspect %s not bound with contract %s", aspectId.Hex(), contract.Hex())
	}
	txAspectBindings = slices.Delete(txAspectBindings, toDelete, toDelete)
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

	k.logger.Debug(
		fmt.Sprintf("setState: UnBindContractAspects"),
		"aspect-id", aspectId.Hex(),
		"contract", contract.String(),
		"txAspectBindings", string(jsonBytes),
	)
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
		return nil, errors.Wrap(err, "unable to deserialize value bytes")
	}
	return bindings, nil
}

func (k *AspectStore) ChangeBoundAspectVersion(ctx sdk.Context, contract common.Address, aspectId common.Address, version uint64) error {
	meta, err := k.GetTxLevelAspects(ctx, contract)
	if err != nil {
		return err
	}
	hasAspect := false
	for _, aspect := range meta {
		if bytes.Equal(aspect.Id.Bytes(), aspectId.Bytes()) {
			aspect.Version = uint256.NewInt(version)
			hasAspect = true
		}
	}
	if !hasAspect {
		return nil
	}
	jsonBytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	// store
	contractBindingStore := k.newPrefixStore(ctx, types.ContractBindKeyPrefix)
	aspectPropertyKey := types.AccountKey(
		contract.Bytes(),
	)
	contractBindingStore.Set(aspectPropertyKey, jsonBytes)

	k.logger.Debug(
		fmt.Sprintf("setState: ChangeBoundAspectVersion"),
		"aspect-id", aspectId.Hex(),
		"contract", contract.String(),
		"aspects", string(jsonBytes),
	)
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

	k.logger.Debug(
		fmt.Sprintf("setState: StoreAspectRefValue"),
		"aspect-id", aspectId.Hex(),
		"context", account.Hex(),
		"aspects", string(jsonBytes),
	)
	return nil
}

func (k *AspectStore) UnbindAspectRefValue(ctx sdk.Context, contract common.Address, aspectId common.Address) error {
	dataSet, err := k.GetAspectRefValue(ctx, aspectId)
	if err != nil {
		return err
	}
	if dataSet == nil {
		return nil
	}
	// remove contract
	dataSet.Remove(contract.String())
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

	k.logger.Debug(
		fmt.Sprintf("setState: UnbindAspectRefValue"),
		"aspect-id", aspectId.Hex(),
		"context", contract.Hex(),
		"aspect-refvalue", string(jsonBytes),
	)
	return nil
}

func (k *AspectStore) UnBindVerificationAspect(ctx sdk.Context, account common.Address, aspectId common.Address) error {

	bindings, err := k.GetVerificationAspects(ctx, account)
	if err != nil {
		return err
	}
	existing := -1
	// check duplicates
	for index, binding := range bindings {
		if bytes.Equal(binding.Id.Bytes(), aspectId.Bytes()) {
			// delete Aspect id
			existing = index
			break
		}
	}
	if existing == -1 {
		return nil
	}
	// delete existing
	newBinding := append(bindings[:existing], bindings[existing+1:]...)

	sort.Slice(newBinding, types.NewBindingPriorityComparator(newBinding))
	jsonBytes, _ := json.Marshal(newBinding)
	if err != nil {
		return err
	}

	// save bindings
	aspectBindingStore := k.newPrefixStore(ctx, types.VerifierBindingKeyPrefix)
	aspectPropertyKey := types.AccountKey(
		account.Bytes(),
	)
	aspectBindingStore.Set(aspectPropertyKey, jsonBytes)

	k.logger.Debug(
		fmt.Sprintf("setState: UnBindVerificationAspect"),
		"aspect-id", aspectId.Hex(),
		"contract", account.Hex(),
		"newBinding", string(jsonBytes),
	)
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
func (k *AspectStore) StoreAspectJP(ctx sdk.Context, aspectId common.Address, version *uint256.Int, point *big.Int) error {
	if version.Uint64() == 0 || point.Int64() == 0 {
		return nil
	}
	//check point
	_, ok := artelasdkType.CheckIsJoinPoint(point)
	if !ok {
		// Default store 0
		point = big.NewInt(0)
	}

	//Default last Aspect version
	if version == nil {
		version = k.GetAspectLastVersion(ctx, aspectId)
	}

	aspectPropertyStore := k.newPrefixStore(ctx, types.AspectJoinPointRunKeyPrefix)
	// store
	aspectPropertyKey := types.AspectArrayKey(
		aspectId.Bytes(),
		version.Bytes(),
		[]byte(types.AspectRunJoinPointKey),
	)
	aspectPropertyStore.Set(aspectPropertyKey, point.Bytes())
	return nil
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
	//Default last Aspect version
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
