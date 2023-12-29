package datactx

import (
	"errors"

	aspctx "github.com/artela-network/aspect-core/context"
	artelatypes "github.com/artela-network/aspect-core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	"google.golang.org/protobuf/proto"

	"github.com/artela-network/artela/x/evm/artela/types"
)

type MessageContextFieldLoader func(ethTxCtx *types.EthTxContext, message *core.Message) proto.Message

type MessageContext struct {
	getSdkCtx            func() sdk.Context
	getEthTxContext      func() *types.EthTxContext
	messageContentLoader map[string]MessageContextFieldLoader
}

func NewMessageContext(getEthTxContext func() *types.EthTxContext,
	getSdkCtx func() sdk.Context) *MessageContext {
	msgCtx := &MessageContext{
		messageContentLoader: make(map[string]MessageContextFieldLoader),
		getEthTxContext:      getEthTxContext,
		getSdkCtx:            getSdkCtx,
	}
	msgCtx.registerLoaders()
	return msgCtx
}

func (c *MessageContext) registerLoaders() {
	loaders := c.messageContentLoader
	loaders[aspctx.MsgIndex] = func(ethTxCtx *types.EthTxContext, message *core.Message) proto.Message {
		index := ethTxCtx.VmTracer().CurrentCallIndex()
		return &artelatypes.UintData{Data: &index}
	}
	loaders[aspctx.MsgFrom] = func(ethTxCtx *types.EthTxContext, message *core.Message) proto.Message {
		return &artelatypes.BytesData{Data: message.From.Bytes()}
	}
	loaders[aspctx.MsgTo] = func(ethTxCtx *types.EthTxContext, message *core.Message) proto.Message {
		return &artelatypes.BytesData{Data: message.To.Bytes()}
	}
	loaders[aspctx.MsgValue] = func(ethTxCtx *types.EthTxContext, message *core.Message) proto.Message {
		return &artelatypes.BytesData{Data: message.Value.Bytes()}
	}
	loaders[aspctx.MsgInput] = func(ethTxCtx *types.EthTxContext, message *core.Message) proto.Message {
		return &artelatypes.BytesData{Data: message.Data}
	}
	loaders[aspctx.MsgGas] = func(ethTxCtx *types.EthTxContext, message *core.Message) proto.Message {

		gasLimit := message.GasLimit
		return &artelatypes.UintData{Data: &gasLimit}
	}
	loaders[aspctx.MsgResultRet] = func(ethTxCtx *types.EthTxContext, message *core.Message) proto.Message {
		tracer := ethTxCtx.VmTracer()
		callIdx := tracer.CurrentCallIndex()
		currentCall := tracer.CallTree().FindCall(callIdx)

		return &artelatypes.BytesData{Data: currentCall.Ret}
	}
	loaders[aspctx.MsgResultGasUsed] = func(ethTxCtx *types.EthTxContext, message *core.Message) proto.Message {
		tracer := ethTxCtx.VmTracer()
		callIdx := tracer.CurrentCallIndex()
		currentCall := tracer.CallTree().FindCall(callIdx)

		result := message.GasLimit - currentCall.RemainingGas
		return &artelatypes.UintData{Data: &result}
	}
	loaders[aspctx.MsgResultError] = func(ethTxCtx *types.EthTxContext, message *core.Message) proto.Message {
		tracer := ethTxCtx.VmTracer()
		callIdx := tracer.CurrentCallIndex()
		currentCall := tracer.CallTree().FindCall(callIdx)

		msg := currentCall.Err.Error()
		return &artelatypes.StringData{Data: &msg}
	}
}

func (c *MessageContext) ValueLoader(key string) ContextLoader {
	return func(ctx *artelatypes.RunnerContext) ([]byte, error) {
		if ctx == nil {
			return nil, errors.New("aspect context error, missing important information")
		}
		txContext := c.getEthTxContext()
		if txContext == nil {
			return nil, errors.New("tx context error, failed to load")
		}
		message := txContext.Message()
		if message == nil {
			return nil, errors.New("message error, failed to load")
		}
		return proto.Marshal(c.messageContentLoader[key](txContext, message))
	}
}
