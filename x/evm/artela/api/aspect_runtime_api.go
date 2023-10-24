package api

import (
	asptypes "github.com/artela-network/aspect-core/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
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
	getEthTxContext    func() *types.EthTxContext
	getAspectContext   func() *types.AspectContext
	getExtBlockContext func() *types.ExtBlockContext
	getCtxByHeight     func(height int64, prove bool) (sdk.Context, error)
	execMap            map[string]datactx.Executor
	aspectStateHostApi *aspectStateHostApi
	// nolint
	app *baseapp.BaseApp
}

func NewAspectRuntime(storeKey storetypes.StoreKey, getEthTxContext func() *types.EthTxContext,
	getAspectContext func() *types.AspectContext,
	getExtBlockContext func() *types.ExtBlockContext,
	getCtxByHeight func(height int64, prove bool) (sdk.Context, error), app *baseapp.BaseApp,
) { //nolint:gofumpt

	aspectRuntimeInstance = &aspectRuntimeHostApi{
		getEthTxContext:    getEthTxContext,
		getAspectContext:   getAspectContext,
		getCtxByHeight:     getCtxByHeight,
		getExtBlockContext: getExtBlockContext,
		execMap:            make(map[string]datactx.Executor),
		aspectStateHostApi: NewAspectState(app, storeKey, getCtxByHeight),
	}
	aspectRuntimeInstance.Register()
}

func (k *aspectRuntimeHostApi) Register() {
	// contexts
	k.execMap[asptypes.TxAspectContext] = datactx.NewTxAspectContent(k.getAspectContext)
	k.execMap[asptypes.TxStateChanges] = datactx.NewStateChanges(k.getEthTxContext)
	k.execMap[asptypes.TxExtProperties] = datactx.NewExtProperties(k.getEthTxContext)
	k.execMap[asptypes.TxContent] = datactx.NewTxContent(k.getEthTxContext)
	k.execMap[asptypes.TxCallTree] = datactx.NewTxCallTree(k.getEthTxContext)
	k.execMap[asptypes.TxReceipt] = datactx.NewTxReceipt(k.getEthTxContext)
	k.execMap[asptypes.TxGasMeter] = datactx.NewTxGasMeter()
	k.execMap[asptypes.EnvConsensusParams] = datactx.NewEnvConsParams()
	k.execMap[asptypes.EnvChainConfig] = datactx.NewEnvChainConfig(k.getEthTxContext)
	k.execMap[asptypes.EnvEvmParams] = datactx.NewEnvEvmParams(k.getEthTxContext)
	k.execMap[asptypes.EnvBaseInfo] = datactx.NewEnvBaseInfo(k.getEthTxContext)
	k.execMap[asptypes.BlockHeader] = datactx.NewEthBlockHeader()
	k.execMap[asptypes.BlockGasMeter] = datactx.NewEthBlockGasMeter()
	k.execMap[asptypes.BlockMinGasPrice] = datactx.NewBlockMinGasPrice()
	k.execMap[asptypes.BlockLastCommit] = datactx.NewBlockLastCommitInfo(k.getExtBlockContext)
	k.execMap[asptypes.BlockTxs] = datactx.NewEthBlockTxs(k.getExtBlockContext)
}

func GetRuntimeInstance() (asptypes.RuntimeHostApi, error) {
	if aspectRuntimeInstance == nil {
		return nil, errors.New("RuntimeHostApi not init")
	}
	return aspectRuntimeInstance, nil
}

func (a aspectRuntimeHostApi) SetAspectContext(ctx *asptypes.RunnerContext, key, value string) bool {
	a.getAspectContext().Add(ctx.AspectId.String(), key, value)
	return true
}

func (a *aspectRuntimeHostApi) GetContext(ctx *asptypes.RunnerContext, key string) *asptypes.ContextQueryResponse {
	has, ctxKey, params := asptypes.HasContextKey(key)
	if has {
		if matcher, ok := a.execMap[ctxKey]; ok {
			sdkCtx, _ := a.getCtxByHeight(ctx.BlockNumber, true)
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
		a.getAspectContext().Add(ctx.AspectId.String(), set.GetKey(), set.GetValue())
	}
	if set.NameSpace == asptypes.SetNameSpace_SetAspectState {
		a.aspectStateHostApi.SetAspectState(ctx, set.GetKey(), set.GetValue())
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
		property := a.aspectStateHostApi.GetProperty(ctx, keyData.Data)
		valueData := &asptypes.StringData{Data: property}
		response.SetData(valueData)
	}
	if query.NameSpace == asptypes.QueryNameSpace_QueryAspectState {
		state := a.aspectStateHostApi.GetAspectState(ctx, keyData.Data)
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
		return a.getAspectContext().Remove(ctx.AspectId.String(), keyData.Data)
	}
	if query.NameSpace == asptypes.RemoveNameSpace_RemoveAspectState {
		return a.aspectStateHostApi.RemoveAspectState(ctx, keyData.Data)
	}
	return false
}
