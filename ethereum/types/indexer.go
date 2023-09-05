package types

import (
	abci "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/ethereum/go-ethereum/common"
)

// EVMTxIndexer defines the interface of custom eth txs indexer.
type EVMTxIndexer interface {
	// LastIndexedBlock returns -1 if indexer db is empty
	LastIndexedBlock() (int64, error)
	IndexBlock(*tmtypes.Block, []*abci.ResponseDeliverTx) error

	// GetByTxHash returns nil if txs not found.
	GetByTxHash(common.Hash) (*TxResult, error)
	// GetByBlockAndIndex returns nil if txs not found.
	GetByBlockAndIndex(int64, int32) (*TxResult, error)
}
