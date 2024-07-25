package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	govmodule "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/artela-network/artela/x/fee/types"
)

// UpdateParams implements the gRPC MsgServer interface. When an UpdateParams
// proposal passes, it updates the module parameters. The update can only be
// performed if the requested authority is the Cosmos SDK governance module
// account.
func (k *Keeper) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority.String() != req.Authority {
		return nil, errorsmod.Wrapf(govmodule.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority.String(), req.Authority)
	}

	ctx := cosmos.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
