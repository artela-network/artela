package support

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethereum "github.com/ethereum/go-ethereum/core/types"

	artela "github.com/artela-network/artela/ethereum/types"
)

// ----------------------------------------------------------------------------
// 							 Transaction Logs
// ----------------------------------------------------------------------------

// NewTransactionLogs creates a new NewTransactionLogs instance.
// TransactionLogs define the logs support from a txs execution with a given hash.
func NewTransactionLogs(hash common.Hash, logs []*Log) TransactionLogs {
	return TransactionLogs{
		Hash: hash.String(),
		Logs: logs,
	}
}

// NewTransactionLogsFromEth creates a new NewTransactionLogs instance using []*ethereum.Log.
func NewTransactionLogsFromEth(hash common.Hash, ethlogs []*ethereum.Log) TransactionLogs {
	return TransactionLogs{
		Hash: hash.String(),
		Logs: NewLogsFromEth(ethlogs),
	}
}

// Validate performs a basic validation of a GenesisAccount fields.
func (tx TransactionLogs) Validate() error {
	if artela.IsEmptyHash(tx.Hash) {
		return fmt.Errorf("hash cannot be the empty %s", tx.Hash)
	}

	for i, log := range tx.Logs {
		if log == nil {
			return fmt.Errorf("log %d cannot be nil", i)
		}
		if err := log.Validate(); err != nil {
			return fmt.Errorf("invalid log %d: %w", i, err)
		}
		if log.TxHash != tx.Hash {
			return fmt.Errorf("log txs hash mismatch (%s ≠ %s)", log.TxHash, tx.Hash)
		}
	}
	return nil
}

// EthLogs returns the Ethereum type Logs from the Transaction Logs.
func (tx TransactionLogs) EthLogs() []*ethereum.Log {
	return LogsToEthereum(tx.Logs)
}

// ----------------------------------------------------------------------------
// 							     Log
// ----------------------------------------------------------------------------

// Validate performs a basic validation of an ethereum Log fields.
// Log represents a protobuf compatible Ethereum Log that defines a contract
// log event.
func (log *Log) Validate() error {
	if err := artela.ValidateAddress(log.Address); err != nil {
		return fmt.Errorf("invalid log address %w", err)
	}
	if artela.IsEmptyHash(log.BlockHash) {
		return fmt.Errorf("block hash cannot be the empty %s", log.BlockHash)
	}
	if log.BlockNumber == 0 {
		return errors.New("block number cannot be zero")
	}
	if artela.IsEmptyHash(log.TxHash) {
		return fmt.Errorf("txs hash cannot be the empty %s", log.TxHash)
	}
	return nil
}

// ToEthereum returns the Ethereum type Log from a artela proto compatible Log.
func (log *Log) ToEthereum() *ethereum.Log {
	topics := make([]common.Hash, len(log.Topics))
	for i, topic := range log.Topics {
		topics[i] = common.HexToHash(topic)
	}

	return &ethereum.Log{
		Address:     common.HexToAddress(log.Address),
		Topics:      topics,
		Data:        log.Data,
		BlockNumber: log.BlockNumber,
		TxHash:      common.HexToHash(log.TxHash),
		TxIndex:     uint(log.TxIndex),
		Index:       uint(log.Index),
		BlockHash:   common.HexToHash(log.BlockHash),
		Removed:     log.Removed,
	}
}

func NewLogsFromEth(ethlogs []*ethereum.Log) []*Log {
	var logs []*Log //nolint: prealloc
	for _, ethlog := range ethlogs {
		logs = append(logs, NewLogFromEth(ethlog))
	}

	return logs
}

// LogsToEthereum casts the Artela Logs to a slice of Ethereum Logs.
func LogsToEthereum(logs []*Log) []*ethereum.Log {
	ethLogs := make([]*ethereum.Log, 0, len(logs))
	for i := range logs {
		ethLogs = append(ethLogs, logs[i].ToEthereum())
	}
	return ethLogs
}

// NewLogFromEth creates a new Log instance from a Ethereum type Log.
func NewLogFromEth(log *ethereum.Log) *Log {
	topics := make([]string, len(log.Topics))
	for i, topic := range log.Topics {
		topics[i] = topic.String()
	}

	return &Log{
		Address:     log.Address.String(),
		Topics:      topics,
		Data:        log.Data,
		BlockNumber: log.BlockNumber,
		TxHash:      log.TxHash.String(),
		TxIndex:     uint64(log.TxIndex),
		Index:       uint64(log.Index),
		BlockHash:   log.BlockHash.String(),
		Removed:     log.Removed,
	}
}
