package datactx

import (
	artelatypes "github.com/artela-network/aspect-core/types"
	"math/big"
)

type ContextLoader func(ctx *artelatypes.RunnerContext) ([]byte, error)

type EVMKeeper interface {
	ChainID() *big.Int
}
