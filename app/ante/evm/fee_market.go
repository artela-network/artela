package evm

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/artela-network/artela/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// GasWantedDecorator keeps track of the gasWanted amount on the current block in transient store
// for BaseFee calculation.
// NOTE: This decorator does not perform any validation
type GasWantedDecorator struct {
	evmKeeper EVMKeeper
	feeKeeper FeeKeeper
}

// NewGasWantedDecorator creates a new NewGasWantedDecorator
func NewGasWantedDecorator(
	evmKeeper EVMKeeper,
	feeKeeper FeeKeeper,
) GasWantedDecorator {
	return GasWantedDecorator{
		evmKeeper,
		feeKeeper,
	}
}

func (gwd GasWantedDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	evmParams := gwd.evmKeeper.GetParams(ctx)
	chainCfg := evmParams.GetChainConfig()
	ethCfg := chainCfg.EthereumConfig(gwd.evmKeeper.ChainID())

	blockHeight := big.NewInt(ctx.BlockHeight())
	isLondon := ethCfg.IsLondon(blockHeight)

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok || !isLondon {
		return next(ctx, tx, simulate)
	}

	gasWanted := feeTx.GetGas()
	// return error if the tx gas is greater than the block limit (max gas)
	blockGasLimit := types.BlockGasLimit(ctx)
	if gasWanted > blockGasLimit {
		return ctx, errorsmod.Wrapf(
			errortypes.ErrOutOfGas,
			"tx gas (%d) exceeds block gas limit (%d)",
			gasWanted,
			blockGasLimit,
		)
	}

	isBaseFeeEnabled := gwd.feeKeeper.GetBaseFeeEnabled(ctx)

	// Add total gasWanted to cumulative in block transientStore in FeeMarket module
	if isBaseFeeEnabled {
		if _, err := gwd.feeKeeper.AddTransientGasWanted(ctx, gasWanted); err != nil {
			return ctx, errorsmod.Wrapf(err, "failed to add gas wanted to transient store")
		}
	}

	return next(ctx, tx, simulate)
}
