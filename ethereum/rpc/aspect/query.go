package aspect

import (
	"context"
	"errors"

	"github.com/artela-network/artela/ethereum/rpc/ethapi"
	"github.com/artela-network/aspect-core/types"
	aspectTypes "github.com/artela-network/aspect-core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// AspectQuery implemets QueryHostAPI interface
var (
	_ aspectTypes.BlockchainAPI = (*AspectQuery)(nil)

	// there should be only one query instance in artela, use a global instance to store it.
	globalQueryInstance *AspectQuery
)

func GetBlockChainAPI(ctx context.Context) (aspectTypes.BlockchainAPI, error) {
	if globalQueryInstance == nil {
		return nil, errors.New("BlockChainAPI instance is not valid")
	}
	return globalQueryInstance, nil
}

type AspectQuery struct {
	backend ethapi.Backend
}

func SetAspectQuery(backend ethapi.Backend) {
	globalQueryInstance = &AspectQuery{backend}
}

func (query AspectQuery) GetTransactionByHash(hash []byte) *aspectTypes.Transaction {
	rpcTx, err := query.backend.GetTransaction(context.Background(), common.BytesToHash(hash))
	if err != nil || rpcTx == nil {
		return &types.Transaction{
			BlockHash:   []byte{},
			BlockNumber: new(uint64),
			Hash:        hash,
		}
	}

	var big2Uint64 = func(in *hexutil.Big) uint64 {
		if in == nil {
			return 0
		}
		return in.ToInt().Uint64()
	}

	var big2Bytes = func(in *hexutil.Big) []byte {
		if in == nil {
			return []byte{}
		}
		return in.ToInt().Bytes()
	}

	blockNumber := big2Uint64(rpcTx.BlockNumber)
	gasPrice := big2Uint64(rpcTx.GasPrice)
	gasFeeCap := big2Uint64(rpcTx.GasFeeCap)
	gasTipCap := big2Uint64(rpcTx.GasTipCap)
	var transactionIndex uint64 = 0
	if rpcTx.TransactionIndex != nil {
		transactionIndex = uint64(*rpcTx.TransactionIndex)
	}
	value := big2Bytes(rpcTx.Value)
	chainID := big2Uint64(rpcTx.ChainID)
	v := big2Bytes(rpcTx.V)
	r := big2Bytes(rpcTx.R)
	s := big2Bytes(rpcTx.S)
	return &aspectTypes.Transaction{
		BlockHash:        rpcTx.BlockHash.Bytes(),
		BlockNumber:      &blockNumber,
		From:             rpcTx.From.Bytes(),
		Gas:              (*uint64)(&rpcTx.Gas),
		GasPrice:         &gasPrice,
		GasFeeCap:        &gasFeeCap,
		GasTipCap:        &gasTipCap,
		Hash:             rpcTx.Hash.Bytes(),
		Input:            rpcTx.Input,
		Nonce:            (*uint64)(&rpcTx.Nonce),
		To:               rpcTx.To.Bytes(),
		TransactionIndex: &transactionIndex,
		Value:            value,
		Type:             (*uint64)(&rpcTx.Type),
		Accesses:         []byte{}, // we don not want this
		ChainId:          &chainID,
		V:                v,
		R:                r,
		S:                s,
	}
}
