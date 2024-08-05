package v1

import (
	"encoding/binary"
	"github.com/artela-network/artela/x/aspect/store"
)

type AspectInfo struct {
	MetaVersion   uint64
	StateVersion  uint64
	AspectVersion uint64
}

func (a *AspectInfo) MarshalText() ([]byte, error) {
	bytes := make([]byte, 26)
	// first 2 bytes saves protocol info version
	protocolVersion := store.ProtocolVersion(1)
	version, _ := protocolVersion.MarshalText()
	copy(bytes, version)
	// next 8 bytes saves meta version
	binary.BigEndian.PutUint64(bytes[2:10], a.MetaVersion)
	// next 8 bytes saves state version
	binary.BigEndian.PutUint64(bytes[10:18], a.StateVersion)
	// next 8 bytes saves aspect version
	binary.BigEndian.PutUint64(bytes[18:26], a.AspectVersion)
	return bytes, nil
}

func (a *AspectInfo) UnmarshalText(text []byte) error {
	var version store.ProtocolVersion
	if err := version.UnmarshalText(text); err != nil {
		return err
	}
	if version != 1 {
		return store.ErrInvalidProtocolInfo
	}
	if len(text) < 26 {
		return store.ErrInvalidProtocolInfo
	}

	a.MetaVersion = binary.BigEndian.Uint64(text[2:10])
	a.StateVersion = binary.BigEndian.Uint64(text[10:18])
	a.AspectVersion = binary.BigEndian.Uint64(text[18:26])

	return nil
}
