package provider

import (
	"context"
	"errors"
	"slices"

	asptypes "github.com/artela-network/aspect-core/types"

	"github.com/artela-network/artela/x/aspect/store"
	aspectmoduletypes "github.com/artela-network/artela/x/aspect/types"
	"github.com/artela-network/artela/x/evm/artela/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/common"
)

var _ asptypes.AspectProvider = (*ArtelaProvider)(nil)

type ArtelaProvider struct {
	getBlockHeight types.GetLastBlockHeight

	evmStoreKey    storetypes.StoreKey
	aspectStoreKey storetypes.StoreKey
}

func NewArtelaProvider(
	evmStoreKey storetypes.StoreKey,
	aspectStoreKey storetypes.StoreKey,
	getBlockHeight types.GetLastBlockHeight,
) *ArtelaProvider {
	return &ArtelaProvider{
		evmStoreKey:    evmStoreKey,
		aspectStoreKey: aspectStoreKey,
		getBlockHeight: getBlockHeight,
	}
}

func (j *ArtelaProvider) GetTxBondAspects(ctx context.Context, address common.Address, point asptypes.PointCut) ([]*asptypes.AspectCode, error) {
	return j.getCodes(ctx, address, point)
}

func (j *ArtelaProvider) GetAccountVerifiers(ctx context.Context, address common.Address) ([]*asptypes.AspectCode, error) {
	return j.getCodes(ctx, address, asptypes.VERIFY_TX)
}

func (j *ArtelaProvider) GetLatestBlock() int64 {
	return j.getBlockHeight()
}

func (j *ArtelaProvider) getCodes(ctx context.Context, address common.Address, point asptypes.PointCut) ([]*asptypes.AspectCode, error) {
	if ctx == nil {
		return nil, errors.New("invalid Context")
	}
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("failed to unwrap AspectRuntimeContext from context.Context")
	}

	accountStore, _, err := store.GetAccountStore(j.buildAccountStoreCtx(aspectCtx, address))
	if err != nil {
		return nil, err
	}

	bindings, err := accountStore.LoadAccountBoundAspects(aspectmoduletypes.NewJoinPointFilter(point))
	if err != nil {
		return nil, err
	}

	codes := make([]*asptypes.AspectCode, 0, len(bindings))
	for _, binding := range bindings {
		metaStore, _, err := store.GetAspectMetaStore(j.buildAspectStoreCtx(aspectCtx, binding.Account))
		if err != nil {
			return nil, err
		}
		code, err := metaStore.GetCode(binding.Version)
		if err != nil {
			return nil, err
		}

		var isExpectedJP bool
		if binding.JoinPoint == 0 {
			meta, err := metaStore.GetVersionMeta(binding.Version)
			if err != nil {
				return nil, err
			}
			isExpectedJP = asptypes.CanExecPoint(int64(meta.JoinPoint), point)
		} else {
			isExpectedJP = asptypes.CanExecPoint(int64(binding.JoinPoint), point)
		}

		// filter matched aspect with given join point
		if !isExpectedJP {
			continue
		}

		codes = append(codes, &asptypes.AspectCode{
			AspectId: binding.Account.Hex(),
			Version:  binding.Version,
			Priority: binding.Priority,
			Code:     code,
		})
	}

	// sort the codes by priority
	slices.SortFunc(codes, func(a, b *asptypes.AspectCode) int {
		if a.Priority == b.Priority {
			return 0
		} else if a.Priority < b.Priority {
			return -1
		} else {
			return 1
		}
	})

	return codes, nil
}

func (j *ArtelaProvider) buildAspectStoreCtx(ctx *types.AspectRuntimeContext, aspectID common.Address) *aspectmoduletypes.AspectStoreContext {
	return &aspectmoduletypes.AspectStoreContext{
		StoreContext: aspectmoduletypes.NewGasFreeStoreContext(ctx.CosmosContext(), j.aspectStoreKey, j.evmStoreKey),
		AspectID:     aspectID,
	}
}

func (j *ArtelaProvider) buildAccountStoreCtx(ctx *types.AspectRuntimeContext, account common.Address) *aspectmoduletypes.AccountStoreContext {
	return &aspectmoduletypes.AccountStoreContext{
		StoreContext: aspectmoduletypes.NewGasFreeStoreContext(ctx.CosmosContext(), j.aspectStoreKey, j.evmStoreKey),
		Account:      account,
	}
}
