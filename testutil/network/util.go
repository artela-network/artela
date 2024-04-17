package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	rpc2 "github.com/artela-network/artela/ethereum/rpc"
	"github.com/artela-network/artela/x/evm/txs/support"
	"google.golang.org/grpc"

	cmtcfg "github.com/cometbft/cometbft/config"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	"github.com/cometbft/cometbft/rpc/client/local"
	"github.com/cometbft/cometbft/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/log"

	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	servergrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	servercmtlog "github.com/cosmos/cosmos-sdk/server/log"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	evmtypes "github.com/artela-network/artela/x/evm/types"
)

func startInProcess(cfg Config, val *Validator) error {
	logger := val.Ctx.Logger
	tmCfg := val.Ctx.Config
	tmCfg.Instrumentation.Prometheus = false

	if err := val.AppConfig.ValidateBasic(); err != nil {
		return err
	}

	nodeKey, err := p2p.LoadOrGenNodeKey(tmCfg.NodeKeyFile())
	if err != nil {
		return err
	}

	app := cfg.AppConstructor(*val)
	cmtApp := sdkserver.NewCometABCIWrapper(app)

	genDocProvider := node.DefaultGenesisDocProviderFunc(tmCfg)

	tmNode, err := node.NewNode(
		tmCfg,
		pvm.LoadOrGenFilePV(tmCfg.PrivValidatorKeyFile(), tmCfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(cmtApp),
		genDocProvider,
		cmtcfg.DefaultDBProvider,
		node.DefaultMetricsProvider(tmCfg.Instrumentation),
		servercmtlog.CometLoggerWrapper{Logger: logger.With("module", val.Moniker)},
	)
	if err != nil {
		return err
	}

	if err := tmNode.Start(); err != nil {
		return err
	}

	val.tmNode = tmNode

	if val.RPCAddress != "" {
		val.RPCClient = local.New(tmNode)
	}

	// We'll need a RPC client if the validator exposes a gRPC or REST endpoint.
	if val.APIAddress != "" || val.AppConfig.GRPC.Enable {
		val.ClientCtx = val.ClientCtx.
			WithClient(val.RPCClient)

		// Add the txs service in the gRPC router.
		app.RegisterTxService(val.ClientCtx)

		// Add the tendermint queries service in the gRPC router.
		app.RegisterTendermintService(val.ClientCtx)
	}

	var grpcServ *grpc.Server

	if val.AppConfig.GRPC.Enable {
		grpcSrv, err := servergrpc.NewGRPCServer(val.ClientCtx, app, val.AppConfig.GRPC)
		if err != nil {
			return err
		}

		servergrpc.StartGRPCServer(context.Background(), logger.With("module", "grpc-server"), val.AppConfig.GRPC, grpcSrv)

		val.grpc = grpcSrv

		/*if val.AppConfig.GRPCWeb.Enable {
			val.grpcWeb, err = servergrpc.StartGRPCWeb(grpcSrv, val.AppConfig.Config)
			if err != nil {
				return err
			}
		}*/
	}

	if val.AppConfig.API.Enable && val.APIAddress != "" {
		apiSrv := api.New(val.ClientCtx, logger.With("module", "api-server"), grpcServ)
		app.RegisterAPIRoutes(apiSrv, val.AppConfig.API)

		errCh := make(chan error)

		go func() {
			if err := apiSrv.Start(context.Background(), val.AppConfig.Config); err != nil {
				errCh <- err
			}
		}()

		select {
		case err := <-errCh:
			return err
		case <-time.After(5 * time.Second): // assume server started successfully
		}

		val.api = apiSrv
	}

	if val.AppConfig.JSONRPC.Enable && val.AppConfig.JSONRPC.Address != "" {
		/* TODO artela rpc server
		if val.Ctx == nil || val.Ctx.Viper == nil {
			return fmt.Errorf("validator %s context is nil", val.Moniker)
		}

		tmEndpoint := "/websocket"
		tmRPCAddr := fmt.Sprintf("tcp://%s", val.AppConfig.GRPC.Address)

		val.jsonrpc, val.jsonrpcDone, err = server.StartJSONRPC(val.Ctx, val.ClientCtx, tmRPCAddr, tmEndpoint, val.AppConfig, nil)
		if err != nil {
			return err
		}

		address := fmt.Sprintf("http://%s", val.AppConfig.JSONRPC.Address)

		val.JSONRPCClient, err = ethclient.Dial(address)
		if err != nil {
			return fmt.Errorf("failed to dial JSON-RPC at %s: %w", val.AppConfig.JSONRPC.Address, err)
		}*/
		cfg := rpc2.DefaultConfig()
		nodeCfg := rpc2.DefaultGethNodeConfig()

		httpEndpoint := strings.Split(val.AppConfig.JSONRPC.Address, ":")
		if len(httpEndpoint) > 0 {
			nodeCfg.HTTPHost = strings.Split(val.AppConfig.JSONRPC.Address, ":")[0]
		}
		if len(httpEndpoint) > 1 {
			nodeCfg.HTTPHost = strings.Split(val.AppConfig.JSONRPC.Address, ":")[1]
		}

		node, err := rpc2.NewNode(nodeCfg)
		if err != nil {
			panic(err)
		}

		val.artelaService = rpc2.NewArtelaService(val.Ctx, val.ClientCtx, nil, cfg, node, accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}), log.Root())
		startErr := val.artelaService.Start()
		if startErr != nil {
			return startErr
		}
	}

	return nil
}

func collectGenFiles(cfg Config, vals []*Validator, outputDir string) error {
	genTime := tmtime.Now()

	for i := 0; i < cfg.NumValidators; i++ {
		tmCfg := vals[i].Ctx.Config

		nodeDir := filepath.Join(outputDir, vals[i].Moniker, "artelad")
		gentxsDir := filepath.Join(outputDir, "gentxs")

		tmCfg.Moniker = vals[i].Moniker
		tmCfg.SetRoot(nodeDir)

		initCfg := genutiltypes.NewInitConfig(cfg.ChainID, gentxsDir, vals[i].NodeID, vals[i].PubKey)

		genFile := tmCfg.GenesisFile()
		appGenesis, err := genutiltypes.AppGenesisFromFile(genFile)
		if err != nil {
			return err
		}

		appState, err := genutil.GenAppStateFromConfig(cfg.Codec, cfg.TxConfig,
			tmCfg, initCfg, appGenesis, banktypes.GenesisBalancesIterator{},
			genutiltypes.DefaultMessageValidator, cfg.TxConfig.SigningContext().ValidatorAddressCodec())
		if err != nil {
			return err
		}

		// overwrite each validator's genesis file to have a canonical genesis time
		if err := genutil.ExportGenesisFileWithTime(genFile, cfg.ChainID, nil, appState, genTime); err != nil {
			return err
		}
	}

	return nil
}

func initGenFiles(cfg Config, genAccounts []authtypes.GenesisAccount, genBalances []banktypes.Balance, genFiles []string) error {
	// set the accounts in the genesis states
	var authGenState authtypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[authtypes.ModuleName], &authGenState)

	accounts, err := authtypes.PackAccounts(genAccounts)
	if err != nil {
		return err
	}

	authGenState.Accounts = append(authGenState.Accounts, accounts...)
	cfg.GenesisState[authtypes.ModuleName] = cfg.Codec.MustMarshalJSON(&authGenState)

	// set the balances in the genesis states
	var bankGenState banktypes.GenesisState
	bankGenState.Balances = genBalances
	cfg.GenesisState[banktypes.ModuleName] = cfg.Codec.MustMarshalJSON(&bankGenState)

	var stakingGenState stakingtypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[stakingtypes.ModuleName], &stakingGenState)

	stakingGenState.Params.BondDenom = cfg.BondDenom
	cfg.GenesisState[stakingtypes.ModuleName] = cfg.Codec.MustMarshalJSON(&stakingGenState)

	var govGenState govv1.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[govtypes.ModuleName], &govGenState)

	// nolint:staticcheck
	govGenState.DepositParams.MinDeposit[0].Denom = cfg.BondDenom
	cfg.GenesisState[govtypes.ModuleName] = cfg.Codec.MustMarshalJSON(&govGenState)

	/* TODO artela
	var inflationGenState inflationtypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[inflationtypes.ModuleName], &inflationGenState)

	inflationGenState.Params.MintDenom = cfg.BondDenom
	cfg.GenesisState[inflationtypes.ModuleName] = cfg.Codec.MustMarshalJSON(&inflationGenState)
	*/

	var crisisGenState crisistypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[crisistypes.ModuleName], &crisisGenState)

	crisisGenState.ConstantFee.Denom = cfg.BondDenom
	cfg.GenesisState[crisistypes.ModuleName] = cfg.Codec.MustMarshalJSON(&crisisGenState)

	var evmGenState support.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[evmtypes.ModuleName], &evmGenState)

	evmGenState.Params.EvmDenom = cfg.BondDenom
	cfg.GenesisState[evmtypes.ModuleName] = cfg.Codec.MustMarshalJSON(&evmGenState)

	appGenStateJSON, err := json.MarshalIndent(cfg.GenesisState, "", "  ")
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		ChainID:    cfg.ChainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < cfg.NumValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}

	return nil
}

func WriteFile(name string, dir string, contents []byte) error {
	file := filepath.Join(dir, name)

	err := tmos.EnsureDir(dir, 0o755)
	if err != nil {
		return err
	}

	return tmos.WriteFile(file, contents, 0o644)
}

// Get a free address for a test CometBFT server
// protocol is either tcp, http, etc
func FreeTCPAddr() (addr, port string, closeFn func() error, err error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", "", nil, err
	}

	closeFn = func() error {
		return l.Close()
	}

	portI := l.Addr().(*net.TCPAddr).Port
	port = fmt.Sprintf("%d", portI)
	addr = fmt.Sprintf("tcp://0.0.0.0:%s", port)
	return
}

func trapSignal(cleanupFunc func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs

		if cleanupFunc != nil {
			cleanupFunc()
		}
		exitCode := 128

		switch sig {
		case syscall.SIGINT:
			exitCode += int(syscall.SIGINT)
		case syscall.SIGTERM:
			exitCode += int(syscall.SIGTERM)
		}

		os.Exit(exitCode)
	}()
}
