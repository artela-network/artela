package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/ethereum/go-ethereum/common/hexutil"

	rpc "github.com/artela-network/artela/ethereum/rpc/types"
	"github.com/artela-network/artela/x/evm/txs"
	"github.com/artela-network/artela/x/evm/types"
)

// GetTxCmd returns the txs commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(NewRawTxCmd())
	return cmd
}

// NewRawTxCmd command build cosmos txs from raw ethereum txs
func NewRawTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raw TX_HEX",
		Short: "Build cosmos txs from raw ethereum txs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := hexutil.Decode(args[0])
			if err != nil {
				return errors.Wrap(err, "failed to decode ethereum txs hex bytes")
			}

			msg := &txs.MsgEthereumTx{}
			if err := msg.UnmarshalBinary(data); err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			rsp, err := rpc.NewQueryClient(clientCtx).Params(cmd.Context(), &txs.QueryParamsRequest{})
			if err != nil {
				return err
			}

			tx, err := msg.BuildTx(clientCtx.TxConfig.NewTxBuilder(), rsp.Params.EvmDenom)
			if err != nil {
				return err
			}

			if clientCtx.GenerateOnly {
				json, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
				if err != nil {
					return err
				}

				return clientCtx.PrintString(fmt.Sprintf("%s\n", json))
			}

			if !clientCtx.SkipConfirm {
				out, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
				if err != nil {
					return err
				}

				_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", out)

				buf := bufio.NewReader(os.Stdin)
				ok, err := input.GetConfirmation("confirm txs before signing and broadcasting", buf, os.Stderr)

				if err != nil || !ok {
					_, _ = fmt.Fprintf(os.Stderr, "%s\n", "canceled txs")
					return err
				}
			}

			txBytes, err := clientCtx.TxConfig.TxEncoder()(tx)
			if err != nil {
				return err
			}

			// broadcast to a Tendermint node
			res, err := clientCtx.BroadcastTx(txBytes)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
