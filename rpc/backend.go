package rpc

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
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

	"github.com/artela-network/artela/rpc/ethapi"
)

// Backend represents the backend object for a artela. It extends the standard
// go-ethereum backend object.
type Backend interface {
	ethapi.Backend
}

// backend represents the backend for the JSON-RPC service.
type backend struct {
	extRPCEnabled bool
	artela        *ArtelaService
	cfg           *Config
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

	// am manage etherum account, any updates of the artela account should also update to am.
	am *accounts.Manager
}

// NewBackend create the backend instance
func NewBackend(
	artela *ArtelaService,
	extRPCEnabled bool,
	cfg *Config,
	am *accounts.Manager,
) Backend {
	b := &backend{
		extRPCEnabled: extRPCEnabled,
		artela:        artela,
		cfg:           cfg,
		logger:        log.Root(),
		am:            am,

		scope: event.SubscriptionScope{},
	}

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
	return &params.ChainConfig{ChainID: big.NewInt(9000)}
}

func (b *backend) FeeHistory(ctx context.Context, blockCount uint64, lastBlock rpc.BlockNumber,
	rewardPercentiles []float64) (*big.Int, [][]*big.Int, []*big.Int, []float64, error) {
	return b.gpo.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
}

func (b *backend) ChainDb() ethdb.Database { //nolint:stylecheck // conforms to interface.
	return ethdb.Database(nil)
}

func (b *backend) AccountManager() *accounts.Manager {
	return b.am
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
	return nil
}

func (b *backend) CurrentBlock() *types.Header {
	return nil
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
	// eth tx -> cosmos tx
	// broadcast tx
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
