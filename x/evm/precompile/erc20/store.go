package erc20

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	evmtypes "github.com/artela-network/artela/x/evm/types"
)

func (c *ERC20Contract) registerNewTokenPairs(ctx sdk.Context, denom string, proxy common.Address) error {
	c.setDenomByProxy(ctx, denom, proxy)
	return c.appendProxyByDenom(ctx, denom, proxy)
}

func (c *ERC20Contract) getProxyByDenom(ctx sdk.Context, denom string) ([]common.Address, error) {
	store := ctx.KVStore(c.storeKey)
	store = prefix.NewStore(store, evmtypes.KeyPrefixERC20Denom)
	data := store.Get([]byte(denom))
	if len(data)%common.AddressLength != 0 {
		return nil, fmt.Errorf("failed to load proxy address, value: %x", data)
	}

	addrs := make([]common.Address, len(data)/common.AddressLength)
	for i := 0; i < len(data); i += common.AddressLength {
		addrs[i] = common.BytesToAddress(data[i : i+common.AddressLength])
	}

	return addrs, nil
}

func (c *ERC20Contract) appendProxyByDenom(ctx sdk.Context, denom string, proxy common.Address) error {
	store := ctx.KVStore(c.storeKey)
	store = prefix.NewStore(store, evmtypes.KeyPrefixERC20Denom)
	data := store.Get([]byte(denom))
	if len(data)%common.AddressLength != 0 {
		return fmt.Errorf("failed to load proxy address, value: %x", data)
	}

	newData := make([]byte, len(data)+common.AddressLength)
	copy(newData, data)
	copy(newData[len(data):], proxy.Bytes())
	store.Set([]byte(denom), newData)

	c.logger.Debug("setState: set ProxyByDenom",
		"denom", denom,
		"addrs", fmt.Sprintf("%x", newData))
	return nil
}

func (c *ERC20Contract) getDenomByProxy(ctx sdk.Context, proxy common.Address) string {
	store := ctx.KVStore(c.storeKey)
	store = prefix.NewStore(store, evmtypes.KeyPrefixERC20Address)
	data := store.Get(proxy.Bytes())
	return string(data)
}

func (c *ERC20Contract) setDenomByProxy(ctx sdk.Context, denom string, proxy common.Address) {
	store := ctx.KVStore(c.storeKey)
	store = prefix.NewStore(store, evmtypes.KeyPrefixERC20Address)
	store.Set(proxy.Bytes(), []byte(denom))

	c.logger.Debug("setState: set DenomByProxy",
		"addr", proxy.String(),
		"denom", denom)
}
