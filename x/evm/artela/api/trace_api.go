package api

import (
	"context"
	"errors"
	"sort"

	"github.com/emirpasic/gods/sets/hashset"
	"google.golang.org/protobuf/proto"

	"github.com/ethereum/go-ethereum/common"

	"github.com/artela-network/artela-evm/vm"
	"github.com/artela-network/artela/x/evm/artela/types"
	artelatypes "github.com/artela-network/aspect-core/types"
)

var (
	_                         artelatypes.AspectTraceHostAPI = (*aspectTraceHostAPI)(nil)
	traceJoinPointConstraints                                = hashset.New(
		artelatypes.PRE_CONTRACT_CALL_METHOD,
		artelatypes.POST_CONTRACT_CALL_METHOD,
		artelatypes.PRE_TX_EXECUTE_METHOD,
		artelatypes.POST_TX_EXECUTE_METHOD,
	)
)

const AccountBalanceMagic = ".balance"

type aspectTraceHostAPI struct {
	aspectRuntimeContext *types.AspectRuntimeContext
}

func (a *aspectTraceHostAPI) QueryStateChange(ctx *artelatypes.RunnerContext, query *artelatypes.StateChangeQuery) ([]byte, error) {
	if !traceJoinPointConstraints.Contains(artelatypes.PointCut(ctx.Point)) {
		return []byte{}, errors.New("cannot query state change in current join point")
	}

	if ctx == nil || len(query.Account) == 0 || query.StateVarName == nil {
		return []byte{}, nil
	}

	txContext := a.aspectRuntimeContext.EthTxContext()

	if txContext == nil || txContext.VmTracer() == nil {
		return []byte{}, nil
	}

	tracer := txContext.VmTracer()
	ethAddr := common.BytesToAddress(query.Account)
	stateVar := *query.StateVarName

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
	switch sc := storageChanges.(type) {
	case *vm.StorageChanges:
		changes := sc.Changes()
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
				return nil, nil
			}
			for _, state := range changes[callIdx] {
				ethStateChanges.All = append(ethStateChanges.All, &artelatypes.EthStateChange{
					Account:   call.From.Bytes(),
					Value:     state,
					CallIndex: &call.Index,
				})
			}
		}
		result = ethStateChanges
	case [][]byte:
		ethStateChangeIndices := &artelatypes.EthStateChangeIndices{
			Indices: sc,
		}

		result = ethStateChangeIndices
	default:
		result = nil
	}

	encoded, err := proto.Marshal(result)
	if err != nil {
		return []byte{}, nil
	}

	return encoded, nil
}

func (a *aspectTraceHostAPI) QueryCallTree(ctx *artelatypes.RunnerContext, query *artelatypes.CallTreeQuery) ([]byte, error) {
	if !traceJoinPointConstraints.Contains(artelatypes.PointCut(ctx.Point)) {
		return []byte{}, errors.New("cannot query call tree in current join point")
	}

	txContext := a.aspectRuntimeContext.EthTxContext()

	if txContext == nil || txContext.VmTracer() == nil || txContext.VmTracer().CallTree() == nil {
		return []byte{}, nil
	}

	var result proto.Message

	tracer := txContext.VmTracer()
	callTree := tracer.CallTree()
	if query.CallIdx == nil || *query.CallIdx < 0 {
		// for negative numbers we return the entire call tree
		ethCallTree := &artelatypes.EthCallTree{Calls: make([]*artelatypes.EthCallMessage, 0, tracer.CurrentCallIndex())}
		traverseEVMCallTree(callTree.Root(), ethCallTree)
		sort.Slice(ethCallTree.Calls, func(i, j int) bool {
			return *ethCallTree.Calls[i].Index < *ethCallTree.Calls[j].Index
		})
		result = ethCallTree
	} else {
		call := tracer.CallTree().FindCall(uint64(*query.CallIdx))
		if call != nil {
			result = callToTrace(call)
		}
	}

	encoded, err := proto.Marshal(result)
	if err != nil {
		return []byte{}, nil
	}
	return encoded, nil
}

func traverseEVMCallTree(ethCall *vm.Call, evmCallTree *artelatypes.EthCallTree) {
	if ethCall == nil {
		return
	}

	evmCallTree.Calls = append(evmCallTree.Calls, callToTrace(ethCall))

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

	gas := ethCall.Gas.Uint64()
	gasUsed := ethCall.Gas.Uint64() - ethCall.RemainingGas

	callErrMsg := ""
	if callErr != nil {
		callErrMsg = callErr.Error()
	}
	parentIndex := ethCall.ParentIndex()
	return &artelatypes.EthCallMessage{
		From:            ethCall.From.Bytes(),
		To:              to,
		Data:            ethCall.Data,
		Value:           value,
		Gas:             &gas,
		Ret:             ethCall.Ret,
		GasUsed:         &gasUsed,
		Error:           &callErrMsg,
		Index:           &ethCall.Index,
		ParentIndex:     &parentIndex,
		ChildrenIndices: ethCall.ChildrenIndices(),
	}
}

func GetAspectTraceHostInstance(ctx context.Context) (artelatypes.AspectTraceHostAPI, error) {
	aspectCtx, ok := ctx.(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetTraceHostInstance: unwrap AspectRuntimeContext failed")
	}
	return &aspectTraceHostAPI{aspectCtx}, nil
}
