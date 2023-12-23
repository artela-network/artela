package api

import (
	"context"
	"github.com/artela-network/artela/x/evm/artela/types"
	artelatypes "github.com/artela-network/aspect-core/types"
	"github.com/pkg/errors"
)

func GetStateDBHostInstance(ctx context.Context) (artelatypes.StateDBHostAPI, error) {
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetStateDBHostInstance: unwrap AspectRuntimeContext failed")
	}
	return aspectCtx.StateDb(), nil
}
