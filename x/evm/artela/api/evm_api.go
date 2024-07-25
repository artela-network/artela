package api

import (
	"context"
	"errors"

	"github.com/emirpasic/gods/sets/hashset"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"

	"github.com/artela-network/artela-evm/vm"
	artelatypes "github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/artela/x/evm/states"
	types "github.com/artela-network/artela/x/evm/txs"
	asptypes "github.com/artela-network/aspect-core/types"
)

var (
	_ asptypes.EVMHostAPI = (*evmHostApi)(nil)

	evmStaticCallConstrainedJoinPoints = hashset.New(
		asptypes.PRE_CONTRACT_CALL_METHOD,
		asptypes.POST_CONTRACT_CALL_METHOD,
		asptypes.PRE_TX_EXECUTE_METHOD,
		asptypes.POST_TX_EXECUTE_METHOD,
		asptypes.OPERATION_METHOD,
		asptypes.INIT_METHOD,
	)
)

type evmHostApi struct {
	aspectCtx *artelatypes.AspectRuntimeContext
}

func (e *evmHostApi) StaticCall(ctx *asptypes.RunnerContext, request *asptypes.StaticCallRequest) (*asptypes.StaticCallResult, error) {
	if !evmStaticCallConstrainedJoinPoints.Contains(asptypes.PointCut(ctx.Point)) {
		return nil, errors.New("cannot execute static call in current join point")
	}

	from := common.BytesToAddress(request.From)
	to := common.BytesToAddress(request.To)

	var evm *vm.EVM
	ethTxCtx := e.aspectCtx.EthTxContext()
	if ethTxCtx != nil {
		// if evm is not nil, it means we are already in a tx,
		// so we can use the last evm to execute the static call
		evm = ethTxCtx.LastEvm()
	}

	// evm is still nil, we need to create a new one
	if evm == nil {
		txConfig := states.NewEmptyTxConfig(common.BytesToHash(e.aspectCtx.CosmosContext().HeaderHash()))
		stateDB := states.New(e.aspectCtx.CosmosContext(), evmKeeper, txConfig)
		evmConfig, err := evmKeeper.EVMConfig(e.aspectCtx.CosmosContext(),
			e.aspectCtx.CosmosContext().BlockHeader().ProposerAddress, evmKeeper.ChainID())
		if err != nil {
			// we need to panic here, since the evm init should not fail here
			panic(err)
		}
		evm = evmKeeper.NewEVM(e.aspectCtx.CosmosContext(), &core.Message{
			From: from,
			To:   &to,
			Data: request.Data,
		}, evmConfig, types.NewNoOpTracer(), stateDB)
	}

	// we cannot create any evm at this stage, return error
	if evm == nil {
		panic("failed to init evm")
	}

	// set the default request gas to current remaining f not specified or out of limit
	if *request.Gas == 0 || *request.Gas > ctx.Gas {
		request.Gas = &ctx.Gas
	}

	ret, gas, err := evm.StaticCall(ctx.Ctx, vm.AccountRef(from), to, request.Data, *request.Gas)
	// update gas
	ctx.Gas = gas

	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	return &asptypes.StaticCallResult{
		Ret:     ret,
		GasLeft: &gas,
		VmError: &errStr,
	}, nil
}

func (e *evmHostApi) JITCall(ctx *asptypes.RunnerContext, request *asptypes.JitInherentRequest) (*asptypes.JitInherentResponse, error) {
	// determine jit call stage
	defBool := false
	switch asptypes.PointCut(ctx.Point) {
	case asptypes.PRE_CONTRACT_CALL_METHOD, asptypes.POST_CONTRACT_CALL_METHOD:
		// FIXME: get leftover gas from last evm
		resp, gas, err := e.aspectCtx.JITManager().Submit(ctx.Ctx, ctx.AspectId, ctx.Gas, request)
		if err != nil {
			if resp == nil {
				resp = &asptypes.JitInherentResponse{}
			}

			resp.Success = &defBool
			msg := err.Error()
			resp.ErrorMsg = &msg

			log.Error("jit inherent submit fail", "err", err)
		}

		ctx.Gas = gas

		return resp, nil
	default:
		log.Error("unsupported join point for jit call", "point", ctx.Point)
		return &asptypes.JitInherentResponse{Success: &defBool}, errors.New("unsupported join point for jit call")
	}
}

func GetEvmHostInstance(ctx context.Context) (asptypes.EVMHostAPI, error) {
	aspectCtx, ok := ctx.(*artelatypes.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("GetEVMHostInstance: unwrap AspectRuntimeContext failed")
	}
	return &evmHostApi{aspectCtx}, nil
}
