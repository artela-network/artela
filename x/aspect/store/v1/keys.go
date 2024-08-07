package v1

import (
	"encoding/binary"
	"encoding/hex"
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

type KeyBuilder struct {
	key []byte
}

func NewKeyBuilder(prefix []byte) *KeyBuilder {
	buffer := make([]byte, len(prefix))
	copy(buffer, prefix)
	return &KeyBuilder{key: buffer}
}

func (k *KeyBuilder) AppendBytes(key []byte) *KeyBuilder {
	return NewKeyBuilder(append(k.key, key...))
}

func (k *KeyBuilder) AppendUint64(key uint64) *KeyBuilder {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, key)
	return NewKeyBuilder(append(k.key, buffer...))
}

func (k *KeyBuilder) AppendString(key string) *KeyBuilder {
	return NewKeyBuilder(append(k.key, key...))
}

func (k *KeyBuilder) AppendByte(key byte) *KeyBuilder {
	return NewKeyBuilder(append(k.key, key))
}

func (k *KeyBuilder) AppendUint8(key uint8) *KeyBuilder {
	return NewKeyBuilder(append(k.key, key))
}

func (k *KeyBuilder) Build() []byte {
	return k.key
}

func (k *KeyBuilder) String() string {
	return hex.EncodeToString(k.key)
}
