package keeper

import (
	"artela/x/evm/types"
)

var _ types.QueryServer = Keeper{}
