package evm

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/baseapp"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/artela-network/artela/app/interfaces"
	"github.com/artela-network/artela/x/evm/txs"
)

// EthSigVerificationDecorator validates an ethereum signatures
type EthSigVerificationDecorator struct {
	evmKeeper interfaces.EVMKeeper
	app       *baseapp.BaseApp
}

// NewEthSigVerificationDecorator creates a new EthSigVerificationDecorator
func NewEthSigVerificationDecorator(app *baseapp.BaseApp, ek interfaces.EVMKeeper) EthSigVerificationDecorator {
	return EthSigVerificationDecorator{
		evmKeeper: ek,
		app:       app,
	}
}

// AnteHandle validates checks that the registered chain id is the same as the one on the message, and
// that the signer address matches the one defined on the message.
// It's not skipped for RecheckTx, because it set `From` address which is critical from other ante handler to work.
// Failure in RecheckTx will prevent tx to be included into block, especially when CheckTx succeed, in which case user
// won't see the error message.
func (esvd EthSigVerificationDecorator) AnteHandle(ctx cosmos.Context, tx cosmos.Tx, simulate bool, next cosmos.AnteHandler) (newCtx cosmos.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*txs.MsgEthereumTx)
		if !ok {
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*txs.MsgEthereumTx)(nil))
		}

		ethTx := msgEthTx.AsTransaction()
		sender, _, err := esvd.evmKeeper.VerifySig(ctx, ethTx)
		if err != nil {
			return ctx, err
		}

		// set up the sender to the transaction field if not already
		msgEthTx.From = sender.Hex()
	}

	return next(ctx, tx, simulate)
}
