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

type (
	contextBuilder  func(height int64, prove bool) (sdk.Context, error)
	heightRetriever func() int64
)

type AspectService struct {
	aspectStore    *AspectStore
	getCtxByHeight contextBuilder
	getHeight      heightRetriever
}

func NewAspectService(storeKey storetypes.StoreKey, getCtxByHeight contextBuilder, getHeight heightRetriever) *AspectService {
	return &AspectService{
		aspectStore:    NewAspectStore(storeKey),
		getCtxByHeight: getCtxByHeight,
		getHeight:      getHeight,
	}
}

func (service *AspectService) GetAspectOf(blockNumber int64, aspectId common.Address) (*treeset.Set, error) {
	sdkCtx, getErr := service.getCtxByHeight(blockNumber, true)
	if getErr != nil {
		return nil, errors.Wrap(getErr, "load context by pre block height failed")
	}
	aspects, err := service.aspectStore.GetAspectRefValue(sdkCtx, aspectId)
	if err != nil {
		return nil, errors.Wrap(err, "load aspect ref failed")
	}
	return aspects, nil
}

func (service *AspectService) GetAspectCode(blockNumber int64, aspectId common.Address) ([]byte, *uint256.Int) {
	sdkCtx, getErr := service.getCtxByHeight(blockNumber, true)
	if getErr != nil {
		return nil, nil
	}
	version := service.aspectStore.GetAspectLastVersion(sdkCtx, aspectId)
	return service.aspectStore.GetAspectCode(sdkCtx, aspectId, version)
}

// GetAspectForAddr BoundAspects get bound Aspects on previous block
func (service *AspectService) GetAspectForAddr(height int64, to common.Address) ([]*artela.AspectCode, error) {
	sdkCtx, getErr := service.getCtxByHeight(height, true)
	if getErr != nil {
		return nil, errors.Wrap(getErr, "load context by pre block height failed")
	}

	aspects, err := service.aspectStore.GetTxLevelAspects(sdkCtx, to)
	if err != nil {
		return nil, errors.Wrap(err, "load contract aspect binding failed")
	}
	aspectCodes := make([]*artela.AspectCode, 0, len(aspects))
	if aspects == nil {
		return aspectCodes, nil
	}
	for _, aspect := range aspects {
		codeBytes, ver := service.aspectStore.GetAspectCode(sdkCtx, aspect.Id, nil)
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

// GetAccountVerifiers gets the bound Aspect verifier for the account
func (service *AspectService) GetAccountVerifiers(height int64, to common.Address) ([]*artela.AspectCode, error) {
	sdkCtx, getErr := service.getCtxByHeight(height, true)
	if getErr != nil {
		return nil, errors.Wrap(getErr, "load context by pre block height failed")
	}

	aspects, err := service.aspectStore.GetVerificationAspects(sdkCtx, to)
	if err != nil {
		return nil, errors.Wrap(err, "load contract aspect binding failed")
	}
	aspectCodes := make([]*artela.AspectCode, 0, len(aspects))
	if aspects == nil {
		return aspectCodes, nil
	}
	for _, aspect := range aspects {
		codeBytes, ver := service.aspectStore.GetAspectCode(sdkCtx, aspect.Id, nil)
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
	aspectMap, err := service.aspectStore.GetBlockLevelAspects(sdkCtx)
	if err != nil {
		return nil, errors.Wrap(err, "load contract aspect binding failed")
	}
	aspectCodes := make([]*artela.AspectCode, 0, len(aspectMap))
	if aspectMap == nil {
		return aspectCodes, nil
	}
	for aspectId, number := range aspectMap {
		aspectAddr := common.HexToAddress(aspectId)
		codeBytes, ver := service.aspectStore.GetAspectCode(sdkCtx, aspectAddr, nil)
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

func (service *AspectService) GetAspectAccount(height int64, aspectId common.Address) (*common.Address, error) {
	sdkCtx, getErr := service.getCtxByHeight(height, true)
	if getErr != nil || sdkCtx.ChainID() == "" {
		return nil, errors.Wrap(getErr, "load context by  block height  failed.")
	}
	value := service.aspectStore.GetAspectPropertyValue(sdkCtx, aspectId, evmtypes.AspectAccountKey)
	if value == "" {
		return nil, errors.New("cannot get aspect account.")
	}
	address := common.HexToAddress(value)
	return &address, nil
}

func (service *AspectService) GetAspectProof(height int64, aspectId common.Address) ([]byte, error) {
	sdkCtx, getErr := service.getCtxByHeight(height, true)
	if getErr != nil || sdkCtx.ChainID() == "" {
		return nil, errors.Wrap(getErr, "load context by  block height  failed.")
	}
	value := service.aspectStore.GetAspectPropertyValue(sdkCtx, aspectId, evmtypes.AspectProofKey)
	if value == "" {
		return nil, errors.New("cannot get aspect proof.")
	}
	address := common.Hex2Bytes(value)
	return address, nil
}

func (service *AspectService) GetBlockHeight() int64 {
	return service.getHeight()
}
