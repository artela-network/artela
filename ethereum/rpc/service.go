package rpc

import (
	rpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/artela/ethereum/rpc/types"
)

type ArtelaService struct {
	clientCtx client.Context
	wsClient  *rpcclient.WSClient
	cfg       *Config
	stack     types.NetworkingStack
	backend   *BackendImpl
	logger    log.Logger
}

func NewArtelaService(
	ctx *server.Context,
	clientCtx client.Context,
	wsClient *rpcclient.WSClient,
	cfg *Config,
	stack types.NetworkingStack,
	logger log.Logger,
) *ArtelaService {
	art := &ArtelaService{
		cfg:       cfg,
		stack:     stack,
		clientCtx: clientCtx,
		wsClient:  wsClient,
		logger:    logger,
	}

	art.backend = NewBackend(ctx, clientCtx, art, stack.ExtRPCEnabled(), cfg, logger)
	return art
}

func (art *ArtelaService) APIs() []rpc.API {
	return GetAPIs(art.clientCtx, art.wsClient, art.logger, art.backend)
}

// Start start the ethereum JsonRPC service
func (art *ArtelaService) Start() error {
	if err := art.registerAPIs(); err != nil {
		return err
	}

	return art.stack.Start()
}

func (art *ArtelaService) Shutdown() error {
	// TODO shut down
	return nil
}

// RegisterAPIs register apis and create graphql instance.
func (art *ArtelaService) registerAPIs() error {
	art.stack.RegisterAPIs(art.APIs())
	return nil
}
