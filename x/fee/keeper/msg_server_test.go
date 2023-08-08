package keeper_test

import (
	"context"
	"testing"

	keepertest "artela/testutil/keeper"
	"artela/x/fee/keeper"
	"artela/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.FeeKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
