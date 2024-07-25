package types

import (
	"math/big"

	cosmos "github.com/cosmos/cosmos-sdk/types"
	authmodule "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramsmodule "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingmodule "github.com/cosmos/cosmos-sdk/x/staking/types"

	feemodule "github.com/artela-network/artela/x/fee/types"
)

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	NewAccountWithAddress(ctx cosmos.Context, addr cosmos.AccAddress) authmodule.AccountI
	GetModuleAddress(moduleName string) cosmos.AccAddress
	GetAllAccounts(ctx cosmos.Context) (accounts []authmodule.AccountI)
	IterateAccounts(ctx cosmos.Context, cb func(account authmodule.AccountI) bool)
	GetSequence(cosmos.Context, cosmos.AccAddress) (uint64, error)
	GetAccount(ctx cosmos.Context, addr cosmos.AccAddress) authmodule.AccountI
	SetAccount(ctx cosmos.Context, account authmodule.AccountI)
	RemoveAccount(ctx cosmos.Context, account authmodule.AccountI)
	GetParams(ctx cosmos.Context) (params authmodule.Params)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	authmodule.BankKeeper
	GetBalance(ctx cosmos.Context, addr cosmos.AccAddress, denom string) cosmos.Coin
	SendCoinsFromModuleToAccount(ctx cosmos.Context, senderModule string, recipientAddr cosmos.AccAddress, amt cosmos.Coins) error
	MintCoins(ctx cosmos.Context, moduleName string, amt cosmos.Coins) error
	BurnCoins(ctx cosmos.Context, moduleName string, amt cosmos.Coins) error
}

// StakingKeeper returns the historical headers kept in store.
type StakingKeeper interface {
	GetHistoricalInfo(ctx cosmos.Context, height int64) (stakingmodule.HistoricalInfo, bool)
	GetValidatorByConsAddr(ctx cosmos.Context, consAddr cosmos.ConsAddress) (validator stakingmodule.Validator, found bool)
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
