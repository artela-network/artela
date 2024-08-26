package rpc

import (
	rpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/artela/ethereum/rpc/api"
	"github.com/artela-network/artela/ethereum/rpc/filters"
)

func GetAPIs(clientCtx client.Context, serverCtx *server.Context, wsClient *rpcclient.WSClient, logger log.Logger, apiBackend *BackendImpl) []rpc.API {
	nonceLock := new(api.AddrLocker)
	return []rpc.API{
		{
			Namespace: "eth",
			Service:   api.NewEthereumAPI(apiBackend, logger),
		}, {
			Namespace: "eth",
			Service:   api.NewBlockChainAPI(apiBackend, logger),
		}, {
			Namespace: "eth",
			Service:   api.NewTransactionAPI(apiBackend, logger, nonceLock),
		}, {
			Namespace: "txpool",
			Service:   api.NewTxPoolAPI(apiBackend, logger),
		}, {
			Namespace: "debug",
			Service:   api.NewDebugAPI(apiBackend, logger, serverCtx),
		}, {
			Namespace: "eth",
			Service:   api.NewEthereumAccountAPI(apiBackend),
		}, {
			Namespace: "personal",
			Service:   api.NewPersonalAccountAPI(apiBackend, logger, nonceLock),
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
