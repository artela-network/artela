package api

import (
	asptypes "github.com/artela-network/aspect-core/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	execMap              map[string]datactx.Executor
	getDeliverStateCtx   func() (sdk.Context, error)
	getStoreByHeight     func(height int64, prefixKey string) (prefix.Store, error)
}

func NewAspectRuntime(aspectRuntimeContext *types.AspectRuntimeContext,
	getDeliverStateCtx func() (sdk.Context, error),
	getStoreByHeight func(height int64, prefixKey string) (prefix.Store, error),
) { //nolint:gofumpt

	aspectRuntimeInstance = &aspectRuntimeHostApi{
		aspectRuntimeContext: aspectRuntimeContext,
		getDeliverStateCtx:   getDeliverStateCtx,
		execMap:              make(map[string]datactx.Executor),
		getStoreByHeight:     getStoreByHeight,
	}
	aspectRuntimeInstance.Register()
}

func (k *aspectRuntimeHostApi) Register() {
	// contexts
	k.execMap[asptypes.TxAspectContext] = datactx.NewTxAspectContent(k.aspectRuntimeContext.AspectContext)
	k.execMap[asptypes.TxStateChanges] = datactx.NewStateChanges(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.TxExtProperties] = datactx.NewExtProperties(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.TxContent] = datactx.NewTxContent(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.TxCallTree] = datactx.NewTxCallTree(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.TxReceipt] = datactx.NewTxReceipt(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.TxGasMeter] = datactx.NewTxGasMeter()
	k.execMap[asptypes.EnvConsensusParams] = datactx.NewEnvConsParams()
	k.execMap[asptypes.EnvChainConfig] = datactx.NewEnvChainConfig(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.EnvEvmParams] = datactx.NewEnvEvmParams(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.EnvBaseInfo] = datactx.NewEnvBaseInfo(k.aspectRuntimeContext.EthTxContext)
	k.execMap[asptypes.BlockHeader] = datactx.NewEthBlockHeader()
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
			sdkCtx, err := a.getDeliverStateCtx()
			if err != nil {
				return asptypes.NewContextQueryResponse(false, err.Error())
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
		a.aspectRuntimeContext.AspectState().SetAspectState(ctx, set.GetKey(), set.GetValue())
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
		storeByHeight, ctxErr := a.getStoreByHeight(ctx.BlockNumber-1, types.AspectPropertyKeyPrefix)
		if ctxErr != nil {
			return asptypes.NewContextQueryResponse(false, "get context by blockHeight error.")
		}
		aspectPropertyKey := types.AspectPropertyKey(
			ctx.AspectId.Bytes(),
			[]byte(keyData.Data),
		)
		get := storeByHeight.Get(aspectPropertyKey)
		data := ""
		if len(get) > 0 {
			data = string(get)
		}
		valueData := &asptypes.StringData{Data: data}
		response.SetData(valueData)
	}
	if query.NameSpace == asptypes.QueryNameSpace_QueryAspectState {
		state := a.aspectRuntimeContext.AspectState().GetAspectState(ctx, keyData.Data)
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
		return a.aspectRuntimeContext.AspectState().RemoveAspectState(ctx, keyData.Data)
	}
	return false
}
