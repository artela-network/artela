package api

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// NetBackend is the collection of methods required to satisfy the net
// RPC DebugAPI.
type NetBackend interface {
	NetAPI
}

// NetAPI is the collection of net RPC DebugAPI methods.
type NetAPI interface {
	PeerCount() hexutil.Uint
	Listening() bool
	Version() string
}

// netAPI offers network related RPC methods.
type netAPI struct {
	b NetBackend
}

// NewNetAPI creates a new net DebugAPI instance.
func NewNetAPI(b NetBackend) NetAPI {
	return &netAPI{b}
}

// Listening returns an indication if the node is listening for network connections.
func (api *netAPI) Listening() bool {
	return api.b.Listening()
}

// PeerCount returns the number of connected peers.
func (api *netAPI) PeerCount() hexutil.Uint {
	return api.b.PeerCount()
}

// Version returns the current ethereum protocol version.
func (api *netAPI) Version() string {
	return api.b.Version()
}
