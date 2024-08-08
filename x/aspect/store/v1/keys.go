package v1

import (
	"github.com/artela-network/artela/x/aspect/store"
)

// v1 keys, first 1 byte for scope, second 2 bytes for version, next 2 bytes for type
var (
	V1AspectMetaKeyPrefix       = []byte{store.AspectScope, 0x00, 0x01, 0x00, 0x01}
	V1AspectBindingKeyPrefix    = []byte{store.AspectScope, 0x00, 0x01, 0x00, 0x02}
	V1AspectCodeKeyPrefix       = []byte{store.AspectScope, 0x00, 0x01, 0x00, 0x03}
	V1AspectPropertiesKeyPrefix = []byte{store.AspectScope, 0x00, 0x01, 0x00, 0x04}
	V1AspectStateKeyPrefix      = []byte{store.AspectScope, 0x00, 0x01, 0x00, 0xff}

	V1AccountBindingKeyPrefix = []byte{store.AccountScope, 0x00, 0x01, 0x00, 0x01}
)

var (
	V1AspectBindingFilterKeyPrefix = byte(0x01)
	V1AspectBindingDataKeyPrefix   = byte(0x02)
)
