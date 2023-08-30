package rpc

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/artela-network/artela/x/evm/txs"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Blockchain API

func (b *backend) SetHead(_ uint64) {
	panic("not implemented")
}

func (b *backend) HeaderByNumber(_ context.Context, number rpc.BlockNumber) (*ethtypes.Header, error) {
	return nil, nil
}

func (b *backend) HeaderByHash(_ context.Context, hash common.Hash) (*ethtypes.Header, error) {
	return nil, nil
}

func (b *backend) HeaderByNumberOrHash(ctx context.Context,
	blockNrOrHash rpc.BlockNumberOrHash,
) (*ethtypes.Header, error) {
	return nil, nil
}

func (b *backend) CurrentHeader() *ethtypes.Header {
	return b.CurrentBlock()
}

func (b *backend) CurrentBlock() *ethtypes.Header {
	var header metadata.MD
	_, err := b.queryClient.Params(b.ctx, &txs.QueryParamsRequest{}, grpc.Header(&header))
	if err != nil {
		return nil
	}

	blockHeightHeader := header.Get(grpctypes.GRPCBlockHeightHeader)
	if headerLen := len(blockHeightHeader); headerLen != 1 {
		return nil
	}

	height, err := strconv.ParseInt(blockHeightHeader[0], 10, 64)
	if err != nil {
		return nil
	}

	if height > math.MaxInt64 {
		return nil
	}

	res, err := b.clientCtx.Client.Block(b.ctx, &height)
	if err != nil {
		return nil
	}
	return &types.Header{
		// TODO fill more header fileds
		ParentHash:      common.BytesToHash(res.Block.LastCommitHash),
		UncleHash:       [32]byte{},
		Coinbase:        common.BytesToAddress(res.Block.ProposerAddress),
		Root:            [32]byte{},
		TxHash:          [32]byte{},
		ReceiptHash:     [32]byte{},
		Bloom:           [256]byte{},
		Difficulty:      &big.Int{},
		Number:          big.NewInt(res.Block.Height),
		GasLimit:        0,
		GasUsed:         0,
		Time:            uint64(res.Block.Time.Unix()),
		Extra:           []byte{},
		MixDigest:       [32]byte{},
		Nonce:           [8]byte{},
		BaseFee:         &big.Int{},
		WithdrawalsHash: nil,
		ExcessDataGas:   &big.Int{},
	}
}

func (b *backend) BlockByNumber(_ context.Context, number rpc.BlockNumber) (*ethtypes.Block, error) {
	return nil, nil
}

func (b *backend) BlockByHash(_ context.Context, hash common.Hash) (*ethtypes.Block, error) {
	return nil, nil
}

func (b *backend) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*ethtypes.Block, error) {
	return nil, nil
}

func (b *backend) StateAndHeaderByNumber(
	ctx context.Context, number rpc.BlockNumber,
) (*state.StateDB, *ethtypes.Header, error) {
	return nil, nil, nil
}

func (b *backend) StateAndHeaderByNumberOrHash(
	ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash,
) (*state.StateDB, *ethtypes.Header, error) {
	return nil, nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *backend) PendingBlockAndReceipts() (*ethtypes.Block, types.Receipts) {
	return nil, nil
}

// GetReceipts get receipts by block hash
func (b *backend) GetReceipts(_ context.Context, hash common.Hash) (types.Receipts, error) {
	return nil, nil
}

// GetTransactionReceipt get receipt by transaction hash
func (b *backend) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	res, err := b.GetTxByEthHash(hash)
	if err != nil {
		return nil, nil
	}
	resBlock, err := b.CosmosBlockByNumber(rpc.BlockNumber(res.Height))
	if err != nil {
		return nil, nil
	}
	tx, err := b.clientCtx.TxConfig.TxDecoder()(resBlock.Block.Txs[res.TxIndex])
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx: %w", err)
	}
	ethMsg := tx.GetMsgs()[res.MsgIndex].(*txs.MsgEthereumTx)

	txData, err := txs.UnpackTxData(ethMsg.Data)
	if err != nil {
		return nil, err
	}

	cumulativeGasUsed := uint64(0)
	txRes, err := b.txResult(ctx, hash, false)
	if err != nil {
		return nil, nil
	}

	// unpack tx data and get the cumulativeGasUsed
	msgTx := &txs.MsgEthereumTxResponse{}
	if err := proto.Unmarshal(txRes.TxResult.Data, msgTx); err != nil {
		return nil, fmt.Errorf("unmarshal TxResult Data failed, %w", err)
	}
	cumulativeGasUsed = msgTx.CumulativeGasUsed

	var status hexutil.Uint
	if res.Failed {
		status = hexutil.Uint(ethtypes.ReceiptStatusFailed)
	} else {
		status = hexutil.Uint(ethtypes.ReceiptStatusSuccessful)
	}

	from, err := ethMsg.GetSender(b.chainID)
	if err != nil {
		return nil, err
	}

	// parse tx logs from events
	msgIndex := int(res.MsgIndex) // #nosec G701 -- checked for int overflow already
	logs, _ := TxLogsFromEvents(txRes.TxResult.Events, msgIndex)

	if res.EthTxIndex == -1 {
		// Fallback to find tx index by iterating all valid eth transactions
		// msgs := b.EthMsgsFromTendermintBlock(resBlock, blockRes)
		// for i := range msgs {
		// 	if msgs[i].Hash == hexTx {
		// 		res.EthTxIndex = int32(i) // #nosec G701
		// 		break
		// 	}
		// }
	}
	// return error if still unable to find the eth tx index
	if res.EthTxIndex == -1 {
		return nil, errors.New("can't find index of ethereum tx")
	}

	receipt := map[string]interface{}{
		// Consensus fields: These fields are defined by the Yellow Paper
		"status":            status,
		"cumulativeGasUsed": hexutil.Uint64(cumulativeGasUsed),
		"logsBloom":         ethtypes.BytesToBloom(ethtypes.LogsBloom(logs)),
		"logs":              logs,

		// Implementation fields: These fields are added by geth when processing a transaction.
		// They are stored in the chain database.
		"transactionHash": hash,
		"contractAddress": nil,
		"gasUsed":         hexutil.Uint64(txRes.TxResult.GasUsed),

		// Inclusion information: These fields provide information about the inclusion of the
		// transaction corresponding to this receipt.
		"blockHash":        common.BytesToHash(resBlock.Block.Header.Hash()).Hex(),
		"blockNumber":      hexutil.Uint64(res.Height),
		"transactionIndex": hexutil.Uint64(res.EthTxIndex),

		// sender and receiver (contract or EOA) addreses
		"from": from,
		"to":   txData.GetTo(),
		"type": hexutil.Uint(ethMsg.AsTransaction().Type()),
	}

	if logs == nil {
		receipt["logs"] = [][]*ethtypes.Log{}
	}

	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if txData.GetTo() == nil {
		receipt["contractAddress"] = crypto.CreateAddress(from, txData.GetNonce())
	}

	if dynamicTx, ok := txData.(*txs.DynamicFeeTx); ok {
		baseFee, err := b.BaseFee(txRes.Height)
		if err == nil {
			receipt["effectiveGasPrice"] = hexutil.Big(*dynamicTx.EffectiveGasPrice(baseFee))
		}
	}

	return receipt, nil
}

func (b *backend) GetTd(_ context.Context, hash common.Hash) *big.Int {
	return nil
}

func (b *backend) GetEVM(ctx context.Context, msg *core.Message, state *state.StateDB,
	header *ethtypes.Header, vmConfig *vm.Config, _ *vm.BlockContext,
) (*vm.EVM, func() error) {
	return nil, func() error {
		return nil
	}
}

func (b *backend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.scope.Track(b.chainFeed.Subscribe(ch))
}

func (b *backend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	b.logger.Debug("called eth.rpc.backend.SubscribeChainHeadEvent", "ch", ch)
	return b.scope.Track(b.chainHeadFeed.Subscribe(ch))
}

func (b *backend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	b.logger.Debug("called eth.rpc.backend.SubscribeChainSideEvent", "ch", ch)
	return b.scope.Track(b.chainSideFeed.Subscribe(ch))
}

func (b *backend) CosmosBlockByHash(blockHash common.Hash) (*tmrpctypes.ResultBlock, error) {
	resBlock, err := b.clientCtx.Client.BlockByHash(b.ctx, blockHash.Bytes())
	if err != nil {
		return nil, err
	}

	if resBlock.Block == nil {
		return nil, nil
	}

	return resBlock, nil
}

func (b *backend) CosmosBlockByNumber(blockNum rpc.BlockNumber) (*tmrpctypes.ResultBlock, error) {
	height := blockNum.Int64()
	if height <= 0 {
		// fetch the latest block number from the app state, more accurate than the tendermint block store state.
		n, err := b.BlockNumber()
		if err != nil {
			return nil, err
		}
		height = int64(n) //#nosec G701 -- checked for int overflow already
	}
	resBlock, err := b.clientCtx.Client.Block(b.ctx, &height)
	if err != nil {
		return nil, err
	}

	if resBlock.Block == nil {
		return nil, nil
	}

	return resBlock, nil
}

// BlockNumberFromTendermint returns the BlockNumber from BlockNumberOrHash
func (b *backend) blockNumberFromCosmos(blockNrOrHash rpc.BlockNumberOrHash) (rpc.BlockNumber, error) {
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
		return *blockNrOrHash.BlockNumber, nil
	default:
		return rpc.EarliestBlockNumber, nil
	}
}

func (b *backend) BlockNumber() (hexutil.Uint64, error) {
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

func (b *backend) BlockTimeByNumber(blockNum int64) (uint64, error) {
	resBlock, err := b.clientCtx.Client.Block(b.ctx, &blockNum)
	if err != nil {
		return 0, err
	}
	return uint64(resBlock.Block.Time.Unix()), nil
}
