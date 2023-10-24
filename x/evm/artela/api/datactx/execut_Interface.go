package datactx

import (
	artelatypes "github.com/artela-network/aspect-core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Executor interface {
	Execute(sdkCtx sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse
}
