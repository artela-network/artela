package keeper

import (
	"context"

	artvmtype "github.com/artela-network/artela/x/evm/artela/types"
	cosmosAspect "github.com/cosmos/cosmos-sdk/aspect/cosmos"
	"github.com/ethereum/go-ethereum/common"
)

func (k Keeper) GetAspectCosmosProvider() cosmosAspect.AspectCosmosProvider {
	return k.aspect
}

func (k Keeper) GetAspectRuntimeContext() *artvmtype.AspectRuntimeContext {
	return k.aspectRuntimeContext
}

func (k Keeper) JITSenderAspectByContext(ctx context.Context, userOpHash common.Hash) common.Address {
	aspectCtx, ok := ctx.(*artvmtype.AspectRuntimeContext)
	if !ok {
		// TODO add log
		// logger.Debug("JITSenderAspectByContext: unwrap AspectRuntimeContext failed")
		return common.Address{}
	}
	return aspectCtx.JITManager().SenderAspect(userOpHash)
}

func (k Keeper) GetAspectContext(ctx context.Context, aspectId string, key string) string {
	aspectCtx, ok := ctx.(*artvmtype.AspectRuntimeContext)
	if !ok {
		// TODO add log
		// logger.Debug("GetAspectContext: unwrap AspectRuntimeContext failed")
		return ""
	}
	return aspectCtx.AspectContext().Get(aspectId, key)
}

func (k Keeper) SetAspectContext(ctx context.Context, aspectId string, key string, value string) {
	aspectCtx, ok := ctx.(*artvmtype.AspectRuntimeContext)
	if !ok {
		// TODO add log
		// logger.Debug("GetAspectContext: unwrap AspectRuntimeContext failed")
		return
	}
	aspectCtx.AspectContext().Add(aspectId, key, value)
}
