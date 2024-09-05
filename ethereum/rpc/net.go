package rpc

import (
	"github.com/cometbft/cometbft/rpc/client"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (b *BackendImpl) Listening() bool {
	tmClient := b.clientCtx.Client.(client.Client)
	netInfo, err := tmClient.NetInfo(b.ctx)
	if err != nil {
		return false
	}
	return netInfo.Listening
}

func (b *BackendImpl) PeerCount() hexutil.Uint {
	tmClient := b.clientCtx.Client.(client.Client)
	netInfo, err := tmClient.NetInfo(b.ctx)
	if err != nil {
		return 0
	}
	return hexutil.Uint(len(netInfo.Peers))
}

// Version returns the current ethereum protocol version.
func (b *BackendImpl) Version() string {
	v, _ := b.version()
	return v
}
