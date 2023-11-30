package datactx

import (
	"encoding/hex"
	"math/big"
	"sort"
	"strconv"

	"github.com/artela-network/artela-evm/vm"
	artelatypes "github.com/artela-network/aspect-core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"

	"github.com/artela-network/artela/x/evm/artela/types"
)

const AccountBalanceMagic = ".balance"

type StateChanges struct {
	getEthTxContext func() *types.EthTxContext
}

func NewStateChanges(getEthTxContext func() *types.EthTxContext) *StateChanges {
	return &StateChanges{getEthTxContext: getEthTxContext}
}

// Execute retrieves the state changes of given state variable and it subordinates indices
func (c StateChanges) Execute(sdkContext sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil || len(keys) < 2 {
		return artelatypes.NewContextQueryResponse(false, "verification failed")
	}
	txContext := c.getEthTxContext()

	if txContext == nil || txContext.VmTracer() == nil {
		return artelatypes.NewContextQueryResponse(false, "tracer is nil")
	}

	tracer := txContext.VmTracer()
	// parse input keys
	addr := keys[0]
	stateVar := keys[1]
	var indices [][]byte
	if len(keys) > 2 {
		indices = make([][]byte, 0, len(keys)-2)
		for i := 2; i < len(keys); i++ {
			decoded, err := hex.DecodeString(keys[i])
			if err != nil {
				return artelatypes.NewContextQueryResponse(false, "invalid index")
			}
			indices = append(indices, decoded)
		}
	}

	ethAddr := common.HexToAddress(addr)

	var storageChanges interface{}
	if stateVar == AccountBalanceMagic {
		storageChanges = txContext.VmTracer().StateChanges().Balance(ethAddr)
	} else {
		storageKey := txContext.VmTracer().StateChanges().FindKeyIndices(ethAddr, stateVar, indices...)
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
					Account:   call.From.Hex(),
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

	contextQueryResponse := artelatypes.NewContextQueryResponse(true, "success")
	contextQueryResponse.SetData(result)
	return contextQueryResponse
}

type TxReceipt struct {
	getEthTxContext func() *types.EthTxContext
}

func NewTxReceipt(getEthTxContext func() *types.EthTxContext) *TxReceipt {
	return &TxReceipt{getEthTxContext: getEthTxContext}
}

func (c TxReceipt) Execute(sdkContext sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}
	txContext := c.getEthTxContext()

	contextQueryResponse := artelatypes.NewContextQueryResponse(true, "basic validate failed.")
	if txContext == nil || txContext.Receipt() == nil {
		return contextQueryResponse
	}
	receipt := txContext.Receipt()

	// set data
	receiptMsg := &artelatypes.EthReceipt{
		Type:              uint32(receipt.Type),
		PostState:         receipt.PostState,
		Status:            receipt.Status,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		LogsBloom:         receipt.Bloom.Bytes(),
		Logs:              artelatypes.ConvertEthLogs(receipt.Logs),
		TxHash: artelatypes.Ternary(receipt.TxHash != common.Hash{}, func() string {
			return receipt.TxHash.String()
		}, ""),
		ContractAddress: artelatypes.Ternary(receipt.ContractAddress != common.Address{}, func() string {
			return receipt.ContractAddress.String()
		}, ""),
		GasUsed: receipt.GasUsed,
		BlockHash: artelatypes.Ternary(receipt.BlockHash != common.Hash{}, func() string {
			return receipt.BlockHash.String()
		}, ""),
		BlockNumber: artelatypes.Ternary(receipt.BlockNumber != nil, func() string {
			return receipt.BlockNumber.String()
		}, ""),
		TransactionIndex: uint32(receipt.TransactionIndex),
	}

	contextQueryResponse.GetResult().Message = "success"
	contextQueryResponse.SetData(receiptMsg)
	return contextQueryResponse
}

type ExtProperties struct {
	getEthTxContext func() *types.EthTxContext
}

func NewExtProperties(getEthTxContext func() *types.EthTxContext) *ExtProperties {
	return &ExtProperties{getEthTxContext: getEthTxContext}
}

func (c ExtProperties) Execute(sdkContext sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}
	txContext := c.getEthTxContext()

	contextQueryResponse := artelatypes.NewContextQueryResponse(true, "basic validate failed.")
	if txContext == nil || txContext.ExtProperties() == nil {
		return contextQueryResponse
	}
	extProperties := txContext.ExtProperties()

	dataMap := make(map[string]string)

	for k, v := range extProperties {
		dataMap[k] = v.(string)
	}
	// set data
	txProperty := &artelatypes.TxExtProperty{
		Property: dataMap,
	}
	contextQueryResponse.GetResult().Message = "success"
	contextQueryResponse.SetData(txProperty)
	return contextQueryResponse
}

type TxGasMeter struct{}

func NewTxGasMeter() *TxGasMeter {
	return &TxGasMeter{}
}

func (c TxGasMeter) Execute(sdkContext sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}

	contextQueryResponse := artelatypes.NewContextQueryResponse(true, "basic validate failed.")
	if sdkContext.GasMeter() == nil {
		return contextQueryResponse
	}
	meter := sdkContext.GasMeter()

	// set data
	gasMsg := &artelatypes.GasMeter{
		GasConsumed:        meter.GasConsumed(),
		GasConsumedToLimit: meter.GasConsumedToLimit(),
		GasRemaining:       meter.GasRemaining(),
		Limit:              meter.Limit(),
	}
	contextQueryResponse.GetResult().Message = "success"
	contextQueryResponse.SetData(gasMsg)
	return contextQueryResponse
}

type TxContent struct {
	getEthTxContext func() *types.EthTxContext
}

func NewTxContent(getEthTxContext func() *types.EthTxContext) *TxContent {
	return &TxContent{getEthTxContext: getEthTxContext}
}

func (c TxContent) Execute(sdkCtx sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}
	txContext := c.getEthTxContext()

	contextQueryResponse := artelatypes.NewContextQueryResponse(true, "basic validate failed.")
	if txContext == nil || txContext.ExtProperties() == nil {
		return contextQueryResponse
	}
	ethTx := txContext.TxContent()
	blockHash := common.BytesToHash(sdkCtx.HeaderHash())
	blockNumber := sdkCtx.BlockHeight()
	index := int64(-1)
	baseFee := big.NewInt(-1)
	chainID := sdkCtx.ChainID()
	// set data
	txMsg, errs := artelatypes.NewEthTransaction(ethTx, blockHash, blockNumber, index, baseFee, chainID)
	if errs != nil {
		contextQueryResponse.GetResult().Message = errs.Error()
		contextQueryResponse.GetResult().Success = false
		return contextQueryResponse
	}
	
	contextQueryResponse.GetResult().Message = "success"
	contextQueryResponse.SetData(txMsg)
	return contextQueryResponse
}

type TxAspectContent struct {
	getAspectContext func() *types.AspectContext
}

func NewTxAspectContent(getAspectContext func() *types.AspectContext) *TxAspectContent {
	return &TxAspectContent{getAspectContext: getAspectContext}
}

// getAspectContext data
func (c TxAspectContent) Execute(sdkCtx sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil || len(keys) == 0 || keys[0] == "" {
		return artelatypes.NewContextQueryResponse(false, "")
	}
	aspId := ctx.AspectId.String()
	if len(keys) >= 2 && keys[1] != "" {
		aspId = keys[1]
	}

	get := c.getAspectContext().Get(aspId, keys[0])
	if get == "" {
		return artelatypes.NewContextQueryResponse(false, "get empty.")
	} else {
		response := artelatypes.NewContextQueryResponse(true, "success")
		data := &artelatypes.StringData{Data: get}
		response.SetData(data)
		return response
	}
}

// TxCallTree
type TxCallTree struct {
	getEthTxContext func() *types.EthTxContext
}

func NewTxCallTree(getEthTxContext func() *types.EthTxContext) *TxCallTree {
	return &TxCallTree{getEthTxContext: getEthTxContext}
}

// getAspectContext data
func (c TxCallTree) Execute(sdkContext sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}
	txContext := c.getEthTxContext()

	contextQueryResponse := artelatypes.NewContextQueryResponse(true, "basic validate failed.")
	if txContext == nil || txContext.VmTracer() == nil || txContext.VmTracer().CallTree() == nil {
		return contextQueryResponse
	}
	tracer := txContext.VmTracer()
	callTree := tracer.CallTree()
	if len(keys) == 1 {
		sIndex := keys[0]
		index, err := strconv.ParseUint(sIndex, 10, 64)
		if err != nil {
			return contextQueryResponse
		}
		call := tracer.CallTree().FindCall(index)
		if call != nil {
			callMap := make(map[uint64]*artelatypes.EthStackTransaction)

			ethStackTx := &artelatypes.EthStackTransaction{
				From: call.From.String(),
				To: artelatypes.Ternary(call.To != common.Address{}, func() string {
					return call.To.String()
				}, ""),
				Data: call.Data,
				Value: artelatypes.Ternary(call.Value != nil, func() string {
					return call.Value.String()
				}, ""),
				Gas:         call.Gas.String(),
				Ret:         call.Ret,
				LeftOverGas: call.RemainingGas,
				Index:       call.Index,
				ParentIndex: artelatypes.Ternary(call.Parent != nil, func() int64 {
					return int64(call.Parent.Index)
				}, -1),
				ChildrenIndex: call.ChildrenIndices(),
			}
			callMap[index] = ethStackTx
			stacks := &artelatypes.EthCallStacks{Calls: callMap}
			contextQueryResponse.GetResult().Message = "success"
			contextQueryResponse.SetData(stacks)
		}

	} else {
		ethCallTree := &artelatypes.EthCallStacks{Calls: make(map[uint64]*artelatypes.EthStackTransaction)}
		traverseEVMCallTree(callTree.Root(), ethCallTree)
		contextQueryResponse.GetResult().Message = "success"
		contextQueryResponse.SetData(ethCallTree)
	}

	return contextQueryResponse
}

func traverseEVMCallTree(innerTx *vm.Call, evmCallTree *artelatypes.EthCallStacks) {
	if evmCallTree == nil {
		evmCallTree = &artelatypes.EthCallStacks{Calls: make(map[uint64]*artelatypes.EthStackTransaction)}
	}
	if innerTx == nil {
		return
	}
	ethStackTx := &artelatypes.EthStackTransaction{
		From: innerTx.From.String(),
		To: artelatypes.Ternary(innerTx.To != common.Address{}, func() string {
			return innerTx.To.String()
		}, ""),
		Data: innerTx.Data,
		Value: artelatypes.Ternary(innerTx.Value != nil, func() string {
			return innerTx.Value.String()
		}, ""),
		Gas:         innerTx.Gas.String(),
		Ret:         innerTx.Ret,
		LeftOverGas: innerTx.RemainingGas,
		Index:       innerTx.Index,
		ParentIndex: artelatypes.Ternary(innerTx.Parent != nil, func() int64 {
			return int64(innerTx.Parent.Index)
		}, -1),
		ChildrenIndex: innerTx.ChildrenIndices(),
	}
	evmCallTree.GetCalls()[innerTx.Index] = ethStackTx
	innerTxs := innerTx.Children
	if innerTxs == nil {
		return
	}
	for _, tx := range innerTxs {
		traverseEVMCallTree(tx, evmCallTree)
	}
}
