package v0

import (
	"math/big"
	"strings"

	artelasdkType "github.com/artela-network/aspect-core/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	aspectutils "github.com/artela-network/artela/x/aspect/common"
	"github.com/artela-network/artela/x/aspect/store"
	"github.com/artela-network/artela/x/aspect/types"
)

var _ store.AspectMetaStore = (*metaStore)(nil)

var reservedPropertyKeys = map[string]struct{}{
	V0AspectAccountKey:           {},
	V0AspectProofKey:             {},
	V0AspectPropertyAllKeyPrefix: {},
}

// metaStore is the version 0 aspect meta, this is no longer maintained.
// Just keeping for backward compatibility.
// Deprecated.
type metaStore struct {
	BaseStore

	latestVersionCache uint64
	ctx                *types.AspectStoreContext
}

// NewAspectMetaStore creates a new instance of aspect meta Store.
// Deprecated
func NewAspectMetaStore(ctx *types.AspectStoreContext, _ []byte) store.AspectMetaStore {
	var meter GasMeter
	if ctx.ChargeGas() {
		meter = NewGasMeter(ctx)
	} else {
		meter = NewNoOpGasMeter(ctx)
	}

	return &metaStore{
		BaseStore: NewBaseStore(meter, ctx.CosmosContext().KVStore(ctx.EVMStoreKey())),
		ctx:       ctx,
	}
}

func (s *metaStore) StoreMeta(meta *types.AspectMeta) error {
	// v0 Store saves paymaster and proof as aspect properties
	paymaster := types.Property{
		Key:   V0AspectAccountKey,
		Value: meta.PayMaster.Bytes(),
	}
	proof := types.Property{
		Key:   V0AspectProofKey,
		Value: meta.Proof,
	}
	return s.storeProperties([]types.Property{paymaster, proof})
}

func (s *metaStore) GetMeta() (*types.AspectMeta, error) {
	paymaster, err := s.GetProperty(0, V0AspectAccountKey)
	if err != nil {
		return nil, err
	}
	proof, err := s.GetProperty(0, V0AspectProofKey)
	if err != nil {
		return nil, err
	}

	return &types.AspectMeta{
		PayMaster: common.BytesToAddress(paymaster),
		Proof:     proof,
	}, nil
}

// StoreBinding stores the binding of the aspect with the given ID to the account.
func (s *metaStore) StoreBinding(account common.Address, _ uint64, _ uint64, _ int8) error {
	return s.saveAspectRef(s.NewPrefixStore(V0AspectRefKeyPrefix), s.ctx.AspectID, account)
}

// RemoveBinding removes the binding of the aspect with the given ID from the account.
func (s *metaStore) RemoveBinding(account common.Address) error {
	return s.removeAspectRef(s.NewPrefixStore(V0AspectRefKeyPrefix), account)
}

func (s *metaStore) MigrateFrom(_ store.AspectMetaStore) error {
	panic("cannot migrate to Store v0")
}

func (s *metaStore) Used() (bool, error) {
	v, err := s.GetLatestVersion()
	if err != nil {
		return false, err
	}
	return v > 0, err
}

func (s *metaStore) Init() error {
	// for v0 Store, we do not need to init anything
	return nil
}

// GetCode returns the code of the aspect with the given ID and version.
// If version is 0 or aspectID is empty, it returns * Store.ErrCodeNotFound.
func (s *metaStore) GetCode(version uint64) ([]byte, error) {
	aspectID := s.ctx.AspectID
	if version == 0 || aspectID == emptyAddress {
		return nil, store.ErrCodeNotFound
	}
	prefixStore := s.NewPrefixStore(V0AspectCodeKeyPrefix)
	storeKey := AspectVersionKey(aspectID.Bytes(), uint256.NewInt(version).Bytes())

	// we do not charge gas for code loading
	code := prefixStore.Get(storeKey)

	// stored code is already validated, so we can ignore the error here
	return aspectutils.ParseByteCode(code)
}

// GetVersionMeta returns the meta of the aspect with the given ID and version.
func (s *metaStore) GetVersionMeta(version uint64) (*types.VersionMeta, error) {
	return s.getMeta(s.ctx.AspectID, version)
}

func (s *metaStore) getMeta(aspectID common.Address, version uint64) (*types.VersionMeta, error) {
	u256Version := uint256.NewInt(version)
	prefixStore := s.NewPrefixStore(V0AspectJoinPointRunKeyPrefix)
	storeKey := AspectArrayKey(
		aspectID.Bytes(),
		u256Version.Bytes(),
		[]byte(V0AspectRunJoinPointKey),
	)

	joinPoint, err := s.Load(prefixStore, storeKey)
	if err != nil {
		return nil, err
	}

	// for v0 pay master / Store / code hash is not stored
	return &types.VersionMeta{
		JoinPoint: big.NewInt(0).SetBytes(joinPoint).Uint64(),
	}, nil
}

func (s *metaStore) GetLatestVersion() (uint64, error) {
	return s.getLatestVersion(s.ctx.AspectID)
}

// GetLatestVersion returns the latest version of the aspect with the given ID.
func (s *metaStore) getLatestVersion(aspectID common.Address) (uint64, error) {
	if s.latestVersionCache > 0 {
		return s.latestVersionCache, nil
	}

	prefixStore := s.NewPrefixStore(V0AspectCodeVersionKeyPrefix)
	storeKey := AspectIDKey(aspectID.Bytes())

	// v0 aspect Store uses uint256 to Store version
	version := uint256.NewInt(0)
	versionBytes, err := s.Load(prefixStore, storeKey)
	if err != nil {
		return 0, err
	}
	if versionBytes != nil || len(versionBytes) > 0 {
		version.SetBytes(versionBytes)
	}

	s.latestVersionCache = version.Uint64()
	return s.latestVersionCache, nil
}

// GetProperty returns the property of the aspect with the given ID and key.
func (s *metaStore) GetProperty(_ uint64, key string) ([]byte, error) {
	aspectID := s.ctx.AspectID
	codeStore := s.NewPrefixStore(V0AspectPropertyKeyPrefix)
	if _, ok := reservedPropertyKeys[key]; ok {
		// for reserved key, we return empty value
		return nil, nil
	}

	aspectPropertyKey := AspectPropertyKey(
		aspectID.Bytes(),
		[]byte(key),
	)

	return s.Load(codeStore, aspectPropertyKey)
}

// BumpVersion bumps the version of the aspect with the given ID.
func (s *metaStore) BumpVersion() (v uint64, err error) {
	aspectID := s.ctx.AspectID
	version, err := s.GetLatestVersion()
	if err != nil {
		s.ctx.Logger().Error("failed to get latest version", "aspect", aspectID.Hex(), "err", err)
		return 0, err
	}

	newVersionU64 := version + 1
	newVersion := uint256.NewInt(newVersionU64)
	prefixStore := s.NewPrefixStore(V0AspectCodeVersionKeyPrefix)
	storeKey := AspectIDKey(aspectID.Bytes())

	defer func() {
		if err == nil {
			s.latestVersionCache = newVersionU64
		}
	}()

	return newVersionU64, s.Store(prefixStore, storeKey, newVersion.Bytes())
}

// StoreVersionMeta stores the meta of the aspect with the given ID and version.
func (s *metaStore) StoreVersionMeta(version uint64, meta *types.VersionMeta) error {
	aspectID := s.ctx.AspectID
	// Store join point
	if meta.JoinPoint > 0 {
		u256Version := uint256.NewInt(version)
		joinPointBig := big.NewInt(0).SetUint64(meta.JoinPoint)
		if _, ok := artelasdkType.CheckIsJoinPoint(joinPointBig); !ok {
			joinPointBig.SetUint64(0)
		}

		prefixStore := s.NewPrefixStore(V0AspectJoinPointRunKeyPrefix)
		storeKey := AspectArrayKey(aspectID.Bytes(), u256Version.Bytes(), []byte(V0AspectRunJoinPointKey))
		if err := s.Store(prefixStore, storeKey, joinPointBig.Bytes()); err != nil {
			return err
		}
	}

	// for v0 pay master / proof is stored with property
	// code hash is not stored
	return nil
}

// StoreCode stores the code of the aspect with the given ID and version.
func (s *metaStore) StoreCode(version uint64, code []byte) error {
	aspectID := s.ctx.AspectID
	prefixStore := s.NewPrefixStore(V0AspectCodeKeyPrefix)
	storeKey := AspectVersionKey(aspectID.Bytes(), uint256.NewInt(version).Bytes())
	return s.Store(prefixStore, storeKey, code)
}

// StoreProperties stores the properties of the aspect with the given ID.
func (s *metaStore) StoreProperties(_ uint64, properties []types.Property) error {
	if len(properties) == 0 {
		return nil
	}

	// check reserved keys
	for _, prop := range properties {
		if _, ok := reservedPropertyKeys[prop.Key]; ok {
			return store.ErrPropertyReserved
		}
	}

	return s.storeProperties(properties)
}

func (s *metaStore) storeProperties(properties []types.Property) error {
	aspectID := s.ctx.AspectID
	prefixStore := s.NewPrefixStore(V0AspectPropertyKeyPrefix)
	propertyKeysKey := AspectPropertyKey(aspectID.Bytes(), []byte(V0AspectPropertyAllKeyPrefix))

	propertyAllKey, err := s.Load(prefixStore, propertyKeysKey)
	if err != nil {
		return err
	}

	keySet := treeset.NewWithStringComparator()
	// add propertyAllKey to keySet
	if len(propertyAllKey) > 0 {
		splitData := strings.Split(string(propertyAllKey), V0AspectPropertyAllKeySplit)
		for _, datum := range splitData {
			keySet.Add(datum)
		}
	}

	for i := range properties {
		// add key and deduplicate
		keySet.Add(properties[i].Key)
	}

	// check key limit
	if keySet.Size() > types.AspectPropertyLimit {
		return store.ErrTooManyProperties
	}

	// Store property key
	for i := range properties {
		key := properties[i].Key
		value := properties[i].Value

		// Store
		aspectPropertyKey := AspectPropertyKey(
			aspectID.Bytes(),
			[]byte(key),
		)

		if err := s.Store(prefixStore, aspectPropertyKey, value); err != nil {
			s.ctx.Logger().Error("failed to Store aspect property", "aspect", aspectID.Hex(), "key", key, "err", err)
			return err
		}

		s.ctx.Logger().Debug("aspect property updated", "aspect", aspectID.Hex(), "key", key)
	}

	// Store AspectPropertyAllKey
	keyAry := make([]string, keySet.Size())
	for i, key := range keySet.Values() {
		keyAry[i] = key.(string)
	}
	allKeys := strings.Join(keyAry, V0AspectPropertyAllKeySplit)
	allKeysKey := AspectPropertyKey(
		aspectID.Bytes(),
		[]byte(V0AspectPropertyAllKeyPrefix),
	)
	if err := s.Store(prefixStore, allKeysKey, []byte(allKeys)); err != nil {
		return err
	}

	return nil
}

// LoadAspectBoundAccounts loads all accounts bound to the given aspect.
func (s *metaStore) LoadAspectBoundAccounts() ([]types.Binding, error) {
	prefixStore := s.NewPrefixStore(V0AspectRefKeyPrefix)
	set, err := s.loadAspectRef(prefixStore)
	if err != nil {
		return nil, err
	}

	bindings := make([]types.Binding, 0)
	if set != nil {
		for _, data := range set.Values() {
			contractAddr := common.HexToAddress(data.(string))
			bindings = append(bindings, types.Binding{
				Account: contractAddr,
			})
		}
	}

	return bindings, nil
}

func (s *metaStore) saveAspectRef(prefixStore prefix.Store, aspectID common.Address, account common.Address) error {
	storeKey := AspectIDKey(aspectID.Bytes())
	set, err := s.loadAspectRef(prefixStore)
	if err != nil {
		return err
	}

	accountHex := account.Hex()
	if set.Contains(accountHex) {
		// already exist
		return store.ErrAlreadyBound
	}

	set.Add(accountHex)
	rawJSON, err := set.MarshalJSON()
	if err != nil {
		return err
	}

	return s.Store(prefixStore, storeKey, rawJSON)
}

func (s *metaStore) removeAspectRef(prefixStore prefix.Store, account common.Address) error {
	aspectID := s.ctx.AspectID
	storeKey := AspectIDKey(aspectID.Bytes())
	set, err := s.loadAspectRef(prefixStore)
	if err != nil {
		return err
	}

	accountHex := account.Hex()
	if !set.Contains(accountHex) {
		// not exist
		return nil
	}

	set.Remove(accountHex)
	rawJSON, err := set.MarshalJSON()
	if err != nil {
		return err
	}

	return s.Store(prefixStore, storeKey, rawJSON)
}

func (s *metaStore) loadAspectRef(prefixStore prefix.Store) (*treeset.Set, error) {
	aspectID := s.ctx.AspectID
	storeKey := AspectIDKey(aspectID.Bytes())

	rawTree, err := s.Load(prefixStore, storeKey)
	if err != nil {
		return nil, err
	}

	if rawTree == nil {
		return treeset.NewWithStringComparator(), nil
	}

	set := treeset.NewWithStringComparator()
	if err := set.UnmarshalJSON(rawTree); err != nil {
		return nil, err
	}

	return set, nil
}
