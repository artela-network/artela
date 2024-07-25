package api

import (
	"context"
	"errors"

	"github.com/emirpasic/gods/sets/hashset"

	"github.com/ethereum/go-ethereum/common"

	"github.com/artela-network/artela/x/evm/artela/types"
	asptypes "github.com/artela-network/aspect-core/types"
)

var (
	_ asptypes.AspectTransientStorageHostAPI = (*aspectTransientStorageHostAPI)(nil)

	transientStorageConstrainedJoinPoints = hashset.New(
		asptypes.PRE_CONTRACT_CALL_METHOD,
		asptypes.POST_CONTRACT_CALL_METHOD,
		asptypes.PRE_TX_EXECUTE_METHOD,
		asptypes.POST_TX_EXECUTE_METHOD,
		asptypes.OPERATION_METHOD,
	)
)

type aspectTransientStorageHostAPI struct {
	aspectRuntimeContext *types.AspectRuntimeContext
}

func (a *aspectTransientStorageHostAPI) Get(ctx *asptypes.RunnerContext, aspectId []byte, key string) ([]byte, error) {
	if !transientStorageConstrainedJoinPoints.Contains(asptypes.PointCut(ctx.Point)) {
		return nil, errors.New("cannot get aspect transient storage in current join point")
	}

	var aspectAddr common.Address
	if len(aspectId) == 0 {
		aspectAddr = ctx.AspectId
	} else {
		aspectAddr = common.BytesToAddress(aspectId)
	}

	return a.aspectRuntimeContext.AspectContext().Get(aspectAddr, key), nil
}

func (a *aspectTransientStorageHostAPI) Set(ctx *asptypes.RunnerContext, key string, value []byte) error {
	if !transientStorageConstrainedJoinPoints.Contains(asptypes.PointCut(ctx.Point)) {
		return errors.New("cannot set aspect transient storage in current join point")
	}

	a.aspectRuntimeContext.AspectContext().Add(ctx.AspectId, key, value)
	return nil
}

func GetAspectTransientStorageHostInstance(ctx context.Context) (asptypes.AspectTransientStorageHostAPI, error) {
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetAspectTransientStorageHostInstance: unwrap AspectRuntimeContext failed")
	}
	return &aspectTransientStorageHostAPI{aspectCtx}, nil
}
