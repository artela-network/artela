package keeper

import (
	"context"
	"errors"

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

func (k Keeper) JITSenderAspectByContext(ctx context.Context, userOpHash common.Hash) (common.Address, error) {
	aspectCtx, ok := ctx.(*artvmtype.AspectRuntimeContext)
	if !ok {
		return common.Address{}, errors.New("JITSenderAspectByContext: unwrap AspectRuntimeContext failed")
	}
	return aspectCtx.JITManager().SenderAspect(userOpHash), nil
}

func (k Keeper) GetAspectContext(ctx context.Context, aspectId string, key string) (string, error) {
	aspectCtx, ok := ctx.(*artvmtype.AspectRuntimeContext)
	if !ok {
		return "", errors.New("GetAspectContext: unwrap AspectRuntimeContext failed")
	}
	return aspectCtx.AspectContext().Get(aspectId, key), nil
}

func (k Keeper) SetAspectContext(ctx context.Context, aspectId string, key string, value string) error {
	aspectCtx, ok := ctx.(*artvmtype.AspectRuntimeContext)
	if !ok {
		return errors.New("SetAspectContext: unwrap AspectRuntimeContext failed")
	}
	aspectCtx.AspectContext().Add(aspectId, key, value)
	return nil
}
