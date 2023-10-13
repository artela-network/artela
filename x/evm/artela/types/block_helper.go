package types

import (
	asptypes "github.com/artela-network/artelasdk/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"strings"

	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	evmtypes "github.com/artela-network/artela/x/evm/txs"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"
)

func GetEthTxsInBlock(clientCtx client.Context, sdkCtx sdk.Context, blockHeight int64, filterHexAddress string) ([]*evmtypes.MsgEthereumTx, error) {
	resBlock, err := clientCtx.Client.Block(sdkCtx, &blockHeight)
	if err != nil {
		return nil, err
	}
	// return if requested block height is greater than the current one
	if resBlock == nil || resBlock.Block == nil {
		return nil, nil
	}
	blockRes, err := clientCtx.Client.BlockResults(sdkCtx, &blockHeight)
	if err != nil {
		return nil, err
	}

	block := resBlock.Block

	txResults := blockRes.TxsResults
	var result []*evmtypes.MsgEthereumTx

	for i, tx := range block.Txs {
		// Check if tx exists on EVM by cross checking with blockResults:
		//  - Include unsuccessful tx that exceeds block gas limit
		//  - Exclude unsuccessful tx with any other error but ExceedBlockGasLimit
		if !rpctypes.TxSuccessOrExceedsBlockGasLimit(txResults[i]) {
			continue
		}

		tx, err := clientCtx.TxConfig.TxDecoder()(tx)
		if err != nil {
			continue
		}

		for _, msg := range tx.GetMsgs() {
			ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
			if !ok {
				continue
			}
			transaction := ethMsg.AsTransaction()
			ethMsg.Hash = ethMsg.AsTransaction().Hash().Hex()

			if filterHexAddress != "" {
				if strings.EqualFold(transaction.To().Hex(), filterHexAddress) {
					result = append(result, ethMsg)
				}
			} else {
				result = append(result, ethMsg)
			}
		}
	}
	return result, nil

}
func ProtoToEthBlockHeader(header *tmproto.Header) *asptypes.EthBlockHeader {
	return &asptypes.EthBlockHeader{
		ParentHash: common.BytesToHash(header.LastBlockId.Hash).String(),
		StateRoot:  common.BytesToHash(header.AppHash).String(),
		Number:     uint64(header.Height),
		Timestamp:  uint64(header.Time.Unix()),
	}
}

func TendermintToEthBlockHeader(header *types.Header) *asptypes.EthBlockHeader {
	return &asptypes.EthBlockHeader{
		ParentHash: common.BytesToHash(header.LastBlockID.Hash.Bytes()).String(),
		StateRoot:  common.BytesToHash(header.AppHash).String(),
		Number:     uint64(header.Height),
		Timestamp:  uint64(header.Time.Unix()),
	}
}

func RPCTxToEthTx(tran *rpctypes.RPCTransaction) *asptypes.EthTransaction {
	if tran == nil {
		return nil
	}
	t := &asptypes.EthTransaction{
		Nonce: uint64(tran.Nonce),
		Gas:   0,
		Input: tran.Input,
		From:  tran.From.String(),
		Hash:  tran.Hash.Bytes(),
		Type:  int32(tran.Type),
		V:     tran.V.ToInt().Bytes(),
		R:     tran.R.ToInt().Bytes(),
		S:     tran.S.ToInt().Bytes(),
	}
	if tran.ChainID != nil {
		t.ChainId = tran.ChainID.String()
	}
	if tran.GasTipCap != nil {
		t.GasTipCap = tran.GasTipCap.String()
	}
	if tran.GasFeeCap != nil {
		t.GasFeeCap = tran.GasFeeCap.String()
	}
	if tran.GasPrice != nil {
		t.GasPrice = tran.GasPrice.String()
	}
	if tran.To != nil {
		t.To = tran.To.String()
	}
	if tran.Value != nil {
		t.Value = tran.Value.String()
	}
	if tran.BlockHash != nil {
		t.BlockHash = tran.BlockHash.Bytes()
	}
	if tran.BlockNumber != nil {
		t.BlockNumber = tran.BlockNumber.ToInt().Int64()
	}
	if tran.V != nil {
		t.V = tran.V.ToInt().Bytes()
	}
	if tran.R != nil {
		t.R = tran.R.ToInt().Bytes()
	}
	if tran.S != nil {
		t.S = tran.S.ToInt().Bytes()
	}
	return t
}

func RPCTxsToEthTxs(trans []*rpctypes.RPCTransaction) []*asptypes.EthTransaction {
	txs := make([]*asptypes.EthTransaction, len(trans))

	for i, tran := range trans {
		t := RPCTxToEthTx(tran)
		txs[i] = t
	}
	return txs
}
