package utils

import (
	cosmos "github.com/cosmos/cosmos-sdk/types"
	stakingmodule "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// BankKeeper defines the exposed interface for using functionality of the bank keeper
// in the context of the AnteHandler utils package.
type BankKeeper interface {
	GetBalance(ctx cosmos.Context, addr cosmos.AccAddress, denom string) cosmos.Coin
}

// DistributionKeeper defines the exposed interface for using functionality of the distribution
// keeper in the context of the AnteHandler utils package.
type DistributionKeeper interface {
	WithdrawDelegationRewards(ctx cosmos.Context, delAddr cosmos.AccAddress, valAddr cosmos.ValAddress) (cosmos.Coins, error)
}

// StakingKeeper defines the exposed interface for using functionality of the staking keeper
// in the context of the AnteHandler utils package.
type StakingKeeper interface {
	BondDenom(ctx cosmos.Context) string
	IterateDelegations(ctx cosmos.Context, delegator cosmos.AccAddress, fn func(index int64, delegation stakingmodule.DelegationI) (stop bool))
}
