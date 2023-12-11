package api

import (
	"fmt"

	ethlog "github.com/ethereum/go-ethereum/log"

	artelatypes "github.com/artela-network/aspect-core/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/artela-network/artela/x/evm/artela/contract"
	"github.com/artela-network/artela/x/evm/artela/types"
	evmtypes "github.com/artela-network/artela/x/evm/types"
)

type aspectStateHostApi struct {
	storeKey           storetypes.StoreKey
	getDeliverStateCtx func() (sdk.Context, error)
	getCtxByHeight     func(height int64, prove bool) (sdk.Context, error)
}

func newAspectState(app *baseapp.BaseApp, storeKey storetypes.StoreKey) *aspectStateHostApi {
	return &aspectStateHostApi{
		storeKey:           storeKey,
		getDeliverStateCtx: app.DeliverStateCtx,
		getCtxByHeight:     app.CreateQueryContext,
	}
}

func (k *aspectStateHostApi) newPrefixStore(fixKey string) prefix.Store {
	sdkCtx, err := k.getDeliverStateCtx()
	if err != nil {
		sdkCtx, _ = k.getCtxByHeight(0, false)
	}
	return prefix.NewStore(sdkCtx.KVStore(k.storeKey), evmtypes.KeyPrefix(fixKey))
}

// func (k *aspectStateHostApi) newTransientStore(blockHeight int64, fixKey string) prefix.Store {
//	ctx, _ := k.getCtxByHeight(blockHeight, false)
//	return prefix.NewStore(ctx.TransientStore(k.storeKey), evmtypes.KeyPrefix(fixKey))
// }

func (k *aspectStateHostApi) GetAspectState(ctx *artelatypes.RunnerContext, key string) string {
	codeStore := k.newPrefixStore(types.AspectStateKeyPrefix)
	aspectPropertyKey := types.AspectArrayKey(
		ctx.AspectId.Bytes(),
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
		[]byte(key),
	)
	_, err := k.getDeliverStateCtx()
	if err != nil {
		return true
	}
	ethlog.Info(fmt.Sprintf("SetAspectState, ---aspectID:%s---, ---key:%s---, ---value:%s---, ctx from deliver: %t", ctx.AspectId.String(), key, value, err == nil))
	codeStore.Set(aspectPropertyKey, []byte(value))
	return true
}

// RemoveAspectState RemoveAspectState( key string) bool
func (k *aspectStateHostApi) RemoveAspectState(ctx *artelatypes.RunnerContext, key string) bool {
	codeStore := k.newPrefixStore(types.AspectStateKeyPrefix)
	aspectPropertyKey := types.AspectArrayKey(
		ctx.AspectId.Bytes(),
		[]byte(key),
	)

	ethlog.Info("RemoveAspectState, key:", key)
	codeStore.Delete(aspectPropertyKey)
	return true
}

// GetProperty GetProperty( key string) string
func (k *aspectStateHostApi) GetProperty(ctx *artelatypes.RunnerContext, key string) string {
	codeStore := contract.NewAspectStore(k.storeKey)
	sdkCtx, _ := k.getCtxByHeight(ctx.BlockNumber-1, false)
	return codeStore.GetAspectPropertyValue(sdkCtx, *ctx.AspectId, key)
}
