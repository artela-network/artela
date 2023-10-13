package types

import (
	"github.com/ethereum/go-ethereum/common"
)

// message type and route constants
const (
	// TypeMsgEthereumTx defines the type string of an Ethereum txs
	TypeMsgEthereumTx = "ethereum_tx"
)

const (
	// ModuleName string name of module
	ModuleName = "evm"

	// StoreKey key for ethereum storage data, account code (StateDB) or block
	// related data for Web3.
	// The EVM module should use a prefix store.
	StoreKey = ModuleName

	// TransientKey is the key to access the EVM transient store, that is reset
	// during the Commit phase.
	TransientKey = "transient_" + ModuleName

	// RouterKey uses module name for routing
	RouterKey = ModuleName
)

// prefix bytes for the EVM persistent store
const (
	prefixCode = iota + 1
	prefixStorage
	prefixParams
)

// prefix bytes for the EVM transient store
const (
	prefixTransientBloom = iota + 1
	prefixTransientTxIndex
	prefixTransientLogSize
	prefixTransientGasUsed
)

// Evm module events
const (
	EventTypeEthereumTx = TypeMsgEthereumTx
	EventTypeBlockBloom = "block_bloom"
	EventTypeTxLog      = "tx_log"

	AttributeKeyContractAddress = "contract"
	AttributeKeyRecipient       = "recipient"
	AttributeKeyTxHash          = "txHash"
	AttributeKeyEthereumTxHash  = "ethereumTxHash"
	AttributeKeyTxIndex         = "txIndex"
	AttributeKeyTxGasUsed       = "txGasUsed"
	AttributeKeyTxType          = "txType"
	AttributeKeyTxLog           = "txLog"
	// AttributeKeyEthereumTxFailed txs failed in evm execution
	AttributeKeyEthereumTxFailed = "ethereumTxFailed"
	AttributeValueCategory       = ModuleName
	AttributeKeyEthereumBloom    = "bloom"

	MetricKeyTransitionDB = "transition_db"
	MetricKeyStaticCall   = "static_call"
)

// KVStore key prefixes
var (
	KeyPrefixCode    = []byte{prefixCode}
	KeyPrefixStorage = []byte{prefixStorage}
	KeyPrefixParams  = []byte{prefixParams}
)

// Transient Store key prefixes
var (
	KeyPrefixTransientBloom   = []byte{prefixTransientBloom}
	KeyPrefixTransientTxIndex = []byte{prefixTransientTxIndex}
	KeyPrefixTransientLogSize = []byte{prefixTransientLogSize}
	KeyPrefixTransientGasUsed = []byte{prefixTransientGasUsed}
)

// AddressStoragePrefix returns a prefix to iterate over a given account storage.
func AddressStoragePrefix(address common.Address) []byte {
	return append(KeyPrefixStorage, address.Bytes()...)
}

// StateKey defines the full key under which an account states is stored.
func StateKey(address common.Address, key []byte) []byte {
	return append(AddressStoragePrefix(address), key...)
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}
