package fee

import (
	"fmt"

	"github.com/artela-network/artela/x/fee/keeper"

	"github.com/artela-network/artela/x/fee/types"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	cosmos "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlock updates base fee
func BeginBlock(ctx cosmos.Context, k *keeper.Keeper) {
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
	ctx.EventManager().EmitEvents(cosmos.Events{
		cosmos.NewEvent(
			types.EventTypeFee,
			cosmos.NewAttribute(types.AttributeKeyBaseFee, baseFee.String()),
		),
	})
}

// EndBlock update block gas wanted.
// The EVM end block logic doesn't update the validator set, thus it returns
// an empty slice.
func EndBlock(ctx cosmos.Context, k *keeper.Keeper) {
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
	minGasMultiplier := k.GetParams(ctx).MinGasMultiplier
	limitedGasWanted := math.LegacyNewDec(gasWanted.Int64()).Mul(minGasMultiplier)
	updatedGasWanted := math.LegacyMaxDec(limitedGasWanted, math.LegacyNewDec(gasUsed.Int64())).TruncateInt().Uint64()
	k.SetBlockGasWanted(ctx, updatedGasWanted)

	defer func() {
		telemetry.SetGauge(float32(updatedGasWanted), "fee", "block_gas")
	}()

	ctx.EventManager().EmitEvent(cosmos.NewEvent(
		"block_gas",
		cosmos.NewAttribute("height", fmt.Sprintf("%d", ctx.BlockHeight())),
		cosmos.NewAttribute("amount", fmt.Sprintf("%d", updatedGasWanted)),
	))
}
