package provider

import (
	"context"
	"errors"

	"github.com/cometbft/cometbft/libs/log"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artela-network/artela/x/evm/artela/contract"
	"github.com/artela-network/artela/x/evm/artela/types"
	asptypes "github.com/artela-network/aspect-core/types"
)

var _ asptypes.AspectProvider = (*ArtelaProvider)(nil)

type ArtelaProvider struct {
	service  *contract.AspectService
	storeKey storetypes.StoreKey
}

func NewArtelaProvider(storeKey storetypes.StoreKey,
	getBlockHeight types.GetLastBlockHeight,
	logger log.Logger,
) *ArtelaProvider {
	service := contract.NewAspectService(storeKey, getBlockHeight, logger)

	return &ArtelaProvider{service, storeKey}
}

func (j *ArtelaProvider) GetTxBondAspects(ctx context.Context, address common.Address, point asptypes.PointCut) ([]*asptypes.AspectCode, error) {
	if ctx == nil {
		return nil, errors.New("invalid Context")
	}
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("failed to unwrap AspectRuntimeContext from context.Context")
	}
	return j.service.GetAspectsForJoinPoint(aspectCtx.CosmosContext(), address, point)
}

func (j *ArtelaProvider) GetAccountVerifiers(ctx context.Context, address common.Address) ([]*asptypes.AspectCode, error) {
	if ctx == nil {
		return nil, errors.New("invalid Context")
	}
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("failed to unwrap AspectRuntimeContext from context.Context")
	}
	return j.service.GetAccountVerifiers(aspectCtx.CosmosContext(), address)
}

func (j *ArtelaProvider) GetLatestBlock() int64 {
	return j.service.GetBlockHeight()
}
