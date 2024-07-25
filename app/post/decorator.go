package post

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/baseapp"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	// vestingtypes "github.com/artela-network/artela/x/vesting/types"
	"github.com/artela-network/artela/app/interfaces"
	evmpost "github.com/artela-network/artela/app/post/evm"
)

// PostDecorators defines the list of module keepers required to run the Artela
// PostHandler decorators.
type PostDecorators struct {
	EvmKeeper interfaces.EVMKeeper
}

// Validate checks if the keepers are defined
func (options PostDecorators) Validate() error {
	if options.EvmKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "evm keeper is required for AnteHandler")
	}
	return nil
}

// newEVMPostHandler creates the default post handler for Ethereum transactions
func newEVMPostHandler(app *baseapp.BaseApp, options PostDecorators) cosmos.PostHandler {
	return cosmos.ChainPostDecorators(evmpost.NewAspectRuntimeContextDecorator(app, options.EvmKeeper))
}

// newCosmosPostHandler creates the default post handler for Cosmos transactions
func newCosmosPostHandler(options PostDecorators) cosmos.PostHandler {
	return cosmos.ChainPostDecorators(cosmos.Terminator{})
}

// newCosmosPostHandlerEip712 creates the post handler for transactions signed with EIP712
func newLegacyCosmosPostHandlerEip712(options PostDecorators) cosmos.PostHandler {
	return cosmos.ChainPostDecorators(cosmos.Terminator{})
}
