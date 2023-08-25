package types

// DefaultGenesisState sets default fee market genesis states.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:   DefaultParams(),
		BlockGas: 0,
	}
}

// NewGenesisState creates a new genesis states.
func NewGenesisState(params Params, blockGas uint64) *GenesisState {
	return &GenesisState{
		Params:   params,
		BlockGas: blockGas,
	}
}

// Validate performs basic genesis states validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
