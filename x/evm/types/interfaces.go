package types

import (
	"context"
	"math/big"

	"cosmossdk.io/core/address"
	cosmos "github.com/cosmos/cosmos-sdk/types"

	authmodule "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramsmodule "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingmodule "github.com/cosmos/cosmos-sdk/x/staking/types"

	feemodule "github.com/artela-network/artela/x/fee/types"
)

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	NewAccountWithAddress(ctx context.Context, addr cosmos.AccAddress) cosmos.AccountI
	GetModuleAddress(moduleName string) cosmos.AccAddress
	GetAllAccounts(ctx context.Context) (accounts []cosmos.AccountI)
	IterateAccounts(ctx context.Context, cb func(account cosmos.AccountI) bool)
	GetSequence(context.Context, cosmos.AccAddress) (uint64, error)
	GetAccount(ctx context.Context, addr cosmos.AccAddress) cosmos.AccountI
	SetAccount(ctx context.Context, account cosmos.AccountI)
	RemoveAccount(ctx context.Context, account cosmos.AccountI)
	GetParams(ctx context.Context) (params authmodule.Params)
	AddressCodec() address.Codec
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	authmodule.BankKeeper
	GetBalance(ctx context.Context, addr cosmos.AccAddress, denom string) cosmos.Coin
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr cosmos.AccAddress, amt cosmos.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt cosmos.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt cosmos.Coins) error
}

// StakingKeeper returns the historical headers kept in store.
type StakingKeeper interface {
	GetHistoricalInfo(ctx context.Context, height int64) (stakingmodule.HistoricalInfo, error)
	GetValidatorByConsAddr(ctx context.Context, consAddr cosmos.ConsAddress) (validator stakingmodule.Validator, err error)
}

// FeeKeeper
type FeeKeeper interface {
	GetBaseFee(ctx cosmos.Context) *big.Int
	GetParams(ctx cosmos.Context) feemodule.Params
	AddTransientGasWanted(ctx cosmos.Context, gasWanted uint64) (uint64, error)
}

type (
	LegacyParams = paramsmodule.ParamSet
	// Subspace defines an interface that implements the legacy Cosmos SDK x/params Subspace type.
	// NOTE: This is used solely for migration of the Cosmos SDK x/params managed parameters.
	Subspace interface {
		GetParamSetIfExists(ctx cosmos.Context, ps LegacyParams)
	}
)
