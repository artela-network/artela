package evm

import (
	"runtime"

	"github.com/artela-network/artela/x/evm/artela/types"
	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/artela-network/artela/x/evm/keeper"

	cosmos "github.com/cosmos/cosmos-sdk/types"

	ethereum "github.com/ethereum/go-ethereum/core/types"
)

// BeginBlock sets the cosmos Context and EIP155 chain id to the Keeper.
func BeginBlock(_ cosmos.Context, k *keeper.Keeper, beginBlock abci.RequestBeginBlock) {
	// We manually call runtime.GC in every block because in certain cases,
	// the store objects created by wasmtime-go are not promptly released,
	// causing the memory allocated by the wasm engine, which consumes a significant
	// amount of resources, to linger in memory for extended periods. This may
	// lead to excessive memory consumption and failure to allocate new memory
	// on some machines.
	// In testing, the time spent by runtime.GC with STW (Stop-The-World) is
	// approximately 100 to 200 microseconds, with very limited impact on blocks.
	runtime.GC()

	// Aspect Runtime Context Lifecycle: create and store ExtBlockContext
	// due to the design of the block context in Cosmos SDK,
	// the extBlockCtx cannot be saved directly to the context of the deliver state
	// using code like ctx = ctx.WithValue(artelatypes.ExtBlockContextKey, extBlockCtx).
	// Instead, it suggests saving it to the keeper.
	k.BlockContext = types.NewEthBlockContextFromABCIBeginBlockReq(beginBlock)
}

// EndBlock also retrieves the bloom filter value from the transient store and commits it to the
// KVStore. The EVM end block logic doesn't update the validator set, thus it returns
// an empty slice.
func EndBlock(ctx cosmos.Context, k *keeper.Keeper, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	// Aspect Runtime Context Lifecycle: destory ExtBlockContext
	k.BlockContext = nil

	// Gas costs are handled within msg handler so costs should be ignored
	infCtx := ctx.WithGasMeter(cosmos.NewInfiniteGasMeter())

	bloom := ethereum.BytesToBloom(k.GetBlockBloomTransient(infCtx).Bytes())
	k.EmitBlockBloomEvent(infCtx, bloom)

	return []abci.ValidatorUpdate{}
}
