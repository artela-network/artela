package evm

import (
	"github.com/artela-network/artela/x/evm/artela/types"

	"github.com/artela-network/artela/x/evm/keeper"

	storetypes "cosmossdk.io/store/types"
	cosmos "github.com/cosmos/cosmos-sdk/types"

	ethereum "github.com/ethereum/go-ethereum/core/types"
)

// BeginBlock sets the cosmos Context and EIP155 chain id to the Keeper.
func BeginBlock(_ cosmos.Context, k *keeper.Keeper) {

	// Aspect Runtime Context Lifecycle: create and store ExtBlockContext
	// due to the design of the block context in Cosmos SDK,
	// the extBlockCtx cannot be saved directly to the context of the deliver state
	// using code like ctx = ctx.WithValue(artelatypes.ExtBlockContextKey, extBlockCtx).
	// Instead, it suggests saving it to the keeper.
	k.BlockContext = types.NewEthBlockContextFromABCIBeginBlockReq()
}

// EndBlock also retrieves the bloom filter value from the transient store and commits it to the
// KVStore. The EVM end block logic doesn't update the validator set, thus it returns
// an empty slice.
func EndBlock(ctx cosmos.Context, k *keeper.Keeper) {
	// Aspect Runtime Context Lifecycle: destory ExtBlockContext
	k.BlockContext = nil

	// Gas costs are handled within msg handler so costs should be ignored
	infCtx := ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())

	bloom := ethereum.BytesToBloom(k.GetBlockBloomTransient(infCtx).Bytes())
	k.EmitBlockBloomEvent(infCtx, bloom)
}
