package txs

import (
	codec "github.com/cosmos/cosmos-sdk/codec/types"
)

var (
	_ codec.UnpackInterfacesMessage = (*QueryTraceTxRequest)(nil)
	_ codec.UnpackInterfacesMessage = (*QueryTraceBlockRequest)(nil)
)

func (m QueryTraceTxRequest) UnpackInterfaces(unPacker codec.AnyUnpacker) error {
	for _, msg := range m.Predecessors {
		if err := msg.UnpackInterfaces(unPacker); err != nil {
			return err
		}
	}
	return m.Msg.UnpackInterfaces(unPacker)
}

func (m QueryTraceBlockRequest) UnpackInterfaces(unPacker codec.AnyUnpacker) error {
	for _, msg := range m.Txs {
		if err := msg.UnpackInterfaces(unPacker); err != nil {
			return err
		}
	}
	return nil
}
