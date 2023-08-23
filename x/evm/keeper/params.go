package keeper

import (
	"github.com/artela-network/artela/x/evm/process"
	"github.com/artela-network/artela/x/evm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params process.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixParams)
	if len(bz) == 0 {
		return k.GetLegacyParams(ctx)
	}
	k.cdc.MustUnmarshal(bz, &params)
	return
}

// SetParams sets the EVM params each in their individual key for better get performance
func (k Keeper) SetParams(ctx sdk.Context, params process.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}

	store.Set(types.KeyPrefixParams, bz)
	return nil
}

// GetLegacyParams returns param set for version before migrate
func (k Keeper) GetLegacyParams(ctx sdk.Context) process.Params {
	var params process.Params
	k.ss.GetParamSetIfExists(ctx, &params)
	return params
}
