package datactx

import (
	"bytes"
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereum "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"google.golang.org/protobuf/proto"

	"github.com/artela-network/artela/x/evm/artela/types"
	aspctx "github.com/artela-network/aspect-core/context"
	artelatypes "github.com/artela-network/aspect-core/types"
)

type TxContextFieldLoader func(ethTxCtx *types.EthTxContext, tx *ethereum.Transaction) proto.Message

type TxContext struct {
	getSdkCtx             func() sdk.Context
	getEthTxContext       func() *types.EthTxContext
	receiptContentLoaders map[string]TxContextFieldLoader
}

func NewTxContext(getEthTxContext func() *types.EthTxContext,
	getSdkCtx func() sdk.Context) *TxContext {
	txContext := &TxContext{
		receiptContentLoaders: make(map[string]TxContextFieldLoader),
		getEthTxContext:       getEthTxContext,
		getSdkCtx:             getSdkCtx,
	}
	txContext.registerLoaders()

	return txContext
}

func (c *TxContext) registerLoaders() {
	loaders := c.receiptContentLoaders
	loaders[aspctx.TxType] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		txType := uint64(tx.Type())
		return &artelatypes.UintData{Data: &txType}
	}
	loaders[aspctx.TxChainId] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		if tx.ChainId() != nil {
			return &artelatypes.BytesData{Data: tx.ChainId().Bytes()}
		}
		return &artelatypes.BytesData{Data: []byte{}}
	}
	loaders[aspctx.TxAccessList] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		res := &artelatypes.EthAccessList{AccessList: make([]*artelatypes.EthAccessTuple, 0, len(tx.AccessList()))}
		if len(tx.AccessList()) == 0 {
			return res
		}
		for _, accessList := range tx.AccessList() {
			accessListMsg := &artelatypes.EthAccessTuple{
				Address: accessList.Address.Bytes(),
				StorageKeys: func() [][]byte {
					res := make([][]byte, 0, len(accessList.StorageKeys))
					for _, key := range accessList.StorageKeys {
						res = append(res, key.Bytes())
					}
					return res
				}(),
			}
			res.AccessList = append(res.AccessList, accessListMsg)
		}
		return res
	}
	loaders[aspctx.TxNonce] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		nonce := tx.Nonce()
		return &artelatypes.UintData{Data: &nonce}
	}
	loaders[aspctx.TxGas] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		gas := tx.Gas()
		return &artelatypes.UintData{Data: &gas}
	}
	loaders[aspctx.TxGasPrice] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		return &artelatypes.BytesData{Data: tx.GasPrice().Bytes()}
	}
	loaders[aspctx.TxGasTipCap] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		return &artelatypes.BytesData{Data: tx.GasTipCap().Bytes()}
	}
	loaders[aspctx.TxGasFeeCap] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		return &artelatypes.BytesData{Data: tx.GasFeeCap().Bytes()}
	}
	loaders[aspctx.TxTo] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		return &artelatypes.BytesData{Data: tx.To().Bytes()}
	}
	loaders[aspctx.TxValue] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		return &artelatypes.BytesData{Data: tx.Value().Bytes()}
	}
	loaders[aspctx.TxData] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		return &artelatypes.BytesData{Data: tx.Data()}
	}
	loaders[aspctx.TxBytes] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		raw, err := tx.MarshalBinary()
		if err != nil {
			panic(err)
		}
		return &artelatypes.BytesData{Data: raw}
	}
	loaders[aspctx.TxHash] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		return &artelatypes.BytesData{Data: tx.Hash().Bytes()}
	}
	loaders[aspctx.TxUnsignedBytes] = func(ethTxCtx *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		config := ethTxCtx.EvmCfg().ChainConfig
		blockNumber := big.NewInt(c.getSdkCtx().BlockHeight())
		blockTime := uint64(c.getSdkCtx().BlockTime().Unix())
		writer := new(bytes.Buffer)
		var err error
		switch {
		case config.IsCancun(blockNumber, blockTime):
			if tx.Type() == ethereum.BlobTxType {
				err = rlp.Encode(writer, []interface{}{
					config.ChainID,
					tx.Nonce(),
					tx.GasTipCap(),
					tx.GasFeeCap(),
					tx.Gas(),
					tx.To(),
					tx.Value(),
					tx.Data(),
					tx.AccessList(),
					tx.BlobGasFeeCap(),
					tx.BlobHashes(),
				})
			}
			fallthrough
		case config.IsLondon(blockNumber):
			if tx.Type() == ethereum.DynamicFeeTxType {
				err = rlp.Encode(writer, []interface{}{
					config.ChainID,
					tx.Nonce(),
					tx.GasTipCap(),
					tx.GasFeeCap(),
					tx.Gas(),
					tx.To(),
					tx.Value(),
					tx.Data(),
					tx.AccessList(),
				})
			}
			fallthrough
		case config.IsBerlin(blockNumber):
			switch tx.Type() {
			case ethereum.LegacyTxType:
				err = rlp.Encode(writer, []interface{}{
					tx.Nonce(),
					tx.GasPrice(),
					tx.Gas(),
					tx.To(),
					tx.Value(),
					tx.Data(),
					config.ChainID, uint(0), uint(0),
				})
			case ethereum.AccessListTxType:
				err = rlp.Encode(writer, []interface{}{
					config.ChainID,
					tx.Nonce(),
					tx.GasPrice(),
					tx.Gas(),
					tx.To(),
					tx.Value(),
					tx.Data(),
					tx.AccessList(),
				})
			}
		case config.IsEIP155(blockNumber):
			err = rlp.Encode(writer, []interface{}{
				tx.Nonce(),
				tx.GasPrice(),
				tx.Gas(),
				tx.To(),
				tx.Value(),
				tx.Data(),
				config.ChainID,
				uint(0),
				uint(0),
			})
		default:
			err = rlp.Encode(writer, []interface{}{
				tx.Nonce(),
				tx.GasPrice(),
				tx.Gas(),
				tx.To(),
				tx.Value(),
				tx.Data(),
			})
		}
		if err != nil {
			panic(err)
		}
		return &artelatypes.BytesData{Data: writer.Bytes()}
	}
	loaders[aspctx.TxUnsignedHash] = func(ethTxCtx *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		blockNumber := big.NewInt(c.getSdkCtx().BlockHeight())
		blockTime := uint64(c.getSdkCtx().BlockTime().Unix())
		config := ethTxCtx.EvmCfg().ChainConfig
		signer := ethereum.MakeSigner(config, blockNumber, blockTime)
		return &artelatypes.BytesData{Data: signer.Hash(tx).Bytes()}
	}
	loaders[aspctx.TxSigV] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		v, _, _ := tx.RawSignatureValues()
		return &artelatypes.BytesData{Data: v.Bytes()}
	}
	loaders[aspctx.TxSigS] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		_, r, _ := tx.RawSignatureValues()
		return &artelatypes.BytesData{Data: r.Bytes()}
	}
	loaders[aspctx.TxSigR] = func(_ *types.EthTxContext, tx *ethereum.Transaction) proto.Message {
		_, _, s := tx.RawSignatureValues()
		return &artelatypes.BytesData{Data: s.Bytes()}
	}
	loaders[aspctx.TxFrom] = func(ethTxCtx *types.EthTxContext, _ *ethereum.Transaction) proto.Message {
		return &artelatypes.BytesData{Data: ethTxCtx.TxFrom().Bytes()}
	}
	loaders[aspctx.TxIndex] = func(ethTxCtx *types.EthTxContext, _ *ethereum.Transaction) proto.Message {
		index := ethTxCtx.TxIndex()
		return &artelatypes.UintData{Data: &index}
	}
}

func (c *TxContext) ValueLoader(key string) ContextLoader {
	return func(ctx *artelatypes.RunnerContext) ([]byte, error) {
		if ctx == nil {
			return nil, errors.New("aspect context error, missing important information")
		}
		txContext := c.getEthTxContext()
		if txContext == nil {
			return nil, errors.New("tx context error, failed to load")
		}
		ethTx := txContext.TxContent()
		return proto.Marshal(c.receiptContentLoaders[key](txContext, ethTx))
	}
}
