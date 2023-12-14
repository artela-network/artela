package cmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	artelatypes "github.com/artela-network/artela/ethereum/types"
)

const (
	// Bech32Prefix defines the Bech32 prefix used for EthAccounts
	Bech32Prefix = "art"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic

	// DisplayDenom defines the denomination displayed to users in client applications.
	DisplayDenom = "art"
	// BaseDenom defines artelad base denonm
	BaseDenom = "uart"
)

// SetBech32Prefixes sets the global prefixes to be used when serializing addresses and public keys to Bech32 strings.
func setBech32Prefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func setBip44CoinType(config *sdk.Config) {
	config.SetCoinType(artelatypes.Bip44CoinType)
	config.SetPurpose(sdk.Purpose)                        // Shared
	config.SetFullFundraiserPath(artelatypes.BIP44HDPath) //nolint: staticcheck
}

// RegisterDenoms registers the base and display denominations to the SDK.
func registerDenoms() {
	if err := sdk.RegisterDenom(DisplayDenom, sdk.OneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(BaseDenom, sdk.NewDecWithPrec(1, artelatypes.BaseDenomUnit)); err != nil {
		panic(err)
	}
}

func initSDKConfig() {
	// set the address prefixes
	config := sdk.GetConfig()
	setBech32Prefixes(config)
	setBip44CoinType(config)
	config.Seal()

	registerDenoms()
}
