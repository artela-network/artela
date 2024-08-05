package api

import (
	"context"
	"errors"
	"github.com/artela-network/artela/x/evm/artela/types"
	asptypes "github.com/artela-network/aspect-core/types"
)

var _ asptypes.AspectPropertyHostAPI = (*aspectPropertyHostAPI)(nil)

type aspectPropertyHostAPI struct {
	aspectRuntimeContext *types.AspectRuntimeContext
}

func (a *aspectPropertyHostAPI) Get(ctx *asptypes.RunnerContext, key string) (ret []byte, err error) {
	ret = a.aspectRuntimeContext.GetAspectProperty(ctx, key)
	return
}

func GetAspectPropertyHostInstance(ctx context.Context) (asptypes.AspectPropertyHostAPI, error) {
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetAspectPropertyHostInstance: unwrap AspectRuntimeContext failed")
	}
	return &aspectPropertyHostAPI{aspectCtx}, nil
}
