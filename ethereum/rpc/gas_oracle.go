package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// BlockByNumber for OracleBackend, should not use it anywhere, unless you know about the block hash diffs.
// Use ArtBlockByNumber instead.
func (b *BackendImpl) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*ethtypes.Block, error) {
	// DO NOT USE THIS METHOD, THIS JUST FOR OracleBackend
	block, err := b.ArtBlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	// BEWARE THE HASH OF THIS BLOCK IS NOT MATCH TO WHAT WAS STORED IN COSMOS DB
	return block.EthBlock(), nil
}

// BlockByNumberOrHash for OracleBackend, should not use it anywhere, unless you know about the block hash diffs.
// Use ArtBlockByNumberOrHash instead.
func (b *BackendImpl) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*ethtypes.Block, error) {
	blockNum, err := b.blockNumberFromCosmos(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	// BEWARE THE HASH OF THIS BLOCK IS NOT MATCH TO WHAT WAS STORED IN COSMOS DB
	return b.BlockByNumber(ctx, blockNum)
}

func (b *BackendImpl) PendingBlockAndReceipts() (*types.Block, types.Receipts) {
	return nil, nil
}
