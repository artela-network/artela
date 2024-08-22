package rpc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strconv"

	sdkmath "cosmossdk.io/math"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/artela-network/artela-evm/vm"
	"github.com/artela-network/artela/ethereum/rpc/api"
	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	"github.com/artela-network/artela/ethereum/rpc/utils"
	"github.com/artela-network/artela/x/evm/txs"
	evmtypes "github.com/artela-network/artela/x/evm/types"
)

// Blockchain API

// GetProof returns an account object with proof and any storage proofs
func (b *BackendImpl) GetProof(address common.Address, storageKeys []string, blockNrOrHash rpctypes.BlockNumberOrHash) (*rpctypes.AccountResult, error) {
	numberOrHash := rpc.BlockNumberOrHash{
		BlockNumber:      (*rpc.BlockNumber)(blockNrOrHash.BlockNumber),
		BlockHash:        blockNrOrHash.BlockHash,
		RequireCanonical: false,
	}
	blockNum, err := b.blockNumberFromCosmos(numberOrHash)
	if err != nil {
		return nil, err
	}

	height := blockNum.Int64()

	_, err = b.CosmosBlockByNumber(blockNum)
	if err != nil {
		// the error message imitates geth behavior
		return nil, fmt.Errorf("block not valid, %v", err)
	}
	ctx := rpctypes.ContextWithHeight(height)

	// if the height is equal to zero, meaning the query condition of the block is either "pending" or "latest"
	if height == 0 {
		bn, err := b.BlockNumber()
		if err != nil {
			return nil, err
		}

		if bn > math.MaxInt64 {
			return nil, fmt.Errorf("not able to query block number greater than MaxInt64")
		}

		height = int64(bn) // #nosec G701 -- checked for int overflow already
	}

	clientCtx := b.clientCtx.WithHeight(height)

	// query storage proofs
	storageProofs := make([]rpctypes.StorageResult, len(storageKeys))

	for i, key := range storageKeys {
		hexKey := common.HexToHash(key)
		valueBz, proof, err := b.queryClient.GetProof(clientCtx, evmtypes.StoreKey, evmtypes.StateKey(address, hexKey.Bytes()))
		if err != nil {
			return nil, err
		}

		storageProofs[i] = rpctypes.StorageResult{
			Key:   key,
			Value: (*hexutil.Big)(new(big.Int).SetBytes(valueBz)),
			Proof: utils.GetHexProofs(proof),
		}
	}

	// query EVM account
	req := &txs.QueryAccountRequest{
		Address: address.String(),
	}

	res, err := b.queryClient.Account(ctx, req)
	if err != nil {
		return nil, err
	}

	// query account proofs
	accountKey := authtypes.AddressStoreKey(sdktypes.AccAddress(address.Bytes()))
	_, proof, err := b.queryClient.GetProof(clientCtx, authtypes.StoreKey, accountKey)
	if err != nil {
		return nil, err
	}

	balance, ok := sdkmath.NewIntFromString(res.Balance)
	if !ok {
		return nil, errors.New("invalid balance")
	}

	return &rpctypes.AccountResult{
		Address:      address,
		AccountProof: utils.GetHexProofs(proof),
		Balance:      (*hexutil.Big)(balance.BigInt()),
		CodeHash:     common.HexToHash(res.CodeHash),
		Nonce:        hexutil.Uint64(res.Nonce),
		StorageHash:  common.Hash{}, // NOTE: Evmos doesn't have a storage hash. TODO: implement?
		StorageProof: storageProofs,
	}, nil
}

func (b *BackendImpl) DoCall(args rpctypes.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash) (*txs.MsgEthereumTxResponse, error) {
	blockNum, err := b.blockNumberFromCosmos(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	bz, err := json.Marshal(&args)
	if err != nil {
		return nil, err
	}
	header, err := b.CosmosBlockByNumber(blockNum)
	if err != nil {
		// the error message imitates geth behavior
		return nil, errors.New("header not found")
	}

	req := txs.EthCallRequest{
		Args:            bz,
		GasCap:          b.RPCGasCap(),
		ProposerAddress: sdktypes.ConsAddress(header.Block.ProposerAddress),
		ChainId:         b.chainID.Int64(),
	}

	// From ContextWithHeight: if the provided height is 0,
	// it will return an empty context and the gRPC query will use
	// the latest block height for querying.
	ctx := rpctypes.ContextWithHeight(blockNum.Int64())
	timeout := b.RPCEVMTimeout()

	// Setup context so it may be canceled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	// Make sure the context is canceled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	res, err := b.queryClient.EthCall(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		if res.VmError != vm.ErrExecutionReverted.Error() {
			return nil, status.Error(codes.Internal, res.VmError)
		}
		return nil, evmtypes.NewExecErrorWithReason(res.Ret)
	}

	return res, nil
}

func (b *BackendImpl) EstimateGas(ctx context.Context, args rpctypes.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	blockNum := rpc.LatestBlockNumber
	if blockNrOrHash != nil {
		blockNum, _ = b.blockNumberFromCosmos(*blockNrOrHash)
	}

	bz, err := json.Marshal(&args)
	if err != nil {
		return 0, err
	}

	header, err := b.CosmosBlockByNumber(blockNum)
	if err != nil {
		// the error message imitates geth behavior
		return 0, errors.New("header not found")
	}

	req := txs.EthCallRequest{
		Args:            bz,
		GasCap:          b.RPCGasCap(),
		ProposerAddress: sdktypes.ConsAddress(header.Block.ProposerAddress),
		ChainId:         b.chainID.Int64(),
	}

	// From ContextWithHeight: if the provided height is 0,
	// it will return an empty context and the gRPC query will use
	// the latest block height for querying.
	res, err := b.queryClient.EstimateGas(rpctypes.ContextWithHeight(blockNum.Int64()), &req)
	if err != nil {
		return 0, err
	}
	return hexutil.Uint64(res.Gas), nil
}

func (b *BackendImpl) HeaderByNumber(_ context.Context, number rpc.BlockNumber) (*ethtypes.Header, error) {
	resBlock, err := b.CosmosBlockByNumber(number)
	if err != nil {
		return nil, err
	}

	if resBlock == nil {
		return nil, fmt.Errorf("block not found for height %d", number)
	}

	blockRes, err := b.CosmosBlockResultByNumber(&resBlock.Block.Height)
	if err != nil {
		return nil, fmt.Errorf("block result not found for height %d", resBlock.Block.Height)
	}

	bloom, err := b.blockBloom(blockRes)
	if err != nil {
		b.logger.Debug("HeaderByNumber BlockBloom failed", "height", resBlock.Block.Height)
	}

	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle the error for pruned node.
		b.logger.Error("failed to fetch Base Fee from prunned block. Check node prunning configuration", "height", resBlock.Block.Height, "error", err)
	}

	ethHeader := rpctypes.EthHeaderFromTendermint(resBlock.Block.Header, bloom, baseFee)
	return ethHeader, nil
}

func (b *BackendImpl) HeaderByHash(ctx context.Context, hash common.Hash) (*ethtypes.Header, error) {
	block, err := b.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, errors.New("block not found")
	}
	return block.Header(), nil
}

func (b *BackendImpl) HeaderByNumberOrHash(ctx context.Context,
	blockNrOrHash rpc.BlockNumberOrHash,
) (*ethtypes.Header, error) {
	block, err := b.ArtBlockByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, errors.New("block not found")
	}
	return block.Header(), nil
}

func (b *BackendImpl) CurrentBlock() *rpctypes.Block {
	block, _ := b.currentBlock()
	return block
}

func (b *BackendImpl) ArtBlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*rpctypes.Block, error) {
	blockNum, err := b.blockNumberFromCosmos(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	return b.ArtBlockByNumber(ctx, blockNum)
}

func (b *BackendImpl) CosmosBlockByHash(blockHash common.Hash) (*tmrpctypes.ResultBlock, error) {
	resBlock, err := b.clientCtx.Client.BlockByHash(b.ctx, blockHash.Bytes())
	if err != nil {
		return nil, err
	}

	if resBlock.Block == nil {
		return nil, fmt.Errorf("failed to query block for hash: %s", blockHash.Hex())
	}

	return resBlock, nil
}

func (b *BackendImpl) CosmosBlockByNumber(blockNum rpc.BlockNumber) (*tmrpctypes.ResultBlock, error) {
	height := blockNum.Int64()
	if height < 0 {
		// fetch the latest block number from the app state, more accurate than the tendermint block store state.
		n, err := b.BlockNumber()
		if err != nil {
			return nil, err
		}
		height = int64(n) // #nosec G701 -- checked for int overflow already
	} else if height == 0 {
		// since in cosmos chain we don't actually have the block 0 (which is the genesis in ethereum),
		// so we decide just return the earliest block (block 1) when the height is 0.
		height = 1
	}
	resBlock, err := b.clientCtx.Client.Block(b.ctx, &height)
	if err != nil {
		return nil, err
	}

	if resBlock.Block == nil {
		return nil, fmt.Errorf("failed to query block for blockNum: %d", blockNum.Int64())
	}

	return resBlock, nil
}

func (b *BackendImpl) StateAndHeaderByNumber(
	_ context.Context, number rpc.BlockNumber,
) (*state.StateDB, *ethtypes.Header, error) {
	return nil, nil, errors.New("StateAndHeaderByNumber is not implemented")
}

func (b *BackendImpl) StateAndHeaderByNumberOrHash(
	_ context.Context, _ rpc.BlockNumberOrHash,
) (*state.StateDB, *ethtypes.Header, error) {
	return nil, nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *BackendImpl) GetEVM(_ context.Context, _ *core.Message, _ *state.StateDB,
	_ *ethtypes.Header, _ *vm.Config, _ *vm.BlockContext,
) (*vm.EVM, func() error) {
	return nil, func() error {
		return errors.New("GetEVM is not impelemted")
	}
}

func (b *BackendImpl) GetCode(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	blockNum, err := b.blockNumberFromCosmos(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	req := &txs.QueryCodeRequest{
		Address: address.String(),
	}

	res, err := b.queryClient.Code(rpctypes.ContextWithHeight(blockNum.Int64()), req)
	if err != nil {
		return nil, err
	}

	return res.Code, nil
}

// GetStorageAt returns the contract storage at the given address, block number, and key.
func (b *BackendImpl) GetStorageAt(address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	blockNum, err := b.blockNumberFromCosmos(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	req := &txs.QueryStorageRequest{
		Address: address.String(),
		Key:     key,
	}

	res, err := b.queryClient.Storage(rpctypes.ContextWithHeight(blockNum.Int64()), req)
	if err != nil {
		return nil, err
	}

	value := common.HexToHash(res.Value)
	return value.Bytes(), nil
}

func (b *BackendImpl) GetCoinbase() (sdktypes.AccAddress, error) {
	node, err := b.clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	status, err := node.Status(b.ctx)
	if err != nil {
		return nil, err
	}

	req := &txs.QueryValidatorAccountRequest{
		ConsAddress: sdktypes.ConsAddress(status.ValidatorInfo.Address).String(),
	}

	res, err := b.queryClient.ValidatorAccount(b.ctx, req)
	if err != nil {
		return nil, err
	}

	address, _ := sdktypes.AccAddressFromBech32(res.AccountAddress) // #nosec G703
	return address, nil
}

// /////////////////////////////////////////////////
// /////////////////////////////////////////////////
// /////////////////////////////////////////////////
// /////////////////////////////////////////////////
func (b *BackendImpl) SetHead(_ uint64) {
	b.logger.Error("SetHead is not implemented")
}

func (b *BackendImpl) currentBlock() (*rpctypes.Block, error) {
	block, err := b.ArtBlockByNumber(context.Background(), rpc.LatestBlockNumber)
	if err != nil {
		b.logger.Error("get CurrentBlock failed", "error", err)
		return nil, err
	}
	return block, nil
}

// GetReceipts get receipts by block hash
func (b *BackendImpl) GetReceipts(ctx context.Context, hash common.Hash) (ethtypes.Receipts, error) {
	resBlock, err := b.CosmosBlockByHash(hash)
	if err != nil || resBlock == nil {
		return nil, fmt.Errorf("query block failed, block hash %s, %w", hash.String(), err)
	}

	receipts := make([]*ethtypes.Receipt, 0, len(resBlock.Block.Txs))
	for _, tx := range resBlock.Block.Txs {
		receipt, err := b.GetTransactionReceipt(ctx, common.Hash(tx.Hash()))
		if err != nil {
			// might be not a eth tx, skip it
			continue
		}
		var contractAddress common.Address
		if receipt["contractAddress"] != nil {
			contractAddress = receipt["contractAddress"].(common.Address)
		}
		var effectiveGasPrice big.Int
		if receipt["effectiveGasPrice"] != nil {
			effectiveGasPrice = big.Int(receipt["effectiveGasPrice"].(hexutil.Big))
		}
		receipts = append(receipts, &ethtypes.Receipt{
			Type:              uint8(receipt["type"].(hexutil.Uint)),
			PostState:         []byte{},
			Status:            uint64(receipt["status"].(hexutil.Uint)),
			CumulativeGasUsed: uint64(receipt["cumulativeGasUsed"].(hexutil.Uint64)),
			Bloom:             receipt["logsBloom"].(ethtypes.Bloom),
			Logs:              receipt["logs"].([]*ethtypes.Log),
			TxHash:            receipt["transactionHash"].(common.Hash),
			ContractAddress:   contractAddress,
			GasUsed:           uint64(receipt["gasUsed"].(hexutil.Uint64)),
			EffectiveGasPrice: &effectiveGasPrice,
			BlockHash:         common.BytesToHash(resBlock.BlockID.Hash.Bytes()),
			BlockNumber:       big.NewInt(resBlock.Block.Height),
			TransactionIndex:  uint(receipt["transactionIndex"].(hexutil.Uint64)),
		})
	}
	return receipts, nil
}

func (b *BackendImpl) GetTd(_ context.Context, _ common.Hash) *big.Int {
	b.logger.Error("GetTd is not implemented")
	return nil
}

func (b *BackendImpl) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.scope.Track(b.chainFeed.Subscribe(ch))
}

func (b *BackendImpl) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	b.logger.Debug("called eth.rpc.rpctypes.SubscribeChainHeadEvent", "ch", ch)
	return b.scope.Track(b.chainHeadFeed.Subscribe(ch))
}

func (b *BackendImpl) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	b.logger.Debug("called eth.rpc.rpctypes.SubscribeChainSideEvent", "ch", ch)
	return b.scope.Track(b.chainSideFeed.Subscribe(ch))
}

// BlockNumberFromTendermint returns the BlockNumber from BlockNumberOrHash
func (b *BackendImpl) blockNumberFromCosmos(blockNrOrHash rpc.BlockNumberOrHash) (rpc.BlockNumber, error) {
	switch {
	case blockNrOrHash.BlockHash == nil && blockNrOrHash.BlockNumber == nil:
		return rpc.EarliestBlockNumber, fmt.Errorf("types BlockHash and BlockNumber cannot be both nil")
	case blockNrOrHash.BlockHash != nil:
		resBlock, err := b.CosmosBlockByHash(*blockNrOrHash.BlockHash)
		if err != nil || resBlock.Block == nil {
			return rpc.EarliestBlockNumber, err
		}
		return rpc.BlockNumber(resBlock.Block.Height), nil
	case blockNrOrHash.BlockNumber != nil:
		if *blockNrOrHash.BlockNumber == rpc.LatestBlockNumber {
			currentHeader, err := b.CurrentHeader()
			if err != nil {
				return rpc.LatestBlockNumber, err
			}
			return rpc.BlockNumber(currentHeader.Number.Int64()), nil
		}
		return *blockNrOrHash.BlockNumber, nil
	default:
		return rpc.EarliestBlockNumber, nil
	}
}

func (b *BackendImpl) BlockNumber() (hexutil.Uint64, error) {
	// do any grpc query, ignore the response and use the returned block height
	var header metadata.MD
	_, err := b.queryClient.Params(b.ctx, &txs.QueryParamsRequest{}, grpc.Header(&header))
	if err != nil {
		return hexutil.Uint64(0), err
	}

	blockHeightHeader := header.Get(grpctypes.GRPCBlockHeightHeader)
	if headerLen := len(blockHeightHeader); headerLen != 1 {
		return 0, fmt.Errorf("unexpected '%s' gRPC header length; got %d, expected: %d", grpctypes.GRPCBlockHeightHeader, headerLen, 1)
	}

	height, err := strconv.ParseUint(blockHeightHeader[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block height: %w", err)
	}

	if height > math.MaxInt64 {
		return 0, fmt.Errorf("block height %d is greater than max uint64", height)
	}

	return hexutil.Uint64(height), nil
}

func (b *BackendImpl) BlockTimeByNumber(blockNum int64) (uint64, error) {
	resBlock, err := b.clientCtx.Client.Block(b.ctx, &blockNum)
	if err != nil {
		return 0, err
	}
	return uint64(resBlock.Block.Time.Unix()), nil
}

func (b *BackendImpl) CosmosBlockResultByNumber(height *int64) (*tmrpctypes.ResultBlockResults, error) {
	return b.clientCtx.Client.BlockResults(b.ctx, height)
}

// BlockBloom query block bloom filter from block results
func (b *BackendImpl) blockBloom(blockRes *tmrpctypes.ResultBlockResults) (ethtypes.Bloom, error) {
	for _, event := range blockRes.EndBlockEvents {
		if event.Type != evmtypes.EventTypeBlockBloom {
			continue
		}

		for _, attr := range event.Attributes {
			if attr.Key == evmtypes.AttributeKeyEthereumBloom {
				encodedBloom, err := base64.StdEncoding.DecodeString(attr.Value)
				if err != nil {
					return ethtypes.Bloom{}, err
				}

				return ethtypes.BytesToBloom(encodedBloom), nil
			}
		}
	}
	return ethtypes.Bloom{}, errors.New("block bloom event is not found")
}

func (b *BackendImpl) BlockFromCosmosBlock(resBlock *tmrpctypes.ResultBlock, blockRes *tmrpctypes.ResultBlockResults) (*rpctypes.Block, error) {
	block := resBlock.Block
	height := block.Height
	bloom, err := b.blockBloom(blockRes)
	if err != nil {
		b.logger.Debug("HeaderByNumber BlockBloom failed", "height", height)
	}

	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		b.logger.Error("failed to fetch Base Fee from prunned block. Check node prunning configuration", "height", height, "error", err)
	}

	ethHeader := rpctypes.EthHeaderFromTendermint(block.Header, bloom, baseFee)
	msgs := b.EthMsgsFromCosmosBlock(resBlock, blockRes)

	txs := make([]*ethtypes.Transaction, len(msgs))
	for i, ethMsg := range msgs {
		txs[i] = ethMsg.AsTransaction()
	}

	gasUsed := uint64(0)
	for _, txsResult := range blockRes.TxsResults {
		// workaround for cosmos-sdk bug. https://github.com/cosmos/cosmos-sdk/issues/10832
		if utils.ShouldIgnoreGasUsed(txsResult) {
			// block gas limit has exceeded, other txs must have failed with same reason.
			break
		}
		gasUsed += uint64(txsResult.GetGasUsed())
	}
	ethHeader.GasUsed = gasUsed

	gasLimit, err := rpctypes.BlockMaxGasFromConsensusParams(context.Background(), b.clientCtx, block.Height)
	if err != nil {
		b.logger.Error("failed to query consensus params", "error", err.Error())
	}
	ethHeader.GasLimit = uint64(gasLimit)

	blockHash := common.BytesToHash(block.Hash().Bytes())

	ethBlock := ethtypes.NewBlock(ethHeader, txs, nil, nil, trie.NewStackTrie(nil))
	res := rpctypes.EthBlockToBlock(ethBlock)
	res.SetHash(blockHash)
	return res, nil
}

func (b *BackendImpl) EthMsgsFromCosmosBlock(resBlock *tmrpctypes.ResultBlock, blockRes *tmrpctypes.ResultBlockResults) []*txs.MsgEthereumTx {
	var result []*txs.MsgEthereumTx
	block := resBlock.Block

	txResults := blockRes.TxsResults

	for i, tx := range block.Txs {
		// Check if tx exists on EVM by cross checking with blockResults:
		//  - Include unsuccessful tx that exceeds block gas limit
		//  - Exclude unsuccessful tx with any other error but ExceedBlockGasLimit
		if !rpctypes.TxSuccessOrExceedsBlockGasLimit(txResults[i]) {
			b.logger.Debug("invalid tx result code", "cosmos-hash", hexutil.Encode(tx.Hash()))
			continue
		}

		tx, err := b.clientCtx.TxConfig.TxDecoder()(tx)
		if err != nil {
			b.logger.Debug("failed to decode transaction in block", "height", block.Height, "error", err.Error())
			continue
		}

		for _, msg := range tx.GetMsgs() {
			ethMsg, ok := msg.(*txs.MsgEthereumTx)
			if !ok {
				continue
			}

			ethMsg.Hash = ethMsg.AsTransaction().Hash().Hex()
			result = append(result, ethMsg)
		}
	}

	return result
}

func (b *BackendImpl) BlockBloom(blockRes *tmrpctypes.ResultBlockResults) (ethtypes.Bloom, error) {
	return b.blockBloom(blockRes)
}

func (b *BackendImpl) GetBlockByNumber(blockNum rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	block, err := b.ArtBlockByNumber(context.Background(), blockNum)
	if err != nil {
		return nil, err
	}

	return api.RPCMarshalHeader(block.Header(), block.Hash()), nil
}

func (b *BackendImpl) processBlock(
	tendermintBlock *tmrpctypes.ResultBlock,
	ethBlock *map[string]interface{},
	rewardPercentiles []float64,
	tendermintBlockResult *tmrpctypes.ResultBlockResults,
) (*rpctypes.OneFeeHistory, error) {
	blockHeight := tendermintBlock.Block.Height
	blockBaseFee, err := b.BaseFee(tendermintBlockResult)
	if err != nil {
		return nil, err
	}

	targetOneFeeHistory := &rpctypes.OneFeeHistory{}
	targetOneFeeHistory.BaseFee = blockBaseFee
	cfg, err := b.chainConfig()
	if err != nil {
		return nil, err
	}
	if cfg.IsLondon(big.NewInt(blockHeight + 1)) {
		header, err := b.CurrentHeader()
		if err != nil {
			return nil, err
		}
		targetOneFeeHistory.NextBaseFee = misc.CalcBaseFee(cfg, header)
	} else {
		targetOneFeeHistory.NextBaseFee = new(big.Int)
	}
	gasLimitUint64, ok := (*ethBlock)["gasLimit"].(hexutil.Uint64)
	if !ok {
		return nil, fmt.Errorf("invalid gas limit type: %T", (*ethBlock)["gasLimit"])
	}

	gasUsed, ok := (*ethBlock)["gasUsed"].(hexutil.Uint64)
	if !ok {
		return nil, fmt.Errorf("invalid gas used type: %T", (*ethBlock)["gasUsed"])
	}

	gasusedfloat, _ := new(big.Float).SetInt(new(big.Int).SetUint64(uint64(gasUsed))).Float64()

	if gasLimitUint64 <= 0 {
		return nil, fmt.Errorf("gasLimit of block height %d should be bigger than 0 , current gaslimit %d", blockHeight, gasLimitUint64)
	}

	gasUsedRatio := gasusedfloat / float64(gasLimitUint64)
	blockGasUsed := gasusedfloat
	targetOneFeeHistory.GasUsedRatio = gasUsedRatio

	rewardCount := len(rewardPercentiles)
	targetOneFeeHistory.Reward = make([]*big.Int, rewardCount)
	for i := 0; i < rewardCount; i++ {
		targetOneFeeHistory.Reward[i] = big.NewInt(0)
	}

	tendermintTxs := tendermintBlock.Block.Txs
	tendermintTxResults := tendermintBlockResult.TxsResults
	tendermintTxCount := len(tendermintTxs)

	var sorter sortGasAndReward

	for i := 0; i < tendermintTxCount; i++ {
		eachTendermintTx := tendermintTxs[i]
		eachTendermintTxResult := tendermintTxResults[i]

		tx, err := b.clientCtx.TxConfig.TxDecoder()(eachTendermintTx)
		if err != nil {
			b.logger.Debug("failed to decode transaction in block", "height", blockHeight, "error", err.Error())
			continue
		}
		txGasUsed := uint64(eachTendermintTxResult.GasUsed) // #nosec G701
		for _, msg := range tx.GetMsgs() {
			ethMsg, ok := msg.(*txs.MsgEthereumTx)
			if !ok {
				continue
			}
			tx := ethMsg.AsTransaction()
			reward := tx.EffectiveGasTipValue(blockBaseFee)
			if reward == nil {
				reward = big.NewInt(0)
			}
			sorter = append(sorter, txGasAndReward{gasUsed: txGasUsed, reward: reward})
		}
	}

	// return an all zero row if there are no transactions to gather data from
	ethTxCount := len(sorter)
	if ethTxCount == 0 {
		return targetOneFeeHistory, nil
	}

	sort.Sort(sorter)

	var txIndex int
	sumGasUsed := sorter[0].gasUsed

	for i, p := range rewardPercentiles {
		thresholdGasUsed := uint64(blockGasUsed * p / 100) // #nosec G701
		for sumGasUsed < thresholdGasUsed && txIndex < ethTxCount-1 {
			txIndex++
			sumGasUsed += sorter[txIndex].gasUsed
		}
		targetOneFeeHistory.Reward[i] = sorter[txIndex].reward
	}

	return targetOneFeeHistory, nil
}

type txGasAndReward struct {
	gasUsed uint64
	reward  *big.Int
}

type sortGasAndReward []txGasAndReward

func (s sortGasAndReward) Len() int { return len(s) }
func (s sortGasAndReward) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortGasAndReward) Less(i, j int) bool {
	return s[i].reward.Cmp(s[j].reward) < 0
}
