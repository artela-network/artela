package rpc

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"runtime"
	"strconv"
	"time"

	sdkmath "cosmossdk.io/math"
	bftclient "github.com/cometbft/cometbft/rpc/client"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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

	"github.com/artela-network/artela/ethereum/rpc/api"
	ethapi2 "github.com/artela-network/artela/ethereum/rpc/ethapi"
	"github.com/artela-network/artela/ethereum/rpc/filters"
	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	"github.com/artela-network/artela/ethereum/rpc/utils"
	"github.com/artela-network/artela/ethereum/server/config"
	ethereumtypes "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/txs"
	evmtypes "github.com/artela-network/artela/x/evm/types"
	feetypes "github.com/artela-network/artela/x/fee/types"
)

var (
	_ gasprice.OracleBackend = (*BackendImpl)(nil)
	_ ethapi2.Backend        = (*BackendImpl)(nil)
	_ filters.Backend        = (*BackendImpl)(nil)
	_ api.NetBackend         = (*BackendImpl)(nil)
	_ api.DebugBackend       = (*BackendImpl)(nil)
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

func (b *BackendImpl) EthBlockByNumber(blockNum rpc.BlockNumber) (*ethtypes.Block, error) {
	// TODO implement me
	panic("implement me")
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
	cfg, err := b.chainConfig()
	if err != nil {
		return nil
	}
	return cfg
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

func (b *BackendImpl) FeeHistory(blockCount uint64, lastBlock rpc.BlockNumber,
	rewardPercentiles []float64,
) (*rpctypes.FeeHistoryResult, error) {
	blockEnd := int64(lastBlock)

	if blockEnd < 0 {
		blockNumber, err := b.BlockNumber()
		if err != nil {
			return nil, err
		}
		blockEnd = int64(blockNumber)
	}

	blocks := int64(blockCount)
	maxBlockCount := int64(b.cfg.AppCfg.JSONRPC.FeeHistoryCap)
	if blocks > maxBlockCount {
		return nil, fmt.Errorf("FeeHistory user block count %d higher than %d", blocks, maxBlockCount)
	}

	if blockEnd+1 < blocks {
		blocks = blockEnd + 1
	}

	blockStart := blockEnd + 1 - blocks
	oldestBlock := (*hexutil.Big)(big.NewInt(blockStart))

	reward := make([][]*hexutil.Big, blocks)
	rewardCount := len(rewardPercentiles)
	for i := 0; i < int(blocks); i++ {
		reward[i] = make([]*hexutil.Big, rewardCount)
	}

	thisBaseFee := make([]*hexutil.Big, blocks+1)
	thisGasUsedRatio := make([]float64, blocks)
	calculateRewards := rewardCount != 0

	for blockID := blockStart; blockID <= blockEnd; blockID++ {
		index := int32(blockID - blockStart) // #nosec G701
		// tendermint block
		tendermintblock, err := b.CosmosBlockByNumber(rpc.BlockNumber(blockID))
		if tendermintblock == nil {
			return nil, err
		}

		// eth block
		ethBlock, err := b.GetBlockByNumber(rpc.BlockNumber(blockID), true)
		if ethBlock == nil {
			return nil, err
		}

		// tendermint block result
		tendermintBlockResult, err := b.CosmosBlockResultByNumber(&tendermintblock.Block.Height)
		if tendermintBlockResult == nil {
			b.logger.Debug("block result not found", "height", tendermintblock.Block.Height, "error", err.Error())
			return nil, err
		}

		oneFeeHistory, err := b.processBlock(tendermintblock, &ethBlock, rewardPercentiles, tendermintBlockResult)
		if err != nil {
			return nil, err
		}

		// copy
		thisBaseFee[index] = (*hexutil.Big)(oneFeeHistory.BaseFee)
		thisBaseFee[index+1] = (*hexutil.Big)(oneFeeHistory.NextBaseFee)
		thisGasUsedRatio[index] = oneFeeHistory.GasUsedRatio
		if calculateRewards {
			for j := 0; j < rewardCount; j++ {
				reward[index][j] = (*hexutil.Big)(oneFeeHistory.Reward[j])
				if reward[index][j] == nil {
					reward[index][j] = (*hexutil.Big)(big.NewInt(0))
				}
			}
		}
	}

	feeHistory := rpctypes.FeeHistoryResult{
		OldestBlock:  oldestBlock,
		BaseFee:      thisBaseFee,
		GasUsedRatio: thisGasUsedRatio,
	}

	if calculateRewards {
		feeHistory.Reward = reward
	}

	return &feeHistory, nil
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

func (b *BackendImpl) RPCTxFeeCap() float64 {
	return b.cfg.RPCTxFeeCap
}

func (b *BackendImpl) UnprotectedAllowed() bool {
	if b.cfg.AppCfg == nil {
		return false
	}
	return b.cfg.AppCfg.JSONRPC.AllowUnprotectedTxs
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

// artela rpc DebugAPI

// artela rpc DebugAPI

func (b *BackendImpl) Listening() bool {
	tmClient := b.clientCtx.Client.(bftclient.Client)
	netInfo, err := tmClient.NetInfo(b.ctx)
	if err != nil {
		return false
	}
	return netInfo.Listening
}

func (b *BackendImpl) PeerCount() hexutil.Uint {
	tmClient := b.clientCtx.Client.(bftclient.Client)
	netInfo, err := tmClient.NetInfo(b.ctx)
	if err != nil {
		return 0
	}
	return hexutil.Uint(len(netInfo.Peers))
}

// ClientVersion returns the current client version.
func (b *BackendImpl) ClientVersion() string {
	return fmt.Sprintf(
		"%s/%s/%s/%s",
		version.AppName,
		version.Version,
		runtime.GOOS+"-"+runtime.GOARCH,
		runtime.Version(),
	)
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

func (b *BackendImpl) PendingTransactions() ([]*sdktypes.Tx, error) {
	return nil, errors.New("PendingTransactions is not implemented")
}

func (b *BackendImpl) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	var (
		result *big.Int
		err    error
	)
	if head, err := b.CurrentHeader(); err == nil && head.BaseFee != nil {
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

func (b *BackendImpl) Syncing() (interface{}, error) {
	status, err := b.clientCtx.Client.Status(b.ctx)
	if err != nil {
		return false, err
	}

	if !status.SyncInfo.CatchingUp {
		return false, nil
	}

	return map[string]interface{}{
		"startingBlock": hexutil.Uint64(status.SyncInfo.EarliestBlockHeight),
		"currentBlock":  hexutil.Uint64(status.SyncInfo.LatestBlockHeight),
		// "highestBlock":  nil, // NA
		// "pulledStates":  nil, // NA
		// "knownStates":   nil, // NA
	}, nil
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
		return nil, errors.New("header not found")
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
