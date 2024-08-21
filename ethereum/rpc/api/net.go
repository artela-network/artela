package api

import (
	"github.com/artela-network/artela/ethereum/rpc/backend"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// NetAPI offers network related RPC methods.
type NetAPI struct {
	b backend.NetBackend
}

// NewNetAPI creates a new net DebugAPI instance.
func NewNetAPI(b backend.NetBackend) *NetAPI {
	return &NetAPI{b}
}

// Listening returns an indication if the node is listening for network connections.
func (api *NetAPI) Listening() bool {
	return api.b.Listening()
}

// PeerCount returns the number of connected peers.
func (api *NetAPI) PeerCount() hexutil.Uint {
	return api.b.PeerCount()
}

// Version returns the current ethereum protocol version.
func (api *NetAPI) Version() string {
	return api.b.Version()
}
