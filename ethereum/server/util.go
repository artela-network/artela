package server

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	db "github.com/cometbft/cometbft-db"
	dbm "github.com/cometbft/cometbft-db"
	tmcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	rpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/cosmos/cosmos-sdk/client"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/version"
	ethlog "github.com/ethereum/go-ethereum/log"

	ethrpc "github.com/artela-network/artela/ethereum/rpc"
	"github.com/artela-network/artela/ethereum/server/config"
	ethNode "github.com/ethereum/go-ethereum/node"
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
func CreateJSONRPC(ctx *sdkserver.Context,
	clientCtx client.Context,
	tmRPCAddr,
	tmEndpoint string,
	config *config.Config,
	db db.DB,
) (*ethrpc.ArtelaService, error) {
	cfg := getRpcConfig(config)

	nodeCfg, err := getNodeConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	stack, err := ethrpc.NewNode(nodeCfg)
	if err != nil {
		return nil, err
	}

	wsClient := ConnectTmWS(tmRPCAddr, tmEndpoint, nodeCfg.Logger)

	serv := ethrpc.NewArtelaService(ctx, clientCtx, wsClient, cfg, stack, nodeCfg.Logger, db)

	// allocate separate WS connection to Tendermint
	tmWsClient := ConnectTmWS(tmRPCAddr, tmEndpoint, nodeCfg.Logger)
	wsSrv := ethrpc.NewWebsocketsServer(clientCtx, tmWsClient, config, nodeCfg.Logger)
	wsSrv.Start()

	return serv, nil
}

func getNodeConfig(ctx *sdkserver.Context, config *config.Config) (*ethNode.Config, error) {
	nodeCfg := ethrpc.DefaultGethNodeConfig()
	// if not define, use default value
	if len(config.JSONRPC.API) > 0 {
		nodeCfg.HTTPModules = config.JSONRPC.API
	}
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
	// do not start websocket
	nodeCfg.WSHost = ""
	return nodeCfg, nil
}

func getRpcConfig(config *config.Config) *ethrpc.Config {
	cfg := ethrpc.DefaultConfig()
	cfg.RPCGasCap = config.JSONRPC.GasCap
	cfg.RPCEVMTimeout = config.JSONRPC.EVMTimeout
	cfg.RPCTxFeeCap = config.JSONRPC.TxFeeCap
	cfg.AppCfg = config
	return cfg
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

func ConnectTmWS(tmRPCAddr, tmEndpoint string, logger ethlog.Logger) *rpcclient.WSClient {
	tmWsClient, err := rpcclient.NewWS(tmRPCAddr, tmEndpoint,
		rpcclient.MaxReconnectAttempts(256),
		rpcclient.ReadWait(120*time.Second),
		rpcclient.WriteWait(120*time.Second),
		rpcclient.PingPeriod(50*time.Second),
		rpcclient.OnReconnect(func() {
			logger.Debug("EVM RPC reconnects to Tendermint WS", "address", tmRPCAddr+tmEndpoint)
		}),
	)

	if err != nil {
		logger.Error(
			"Tendermint WS client could not be created",
			"address", tmRPCAddr+tmEndpoint,
			"error", err,
		)
	} else if err := tmWsClient.OnStart(); err != nil {
		logger.Error(
			"Tendermint WS client could not start",
			"address", tmRPCAddr+tmEndpoint,
			"error", err,
		)
	}

	return tmWsClient
}
