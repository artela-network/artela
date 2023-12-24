package api

import (
	"context"
	"errors"
	"github.com/artela-network/artela/x/evm/artela/types"
	asptypes "github.com/artela-network/aspect-core/types"
	"github.com/emirpasic/gods/sets/hashset"
)

var (
	_              asptypes.AspectStateHostAPI = (*aspectStateHostAPI)(nil)
	setConstraints                             = hashset.New(
		asptypes.PRE_CONTRACT_CALL_METHOD,
		asptypes.POST_CONTRACT_CALL_METHOD,
		asptypes.PRE_TX_EXECUTE_METHOD,
		asptypes.POST_TX_EXECUTE_METHOD,
		asptypes.OPERATION_METHOD,
	)
)

type aspectStateHostAPI struct {
	aspectRuntimeContext *types.AspectRuntimeContext
}

func (a *aspectStateHostAPI) Get(ctx *asptypes.RunnerContext, key string) []byte {
	return a.aspectRuntimeContext.GetAspectState(ctx, key)
}

func (a *aspectStateHostAPI) Set(ctx *asptypes.RunnerContext, key string, value []byte) {
	if !setConstraints.Contains(asptypes.PointCut(ctx.Point)) {
		panic("cannot set aspect state in current join point")
	}
	a.aspectRuntimeContext.SetAspectState(ctx, key, value)
}

func GetAspectStateHostInstance(ctx context.Context) (asptypes.AspectStateHostAPI, error) {
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetAspectRuntimeContextHostInstance: unwrap AspectRuntimeContext failed")
	}
	return &aspectStateHostAPI{aspectCtx}, nil
}
