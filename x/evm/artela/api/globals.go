package api

import (
	"math/big"

	cosmos "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"

	"github.com/artela-network/artela-evm/vm"
	"github.com/artela-network/artela/x/evm/states"
)

// globals, we need to manage the lifecycle carefully for the following hooks / variables
var (
	evmKeeper EVMKeeper
)

type EVMKeeper interface {
	NewEVM(
		ctx cosmos.Context,
		msg *core.Message,
		cfg *states.EVMConfig,
		tracer vm.EVMLogger,
		stateDB vm.StateDB,
	) *vm.EVM
	EVMConfig(ctx cosmos.Context, proposerAddress cosmos.ConsAddress, chainID *big.Int) (*states.EVMConfig, error)
	ChainID() *big.Int
	GetAccount(ctx cosmos.Context, addr common.Address) *states.StateAccount
	GetState(ctx cosmos.Context, addr common.Address, key common.Hash) common.Hash
	GetCode(ctx cosmos.Context, codeHash common.Hash) []byte
	ForEachStorage(ctx cosmos.Context, addr common.Address, cb func(key, value common.Hash) bool)

	SetAccount(ctx cosmos.Context, addr common.Address, account states.StateAccount) error
	SetState(ctx cosmos.Context, addr common.Address, key common.Hash, value []byte)
	SetCode(ctx cosmos.Context, codeHash []byte, code []byte)
	DeleteAccount(ctx cosmos.Context, addr common.Address) error
}

func InitAspectGlobals(keeper EVMKeeper) {
	evmKeeper = keeper
}
