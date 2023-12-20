package api

import (
	"context"
	"strconv"

	"github.com/artela-network/aspect-core/integration"
	coretypes "github.com/artela-network/aspect-core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	artela "github.com/artela-network/artela/ethereum/types"
	artelatypes "github.com/artela-network/artela/x/evm/artela/types"
	types "github.com/artela-network/artela/x/evm/txs"
	"github.com/artela-network/artela/x/evm/txs/support"
)

var (
	_               coretypes.EvmHostApi = (*evmHostApi)(nil)
	evmHostInstance *evmHostApi          // only use this to cache the handler of getCtxByHeight, ethCall
)

type evmHostApi struct {
	aspectCtx *artelatypes.AspectRuntimeContext

	getCtxByHeight func(height int64, prove bool) (sdk.Context, error)
	ethCall        func(c context.Context, req *types.EthCallRequest) (*types.MsgEthereumTxResponse, error)
}

func NewEvmHostInstance(getCtxByHeight func(height int64, prove bool) (sdk.Context, error),
	ethCall func(c context.Context, req *types.EthCallRequest) (*types.MsgEthereumTxResponse, error),
) {
	evmHostInstance = &evmHostApi{
		getCtxByHeight: getCtxByHeight,
		ethCall:        ethCall,
	}
}

func GetEvmHostInstance(ctx context.Context) (coretypes.EvmHostApi, error) {
	aspectCtx, ok := ctx.(*artelatypes.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetEvmHostInstance: unwrap AspectRuntimeContext failed")
	}
	return &evmHostApi{
		aspectCtx:      aspectCtx,
		getCtxByHeight: evmHostInstance.getCtxByHeight,
		ethCall:        evmHostInstance.ethCall,
	}, nil
}

func (e evmHostApi) StaticCall(ctx *coretypes.RunnerContext, request *coretypes.EthMessage) *coretypes.EthMessageCallResult {
	sdkCtx, err := e.getCtxByHeight(ctx.BlockNumber, false)
	if err != nil {
		return coretypes.ErrEthMessageCallResult(err)
	}
	marshal, jsonErr := jsoniter.Marshal(request)
	if jsonErr != nil {
		return coretypes.ErrEthMessageCallResult(jsonErr)
	}
	parseUint, parseErr := strconv.ParseUint(request.GasFeeCap, 10, 64)
	if parseErr != nil {
		return coretypes.ErrEthMessageCallResult(parseErr)
	}
	chainID, chainErr := artela.ParseChainID(sdkCtx.ChainID())
	if chainErr != nil {
		return coretypes.ErrEthMessageCallResult(chainErr)
	}

	ethRequest := &types.EthCallRequest{
		Args:            marshal,
		GasCap:          parseUint,
		ProposerAddress: nil,
		ChainId:         chainID.Int64(),
	}
	call, ethErr := e.ethCall(sdkCtx.Context(), ethRequest)
	if ethErr != nil {
		return coretypes.ErrEthMessageCallResult(ethErr)
	}
	return &coretypes.EthMessageCallResult{
		Hash:    call.Hash,
		Logs:    ConvertEthLogs(call.Logs),
		Ret:     call.Ret,
		VmError: call.VmError,
		GasUsed: call.GasUsed,
	}
}

func ConvertEthLogs(logs []*support.Log) []*coretypes.EthLog {
	if logs == nil {
		return nil
	}
	ethLogs := make([]*coretypes.EthLog, len(logs))
	for i, log := range logs {
		ethLogs[i] = ConvertEthLog(log)
	}
	return ethLogs
}

func ConvertEthLog(logs *support.Log) *coretypes.EthLog {
	if logs == nil {
		return nil
	}
	topicStrArray := make([]string, len(logs.Topics))
	copy(topicStrArray, logs.Topics)

	return &coretypes.EthLog{
		Address:     logs.Address,
		Topics:      topicStrArray,
		Data:        logs.Data,
		BlockNumber: logs.BlockNumber,
		TxHash:      logs.TxHash,
		TxIndex:     logs.TxIndex,
		BlockHash:   logs.BlockHash,
		Index:       logs.Index,
		Removed:     logs.Removed,
	}
}

func (e evmHostApi) JITCall(ctx *coretypes.RunnerContext, request *coretypes.JitInherentRequest) *coretypes.JitInherentResponse {
	// determine jit call stage
	var stage integration.JoinPointStage
	switch coretypes.PointCut(ctx.Point) {
	case coretypes.PRE_TX_EXECUTE_METHOD, coretypes.POST_TX_EXECUTE_METHOD,
		coretypes.PRE_CONTRACT_CALL_METHOD, coretypes.POST_CONTRACT_CALL_METHOD:
		stage = integration.TransactionExecution
	case coretypes.VERIFY_TX, coretypes.ON_ACCOUNT_VERIFY_METHOD:
		stage = integration.PreTransactionExecution
	case coretypes.POST_TX_COMMIT:
		stage = integration.PostTransactionExecution
	case coretypes.ON_BLOCK_INITIALIZE_METHOD:
		stage = integration.BlockInitialization
	case coretypes.ON_BLOCK_FINALIZE_METHOD:
		stage = integration.BlockFinalization
	default:
		log.Error("unsupported join point for jit call", "point", ctx.Point)
		return &coretypes.JitInherentResponse{Success: false}
	}

	// convert aspect id to address
	aspect := *ctx.AspectId

	// FIXME: get leftover gas from last evm
	resp, gas, err := e.aspectCtx.JITManager().Submit(ctx.Ctx, aspect, ctx.Gas, stage, request)
	if err != nil {
		if resp == nil {
			resp = &coretypes.JitInherentResponse{}
		}

		resp.Success = false
		resp.ErrorMsg = err.Error()
		log.Error("jit inherent submit fail", "err", err)
	}

	ctx.Gas = gas

	return resp
}
