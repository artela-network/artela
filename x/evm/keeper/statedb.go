package keeper

import (
	"fmt"
	"math/big"

	artela "github.com/artela-network/artela/ethereum/types"
	"github.com/status-im/keycard-go/hexutils"

	sdkmath "cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	"github.com/artela-network/artela/x/evm/states"
	"github.com/artela-network/artela/x/evm/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethlog "github.com/ethereum/go-ethereum/log"
)

var _ states.Keeper = &Keeper{}

// ----------------------------------------------------------------------------
// 								   Getter
// ----------------------------------------------------------------------------

// GetAccount returns nil if account is not exist, returns error if it's not `EthAccountI`
func (k *Keeper) GetAccount(ctx cosmos.Context, addr common.Address) *states.StateAccount {
	acct := k.GetAccountWithoutBalance(ctx, addr)
	if acct == nil {
		return nil
	}

	acct.Balance = k.GetBalance(ctx, addr)
	return acct
}

// GetState loads contract states from database, implements `states.Keeper` interface.
func (k *Keeper) GetState(ctx cosmos.Context, addr common.Address, key common.Hash) common.Hash {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AddressStoragePrefix(addr))

	value := store.Get(key.Bytes())
	if len(value) == 0 {
		return common.Hash{}
	}

	return common.BytesToHash(value)
}

// GetCode loads contract code from database, implements `states.Keeper` interface.
func (k *Keeper) GetCode(ctx cosmos.Context, codeHash common.Hash) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCode)
	return store.Get(codeHash.Bytes())
}

// ----------------------------------------------------------------------------
// 								   Setter
// ----------------------------------------------------------------------------

// SetBalance update account's balance, compare with current balance first, then decide to mint or burn.
func (k *Keeper) SetBalance(ctx cosmos.Context, addr common.Address, amount *big.Int) error {
	cosmosAddr := cosmos.AccAddress(addr.Bytes())

	params := k.GetParams(ctx)
	coin := k.bankKeeper.GetBalance(ctx, cosmosAddr, params.EvmDenom)
	balance := coin.Amount.BigInt()
	delta := new(big.Int).Sub(amount, balance)
	switch delta.Sign() {
	case 1:
		// mint
		coins := cosmos.NewCoins(cosmos.NewCoin(params.EvmDenom, sdkmath.NewIntFromBigInt(delta)))
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
			return err
		}
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosAddr, coins); err != nil {
			return err
		}
	case -1:
		// burn
		coins := cosmos.NewCoins(cosmos.NewCoin(params.EvmDenom, sdkmath.NewIntFromBigInt(new(big.Int).Neg(delta))))
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, cosmosAddr, types.ModuleName, coins); err != nil {
			return err
		}
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
			return err
		}
	default:
		// not changed
	}
	return nil
}

// SetAccount updates nonce/balance/codeHash together.
func (k *Keeper) SetAccount(ctx cosmos.Context, addr common.Address, account states.StateAccount) error {
	// update account
	cosmosAddr := cosmos.AccAddress(addr.Bytes())
	acct := k.accountKeeper.GetAccount(ctx, cosmosAddr)
	if acct == nil {
		acct = k.accountKeeper.NewAccountWithAddress(ctx, cosmosAddr)
	}

	if err := acct.SetSequence(account.Nonce); err != nil {
		return err
	}

	codeHash := common.BytesToHash(account.CodeHash)

	if ethAcct, ok := acct.(artela.EthAccountI); ok {
		if err := ethAcct.SetCodeHash(codeHash); err != nil {
			return err
		}
	}

	k.accountKeeper.SetAccount(ctx, acct)

	if err := k.SetBalance(ctx, addr, account.Balance); err != nil {
		return err
	}

	k.Logger(ctx).Debug(
		"account updated",
		"ethereum-address", addr.Hex(),
		"nonce", account.Nonce,
		"codeHash", codeHash.Hex(),
		"balance", account.Balance,
	)
	return nil
}

// SetState update contract storage, delete if value is empty.
func (k *Keeper) SetState(ctx cosmos.Context, addr common.Address, key common.Hash, value []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AddressStoragePrefix(addr))
	action := "updated"
	if len(value) == 0 {
		ethlog.Info("SetState, delete key:", key.String())
		store.Delete(key.Bytes())
		action = "deleted"
	} else {
		ethlog.Info("SetState, key:", key.String(), "value:", hexutils.BytesToHex(value))
		store.Set(key.Bytes(), value)
	}
	k.Logger(ctx).Debug(
		fmt.Sprintf("states %s", action),
		"ethereum-address", addr.Hex(),
		"key", key.Hex(),
	)
}

// SetCode set contract code, delete if code is empty.
func (k *Keeper) SetCode(ctx cosmos.Context, codeHash, code []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCode)

	// store or delete code
	action := "updated"
	if len(code) == 0 {
		ethlog.Info("SetCode, delete key:", hexutils.BytesToHex(codeHash))
		store.Delete(codeHash)
		action = "deleted"
	} else {
		ethlog.Info("SetCode, key:", hexutils.BytesToHex(codeHash), "value:", hexutils.BytesToHex(code))
		store.Set(codeHash, code)
	}
	k.Logger(ctx).Debug(
		fmt.Sprintf("code %s", action),
		"code-hash", common.BytesToHash(codeHash).Hex(),
	)
}

// ForEachStorage iterate contract storage, callback return false to break early
func (k *Keeper) ForEachStorage(ctx cosmos.Context, addr common.Address, cb func(key, value common.Hash) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.AddressStoragePrefix(addr)

	iterator := cosmos.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := common.BytesToHash(iterator.Key())
		value := common.BytesToHash(iterator.Value())

		// check if iteration stops
		if !cb(key, value) {
			return
		}
	}
}

// DeleteAccount handles contract's suicide call:
// - clear balance
// - remove code
// - remove states
// - remove auth account
func (k *Keeper) DeleteAccount(ctx cosmos.Context, addr common.Address) error {
	cosmosAddr := cosmos.AccAddress(addr.Bytes())
	acct := k.accountKeeper.GetAccount(ctx, cosmosAddr)
	if acct == nil {
		return nil
	}

	// NOTE: only Ethereum accounts (contracts) can be selfdestructed
	_, ok := acct.(artela.EthAccountI)
	if !ok {
		return errorsmod.Wrapf(types.ErrInvalidAccount, "type %T, address %s", acct, addr)
	}

	// clear balance
	if err := k.SetBalance(ctx, addr, new(big.Int)); err != nil {
		return err
	}

	// clear storage
	k.ForEachStorage(ctx, addr, func(key, _ common.Hash) bool {
		k.SetState(ctx, addr, key, nil)
		return true
	})

	// remove auth account
	k.accountKeeper.RemoveAccount(ctx, acct)

	k.Logger(ctx).Debug(
		"account suicided",
		"ethereum-address", addr.Hex(),
		"cosmos-address", cosmosAddr.String(),
	)

	return nil
}
