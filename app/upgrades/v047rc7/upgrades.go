// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package v047rc7

import (
	"math/big"

	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v17.0.0
func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, keeper bankkeeper.Keeper, accountKeeper authkeeper.AccountKeeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		mintAart(ctx, accountKeeper, keeper, logger)

		// Leave modules are as-is to avoid running InitGenesis.
		logger.Debug("v047rc7 running module migrations ...")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}

func mintAart(ctx sdk.Context, accountKeeper authkeeper.AccountKeeper, bankKeeper bankkeeper.Keeper, logger log.Logger) {
	accountKeeper.IterateAccounts(ctx, func(account authtypes.AccountI) (stop bool) {

		queryBalanceRequest := &banktypes.QueryBalanceRequest{
			Address: account.GetAddress().String(),
			Denom:   "uart",
		}
		if response, err := bankKeeper.Balance(ctx, queryBalanceRequest); err != nil {
			logger.Error("MintAart get balance error: ", err, " for account: ", account.GetAddress().String())
			return false
		} else {

			amount := response.GetBalance().Amount

			if amount.BigInt().Cmp(big.NewInt(0)) <= 0 {
				// skip balance <= 0
				return false
			}
			logger.Info("MintAart get balance for account: ", account.GetAddress().String(), " : ", amount.BigInt().String())

			var mintCoins sdk.Coins
			mintCoins = mintCoins.Add(sdk.NewCoin("aart", amount))
			mintErr := bankKeeper.MintCoins(ctx, "evm", mintCoins)
			if mintErr != nil {
				logger.Error("MintAart mint coins error ", err, "aart: ", mintCoins.String())
				return false
			}
			moduleAcct := accountKeeper.GetModuleAddress("evm")

			recipientAddress, _ := sdk.AccAddressFromBech32(account.GetAddress().String())

			if sendErr := bankKeeper.SendCoins(ctx, moduleAcct, recipientAddress, mintCoins); sendErr != nil {
				logger.Error("MintAart send coins error ", sendErr, "aart: ", mintCoins.String(), " to: ", account.GetAddress().String())
				return false
			}
			logger.Info("MintAart send coins success ", account.GetAddress().String(), "aart: ", mintCoins.String())
			var burnCoins sdk.Coins

			burnCoins = burnCoins.Add(sdk.NewCoin("uart", amount))
			if sendErr := bankKeeper.SendCoins(ctx, recipientAddress, moduleAcct, burnCoins); sendErr != nil {
				logger.Error("MintAart send coins error ", sendErr, "aart: ", mintCoins.String(), " to: ", account.GetAddress().String())
				return false
			}
			burnErr := bankKeeper.BurnCoins(ctx, "evm", burnCoins)
			if burnErr != nil {
				logger.Error("MintAart burn coins error ", burnErr, "uart: ", burnCoins.String(), " to: ", account.GetAddress().String())
				return false
			}
		}

		return false
	})
}
