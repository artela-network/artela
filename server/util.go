package server

import (
	"io"
	"os"
	"path/filepath"

	"github.com/artela-network/artela/rpc"
	"github.com/artela-network/artela/server/config"
	dbm "github.com/cometbft/cometbft-db"
	tmcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/spf13/cobra"
)

// add server commands
func AddCommands(rootCmd *cobra.Command, defaultNodeHome string, appCreator types.AppCreator, appExport types.AppExporter, addStartFlags types.ModuleInitFlags) {
	tendermintCmd := &cobra.Command{
		Use:     "tendermint",
		Aliases: []string{"comet", "cometbft"},
		Short:   "Tendermint subcommands",
	}

	tendermintCmd.AddCommand(
		sdkserver.ShowNodeIDCmd(),
		sdkserver.ShowValidatorCmd(),
		sdkserver.ShowAddressCmd(),
		sdkserver.VersionCmd(),
		tmcmd.ResetAllCmd,
		tmcmd.ResetStateCmd,
		sdkserver.BootstrapStateCmd(appCreator),
	)

	startCmd := StartCmd(appCreator, defaultNodeHome)
	addStartFlags(startCmd)

	rootCmd.AddCommand(
		startCmd,
		tendermintCmd,
		sdkserver.ExportCmd(appExport, defaultNodeHome),
		version.NewVersionCommand(),
		sdkserver.NewRollbackCmd(appCreator, defaultNodeHome),
	)
}

// StartJSONRPC starts the JSON-RPC server
func CreateJSONRPC(ctx *server.Context,
	clientCtx client.Context,
	tmRPCAddr,
	tmEndpoint string,
	config *config.Config,
) (*rpc.ArtelaService, error) {
	cfg := rpc.DefaultConfig()
	cfg.RPCGasCap = config.JSONRPC.GasCap
	cfg.RPCEVMTimeout = config.JSONRPC.EVMTimeout
	cfg.RPCTxFeeCap = config.JSONRPC.TxFeeCap

	nodeCfg := rpc.DefaultGethNodeConfig()
	// TODO, config end point
	// nodeCfg.HTTPHost = tmRPCAddr
	// nodeCfg. =
	stack, err := rpc.NewNode(nodeCfg)
	if err != nil {
		return nil, err
	}

	am := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false})
	serv := rpc.NewArtelaService(ctx, clientCtx, cfg, stack, am)

	return serv, nil
}

func openDB(rootDir string, backendType dbm.BackendType) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return dbm.NewDB("application", backendType, dataDir)
}

func openTraceWriter(traceWriterFile string) (w io.WriteCloser, err error) {
	if traceWriterFile == "" {
		return
	}
	return os.OpenFile(
		traceWriterFile,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0o666,
	)
}
