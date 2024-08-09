package store

import (
	"encoding/binary"
	"encoding/hex"
)

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
