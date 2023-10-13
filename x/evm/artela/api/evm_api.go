package api

import (
	"context"
	ethermint "github.com/artela-network/artela/ethereum/types"
	txs "github.com/artela-network/artela/x/evm/txs"
	"github.com/artela-network/artela/x/evm/txs/support"
	inherent "github.com/artela-network/artelasdk/chaincoreext/jit_inherent"
	"github.com/artela-network/artelasdk/integration"
	artelatypes "github.com/artela-network/artelasdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"strconv"
)

var (
	_               artelatypes.EvmHostApi = (*evmHostApi)(nil)
	evmHostInstance *evmHostApi
)

type evmHostApi struct {
	getCtxByHeight func(height int64, prove bool) (sdk.Context, error)
	ethCall        func(c context.Context, req *txs.EthCallRequest) (*txs.MsgEthereumTxResponse, error)
}

func NewEvmHostInstance(getCtxByHeight func(height int64, prove bool) (sdk.Context, error),
	ethCall func(c context.Context, req *txs.EthCallRequest) (*txs.MsgEthereumTxResponse, error),
) {
	evmHostInstance = &evmHostApi{
		getCtxByHeight: getCtxByHeight,
		ethCall:        ethCall,
	}
}

func GetEvmHostInstance() (artelatypes.EvmHostApi, error) {
	if evmHostInstance == nil {
		return nil, errors.New("AspectStateHostApi not init")
	}
	return evmHostInstance, nil
}

func (e evmHostApi) StaticCall(ctx *artelatypes.RunnerContext, request *artelatypes.EthMessage) *artelatypes.EthMessageCallResult {

	sdkCtx, err := e.getCtxByHeight(ctx.BlockNumber, false)
	if err != nil {
		return artelatypes.ErrEthMessageCallResult(err)
	}
	marshal, jsonErr := jsoniter.Marshal(request)
	if jsonErr != nil {
		return artelatypes.ErrEthMessageCallResult(jsonErr)
	}
	parseUint, parseErr := strconv.ParseUint(request.GasFeeCap, 10, 64)
	if parseErr != nil {
		return artelatypes.ErrEthMessageCallResult(parseErr)
	}
	chainID, chainErr := ethermint.ParseChainID(sdkCtx.ChainID())
	if chainErr != nil {
		return artelatypes.ErrEthMessageCallResult(chainErr)
	}

	ethRequest := &txs.EthCallRequest{
		Args:            marshal,
		GasCap:          parseUint,
		ProposerAddress: nil,
		ChainId:         chainID.Int64(),
	}
	call, ethErr := e.ethCall(sdkCtx.Context(), ethRequest)
	if ethErr != nil {
		return artelatypes.ErrEthMessageCallResult(ethErr)
	}
	return &artelatypes.EthMessageCallResult{
		Hash:    call.Hash,
		Logs:    ConvertEthLogs(call.Logs),
		Ret:     call.Ret,
		VmError: call.VmError,
		GasUsed: call.GasUsed,
	}

}

func ConvertEthLogs(logs []*support.Log) []*artelatypes.EthLog {
	if logs == nil {
		return nil
	}
	ethLogs := make([]*artelatypes.EthLog, len(logs))
	for i, log := range logs {
		ethLogs[i] = ConvertEthLog(log)
	}
	return ethLogs
}
func ConvertEthLog(logs *support.Log) *artelatypes.EthLog {
	if logs == nil {
		return nil
	}
	topicStrArray := make([]string, len(logs.Topics))
	for i, topic := range logs.Topics {
		topicStrArray[i] = topic
	}
	return &artelatypes.EthLog{
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

func (e evmHostApi) JITCall(ctx *artelatypes.RunnerContext, request *artelatypes.JitInherentRequest) *artelatypes.JitInherentResponse {
	// determine jit call stage
	var stage integration.JoinPointStage
	switch artelatypes.PointCut(ctx.Point) {
	case artelatypes.PRE_TX_EXECUTE_METHOD, artelatypes.POST_TX_EXECUTE_METHOD,
		artelatypes.PRE_CONTRACT_CALL_METHOD, artelatypes.POST_CONTRACT_CALL_METHOD:
		stage = integration.TransactionExecution
	case artelatypes.ON_TX_VERIFY_METHOD, artelatypes.ON_ACCOUNT_VERIFY_METHOD:
		stage = integration.PreTransactionExecution
	case artelatypes.ON_TX_COMMIT_METHOD:
		stage = integration.PostTransactionExecution
	case artelatypes.ON_BLOCK_INITIALIZE_METHOD:
		stage = integration.BlockInitialization
	case artelatypes.ON_BLOCK_FINALIZE_METHOD:
		stage = integration.BlockFinalization
	default:
		log.Error("unsupported join point for jit call", "point", ctx.Point)
		return &artelatypes.JitInherentResponse{Success: false}
	}

	// convert aspect id to address
	aspect := *ctx.AspectId

	// FIXME: get leftover gas from last evm
	resp, err := inherent.Get().Submit(aspect, ctx.Gas, stage, request)

	if err != nil {
		//if errors.Is(err, vm.ErrOutOfGas) {
		resp.Success = false
		resp.ErrorMsg = err.Error()
		//	}

		log.Error("jit inherent submit fail", "err", err)
	}

	ctx.Gas = resp.LeftoverGas

	return resp
}
