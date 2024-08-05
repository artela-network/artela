package store

// scope prefixes
var (
	GlobalScope  = byte(0xff)
	AccountScope = byte(0xee)
	AspectScope  = byte(0xdd)
)

// global keys, shouldn't be changed in the future
var (
	AspectProtocolInfoKeyPrefix = []byte{GlobalScope, 0x01}
)
