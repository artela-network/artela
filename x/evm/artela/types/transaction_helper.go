package types

import (
	statedb "github.com/artela-network/artela/x/evm/states"
	evmtxs "github.com/artela-network/artela/x/evm/txs"

	asptypes "github.com/artela-network/artelasdk/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func EthTransactionByMsg(sdkCtx sdk.Context, message core.Message, txConfig statedb.TxConfig) *asptypes.EthTransaction {
	return &asptypes.EthTransaction{
		ChainId:     sdkCtx.ChainID(),
		Nonce:       message.Nonce,
		GasTipCap:   asptypes.Ternary(message.GasTipCap != nil, func() string { return message.GasTipCap.String() }, "0"),
		GasFeeCap:   asptypes.Ternary(message.GasFeeCap != nil, func() string { return message.GasFeeCap.String() }, "0"),
		Gas:         message.GasLimit,
		GasPrice:    asptypes.Ternary(message.GasPrice != nil, func() string { return message.GasPrice.String() }, "0"),
		To:          asptypes.Ternary(message.To != nil, func() string { return message.To.String() }, ""),
		Value:       asptypes.Ternary(message.Value != nil, func() string { return message.Value.String() }, "0"),
		Input:       message.Data,
		AccessList:  asptypes.ConvertTuples(message.AccessList),
		BlockHash:   sdkCtx.HeaderHash(),
		BlockNumber: sdkCtx.BlockHeight(),
		From:        message.From.String(),
		Hash:        txConfig.TxHash.Bytes(),
		Type:        int32(txConfig.TxType),
	}
}

func ConvertEthTransaction(sdkCtx sdk.Context, ethMsg *evmtxs.MsgEthereumTx) *asptypes.EthTransaction {
	ethTx := ethMsg.AsTransaction()
	v, r, s := ethTx.RawSignatureValues()

	return &asptypes.EthTransaction{
		ChainId:     sdkCtx.ChainID(),
		Nonce:       ethTx.Nonce(),
		GasTipCap:   asptypes.Ternary(ethTx.GasTipCap() != nil, func() string { return ethTx.GasTipCap().String() }, "0"),
		GasFeeCap:   asptypes.Ternary(ethTx.GasFeeCap() != nil, func() string { return ethTx.GasFeeCap().String() }, "0"),
		Gas:         ethTx.Gas(),
		GasPrice:    asptypes.Ternary(ethTx.GasPrice() != nil, func() string { return ethTx.GasPrice().String() }, "0"),
		To:          asptypes.Ternary(ethTx.To() != nil, func() string { return ethTx.To().String() }, ""),
		Value:       asptypes.Ternary(ethTx.Value() != nil, func() string { return ethTx.Value().String() }, "0"),
		Input:       ethTx.Data(),
		AccessList:  asptypes.ConvertTuples(ethTx.AccessList()),
		BlockHash:   asptypes.Ternary(sdkCtx.HeaderHash() != nil, func() []byte { return sdkCtx.HeaderHash().Bytes() }, []byte{0}),
		BlockNumber: sdkCtx.BlockHeight(),
		From:        ethMsg.From,
		Hash:        common.Hex2Bytes(ethMsg.Hash),
		Type:        int32(ethTx.Type()),
		V:           asptypes.Ternary(v != nil, func() []byte { return v.Bytes() }, []byte{0}),
		R:           asptypes.Ternary(r != nil, func() []byte { return r.Bytes() }, []byte{0}),
		S:           asptypes.Ternary(s != nil, func() []byte { return s.Bytes() }, []byte{0}),
	}
}

func ConvertEthBlockHeader(header tmproto.Header) *asptypes.EthBlockHeader {

	return &asptypes.EthBlockHeader{
		ParentHash:       asptypes.Ternary(header.LastBlockId.Hash != nil && len(header.LastBlockId.Hash) > 0, func() string { return common.BytesToHash(header.LastBlockId.Hash).String() }, ""),
		UncleHash:        "",
		Coinbase:         "",
		StateRoot:        asptypes.Ternary(header.AppHash != nil && len(header.AppHash) > 0, func() string { return common.BytesToHash(header.AppHash).String() }, ""),
		TransactionsRoot: asptypes.Ternary(header.DataHash != nil && len(header.DataHash) > 0, func() string { return common.BytesToHash(header.DataHash).String() }, ""),
		ReceiptHash:      "",
		Difficulty:       0,
		Number:           uint64(header.Height),
		GasLimit:         0,
		GasUsed:          0,
		Timestamp:        uint64(header.Time.UTC().Unix()),
		ExtraData:        nil,
		MixHash:          nil,
		Nonce:            0,
		BaseFee:          0,
	}

}
func ConvertEthTx(msg sdk.Msg) *ethtypes.Transaction {
	if IsEthTx(msg) {
		ethTx, _ := msg.(*evmtxs.MsgEthereumTx)
		return ethTx.AsTransaction()
	}
	return nil
}
func ConvertMsgEthereumTx(msg sdk.Msg) *evmtxs.MsgEthereumTx {
	if IsEthTx(msg) {
		ethTx, _ := msg.(*evmtxs.MsgEthereumTx)
		return ethTx
	}
	return nil
}
func IsEthTx(msg sdk.Msg) bool {
	_, ok := msg.(*evmtxs.MsgEthereumTx)
	return ok
}
