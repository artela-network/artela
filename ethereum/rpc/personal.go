package rpc

import (
	"fmt"
	"time"

	"github.com/artela-network/artela/ethereum/crypto/ethsecp256k1"
	"github.com/artela-network/artela/ethereum/crypto/hd"
	types2 "github.com/artela-network/artela/ethereum/types"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (b *BackendImpl) NewAccount(password string) (common.AddressEIP55, error) {
	name := "key_" + time.Now().UTC().Format(time.RFC3339)

	cfg := sdktypes.GetConfig()
	basePath := cfg.GetFullBIP44Path()

	hdPathIter, err := types2.NewHDPathIterator(basePath, true)
	if err != nil {
		b.logger.Info("NewHDPathIterator failed", "error", err)
		return common.AddressEIP55{}, err
	}
	// create the mnemonic and save the account
	hdPath := hdPathIter()

	info, _, err := b.clientCtx.Keyring.NewMnemonic(name, keyring.English, hdPath.String(), password, hd.EthSecp256k1)
	if err != nil {
		b.logger.Info("NewMnemonic failed", "error", err)
		return common.AddressEIP55{}, err
	}

	pubKey, err := info.GetPubKey()
	if err != nil {
		b.logger.Info("GetPubKey failed", "error", err)
		return common.AddressEIP55{}, err
	}
	addr := common.BytesToAddress(pubKey.Address().Bytes())
	return common.AddressEIP55(addr), nil
}

func (b *BackendImpl) ImportRawKey(privkey, password string) (common.Address, error) {
	priv, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return common.Address{}, err
	}

	privKey := &ethsecp256k1.PrivKey{Key: crypto.FromECDSA(priv)}

	addr := sdktypes.AccAddress(privKey.PubKey().Address().Bytes())
	ethereumAddr := common.BytesToAddress(addr)

	// return if the key has already been imported
	if _, err := b.clientCtx.Keyring.KeyByAddress(addr); err == nil {
		return ethereumAddr, nil
	}

	// ignore error as we only care about the length of the list
	list, _ := b.clientCtx.Keyring.List() // #nosec G703
	privKeyName := fmt.Sprintf("personal_%d", len(list))

	armor := sdkcrypto.EncryptArmorPrivKey(privKey, password, ethsecp256k1.KeyType)

	if err := b.clientCtx.Keyring.ImportPrivKey(privKeyName, armor, password); err != nil {
		return common.Address{}, err
	}

	return ethereumAddr, nil
}
