package evm

import (
	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/artela-network/artela/x/evm/keeper"

	cosmos "github.com/cosmos/cosmos-sdk/types"

	ethereum "github.com/ethereum/go-ethereum/core/types"
)

// BeginBlock sets the cosmos Context and EIP155 chain id to the Keeper.
func BeginBlock(ctx cosmos.Context, k *keeper.Keeper, beginBlock abci.RequestBeginBlock) {
	//k.Lock.Lock()
	//defer k.Lock.Unlock()

	k.GetAspectRuntimeContext().WithCosmosContext(ctx)
	k.GetAspectRuntimeContext().ExtBlockContext().WithEvidenceList(beginBlock.ByzantineValidators).WithLastCommit(beginBlock.LastCommitInfo).WithRpcClient(k.GetClientContext())
	k.WithChainID(ctx)

	// Create a Aspect State Object for OnBlockInitialize
	//k.GetAspectRuntimeContext().CreateStateObject(true, ctx.BlockHeight(), types.AspectStateBeginBlock)
	//defer k.GetAspectRuntimeContext().RefreshState(true, ctx.BlockHeight(), types.AspectStateBeginBlock)

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

	// delete all aspect state objects in the block
	defer k.GetAspectRuntimeContext().RefreshState(false, ctx.BlockHeight(), "")

	// Gas costs are handled within msg handler so costs should be ignored
	infCtx := ctx.WithGasMeter(cosmos.NewInfiniteGasMeter())

	bloom := ethereum.BytesToBloom(k.GetBlockBloomTransient(infCtx).Bytes())
	k.EmitBlockBloomEvent(infCtx, bloom)

	k.GetAspectRuntimeContext().WithCosmosContext(ctx)
	// Create a Aspect State Object for OnBlockFinalize
	//	k.GetAspectRuntimeContext().CreateStateObject(true, ctx.BlockHeight(), types.AspectStateEndBlock)

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

	// clear aspect  block context
	//	k.GetAspectRuntimeContext().RefreshState(true, ctx.BlockHeight(), types.AspectStateEndBlock)
	k.GetAspectRuntimeContext().ClearBlockContext()

	return []abci.ValidatorUpdate{}
}
