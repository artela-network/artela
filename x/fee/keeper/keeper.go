package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	paramsmodule "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/artela-network/artela/x/fee/types"
)

// Keeper grants access to the Fee Market module states.
type Keeper struct {
	// Protobuf codec
	cdc codec.BinaryCodec
	// Store key required for the Fee Market Prefix KVStore.
	storeKey     storetypes.StoreKey
	transientKey storetypes.StoreKey
	// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
	authority cosmos.AccAddress
	// Legacy subspace
	ss paramsmodule.Subspace
}

// NewKeeper generates new fee market module keeper
func NewKeeper(
	cdc codec.BinaryCodec, authority cosmos.AccAddress, storeKey, transientKey storetypes.StoreKey, ss paramsmodule.Subspace,
) *Keeper {
	// ensure authority account is correctly formatted
	if err := cosmos.VerifyAddressFormat(authority); err != nil {
		panic(err)
	}

	return &Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		authority:    authority,
		transientKey: transientKey,
		ss:           ss,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx cosmos.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// ----------------------------------------------------------------------------
// Parent Block Gas Used
// Required by EIP1559 base fee calculation.
// ----------------------------------------------------------------------------

// SetBlockGasWanted sets the block gas wanted to the store.
// CONTRACT: this should be only called during EndBlock.
func (k Keeper) SetBlockGasWanted(ctx cosmos.Context, gas uint64) {
	store := ctx.KVStore(k.storeKey)
	gasBz := cosmos.Uint64ToBigEndian(gas)
	store.Set(types.KeyPrefixBlockGasWanted, gasBz)

	k.Logger(ctx).Debug("setState: SetBlockGasWanted",
		"key", "KeyPrefixBlockGasWanted",
		"gas", fmt.Sprintf("%d", gas))
}

// GetBlockGasWanted returns the last block gas wanted value from the store.
func (k Keeper) GetBlockGasWanted(ctx cosmos.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixBlockGasWanted)
	if len(bz) == 0 {
		return 0
	}

	return cosmos.BigEndianToUint64(bz)
}

// GetTransientGasWanted returns the gas wanted in the current block from transient store.
func (k Keeper) GetTransientGasWanted(ctx cosmos.Context) uint64 {
	store := ctx.TransientStore(k.transientKey)
	bz := store.Get(types.KeyPrefixTransientBlockGasWanted)
	if len(bz) == 0 {
		return 0
	}
	return cosmos.BigEndianToUint64(bz)
}

// SetTransientBlockGasWanted sets the block gas wanted to the transient store.
func (k Keeper) SetTransientBlockGasWanted(ctx cosmos.Context, gasWanted uint64) {
	store := ctx.TransientStore(k.transientKey)
	gasBz := cosmos.Uint64ToBigEndian(gasWanted)
	store.Set(types.KeyPrefixTransientBlockGasWanted, gasBz)

	k.Logger(ctx).Debug("setState: SetTransientBlockGasWanted",
		"key", "KeyPrefixTransientBlockGasWanted",
		"gasWanted", fmt.Sprintf("%d", gasWanted))
}

// AddTransientGasWanted adds the cumulative gas wanted in the transient store
func (k Keeper) AddTransientGasWanted(ctx cosmos.Context, gasWanted uint64) (uint64, error) {
	result := k.GetTransientGasWanted(ctx) + gasWanted
	k.SetTransientBlockGasWanted(ctx, result)
	return result, nil
}
