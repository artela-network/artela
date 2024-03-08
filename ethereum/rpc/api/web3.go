package api

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Web3Backend is the collection of methods required to satisfy the net
// RPC DebugAPI.
type Web3Backend interface {
	ClientVersion() string
}

// Web3API is the collection of net RPC DebugAPI methods.
type Web3API interface {
	ClientVersion() string
	Sha3(input hexutil.Bytes) hexutil.Bytes
}

// web3Api offers network related RPC methods.
type web3API struct {
	b Web3Backend
}

// NewWeb3API creates a new web3 DebugAPI instance.
func NewWeb3API(b Web3Backend) Web3Backend {
	return &web3API{b}
}

// ClientVersion returns the node name.
func (api *web3API) ClientVersion() string {
	return api.b.ClientVersion()
}

// Sha3 applies the ethereum sha3 implementation on the input.
// It assumes the input is hex encoded.
func (*web3API) Sha3(input hexutil.Bytes) hexutil.Bytes {
	return crypto.Keccak256(input)
}
