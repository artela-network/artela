package erc20

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/artela-network/artela-evm/vm"
	artelatypes "github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/artela/x/evm/precompile/erc20/proxy"
	"github.com/artela-network/artela/x/evm/precompile/erc20/types"
	evmtypes "github.com/artela-network/artela/x/evm/types"
	"github.com/cometbft/cometbft/libs/log"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ vm.PrecompiledContract = (*ERC20Contract)(nil)

	// global precompiled contracts
	GlobalERC20Contract *ERC20Contract
)

type ERC20Contract struct {
	logger   log.Logger
	storeKey storetypes.StoreKey

	bankKeeper evmtypes.BankKeeper
}

func InitERC20Contract(logger log.Logger, storeKey storetypes.StoreKey, bankKeeper evmtypes.BankKeeper) {
	GlobalERC20Contract = &ERC20Contract{
		logger:     logger,
		storeKey:   storeKey,
		bankKeeper: bankKeeper,
	}
}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *ERC20Contract) RequiredGas(input []byte) uint64 {
	return 21000
}

func (c *ERC20Contract) Run(ctx context.Context, input []byte) ([]byte, error) {
	var sdkCtx sdk.Context
	if aspectCtx, ok := ctx.(*artelatypes.AspectRuntimeContext); !ok {
		return nil, errors.New("failed to unwrap AspectRuntimeContext from context.Context")
	} else {
		sdkCtx = aspectCtx.CosmosContext()
	}

	if len(input) < 4 {
		return nil, errors.New("invalid input")
	}

	parsedABI, err := abi.JSON(strings.NewReader(proxy.ERC20ProxyAbi))
	if err != nil {
		return nil, err
	}

	var (
		methodID  = input[:4]
		inputData = input[4:]
	)

	method, err := parsedABI.MethodById(methodID)
	if err != nil {
		return nil, err
	}

	args := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(args, inputData); err != nil {
		return nil, err
	}

	var caller common.Address // TODO
	switch method.Name {
	case types.Method_BalanceOf:
		return c.handleBalanceOf(sdkCtx, input[4:], caller)
	case types.Method_Transfer:
		return c.handleTransfer(sdkCtx, input[4:], caller)
	default:
		return nil, errors.New("unknown method")
	}
}

func (c *ERC20Contract) handleBalanceOf(ctx sdk.Context, input []byte, caller common.Address) ([]byte, error) {
	if len(input) != 32 {
		return nil, errors.New("invalid input length")
	}
	address := common.BytesToAddress(input[12:])
	_ = address
	var balance *big.Int // TODO
	if balance == nil {
		balance = big.NewInt(0)
	}
	return balance.FillBytes(make([]byte, 32)), nil
}

func (c *ERC20Contract) handleTransfer(ctx sdk.Context, input []byte, caller common.Address) ([]byte, error) {
	if len(input) != 64 {
		return nil, errors.New("invalid input length")
	}

	to := common.BytesToAddress(input[12:32])
	amount := new(big.Int).SetBytes(input[32:64])

	// TODO
	_ = to
	_ = amount

	return nil, nil
}
