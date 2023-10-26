package rpc

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"

	rpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/cosmos/cosmos-sdk/client"

	"github.com/artela-network/artela/ethereum/rpc/ethapi"
	"github.com/artela-network/artela/ethereum/rpc/filters"
	"github.com/artela-network/artela/ethereum/types"
)

func GetAPIs(clientCtx client.Context, wsClient *rpcclient.WSClient, logger log.Logger, apiBackend *BackendImpl) []rpc.API {
	chainID, err := types.ParseChainID(clientCtx.ChainID)
	if err != nil {
		panic(err)
	}

	nonceLock := new(ethapi.AddrLocker)
	return []rpc.API{
		{
			Namespace: "eth",
			Service:   ethapi.NewEthereumAPI(apiBackend),
		}, {
			Namespace: "eth",
			Service:   ethapi.NewBlockChainAPI(apiBackend),
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
			Namespace: "eth",
			Service:   ethapi.NewEthereumAccountAPI(apiBackend),
		}, {
			Namespace: "personal",
			Service:   ethapi.NewPersonalAccountAPI(apiBackend, logger, nonceLock),
		}, {
			Namespace: "net",
			Service:   ethapi.NewNetAPI(nil, chainID.Uint64()),
		}, {
			Namespace: "eth",
			Service:   filters.NewPublicFilterAPI(logger, clientCtx, wsClient, apiBackend),
		},
	}
}
