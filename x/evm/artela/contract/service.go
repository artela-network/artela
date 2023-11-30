package contract

import (
	"sort"

	artela "github.com/artela-network/aspect-core/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"

	evmtypes "github.com/artela-network/artela/x/evm/artela/types"
)

type AspectService struct {
	storeKey       storetypes.StoreKey
	aspectStore    *AspectStore
	getCtxByHeight func(height int64, prove bool) (sdk.Context, error)
}

func NewAspectService(storeKey storetypes.StoreKey, getCtxByHeight func(height int64, prove bool) (sdk.Context, error)) *AspectService {
	return &AspectService{
		storeKey:       storeKey,
		aspectStore:    NewAspectStore(),
		getCtxByHeight: getCtxByHeight,
	}
}

func (service *AspectService) GetAspectOf(ctx sdk.Context, aspectId common.Address, commit bool) (*treeset.Set, error) {
	if commit {
		sdkCtx, getErr := service.getCtxByHeight(ctx.BlockHeight()-1, true)
		if getErr != nil {
			return nil, errors.Wrap(getErr, "load context by pre block height failed")
		}
		ctx = sdkCtx
	}
	aspects, err := service.aspectStore.GetAspectRefValue(ctx, service.storeKey, aspectId)
	if err != nil {
		return nil, errors.Wrap(err, "load aspect ref failed")
	}
	return aspects, nil
}

func (service *AspectService) GetAspectCode(ctx sdk.Context, aspectId common.Address, version *uint256.Int, commit bool) ([]byte, *uint256.Int) {
	if commit {
		sdkCtx, getErr := service.getCtxByHeight(ctx.BlockHeight()-1, true)
		if getErr != nil {
			return nil, nil
		}
		ctx = sdkCtx
	}
	if version == nil {
		version = service.aspectStore.GetAspectLastVersion(ctx, service.storeKey, aspectId)
	}
	return service.aspectStore.GetAspectCode(ctx, service.storeKey, aspectId, version)
}

// GetAspectForAddr BoundAspects get bound Aspects on previous block
func (service *AspectService) GetAspectForAddr(ctx sdk.Context, to common.Address, commit bool) ([]*artela.AspectCode, error) {
	if commit {
		sdkCtx, getErr := service.getCtxByHeight(ctx.BlockHeight()-1, true)
		if getErr != nil {
			return nil, errors.Wrap(getErr, "load context by pre block height failed")
		}
		ctx = sdkCtx
	}

	aspects, err := service.aspectStore.GetContractBondAspects(ctx, service.storeKey, to)
	if err != nil {
		return nil, errors.Wrap(err, "load contract aspect binding failed")
	}
	aspectCodes := make([]*artela.AspectCode, 0, len(aspects))
	if aspects == nil {
		return aspectCodes, nil
	}
	for _, aspect := range aspects {
		codeBytes, ver := service.aspectStore.GetAspectCode(ctx, service.storeKey, aspect.Id, nil)
		aspectCode := &artela.AspectCode{
			AspectId: aspect.Id.String(),
			Priority: uint32(aspect.Priority),
			Version:  ver.Uint64(),
			Code:     codeBytes,
		}
		aspectCodes = append(aspectCodes, aspectCode)
	}
	return aspectCodes, nil
}

func (service *AspectService) GetAspectForBlock(height int64) ([]*artela.AspectCode, error) {
	sdkCtx, getErr := service.getCtxByHeight(height, true)
	if getErr != nil || sdkCtx.ChainID() == "" {
		return nil, errors.Wrap(getErr, "load context by block failed GetAspectForBlock")
	}
	aspectMap, err := service.aspectStore.GetBlockLevelAspects(sdkCtx, service.storeKey)
	if err != nil {
		return nil, errors.Wrap(err, "load contract aspect binding failed")
	}
	aspectCodes := make([]*artela.AspectCode, 0, len(aspectMap))
	if aspectMap == nil {
		return aspectCodes, nil
	}
	for aspectId, number := range aspectMap {
		aspectAddr := common.HexToAddress(aspectId)
		codeBytes, ver := service.aspectStore.GetAspectCode(sdkCtx, service.storeKey, aspectAddr, nil)
		aspectCode := &artela.AspectCode{
			AspectId: aspectAddr.String(),
			Priority: uint32(number),
			Version:  ver.Uint64(),
			Code:     codeBytes,
		}
		aspectCodes = append(aspectCodes, aspectCode)
	}
	sort.Slice(aspectCodes, evmtypes.NewBindingAspectPriorityComparator(aspectCodes))
	return aspectCodes, nil
}

func (service *AspectService) GetAspectAccount(ctx sdk.Context, aspectId common.Address, commit bool) (*common.Address, error) {
	if commit {
		sdkCtx, getErr := service.getCtxByHeight(ctx.BlockHeight()-1, true)
		if getErr != nil {
			return nil, errors.Wrap(getErr, "load context by pre block height failed")
		}
		ctx = sdkCtx
	}
	if ctx.ChainID() == "" {
		return nil, errors.New("load context by  block height  failed.")
	}
	value := service.aspectStore.GetAspectPropertyValue(ctx, service.storeKey, aspectId, evmtypes.AspectAccountKey)
	if value == "" {
		return nil, errors.New("cannot get aspect account.")
	}
	address := common.HexToAddress(value)
	return &address, nil
}

func (service *AspectService) GetAspectProof(ctx sdk.Context, aspectId common.Address, commit bool) ([]byte, error) {
	if commit {
		sdkCtx, getErr := service.getCtxByHeight(ctx.BlockHeight()-1, true)
		if getErr != nil {
			return nil, errors.Wrap(getErr, "load context by pre block height failed")
		}
		ctx = sdkCtx
	}

	if ctx.ChainID() == "" {
		return nil, errors.New("load context by  block height  failed.")
	}
	value := service.aspectStore.GetAspectPropertyValue(ctx, service.storeKey, aspectId, evmtypes.AspectProofKey)
	if value == "" {
		return nil, errors.New("cannot get aspect proof.")
	}
	address := common.Hex2Bytes(value)
	return address, nil
}
