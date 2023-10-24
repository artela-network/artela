package contract

import (
	"bytes"
	"encoding/json"
	"sort"

	"cosmossdk.io/errors"

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

type AspectStore struct {
	storeKey storetypes.StoreKey
}

func NewAspectStore(storeKey storetypes.StoreKey) *AspectStore {
	return &AspectStore{
		storeKey: storeKey,
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
func (k *AspectStore) StoreAspectCode(ctx sdk.Context, aspectId common.Address, code []byte) {
	// get last value
	version := k.GetAspectLastVersion(ctx, aspectId)

	// store code
	codeStore := k.newPrefixStore(ctx, types.AspectCodeKeyPrefix)
	newVersion := version.Add(version, uint256.NewInt(1))
	versionKey := types.AspectVersionKey(
		aspectId.Bytes(),
		newVersion.Bytes(),
	)
	codeStore.Set(versionKey, code)
	// update last version
	k.StoreAspectVersion(ctx, aspectId, newVersion)
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

// StoreAspectProperty Property
func (k *AspectStore) StoreAspectProperty(ctx sdk.Context, aspectId common.Address, prop []types.Property) {
	// get treemap value
	aspectPropertyStore := k.newPrefixStore(ctx, types.AspectPropertyKeyPrefix)
	for i := range prop {
		key := prop[i].Key
		value := prop[i].Value
		// store
		aspectPropertyKey := types.AspectPropertyKey(
			aspectId.Bytes(),
			[]byte(key),
		)
		aspectPropertyStore.Set(aspectPropertyKey, []byte(value))
	}
}

func (k *AspectStore) GetAspectPropertyValue(ctx sdk.Context, aspectId common.Address, propertyKey string) string {
	if types.AspectProofKey == propertyKey || types.AspectAccountKey == propertyKey {
		// Block query of account and Proof
		return ""
	}
	codeStore := k.newPrefixStore(ctx, types.AspectPropertyKeyPrefix)
	aspectPropertyKey := types.AspectPropertyKey(
		aspectId.Bytes(),
		[]byte(propertyKey),
	)
	get := codeStore.Get(aspectPropertyKey)
	if nil != get && len(get) > 0 {
		return string(get)
	} else {
		return ""
	}
}

func (k *AspectStore) BindContractAspects(ctx sdk.Context, contract common.Address, aspectId common.Address, aspectVersion *uint256.Int, priority int8) error {
	// get treemap value
	bindings, err := k.GetContractBondAspects(ctx, contract)
	if err != nil {
		return err
	}

	code, version := k.GetAspectCode(ctx, aspectId, aspectVersion)
	if code == nil || version == nil {
		return errors.Wrap(nil, "aspect not exist")
	}

	newAspect := &types.AspectMeta{
		Id:       aspectId,
		Version:  version,
		Priority: int64(priority),
	}

	bindings = append(bindings, newAspect)
	sort.Slice(bindings, types.NewBindingPriorityComparator(bindings))

	jsonBytes, err := json.Marshal(bindings)
	if err != nil {
		return err
	}
	// store
	contractBindingStore := k.newPrefixStore(ctx, types.ContractBindKeyPrefix)
	aspectPropertyKey := types.ContractKey(
		contract.Bytes(),
	)
	contractBindingStore.Set(aspectPropertyKey, jsonBytes)

	return nil
}

func (k *AspectStore) UnBindContractAspects(ctx sdk.Context, contract common.Address, aspectId common.Address) error {
	bindings, err := k.GetContractBondAspects(ctx, contract)
	if err != nil {
		return err
	}
	toDelete := slices.IndexFunc(bindings, func(meta *types.AspectMeta) bool {
		return bytes.Equal(meta.Id.Bytes(), aspectId.Bytes())
	})
	if toDelete < 0 {
		return errors.Wrapf(nil, "aspect %s not bound with contract %s", aspectId.Hex(), contract.Hex())
	}
	bindings = slices.Delete(bindings, toDelete, toDelete)
	jsonBytes, err := json.Marshal(bindings)
	if err != nil {
		return err
	}
	// store
	contractBindingStore := k.newPrefixStore(ctx, types.ContractBindKeyPrefix)

	aspectPropertyKey := types.ContractKey(
		contract.Bytes(),
	)
	contractBindingStore.Set(aspectPropertyKey, jsonBytes)
	return nil
}

func (k *AspectStore) GetContractBondAspects(ctx sdk.Context, contract common.Address) ([]*types.AspectMeta, error) {
	// retrieve raw binding store
	contractBindingStore := k.newPrefixStore(ctx, types.ContractBindKeyPrefix)
	contractKey := types.ContractKey(
		contract.Bytes(),
	)
	sortAry := contractBindingStore.Get(contractKey)

	var bindings []*types.AspectMeta
	if len(sortAry) == 0 {
		return bindings, nil
	}
	if err := json.Unmarshal(sortAry, &bindings); err != nil {
		return nil, errors.Wrap(err, "unable to deserialize value bytes")
	}
	return bindings, nil
}

func (k *AspectStore) ChangeBoundAspectVersion(ctx sdk.Context, contract common.Address, aspectId common.Address, version uint64) error {
	meta, err := k.GetContractBondAspects(ctx, contract)
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
	aspectPropertyKey := types.ContractKey(
		contract.Bytes(),
	)
	contractBindingStore.Set(aspectPropertyKey, jsonBytes)
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

func (k *AspectStore) StoreAspectRefValue(ctx sdk.Context, contract common.Address, aspectId common.Address) error {
	dataSet, err := k.GetAspectRefValue(ctx, contract)
	if err != nil {
		return err
	}
	if dataSet == nil {
		dataSet = treeset.NewWithStringComparator()
	}
	dataSet.Add(contract.String())
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
	return nil
}

func (k *AspectStore) UnbindAspectRefValue(ctx sdk.Context, contract common.Address, aspectId common.Address) error {
	dataSet, err := k.GetAspectRefValue(ctx, contract)
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
	return nil
}
