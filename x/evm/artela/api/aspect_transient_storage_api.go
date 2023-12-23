package api

import (
	"context"
	"errors"
	"github.com/artela-network/artela/x/evm/artela/types"
	asptypes "github.com/artela-network/aspect-core/types"
)

var (
	_ asptypes.AspectTransientStorageHostAPI = (*aspectTransientStorageHostAPI)(nil)
)

type aspectTransientStorageHostAPI struct {
	aspectRuntimeContext *types.AspectRuntimeContext
}

func (a *aspectTransientStorageHostAPI) Get(ctx *asptypes.RunnerContext, key string) []byte {
	return a.aspectRuntimeContext.AspectContext().Get(ctx.AspectId, key)
}

func (a *aspectTransientStorageHostAPI) Set(ctx *asptypes.RunnerContext, key string, value []byte) {
	a.aspectRuntimeContext.AspectContext().Add(ctx.AspectId, key, value)
}

func GetAspectTransientStorageHostInstance(ctx context.Context) (asptypes.AspectTransientStorageHostAPI, error) {
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetAspectRuntimeContextHostInstance: unwrap AspectRuntimeContext failed")
	}
	return &aspectTransientStorageHostAPI{aspectCtx}, nil
}
