package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/artela-network/artela/x/evm/migrations/v047rc7"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		keeper: keeper,
	}
}

// Migrate3to4 migrates the store from consensus version 3 to 4
func (m Migrator) Migrate5to6(ctx sdk.Context) error {
	return v047rc7.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc, m.keeper.logger)
}
