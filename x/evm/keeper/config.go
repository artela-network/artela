package keeper

import (
	"github.com/artela-network/artela/x/evm/txs/support"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/artela-network/artela/x/evm/states"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
)

// EVMConfig creates the EVMConfig based on current states
func (k *Keeper) EVMConfig(ctx sdk.Context, proposerAddress sdk.ConsAddress, chainID *big.Int) (*states.EVMConfig, error) {
	params := k.GetParams(ctx)
	ethCfg := params.ChainConfig.EthereumConfig(chainID)

	// get the coinbase address from the block proposer
	coinbase, err := k.GetProposerAddress(ctx, proposerAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to obtain coinbase address")
	}

	baseFee := k.GetBaseFee(ctx, ethCfg)
	return &states.EVMConfig{
		Params:      params,
		ChainConfig: ethCfg,
		CoinBase:    coinbase,
		BaseFee:     baseFee,
	}, nil
}

// TxConfig loads `TxConfig` from current transient storage
func (k *Keeper) TxConfig(ctx sdk.Context, txHash common.Hash) states.TxConfig {
	return states.NewTxConfig(
		common.BytesToHash(ctx.HeaderHash()), // BlockHash
		txHash,                               // TxHash
		uint(k.GetTxIndexTransient(ctx)),     // TxIndex
		uint(k.GetLogSizeTransient(ctx)),     // LogIndex
	)
}

// VMConfig creates an EVM configuration from the debug setting and the extra EIPs enabled on the
// module parameters. The config support uses the default JumpTable from the EVM.
func (k Keeper) VMConfig(ctx sdk.Context, _ core.Message, cfg *states.EVMConfig, tracer vm.EVMLogger) vm.Config {
	noBaseFee := true
	if support.IsLondon(cfg.ChainConfig, ctx.BlockHeight()) {
		noBaseFee = k.feeKeeper.GetParams(ctx).NoBaseFee
	}

	// var debug bool
	// if _, ok := tracer.(types.NoOpTracer); !ok {
	// 	debug = true
	// }

	return vm.Config{
		// Debug:     debug,
		Tracer:    tracer,
		NoBaseFee: noBaseFee,
		ExtraEips: cfg.Params.EIPs(),
	}
}
