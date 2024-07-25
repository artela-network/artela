package types

import (
	"fmt"
	"strconv"

	abci "github.com/cometbft/cometbft/abci/types"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/txs"
	evmtypes "github.com/artela-network/artela/x/evm/types"
)

// EventFormat is the format version of the events.
//
// To fix the issue of txs exceeds block gas limit, we changed the event format in a breaking way.
// But to avoid forcing clients to re-sync from scatch, we make json-rpc logic to be compatible with both formats.
type EventFormat int

const (
	eventFormatUnknown EventFormat = iota

	// Event Format 1 (the format used before PR #1062):
	// ```
	// ethereum_tx(amount, ethereumTxHash, [txIndex, txGasUsed], txHash, [receipient], ethereumTxFailed)
	// tx_log(txLog, txLog, ...)
	// ethereum_tx(amount, ethereumTxHash, [txIndex, txGasUsed], txHash, [receipient], ethereumTxFailed)
	// tx_log(txLog, txLog, ...)
	// ...
	// ```
	eventFormat1

	// Event Format 2 (the format used after PR #1062):
	// ```
	// ethereum_tx(ethereumTxHash, txIndex)
	// ethereum_tx(ethereumTxHash, txIndex)
	// ...
	// ethereum_tx(amount, ethereumTxHash, txIndex, txGasUsed, txHash, [receipient], ethereumTxFailed)
	// tx_log(txLog, txLog, ...)
	// ethereum_tx(amount, ethereumTxHash, txIndex, txGasUsed, txHash, [receipient], ethereumTxFailed)
	// tx_log(txLog, txLog, ...)
	// ...
	// ```
	// If the txs exceeds block gas limit, it only emits the first part.
	eventFormat2
)

// ParsedTx is the txs infos parsed from events.
type ParsedTx struct {
	MsgIndex int

	// the following fields are parsed from events
	From common.Address

	Hash common.Hash
	// -1 means uninitialized
	EthTxIndex int32
	GasUsed    uint64
	Failed     bool
}

// NewParsedTx initialize a ParsedTx
func NewParsedTx(msgIndex int) ParsedTx {
	return ParsedTx{MsgIndex: msgIndex, EthTxIndex: -1}
}

// ParsedTxs is the txs infos parsed from eth txs events.
type ParsedTxs struct {
	// one item per message
	Txs []ParsedTx
	// map txs hash to msg index
	TxHashes map[common.Hash]int
}

// ParseTxResult parse eth txs infos from cosmos-sdk events.
// It supports two event formats, the formats are described in the comments of the format constants.
func ParseTxResult(result *abci.ResponseDeliverTx, tx sdk.Tx) (*ParsedTxs, error) {
	format := eventFormatUnknown
	// the index of current ethereum_tx event in format 1 or the second part of format 2
	eventIndex := -1

	p := &ParsedTxs{
		TxHashes: make(map[common.Hash]int),
	}

	for _, event := range result.Events {
		if event.Type != evmtypes.EventTypeEthereumTx {
			continue
		}

		if format == eventFormatUnknown {
			// discover the format version by inspect the first ethereum_tx event.
			if len(event.Attributes) > 2 {
				format = eventFormat1
			} else {
				format = eventFormat2
			}
		}

		if len(event.Attributes) == 2 {
			// the first part of format 2
			if err := p.newTx(event.Attributes); err != nil {
				return nil, err
			}
		} else {
			// format 1 or second part of format 2
			eventIndex++
			if format == eventFormat1 {
				// append txs
				if err := p.newTx(event.Attributes); err != nil {
					return nil, err
				}
			} else {
				// the second part of format 2, update txs fields
				if err := p.updateTx(eventIndex, event.Attributes); err != nil {
					return nil, err
				}
			}
		}
	}

	// loop the second time to update the parsed tx from
	for _, event := range result.Events {
		if event.Type != sdk.EventTypeMessage {
			continue
		}

		// filter only the event module message
		fromEvmModule := false
		for _, attr := range event.Attributes {
			if attr.Key == sdk.AttributeKeyModule {
				if attr.Value == evmtypes.ModuleName {
					fromEvmModule = true
				}
			}
		}

		// ignore non-evm module messages
		if !fromEvmModule {
			continue
		}

		// fill the from address
		for _, attr := range event.Attributes {
			if attr.Key == sdk.AttributeKeySender {
				// update all parsed txs
				for i := 0; i < len(p.Txs); i++ {
					p.Txs[i].From = common.HexToAddress(attr.Value)
				}
				break
			}
		}
	}

	// some old versions miss some events, fill it with txs result
	gasUsed := uint64(result.GasUsed) // #nosec G701
	if len(p.Txs) == 1 {
		p.Txs[0].GasUsed = gasUsed
	}

	// this could only happen if txs exceeds block gas limit
	if result.Code != 0 && tx != nil {
		for i := 0; i < len(p.Txs); i++ {
			p.Txs[i].Failed = true

			// replace gasUsed with gasLimit because that's what's actually deducted.
			gasLimit := tx.GetMsgs()[i].(*txs.MsgEthereumTx).GetGas()
			p.Txs[i].GasUsed = gasLimit
		}
	}
	return p, nil
}

// ParseTxIndexerResult parse tm txs result to a format compatible with the custom txs indexer.
func ParseTxIndexerResult(txResult *tmrpctypes.ResultTx, tx sdk.Tx, getter func(*ParsedTxs) *ParsedTx) (*types.TxResult, error) {
	txs, err := ParseTxResult(&txResult.TxResult, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse txs events: block %d, index %d, %v", txResult.Height, txResult.Index, err)
	}

	parsedTx := getter(txs)
	if parsedTx == nil {
		return nil, fmt.Errorf("ethereum txs not found in msgs: block %d, index %d", txResult.Height, txResult.Index)
	}
	index := uint32(parsedTx.MsgIndex) // #nosec G701
	return &types.TxResult{
		Sender:            parsedTx.From.Hex(),
		Height:            txResult.Height,
		TxIndex:           txResult.Index,
		MsgIndex:          index,
		EthTxIndex:        parsedTx.EthTxIndex,
		Failed:            parsedTx.Failed,
		GasUsed:           parsedTx.GasUsed,
		CumulativeGasUsed: txs.AccumulativeGasUsed(parsedTx.MsgIndex),
	}, nil
}

// newTx parse a new txs from events, called during parsing.
func (p *ParsedTxs) newTx(attrs []abci.EventAttribute) error {
	msgIndex := len(p.Txs)
	tx := NewParsedTx(msgIndex)
	if err := fillTxAttributes(&tx, attrs); err != nil {
		return err
	}
	p.Txs = append(p.Txs, tx)
	p.TxHashes[tx.Hash] = msgIndex
	return nil
}

// updateTx updates an exiting txs from events, called during parsing.
// In event format 2, we update the txs with the attributes of the second `ethereum_tx` event,
// so we prefer the second event and override the first one.
func (p *ParsedTxs) updateTx(eventIndex int, attrs []abci.EventAttribute) error {
	tx := NewParsedTx(eventIndex)
	if err := fillTxAttributes(&tx, attrs); err != nil {
		return err
	}
	if tx.Hash != p.Txs[eventIndex].Hash {
		// if hash is different, index the new one too
		p.TxHashes[tx.Hash] = eventIndex
	}
	// override the txs because the second event is more trustworthy
	p.Txs[eventIndex] = tx
	return nil
}

// GetTxByHash find ParsedTx by txs hash, returns nil if not exists.
func (p *ParsedTxs) GetTxByHash(hash common.Hash) *ParsedTx {
	if idx, ok := p.TxHashes[hash]; ok {
		return &p.Txs[idx]
	}
	return nil
}

// GetTxByMsgIndex returns ParsedTx by msg index
func (p *ParsedTxs) GetTxByMsgIndex(i int) *ParsedTx {
	if i < 0 || i >= len(p.Txs) {
		return nil
	}
	return &p.Txs[i]
}

// GetTxByTxIndex returns ParsedTx by txs index
func (p *ParsedTxs) GetTxByTxIndex(txIndex int) *ParsedTx {
	if len(p.Txs) == 0 {
		return nil
	}
	// assuming the `EthTxIndex` increase continuously,
	// convert TxIndex to MsgIndex by subtract the begin TxIndex.
	msgIndex := txIndex - int(p.Txs[0].EthTxIndex)
	// GetTxByMsgIndex will check the bound
	return p.GetTxByMsgIndex(msgIndex)
}

// AccumulativeGasUsed calculates the accumulated gas used within the batch of txs
func (p *ParsedTxs) AccumulativeGasUsed(msgIndex int) (result uint64) {
	for i := 0; i <= msgIndex; i++ {
		result += p.Txs[i].GasUsed
	}
	return result
}

// fillTxAttribute parse attributes by name, less efficient than hardcode the index, but more stable against event
// format changes.
func fillTxAttribute(tx *ParsedTx, key []byte, value []byte) error {
	switch string(key) {
	case evmtypes.AttributeKeyEthereumTxHash:
		tx.Hash = common.HexToHash(string(value))
	case evmtypes.AttributeKeyTxIndex:
		txIndex, err := strconv.ParseUint(string(value), 10, 31) // #nosec G701
		if err != nil {
			return err
		}
		tx.EthTxIndex = int32(txIndex) // #nosec G701
	case evmtypes.AttributeKeyTxGasUsed:
		gasUsed, err := strconv.ParseUint(string(value), 10, 64)
		if err != nil {
			return err
		}
		tx.GasUsed = gasUsed
	case evmtypes.AttributeKeyEthereumTxFailed:
		tx.Failed = len(value) > 0
	}
	return nil
}

func fillTxAttributes(tx *ParsedTx, attrs []abci.EventAttribute) error {
	for _, attr := range attrs {
		// TODO: mark-artela
		if err := fillTxAttribute(tx, []byte(attr.Key), []byte(attr.Value)); err != nil {
			return err
		}
	}
	return nil
}
