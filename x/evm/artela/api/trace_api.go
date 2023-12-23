package api

import (
	"context"
	"errors"
	"github.com/artela-network/artela-evm/vm"
	"github.com/artela-network/artela/x/evm/artela/types"
	artelatypes "github.com/artela-network/aspect-core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/proto"
	"sort"
)

var (
	_ artelatypes.AspectTraceHostAPI = (*aspectTraceHostAPI)(nil)
)

const AccountBalanceMagic = ".balance"

type aspectTraceHostAPI struct {
	aspectRuntimeContext *types.AspectRuntimeContext
}

func (a *aspectTraceHostAPI) QueryStateChange(ctx *artelatypes.RunnerContext, query *artelatypes.StateChangeQuery) []byte {
	if ctx == nil || len(query.Account) == 0 || len(query.StateVarName) == 0 {
		return []byte{}
	}

	txContext := a.aspectRuntimeContext.EthTxContext()

	if txContext == nil || txContext.VmTracer() == nil {
		return []byte{}
	}

	tracer := txContext.VmTracer()
	ethAddr := common.BytesToAddress(query.Account)
	stateVar := query.StateVarName

	var storageChanges interface{}
	if stateVar == AccountBalanceMagic {
		storageChanges = txContext.VmTracer().StateChanges().Balance(ethAddr)
	} else {
		storageKey := txContext.VmTracer().StateChanges().FindKeyIndices(ethAddr, stateVar, query.Indices...)
		if storageKey == nil {
			storageChanges = nil
		} else if storageKey.NodeType() == vm.DataNode {
			storageChanges = storageKey.Changes()
		} else if storageKey.NodeType() == vm.BranchNode {
			storageChanges = storageKey.ChildrenIndices()
		}
	}

	var result proto.Message
	// nolint:gosimple
	switch storageChanges.(type) {
	// nolint:gosimple
	case *vm.StorageChanges:
		changes := storageChanges.(*vm.StorageChanges).Changes()
		ethStateChanges := &artelatypes.EthStateChanges{
			All: make([]*artelatypes.EthStateChange, 0, len(changes)),
		}

		callIndices := make([]uint64, 0, len(changes))
		for callIdx := range changes {
			callIndices = append(callIndices, callIdx)
		}
		sort.Slice(callIndices, func(i, j int) bool {
			return callIndices[i] < callIndices[j]
		})

		for _, callIdx := range callIndices {
			call := tracer.CallTree().FindCall(callIdx)
			if call == nil {
				return nil
			}
			for _, state := range changes[callIdx] {
				ethStateChanges.All = append(ethStateChanges.All, &artelatypes.EthStateChange{
					Account:   call.From.Bytes(),
					Value:     state,
					CallIndex: call.Index,
				})
			}
		}
		result = ethStateChanges
	// nolint:gosimple
	case [][]byte:
		indices := storageChanges.([][]byte)
		ethStateChangeIndices := &artelatypes.EthStateChangeIndices{
			Indices: indices,
		}

		result = ethStateChangeIndices
	default:
		result = nil
	}

	encoded, err := proto.Marshal(result)
	if err != nil {
		return []byte{}
	}

	return encoded
}

func (a *aspectTraceHostAPI) QueryCallTree(ctx *artelatypes.RunnerContext, query *artelatypes.CallTreeQuery) []byte {
	if ctx == nil {
		return []byte{}
	}
	txContext := a.aspectRuntimeContext.EthTxContext()

	if txContext == nil || txContext.VmTracer() == nil || txContext.VmTracer().CallTree() == nil {
		return []byte{}
	}

	var result proto.Message

	tracer := txContext.VmTracer()
	callTree := tracer.CallTree()
	if query.CallIdx < 0 {
		// for negative numbers we return the entire call tree
		ethCallTree := &artelatypes.EthCallTree{Calls: make(map[uint64]*artelatypes.EthCallMessage)}
		traverseEVMCallTree(callTree.Root(), ethCallTree)
		result = ethCallTree
	} else {
		call := tracer.CallTree().FindCall(uint64(query.CallIdx))
		if call != nil {
			result = callToTrace(call)
		}
	}

	encoded, err := proto.Marshal(result)
	if err != nil {
		return []byte{}
	}
	return encoded
}

func traverseEVMCallTree(ethCall *vm.Call, evmCallTree *artelatypes.EthCallTree) {
	if evmCallTree == nil {
		evmCallTree = &artelatypes.EthCallTree{Calls: make(map[uint64]*artelatypes.EthCallMessage)}
	}
	if ethCall == nil {
		return
	}

	evmCallTree.GetCalls()[ethCall.Index] = callToTrace(ethCall)

	children := ethCall.Children
	if children == nil {
		return
	}
	for _, call := range children {
		traverseEVMCallTree(call, evmCallTree)
	}
}

func callToTrace(ethCall *vm.Call) *artelatypes.EthCallMessage {
	var to []byte
	if ethCall.To != nil {
		to = ethCall.To.Bytes()
	}

	var value []byte
	if ethCall.Value != nil {
		value = ethCall.Value.Bytes()
	}

	var callErr error
	if ethCall.Err != nil {
		callErr = ethCall.Err
	}

	return &artelatypes.EthCallMessage{
		From:            ethCall.From.Bytes(),
		To:              to,
		Data:            ethCall.Data,
		Value:           value,
		Gas:             ethCall.Gas.Uint64(),
		Ret:             ethCall.Ret,
		GasUsed:         ethCall.Gas.Uint64() - ethCall.RemainingGas,
		Error:           callErr.Error(),
		Index:           ethCall.Index,
		ParentIndex:     ethCall.ParentIndex(),
		ChildrenIndices: ethCall.ChildrenIndices(),
	}
}

func GetAspectTraceHostInstance(ctx context.Context) (artelatypes.AspectTraceHostAPI, error) {
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetAspectRuntimeContextHostInstance: unwrap AspectRuntimeContext failed")
	}
	return &aspectTraceHostAPI{aspectCtx}, nil
}
