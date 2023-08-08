package fee_test

import (
	"testing"

	keepertest "artela/testutil/keeper"
	"artela/testutil/nullify"
	"artela/x/fee"
	"artela/x/fee/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.FeeKeeper(t)
	fee.InitGenesis(ctx, *k, genesisState)
	got := fee.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
