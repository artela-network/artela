package evm

import (
	"github.com/artela-network/artela/x/evm/states"
	"math/big"

	"github.com/artela-network/evm/vm"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"

	evmtypes "github.com/artela-network/artela/x/evm/txs/support"
	feemodule "github.com/artela-network/artela/x/fee/types"
)

// EVMKeeper defines the expected keeper interface used on the AnteHandler
type EVMKeeper interface { //nolint: revive
	states.Keeper
	DynamicFeeEVMKeeper

	NewEVM(ctx cosmos.Context, msg *core.Message, cfg *states.EVMConfig, tracer vm.EVMLogger, stateDB vm.StateDB) *vm.EVM
	DeductTxCostsFromUserBalance(ctx cosmos.Context, fees cosmos.Coins, from common.Address) error
	GetBalance(ctx cosmos.Context, addr common.Address) *big.Int
	ResetTransientGasUsed(ctx cosmos.Context)
	GetTxIndexTransient(ctx cosmos.Context) uint64
	GetParams(ctx cosmos.Context) evmtypes.Params
}

type FeeKeeper interface {
	GetParams(ctx cosmos.Context) (params feemodule.Params)
	AddTransientGasWanted(ctx cosmos.Context, gasWanted uint64) (uint64, error)
	GetBaseFeeEnabled(ctx cosmos.Context) bool
}

// DynamicFeeEVMKeeper is a subset of EVMKeeper interface that supports dynamic fee checker
type DynamicFeeEVMKeeper interface {
	ChainID() *big.Int
	GetParams(ctx cosmos.Context) evmtypes.Params
	GetBaseFee(ctx cosmos.Context, ethCfg *params.ChainConfig) *big.Int
}

type protoTxProvider interface {
	GetProtoTx() *tx.Tx
}
