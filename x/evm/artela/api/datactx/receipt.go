package datactx

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereum "github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/protobuf/proto"

	"github.com/artela-network/artela/x/evm/artela/types"
	aspctx "github.com/artela-network/aspect-core/context"
	artelatypes "github.com/artela-network/aspect-core/types"
)

type ReceiptContextFieldLoader func(ethTxCtx *types.EthTxContext, receipt *ethereum.Receipt) proto.Message

type ReceiptContext struct {
	getSdkCtx             func() sdk.Context
	getEthTxContext       func() *types.EthTxContext
	receiptContentLoaders map[string]ReceiptContextFieldLoader
}

func NewReceiptContext(getEthTxContext func() *types.EthTxContext,
	getSdkCtx func() sdk.Context) *ReceiptContext {
	receiptCtx := &ReceiptContext{
		receiptContentLoaders: make(map[string]ReceiptContextFieldLoader),
		getEthTxContext:       getEthTxContext,
		getSdkCtx:             getSdkCtx,
	}
	receiptCtx.registerLoaders()
	return receiptCtx
}

func (c *ReceiptContext) registerLoaders() {
	loaders := c.receiptContentLoaders
	loaders[aspctx.ReceiptStatus] = func(_ *types.EthTxContext, receipt *ethereum.Receipt) proto.Message {
		return &artelatypes.UintData{Data: &receipt.Status}
	}
	loaders[aspctx.ReceiptLogs] = func(_ *types.EthTxContext, receipt *ethereum.Receipt) proto.Message {
		logs := make([]*artelatypes.EthLog, 0, len(receipt.Logs))
		for _, log := range receipt.Logs {
			topics := make([][]byte, 0, len(log.Topics))
			for _, topic := range log.Topics {
				topics = append(topics, topic.Bytes())
			}
			index := uint64(log.Index)
			logs = append(logs, &artelatypes.EthLog{
				Address: log.Address.Bytes(),
				Topics:  topics,
				Data:    log.Data,
				Index:   &index,
			})
		}
		return &artelatypes.EthLogs{Logs: logs}
	}
	loaders[aspctx.ReceiptGasUsed] = func(_ *types.EthTxContext, receipt *ethereum.Receipt) proto.Message {
		return &artelatypes.UintData{Data: &receipt.GasUsed}
	}
	loaders[aspctx.ReceiptCumulativeGasUsed] = func(_ *types.EthTxContext, receipt *ethereum.Receipt) proto.Message {
		return &artelatypes.UintData{Data: &receipt.CumulativeGasUsed}
	}
	loaders[aspctx.ReceiptBloom] = func(_ *types.EthTxContext, receipt *ethereum.Receipt) proto.Message {
		return &artelatypes.BytesData{Data: receipt.Bloom.Bytes()}
	}
}

func (c *ReceiptContext) ValueLoader(key string) ContextLoader {
	return func(ctx *artelatypes.RunnerContext) ([]byte, error) {
		if ctx == nil {
			return nil, errors.New("aspect context error, missing important information")
		}
		txContext := c.getEthTxContext()
		if txContext == nil {
			return nil, errors.New("tx context error, failed to load")
		}
		receipt := txContext.Receipt()
		if receipt == nil {
			return nil, errors.New("receipt error, failed to load")
		}
		return proto.Marshal(c.receiptContentLoaders[key](txContext, receipt))
	}
}
