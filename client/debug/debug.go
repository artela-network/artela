package debug

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cometbft/cometbft/libs/bytes"
	cometbftTypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/ethereum/go-ethereum/common"

	appparams "github.com/artela-network/artela/app/params"
	"github.com/artela-network/artela/ethereum/eip712"
	artela "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/txs"
)

// Cmd creates a main CLI command
func Cmd(encodingConfig appparams.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Tool for helping with debugging your application",
		RunE:  client.ValidateCmd,
	}

	cmd.AddCommand(PubkeyCmd())
	cmd.AddCommand(AddrCmd())
	cmd.AddCommand(RawBytesCmd())
	cmd.AddCommand(LegacyEIP712Cmd())
	cmd.AddCommand(CosmosTxHash(encodingConfig))

	return cmd
}

// getPubKeyFromString decodes SDK PubKey using JSON marshaler.
func getPubKeyFromString(ctx client.Context, pkstr string) (cryptotypes.PubKey, error) {
	var pk cryptotypes.PubKey
	err := ctx.Codec.UnmarshalInterfaceJSON([]byte(pkstr), &pk)
	return pk, err
}

func PubkeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pubkey [pubkey]",
		Short: "Decode a pubkey from proto JSON",
		Long:  "Decode a pubkey from proto JSON and display it's address",
		Example: fmt.Sprintf(
			`"$ %s debug pubkey '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AurroA7jvfPd1AadmmOvWM2rJSwipXfRf8yD6pLbA2DJ"}'`,
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			pk, err := getPubKeyFromString(clientCtx, args[0])
			if err != nil {
				return err
			}

			addr := pk.Address()
			cmd.Printf("Address (EIP-55): %s\n", common.BytesToAddress(addr))
			cmd.Printf("Bech32 Acc: %s\n", sdk.AccAddress(addr))
			cmd.Println("PubKey Hex:", hex.EncodeToString(pk.Bytes()))
			return nil
		},
	}
}

func AddrCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "addr [address]",
		Short: "Convert an address between hex and bech32",
		Long:  "Convert an address between hex encoding and bech32.",
		Example: fmt.Sprintf(
			`$ %s debug addr ethm10jmp6sgh4cc6zt3e8gw05wavvejgr5pw2unfju
$ %s debug addr 0xA588C66983a81e800Db4dF74564F09f91c026351`, version.AppName, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addrString := args[0]
			cfg := sdk.GetConfig()

			var addr []byte
			switch {
			case common.IsHexAddress(addrString):
				addr = common.HexToAddress(addrString).Bytes()
			case strings.HasPrefix(addrString, cfg.GetBech32ValidatorAddrPrefix()):
				addr, _ = sdk.ValAddressFromBech32(addrString)
			case strings.HasPrefix(addrString, cfg.GetBech32AccountAddrPrefix()):
				addr, _ = sdk.AccAddressFromBech32(addrString)
			default:
				return fmt.Errorf("expected a valid hex or bech32 address (acc prefix %s), got '%s'", cfg.GetBech32AccountAddrPrefix(), addrString)
			}

			cmd.Println("Address bytes:", addr)
			cmd.Printf("Address (hex): %s\n", bytes.HexBytes(addr))
			cmd.Printf("Address (EIP-55): %s\n", common.BytesToAddress(addr))
			cmd.Printf("Bech32 Acc: %s\n", sdk.AccAddress(addr))
			cmd.Printf("Bech32 Val: %s\n", sdk.ValAddress(addr))
			return nil
		},
	}
}

func RawBytesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "raw-bytes [raw-bytes]",
		Short:   "Convert raw bytes output (eg. [10 21 13 255]) to hex",
		Example: fmt.Sprintf(`$ %s debug raw-bytes [72 101 108 108 111 44 32 112 108 97 121 103 114 111 117 110 100]`, version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			stringBytes := args[0]
			stringBytes = strings.Trim(stringBytes, "[")
			stringBytes = strings.Trim(stringBytes, "]")
			spl := strings.Split(stringBytes, " ")

			byteArray := []byte{}
			for _, s := range spl {
				b, err := strconv.ParseInt(s, 10, 8)
				if err != nil {
					return err
				}
				byteArray = append(byteArray, byte(b))
			}
			fmt.Printf("%X\n", byteArray)
			return nil
		},
	}
}

// LegacyEIP712Cmd outputs types of legacy EIP712 typed data
func LegacyEIP712Cmd() *cobra.Command {
	return &cobra.Command{
		Use:     "legacy-eip712 [file]",
		Short:   "Output types of legacy eip712 typed data according to the given txs",
		Example: fmt.Sprintf(`$ %s debug legacy-eip712 txs.json --chain-id artelad_9000-1`, version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			stdTx, err := authclient.ReadTxFromFile(clientCtx, args[0])
			if err != nil {
				return errors.Wrap(err, "read txs from file")
			}

			txBytes, err := clientCtx.TxConfig.TxJSONEncoder()(stdTx)
			if err != nil {
				return errors.Wrap(err, "encode txs")
			}

			chainID, err := artela.ParseChainID(clientCtx.ChainID)
			if err != nil {
				return errors.Wrap(err, "invalid chain ID passed as argument")
			}

			td, err := eip712.LegacyWrapTxToTypedData(clientCtx.Codec, chainID.Uint64(), stdTx.GetMsgs()[0], txBytes, nil)
			if err != nil {
				return errors.Wrap(err, "wrap txs to typed data")
			}

			bz, err := json.Marshal(td.Map()["types"])
			if err != nil {
				return err
			}

			fmt.Println(string(bz))
			return nil
		},
	}
}

func CosmosTxHash(encodingConfig appparams.EncodingConfig) *cobra.Command {
	return &cobra.Command{
		Use:     "tx [base64 encoded tx]",
		Short:   "Get the ethereum tx of cosmos tx",
		Example: "",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBytes, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				log.Fatalf("Failed to decode base64: %v", err)
			}

			cometbftTx := cometbftTypes.Tx(txBytes)

			tx, err := encodingConfig.TxConfig.TxDecoder()(cometbftTx)
			if err != nil {
				fmt.Println(err)
			}

			for i, msg := range tx.GetMsgs() {
				ethMsg, ok := msg.(*txs.MsgEthereumTx)
				if !ok {
					fmt.Printf("message %d is not etherum tx\n", i)
					continue
				}

				fmt.Println("this is a ethereum tx:")
				ethMsg.Hash = ethMsg.AsTransaction().Hash().Hex()
				// result = append(result, ethMsg)
				ethTx := ethMsg.AsTransaction()
				fmt.Printf("	hash: %s\n	to: %s\n	value: %s\n	data: %s\n",
					ethTx.Hash().String(),
					ethTx.To().String(),
					ethTx.Value().String(),
					common.Bytes2Hex(ethTx.Data()),
				)
			}
			return nil
		},
	}
}
