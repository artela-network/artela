package rpc

import (
	"context"
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/bloombits"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/gasprice"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"

	ethapi2 "github.com/artela-network/artela/ethereum/rpc/ethapi"
	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	ethereumtypes "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/txs"
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

// This is copied from filters.Backend
// eth/filters needs to be initialized from this backend type, so methods needed by
// it must also be included here.

// GetBody retrieves the block body.
func (b *backend) GetBody(ctx context.Context, hash common.Hash,
	number rpc.BlockNumber,
) (*ethtypes.Body, error) {
	return nil, nil
}

// GetLogs returns the logs.
func (b *backend) GetLogs(
	_ context.Context, blockHash common.Hash, number uint64,
) ([][]*ethtypes.Log, error) {
	return nil, nil
}

func (b *backend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.scope.Track(b.rmLogsFeed.Subscribe(ch))
}

func (b *backend) SubscribeLogsEvent(ch chan<- []*ethtypes.Log) event.Subscription {
	return b.scope.Track(b.logsFeed.Subscribe(ch))
}

func (b *backend) SubscribePendingLogsEvent(ch chan<- []*ethtypes.Log) event.Subscription {
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
// 	_ context.Context, header *ethtypes.Header,
// ) *vm.BlockContext {
// 	return nil
// }

func (b *backend) BaseFee(height int64) (*big.Int, error) {
	// return BaseFee if London hard fork is activated and feemarket is enabled
	res, err := b.queryClient.BaseFee(rpctypes.ContextWithHeight(height), &txs.QueryBaseFeeRequest{})
	if err != nil || res.BaseFee == nil {
		return nil, err
	}

	return res.BaseFee.BigInt(), nil
}

func (b *backend) PendingTransactions() ([]*sdktypes.Tx, error) {
	return nil, nil
}
