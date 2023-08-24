package types

const (
	// ModuleName string name of module
	ModuleName = "feeMeter"

	// StoreKey key for base fee and block gas used.
	// The Fee Market module should use a prefix store.
	StoreKey = ModuleName

	// RouterKey uses module name for routing
	RouterKey = ModuleName

	// TransientKey is the key to access the Fee transient store, that is reset
	// during the Commit phase.
	TransientKey = "transient_" + ModuleName
)

// prefix bytes for the fee persistent store
const (
	prefixBlockGasWanted    = iota + 1
	deprecatedPrefixBaseFee // unused
)

const (
	prefixTransientBlockGasUsed = iota + 1
)

// KVStore key prefixes
var (
	KeyPrefixBlockGasWanted = []byte{prefixBlockGasWanted}
)

// Transient Store key prefixes
var (
	KeyPrefixTransientBlockGasWanted = []byte{prefixTransientBlockGasUsed}
)

// fee module events
const (
	EventTypeFee = "fee"

	AttributeKeyBaseFee = "base_fee"
)
