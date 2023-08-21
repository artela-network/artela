package rpc

import (
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/artela/rpc/ethapi"
	"github.com/artela-network/artela/rpc/graphql"
	"github.com/artela-network/artela/rpc/types"
)

var defaultEthConfig = ethconfig.Config{
	SyncMode:           0,
	FilterLogCacheSize: 0,
}

type ArtelaService struct {
	clientCtx    client.Context
	cfg          *Config
	stack        types.NetworkingStack
	backend      ethapi.Backend
	filterSystem *filters.FilterSystem
}

func NewArtelaService(
	clientCtx client.Context,
	cfg *Config,
	stack types.NetworkingStack,
	am *accounts.Manager,
) *ArtelaService {
	art := &ArtelaService{
		cfg:       cfg,
		stack:     stack,
		clientCtx: clientCtx,
	}

	// TODO, add wallet for backend
	// accts :=
	// wallets := make([]accounts.Wallet, len(accs))
	// for i := 0; i < len(accs); i++ {
	// 	wallets[i] = &keystoreWallet{account: accs[i], keystore: ks}
	// }

	// Set the Backend.
	art.backend = NewBackend(art, stack.ExtRPCEnabled(), cfg, am)
	return art
}

func (art *ArtelaService) APIs() []rpc.API {
	return ethapi.GetAPIs(art.backend)
}

// Start start the ethereum services
func (art *ArtelaService) Start() error {
	if err := art.registerAPIs(); err != nil {
		return err
	}

	go func() {
		// wait for the start of the node.
		time.Sleep(2 * time.Second)

		if art.stack.Start() != nil {
			os.Exit(1)
		}
	}()

	return nil
}

func (art *ArtelaService) Shutdown() error {
	// TODO shut down
	return nil
}

// RegisterAPIs register apis and create graphql instance.
func (art *ArtelaService) registerAPIs() error {
	art.stack.RegisterAPIs(art.APIs())
	art.filterSystem = RegisterFilterAPI(art.stack, art.backend, &defaultEthConfig)

	// create graphql
	if err := graphql.New(art.stack, art.backend, art.filterSystem, []string{"*"}, []string{"*"}); err != nil {
		return err
	}

	return nil
}

func RegisterFilterAPI(stack types.NetworkingStack, backend ethapi.Backend, ethcfg *ethconfig.Config) *filters.FilterSystem {
	isLightClient := ethcfg.SyncMode == downloader.LightSync
	filterSystem := filters.NewFilterSystem(backend, filters.Config{
		LogCacheSize: ethcfg.FilterLogCacheSize,
	})
	stack.RegisterAPIs([]rpc.API{{
		Namespace: "eth",
		Service:   filters.NewFilterAPI(filterSystem, isLightClient),
	}})
	return filterSystem
}
