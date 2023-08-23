package process

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// UnPackInterfaces implements UnpackInterfacesMesssage.UnPackInterfaces
func (m QueryTraceTxRequest) UnPackInterfaces(unPacker codectypes.AnyUnpacker) error {
	for _, msg := range m.Predecessors {
		if err := msg.UnpackInterfaces(unPacker); err != nil {
			return err
		}
	}
	return m.Msg.UnpackInterfaces(unPacker)
}

func (m QueryTraceBlockRequest) UnPackInterfaces(unPacker codectypes.AnyUnpacker) error {
	for _, msg := range m.Txs {
		if err := msg.UnpackInterfaces(unPacker); err != nil {
			return err
		}
	}
	return nil
}
