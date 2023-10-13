package api

import (
	"github.com/artela-network/artela/x/evm/artela/contract"
	"github.com/artela-network/artela/x/evm/artela/types"
	evmtypes "github.com/artela-network/artela/x/evm/types"
	artelatypes "github.com/artela-network/artelasdk/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type aspectStateHostApi struct {
	app            *baseapp.BaseApp
	storeKey       storetypes.StoreKey
	getCtxByHeight func(height int64, prove bool) (sdk.Context, error)
}

func NewAspectState(app *baseapp.BaseApp, storeKey storetypes.StoreKey, getCtxByHeight func(height int64, prove bool) (sdk.Context, error)) *aspectStateHostApi {
	return &aspectStateHostApi{
		app:            app,
		storeKey:       storeKey,
		getCtxByHeight: getCtxByHeight,
	}
}

func (k *aspectStateHostApi) newPrefixStore(fixKey string) prefix.Store {
	sdkCtx, err := k.app.DeliverStateCtx()
	if err != nil {
		sdkCtx, _ = k.getCtxByHeight(0, false)
	}
	return prefix.NewStore(sdkCtx.KVStore(k.storeKey), evmtypes.KeyPrefix(fixKey))
}

//func (k *aspectStateHostApi) newTransientStore(blockHeight int64, fixKey string) prefix.Store {
//	ctx, _ := k.getCtxByHeight(blockHeight, false)
//	return prefix.NewStore(ctx.TransientStore(k.storeKey), evmtypes.KeyPrefix(fixKey))
//}

func (k *aspectStateHostApi) GetAspectState(ctx *artelatypes.RunnerContext, key string) string {
	codeStore := k.newPrefixStore(types.AspectStateKeyPrefix)
	aspectPropertyKey := types.AspectArrayKey(
		ctx.AspectId.Bytes(),
		ctx.ContractAddr.Bytes(),
		[]byte(key),
	)
	get := codeStore.Get(aspectPropertyKey)
	return artelatypes.Ternary(get != nil, func() string {
		return string(get)
	}, "")
}

func (k *aspectStateHostApi) SetAspectState(ctx *artelatypes.RunnerContext, key, value string) bool {
	codeStore := k.newPrefixStore(types.AspectStateKeyPrefix)
	aspectPropertyKey := types.AspectArrayKey(
		ctx.AspectId.Bytes(),
		ctx.ContractAddr.Bytes(),
		[]byte(key),
	)
	codeStore.Set(aspectPropertyKey, []byte(value))
	return true
}

// RemoveAspectState RemoveAspectState( key string) bool
func (k *aspectStateHostApi) RemoveAspectState(ctx *artelatypes.RunnerContext, key string) bool {
	codeStore := k.newPrefixStore(types.AspectStateKeyPrefix)
	aspectPropertyKey := types.AspectArrayKey(
		ctx.AspectId.Bytes(),
		ctx.ContractAddr.Bytes(),
		[]byte(key),
	)

	codeStore.Delete(aspectPropertyKey)
	return true
}

// GetProperty GetProperty( key string) string
func (k *aspectStateHostApi) GetProperty(ctx *artelatypes.RunnerContext, key string) string {
	codeStore := contract.NewAspectStore(k.storeKey)
	sdkCtx, _ := k.getCtxByHeight(ctx.BlockNumber-1, false)
	return codeStore.GetAspectPropertyValue(sdkCtx, *ctx.AspectId, key)

}
