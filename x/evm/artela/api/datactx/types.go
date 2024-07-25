package datactx

import (
	"math/big"

	artelatypes "github.com/artela-network/aspect-core/types"
)

type ContextLoader func(ctx *artelatypes.RunnerContext) ([]byte, error)

type EVMKeeper interface {
	ChainID() *big.Int
}
