package contract

import (
	"math/big"
	"sort"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	artela "github.com/artela-network/aspect-core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"

	evmtypes "github.com/artela-network/artela/x/evm/artela/types"
)

type (
	heightRetriever func() int64
)

type AspectService struct {
	aspectStore *AspectStore
	getHeight   evmtypes.GetLastBlockHeight
}

func NewAspectService(storeKey storetypes.StoreKey,
	getHeight evmtypes.GetLastBlockHeight, logger log.Logger) *AspectService {
	return &AspectService{
		aspectStore: NewAspectStore(storeKey, logger),
		getHeight:   getHeight,
	}
}

func (service *AspectService) GetAspectOf(ctx sdk.Context, aspectId common.Address) (*treeset.Set, error) {

	aspects, err := service.aspectStore.GetAspectRefValue(ctx, aspectId)
	if err != nil {
		return nil, errors.Wrap(err, "load aspect ref failed")
	}
	return aspects, nil
}

func (service *AspectService) GetAspectCode(ctx sdk.Context, aspectId common.Address, version *uint256.Int) ([]byte, *uint256.Int) {

	if version == nil {
		version = service.aspectStore.GetAspectLastVersion(ctx, aspectId)
	}
	return service.aspectStore.GetAspectCode(ctx, aspectId, version)
}
func (service *AspectService) GetBoundAspectForAddr(sdkCtx sdk.Context, to common.Address) ([]*artela.AspectCode, error) {
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

// GetAspectsForJoinPoint BoundAspects get bound Aspects on previous block
func (service *AspectService) GetAspectsForJoinPoint(ctx sdk.Context, to common.Address, cut artela.PointCut) ([]*artela.AspectCode, error) {

	aspects, err := service.aspectStore.GetTxLevelAspects(ctx, to)

	if err != nil {
		return nil, errors.Wrap(err, "load contract aspect binding failed")
	}

	aspectCodes := make([]*artela.AspectCode, 0, len(aspects))
	if aspects == nil {
		return aspectCodes, nil
	}
	for _, aspect := range aspects {
		// check if the Join point has run permissions
		jp := service.aspectStore.GetAspectJP(ctx, aspect.Id, nil)
		if !artela.CanExecPoint(jp.Int64(), cut) {
			continue
		}
		codeBytes, ver := service.aspectStore.GetAspectCode(ctx, aspect.Id, nil)
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
func (service *AspectService) GetAccountVerifiers(ctx sdk.Context, to common.Address) ([]*artela.AspectCode, error) {

	aspects, err := service.aspectStore.GetVerificationAspects(ctx, to)
	if err != nil {
		return nil, errors.Wrap(err, "load contract aspect binding failed")
	}
	aspectCodes := make([]*artela.AspectCode, 0, len(aspects))
	if aspects == nil {
		return aspectCodes, nil
	}
	for _, aspect := range aspects {
		// check if the verify point has run permissions
		jp := service.aspectStore.GetAspectJP(ctx, aspect.Id, nil)
		if !artela.CanExecPoint(jp.Int64(), artela.VERIFY_TX) {
			continue
		}
		codeBytes, ver := service.aspectStore.GetAspectCode(ctx, aspect.Id, nil)
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

func (service *AspectService) GetAspectForBlock(ctx sdk.Context) ([]*artela.AspectCode, error) {

	aspectMap, err := service.aspectStore.GetBlockLevelAspects(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "load contract aspect binding failed")
	}
	aspectCodes := make([]*artela.AspectCode, 0, len(aspectMap))
	if aspectMap == nil {
		return aspectCodes, nil
	}
	for aspectId, number := range aspectMap {
		aspectAddr := common.HexToAddress(aspectId)

		// check if the join point has run permissions
		jp := service.aspectStore.GetAspectJP(ctx, aspectAddr, nil)
		blockInitCheck := artela.CanExecPoint(jp.Int64(), artela.ON_BLOCK_INITIALIZE_METHOD)
		blockFinalCheck := artela.CanExecPoint(jp.Int64(), artela.ON_BLOCK_FINALIZE_METHOD)
		if !(blockInitCheck || blockFinalCheck) {
			continue
		}

		codeBytes, ver := service.aspectStore.GetAspectCode(ctx, aspectAddr, nil)
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

func (service *AspectService) GetAspectAccount(ctx sdk.Context, aspectId common.Address) (*common.Address, error) {

	if ctx.ChainID() == "" {
		return nil, errors.New("chainID is empty.")
	}
	value := service.aspectStore.GetAspectPropertyValue(ctx, aspectId, evmtypes.AspectAccountKey)
	if value == nil {
		return nil, errors.New("cannot get aspect account.")
	}
	address := common.HexToAddress(string(value))
	return &address, nil
}

func (service *AspectService) GetAspectProof(ctx sdk.Context, aspectId common.Address) ([]byte, error) {

	if ctx.ChainID() == "" {
		return nil, errors.New("chainID is empty.")
	}
	value := service.aspectStore.GetAspectPropertyValue(ctx, aspectId, evmtypes.AspectProofKey)
	if value == nil {
		return nil, errors.New("cannot get aspect proof.")
	}
	address := common.Hex2Bytes(string(value))
	return address, nil
}

func (service *AspectService) GetBlockHeight() int64 {
	return service.getHeight()
}

func (service *AspectService) GetAspectJoinPoint(ctx sdk.Context, aspectId common.Address, version *uint256.Int) *big.Int {

	if version == nil {
		version = service.aspectStore.GetAspectLastVersion(ctx, aspectId)
	}
	return service.aspectStore.GetAspectJP(ctx, aspectId, version)
}
