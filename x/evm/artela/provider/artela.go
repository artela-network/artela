package provider

import (
	"github.com/artela-network/artela/x/evm/artela/contract"
	"github.com/artela-network/artela/x/evm/artela/types"
	asptypes "github.com/artela-network/aspect-core/types"
	"github.com/cometbft/cometbft/libs/log"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/common"
)

var _ asptypes.AspectProvider = (*ArtelaProvider)(nil)

type ArtelaProvider struct {
	service        *contract.AspectService
	getCtxByHeight types.ContextBuilder
	storeKey       storetypes.StoreKey
}

func NewArtelaProvider(storeKey storetypes.StoreKey,
	getCtxByHeight types.ContextBuilder,
	getBlockHeight types.GetLastBlockHeight,
	logger log.Logger,
) *ArtelaProvider {
	service := contract.NewAspectService(storeKey, getCtxByHeight, getBlockHeight, logger)

	return &ArtelaProvider{service, getCtxByHeight, storeKey}
}

func (j *ArtelaProvider) GetTxBondAspects(blockNum int64, address common.Address, point asptypes.PointCut) ([]*asptypes.AspectCode, error) {
	heightCtx, err := j.getCtxByHeight(blockNum-1, true)
	if err != nil {
		return nil, err
	}
	return j.service.GetAspectsForJoinPoint(heightCtx, address, point, false)
}

func (j *ArtelaProvider) GetAccountVerifiers(blockNum int64, address common.Address) ([]*asptypes.AspectCode, error) {
	heightCtx, err := j.getCtxByHeight(blockNum-1, true)
	if err != nil {
		return nil, err
	}
	return j.service.GetAccountVerifiers(heightCtx, address, false)
}

func (j *ArtelaProvider) GetBlockBondAspects(blockNum int64) ([]*asptypes.AspectCode, error) {
	heightCtx, err := j.getCtxByHeight(blockNum-1, true)
	if err != nil {
		return nil, err
	}
	return j.service.GetAspectForBlock(heightCtx, false)
}

func (j *ArtelaProvider) GetAspectAccount(blockNum int64, aspectId common.Address) (*common.Address, error) {
	heightCtx, err := j.getCtxByHeight(blockNum-1, true)
	if err != nil {
		return nil, err
	}
	return j.service.GetAspectAccount(heightCtx, aspectId, false)
}

func (j *ArtelaProvider) GetLatestBlock() int64 {
	return j.service.GetBlockHeight()
}
