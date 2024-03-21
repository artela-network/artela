package evm

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/baseapp"
	cosmos "github.com/cosmos/cosmos-sdk/types"

	"github.com/artela-network/artela/app/interfaces"
	"github.com/artela-network/artela/x/evm/artela/types"
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

func (aspd AspectRuntimeContextDecorator) PostHandle(ctx cosmos.Context, tx cosmos.Tx, simulate, success bool, next cosmos.PostHandler) (newCtx cosmos.Context, err error) {
	// Aspect Runtime Context Lifecycle: destroy AspectRuntimeContext
	aspectCtx, ok := ctx.Value(types.AspectContextKey).(*types.AspectRuntimeContext)
	if !ok {
		return ctx, errors.New("EthereumTx: unwrap AspectRuntimeContext failed")
	}

	aspectCtx.Destroy()

	return next(ctx, tx, simulate, success)
}
