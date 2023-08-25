package support

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/artela-network/artela/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

// Validate performs a basic validation of the State fields.
// NOTE: states value can be empty
func (s State) Validate() error {
	if strings.TrimSpace(s.Key) == "" {
		return errorsmod.Wrap(types.ErrInvalidState, "states key hash cannot be blank")
	}

	return nil
}

// NewState creates a new State instance
func NewState(key, value common.Hash) State {
	return State{
		Key:   key.String(),
		Value: value.String(),
	}
}
