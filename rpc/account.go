package rpc

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/common"
)

type Account struct {
	clientCtx client.Context
}

func (acct *Account) Accounts() []common.Address {
	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	infos, err := acct.clientCtx.Keyring.List()
	if err != nil {
		return nil
	}

	for _, info := range infos {
		pubKey, err := info.GetPubKey()
		if err != nil {
			return nil
		}
		addressBytes := pubKey.Address().Bytes()
		addresses = append(addresses, common.BytesToAddress(addressBytes))
	}

	return addresses
}
