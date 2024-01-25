package rpc

import (
	rpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/artela/ethereum/rpc/aspect"
	"github.com/artela-network/artela/ethereum/rpc/types"
	aspectTypes "github.com/artela-network/aspect-core/types"
)

// nolint:unused
var defaultEthConfig = ethconfig.Config{
	SyncMode:           0,
	FilterLogCacheSize: 0,
}

type ArtelaService struct {
	clientCtx client.Context
	wsClient  *rpcclient.WSClient
	cfg       *Config
	stack     types.NetworkingStack
	backend   *BackendImpl
	// nolint:unused
	filterSystem *filters.FilterSystem
	logger       log.Logger
}

func NewArtelaService(
	ctx *server.Context,
	clientCtx client.Context,
	wsClient *rpcclient.WSClient,
	cfg *Config,
	stack types.NetworkingStack,
	am *accounts.Manager,
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
	aspect.SetAspectQuery(art.backend)
	aspectTypes.GetBlockchainHook = aspect.GetBlockChainAPI
	return art
}

func Accounts(clientCtx client.Context) ([]common.Address, error) {
	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	infos, err := clientCtx.Keyring.List()
	if err != nil {
		return addresses, err
	}

	for _, info := range infos {
		pubKey, err := info.GetPubKey()
		if err != nil {
			return nil, err
		}
		addressBytes := pubKey.Address().Bytes()
		addresses = append(addresses, common.BytesToAddress(addressBytes))
	}

	return addresses, nil
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
	// art.filterSystem = RegisterFilterAPI(art.stack, art.backend, &defaultEthConfig)

	// create graphql
	// if err := graphql.New(art.stack, art.backend, art.filterSystem, []string{"*"}, []string{"*"}); err != nil {
	// 	return err
	// }

	return nil
}

// func RegisterFilterAPI(stack types.NetworkingStack, backend ethapi.Backend, ethcfg *ethconfig.Config) *filters.FilterSystem {
// 	isLightClient := ethcfg.SyncMode == downloader.LightSync
// 	filterSystem := filters.NewFilterSystem(backend, filters.Config{
// 		LogCacheSize: ethcfg.FilterLogCacheSize,
// 	})
// 	stack.RegisterAPIs([]rpc.API{{
// 		Namespace: "eth",
// 		Service:   filters.NewFilterAPI(filterSystem, isLightClient),
// 	}})
// 	return filterSystem
// }
