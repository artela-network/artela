package datactx

import (
	"github.com/artela-network/artela/x/evm/artela/types"
	artelatypes "github.com/artela-network/artelasdk/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EthBlockTxs struct {
	getExtBlockContext func() *types.ExtBlockContext
}

func NewEthBlockTxs(getExtBlockContext func() *types.ExtBlockContext) *EthBlockTxs {
	return &EthBlockTxs{getExtBlockContext: getExtBlockContext}
}
func (c EthBlockTxs) Execute(sdkCtx sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}

	txsInBlock, err := types.GetEthTxsInBlock(c.getExtBlockContext().RpcClient(), sdkCtx, ctx.BlockNumber, ctx.ContractAddr.String())
	if err != nil {
		return artelatypes.NewContextQueryResponse(false, err.Error())
	}
	ethTransactions := make([]*artelatypes.EthTransaction, len(txsInBlock))
	for i, tx := range txsInBlock {
		ethTransactions[i] = types.ConvertEthTransaction(sdkCtx, tx)
	}
	array := &artelatypes.EthTxArray{Tx: ethTransactions}
	response := artelatypes.NewContextQueryResponse(true, "success")
	response.SetData(array)
	return response
}

type EthBlockEvidence struct {
	getExtBlockContext func() *types.ExtBlockContext
}

func NewEthBlockEvidence(getExtBlockContext func() *types.ExtBlockContext) *EthBlockEvidence {
	return &EthBlockEvidence{getExtBlockContext: getExtBlockContext}
}
func (c EthBlockEvidence) Execute(sdkCtx sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}
	blockEvidence := c.getExtBlockContext().EvidenceList()
	if len(blockEvidence) == 0 {
		return artelatypes.NewContextQueryResponse(true, "empty data")
	}

	evidenceList := make([]*artelatypes.Evidence, len(blockEvidence))
	for i, evidence := range blockEvidence {
		evidenceList[i] = &artelatypes.Evidence{
			Type: artelatypes.EvidenceType(evidence.Type),
			Validator: &artelatypes.Validator{
				Address: evidence.Validator.Address,
				Power:   evidence.Validator.Power,
			},
			Height:           evidence.Height,
			Time:             evidence.Time.UnixMicro(),
			TotalVotingPower: evidence.TotalVotingPower,
		}
	}
	list := &artelatypes.EvidenceList{Evidences: evidenceList}
	response := artelatypes.NewContextQueryResponse(true, "success")
	response.SetData(list)
	return response
}

type EthBlockId struct {
	getExtBlockContext func() *types.ExtBlockContext
}

func NewEthBlockId(getExtBlockContext func() *types.ExtBlockContext) *EthBlockId {
	return &EthBlockId{getExtBlockContext: getExtBlockContext}
}
func (c EthBlockId) Execute(sdkCtx sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}
	block, err := c.getExtBlockContext().RpcClient().Client.Block(sdkCtx.Context(), &ctx.BlockNumber)
	if err != nil {
		return artelatypes.NewContextQueryResponse(true, err.Error())
	}
	bId := block.BlockID
	blockID := &artelatypes.BlockID{
		Hash: bId.Hash.Bytes(),
		PartSetHeader: &artelatypes.PartSetHeader{
			Total: bId.PartSetHeader.Total,
			Hash:  bId.PartSetHeader.Hash.Bytes(),
		},
	}
	response := artelatypes.NewContextQueryResponse(true, "success")
	response.SetData(blockID)
	return response
}

type EthBlockHeader struct {
}

func NewEthBlockHeader() *EthBlockHeader {
	return &EthBlockHeader{}
}

// getAspectContext data
func (c EthBlockHeader) Execute(sdkCtx sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {

	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}
	header := sdkCtx.BlockHeader()

	if len(header.GetAppHash()) == 0 {
		return artelatypes.NewContextQueryResponse(false, "get empty.")
	} else {
		blockHeader := types.ProtoToEthBlockHeader(&header)
		response := artelatypes.NewContextQueryResponse(true, "success")
		response.SetData(blockHeader)
		return response
	}

}

// block last commit info
type BlockLastCommitInfo struct {
	getExtBlockContext func() *types.ExtBlockContext
}

func NewBlockLastCommitInfo(getExtBlockContext func() *types.ExtBlockContext) *BlockLastCommitInfo {
	return &BlockLastCommitInfo{getExtBlockContext: getExtBlockContext}
}

func (c BlockLastCommitInfo) Execute(sdkCtx sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}

	contextQueryResponse := artelatypes.NewContextQueryResponse(true, "basic validate failed.")
	if c.getExtBlockContext() == nil {
		return contextQueryResponse
	}
	blockCtx := c.getExtBlockContext().LastCommitInfo()
	info := ConvertVoteInfos(blockCtx.Votes)
	commitInfo := &artelatypes.LastCommitInfo{
		Round: blockCtx.Round,
		Votes: info,
	}
	contextQueryResponse.GetResult().Message = "success"
	contextQueryResponse.SetData(commitInfo)
	return contextQueryResponse
}

func ConvertVoteInfos(infos []abci.VoteInfo) []*artelatypes.VoteInfo {
	if len(infos) == 0 {
		return nil
	}
	voteInfos := make([]*artelatypes.VoteInfo, len(infos))
	for i, info := range infos {
		voteInfos[i] = ConvertVoteInfo(info)
	}
	return voteInfos
}
func ConvertVoteInfo(info abci.VoteInfo) *artelatypes.VoteInfo {

	return &artelatypes.VoteInfo{
		Validator: &artelatypes.Validator{
			Address: info.Validator.Address,
			Power:   info.Validator.Power,
		},
		SignedLastBlock: info.SignedLastBlock,
	}
}

// block minGasPrice
type BlockMinGasPrice struct {
}

func NewBlockMinGasPrice() *BlockMinGasPrice {
	return &BlockMinGasPrice{}
}

// getAspectContext data
func (c BlockMinGasPrice) Execute(sdkContext sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}

	contextQueryResponse := artelatypes.NewContextQueryResponse(true, "basic validate failed.")
	if sdkContext.MinGasPrices() == nil {
		return contextQueryResponse
	}
	meter := sdkContext.MinGasPrices()
	//convert
	coins := make([]*artelatypes.DecCoin, meter.Len())
	sort := meter.Sort()
	for i, coin := range sort {
		coins[i] = &artelatypes.DecCoin{
			Denom:  coin.Denom,
			Amount: coin.Amount.String(),
		}
	}
	//set data
	gasMsg := &artelatypes.MinGasPrice{
		Coins: coins,
	}
	contextQueryResponse.GetResult().Message = "success"
	contextQueryResponse.SetData(gasMsg)
	return contextQueryResponse
}

// /block gas meter

type EthBlockGasMeter struct {
}

func NewEthBlockGasMeter() *EthBlockGasMeter {
	return &EthBlockGasMeter{}
}

// getAspectContext data
func (c EthBlockGasMeter) Execute(sdkContext sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}

	contextQueryResponse := artelatypes.NewContextQueryResponse(true, "basic validate failed.")
	if sdkContext.BlockGasMeter() == nil {
		return contextQueryResponse
	}
	meter := sdkContext.BlockGasMeter()

	//set data
	gasMsg := &artelatypes.GasMeter{
		GasConsumed:        meter.GasConsumed(),
		GasConsumedToLimit: meter.GasConsumedToLimit(),
		GasRemaining:       meter.GasRemaining(),
		Limit:              meter.Limit(),
	}
	contextQueryResponse.GetResult().Message = "success"
	contextQueryResponse.SetData(gasMsg)
	return contextQueryResponse

}
