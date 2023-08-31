package rpc

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"

	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	"github.com/artela-network/artela/x/evm/txs"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
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
		if *blockNrOrHash.BlockNumber == rpc.LatestBlockNumber {
			currentHeight := b.CurrentHeader().Number
			return rpc.BlockNumber(currentHeight.Int64()), nil
		}
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

func (b *backend) CosmosBlockResultByNumber(height *int64) (*tmrpctypes.ResultBlockResults, error) {
	return b.clientCtx.Client.BlockResults(b.ctx, height)
}

func (b *backend) GetCode(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
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
