package erc20

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"cosmossdk.io/math"
	"github.com/artela-network/artela-evm/vm"
	artelatypes "github.com/artela-network/artela/x/evm/artela/types"
	precompiled "github.com/artela-network/artela/x/evm/precompile"
	"github.com/artela-network/artela/x/evm/precompile/erc20/proxy"
	"github.com/artela-network/artela/x/evm/precompile/erc20/types"
	evmtypes "github.com/artela-network/artela/x/evm/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ vm.PrecompiledContract = (*ERC20Contract)(nil)
)

type APIMethod func(sdk.Context, common.Address, common.Address, map[string]interface{}) ([]byte, error)

type ERC20Contract struct {
	logger   log.Logger
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	tokenPairs types.TokenPairs // TODO cache the token pairs
	bankKeeper evmtypes.BankKeeper
	methods    map[string]APIMethod
}

func InitERC20Contract(logger log.Logger, cdc codec.BinaryCodec, storeKey storetypes.StoreKey, bankKeeper evmtypes.BankKeeper) *ERC20Contract {
	contract := &ERC20Contract{
		logger:     logger,
		cdc:        cdc,
		storeKey:   storeKey,
		bankKeeper: bankKeeper,
		methods:    make(map[string]APIMethod),
	}

	contract.methods[types.Method_BalanceOf] = contract.handleBalanceOf
	contract.methods[types.Method_Register] = contract.handleRegister
	contract.methods[types.Method_Transfer] = contract.handleTransfer

	precompiled.RegisterPrecompiles(types.PrecompiledAddress, contract)
	return contract
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

	// get tx.from, which is the proxy contract address
	caller, ok := sdkCtx.Value("msgFrom").(common.Address)
	if !ok {
		return nil, errors.New("from address not valiad")
	}

	// get tx.To, which is the proxy contract address
	msgTo, ok := sdkCtx.Value("msgTo").(common.Address)
	if !ok {
		return nil, errors.New("to address not valiad")
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

	fn, ok := c.methods[method.Name]
	if !ok {
		return nil, errors.New("unknown method")
	}

	args := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(args, inputData); err != nil {
		return nil, err
	}

	return fn(sdkCtx, msgTo, caller, args)
}

func (c *ERC20Contract) handleRegister(ctx sdk.Context, proxy common.Address, _ common.Address, args map[string]interface{}) ([]byte, error) {
	if len(args) != 1 {
		return types.False32Byte, errors.New("invalid input")
	}

	denom, ok := args["denom"].(string)
	if !ok || len(denom) == 0 {
		return types.False32Byte, errors.New("invalid input denom")
	}

	if d := c.GetDenomByProxy(ctx, proxy); len(d) > 0 {
		return types.False32Byte, errors.New("proxy has been registered")
	}

	if err := c.registerNewTokenPairs(ctx, denom, proxy); err != nil {
		return types.False32Byte, err
	}
	return types.True32Byte, nil
}

func (c *ERC20Contract) handleBalanceOf(ctx sdk.Context, proxy common.Address, _ common.Address, args map[string]interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("invalid input")
	}

	denom, err := c.getDenom(ctx, proxy)
	if err != nil {
		return types.False32Byte, err
	}

	addr, ok := args["account"].(common.Address)
	if !ok {
		return nil, errors.New("invalid input account")
	}

	accAddr := sdk.AccAddress(addr.Bytes())

	coin := c.bankKeeper.GetBalance(ctx, accAddr, denom)
	balance := coin.Amount.BigInt()
	if balance == nil {
		balance = big.NewInt(0)
	}
	return balance.FillBytes(make([]byte, 32)), nil
}

func (c *ERC20Contract) handleTransfer(ctx sdk.Context, proxy common.Address, caller common.Address, args map[string]interface{}) ([]byte, error) {
	if len(args) != 2 {
		return types.False32Byte, errors.New("invalid input")
	}

	denom, err := c.getDenom(ctx, proxy)
	if err != nil {
		return types.False32Byte, err
	}

	fromAccAddr := sdk.AccAddress(caller.Bytes())

	to, ok := args["to"].(common.Address)
	if !ok {
		return types.False32Byte, errors.New("invalid input address")
	}
	toAccAddr := sdk.AccAddress(to.Bytes())

	amount, ok := args["amount"].(*big.Int)
	if !ok {
		return types.False32Byte, errors.New("invalid input amount")
	}

	coins := sdk.NewCoins(sdk.NewCoin(denom, math.NewIntFromBigInt(amount)))
	if err := c.bankKeeper.IsSendEnabledCoins(ctx, coins...); err != nil {
		return types.False32Byte, err
	}

	err = c.bankKeeper.SendCoins(
		ctx, fromAccAddr, toAccAddr, coins)
	if err != nil {
		return types.False32Byte, err
	}

	return types.True32Byte, nil
}

func (c *ERC20Contract) getDenom(ctx sdk.Context, proxy common.Address) (string, error) {
	// get registered denom for the proxy address
	denom := c.GetDenomByProxy(ctx, proxy)
	if len(denom) == 0 {
		return "", errors.New("no registered coin found")
	}

	return denom, nil
}
