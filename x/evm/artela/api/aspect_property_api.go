package api

import (
	"context"
	"errors"
	"github.com/artela-network/artela/x/evm/artela/contract"
	"github.com/artela-network/artela/x/evm/artela/types"
	asptypes "github.com/artela-network/aspect-core/types"
)

var (
	_ asptypes.AspectPropertyHostAPI = (*aspectPropertyHostAPI)(nil)
)

type aspectPropertyHostAPI struct {
	aspectRuntimeContext *types.AspectRuntimeContext
}

func (a *aspectPropertyHostAPI) Get(ctx *asptypes.RunnerContext, key string) []byte {
	// TODO: this part looks weird,
	//       but due to the time issue, we just migrate the old logics for now
	nativeContractStore := contract.NewAspectStore(a.aspectRuntimeContext.StoreKey(), a.aspectRuntimeContext.Logger())
	return nativeContractStore.GetAspectPropertyValue(a.aspectRuntimeContext.CosmosContext(), ctx.AspectId, key)
}

func GetAspectPropertyHostInstance(ctx context.Context) (asptypes.AspectPropertyHostAPI, error) {
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetAspectPropertyHostInstance: unwrap AspectRuntimeContext failed")
	}
	return &aspectPropertyHostAPI{aspectCtx}, nil
}
