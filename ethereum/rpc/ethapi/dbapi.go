package ethapi

import (
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// DbGet returns the raw value of a key stored in the database.
func (api *DebugAPI) DbGet(_ string) (hexutil.Bytes, error) {
	// not implement
	return nil, errors.New("not implemented")
}

// DbAncient retrieves an ancient binary blob from the append-only immutable files.
// It is a mapping to the `AncientReaderOp.Ancient` method
func (api *DebugAPI) DbAncient(_ string, _ uint64) (hexutil.Bytes, error) {
	// not implement
	return nil, errors.New("not implemented")
}

// DbAncients returns the ancient item numbers in the ancient store.
// It is a mapping to the `AncientReaderOp.Ancients` method
func (api *DebugAPI) DbAncients() (uint64, error) {
	return 0, errors.New("not implemented")
}
