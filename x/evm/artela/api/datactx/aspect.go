package datactx

import (
	"errors"

	"google.golang.org/protobuf/proto"

	aspctx "github.com/artela-network/aspect-core/context"
	artelatypes "github.com/artela-network/aspect-core/types"
)

type AspectContextFieldLoader func(ctx *artelatypes.RunnerContext) proto.Message

type AspectContext struct {
	aspectContentLoader map[string]AspectContextFieldLoader
}

func NewAspectContext() *AspectContext {
	aspectCtx := &AspectContext{
		aspectContentLoader: make(map[string]AspectContextFieldLoader),
	}
	aspectCtx.registerLoaders()
	return aspectCtx
}

func (c *AspectContext) registerLoaders() {
	loaders := c.aspectContentLoader
	loaders[aspctx.AspectId] = func(ctx *artelatypes.RunnerContext) proto.Message {
		return &artelatypes.BytesData{Data: ctx.AspectId.Bytes()}
	}
	loaders[aspctx.AspectVersion] = func(ctx *artelatypes.RunnerContext) proto.Message {
		return &artelatypes.UintData{Data: &ctx.AspectVersion}
	}
}

func (c *AspectContext) ValueLoader(key string) ContextLoader {
	return func(ctx *artelatypes.RunnerContext) ([]byte, error) {
		if ctx == nil {
			return nil, errors.New("aspect context error, missing important information")
		}
		return proto.Marshal(c.aspectContentLoader[key](ctx))
	}
}
