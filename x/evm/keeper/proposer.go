package keeper

import (
	errorsmod "cosmossdk.io/errors"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	stakingmodule "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
)

// GetProposerAddress returns the block proposer's validator operator address.
func (k Keeper) GetProposerAddress(ctx cosmos.Context, proposerAddress cosmos.ConsAddress) (common.Address, error) {
	validator, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, GetProposerAddress(ctx, proposerAddress))
	if !found {
		return common.Address{}, errorsmod.Wrapf(
			stakingmodule.ErrNoValidatorFound,
			"failed to retrieve validator from block proposer address %s",
			proposerAddress.String(),
		)
	}

	return common.BytesToAddress(validator.GetOperator()), nil
}

// GetProposerAddress returns current block proposer's address when provided proposer address is empty.
func GetProposerAddress(ctx cosmos.Context, proposerAddress cosmos.ConsAddress) cosmos.ConsAddress {
	if len(proposerAddress) > 0 {
		return proposerAddress
	}
	return ctx.BlockHeader().ProposerAddress
}
