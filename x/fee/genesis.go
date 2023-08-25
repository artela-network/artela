package fee

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/artela-network/artela/x/fee/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/artela-network/artela/x/fee/types"
)

// InitGenesis initializes genesis states based on exported genesis
func InitGenesis(
	ctx sdk.Context,
	k *keeper.Keeper,
	genState types.GenesisState,
) []abci.ValidatorUpdate {
	err := k.SetParams(ctx, genState.Params)
	if err != nil {
		panic(errorsmod.Wrap(err, "could not set parameters at genesis"))
	}

	k.SetBlockGasWanted(ctx, genState.BlockGas)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports genesis states of the fee market module
func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:   k.GetParams(ctx),
		BlockGas: k.GetBlockGasWanted(ctx),
	}
}
