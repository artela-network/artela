package datactx

import (
	"strings"

	artelatypes "github.com/artela-network/aspect-core/types"
	"google.golang.org/protobuf/proto"
)

func toMessage(tx *artelatypes.EthTransaction, keys []string) proto.Message {
	if len(keys) == 0 {
		return tx
	}

	switch strings.ToLower(keys[0]) {
	case "nonce":
		return &artelatypes.IntData{Data: int64(tx.Nonce)}
	case "blockhash":
		return &artelatypes.BytesData{Data: tx.BlockHash}
	case "blocknumber":
		return &artelatypes.IntData{Data: tx.BlockNumber}
	case "from":
		return &artelatypes.StringData{Data: tx.From}
	case "gas":
		return &artelatypes.IntData{Data: int64(tx.Gas)}
	}
	return nil

}
