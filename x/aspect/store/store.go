package store

import (
	"encoding/binary"
	"encoding/hex"
	aspectmoduletypes "github.com/artela-network/artela/x/aspect/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/ethereum/go-ethereum/common"
)

const (
	ProtocolVersionLen = 2
	ProtocolInfoLen    = 4
)

type ProtocolVersion uint16

func (p ProtocolVersion) MarshalText() ([]byte, error) {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(p))
	return bytes, nil
}

func (p *ProtocolVersion) UnmarshalText(text []byte) error {
	if len(text) < 2 {
		*p = 0
	} else {
		*p = ProtocolVersion(binary.BigEndian.Uint16(text[:2]))
	}

	return nil
}

func (p *ProtocolVersion) Offset() uint64 {
	if *p == 0 {
		return 0
	} else {
		return ProtocolVersionLen
	}
}

type AspectInfo struct {
	MetaVersion  ProtocolVersion
	StateVersion ProtocolVersion

	offset uint64
}

func (a *AspectInfo) Offset() uint64 {
	return a.offset
}

func (a AspectInfo) MarshalText() ([]byte, error) {
	bytes := make([]byte, ProtocolInfoLen)
	// next 2 bytes saves meta version
	binary.BigEndian.PutUint16(bytes[0:2], uint16(a.MetaVersion))
	// next 2 bytes saves state version
	binary.BigEndian.PutUint16(bytes[2:4], uint16(a.StateVersion))
	return bytes, nil
}

func (a *AspectInfo) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		// v0 store does not have aspect info, so we just return nil here
		// so that all versions of aspect info is 0, which is compatible with v0 store
		a.offset = 0
		return nil
	}

	if len(text) < ProtocolInfoLen {
		return ErrInvalidProtocolInfo
	}

	a.MetaVersion = ProtocolVersion(binary.BigEndian.Uint16(text[0:2]))
	a.StateVersion = ProtocolVersion(binary.BigEndian.Uint16(text[2:4]))
	a.offset = ProtocolInfoLen

	return nil
}

type (
	AspectStateStoreConstructor = func(ctx *aspectmoduletypes.AspectStoreContext) AspectStateStore
	AspectMetaStoreConstructor  = func(ctx *aspectmoduletypes.AspectStoreContext, protocolExtension []byte) AspectMetaStore
	AccountStoreConstructor     = func(ctx *aspectmoduletypes.AccountStoreContext) AccountStore
)

// aspect initializer registry
var (
	aspectStateStoreRegistry = make(map[ProtocolVersion]AspectStateStoreConstructor)
	aspectMetaStoreRegistry  = make(map[ProtocolVersion]AspectMetaStoreConstructor)
	accountStoreRegistry     = make(map[ProtocolVersion]AccountStoreConstructor)
)

func RegisterAspectStateStore(version ProtocolVersion, constructor AspectStateStoreConstructor) {
	aspectStateStoreRegistry[version] = constructor
}

func RegisterAspectMetaStore(version ProtocolVersion, constructor AspectMetaStoreConstructor) {
	aspectMetaStoreRegistry[version] = constructor
}

func RegisterAccountStore(version ProtocolVersion, constructor AccountStoreConstructor) {
	accountStoreRegistry[version] = constructor
}

type GasMeteredStore interface {
	// Gas returns the gas remains in the store
	Gas() uint64
	// TransferGasFrom transfers the gas from another store
	TransferGasFrom(store GasMeteredStore)
}

// AccountStore is the store for each account that using aspect
type AccountStore interface {
	// LoadAccountBoundAspects returns the aspects bound to the account,
	LoadAccountBoundAspects(filter aspectmoduletypes.BindingFilter) ([]aspectmoduletypes.Binding, error)
	// StoreBinding adds the binding of the given aspect to the account
	StoreBinding(aspectID common.Address, version uint64, joinPoint uint64, priority int8, isCA bool) error
	// RemoveBinding removes the binding of the given aspect from the account
	RemoveBinding(aspectID common.Address, joinPoint uint64, isCA bool) error

	// Used returns true if this given version of store has been used before
	Used() (bool, error)
	// MigrateFrom migrates the data from the old store to the new store
	MigrateFrom(old AccountStore) error
	// Init initializes the store
	Init() error
	// Version returns the version of the store
	Version() ProtocolVersion

	GasMeteredStore
}

// AspectStateStore is the store for aspect state related info
type AspectStateStore interface {
	// GetState returns the value for the given key
	GetState(key []byte) []byte
	// SetState sets the value for the given key
	SetState(key []byte, value []byte)
	// Version returns the version of the store
	Version() ProtocolVersion
}

// AspectMetaStore is the store for aspect metadata
type AspectMetaStore interface {
	// GetCode returns the code for the given version
	GetCode(version uint64) ([]byte, error)
	// GetVersionMeta returns the meta for the given version
	GetVersionMeta(version uint64) (*aspectmoduletypes.VersionMeta, error)
	// GetMeta returns the meta for the aspect
	GetMeta() (*aspectmoduletypes.AspectMeta, error)
	// GetLatestVersion returns the latest version of the aspect
	GetLatestVersion() (uint64, error)
	// GetProperty returns the properties for the given version
	GetProperty(version uint64, key string) ([]byte, error)
	// LoadAspectBoundAccounts returns the accounts bound to the aspect
	LoadAspectBoundAccounts() ([]aspectmoduletypes.Binding, error)

	// BumpVersion bumps the version of the aspect
	BumpVersion() (uint64, error)
	// StoreVersionMeta stores the meta for the given version
	StoreVersionMeta(version uint64, meta *aspectmoduletypes.VersionMeta) error
	// StoreMeta stores the meta for the aspect
	StoreMeta(meta *aspectmoduletypes.AspectMeta) error
	// StoreCode stores the code for the given version
	StoreCode(version uint64, code []byte) error
	// StoreProperties stores the properties for the given version
	StoreProperties(version uint64, properties []aspectmoduletypes.Property) error
	// StoreBinding stores the binding for the given account
	StoreBinding(account common.Address, version uint64, joinPoint uint64, priority int8) error
	// RemoveBinding removes the binding for the given account
	RemoveBinding(account common.Address) error

	// Version returns the version of the store
	Version() ProtocolVersion
	// MigrateFrom migrates the data from the old store to the new store
	MigrateFrom(old AspectMetaStore) error
	// Used returns true if this given version of store has been used before
	Used() (bool, error)
	// Init initializes the store
	Init() error

	GasMeteredStore
}

// loadProtocolInfo loads the protocol info for the given address
func loadProtocolInfo(ctx aspectmoduletypes.StoreContext) (ProtocolVersion, []byte, error) {
	aspectKV := ctx.CosmosContext().KVStore(ctx.AspectStoreKey())
	store := prefix.NewStore(aspectKV, AspectProtocolInfoKeyPrefix)

	var address common.Address
	switch ctx := ctx.(type) {
	case *aspectmoduletypes.AccountStoreContext:
		address = ctx.Account
	case *aspectmoduletypes.AspectStoreContext:
		address = ctx.AspectID
	default:
		return 0, nil, ErrInvalidStoreContext
	}

	protoInfo := store.Get(address.Bytes())

	var protocolVersion ProtocolVersion
	if err := protocolVersion.UnmarshalText(protoInfo); err != nil {
		ctx.Logger().Error("unmarshal aspect protocol info failed",
			"address", address.Hex(), "data", hex.EncodeToString(protoInfo))
		return 0, protoInfo, err
	}

	return protocolVersion, protoInfo[protocolVersion.Offset():], nil
}

// parseAspectInfo parses
func parseAspectInfo(raw []byte) (*AspectInfo, error) {
	aspectInfo := &AspectInfo{}
	if err := aspectInfo.UnmarshalText(raw); err != nil {
		return nil, err
	}

	return aspectInfo, nil
}

// GetAccountStore returns the account store for the given account,
// account store is used to store account related info like bound aspects.
// This function will return 2 stores, the current store and the new store.
// New store will be nil if no migration needed, otherwise it will be the instance of the new version store.
func GetAccountStore(ctx *aspectmoduletypes.AccountStoreContext) (current AccountStore, new AccountStore, err error) {
	// load protocol version
	protocolVersion, _, err := loadProtocolInfo(ctx)
	if err != nil {
		return nil, nil, err
	}

	// load protocol storage constructor
	constructor, ok := accountStoreRegistry[protocolVersion]
	if !ok {
		ctx.Logger().Error("unsupported protocol version", "version", protocolVersion)
		return nil, nil, ErrUnknownProtocolVersion
	}

	// if latest version is greater than the used version,
	// we also init an instance of the latest version of store to let the caller func
	// decides whether to migrate the data or not
	var latestStore AccountStore
	latestVersion := latestStoreVersion(accountStoreRegistry)
	if latestVersion > protocolVersion {
		// init the latest store version
		latestStore = accountStoreRegistry[latestVersion](ctx)
	}

	// if this protocol version is 0, we have 2 cases here:
	// 1. the account has never used aspect before
	// 2. the account was using aspect at protocol version 0, but not migrated to the new version
	// the following is just the special case for processing v0 store
	if protocolVersion == 0 {
		// init v0 store first
		v0Store := constructor(ctx)
		// so first we need to check whether this account has used aspect before
		if used, err := v0Store.Used(); err != nil {
			return nil, nil, err
		} else if used {
			// if v0 store used, we need to migrate the data to the latest version
			return v0Store, latestStore, nil
		} else {
			// otherwise, we just return the latest store
			if latestStore == nil {
				// set the latest store to v0 store if no migration needed
				latestStore = v0Store
			}

			return latestStore, nil, nil
		}
	}

	// build the aspect store
	return constructor(ctx), latestStore, nil
}

// GetAspectMetaStore returns the aspect meta store for the given aspect id
func GetAspectMetaStore(ctx *aspectmoduletypes.AspectStoreContext) (current AspectMetaStore, new AspectMetaStore, err error) {
	// load protocol version
	protocolVersion, rawAspectInfo, err := loadProtocolInfo(ctx)
	if err != nil {
		return nil, nil, err
	}

	return getAspectMetaStore(ctx, protocolVersion, rawAspectInfo)
}

func getAspectMetaStore(ctx *aspectmoduletypes.AspectStoreContext, protocolVersion ProtocolVersion, rawAspectInfo []byte) (current AspectMetaStore, new AspectMetaStore, err error) {
	// load aspect info if protocol version is not 0
	aspectInfo, err := parseAspectInfo(rawAspectInfo)
	if err != nil {
		ctx.Logger().Error("parse aspect info failed", "aspectId", ctx.AspectID.Hex())
		return nil, nil, err
	}
	protocolExtension := rawAspectInfo[aspectInfo.Offset():]
	metaVersion := aspectInfo.MetaVersion

	// load protocol storage constructor
	constructor, ok := aspectMetaStoreRegistry[metaVersion]
	if !ok {
		ctx.Logger().Error("unsupported meta version", "version", metaVersion)
		return nil, nil, ErrUnknownProtocolVersion
	}

	// if latest version is greater than the used version,
	// we also init an instance of the latest version of store to let the caller func
	// decides whether to migrate the data or not
	var latestStore AspectMetaStore
	latestVersion := latestStoreVersion(aspectMetaStoreRegistry)
	if latestVersion > protocolVersion {
		// init the latest store version
		latestStore = aspectMetaStoreRegistry[latestVersion](ctx, protocolExtension)
	}

	// if this protocol version is 0, we have 2 cases here:
	// 1. the aspect has not been deployed yet
	// 2. the aspect is deployed at protocol version 0, but not migrated to the new version
	if protocolVersion == 0 {
		// init v0 store first
		v0Store := constructor(ctx, protocolExtension)
		// so first we need to check whether this account has used aspect before
		if used, err := v0Store.Used(); err != nil {
			// check fail
			return nil, nil, err
		} else if used {
			// if v0 store used, we need to migrate the data to the latest version
			return v0Store, latestStore, nil
		} else {
			// otherwise, we just return the latest store
			if latestStore == nil {
				// set the latest store to v0 store if no migration needed
				latestStore = v0Store
			}

			return latestStore, nil, nil
		}
	}

	// build the aspect store
	return constructor(ctx, protocolExtension), latestStore, nil
}

func GetAspectStateStore(ctx *aspectmoduletypes.AspectStoreContext) (AspectStateStore, error) {
	// load protocol version
	_, rawAspectInfo, err := loadProtocolInfo(ctx)
	if err != nil {
		return nil, err
	}

	// aspect info must have been initialized before the state store initialized,
	// if rawAspectInfo is empty, we just go with init v0 store
	aspectInfo, err := parseAspectInfo(rawAspectInfo)
	if err != nil {
		ctx.Logger().Error("parse aspect info failed", "aspectId", ctx.AspectID.Hex())
		return nil, err
	}

	stateVersion := aspectInfo.StateVersion

	// load protocol state constructor
	constructor, ok := aspectStateStoreRegistry[stateVersion]
	if !ok {
		ctx.Logger().Error("unsupported state version", "version", stateVersion)
		return nil, ErrUnknownProtocolVersion
	}

	return constructor(ctx), nil
}

func latestStoreVersion[T any](registry map[ProtocolVersion]T) ProtocolVersion {
	return ProtocolVersion(len(registry) - 1)
}
