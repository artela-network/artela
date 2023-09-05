package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"

	support "github.com/artela-network/artela/x/evm/txs/support"
	evmtypes "github.com/artela-network/artela/x/evm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func TxLogsFromEvents(events []abci.Event, msgIndex int) ([]*ethtypes.Log, error) {
	for _, event := range events {
		if event.Type != evmtypes.EventTypeTxLog {
			continue
		}

		if msgIndex > 0 {
			// not the eth tx we want
			msgIndex--
			continue
		}

		return ParseTxLogsFromEvent(event)
	}
	return nil, fmt.Errorf("eth tx logs not found for message index %d", msgIndex)
}

// ParseTxLogsFromEvent parse tx logs from one event
func ParseTxLogsFromEvent(event abci.Event) ([]*ethtypes.Log, error) {
	logs := make([]*support.Log, 0, len(event.Attributes))
	for _, attr := range event.Attributes {
		if !bytes.Equal([]byte(attr.Key), []byte(evmtypes.AttributeKeyTxLog)) {
			continue
		}

		var log support.Log
		if err := json.Unmarshal([]byte(attr.Value), &log); err != nil {
			return nil, err
		}

		logs = append(logs, &log)
	}
	return support.LogsToEthereum(logs), nil
}
