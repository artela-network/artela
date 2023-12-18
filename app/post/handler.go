package post

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/baseapp"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
)

// NewPostHandler returns an post handler responsible for attempting to route an
// Ethereum or SDK transaction to an internal post handler for performing
// transaction-level processing (e.g. fee payment, signature verification) before
// being passed onto it's respective handler.
func NewPostHandler(app *baseapp.BaseApp, options PostDecorators) cosmos.PostHandler {
	return func(
		ctx cosmos.Context, tx cosmos.Tx, sim, success bool,
	) (newCtx cosmos.Context, err error) {
		var postHandler cosmos.PostHandler

		txWithExtensions, ok := tx.(authante.HasExtensionOptionsTx)
		if ok {
			opts := txWithExtensions.GetExtensionOptions()
			if len(opts) > 0 {
				switch typeURL := opts[0].GetTypeUrl(); typeURL {
				case "/artela.evm.v1.ExtensionOptionsEthereumTx":
					// handle as *evmtypes.MsgEthereumTx
					postHandler = newEVMPostHandler(app, options)
				case "/artela.types.v1.ExtensionOptionsWeb3Tx":
					// handle as normal Cosmos SDK tx, except signature is checked for EIP712 representation
					postHandler = newLegacyCosmosPostHandlerEip712(options)
				case "/artela.types.v1.ExtensionOptionDynamicFeeTx":
					// cosmos-cosmos tx with dynamic fee extension
					postHandler = newCosmosPostHandler(options)
				default:
					return ctx, errorsmod.Wrapf(
						errortypes.ErrUnknownExtensionOptions,
						"rejecting tx with unsupported extension option: %s", typeURL,
					)
				}

				return postHandler(ctx, tx, sim, success)
			}
		}

		// handle as totally normal Cosmos SDK tx
		switch tx.(type) {
		case cosmos.Tx:
			postHandler = newCosmosPostHandler(options)
		default:
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid transaction type: %T", tx)
		}

		return postHandler(ctx, tx, sim, success)
	}
}
