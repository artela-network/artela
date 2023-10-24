package rpc

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"time"

	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
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
	"github.com/artela-network/artela/ethereum/rpc/filters"
	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	"github.com/artela-network/artela/ethereum/server/config"
	"github.com/artela-network/artela/ethereum/types"
	ethereumtypes "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/txs"
	feetypes "github.com/artela-network/artela/x/fee/types"
)

var (
	_ gasprice.OracleBackend = (*BackendImpl)(nil)
	_ ethapi2.Backend        = (*BackendImpl)(nil)
	_ filters.Backend        = (*BackendImpl)(nil)
)

// backend represents the backend for the JSON-RPC service.
type BackendImpl struct {
	extRPCEnabled bool
	artela        *ArtelaService
	cfg           *Config
	appConf       config.Config
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
	ctx *server.Context,
	clientCtx client.Context,
	artela *ArtelaService,
	extRPCEnabled bool,
	cfg *Config,
	logger log.Logger,
) *BackendImpl {
	b := &BackendImpl{
		ctx:           context.Background(),
		extRPCEnabled: extRPCEnabled,
		artela:        artela,
		cfg:           cfg,
		logger:        logger,
		clientCtx:     clientCtx,
		queryClient:   rpctypes.NewQueryClient(clientCtx),

		scope: event.SubscriptionScope{},
	}

	var err error
	b.appConf, err = config.GetConfig(ctx.Viper)
	if err != nil {
		panic(err)
	}

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

func (b *BackendImpl) SyncProgress() ethereum.SyncProgress {
	return ethereum.SyncProgress{
		CurrentBlock: 0,
		HighestBlock: 0,
	}
}

func (b *BackendImpl) SuggestGasTipCap(baseFee *big.Int) (*big.Int, error) {
	if baseFee == nil {
		// london hardfork not enabled or feemarket not enabled
		return big.NewInt(0), nil
	}

	params, err := b.queryClient.FeeMarket.Params(b.ctx, &feetypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}
	// calculate the maximum base fee delta in current block, assuming all block gas limit is consumed
	// ```
	// GasTarget = GasLimit / ElasticityMultiplier
	// Delta = BaseFee * (GasUsed - GasTarget) / GasTarget / Denominator
	// ```
	// The delta is at maximum when `GasUsed` is equal to `GasLimit`, which is:
	// ```
	// MaxDelta = BaseFee * (GasLimit - GasLimit / ElasticityMultiplier) / (GasLimit / ElasticityMultiplier) / Denominator
	//          = BaseFee * (ElasticityMultiplier - 1) / Denominator
	// ```t
	maxDelta := baseFee.Int64() * (int64(params.Params.ElasticityMultiplier) - 1) / int64(params.Params.BaseFeeChangeDenominator) // #nosec G701
	if maxDelta < 0 {
		// impossible if the parameter validation passed.
		maxDelta = 0
	}
	return big.NewInt(maxDelta), nil
}

func (b *BackendImpl) ChainConfig() *params.ChainConfig {
	params, err := b.queryClient.Params(b.ctx, &txs.QueryParamsRequest{})
	if err != nil {
		b.logger.Info("queryClient.Params failed", err)
		return nil
	}

	return params.Params.ChainConfig.EthereumConfig(b.chainID)
}

func (b *BackendImpl) FeeHistory(ctx context.Context, blockCount uint64, lastBlock rpc.BlockNumber,
	rewardPercentiles []float64,
) (*big.Int, [][]*big.Int, []*big.Int, []float64, error) {
	return b.gpo.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
}

func (b *BackendImpl) ChainDb() ethdb.Database { //nolint:stylecheck // conforms to interface.
	return ethdb.Database(nil)
}

func (b *BackendImpl) ExtRPCEnabled() bool {
	return b.extRPCEnabled
}

func (b *BackendImpl) RPCGasCap() uint64 {
	return b.cfg.RPCGasCap
}

func (b *BackendImpl) RPCEVMTimeout() time.Duration {
	return b.cfg.RPCEVMTimeout
}

func (b *BackendImpl) RPCTxFeeCap() float64 {
	return b.cfg.RPCTxFeeCap
}

func (b *BackendImpl) UnprotectedAllowed() bool {
	return false
}

// This is copied from filters.Backend
// eth/filters needs to be initialized from this backend type, so methods needed by
// it must also be included here.

// GetBody retrieves the block body.
func (b *BackendImpl) GetBody(ctx context.Context, hash common.Hash,
	number rpc.BlockNumber,
) (*ethtypes.Body, error) {
	return nil, nil
}

// GetLogs returns the logs.
func (b *BackendImpl) GetLogs(
	_ context.Context, blockHash common.Hash, number uint64,
) ([][]*ethtypes.Log, error) {
	return nil, nil
}

func (b *BackendImpl) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.scope.Track(b.rmLogsFeed.Subscribe(ch))
}

func (b *BackendImpl) SubscribeLogsEvent(ch chan<- []*ethtypes.Log) event.Subscription {
	return b.scope.Track(b.logsFeed.Subscribe(ch))
}

func (b *BackendImpl) SubscribePendingLogsEvent(ch chan<- []*ethtypes.Log) event.Subscription {
	return b.scope.Track(b.pendingLogsFeed.Subscribe(ch))
}

func (b *BackendImpl) BloomStatus() (uint64, uint64) {
	return 0, 0
}

func (b *BackendImpl) ServiceFilter(_ context.Context, _ *bloombits.MatcherSession) {
}

// artela rpc API

func (b *BackendImpl) Listening() bool {
	return true
}

func (b *BackendImpl) PeerCount() hexutil.Uint {
	return 1
}

// ClientVersion returns the current client version.
func (b *BackendImpl) ClientVersion() string {
	return ""
}

// func (b *BackendImpl) GetBlockContext(
// 	_ context.Context, header *ethtypes.Header,
// ) *vm.BlockContext {
// 	return nil
// }

func (b *BackendImpl) BaseFee(blockRes *tmrpctypes.ResultBlockResults) (*big.Int, error) {
	// return BaseFee if London hard fork is activated and feemarket is enabled
	res, err := b.queryClient.BaseFee(rpctypes.ContextWithHeight(blockRes.Height), &txs.QueryBaseFeeRequest{})
	if err != nil || res.BaseFee == nil {
		// we can't tell if it's london HF not enabled or the state is pruned,
		// in either case, we'll fallback to parsing from begin blocker event,
		// faster to iterate reversely
		for i := len(blockRes.BeginBlockEvents) - 1; i >= 0; i-- {
			evt := blockRes.BeginBlockEvents[i]
			if evt.Type == feetypes.EventTypeFee && len(evt.Attributes) > 0 {
				baseFee, err := strconv.ParseInt(string(evt.Attributes[0].Value), 10, 64)
				if err == nil {
					return big.NewInt(baseFee), nil
				}
				break
			}
		}
		return nil, err
	}

	if res.BaseFee == nil {
		b.logger.Debug("res.BaseFee is nil")
		return nil, nil
	}

	return res.BaseFee.BigInt(), nil
}

func (b *BackendImpl) PendingTransactions() ([]*sdktypes.Tx, error) {
	return nil, errors.New("PendingTransactions is not implemented")
}

func (b *BackendImpl) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	var (
		result *big.Int
		err    error
	)
	if head := b.CurrentHeader(); head.BaseFee != nil {
		result, err = b.SuggestGasTipCap(head.BaseFee)
		if err != nil {
			return nil, err
		}
		result = result.Add(result, head.BaseFee)
	} else {
		result = big.NewInt(b.RPCMinGasPrice())
	}

	// return at least GlobalMinGasPrice from FeeMarket module
	minGasPrice, err := b.GlobalMinGasPrice()
	if err != nil {
		return nil, err
	}
	minGasPriceInt := minGasPrice.TruncateInt().BigInt()
	if result.Cmp(minGasPriceInt) < 0 {
		result = minGasPriceInt
	}

	return (*hexutil.Big)(result), nil
}

func (b *BackendImpl) RPCMinGasPrice() int64 {
	evmParams, err := b.queryClient.Params(b.ctx, &txs.QueryParamsRequest{})
	if err != nil {
		return types.DefaultGasPrice
	}

	minGasPrice := b.appConf.GetMinGasPrices()
	amt := minGasPrice.AmountOf(evmParams.Params.EvmDenom).TruncateInt64()
	if amt == 0 {
		b.logger.Debug("amt is 0, return DefaultGasPrice")
		return types.DefaultGasPrice
	}

	return amt
}

// GlobalMinGasPrice returns MinGasPrice param from FeeMarket
func (b *BackendImpl) GlobalMinGasPrice() (sdktypes.Dec, error) {
	res, err := b.queryClient.FeeMarket.Params(b.ctx, &feetypes.QueryParamsRequest{})
	if err != nil {
		return sdktypes.ZeroDec(), err
	}
	return res.Params.MinGasPrice, nil
}

// RPCBlockRangeCap defines the max block range allowed for `eth_getLogs` query.
func (b *BackendImpl) RPCBlockRangeCap() int32 {
	return b.appConf.JSONRPC.BlockRangeCap
}

// RPCFilterCap is the limit for total number of filters that can be created
func (b *BackendImpl) RPCFilterCap() int32 {
	return b.appConf.JSONRPC.FilterCap
}

// RPCLogsCap defines the max number of results can be returned from single `eth_getLogs` query.
func (b *BackendImpl) RPCLogsCap() int32 {
	return b.appConf.JSONRPC.LogsCap
}
