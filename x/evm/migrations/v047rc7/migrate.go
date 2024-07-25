package v047rc7

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateStore(_ sdk.Context, _ storetypes.StoreKey, _ codec.BinaryCodec, logger log.Logger) error {
	logger.Error("v047rc7 Migrate Store is updating")
	return nil
}
