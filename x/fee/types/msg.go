package types

import (
	errorsmod "cosmossdk.io/errors"
	cosmos "github.com/cosmos/cosmos-sdk/types"
)

var _ cosmos.Msg = &MsgUpdateParams{}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParams) GetSigners() []cosmos.AccAddress {
	addr := cosmos.MustAccAddressFromBech32(m.Authority)
	return []cosmos.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := cosmos.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	return m.Params.Validate()
}

// GetSignBytes implements the LegacyMsg interface.
func (m MsgUpdateParams) GetSignBytes() []byte {
	return cosmos.MustSortJSON(AminoCdc.MustMarshalJSON(&m))
}
