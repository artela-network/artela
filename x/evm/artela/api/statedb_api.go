package api

import (
	"context"

	"github.com/artela-network/artela-evm/vm"
	"github.com/artela-network/artela/x/evm/artela/types"
	artelatypes "github.com/artela-network/aspect-core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	_ artelatypes.StateDbHostApi = (*stateDbHostApi)(nil)
	// stateDbInstance *stateDbHostApi
)

type stateDbHostApi struct {
	getLastStateDB func() vm.StateDB
}

func NewStateDbApi(getLastStateDB func() vm.StateDB) artelatypes.StateDbHostApi {
	stateDbInstance := &stateDbHostApi{
		getLastStateDB: getLastStateDB,
	}
	return stateDbInstance
}

func GetStateApiInstance(ctx context.Context) (artelatypes.StateDbHostApi, error) {
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetStateApiInstance: unwrap AspectRuntimeContext failed")
	}
	return NewStateDbApi(aspectCtx.StateDb), nil
}

// GetBalance(request AddressQueryRequest) StringDataResponse
func (s *stateDbHostApi) GetBalance(ctx *artelatypes.RunnerContext, addressEquals string) string {
	if addressEquals == "" {
		return ""
	}
	address := common.HexToAddress(addressEquals)
	balance := s.getLastStateDB().GetBalance(address)
	balanceStr := artelatypes.Ternary(balance != nil, func() string { return balance.String() }, "0")
	return balanceStr
}

// GetState retrieves a value from the given account's storage trie.
func (s *stateDbHostApi) GetState(ctx *artelatypes.RunnerContext, addressEquals, hashEquals string) string {
	if hashEquals == "" || addressEquals == "" {
		return ""
	}
	address := common.HexToAddress(addressEquals)
	hash := common.HexToHash(hashEquals)

	result := s.getLastStateDB().GetState(address, hash)
	return result.String()
}

// GetRefund() IntDataResponse
func (s *stateDbHostApi) GetRefund(ctx *artelatypes.RunnerContext) uint64 {
	return s.getLastStateDB().GetRefund()
}

// GetCodeHash(request AddressQueryRequest) StringDataResponse
func (s *stateDbHostApi) GetCodeHash(ctx *artelatypes.RunnerContext, addressEquals string) string {
	if addressEquals == "" {
		return ""
	}
	address := common.HexToAddress(addressEquals)
	result := s.getLastStateDB().GetCodeHash(address)
	return result.String()
}

// GetNonce(request AddressQueryRequest) IntDataResponse
func (s *stateDbHostApi) GetNonce(ctx *artelatypes.RunnerContext, addressEquals string) uint64 {
	if addressEquals == "" {
		return 0
	}
	address := common.HexToAddress(addressEquals)
	return s.getLastStateDB().GetNonce(address)
}
