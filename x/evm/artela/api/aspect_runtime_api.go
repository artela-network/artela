package api

import (
	"github.com/artela-network/artela/x/evm/artela/contract"
	asptypes "github.com/artela-network/aspect-core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"

	"github.com/artela-network/artela/x/evm/artela/api/datactx"
	"github.com/artela-network/artela/x/evm/artela/types"
)

var (
	_                     asptypes.RuntimeHostApi = (*aspectRuntimeHostApi)(nil)
	aspectRuntimeInstance *aspectRuntimeHostApi
)

type aspectRuntimeHostApi struct {
	aspectRuntimeContext *types.AspectRuntimeContext
	getEthereumConfig    func(ctx sdk.Context) *ethparams.ChainConfig
	getCtxByHeight       types.ContextBuilder
	execMap              map[string]datactx.Executor
}

func NewAspectRuntime(
	aspectRuntimeContext *types.AspectRuntimeContext,
	getCtxByHeight types.ContextBuilder,
	getEthereumConfig func(ctx sdk.Context) *ethparams.ChainConfig,
) { //nolint:gofumpt
	aspectRuntimeInstance = &aspectRuntimeHostApi{
		aspectRuntimeContext: aspectRuntimeContext,
		getCtxByHeight:       getCtxByHeight,
		getEthereumConfig:    getEthereumConfig,
		execMap:              make(map[string]datactx.Executor),
	}
	aspectRuntimeInstance.Register()
}

func (k *aspectRuntimeHostApi) Register() {
	// contexts
	k.execMap[asptypes.TxAspectContext] = datactx.NewTxAspectContent(k.aspectRuntimeContext.AspectContext)
	k.execMap[asptypes.TxStateChanges] = datactx.NewStateChanges(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.TxExtProperties] = datactx.NewExtProperties(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.TxContent] = datactx.NewTxContent(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.TxMsgHash] = datactx.NewTxMsgHash(k.aspectRuntimeContext.EthTxContext, k.getEthereumConfig)
	k.execMap[asptypes.TxCallTree] = datactx.NewTxCallTree(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.TxReceipt] = datactx.NewTxReceipt(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.TxGasMeter] = datactx.NewTxGasMeter()
	k.execMap[asptypes.EnvConsensusParams] = datactx.NewEnvConsParams()
	k.execMap[asptypes.EnvChainConfig] = datactx.NewEnvChainConfig(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.EnvEvmParams] = datactx.NewEnvEvmParams(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.EnvBaseInfo] = datactx.NewEnvBaseInfo(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.BlockHeader] = datactx.NewEthBlockHeader(k.getCtxByHeight)
	k.execMap[asptypes.BlockGasMeter] = datactx.NewEthBlockGasMeter()
	k.execMap[asptypes.BlockMinGasPrice] = datactx.NewBlockMinGasPrice()
	k.execMap[asptypes.BlockLastCommit] = datactx.NewBlockLastCommitInfo(k.aspectRuntimeContext.ExtBlockContext)
	k.execMap[asptypes.BlockTxs] = datactx.NewEthBlockTxs(k.aspectRuntimeContext.ExtBlockContext)
}

func GetRuntimeInstance() (asptypes.RuntimeHostApi, error) {
	if aspectRuntimeInstance == nil {
		return nil, errors.New("RuntimeHostApi not init")
	}
	return aspectRuntimeInstance, nil
}

func (a aspectRuntimeHostApi) SetAspectContext(ctx *asptypes.RunnerContext, key, value string) bool {
	a.aspectRuntimeContext.AspectContext().Add(ctx.AspectId.String(), key, value)
	return true
}

func (a *aspectRuntimeHostApi) GetContext(ctx *asptypes.RunnerContext, key string) *asptypes.ContextQueryResponse {
	has, ctxKey, params := asptypes.HasContextKey(key)
	if has {
		if matcher, ok := a.execMap[ctxKey]; ok {
			sdkCtx, err := a.getCtxByHeight(ctx.BlockNumber, true)
			if err != nil {
				sdkCtx, err = a.getCtxByHeight(ctx.BlockNumber-1, true)
				if err != nil {
					return asptypes.NewContextQueryResponse(false, "Failed to retrieve chain context")
				}
			}
			execute := matcher.Execute(sdkCtx, ctx, params)
			if execute == nil {
				return asptypes.NewContextQueryResponse(false, "Get fail.")
			}
			return execute
		}
	}
	return asptypes.NewContextQueryResponse(false, "not supported key.")
}

func (a *aspectRuntimeHostApi) Set(ctx *asptypes.RunnerContext, set *asptypes.ContextSetRequest) bool {
	if set.NameSpace == asptypes.SetNameSpace_SetAspectContext {
		a.aspectRuntimeContext.AspectContext().Add(ctx.AspectId.String(), set.GetKey(), set.GetValue())
	}
	if set.NameSpace == asptypes.SetNameSpace_SetAspectState {
		a.aspectRuntimeContext.SetAspectState(ctx, set.GetKey(), set.GetValue())
	}
	return true
}

func (a *aspectRuntimeHostApi) Query(ctx *asptypes.RunnerContext, query *asptypes.ContextQueryRequest) *asptypes.ContextQueryResponse {
	keyData := &asptypes.StringData{}
	err := query.Query.UnmarshalTo(keyData)
	if err != nil {
		return asptypes.NewContextQueryResponse(false, "input unmarshal error.")
	}

	response := asptypes.NewContextQueryResponse(true, "success.")
	if query.NameSpace == asptypes.QueryNameSpace_QueryAspectProperty {
		codeStore := contract.NewAspectStore(a.aspectRuntimeContext.StoreKey())
		property := codeStore.GetAspectPropertyValue(a.aspectRuntimeContext.CosmosContext(), *ctx.AspectId, keyData.Data)
		valueData := &asptypes.StringData{Data: property}
		response.SetData(valueData)
	}
	if query.NameSpace == asptypes.QueryNameSpace_QueryAspectState {
		state := a.aspectRuntimeContext.GetAspectState(ctx, keyData.Data)
		valueData := &asptypes.StringData{Data: state}
		response.SetData(valueData)
	}
	return response
}

func (a *aspectRuntimeHostApi) Remove(ctx *asptypes.RunnerContext, query *asptypes.ContextRemoveRequest) bool {
	keyData := &asptypes.StringData{}
	err := query.Query.UnmarshalTo(keyData)
	if err != nil {
		return false
	}
	if query.NameSpace == asptypes.RemoveNameSpace_RemoveAspectContext {
		return a.aspectRuntimeContext.AspectContext().Remove(ctx.AspectId.String(), keyData.Data)
	}
	if query.NameSpace == asptypes.RemoveNameSpace_RemoveAspectState {
		return a.aspectRuntimeContext.RemoveAspectState(ctx, keyData.Data)
	}
	return false
}
