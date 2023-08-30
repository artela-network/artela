package rpc

import (
	"context"
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/artela-network/artela/ethereum/types"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	ctypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"

	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	"github.com/artela-network/artela/x/evm/txs"
	evmtypes "github.com/artela-network/artela/x/evm/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

// Transaction pool API

func (b *backend) SendTx(ctx context.Context, signedTx *ethtypes.Transaction) error {
	// verify the ethereum tx
	ethereumTx := &txs.MsgEthereumTx{}
	if err := ethereumTx.FromEthereumTx(signedTx); err != nil {
		b.logger.Error("transaction converting failed", "error", err.Error())
		return err
	}

	if err := ethereumTx.ValidateBasic(); err != nil {
		b.logger.Debug("tx failed basic validation", "error", err.Error())
		return err
	}

	// Query params to use the EVM denomination
	res, err := b.queryClient.QueryClient.Params(b.ctx, &txs.QueryParamsRequest{})
	if err != nil {
		b.logger.Error("failed to query evm params", "error", err.Error())
		return err
	}

	cosmosTx, err := ethereumTx.BuildTx(b.clientCtx.TxConfig.NewTxBuilder(), res.Params.EvmDenom)
	if err != nil {
		b.logger.Error("failed to build cosmos tx", "error", err.Error())
		return err
	}

	// Encode transaction by default Tx encoder
	txBytes, err := b.clientCtx.TxConfig.TxEncoder()(cosmosTx)
	if err != nil {
		b.logger.Error("failed to encode eth tx using default encoder", "error", err.Error())
		return err
	}

	// txHash := ethereumTx.AsTransaction().Hash()

	syncCtx := b.clientCtx.WithBroadcastMode(flags.BroadcastSync)
	rsp, err := syncCtx.BroadcastTx(txBytes)
	if rsp != nil && rsp.Code != 0 {
		err = errorsmod.ABCIError(rsp.Codespace, rsp.Code, rsp.RawLog)
	}
	if err != nil {
		b.logger.Error("failed to broadcast tx", "error", err.Error())
		return err
	}

	return nil
}

func (b *backend) GetTransaction(
	_ context.Context, txHash common.Hash,
) (*ethtypes.Transaction, common.Hash, uint64, uint64, error) {
	b.logger.Debug("called eth.rpc.backend.GetTransaction", "tx_hash", txHash)
	return nil, common.Hash{}, 0, 0, nil
}

func (b *backend) GetPoolTransactions() (ctypes.Transactions, error) {
	b.logger.Debug("called eth.rpc.backend.GetPoolTransactions")
	return nil, nil
}

func (b *backend) GetPoolTransaction(txHash common.Hash) *ethtypes.Transaction {
	return nil
}

func (b *backend) GetPoolNonce(_ context.Context, addr common.Address) (uint64, error) {
	return 0, nil
}

func (b *backend) Stats() (int, int) {
	return 0, 0
}

func (b *backend) TxPoolContent() (
	map[common.Address]ctypes.Transactions, map[common.Address]ctypes.Transactions,
) {
	return nil, nil
}

func (b *backend) TxPoolContentFrom(addr common.Address) (
	ctypes.Transactions, ctypes.Transactions,
) {
	return nil, nil
}

func (b *backend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.scope.Track(b.newTxsFeed.Subscribe(ch))
}

// Version returns the current ethereum protocol version.
func (b *backend) Version() string {
	chainID := b.ChainConfig().ChainID
	if chainID == nil {
		b.logger.Error("eth.rpc.backend.Version", "ChainID is nil")
		return "-1"
	}
	return chainID.String()
}

func (b *backend) Engine() consensus.Engine {
	return nil
}

func (b *backend) GetTxByEthHash(hash common.Hash) (*types.TxResult, error) {
	// if b.indexer != nil {
	// 	return b.indexer.GetByTxHash(hash)
	// }

	// fallback to tendermint tx indexer
	query := fmt.Sprintf("%s.%s='%s'", evmtypes.TypeMsgEthereumTx, evmtypes.AttributeKeyEthereumTxHash, hash.Hex())
	txResult, err := b.queryCosmosTxIndexer(query, func(txs *rpctypes.ParsedTxs) *rpctypes.ParsedTx {
		return txs.GetTxByHash(hash)
	})
	if err != nil {
		return nil, fmt.Errorf("GetTxByEthHash %s, %w", hash.Hex(), err)
	}
	return txResult, nil
}

func (b *backend) queryCosmosTxIndexer(query string, txGetter func(*rpctypes.ParsedTxs) *rpctypes.ParsedTx) (*types.TxResult, error) {
	resTxs, err := b.clientCtx.Client.TxSearch(b.ctx, query, false, nil, nil, "")
	if err != nil {
		return nil, err
	}
	if len(resTxs.Txs) == 0 {
		return nil, errors.New("ethereum tx not found")
	}
	txResult := resTxs.Txs[0]
	if !rpctypes.TxSuccessOrExceedsBlockGasLimit(&txResult.TxResult) {
		return nil, errors.New("invalid ethereum tx")
	}

	var tx sdktypes.Tx
	if txResult.TxResult.Code != 0 {
		// it's only needed when the tx exceeds block gas limit
		tx, err = b.clientCtx.TxConfig.TxDecoder()(txResult.Tx)
		if err != nil {
			return nil, fmt.Errorf("invalid ethereum tx, %w", err)
		}
	}

	return rpctypes.ParseTxIndexerResult(txResult, tx, txGetter)
}

func (b *backend) txResult(ctx context.Context, hash common.Hash, prove bool) (*tmrpctypes.ResultTx, error) {
	return b.clientCtx.Client.Tx(ctx, hash.Bytes(), prove)
}
