package keeper

import (
	"artelad/x/evm/types"
)

var _ types.QueryServer = Keeper{}
