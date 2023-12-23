package api

import (
	"context"
	"errors"
	"github.com/artela-network/artela/x/evm/artela/types"
	asptypes "github.com/artela-network/aspect-core/types"
)

var (
	_ asptypes.AspectStateHostAPI = (*aspectStateHostAPI)(nil)
)

type aspectStateHostAPI struct {
	aspectRuntimeContext *types.AspectRuntimeContext
}

func (a *aspectStateHostAPI) Get(ctx *asptypes.RunnerContext, key string) []byte {
	return a.aspectRuntimeContext.GetAspectState(ctx, key)
}

func (a *aspectStateHostAPI) Set(ctx *asptypes.RunnerContext, key string, value []byte) {
	a.aspectRuntimeContext.SetAspectState(ctx, key, value)
}

func GetAspectStateHostInstance(ctx context.Context) (asptypes.AspectStateHostAPI, error) {
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetAspectRuntimeContextHostInstance: unwrap AspectRuntimeContext failed")
	}
	return &aspectStateHostAPI{aspectCtx}, nil
}
