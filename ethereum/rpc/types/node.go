package types

import (
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"
)

// NetworkingStack defines methods that allow a Polaris chain to build and expose JSON-RPC apis.
type NetworkingStack interface {
	// IsExtRPCEnabled returns true if the networking stack is configured to expose JSON-RPC APIs.
	ExtRPCEnabled() bool

	// RegisterHandler manually registers a new handler into the networking stack.
	RegisterHandler(string, string, http.Handler)

	// RegisterAPIs registers JSON-RPC handlers for the networking stack.
	RegisterAPIs([]rpc.API)

	// Start starts the networking stack.
	Start() error
}
