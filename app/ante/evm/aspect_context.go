package evm

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/baseapp"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artela-network/artela/app/interfaces"
	"github.com/artela-network/artela/x/evm/artela/provider"
	"github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/artela/x/evm/states"
	"github.com/artela-network/artela/x/evm/txs"
	inherent "github.com/artela-network/aspect-core/chaincoreext/jit_inherent"
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

		// create a temporary state db for check tx use
		txConfig := states.NewEmptyTxConfig(common.BytesToHash(ctx.BlockHeader().DataHash))
		stateDB := states.New(ctx, aspd.evmKeeper, txConfig)

		// Aspect Runtime Context Lifecycle: create AspectRuntimeContext
		// this handler is for eth transaction only, should be 1 message for eth transaction.
		evmConfig, err := aspd.evmKeeper.EVMConfigFromCtx(ctx)
		if err != nil {
			return ctx, fmt.Errorf("failed to get evm config from context: %w", err)
		}
		ethTxContext := types.NewEthTxContext(msgEthTx.AsEthCallTransaction()).WithEVMConfig(evmConfig).WithStateDB(stateDB)
		aspectCtx := types.NewAspectRuntimeContext()
		protocol := provider.NewAspectProtocolProvider(aspectCtx.EthTxContext)
		jitManager := inherent.NewManager(protocol)
		aspectCtx.SetEthTxContext(ethTxContext, jitManager)

		if ctx.IsCheckTx() {
			// at check tx stage since current block proposal is not prepared,
			// so we can only initialize the block context with a height.
			// also that we should not provide jit manager at this stage.
			// use the height of next block as the block height is because if it is in check tx,
			// the state passed to us is still the state from last commit, so we need to increase the block height by one,
			// to make sure the number is the same as the one in deliver tx.
			aspectCtx.SetEthBlockContext(types.NewEthBlockContextFromHeight(ctx.BlockHeight() + 1))
		} else {
			// at deliver tx stage,
			// we can initialize the block context with the pending block proposal from
			// begin block.
			aspectCtx.SetEthBlockContext(aspd.evmKeeper.GetBlockContext())
		}
		aspectCtx.WithCosmosContext(ctx)
		aspectCtx.CreateStateObject()

		ctx = ctx.WithValue(types.AspectContextKey, aspectCtx)
	}

	return next(ctx, tx, simulate)
}
