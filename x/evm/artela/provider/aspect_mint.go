package provider

import (
	"math/big"

	asptypes "github.com/artela-network/aspect-core/types"
	"github.com/cosmos/cosmos-sdk/aspect/cosmos"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethereum "github.com/ethereum/go-ethereum/core/types"

	"github.com/artela-network/artela/x/evm/artela/contract"
	"github.com/artela-network/artela/x/evm/artela/types"
	statedb "github.com/artela-network/artela/x/evm/states"
)

var _ cosmos.AspectCosmosProvider = (*AspectMintProvider)(nil)

type AspectMintProvider struct {
	service *contract.AspectService
}

func NewAspectMint(storeKey storetypes.StoreKey, getCtxByHeight func(height int64, prove bool) (sdk.Context, error)) *AspectMintProvider {
	service := contract.NewAspectService(storeKey, getCtxByHeight)

	return &AspectMintProvider{service: service}
}

func (j *AspectMintProvider) TxToPointRequest(sdkCtx sdk.Context, transaction *ethereum.Transaction, txIndex int64, baseFee *big.Int, innerTx *asptypes.EthStackTransaction) (*asptypes.EthTxAspect, error) {
	ethTransaction, err := asptypes.NewEthTransaction(transaction, common.BytesToHash(sdkCtx.HeaderHash().Bytes()), sdkCtx.BlockHeight(), txIndex, baseFee, sdkCtx.ChainID())
	if err != nil {
		return nil, err
	}
	return &asptypes.EthTxAspect{
		Tx:          ethTransaction,
		CurrInnerTx: innerTx,
		GasInfo:     &asptypes.GasInfo{},
	}, nil
}

func (j *AspectMintProvider) CreateTxPointRequest(sdkCtx sdk.Context, msg sdk.Msg, txIndex int64, baseFee *big.Int, innerTx *asptypes.EthStackTransaction) (*asptypes.EthTxAspect, error) {
	ethMsg := types.ConvertMsgEthereumTx(msg)
	transaction := ethMsg.AsTransaction()
	ethTransaction, err := asptypes.NewEthTransaction(transaction, common.BytesToHash(sdkCtx.HeaderHash().Bytes()), sdkCtx.BlockHeight(), txIndex, baseFee, sdkCtx.ChainID())
	if err != nil {
		return nil, err
	}
	return &asptypes.EthTxAspect{
		Tx:          ethTransaction,
		CurrInnerTx: innerTx,
		GasInfo:     &asptypes.GasInfo{},
	}, nil
}

func (j *AspectMintProvider) CreateBlockPointRequest(sdkCtx sdk.Context) *asptypes.EthBlockAspect {
	header := types.ConvertEthBlockHeader(sdkCtx.BlockHeader())
	return &asptypes.EthBlockAspect{Header: header, GasInfo: &asptypes.GasInfo{
		GasWanted: 0,
		GasUsed:   0,
		Gas:       0,
	}}
}

func (j *AspectMintProvider) CreateTxPointRequestInEvm(sdkCtx sdk.Context, msg core.Message, txConfig statedb.TxConfig, innerTx *asptypes.EthStackTransaction) *asptypes.EthTxAspect {
	chainId := sdkCtx.ChainID()
	blockHash := common.BytesToHash(sdkCtx.HeaderHash().Bytes())
	blockHeight := sdkCtx.BlockHeight()
	ethTx := asptypes.NewEthTransactionByMessage(&msg, txConfig.TxHash, chainId, blockHash, blockHeight, uint8(txConfig.TxType))
	return &asptypes.EthTxAspect{
		Tx:          ethTx,
		CurrInnerTx: innerTx,
		GasInfo: &asptypes.GasInfo{
			GasWanted: 0,
			GasUsed:   0,
			Gas:       0,
		},
	}
}

func (AspectMintProvider) FilterAspectTx(tx sdk.Msg) bool {
	if tx.ValidateBasic() != nil {
		return false
	}
	isEthTx := types.IsEthTx(tx)
	if !isEthTx {
		return false
	}
	ethTx := types.ConvertEthTx(tx)
	if ethTx == nil || ethTx.To() == nil || asptypes.IsAspectContractAddr(ethTx.To()) {
		return false
	}

	return true
}

func (j *AspectMintProvider) GetTxBondAspects(blockNum int64, address common.Address) ([]*asptypes.AspectCode, error) {
	return j.service.GetAspectForAddr(blockNum, address)
}

func (j *AspectMintProvider) GetBlockBondAspects(blockNum int64) ([]*asptypes.AspectCode, error) {
	return j.service.GetAspectForBlock(blockNum)
}

func (j *AspectMintProvider) GetAspectAccount(blockNum int64, aspectId common.Address) (*common.Address, error) {
	return j.service.GetAspectAccount(blockNum, aspectId)
}
