package keyring

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/artela-network/artela/ethereum/crypto/ethsecp256k1"
	"github.com/artela-network/artela/ethereum/crypto/hd"
)

// AppName defines the Ledger app used for signing. Artela uses the Ethereum app
const AppName = "Ethereum"

var (
	// CreatePubkey uses the ethsecp256k1 pubkey with Ethereum address generation and keccak hashing
	CreatePubkey = func(key []byte) types.PubKey { return &ethsecp256k1.PubKey{Key: key} }
	// SkipDERConversion represents whether the signed Ledger output should skip conversion from DER to BER.
	// This is set to true for signing performed by the Ledger Ethereum app.
	SkipDERConversion = true
)

// Option defines a function keys options for the ethereum Secp256k1 curve.
// It supports eth_secp256k1 keys for accounts.
func Option() keyring.Option {
	return func(options *keyring.Options) {
		options.SupportedAlgos = hd.SupportedAlgorithms
		options.SupportedAlgosLedger = hd.SupportedAlgorithmsLedger
		options.LedgerCreateKey = CreatePubkey
		options.LedgerAppName = AppName
		options.LedgerSigSkipDERConv = SkipDERConversion
	}
}
