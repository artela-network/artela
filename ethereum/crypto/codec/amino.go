package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	ethsecp256k12 "github.com/artela-network/artela/ethereum/crypto/ethsecp256k1"
)

var KeysCdc *codec.LegacyAmino

// RegisterCrypto registers all crypto dependency types with the provided Amino
// codec.
func RegisterCrypto(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&ethsecp256k12.PubKey{},
		ethsecp256k12.PubKeyName, nil)
	cdc.RegisterConcrete(&ethsecp256k12.PrivKey{},
		ethsecp256k12.PrivKeyName, nil)

	keyring.RegisterLegacyAminoCodec(cdc)
	cryptocodec.RegisterCrypto(cdc)

	// NOTE: update SDK's amino codec to include the ethsecp256k1 keys.
	// DO NOT REMOVE unless deprecated on the SDK.
	legacy.Cdc = cdc
	KeysCdc = cdc
}
