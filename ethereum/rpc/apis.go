package rpc

import (
	rpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/artela/ethereum/rpc/api"
	"github.com/artela-network/artela/ethereum/rpc/ethapi"
	"github.com/artela-network/artela/ethereum/rpc/filters"
)

func GetAPIs(clientCtx client.Context, wsClient *rpcclient.WSClient, logger log.Logger, apiBackend *BackendImpl) []rpc.API {
	nonceLock := new(ethapi.AddrLocker)
	return []rpc.API{
		{
			Namespace: "eth",
			Service:   ethapi.NewEthereumAPI(apiBackend, logger),
		}, {
			Namespace: "eth",
			Service:   ethapi.NewBlockChainAPI(apiBackend, logger),
		}, {
			Namespace: "eth",
			Service:   ethapi.NewTransactionAPI(apiBackend, logger, nonceLock),
		}, {
			Namespace: "txpool",
			Service:   ethapi.NewTxPoolAPI(apiBackend),
		}, {
			Namespace: "debug",
			Service:   ethapi.NewDebugAPI(apiBackend),
		}, {
			Namespace: "debug",
			Service:   api.NewDebugAPI(apiBackend),
		}, {
			Namespace: "eth",
			Service:   ethapi.NewEthereumAccountAPI(apiBackend),
		}, {
			Namespace: "personal",
			Service:   ethapi.NewPersonalAccountAPI(apiBackend, logger, nonceLock),
		}, {
			Namespace: "net",
			Service:   api.NewNetAPI(apiBackend),
		}, {
			Namespace: "eth",
			Service:   filters.NewPublicFilterAPI(logger, clientCtx, wsClient, apiBackend),
		}, {
			Namespace: "web3",
			Service:   api.NewWeb3API(apiBackend),
		},
	}
}
