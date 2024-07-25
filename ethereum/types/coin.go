package types

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// AttoArtela defines the default coin denomination used in Artela in:
	//
	// - Staking parameters: denomination used as stake in the dPoS chain
	// - Mint parameters: denomination minted due to fee distribution rewards
	// - Governance parameters: denomination used for spam prevention in proposal deposits
	// - Crisis parameters: constant fee denomination used for spam prevention to check broken invariant
	// - EVM parameters: denomination used for running EVM states transitions in Artela.
	AttoArtela string = "uart"

	// BaseDenomUnit defines the base denomination unit for Artela.
	// 1 art = 1x10^{BaseDenomUnit} uart
	BaseDenomUnit = 18

	// DefaultGasPrice is default gas price for evm transactions
	DefaultGasPrice = 20
)

// PowerReduction defines the default power reduction value for staking
var PowerReduction = sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(BaseDenomUnit), nil))

// NewArtelaCoin is a utility function that returns an "uart" coin with the given sdkmath.Int amount.
// The function will panic if the provided amount is negative.
func NewArtelaCoin(amount sdkmath.Int) sdk.Coin {
	return sdk.NewCoin(AttoArtela, amount)
}

// NewArtelaDecCoin is a utility function that returns an "uart" decimal coin with the given sdkmath.Int amount.
// The function will panic if the provided amount is negative.
func NewArtelaDecCoin(amount sdkmath.Int) sdk.DecCoin {
	return sdk.NewDecCoin(AttoArtela, amount)
}

// NewArtelaCoinInt64 is a utility function that returns an "uart" coin with the given int64 amount.
// The function will panic if the provided amount is negative.
func NewArtelaCoinInt64(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(AttoArtela, amount)
}
