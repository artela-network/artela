package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/artela/common/aspect"
	"github.com/artela-network/artela/ethereum/rpc/ethapi"
	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	"github.com/artela-network/artela/ethereum/rpc/utils"
	"github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/txs"
	evmtypes "github.com/artela-network/artela/x/evm/types"
)

// Transaction pool API

func (b *BackendImpl) SendTx(ctx context.Context, signedTx *ethtypes.Transaction) error {
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

func (b *BackendImpl) GetTransaction(ctx context.Context, txHash common.Hash) (*ethapi.RPCTransaction, error) {
	res, err := b.GetTxByEthHash(txHash)
	hexTx := txHash.Hex()

	if err != nil {
		b.logger.Debug("GetTxByEthHash failed, try to getTransactionByHashPending", "error", err)
		return b.getTransactionByHashPending(txHash)
	}

	block, err := b.CosmosBlockByNumber(rpc.BlockNumber(res.Height))
	if err != nil {
		return nil, err
	}

	tx, err := b.clientCtx.TxConfig.TxDecoder()(block.Block.Txs[res.TxIndex])
	if err != nil {
		return nil, err
	}

	// the `res.MsgIndex` is inferred from tx index, should be within the bound.
	msg, ok := tx.GetMsgs()[res.MsgIndex].(*txs.MsgEthereumTx)
	if !ok {
		return nil, errors.New("invalid ethereum tx")
	}
	msg.From = res.Sender

	blockRes, err := b.CosmosBlockResultByNumber(&block.Block.Height)
	if err != nil {
		b.logger.Debug("block result not found", "height", block.Block.Height, "error", err.Error())
		return nil, nil
	}

	if res.EthTxIndex == -1 {
		// Fallback to find tx index by iterating all valid eth transactions
		msgs := b.EthMsgsFromCosmosBlock(block, blockRes)
		for i := range msgs {
			if msgs[i].Hash == hexTx {
				res.EthTxIndex = int32(i)
				break
			}
		}
	}
	// if we still unable to find the eth tx index, return error, shouldn't happen.
	if res.EthTxIndex == -1 {
		return nil, errors.New("can't find index of ethereum tx")
	}

	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle the error for pruned node.
		b.logger.Error("failed to fetch Base Fee from prunned block. Check node prunning configuration", "height", blockRes.Height, "error", err)
	}

	cfg, err := b.chainConfig()
	if err != nil {
		return nil, err
	}

	return ethapi.NewTransactionFromMsg(
		msg,
		common.BytesToHash(block.BlockID.Hash.Bytes()),
		uint64(res.Height),
		uint64(res.EthTxIndex),
		baseFee,
		cfg,
	), nil
}

func (b *BackendImpl) GetPoolTransactions() (ethtypes.Transactions, error) {
	b.logger.Debug("called eth.rpc.backend.GetPoolTransactions")
	return nil, errors.New("GetPoolTransactions is not implemented")
}

func (b *BackendImpl) GetPoolTransaction(txHash common.Hash) *ethtypes.Transaction {
	b.logger.Error("GetPoolTransaction is not implemented")
	return nil
}

func (b *BackendImpl) GetPoolNonce(_ context.Context, addr common.Address) (uint64, error) {
	return 0, errors.New("GetPoolNonce is not implemented")
}

func (b *BackendImpl) Stats() (int, int) {
	b.logger.Error("Stats is not implemented")
	return 0, 0
}

func (b *BackendImpl) TxPoolContent() (
	map[common.Address]ethtypes.Transactions, map[common.Address]ethtypes.Transactions,
) {
	b.logger.Error("TxPoolContent is not implemented")
	return nil, nil
}

func (b *BackendImpl) TxPoolContentFrom(addr common.Address) (
	ethtypes.Transactions, ethtypes.Transactions,
) {
	return nil, nil
}

func (b *BackendImpl) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.scope.Track(b.newTxsFeed.Subscribe(ch))
}

// Version returns the current ethereum protocol version.
func (b *BackendImpl) Version() string {
	v, _ := b.version()
	return v
}

func (b *BackendImpl) version() (string, error) {
	cfg, err := b.chainConfig()
	if err != nil {
		return "", err
	}

	if cfg.ChainID == nil {
		b.logger.Error("eth.rpc.backend.Version", "ChainID is nil")
		return "", errors.New("chain id is not valid")
	}
	return cfg.ChainID.String(), nil
}

func (b *BackendImpl) Engine() consensus.Engine {
	b.logger.Error("Engine is not implemented")
	return nil
}

func (b *BackendImpl) GetTxByEthHash(hash common.Hash) (*types.TxResult, error) {
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

// GetTransactionReceipt get receipt by transaction hash
func (b *BackendImpl) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	res, err := b.GetTxByEthHash(hash)
	if err != nil {
		b.logger.Debug("GetTransactionReceipt failed", "error", err)
		return nil, nil
	}
	resBlock, err := b.CosmosBlockByNumber(rpc.BlockNumber(res.Height))
	if err != nil {
		b.logger.Debug("GetTransactionReceipt failed", "error", err)
		return nil, nil
	}
	tx, err := b.clientCtx.TxConfig.TxDecoder()(resBlock.Block.Txs[res.TxIndex])
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx: %w", err)
	}
	ethMsg := tx.GetMsgs()[res.MsgIndex].(*txs.MsgEthereumTx)

	txData, err := txs.UnpackTxData(ethMsg.Data)
	if err != nil {
		return nil, err
	}

	cumulativeGasUsed := uint64(0)
	blockRes, err := b.CosmosBlockResultByNumber(&res.Height)
	if err != nil {
		b.logger.Debug("GetTransactionReceipt failed", "error", err)
		return nil, nil
	}
	for _, txResult := range blockRes.TxsResults[0:res.TxIndex] {
		cumulativeGasUsed += uint64(txResult.GasUsed)
	}
	cumulativeGasUsed += res.CumulativeGasUsed

	var status hexutil.Uint
	if res.Failed {
		status = hexutil.Uint(ethtypes.ReceiptStatusFailed)
	} else {
		status = hexutil.Uint(ethtypes.ReceiptStatusSuccessful)
	}

	// parse tx logs from events
	msgIndex := int(res.MsgIndex)
	logs, _ := utils.TxLogsFromEvents(blockRes.TxsResults[res.TxIndex].Events, msgIndex)

	if res.EthTxIndex == -1 {
		// Fallback to find tx index by iterating all valid eth transactions
		msgs := b.EthMsgsFromCosmosBlock(resBlock, blockRes)
		for i := range msgs {
			if msgs[i].Hash == hash.Hex() {
				res.EthTxIndex = int32(i) // #nosec G701
				break
			}
		}
	}
	// return error if still unable to find the eth tx index
	if res.EthTxIndex == -1 {
		return nil, errors.New("can't find index of ethereum tx")
	}

	receipt := map[string]interface{}{
		// Consensus fields: These fields are defined by the Yellow Paper
		"status":            status,
		"cumulativeGasUsed": hexutil.Uint64(cumulativeGasUsed),
		"logsBloom":         ethtypes.BytesToBloom(ethtypes.LogsBloom(logs)),
		"logs":              logs,

		// Implementation fields: These fields are added by geth when processing a transaction.
		// They are stored in the chain database.
		"transactionHash": hash,
		"contractAddress": nil,
		"gasUsed":         hexutil.Uint64(txData.GetGas()),

		// Inclusion information: These fields provide information about the inclusion of the
		// transaction corresponding to this receipt.
		"blockHash":        common.BytesToHash(resBlock.Block.Header.Hash()).Hex(),
		"blockNumber":      hexutil.Uint64(res.Height),
		"transactionIndex": hexutil.Uint64(res.EthTxIndex),

		// sender and receiver (contract or EOA) addreses
		"from": res.Sender,
		"to":   txData.GetTo(),
		"type": hexutil.Uint(ethMsg.AsTransaction().Type()),
	}

	if logs == nil {
		receipt["logs"] = [][]*ethtypes.Log{}
	}

	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if txData.GetTo() == nil || aspect.IsAspectDeploy(txData.GetTo(), txData.GetData()) {
		receipt["contractAddress"] = crypto.CreateAddress(common.HexToAddress(res.Sender), txData.GetNonce())
	}

	if dynamicTx, ok := txData.(*txs.DynamicFeeTx); ok {
		baseFee, err := b.BaseFee(blockRes)
		if err == nil {
			receipt["effectiveGasPrice"] = hexutil.Big(*dynamicTx.EffectiveGasPrice(baseFee))
		}
	}

	return receipt, nil
}

func (b *BackendImpl) queryCosmosTxIndexer(query string, txGetter func(*rpctypes.ParsedTxs) *rpctypes.ParsedTx) (*types.TxResult, error) {
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

// getTransactionByHashPending find pending tx from mempool
func (b *BackendImpl) getTransactionByHashPending(txHash common.Hash) (*ethapi.RPCTransaction, error) {
	hexTx := txHash.Hex()
	// try to find tx in mempool
	ptxs, err := b.PendingTransactions()
	if err != nil {
		b.logger.Debug("tx not found", "hash", hexTx, "error", err.Error())
		return nil, nil
	}

	for _, tx := range ptxs {
		msg, err := txs.UnwrapEthereumMsg(tx, txHash)
		if err != nil {
			// not ethereum tx
			continue
		}

		cfg, err := b.chainConfig()
		if err != nil {
			return nil, err
		}
		if msg.Hash == hexTx {
			// use zero block values since it's not included in a block yet
			rpctx := ethapi.NewTransactionFromMsg(
				msg,
				common.Hash{},
				uint64(0),
				uint64(0),
				nil,
				cfg,
			)
			return rpctx, nil
		}
	}

	b.logger.Debug("tx not found", "hash", hexTx)
	return nil, nil
}

func (b *BackendImpl) EstimateGas(ctx context.Context, args ethapi.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	blockNum := rpc.LatestBlockNumber
	if blockNrOrHash != nil {
		blockNum, _ = b.blockNumberFromCosmos(*blockNrOrHash)
	}

	bz, err := json.Marshal(&args)
	if err != nil {
		return 0, err
	}

	header, err := b.CosmosBlockByNumber(blockNum)
	if err != nil {
		// the error message imitates geth behavior
		return 0, errors.New("header not found")
	}

	req := txs.EthCallRequest{
		Args:            bz,
		GasCap:          b.RPCGasCap(),
		ProposerAddress: sdktypes.ConsAddress(header.Block.ProposerAddress),
		ChainId:         b.chainID.Int64(),
	}

	// From ContextWithHeight: if the provided height is 0,
	// it will return an empty context and the gRPC query will use
	// the latest block height for querying.
	res, err := b.queryClient.EstimateGas(rpctypes.ContextWithHeight(blockNum.Int64()), &req)
	if err != nil {
		return 0, err
	}
	return hexutil.Uint64(res.Gas), nil
}
