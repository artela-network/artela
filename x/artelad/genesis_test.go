package artelad_test

import (
	"testing"

	keepertest "artelad/testutil/keeper"
	"artelad/testutil/nullify"
	"artelad/x/artelad"
	"artelad/x/artelad/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ArteladKeeper(t)
	artelad.InitGenesis(ctx, *k, genesisState)
	got := artelad.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
