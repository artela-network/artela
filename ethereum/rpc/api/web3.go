package api

import (
	"github.com/artela-network/artela/ethereum/rpc/backend"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// web3Api offers network related RPC methods.
type Web3API struct {
	b backend.Web3Backend
}

// NewWeb3API creates a new web3 DebugAPI instance.
func NewWeb3API(b backend.Web3Backend) backend.Web3Backend {
	return &Web3API{b}
}

// ClientVersion returns the node name.
func (api *Web3API) ClientVersion() string {
	return api.b.ClientVersion()
}

// Sha3 applies the ethereum sha3 implementation on the input.
// It assumes the input is hex encoded.
func (*Web3API) Sha3(input hexutil.Bytes) hexutil.Bytes {
	return crypto.Keccak256(input)
}
