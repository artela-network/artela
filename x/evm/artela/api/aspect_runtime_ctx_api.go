package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/emirpasic/gods/sets/hashset"

	"github.com/cosmos/gogoproto/proto"

	"github.com/artela-network/artela/x/evm/artela/api/datactx"
	"github.com/artela-network/artela/x/evm/artela/types"
	aspctx "github.com/artela-network/aspect-core/context"
	asptypes "github.com/artela-network/aspect-core/types"
)

var (
	_                 asptypes.RuntimeContextHostAPI = (*aspectRuntimeContextHostAPI)(nil)
	ctxKeyConstraints                                = map[asptypes.PointCut]*hashset.Set{
		asptypes.INIT_METHOD:               hashset.New(aspctx.InitKeys...),
		asptypes.OPERATION_METHOD:          hashset.New(aspctx.OperationKeys...),
		asptypes.VERIFY_TX:                 hashset.New(aspctx.VerifyTxCtxKeys...),
		asptypes.PRE_TX_EXECUTE_METHOD:     hashset.New(aspctx.PreTxCtxKeys...),
		asptypes.POST_TX_EXECUTE_METHOD:    hashset.New(aspctx.PostTxCtxKeys...),
		asptypes.PRE_CONTRACT_CALL_METHOD:  hashset.New(aspctx.PreCallCtxKeys...),
		asptypes.POST_CONTRACT_CALL_METHOD: hashset.New(aspctx.PostCallCtxKeys...),
	}
)

type aspectRuntimeContextHostAPI struct {
	aspectRuntimeContext *types.AspectRuntimeContext
	execMap              map[string]datactx.ContextLoader
}

func newAspectRuntimeContextHostAPI(aspectRuntimeContext *types.AspectRuntimeContext) *aspectRuntimeContextHostAPI {
	instance := &aspectRuntimeContextHostAPI{
		aspectRuntimeContext: aspectRuntimeContext,
		execMap:              make(map[string]datactx.ContextLoader),
	}

	instance.Register()
	return instance
}

func (a *aspectRuntimeContextHostAPI) Register() {
	// tx contexts
	txCtx := datactx.NewTxContext(a.aspectRuntimeContext.EthTxContext, a.aspectRuntimeContext.CosmosContext)
	a.execMap[aspctx.TxType] = txCtx.ValueLoader(aspctx.TxType)
	a.execMap[aspctx.TxChainId] = txCtx.ValueLoader(aspctx.TxChainId)
	a.execMap[aspctx.TxAccessList] = txCtx.ValueLoader(aspctx.TxAccessList)
	a.execMap[aspctx.TxNonce] = txCtx.ValueLoader(aspctx.TxNonce)
	a.execMap[aspctx.TxGasPrice] = txCtx.ValueLoader(aspctx.TxGasPrice)
	a.execMap[aspctx.TxGas] = txCtx.ValueLoader(aspctx.TxGas)
	a.execMap[aspctx.TxGasTipCap] = txCtx.ValueLoader(aspctx.TxGasTipCap)
	a.execMap[aspctx.TxGasFeeCap] = txCtx.ValueLoader(aspctx.TxGasFeeCap)
	a.execMap[aspctx.TxTo] = txCtx.ValueLoader(aspctx.TxTo)
	a.execMap[aspctx.TxValue] = txCtx.ValueLoader(aspctx.TxValue)
	a.execMap[aspctx.TxData] = txCtx.ValueLoader(aspctx.TxData)
	a.execMap[aspctx.TxBytes] = txCtx.ValueLoader(aspctx.TxBytes)
	a.execMap[aspctx.TxHash] = txCtx.ValueLoader(aspctx.TxHash)
	a.execMap[aspctx.TxUnsignedBytes] = txCtx.ValueLoader(aspctx.TxUnsignedBytes)
	a.execMap[aspctx.TxUnsignedHash] = txCtx.ValueLoader(aspctx.TxUnsignedHash)
	a.execMap[aspctx.TxSigV] = txCtx.ValueLoader(aspctx.TxSigV)
	a.execMap[aspctx.TxSigR] = txCtx.ValueLoader(aspctx.TxSigR)
	a.execMap[aspctx.TxSigS] = txCtx.ValueLoader(aspctx.TxSigS)
	a.execMap[aspctx.TxFrom] = txCtx.ValueLoader(aspctx.TxFrom)
	a.execMap[aspctx.TxIndex] = txCtx.ValueLoader(aspctx.TxIndex)

	// msg contexts
	msgCtx := datactx.NewMessageContext(a.aspectRuntimeContext.EthTxContext, a.aspectRuntimeContext.CosmosContext)
	a.execMap[aspctx.MsgIndex] = msgCtx.ValueLoader(aspctx.MsgIndex)
	a.execMap[aspctx.MsgFrom] = msgCtx.ValueLoader(aspctx.MsgFrom)
	a.execMap[aspctx.MsgTo] = msgCtx.ValueLoader(aspctx.MsgTo)
	a.execMap[aspctx.MsgValue] = msgCtx.ValueLoader(aspctx.MsgValue)
	a.execMap[aspctx.MsgInput] = msgCtx.ValueLoader(aspctx.MsgInput)
	a.execMap[aspctx.MsgGas] = msgCtx.ValueLoader(aspctx.MsgGas)
	a.execMap[aspctx.MsgResultRet] = msgCtx.ValueLoader(aspctx.MsgResultRet)
	a.execMap[aspctx.MsgResultGasUsed] = msgCtx.ValueLoader(aspctx.MsgResultGasUsed)
	a.execMap[aspctx.MsgResultError] = msgCtx.ValueLoader(aspctx.MsgResultError)

	// receipt contexts
	receiptCtx := datactx.NewReceiptContext(a.aspectRuntimeContext.EthTxContext, a.aspectRuntimeContext.CosmosContext)
	a.execMap[aspctx.ReceiptStatus] = receiptCtx.ValueLoader(aspctx.ReceiptStatus)
	a.execMap[aspctx.ReceiptLogs] = receiptCtx.ValueLoader(aspctx.ReceiptLogs)
	a.execMap[aspctx.ReceiptGasUsed] = receiptCtx.ValueLoader(aspctx.ReceiptGasUsed)
	a.execMap[aspctx.ReceiptCumulativeGasUsed] = receiptCtx.ValueLoader(aspctx.ReceiptCumulativeGasUsed)
	a.execMap[aspctx.ReceiptBloom] = receiptCtx.ValueLoader(aspctx.ReceiptBloom)

	// block contexts
	blockCtx := datactx.NewBlockContext(a.aspectRuntimeContext)
	a.execMap[aspctx.BlockHeaderParentHash] = blockCtx.ValueLoader(aspctx.BlockHeaderParentHash)
	a.execMap[aspctx.BlockHeaderMiner] = blockCtx.ValueLoader(aspctx.BlockHeaderMiner)
	a.execMap[aspctx.BlockHeaderTransactionsRoot] = blockCtx.ValueLoader(aspctx.BlockHeaderTransactionsRoot)
	a.execMap[aspctx.BlockHeaderNumber] = blockCtx.ValueLoader(aspctx.BlockHeaderNumber)
	a.execMap[aspctx.BlockHeaderTimestamp] = blockCtx.ValueLoader(aspctx.BlockHeaderTimestamp)

	// env contexts
	envCtx := datactx.NewEnvContext(a.aspectRuntimeContext, evmKeeper)
	a.execMap[aspctx.EnvExtraEIPs] = envCtx.ValueLoader(aspctx.EnvExtraEIPs)
	a.execMap[aspctx.EnvEnableCreate] = envCtx.ValueLoader(aspctx.EnvEnableCreate)
	a.execMap[aspctx.EnvEnableCall] = envCtx.ValueLoader(aspctx.EnvEnableCall)
	a.execMap[aspctx.EnvAllowUnprotectedTxs] = envCtx.ValueLoader(aspctx.EnvAllowUnprotectedTxs)
	a.execMap[aspctx.EnvChainChainId] = envCtx.ValueLoader(aspctx.EnvChainChainId)
	a.execMap[aspctx.EnvChainHomesteadBlock] = envCtx.ValueLoader(aspctx.EnvChainHomesteadBlock)
	a.execMap[aspctx.EnvChainDaoForkBlock] = envCtx.ValueLoader(aspctx.EnvChainDaoForkBlock)
	a.execMap[aspctx.EnvChainDaoForkSupport] = envCtx.ValueLoader(aspctx.EnvChainDaoForkSupport)
	a.execMap[aspctx.EnvChainEip150Block] = envCtx.ValueLoader(aspctx.EnvChainEip150Block)
	a.execMap[aspctx.EnvChainEip155Block] = envCtx.ValueLoader(aspctx.EnvChainEip155Block)
	a.execMap[aspctx.EnvChainEip158Block] = envCtx.ValueLoader(aspctx.EnvChainEip158Block)
	a.execMap[aspctx.EnvChainByzantiumBlock] = envCtx.ValueLoader(aspctx.EnvChainByzantiumBlock)
	a.execMap[aspctx.EnvChainConstantinopleBlock] = envCtx.ValueLoader(aspctx.EnvChainConstantinopleBlock)
	a.execMap[aspctx.EnvChainPetersburgBlock] = envCtx.ValueLoader(aspctx.EnvChainPetersburgBlock)
	a.execMap[aspctx.EnvChainIstanbulBlock] = envCtx.ValueLoader(aspctx.EnvChainIstanbulBlock)
	a.execMap[aspctx.EnvChainMuirGlacierBlock] = envCtx.ValueLoader(aspctx.EnvChainMuirGlacierBlock)
	a.execMap[aspctx.EnvChainBerlinBlock] = envCtx.ValueLoader(aspctx.EnvChainBerlinBlock)
	a.execMap[aspctx.EnvChainLondonBlock] = envCtx.ValueLoader(aspctx.EnvChainLondonBlock)
	a.execMap[aspctx.EnvChainArrowGlacierBlock] = envCtx.ValueLoader(aspctx.EnvChainArrowGlacierBlock)
	a.execMap[aspctx.EnvChainGrayGlacierBlock] = envCtx.ValueLoader(aspctx.EnvChainGrayGlacierBlock)
	a.execMap[aspctx.EnvChainMergeNetSplitBlock] = envCtx.ValueLoader(aspctx.EnvChainMergeNetSplitBlock)
	a.execMap[aspctx.EnvChainShanghaiTime] = envCtx.ValueLoader(aspctx.EnvChainShanghaiTime)
	a.execMap[aspctx.EnvChainCancunTime] = envCtx.ValueLoader(aspctx.EnvChainCancunTime)
	a.execMap[aspctx.EnvChainPragueTime] = envCtx.ValueLoader(aspctx.EnvChainPragueTime)
	a.execMap[aspctx.EnvConsensusParamsBlockMaxGas] = envCtx.ValueLoader(aspctx.EnvConsensusParamsBlockMaxGas)
	a.execMap[aspctx.EnvConsensusParamsBlockMaxBytes] = envCtx.ValueLoader(aspctx.EnvConsensusParamsBlockMaxBytes)
	a.execMap[aspctx.EnvConsensusParamsEvidenceMaxAgeDuration] = envCtx.ValueLoader(aspctx.EnvConsensusParamsEvidenceMaxAgeDuration)
	a.execMap[aspctx.EnvConsensusParamsEvidenceMaxAgeNumBlocks] = envCtx.ValueLoader(aspctx.EnvConsensusParamsEvidenceMaxAgeNumBlocks)
	a.execMap[aspctx.EnvConsensusParamsEvidenceMaxBytes] = envCtx.ValueLoader(aspctx.EnvConsensusParamsEvidenceMaxBytes)
	a.execMap[aspctx.EnvConsensusParamsValidatorPubKeyTypes] = envCtx.ValueLoader(aspctx.EnvConsensusParamsValidatorPubKeyTypes)
	a.execMap[aspctx.EnvConsensusParamsAppVersion] = envCtx.ValueLoader(aspctx.EnvConsensusParamsAppVersion)

	// aspect contexts
	aspectCtx := datactx.NewAspectContext()
	a.execMap[aspctx.AspectId] = aspectCtx.ValueLoader(aspctx.AspectId)
	a.execMap[aspctx.AspectVersion] = aspectCtx.ValueLoader(aspctx.AspectVersion)

	// isCall context
	a.execMap[aspctx.IsCall] = func(ctx *asptypes.RunnerContext) ([]byte, error) {
		// verify tx can only be triggered when it is a transaction, so we can safely assume that
		// if the point is verify-tx, it must be a transaction.
		// otherwise if the request is not being committed, we can assume that it is a call.
		result := ctx.Point != string(asptypes.VERIFY_TX) &&
			!a.aspectRuntimeContext.EthTxContext().Commit()
		return proto.Marshal(&asptypes.BoolData{Data: &result})
	}
}

func (a *aspectRuntimeContextHostAPI) Get(ctx *asptypes.RunnerContext, key string) ([]byte, error) {
	joinPointCtxKeyConstraints, ok := ctxKeyConstraints[asptypes.PointCut(ctx.Point)]
	if !ok || !joinPointCtxKeyConstraints.Contains(key) {
		return nil, fmt.Errorf("key %s is not available at join point %s", key, ctx.Point)
	}

	res, err := a.execMap[key](ctx)
	if err != nil {
		// error returned here usually should not happen, but if it does, it is potentially a bug or
		// something wrong with the node, we need to panic here to avoid any potential security issues.
		panic(err)
	}

	return res, nil
}

func GetAspectRuntimeContextHostInstance(ctx context.Context) (asptypes.RuntimeContextHostAPI, error) {
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetAspectRuntimeContextHostInstance: unwrap AspectRuntimeContext failed")
	}
	return newAspectRuntimeContextHostAPI(aspectCtx), nil
}
