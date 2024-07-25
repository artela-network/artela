package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/artela-network/artela/x/evm/migrations/v047rc7"
	"github.com/artela-network/artela/x/evm/migrations/v048rc8"
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

// Migrate5to6 migrates the store from consensus version 5 to 6
func (m Migrator) Migrate5to6(ctx sdk.Context) error {
	return v047rc7.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc, m.keeper.logger)
}

// Migrate6to7 migrates the store from consensus version 6 to 7
func (m Migrator) Migrate6to7(ctx sdk.Context) error {
	return v048rc8.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc, m.keeper.logger)
}
