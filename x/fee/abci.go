package fee

import (
	"fmt"
	"github.com/artela-network/artela/x/fee/keeper"

	"github.com/artela-network/artela/x/fee/types"
	abci "github.com/cometbft/cometbft/abci/types"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlock updates base fee
func BeginBlock(ctx sdk.Context, k *keeper.Keeper, _ abci.RequestBeginBlock) {
	baseFee := k.CalculateBaseFee(ctx)

	// return immediately if base fee is nil
	if baseFee == nil {
		return
	}

	k.SetBaseFee(ctx, baseFee)

	defer func() {
		telemetry.SetGauge(float32(baseFee.Int64()), "fee", "base_fee")
	}()

	// Store current base fee in event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFee,
			sdk.NewAttribute(types.AttributeKeyBaseFee, baseFee.String()),
		),
	})
}

// EndBlock update block gas wanted.
// The EVM end block logic doesn't update the validator set, thus it returns
// an empty slice.
func EndBlock(ctx sdk.Context, k *keeper.Keeper, _ abci.RequestEndBlock) {
	if ctx.BlockGasMeter() == nil {
		k.Logger(ctx).Error("block gas meter is nil when setting block gas wanted")
		return
	}

	gasWanted := sdkmath.NewIntFromUint64(k.GetTransientGasWanted(ctx))
	gasUsed := sdkmath.NewIntFromUint64(ctx.BlockGasMeter().GasConsumedToLimit())

	if !gasWanted.IsInt64() {
		k.Logger(ctx).Error("integer overflow by integer type conversion. Gas wanted > MaxInt64", "gas wanted", gasWanted.String())
		return
	}

	if !gasUsed.IsInt64() {
		k.Logger(ctx).Error("integer overflow by integer type conversion. Gas used > MaxInt64", "gas used", gasUsed.String())
		return
	}

	// to prevent BaseFee manipulation we limit the gasWanted so that
	// gasWanted = max(gasWanted * MinGasMultiplier, gasUsed)
	// this will be keep BaseFee protected from un-penalized manipulation
	// more info here https://github.com/evmos/ethermint/pull/1105#discussion_r888798925
	minGasMultiplier := k.GetParams(ctx).MinGasMultiplier
	limitedGasWanted := sdk.NewDec(gasWanted.Int64()).Mul(minGasMultiplier)
	updatedGasWanted := sdk.MaxDec(limitedGasWanted, sdk.NewDec(gasUsed.Int64())).TruncateInt().Uint64()
	k.SetBlockGasWanted(ctx, updatedGasWanted)

	defer func() {
		telemetry.SetGauge(float32(updatedGasWanted), "fee", "block_gas")
	}()

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"block_gas",
		sdk.NewAttribute("height", fmt.Sprintf("%d", ctx.BlockHeight())),
		sdk.NewAttribute("amount", fmt.Sprintf("%d", updatedGasWanted)),
	))
}
