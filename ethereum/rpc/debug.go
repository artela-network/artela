package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/pkg/errors"

	tmrpcclient "github.com/cometbft/cometbft/rpc/client"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	evmtxs "github.com/artela-network/artela/x/evm/txs"
	"github.com/artela-network/artela/x/evm/txs/support"
)

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (b *BackendImpl) TraceTransaction(hash common.Hash, config *support.TraceConfig) (interface{}, error) {
	// Get transaction by hash
	transaction, err := b.GetTxByEthHash(hash)
	if err != nil {
		b.logger.Debug("tx not found", "hash", hash)
		return nil, err
	}

	// check if block number is 0
	if transaction.Height == 0 {
		return nil, errors.New("genesis is not traceable")
	}

	blk, err := b.CosmosBlockByNumber(rpc.BlockNumber(transaction.Height))
	if err != nil {
		b.logger.Debug("block not found", "height", transaction.Height)
		return nil, err
	}

	// check tx index is not out of bound
	if len(blk.Block.Txs) > math.MaxUint32 {
		return nil, fmt.Errorf("tx count %d is overfloing", len(blk.Block.Txs))
	}
	txsLen := uint32(len(blk.Block.Txs)) // #nosec G701 -- checked for int overflow already
	if txsLen < transaction.TxIndex {
		b.logger.Debug("tx index out of bounds", "index", transaction.TxIndex, "hash", hash.String(), "height", blk.Block.Height)
		return nil, fmt.Errorf("transaction not included in block %v", blk.Block.Height)
	}

	var predecessors []*evmtxs.MsgEthereumTx
	for _, txBz := range blk.Block.Txs[:transaction.TxIndex] {
		tx, err := b.clientCtx.TxConfig.TxDecoder()(txBz)
		if err != nil {
			b.logger.Debug("failed to decode transaction in block", "height", blk.Block.Height, "error", err.Error())
			continue
		}
		for _, msg := range tx.GetMsgs() {
			ethMsg, ok := msg.(*evmtxs.MsgEthereumTx)
			if !ok {
				continue
			}

			predecessors = append(predecessors, ethMsg)
		}
	}

	tx, err := b.clientCtx.TxConfig.TxDecoder()(blk.Block.Txs[transaction.TxIndex])
	if err != nil {
		b.logger.Debug("tx not found", "hash", hash)
		return nil, err
	}

	// add predecessor messages in current cosmos tx
	index := int(transaction.MsgIndex) // #nosec G701
	for i := 0; i < index; i++ {
		ethMsg, ok := tx.GetMsgs()[i].(*evmtxs.MsgEthereumTx)
		if !ok {
			continue
		}
		predecessors = append(predecessors, ethMsg)
	}

	ethMessage, ok := tx.GetMsgs()[transaction.MsgIndex].(*evmtxs.MsgEthereumTx)
	if !ok {
		b.logger.Debug("invalid transaction type", "type", fmt.Sprintf("%T", tx))
		return nil, fmt.Errorf("invalid transaction type %T", tx)
	}

	nc, ok := b.clientCtx.Client.(tmrpcclient.NetworkClient)
	if !ok {
		return nil, errors.New("invalid rpc client")
	}

	cp, err := nc.ConsensusParams(b.ctx, &blk.Block.Height)
	if err != nil {
		return nil, err
	}

	traceTxRequest := evmtxs.QueryTraceTxRequest{
		Msg:             ethMessage,
		Predecessors:    predecessors,
		BlockNumber:     blk.Block.Height,
		BlockTime:       blk.Block.Time,
		BlockHash:       common.Bytes2Hex(blk.BlockID.Hash),
		ProposerAddress: sdk.ConsAddress(blk.Block.ProposerAddress),
		ChainId:         b.chainID.Int64(),
		BlockMaxGas:     cp.ConsensusParams.Block.MaxGas,
	}

	if config != nil {
		traceTxRequest.TraceConfig = config
	}

	// minus one to get the context of block beginning
	contextHeight := transaction.Height - 1
	if contextHeight < 1 {
		// 0 is a special value in `ContextWithHeight`
		contextHeight = 1
	}
	traceResult, err := b.queryClient.TraceTx(rpctypes.ContextWithHeight(contextHeight), &traceTxRequest)
	if err != nil {
		return nil, err
	}

	// Response format is unknown due to custom tracer config param
	// More information can be found here https://geth.ethereum.org/docs/dapp/tracing-filtered
	var decodedResult interface{}
	err = json.Unmarshal(traceResult.Data, &decodedResult)
	if err != nil {
		return nil, err
	}

	return decodedResult, nil
}

// TraceBlock configures a new tracer according to the provided configuration, and
// executes all the transactions contained within. The return value will be one item
// per transaction, dependent on the requested tracer.
func (b *BackendImpl) TraceBlock(height rpc.BlockNumber,
	config *support.TraceConfig,
	block *tmrpctypes.ResultBlock,
) ([]*evmtxs.TxTraceResult, error) {
	txs := block.Block.Txs
	txsLength := len(txs)

	if txsLength == 0 {
		// If there are no transactions return empty array
		return []*evmtxs.TxTraceResult{}, nil
	}

	txDecoder := b.clientCtx.TxConfig.TxDecoder()

	var txsMessages []*evmtxs.MsgEthereumTx
	for i, tx := range txs {
		decodedTx, err := txDecoder(tx)
		if err != nil {
			b.logger.Error("failed to decode transaction", "hash", txs[i].Hash(), "error", err.Error())
			continue
		}

		for _, msg := range decodedTx.GetMsgs() {
			ethMessage, ok := msg.(*evmtxs.MsgEthereumTx)
			if !ok {
				// Just considers Ethereum transactions
				continue
			}
			txsMessages = append(txsMessages, ethMessage)
		}
	}

	// minus one to get the context at the beginning of the block
	contextHeight := height - 1
	if contextHeight < 1 {
		// 0 is a special value for `ContextWithHeight`.
		contextHeight = 1
	}

	nc, ok := b.clientCtx.Client.(tmrpcclient.NetworkClient)
	if !ok {
		return nil, errors.New("invalid rpc client")
	}

	cp, err := nc.ConsensusParams(b.ctx, &block.Block.Height)
	if err != nil {
		return nil, err
	}

	traceBlockRequest := &evmtxs.QueryTraceBlockRequest{
		Txs:             txsMessages,
		TraceConfig:     config,
		BlockNumber:     block.Block.Height,
		BlockTime:       block.Block.Time,
		BlockHash:       common.Bytes2Hex(block.BlockID.Hash),
		ProposerAddress: sdk.ConsAddress(block.Block.ProposerAddress),
		ChainId:         b.chainID.Int64(),
		BlockMaxGas:     cp.ConsensusParams.Block.MaxGas,
	}

	res, err := b.queryClient.TraceBlock(rpctypes.ContextWithHeight(int64(contextHeight)), traceBlockRequest)
	if err != nil {
		return nil, err
	}

	decodedResults := make([]*evmtxs.TxTraceResult, txsLength)
	if err := json.Unmarshal(res.Data, &decodedResults); err != nil {
		return nil, err
	}

	return decodedResults, nil
}

// GetReceipts get receipts by block hash
func (b *BackendImpl) GetReceipts(ctx context.Context, hash common.Hash) (ethtypes.Receipts, error) {
	resBlock, err := b.CosmosBlockByHash(hash)
	if err != nil || resBlock == nil || resBlock.Block == nil {
		return nil, fmt.Errorf("query block failed, block hash %s, %w", hash.String(), err)
	}

	blockRes, err := b.CosmosBlockResultByNumber(&resBlock.Block.Height)
	if err != nil {
		b.logger.Debug("GetTransactionReceipt failed", "error", err)
		return nil, nil
	}

	msgs := b.EthMsgsFromCosmosBlock(resBlock, blockRes)

	receipts := make([]*ethtypes.Receipt, 0, len(msgs))
	for _, msg := range msgs {
		receipt, err := b.GetTransactionReceipt(ctx, common.HexToHash(msg.Hash))
		if err != nil || receipt == nil {
			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			}
			b.logger.Error("eth_getReceipts failed", "error", errMsg, "txHash", msg.Hash)
			continue
		}
		var contractAddress common.Address
		if receipt["contractAddress"] != nil {
			contractAddress = receipt["contractAddress"].(common.Address)
		}
		var effectiveGasPrice big.Int
		if receipt["effectiveGasPrice"] != nil {
			effectiveGasPrice = big.Int(receipt["effectiveGasPrice"].(hexutil.Big))
		}
		receipts = append(receipts, &ethtypes.Receipt{
			Type:              uint8(receipt["type"].(hexutil.Uint)),
			PostState:         []byte{},
			Status:            uint64(receipt["status"].(hexutil.Uint)),
			CumulativeGasUsed: uint64(receipt["cumulativeGasUsed"].(hexutil.Uint64)),
			Bloom:             receipt["logsBloom"].(ethtypes.Bloom),
			Logs:              receipt["logs"].([]*ethtypes.Log),
			TxHash:            receipt["transactionHash"].(common.Hash),
			ContractAddress:   contractAddress,
			GasUsed:           uint64(receipt["gasUsed"].(hexutil.Uint64)),
			EffectiveGasPrice: &effectiveGasPrice,
			BlockHash:         common.BytesToHash(resBlock.BlockID.Hash.Bytes()),
			BlockNumber:       big.NewInt(resBlock.Block.Height),
			TransactionIndex:  uint(receipt["transactionIndex"].(hexutil.Uint64)),
		})
	}
	return receipts, nil
}

func (b *BackendImpl) DBProperty(property string) (string, error) {
	if b.db == nil || b.db.Stats() == nil {
		return "", errors.New("property is not valid")
	}
	if property == "" {
		property = "leveldb.stats"
	} else if !strings.HasPrefix(property, "leveldb.") {
		property = "leveldb." + property
	}

	return (b.db.Stats())[property], nil
}

func (b *BackendImpl) DBCompact(start []byte, limit []byte) error {
	if b.db == nil {
		return errors.New("compact is not valid")
	}
	return b.db.ForceCompact(start, limit)
}
