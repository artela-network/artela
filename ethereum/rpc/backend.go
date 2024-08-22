package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	sdkmath "cosmossdk.io/math"
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

	"github.com/artela-network/artela/ethereum/rpc/filters"
	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	"github.com/artela-network/artela/ethereum/server/config"
	ethereumtypes "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/txs"
	feetypes "github.com/artela-network/artela/x/fee/types"
)

var (
	_ gasprice.OracleBackend = (*BackendImpl)(nil)
	_ filters.Backend        = (*BackendImpl)(nil)

	_ rpctypes.Backend             = (*BackendImpl)(nil)
	_ rpctypes.EthereumBackend     = (*BackendImpl)(nil)
	_ rpctypes.BlockChainBackend   = (*BackendImpl)(nil)
	_ rpctypes.TrancsactionBackend = (*BackendImpl)(nil)
	_ rpctypes.DebugBackend        = (*BackendImpl)(nil)
	_ rpctypes.PersonalBackend     = (*BackendImpl)(nil)
	_ rpctypes.TxPoolBackend       = (*BackendImpl)(nil)
	_ rpctypes.NetBackend          = (*BackendImpl)(nil)
	_ rpctypes.Web3Backend         = (*BackendImpl)(nil)
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

func (b *BackendImpl) CurrentHeader() (*ethtypes.Header, error) {
	block, err := b.ArtBlockByNumber(context.Background(), rpc.LatestBlockNumber)
	if err != nil {
		return nil, err
	}
	if block == nil || block.Header() == nil {
		return nil, errors.New("current block header not found")
	}
	return block.Header(), nil
}

func (b *BackendImpl) Accounts() []common.Address {
	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	infos, err := b.clientCtx.Keyring.List()
	if err != nil {
		b.logger.Info("keying list failed", "error", err)
		return nil
	}

	for _, info := range infos {
		pubKey, err := info.GetPubKey()
		if err != nil {
			b.logger.Info("getPubKey failed", "info", info, "error", err)
			return nil
		}
		addressBytes := pubKey.Address().Bytes()
		addresses = append(addresses, common.BytesToAddress(addressBytes))
	}

	return addresses
}

func (b *BackendImpl) GetBalance(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	blockNum, err := b.blockNumberFromCosmos(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	req := &txs.QueryBalanceRequest{
		Address: address.String(),
	}

	_, err = b.CosmosBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	res, err := b.queryClient.Balance(rpctypes.ContextWithHeight(blockNum.Int64()), req)
	if err != nil {
		return nil, err
	}

	val, ok := sdkmath.NewIntFromString(res.Balance)
	if !ok {
		return nil, errors.New("invalid balance")
	}

	if val.IsNegative() {
		return nil, errors.New("couldn't fetch balance. Node state is pruned")
	}

	return (*hexutil.Big)(val.BigInt()), nil
}

func (b *BackendImpl) ArtBlockByNumber(_ context.Context, number rpc.BlockNumber) (*rpctypes.Block, error) {
	resBlock, err := b.CosmosBlockByNumber(number)
	if err != nil || resBlock == nil {
		return nil, fmt.Errorf("query block failed, block number %d, %w", number, err)
	}

	blockRes, err := b.CosmosBlockResultByNumber(&resBlock.Block.Height)
	if err != nil {
		return nil, fmt.Errorf("block result not found for height %d", resBlock.Block.Height)
	}

	return b.BlockFromCosmosBlock(resBlock, blockRes)
}

func (b *BackendImpl) BlockByHash(_ context.Context, hash common.Hash) (*rpctypes.Block, error) {
	resBlock, err := b.CosmosBlockByHash(hash)
	if err != nil || resBlock == nil {
		return nil, fmt.Errorf("failed to get block by hash %s, %w", hash.Hex(), err)
	}

	blockRes, err := b.CosmosBlockResultByNumber(&resBlock.Block.Height)
	if err != nil {
		return nil, fmt.Errorf("block result not found for height %d", resBlock.Block.Height)
	}

	return b.BlockFromCosmosBlock(resBlock, blockRes)
}

func (b *BackendImpl) ChainConfig() *params.ChainConfig {
	cfg, err := b.chainConfig()
	if err != nil {
		return nil
	}
	return cfg
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

// General Ethereum DebugAPI

func (b *BackendImpl) SyncProgress() ethereum.SyncProgress {
	return ethereum.SyncProgress{
		CurrentBlock: 0,
		HighestBlock: 0,
	}
}

func (b *BackendImpl) chainConfig() (*params.ChainConfig, error) {
	params, err := b.queryClient.Params(b.ctx, &txs.QueryParamsRequest{})
	if err != nil {
		b.logger.Info("queryClient.Params failed", err)
		return nil, err
	}

	currentHeader, err := b.CurrentHeader()
	if err != nil {
		return nil, err
	}

	return params.Params.ChainConfig.EthereumConfig(currentHeader.Number.Int64(), b.chainID), nil
}

func (b *BackendImpl) ChainDb() ethdb.Database {
	return nil
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
				baseFee, err := strconv.ParseInt(evt.Attributes[0].Value, 10, 64)
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

func (b *BackendImpl) RPCMinGasPrice() int64 {
	evmParams, err := b.queryClient.Params(b.ctx, &txs.QueryParamsRequest{})
	if err != nil {
		return ethereumtypes.DefaultGasPrice
	}

	minGasPrice := b.appConf.GetMinGasPrices()
	amt := minGasPrice.AmountOf(evmParams.Params.EvmDenom).TruncateInt64()
	if amt == 0 {
		b.logger.Debug("amt is 0, return DefaultGasPrice")
		return ethereumtypes.DefaultGasPrice
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
