package keeper

import (
	"fmt"
	"math/big"

	cosmos "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"

	"github.com/artela-network/artela/x/fee/types"
)

// CalculateBaseFee calculates the base fee for the current block. This is only calculated once per
// block during BeginBlock. If the NoBaseFee parameter is enabled or below activation height, this function returns nil.
// NOTE: This code is inspired from the go-ethereum EIP1559 implementation and adapted to Cosmos SDK-based
// chains. For the canonical code refer to: https://github.com/ethereum/go-ethereum/blob/master/consensus/misc/eip1559.go
func (k Keeper) CalculateBaseFee(ctx cosmos.Context) *big.Int {
	params := k.GetParams(ctx)

	// Ignore the calculation if not enabled
	if !params.IsBaseFeeEnabled(ctx.BlockHeight()) {
		return nil
	}

	consParams := ctx.ConsensusParams()

	// If the current block is the first EIP-1559 block, return the base fee
	// defined in the parameters (DefaultBaseFee if it hasn't been changed by
	// governance).
	if ctx.BlockHeight() == params.EnableHeight {
		return params.BaseFee.BigInt()
	}

	// get the block gas used and the base fee values for the parent block.
	// NOTE: this is not the parent's base fee but the current block's base fee,
	// as it is retrieved from the transient store, which is committed to the
	// persistent KVStore after EndBlock (ABCI Commit).
	parentBaseFee := params.BaseFee.BigInt()
	if parentBaseFee == nil {
		return nil
	}

	parentGasUsed := k.GetBlockGasWanted(ctx)

	gasLimit := new(big.Int).SetUint64(math.MaxUint64)

	// NOTE: a MaxGas equal to -1 means that block gas is unlimited
	if consParams != nil && consParams.Block.MaxGas > -1 {
		gasLimit = big.NewInt(consParams.Block.MaxGas)
	}

	// CONTRACT: ElasticityMultiplier cannot be 0 as it's checked in the params
	// validation
	parentGasTargetBig := new(big.Int).Div(gasLimit, new(big.Int).SetUint64(uint64(params.ElasticityMultiplier)))
	if !parentGasTargetBig.IsUint64() {
		return nil
	}

	parentGasTarget := parentGasTargetBig.Uint64()
	baseFeeChangeDenominator := new(big.Int).SetUint64(uint64(params.BaseFeeChangeDenominator))

	// If the parent gasUsed is the same as the target, the baseFee remains
	// unchanged.
	if parentGasUsed == parentGasTarget {
		return new(big.Int).Set(parentBaseFee)
	}

	if parentGasUsed > parentGasTarget {
		// If the parent block used more gas than its target, the baseFee should
		// increase.
		gasUsedDelta := new(big.Int).SetUint64(parentGasUsed - parentGasTarget)
		x := new(big.Int).Mul(parentBaseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := math.BigMax(
			x.Div(y, baseFeeChangeDenominator),
			common.Big1,
		)

		return x.Add(parentBaseFee, baseFeeDelta)
	}

	// Otherwise if the parent block used less gas than its target, the baseFee
	// should decrease.
	gasUsedDelta := new(big.Int).SetUint64(parentGasTarget - parentGasUsed)
	x := new(big.Int).Mul(parentBaseFee, gasUsedDelta)
	y := x.Div(x, parentGasTargetBig)
	baseFeeDelta := x.Div(y, baseFeeChangeDenominator)

	// Set global min gas price as lower bound of the base fee, transactions below
	// the min gas price don't even reach the mempool.
	minGasPrice := params.MinGasPrice.TruncateInt().BigInt()
	return math.BigMax(x.Sub(parentBaseFee, baseFeeDelta), minGasPrice)
}

// GetParams returns the total set of fee market parameters.
func (k Keeper) GetParams(ctx cosmos.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if len(bz) == 0 {
		var p types.Params
		k.ss.GetParamSetIfExists(ctx, &p)
		return p
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the fee market params in a single key
func (k Keeper) SetParams(ctx cosmos.Context, params types.Params) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}

	store.Set(types.ParamsKey, bz)
	k.Logger(ctx).Debug("setState: SetBlockGasWanted",
		"key", string(types.ParamsKey),
		"params", fmt.Sprintf("%+v", params))

	return nil
}

// ----------------------------------------------------------------------------
// Parent Base Fee
// Required by EIP1559 base fee calculation.
// ----------------------------------------------------------------------------

// GetBaseFeeEnabled returns true if base fee is enabled
func (k Keeper) GetBaseFeeEnabled(ctx cosmos.Context) bool {
	params := k.GetParams(ctx)
	return !params.NoBaseFee && ctx.BlockHeight() >= params.EnableHeight
}

// GetBaseFee gets the base fee from the store
func (k Keeper) GetBaseFee(ctx cosmos.Context) *big.Int {
	params := k.GetParams(ctx)
	if params.NoBaseFee {
		return nil
	}

	baseFee := params.BaseFee.BigInt()
	return baseFee
}

// SetBaseFee set's the base fee in the store
func (k Keeper) SetBaseFee(ctx cosmos.Context, baseFee *big.Int) {
	params := k.GetParams(ctx)
	params.BaseFee = cosmos.NewIntFromBigInt(baseFee)
	err := k.SetParams(ctx, params)
	if err != nil {
		return
	}
}
