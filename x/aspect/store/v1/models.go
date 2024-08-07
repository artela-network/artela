package v1

import (
	"encoding/binary"
	"github.com/artela-network/artela/x/aspect/store"
	"github.com/ethereum/go-ethereum/common"
)

var emptyHash = common.Hash{}

type DataLength uint64

func (l *DataLength) UnmarshalText(text []byte) error {
	if len(text) < 8 {
		*l = 0
	} else {
		*l = DataLength(binary.BigEndian.Uint64(text[:8]))
	}
	return nil
}

func (l DataLength) MarshalText() (text []byte, err error) {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, uint64(l))
	return result, nil
}

type Bindings []Binding

func (b *Bindings) UnmarshalText(text []byte) error {
	if len(text)%32 != 0 {
		return store.ErrInvalidBinding
	}
	for i := 0; i < len(text); i += 32 {
		data := text[i : i+32]
		if common.Hash(data) == emptyHash {
			// EOF
			return nil
		}

		var binding Binding
		if err := binding.UnmarshalText(data); err != nil {
			return err
		}
		*b = append(*b, binding)
	}
	return nil
}

func (b Bindings) MarshalText() (text []byte, err error) {
	result := make([]byte, 0, len(b)*32)
	for _, binding := range b {
		marshaled, err := binding.MarshalText()
		if err != nil {
			return nil, err
		}
		result = append(result, marshaled...)
	}
	return result, nil
}

type Binding struct {
	Account   common.Address
	Version   uint64
	Priority  int8
	JoinPoint uint16
}

func (b *Binding) UnmarshalText(text []byte) error {
	if len(text) < 32 {
		return store.ErrInvalidBinding
	}
	b.Account.SetBytes(text[:20])
	b.Version = binary.BigEndian.Uint64(text[20:28])
	b.Priority = int8(text[28])
	b.JoinPoint = binary.BigEndian.Uint16(text[29:31])
	// last byte is reserved for future use
	return nil
}

func (b Binding) MarshalText() (text []byte, err error) {
	result := make([]byte, 32)
	copy(result[:20], b.Account.Bytes())
	binary.BigEndian.PutUint64(result[20:28], b.Version)
	result[28] = byte(b.Priority)
	binary.BigEndian.PutUint16(result[29:31], b.JoinPoint)
	// last byte is reserved for future use
	return result, nil
}

type Extension struct {
	AspectVersion uint64
	PayMaster     common.Address
	Proof         []byte
}

func (e *Extension) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		// not inited
		return nil
	}
	if len(text) < 36 {
		return store.ErrInvalidExtension
	}
	e.AspectVersion = binary.BigEndian.Uint64(text[:8])
	e.PayMaster.SetBytes(text[8:28])
	proofLen := binary.BigEndian.Uint64(text[28:36])
	if uint64(len(text)) != 36+proofLen {
		return store.ErrInvalidExtension
	}
	e.Proof = make([]byte, proofLen)
	copy(e.Proof, text[36:])
	return nil
}

func (e Extension) MarshalText() (text []byte, err error) {
	result := make([]byte, 8+20+8+len(e.Proof))
	binary.BigEndian.PutUint64(result[:8], e.AspectVersion)
	copy(result[8:28], e.PayMaster.Bytes())
	binary.BigEndian.PutUint64(result[28:36], uint64(len(e.Proof)))
	copy(result[36:], e.Proof)
	return result, nil
}

// VersionMeta is the data model for holding the metadata of a specific version of aspect
type VersionMeta struct {
	JoinPoint uint64
	CodeHash  common.Hash
}

func (v *VersionMeta) UnmarshalText(text []byte) error {
	if len(text) < 40 {
		return store.ErrInvalidVersionMeta
	}
	v.JoinPoint = binary.BigEndian.Uint64(text[:8])
	v.CodeHash.SetBytes(text[8:40])
	return nil
}

func (v VersionMeta) MarshalText() (text []byte, err error) {
	result := make([]byte, 40)
	binary.BigEndian.PutUint64(result[:8], v.JoinPoint)
	copy(result[8:], v.CodeHash.Bytes())
	return result, nil
}
