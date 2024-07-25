package datactx

import (
	"errors"

	"google.golang.org/protobuf/proto"

	"github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/aspect-core/context"
	artelatypes "github.com/artela-network/aspect-core/types"
)

type BlockContextFieldLoader func(blockCtx *types.EthBlockContext) proto.Message

type BlockContext struct {
	blockContentLoaders map[string]BlockContextFieldLoader
	ctx                 *types.AspectRuntimeContext
}

func NewBlockContext(ctx *types.AspectRuntimeContext) *BlockContext {
	blockCtx := &BlockContext{
		blockContentLoaders: make(map[string]BlockContextFieldLoader),
		ctx:                 ctx,
	}
	blockCtx.registerLoaders()
	return blockCtx
}

func (c *BlockContext) registerLoaders() {
	loaders := c.blockContentLoaders
	loaders[context.BlockHeaderParentHash] = func(blockCtx *types.EthBlockContext) proto.Message {
		return &artelatypes.BytesData{Data: blockCtx.BlockHeader().ParentHash.Bytes()}
	}
	loaders[context.BlockHeaderMiner] = func(blockCtx *types.EthBlockContext) proto.Message {
		return &artelatypes.BytesData{Data: blockCtx.BlockHeader().Coinbase.Bytes()}
	}
	loaders[context.BlockHeaderTransactionsRoot] = func(blockCtx *types.EthBlockContext) proto.Message {
		return &artelatypes.BytesData{Data: blockCtx.BlockHeader().TxHash.Bytes()}
	}
	loaders[context.BlockHeaderNumber] = func(blockCtx *types.EthBlockContext) proto.Message {
		number := blockCtx.BlockHeader().Number.Uint64()
		return &artelatypes.UintData{Data: &number}
	}
	loaders[context.BlockHeaderTimestamp] = func(blockCtx *types.EthBlockContext) proto.Message {
		time := blockCtx.BlockHeader().Time
		return &artelatypes.UintData{Data: &time}
	}
}

func (c *BlockContext) ValueLoader(key string) ContextLoader {
	return func(ctx *artelatypes.RunnerContext) ([]byte, error) {
		if ctx == nil {
			return nil, errors.New("aspect context error, missing important information")
		}
		return proto.Marshal(c.blockContentLoaders[key](c.ctx.EthBlockContext()))
	}
}
