package provider

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/aspect-core/integration"
)

var _ integration.AspectProtocol = (*AspectProtocolProvider)(nil)

type AspectProtocolProvider struct {
	getEthTxContext func() *types.EthTxContext
}

func NewAspectProtocolProvider(getEthTxContext func() *types.EthTxContext) *AspectProtocolProvider {
	return &AspectProtocolProvider{
		getEthTxContext: getEthTxContext,
	}
}

func (a *AspectProtocolProvider) ChainId() *big.Int {
	return a.getEthTxContext().EvmCfg().ChainConfig.ChainID
}

func (a *AspectProtocolProvider) VMFromSnapshotState() (integration.VM, error) {
	txContext := a.getEthTxContext()
	if txContext == nil || txContext.LastEvm() == nil {
		return nil, nil
	}
	evm := a.getEthTxContext().LastEvm()
	return evm, nil
}

func (a AspectProtocolProvider) VMFromCanonicalState() (integration.VM, error) {
	// TODO implement me
	panic("implement me")
}

func (a AspectProtocolProvider) ConvertProtocolTx(txData integration.TxData) (integration.BaseLayerTx, error) {
	// TODO implement me
	panic("implement me")
}

func (a AspectProtocolProvider) EstimateGas(txData integration.TxData) (uint64, error) {
	// TODO implement me
	panic("implement me")
}

func (a AspectProtocolProvider) GasPrice() (*big.Int, error) {
	// TODO implement me
	panic("implement me")
}

func (a AspectProtocolProvider) LastBlockHeader() (integration.BlockHeader, error) {
	// TODO implement me
	panic("implement me")
}

func (a AspectProtocolProvider) NonceOf(address common.Address) (uint64, error) {
	// TODO implement me
	panic("implement me")
}

func (a AspectProtocolProvider) SubmitTxToCurrentProposal(tx integration.BaseLayerTx) error {
	// TODO implement me
	panic("implement me")
}

func (a AspectProtocolProvider) InitSystemContract(addr common.Address, code []byte, storage map[common.Hash][]byte, contractType integration.SystemContractType) error {
	// TODO implement me
	panic("implement me")
}

func (a AspectProtocolProvider) BalanceOf(address common.Address) *big.Int {
	// TODO implement me
	panic("implement me")
}
