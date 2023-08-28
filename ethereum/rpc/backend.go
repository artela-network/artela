package rpc

import (
	"context"
	"errors"
	"math"
	"math/big"
	"strconv"
	"time"

	ethapi2 "github.com/artela-network/artela/ethereum/rpc/ethapi"
	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	ethereumtypes "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/txs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/bloombits"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/gasprice"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

// Backend represents the backend object for a artela. It extends the standard
// go-ethereum backend object.
type Backend interface {
	ethapi2.Backend
}

// backend represents the backend for the JSON-RPC service.
type backend struct {
	extRPCEnabled bool
	artela        *ArtelaService
	cfg           *Config
	chainID       *big.Int
	gpo           *gasprice.Oracle
	logger        log.Logger

	scope           event.SubscriptionScope
	chainFeed       event.Feed
	chainHeadFeed   event.Feed
	logsFeed        event.Feed
	pendingLogsFeed event.Feed
	rmLogsFeed      event.Feed
	chainSideFeed   event.Feed
	newTxsFeed      event.Feed

	ctx         context.Context
	clientCtx   client.Context
	queryClient *rpctypes.QueryClient
	// am manage etherum account, any updates of the artela account should also update to am.
	ab ethapi2.AccountBackend
}

// NewBackend create the backend instance
func NewBackend(
	clientCtx client.Context,
	artela *ArtelaService,
	extRPCEnabled bool,
	cfg *Config,
) Backend {
	b := &backend{
		ctx:           context.Background(),
		extRPCEnabled: extRPCEnabled,
		artela:        artela,
		cfg:           cfg,
		logger:        log.Root(),
		clientCtx:     clientCtx,
		queryClient:   rpctypes.NewQueryClient(clientCtx),

		scope: event.SubscriptionScope{},
	}

	var err error
	b.chainID, err = ethereumtypes.ParseChainID(clientCtx.ChainID)
	if err != nil {
		panic(err)
	}

	// Set the Backend.
	b.ab = NewAccountBackend(b.ctx, clientCtx, b.queryClient)

	if cfg.GPO.Default == nil {
		panic("cfg.GPO.Default is nil")
	}
	b.gpo = gasprice.NewOracle(b, *cfg.GPO)
	return b
}

// General Ethereum API

func (b *backend) SyncProgress() ethereum.SyncProgress {
	return ethereum.SyncProgress{
		CurrentBlock: 0,
		HighestBlock: 0,
	}
}

func (b *backend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestTipCap(ctx)
}

func (b *backend) ChainConfig() *params.ChainConfig {
	params, err := b.queryClient.Params(b.ctx, &txs.QueryParamsRequest{})
	if err != nil {
		return nil
	}

	return params.Params.ChainConfig.EthereumConfig(b.chainID)
}

func (b *backend) FeeHistory(ctx context.Context, blockCount uint64, lastBlock rpc.BlockNumber,
	rewardPercentiles []float64) (*big.Int, [][]*big.Int, []*big.Int, []float64, error) {
	return b.gpo.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
}

func (b *backend) ChainDb() ethdb.Database { //nolint:stylecheck // conforms to interface.
	return ethdb.Database(nil)
}

func (b *backend) AccountManager() ethapi2.AccountBackend {
	return b.ab
}

func (b *backend) ExtRPCEnabled() bool {
	return b.extRPCEnabled
}

func (b *backend) RPCGasCap() uint64 {
	return b.cfg.RPCGasCap
}

func (b *backend) RPCEVMTimeout() time.Duration {
	return b.cfg.RPCEVMTimeout
}

func (b *backend) RPCTxFeeCap() float64 {
	return b.cfg.RPCTxFeeCap
}

func (b *backend) UnprotectedAllowed() bool {
	return false
}

// Blockchain API

func (b *backend) SetHead(_ uint64) {
	panic("not implemented")
}

func (b *backend) HeaderByNumber(_ context.Context, number rpc.BlockNumber) (*types.Header, error) {
	return nil, nil
}

func (b *backend) HeaderByHash(_ context.Context, hash common.Hash) (*types.Header, error) {
	return nil, nil
}

func (b *backend) HeaderByNumberOrHash(ctx context.Context,
	blockNrOrHash rpc.BlockNumberOrHash,
) (*types.Header, error) {
	return nil, nil
}

func (b *backend) CurrentHeader() *types.Header {
	return b.CurrentBlock()
}

func (b *backend) CurrentBlock() *types.Header {
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

func (b *backend) BlockByNumber(_ context.Context, number rpc.BlockNumber) (*types.Block, error) {
	return nil, nil
}

func (b *backend) BlockByHash(_ context.Context, hash common.Hash) (*types.Block, error) {
	return nil, nil
}

func (b *backend) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block, error) {
	return nil, nil
}

func (b *backend) StateAndHeaderByNumber(
	ctx context.Context, number rpc.BlockNumber,
) (*state.StateDB, *types.Header, error) {
	return nil, nil, nil
}

func (b *backend) StateAndHeaderByNumberOrHash(
	ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash,
) (*state.StateDB, *types.Header, error) {
	return nil, nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *backend) PendingBlockAndReceipts() (*types.Block, types.Receipts) {
	return nil, nil
}

func (b *backend) GetReceipts(_ context.Context, hash common.Hash) (types.Receipts, error) {
	return nil, nil
}

func (b *backend) GetTd(_ context.Context, hash common.Hash) *big.Int {
	return nil
}

func (b *backend) GetEVM(ctx context.Context, msg *core.Message, state *state.StateDB,
	header *types.Header, vmConfig *vm.Config, _ *vm.BlockContext,
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

// Transaction pool API

func (b *backend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	// verify the ethereum tx
	ethereumTx := &txs.MsgEthereumTx{}
	if err := ethereumTx.FromEthereumTx(signedTx); err != nil {
		b.logger.Error("transaction converting failed", "error", err.Error())
		return err
	}

	if err := ethereumTx.ValidateBasic(); err != nil {
		b.logger.Debug("tx failed basic validation", "error", err.Error())
		return err
	}

	// Query params to use the EVM denomination
	res, err := b.queryClient.QueryClient.Params(b.ctx, &txs.QueryParamsRequest{})
	if err != nil {
		b.logger.Error("failed to query evm params", "error", err.Error())
		return err
	}

	cosmosTx, err := ethereumTx.BuildTx(b.clientCtx.TxConfig.NewTxBuilder(), res.Params.EvmDenom)
	if err != nil {
		b.logger.Error("failed to build cosmos tx", "error", err.Error())
		return err
	}

	// Encode transaction by default Tx encoder
	txBytes, err := b.clientCtx.TxConfig.TxEncoder()(cosmosTx)
	if err != nil {
		b.logger.Error("failed to encode eth tx using default encoder", "error", err.Error())
		return err
	}

	// txHash := ethereumTx.AsTransaction().Hash()

	syncCtx := b.clientCtx.WithBroadcastMode(flags.BroadcastSync)
	rsp, err := syncCtx.BroadcastTx(txBytes)
	if rsp != nil && rsp.Code != 0 {
		err = errorsmod.ABCIError(rsp.Codespace, rsp.Code, rsp.RawLog)
	}
	if err != nil {
		b.logger.Error("failed to broadcast tx", "error", err.Error())
		return err
	}

	return nil
}

func (b *backend) GetTransaction(
	_ context.Context, txHash common.Hash,
) (*types.Transaction, common.Hash, uint64, uint64, error) {
	b.logger.Debug("called eth.rpc.backend.GetTransaction", "tx_hash", txHash)
	return nil, common.Hash{}, 0, 0, nil
}

func (b *backend) GetPoolTransactions() (types.Transactions, error) {
	b.logger.Debug("called eth.rpc.backend.GetPoolTransactions")
	return nil, nil
}

func (b *backend) GetPoolTransaction(txHash common.Hash) *types.Transaction {
	return nil
}

func (b *backend) GetPoolNonce(_ context.Context, addr common.Address) (uint64, error) {
	return 0, nil
}

func (b *backend) Stats() (int, int) {
	return 0, 0
}

func (b *backend) TxPoolContent() (
	map[common.Address]types.Transactions, map[common.Address]types.Transactions,
) {
	return nil, nil
}

func (b *backend) TxPoolContentFrom(addr common.Address) (
	types.Transactions, types.Transactions,
) {
	return nil, nil
}

func (b *backend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.scope.Track(b.newTxsFeed.Subscribe(ch))
}

// Version returns the current ethereum protocol version.
func (b *backend) Version() string {
	chainID := b.ChainConfig().ChainID
	if chainID == nil {
		b.logger.Error("eth.rpc.backend.Version", "ChainID is nil")
		return "-1"
	}
	return chainID.String()
}

func (b *backend) Engine() consensus.Engine {
	return nil
}

// This is copied from filters.Backend
// eth/filters needs to be initialized from this backend type, so methods needed by
// it must also be included here.

// GetBody retrieves the block body.
func (b *backend) GetBody(ctx context.Context, hash common.Hash,
	number rpc.BlockNumber,
) (*types.Body, error) {
	return nil, nil
}

// GetLogs returns the logs.
func (b *backend) GetLogs(
	_ context.Context, blockHash common.Hash, number uint64,
) ([][]*types.Log, error) {
	return nil, nil
}

func (b *backend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.scope.Track(b.rmLogsFeed.Subscribe(ch))
}

func (b *backend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.scope.Track(b.logsFeed.Subscribe(ch))
}

func (b *backend) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.scope.Track(b.pendingLogsFeed.Subscribe(ch))
}

func (b *backend) BloomStatus() (uint64, uint64) {
	return 0, 0
}

func (b *backend) ServiceFilter(_ context.Context, _ *bloombits.MatcherSession) {
}

// artela rpc API

func (b *backend) Listening() bool {
	return true
}

func (b *backend) PeerCount() hexutil.Uint {
	return 1
}

// ClientVersion returns the current client version.
func (b *backend) ClientVersion() string {
	return ""
}

// func (b *backend) GetBlockContext(
// 	_ context.Context, header *types.Header,
// ) *vm.BlockContext {
// 	return nil
// }
