package evm

import (
	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/artela-network/artela/x/evm/keeper"

	cosmos "github.com/cosmos/cosmos-sdk/types"

	artelatypes "github.com/artela-network/artela/x/evm/artela/types"
	ethereum "github.com/ethereum/go-ethereum/core/types"
)

// BeginBlock sets the cosmos Context and EIP155 chain id to the Keeper.
func BeginBlock(ctx cosmos.Context, k *keeper.Keeper, beginBlock abci.RequestBeginBlock) {

	// Aspect Runtime Context Lifecycle: create and store ExtBlockContext
	// due to the design of the block context in Cosmos SDK,
	// the extBlockCtx cannot be saved directly to the context of the deliver state
	// using code like ctx = ctx.WithValue(artelatypes.ExtBlockContextKey, extBlockCtx).
	// Instead, it suggests saving it to the keeper.
	extBlockCtx := artelatypes.NewExtBlockContext()
	rpcContext := k.GetClientContext()
	extBlockCtx = extBlockCtx.WithBlockConfig(beginBlock.ByzantineValidators, beginBlock.LastCommitInfo, rpcContext)
	k.ExtBlockContext = extBlockCtx

	k.WithChainID(ctx)

	// --------aspect OnBlockInitialize start ---  //
	/*header := types.ConvertEthBlockHeader(ctx.BlockHeader())
	request := &asptypes.EthBlockAspect{Header: header, GasInfo: &asptypes.GasInfo{
		GasWanted: 0,
		GasUsed:   0,
		Gas:       0,
	}}

	receive := djpm.AspectInstance().OnBlockInitialize(request)
	hasErr, receiveErr := receive.HasErr()
	if hasErr {
		ctx.Logger().Error("Aspect.OnBlockInitialize Return Error ", receiveErr.Error(), "height", request.Header.Number)
	}*/
	// --------aspect OnBlockInitialize end ---  //
}

// EndBlock also retrieves the bloom filter value from the transient store and commits it to the
// KVStore. The EVM end block logic doesn't update the validator set, thus it returns
// an empty slice.
func EndBlock(ctx cosmos.Context, k *keeper.Keeper, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	// Aspect Runtime Context Lifecycle: destory ExtBlockContext
	k.ExtBlockContext = nil

	// Gas costs are handled within msg handler so costs should be ignored
	infCtx := ctx.WithGasMeter(cosmos.NewInfiniteGasMeter())

	bloom := ethereum.BytesToBloom(k.GetBlockBloomTransient(infCtx).Bytes())
	k.EmitBlockBloomEvent(infCtx, bloom)

	// --------aspect OnBlockFinalize start ---  //
	/*header := types.ConvertEthBlockHeader(ctx.BlockHeader())
	request := &asptypes.EthBlockAspect{Header: header, GasInfo: &asptypes.GasInfo{
		GasWanted: 0,
		GasUsed:   0,
		Gas:       0,
	}}

	receive := djpm.AspectInstance().OnBlockFinalize(request)
	hasErr, receiveErr := receive.HasErr()
	if hasErr {
		ctx.Logger().Error("Aspect.OnBlockFinalize Return Error ", receiveErr.Error(), "height", request.Header.Number)
	}
	*/
	// --------aspect OnBlockFinalize end ---  //

	return []abci.ValidatorUpdate{}
}
