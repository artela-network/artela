package rpc

import (
	"context"
	"fmt"
	"math/big"

	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	feetypes "github.com/artela-network/artela/x/fee/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/rpc"
)

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

func (b *BackendImpl) Engine() consensus.Engine {
	b.logger.Error("Engine is not implemented")
	return nil
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
