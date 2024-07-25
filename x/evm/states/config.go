package states

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"

	"github.com/artela-network/artela/x/evm/txs/support"
)

// EVMConfig encapsulates common parameters needed to create an EVM to execute a message
// It's mainly to reduce the number of method parameters
type EVMConfig struct {
	Params      support.Params
	ChainConfig *params.ChainConfig
	CoinBase    common.Address
	BaseFee     *big.Int
}

// TxConfig encapulates the readonly information of current txs for `StateDB`.
type TxConfig struct {
	BlockHash common.Hash // hash of current block
	TxHash    common.Hash // hash of current txs
	TxIndex   uint        // the index of current txs
	LogIndex  uint        // the index of next log within current block
	TxType    uint        // the index of next log within current block
}

// NewTxConfig returns a TxConfig
func NewTxConfig(blockHash, txHash common.Hash, txIndex, logIndex uint, txType uint) TxConfig {
	return TxConfig{
		BlockHash: blockHash,
		TxHash:    txHash,
		TxIndex:   txIndex,
		LogIndex:  logIndex,
		TxType:    txType,
	}
}

// NewEmptyTxConfig construct an empty TxConfig,
// used in context where there's no txs, e.g. `eth_call`/`eth_estimateGas`.
func NewEmptyTxConfig(blockHash common.Hash) TxConfig {
	return TxConfig{
		BlockHash: blockHash,
		TxHash:    common.Hash{},
		TxIndex:   0,
		LogIndex:  0,
	}
}
