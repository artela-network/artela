package evm

import (
	"github.com/artela-network/artela/app/interfaces"
	"github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/artela/x/evm/txs"
	"github.com/cosmos/cosmos-sdk/baseapp"

	errorsmod "cosmossdk.io/errors"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// CreateAspectRuntimeContextDecorator prepare the aspect runtime context
type AspectRuntimeContextDecorator struct {
	evmKeeper interfaces.EVMKeeper
	app       *baseapp.BaseApp
}

// NewAspectRuntimeContextDecorator creates a new AspectRuntimeContextDecorator
func NewAspectRuntimeContextDecorator(app *baseapp.BaseApp, ek interfaces.EVMKeeper) AspectRuntimeContextDecorator {
	return AspectRuntimeContextDecorator{
		evmKeeper: ek,
		app:       app,
	}
}

// AnteHandle validates checks that the registered chain id is the same as the one on the message, and
// that the signer address matches the one defined on the message.
// It's not skipped for RecheckTx, because it set `From` address which is critical from other ante handler to work.
// Failure in RecheckTx will prevent tx to be included into block, especially when CheckTx succeed, in which case user
// won't see the error message.
func (aspd AspectRuntimeContextDecorator) AnteHandle(ctx cosmos.Context, tx cosmos.Tx, simulate bool, next cosmos.AnteHandler) (newCtx cosmos.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*txs.MsgEthereumTx)
		if !ok {
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*txs.MsgEthereumTx)(nil))
		}

		// Aspect Runtime Context Lifecycle: create AspectRuntimeContext
		// this handler is for eth transaction only, should be 1 message for eth transaction.
		ethTxContext := types.NewEthTxContext(msgEthTx.AsEthCallTransaction())
		aspectCtx := types.NewAspectRuntimeContext()
		aspectCtx.SetEthTxContext(ethTxContext)
		aspectCtx.WithCosmosContext(ctx)

		ctx = ctx.WithValue(types.AspectContextKey, aspectCtx)
	}

	return next(ctx, tx, simulate)
}
