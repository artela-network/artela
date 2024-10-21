package erc20

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/artela-network/artela/x/evm/precompile/erc20/types"
	evmtypes "github.com/artela-network/artela/x/evm/types"
)

func (c *ERC20Contract) loadTokenPairs(ctx sdk.Context) (types.TokenPairs, error) {
	store := ctx.KVStore(c.storeKey)
	bz := store.Get(evmtypes.KeyPrefixPrecompile)
	var tokenPairs types.TokenPairs
	err := c.cdc.Unmarshal(bz, &tokenPairs)
	if err != nil {
		return types.TokenPairs{}, err
	}

	c.tokenPairs = tokenPairs
	return c.tokenPairs, nil
}

func (c *ERC20Contract) storeTokenPairs(ctx sdk.Context) error {
	store := ctx.KVStore(c.storeKey)
	bz, err := c.cdc.Marshal(&c.tokenPairs)
	if err != nil {
		return err
	}

	store.Set(evmtypes.KeyPrefixPrecompile, bz)
	c.logger.Debug("setState: set token pair",
		"key", "KeyPrefixPrecompile",
		"value", fmt.Sprintf("%+v", c.tokenPairs))
	return nil
}
