package keeper

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"

	artvmtype "github.com/artela-network/artela/x/evm/artela/types"
)

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

func (k Keeper) IsCommit(ctx context.Context) bool {
	aspectCtx, ok := ctx.(*artvmtype.AspectRuntimeContext)
	if !ok || aspectCtx.EthTxContext() == nil {
		return false
	}

	return aspectCtx.EthTxContext().Commit()
}

func (k Keeper) GetAspectContext(ctx context.Context, address common.Address, key string) ([]byte, error) {
	aspectCtx, ok := ctx.(*artvmtype.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetAspectContext: unwrap AspectRuntimeContext failed")
	}
	return aspectCtx.AspectContext().Get(address, key), nil
}

func (k Keeper) SetAspectContext(ctx context.Context, address common.Address, key string, value []byte) error {
	aspectCtx, ok := ctx.(*artvmtype.AspectRuntimeContext)
	if !ok {
		return errors.New("SetAspectContext: unwrap AspectRuntimeContext failed")
	}
	aspectCtx.AspectContext().Add(address, key, value)
	return nil
}

func (k Keeper) GetBlockContext() *artvmtype.EthBlockContext {
	return k.BlockContext
}
