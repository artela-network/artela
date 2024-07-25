package cmd

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	ethtypes "github.com/artela-network/artela/ethereum/types"
	evmstate "github.com/artela-network/artela/x/evm/txs/support"
	evmtypes "github.com/artela-network/artela/x/evm/types"
)

// AddGenesisContractCmd returns add-genesis-contract cobra Command.
func AddGenesisContractCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-contract [contract_address_hex] [bin_runtime_code_hex/bin_runtime_code_path]",
		Short: "Add a genesis contract to genesis.json",
		Long: `Add a genesis contract to genesis.json. The provided contract must specify
the address and the bin-runtime code.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			rawAddr, err := hexutil.Decode(args[0])
			if err != nil {
				return errors.New("invalid address, please input a valid ethereum format address")
			}

			addr, err := sdk.AccAddressFromHexUnsafe(hex.EncodeToString(rawAddr))
			if err != nil {
				return errors.New("unable to parse address")
			}

			// load contract bin-runtime code
			contractBin, err := hexutil.Decode(args[1])
			if err != nil {
				contractBin, err = os.ReadFile(args[1])
				if err != nil {
					return errors.New("failed to load contract bytecode")
				}
			}

			// create concrete account type based on input parameters
			var genAccount *ethtypes.EthAccount

			balances := banktypes.Balance{Address: addr.String(), Coins: []sdk.Coin{}}
			baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)

			genAccount = &ethtypes.EthAccount{
				BaseAccount: baseAccount,
				CodeHash:    crypto.Keccak256Hash(contractBin).Hex(),
			}

			if err := genAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis contract: %w", err)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			authGenState := authtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)

			accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}

			if accs.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}

			// Add the new account to the set of genesis accounts and sanitize the
			// accounts afterward.
			accs = append(accs, genAccount)
			accs = authtypes.SanitizeGenesisAccounts(accs)

			genAccs, err := authtypes.PackAccounts(accs)
			if err != nil {
				return fmt.Errorf("failed to convert accounts into any's: %w", err)
			}
			authGenState.Accounts = genAccs

			authGenStateBz, err := clientCtx.Codec.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[authtypes.ModuleName] = authGenStateBz

			bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			bankGenState.Balances = append(bankGenState.Balances, balances)
			bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)
			bankGenState.Supply = bankGenState.Supply.Add(balances.Coins...)

			bankGenStateBz, err := clientCtx.Codec.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal bank genesis state: %w", err)
			}
			appState[banktypes.ModuleName] = bankGenStateBz

			// generate evm state
			var evmGenState evmstate.GenesisState
			if err := clientCtx.Codec.UnmarshalJSON(appState[evmtypes.ModuleName], &evmGenState); err != nil {
				return err
			}

			evmGenState.Accounts = append(evmGenState.Accounts, evmstate.GenesisAccount{
				Address: genAccount.EthAddress().Hex(),
				Code:    hex.EncodeToString(contractBin),
			})

			evmGenStateBz, err := clientCtx.Codec.MarshalJSON(&evmGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal bank genesis state: %w", err)
			}
			appState[evmtypes.ModuleName] = evmGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
