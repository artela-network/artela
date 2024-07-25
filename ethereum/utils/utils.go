package utils

import (
	"bytes"
	"math/big"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	ethereum "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/artela-network/artela/ethereum/crypto/ethsecp256k1"
	"github.com/artela-network/aspect-core/djpm"
)

const (
	// MainnetChainID defines the Artela EIP155 chain ID for mainnet
	MainnetChainID = "artela_11821"
	// TestnetChainID defines the Artela EIP155 chain ID for testnet
	TestnetChainID = "artela_11822"

	DevnetChainID = "artela_11823"

	LocalChainID = "artela_11820"
	// BaseDenom defines the Artela mainnet denomination
	BaseDenom = "uart"
)

// IsMainnet returns true if the chain-id has the Artela mainnet EIP155 chain prefix.
func IsMainnet(chainID string) bool {
	return strings.HasPrefix(chainID, MainnetChainID)
}

// IsTestnet returns true if the chain-id has the Artela testnet EIP155 chain prefix.
func IsTestnet(chainID string) bool {
	return strings.HasPrefix(chainID, TestnetChainID)
}

// IsTestnet returns true if the chain-id has the Artela testnet EIP155 chain prefix.
func IsDevnet(chainID string) bool {
	return strings.HasPrefix(chainID, DevnetChainID)
}

func IsLocal(chainID string) bool {
	return strings.HasPrefix(chainID, LocalChainID)
}

// IsSupportedKey returns true if the pubkey type is supported by the chain
// (i.e eth_secp256k1, amino multisig, ed25519).
// NOTE: Nested multisigs are not supported.
func IsSupportedKey(pubkey cryptotypes.PubKey) bool {
	switch pubkey := pubkey.(type) {
	case *ethsecp256k1.PubKey, *ed25519.PubKey:
		return true
	case multisig.PubKey:
		if len(pubkey.GetPubKeys()) == 0 {
			return false
		}

		for _, pk := range pubkey.GetPubKeys() {
			switch pk.(type) {
			case *ethsecp256k1.PubKey, *ed25519.PubKey:
				continue
			default:
				// Nested multisigs are unsupported
				return false
			}
		}

		return true
	default:
		return false
	}
}

// GetArtelaAddressFromBech32 returns the sdk.Account address of given address,
// while also changing bech32 human readable prefix (HRP) to the value set on
// the global sdk.Config (eg: `artela`).
// The function fails if the provided bech32 address is invalid.
func GetArtelaAddressFromBech32(address string) (sdk.AccAddress, error) {
	bech32Prefix := strings.SplitN(address, "1", 2)[0]
	if bech32Prefix == address {
		return nil, errorsmod.Wrapf(errortypes.ErrInvalidAddress, "invalid bech32 address: %s", address)
	}

	addressBz, err := sdk.GetFromBech32(address, bech32Prefix)
	if err != nil {
		return nil, errorsmod.Wrapf(errortypes.ErrInvalidAddress, "invalid address %s, %s", address, err.Error())
	}

	// safety check: shouldn't happen
	if err := sdk.VerifyAddressFormat(addressBz); err != nil {
		return nil, err
	}

	return sdk.AccAddress(addressBz), nil
}

func IsCustomizedVerification(tx *ethereum.Transaction) bool {
	v, r, s := tx.RawSignatureValues()
	zero := big.NewInt(0)
	noSigCallToContract := (v == nil || r == nil || s == nil || (v.Cmp(zero) == 0 && r.Cmp(zero) == 0 && s.Cmp(zero) == 0)) &&
		(tx.To() != nil && *tx.To() != common.Address{})

	// ignore transactions with signature or contract creation transactions
	if !noSigCallToContract {
		return false
	}

	// check data
	data := tx.Data()

	// the customized data layout will be [4B Header][4B Checksum][NB ABI.Encode(ValidationData, CallData)]
	if len(data) < 8 {
		return false
	}

	// check prefix
	if !bytes.Equal(data[:4], djpm.CustomVerificationPrefix) {
		return false
	}

	// compute checksum and check
	dataHash := crypto.Keccak256(data[8:])
	return bytes.Equal(data[4:8], dataHash[:4])
}
