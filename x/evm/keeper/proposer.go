package keeper

import (
	errorsmod "cosmossdk.io/errors"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	stakingmodule "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
)

// GetProposerAddress returns the block proposer's validator operator address.
func (k Keeper) GetProposerAddress(ctx cosmos.Context, proposerAddress cosmos.ConsAddress) (common.Address, error) {
	validator, err := k.stakingKeeper.GetValidatorByConsAddr(ctx, GetProposerAddress(ctx, proposerAddress))
	if err != nil {
		return common.Address{}, errorsmod.Wrapf(
			stakingmodule.ErrNoValidatorFound,
			"failed to retrieve validator from block proposer address %s, %v",
			proposerAddress.String(),
			err,
		)
	}

	addr, err := cosmos.ValAddressFromBech32(validator.GetOperator())
	if err != nil {
		panic(err)
	}

	return common.BytesToAddress(addr), nil
}

// GetProposerAddress returns current block proposer's address when provided proposer address is empty.
func GetProposerAddress(ctx cosmos.Context, proposerAddress cosmos.ConsAddress) cosmos.ConsAddress {
	if len(proposerAddress) > 0 {
		return proposerAddress
	}
	return ctx.BlockHeader().ProposerAddress
}
