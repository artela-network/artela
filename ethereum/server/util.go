package server

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	rpc2 "github.com/artela-network/artela/ethereum/rpc"
	"github.com/artela-network/artela/ethereum/server/config"
	ethlog "github.com/ethereum/go-ethereum/log"

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

// CreateJSONRPC starts the JSON-RPC server
func CreateJSONRPC(ctx *server.Context,
	clientCtx client.Context,
	tmRPCAddr,
	tmEndpoint string,
	config *config.Config,
) (*rpc2.ArtelaService, error) {
	cfg := rpc2.DefaultConfig()
	cfg.RPCGasCap = config.JSONRPC.GasCap
	cfg.RPCEVMTimeout = config.JSONRPC.EVMTimeout
	cfg.RPCTxFeeCap = config.JSONRPC.TxFeeCap

	nodeCfg := rpc2.DefaultGethNodeConfig()
	address := strings.Split(config.JSONRPC.Address, ":")
	if len(address) > 0 {
		nodeCfg.HTTPHost = address[0]
	}
	if len(address) > 1 {
		port, err := strconv.Atoi(address[1])
		if err != nil {
			return nil, fmt.Errorf("address of JSON RPC Configuration is not valid, %w", err)
		}
		nodeCfg.HTTPPort = port
	}

	logger := ctx.Logger.With("module", "geth")
	nodeCfg.Logger = ethlog.New()
	nodeCfg.Logger.SetHandler(ethlog.FuncHandler(func(r *ethlog.Record) error {
		switch r.Lvl {
		case ethlog.LvlTrace, ethlog.LvlDebug:
			logger.Debug(r.Msg, r.Ctx...)
		case ethlog.LvlInfo, ethlog.LvlWarn:
			logger.Info(r.Msg, r.Ctx...)
		case ethlog.LvlError, ethlog.LvlCrit:
			logger.Error(r.Msg, r.Ctx...)
		}
		return nil
	}))

	stack, err := rpc2.NewNode(nodeCfg)
	if err != nil {
		return nil, err
	}

	am := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false})
	serv := rpc2.NewArtelaService(ctx, clientCtx, cfg, stack, am, nodeCfg.Logger)

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
