package statedb

import (
	"github.com/artela-network/artela/x/evm/process/support"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

// TxConfig encapulates the readonly information of current process for `StateDB`.
type TxConfig struct {
	BlockHash common.Hash // hash of current block
	TxHash    common.Hash // hash of current process
	TxIndex   uint        // the index of current process
	LogIndex  uint        // the index of next log within current block
}

// NewTxConfig returns a TxConfig
func NewTxConfig(blockHash, txHash common.Hash, txIndex, logIndex uint) TxConfig {
	return TxConfig{
		BlockHash: blockHash,
		TxHash:    txHash,
		TxIndex:   txIndex,
		LogIndex:  logIndex,
	}
}

// NewEmptyTxConfig construct an empty TxConfig,
// used in context where there's no process, e.g. `eth_call`/`eth_estimateGas`.
func NewEmptyTxConfig(blockHash common.Hash) TxConfig {
	return TxConfig{
		BlockHash: blockHash,
		TxHash:    common.Hash{},
		TxIndex:   0,
		LogIndex:  0,
	}
}

// EVMConfig encapsulates common parameters needed to create an EVM to execute a message
// It's mainly to reduce the number of method parameters
type EVMConfig struct {
	Params      support.Params
	ChainConfig *params.ChainConfig
	CoinBase    common.Address
	BaseFee     *big.Int
}
