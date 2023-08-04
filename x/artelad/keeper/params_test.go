package keeper_test

import (
	"testing"

	testkeeper "artelad/testutil/keeper"
	"artelad/x/artelad/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.ArteladKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
