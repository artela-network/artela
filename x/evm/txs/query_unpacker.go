package txs

import (
	codec "github.com/cosmos/cosmos-sdk/codec/types"
)

func (m QueryTraceTxRequest) UnPackInterfaces(unPacker codec.AnyUnpacker) error {
	for _, msg := range m.Predecessors {
		if err := msg.UnpackInterfaces(unPacker); err != nil {
			return err
		}
	}
	return m.Msg.UnpackInterfaces(unPacker)
}

func (m QueryTraceBlockRequest) UnPackInterfaces(unPacker codec.AnyUnpacker) error {
	for _, msg := range m.Txs {
		if err := msg.UnpackInterfaces(unPacker); err != nil {
			return err
		}
	}
	return nil
}
