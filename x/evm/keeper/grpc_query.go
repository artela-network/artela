package keeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/math"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	ethereum "github.com/ethereum/go-ethereum/core/types"
	ethparams "github.com/ethereum/go-ethereum/params"

	"github.com/artela-network/artela-evm/tracers"
	"github.com/artela-network/artela-evm/tracers/logger"
	"github.com/artela-network/artela-evm/vm"
	artela "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/artela/provider"
	artelatypes "github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/artela/x/evm/states"
	"github.com/artela-network/artela/x/evm/txs"
	"github.com/artela-network/artela/x/evm/txs/support"
	"github.com/artela-network/artela/x/evm/types"
	inherent "github.com/artela-network/aspect-core/chaincoreext/jit_inherent"
)

var _ txs.QueryServer = Keeper{}

const (
	defaultTraceTimeout = 5 * time.Second
)

// Account implements the Query/StateAccount gRPC method
func (k Keeper) Account(c context.Context, req *txs.QueryAccountRequest) (*txs.QueryAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := artela.ValidateAddress(req.Address); err != nil {
		return nil, status.Error(
			codes.InvalidArgument, err.Error(),
		)
	}

	addr := common.HexToAddress(req.Address)

	ctx := cosmos.UnwrapSDKContext(c)
	acct := k.GetAccountOrEmpty(ctx, addr)

	return &txs.QueryAccountResponse{
		Balance:  acct.Balance.String(),
		CodeHash: common.BytesToHash(acct.CodeHash).Hex(),
		Nonce:    acct.Nonce,
	}, nil
}

func (k Keeper) CosmosAccount(c context.Context, req *txs.QueryCosmosAccountRequest) (*txs.QueryCosmosAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := artela.ValidateAddress(req.Address); err != nil {
		return nil, status.Error(
			codes.InvalidArgument, err.Error(),
		)
	}

	ctx := cosmos.UnwrapSDKContext(c)

	ethAddr := common.HexToAddress(req.Address)
	cosmosAddr := cosmos.AccAddress(ethAddr.Bytes())

	account := k.accountKeeper.GetAccount(ctx, cosmosAddr)
	res := txs.QueryCosmosAccountResponse{
		CosmosAddress: cosmosAddr.String(),
	}

	if account != nil {
		res.Sequence = account.GetSequence()
		res.AccountNumber = account.GetAccountNumber()
	}

	return &res, nil
}

// ValidatorAccount implements the Query/Balance gRPC method
func (k Keeper) ValidatorAccount(c context.Context, req *txs.QueryValidatorAccountRequest) (*txs.QueryValidatorAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	consAddr, err := cosmos.ConsAddressFromBech32(req.ConsAddress)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument, err.Error(),
		)
	}

	ctx := cosmos.UnwrapSDKContext(c)

	validator, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		return nil, fmt.Errorf("validator not found for %s", consAddr.String())
	}

	accAddr := cosmos.AccAddress(validator.GetOperator())

	res := txs.QueryValidatorAccountResponse{
		AccountAddress: accAddr.String(),
	}

	account := k.accountKeeper.GetAccount(ctx, accAddr)
	if account != nil {
		res.Sequence = account.GetSequence()
		res.AccountNumber = account.GetAccountNumber()
	}

	return &res, nil
}

// Balance implements the Query/Balance gRPC method
func (k Keeper) Balance(c context.Context, req *txs.QueryBalanceRequest) (*txs.QueryBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := artela.ValidateAddress(req.Address); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			types.ErrZeroAddress.Error(),
		)
	}

	ctx := cosmos.UnwrapSDKContext(c)

	balanceInt := k.GetBalance(ctx, common.HexToAddress(req.Address))

	return &txs.QueryBalanceResponse{
		Balance: balanceInt.String(),
	}, nil
}

// Storage implements the Query/Storage gRPC method
func (k Keeper) Storage(c context.Context, req *txs.QueryStorageRequest) (*txs.QueryStorageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := artela.ValidateAddress(req.Address); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			types.ErrZeroAddress.Error(),
		)
	}

	ctx := cosmos.UnwrapSDKContext(c)

	address := common.HexToAddress(req.Address)
	key := common.HexToHash(req.Key)

	state := k.GetState(ctx, address, key)
	stateHex := state.Hex()

	return &txs.QueryStorageResponse{
		Value: stateHex,
	}, nil
}

// Code implements the Query/Code gRPC method
func (k Keeper) Code(c context.Context, req *txs.QueryCodeRequest) (*txs.QueryCodeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := artela.ValidateAddress(req.Address); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			types.ErrZeroAddress.Error(),
		)
	}

	ctx := cosmos.UnwrapSDKContext(c)

	address := common.HexToAddress(req.Address)
	acct := k.GetAccountWithoutBalance(ctx, address)

	var code []byte
	if acct != nil && acct.IsContract() {
		code = k.GetCode(ctx, common.BytesToHash(acct.CodeHash))
	}

	return &txs.QueryCodeResponse{
		Code: code,
	}, nil
}

// Params implements the Query/Params gRPC method
func (k Keeper) Params(c context.Context, _ *txs.QueryParamsRequest) (*txs.QueryParamsResponse, error) {
	ctx := cosmos.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &txs.QueryParamsResponse{
		Params: params,
	}, nil
}

// EthCall implements eth_call rpc api.
func (k Keeper) EthCall(c context.Context, req *txs.EthCallRequest) (*txs.MsgEthereumTxResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			k.logger.Error("EthCall panic", "err", r)
		}
	}()
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := cosmos.UnwrapSDKContext(c)

	var args txs.TransactionArgs
	err := json.Unmarshal(req.Args, &args)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	chainID, err := getChainID(ctx, req.ChainId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	cfg, err := k.EVMConfig(ctx, GetProposerAddress(ctx, req.ProposerAddress), chainID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// ApplyMessageWithConfig expect correct nonce set in msg
	nonce := k.GetNonce(ctx, args.GetFrom())
	args.Nonce = (*hexutil.Uint64)(&nonce)

	msg, err := args.ToMessage(req.GasCap, cfg.BaseFee)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	txConfig := states.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash()))
	// Aspect Runtime Context Lifecycle: create aspect context.
	// This marks the beginning of running an aspect of EthCall, creating the aspect context,
	// and establishing the link with the SDK context.
	ctx, aspectCtx := k.WithAspectContext(ctx, args.ToTransaction().AsEthCallTransaction(), cfg,
		artelatypes.NewEthBlockContextFromQuery(ctx, k.clientContext))
	defer aspectCtx.Destroy()

	// pass false to not commit StateDB
	isCustomVerification := len(args.GetValidationData()) > 0
	res, err := k.ApplyMessageWithConfig(ctx, aspectCtx, msg, nil, false, cfg, txConfig, isCustomVerification)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

// EstimateGas implements eth_estimateGas rpc api.
func (k Keeper) EstimateGas(c context.Context, req *txs.EthCallRequest) (*txs.EstimateGasResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := cosmos.UnwrapSDKContext(c)
	chainID, err := getChainID(ctx, req.ChainId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if req.GasCap < ethparams.TxGas {
		return nil, status.Error(codes.InvalidArgument, "gas cap cannot be lower than 21,000")
	}

	var args txs.TransactionArgs
	err = json.Unmarshal(req.Args, &args)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Binary search the gas requirement, as it may be higher than the amount used
	var (
		lo     = ethparams.TxGas - 1
		hi     uint64
		gasCap uint64
	)

	// Determine the highest gas limit can be used during the estimation.
	if args.Gas != nil && uint64(*args.Gas) >= ethparams.TxGas {
		hi = uint64(*args.Gas)
	} else {
		// Query block gas limit
		params := ctx.ConsensusParams()
		if params != nil && params.Block != nil && params.Block.MaxGas > 0 {
			hi = uint64(params.Block.MaxGas)
		} else {
			hi = req.GasCap
		}
	}

	// Recap the highest gas allowance with specified gascap.
	if req.GasCap != 0 && hi > req.GasCap {
		hi = req.GasCap
	}
	txMsg := args.ToTransaction()

	gasCap = hi
	cfg, err := k.EVMConfig(ctx, GetProposerAddress(ctx, req.ProposerAddress), chainID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to load evm config")
	}

	// ApplyMessageWithConfig expect correct nonce set in msg
	nonce := k.GetNonce(ctx, args.GetFrom())
	args.Nonce = (*hexutil.Uint64)(&nonce)

	txConfig := states.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash().Bytes()))

	// convert the txs args to an ethereum message
	msg, err := args.ToMessage(req.GasCap, cfg.BaseFee)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// NOTE: the errors from the executable below should be consistent with go-ethereum,
	// so we don't wrap them with the gRPC status code

	isCustomVerification := len(args.GetValidationData()) > 0

	// Create a helper to check if a gas allowance results in an executable txs
	executable := func(gas uint64) (vmError bool, rsp *txs.MsgEthereumTxResponse, err error) {
		// need to create a cache context here to avoid state change affecting each other
		tmpCtx, _ := ctx.CacheContext()
		// Aspect Runtime Context Lifecycle: create aspect context.
		// This marks the beginning of running an aspect of EstimateGas, creating the aspect context,
		// and establishing the link with the SDK context.
		cosmosCtx, aspectCtx := k.WithAspectContext(tmpCtx, txMsg.AsTransaction(), cfg,
			artelatypes.NewEthBlockContextFromQuery(tmpCtx, k.clientContext))
		defer aspectCtx.Destroy()

		// update the message with the new gas value
		msg.GasLimit = gas
		// pass false to not commit StateDB
		rsp, err = k.ApplyMessageWithConfig(cosmosCtx, aspectCtx, msg, nil, false, cfg, txConfig, isCustomVerification)
		if err != nil {
			if errors.Is(err, core.ErrIntrinsicGas) {
				return true, nil, nil // Special case, raise gas limit
			}
			return true, nil, err // Bail out
		}
		return len(rsp.VmError) > 0, rsp, nil
	}

	// Execute the binary search and hone in on an executable gas limit
	hi, err = txs.BinSearch(lo, hi, executable)
	if err != nil {
		return nil, err
	}

	// Reject the txs as invalid if it still fails at the highest allowance
	if hi == gasCap {
		failed, result, err := executable(hi)
		if err != nil {
			return nil, err
		}

		if failed {
			if result != nil && result.VmError != vm.ErrOutOfGas.Error() {
				if result.VmError == vm.ErrExecutionReverted.Error() {
					return nil, types.NewExecErrorWithReason(result.Ret)
				}
				return nil, errors.New(result.VmError)
			}
			// Otherwise, the specified gas cap is too low
			return nil, fmt.Errorf("gas required exceeds allowance (%d)", gasCap)
		}
	}
	return &txs.EstimateGasResponse{Gas: hi}, nil
}

// TraceTx configures a new tracer according to the provided configuration, and
// executes the given message in the provided environment. The return value will
// be tracer dependent.
func (k Keeper) TraceTx(c context.Context, req *txs.QueryTraceTxRequest) (*txs.QueryTraceTxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.TraceConfig != nil && req.TraceConfig.Limit < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "output limit cannot be negative, got %d", req.TraceConfig.Limit)
	}

	// minus one to get the context of block beginning
	contextHeight := req.BlockNumber - 1
	if contextHeight < 1 {
		// 0 is a special value in `ContextWithHeight`
		contextHeight = 1
	}

	ctx := cosmos.UnwrapSDKContext(c)
	ctx = ctx.WithBlockHeight(contextHeight)
	ctx = ctx.WithBlockTime(req.BlockTime)
	ctx = ctx.WithHeaderHash(common.Hex2Bytes(req.BlockHash))
	chainID, err := getChainID(ctx, req.ChainId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	cfg, err := k.EVMConfig(ctx, GetProposerAddress(ctx, req.ProposerAddress), chainID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load evm config: %s", err.Error())
	}
	signer := ethereum.MakeSigner(cfg.ChainConfig, big.NewInt(ctx.BlockHeight()), uint64(ctx.BlockTime().Unix()))

	txConfig := states.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash().Bytes()))
	for i, tx := range req.Predecessors {
		ethTx := tx.AsTransaction()

		// Aspect Runtime Context Lifecycle: create aspect context.
		// This marks the beginning of running an aspect of TraceTx, creating the aspect context,
		// and establishing the link with the SDK context.
		ctx, aspectCtx := k.WithAspectContext(ctx, ethTx, cfg,
			artelatypes.NewEthBlockContextFromQuery(ctx, k.clientContext))

		msg, err := txs.ToMessage(ethTx, signer, cfg.BaseFee)
		if err != nil {
			aspectCtx.Destroy()
			continue
		}
		txConfig.TxHash = ethTx.Hash()
		txConfig.TxIndex = uint(i)

		isCustomVerification := k.isCustomizedVerification(ethTx)
		rsp, err := k.ApplyMessageWithConfig(ctx, aspectCtx, msg, txs.NewNoOpTracer(), true, cfg, txConfig, isCustomVerification)
		if err != nil {
			aspectCtx.Destroy()
			continue
		}

		aspectCtx.Destroy()
		txConfig.LogIndex += uint(len(rsp.Logs))
	}

	tx := req.Msg.AsTransaction()
	txConfig.TxHash = tx.Hash()
	if len(req.Predecessors) > 0 {
		txConfig.TxIndex++
	}

	var tracerConfig json.RawMessage
	if req.TraceConfig != nil && req.TraceConfig.TracerJsonConfig != "" {
		// ignore error. default to no traceConfig
		_ = json.Unmarshal([]byte(req.TraceConfig.TracerJsonConfig), &tracerConfig)
	}

	result, _, err := k.traceTx(ctx, cfg, txConfig, signer, tx, req.TraceConfig, false, tracerConfig)
	if err != nil {
		// error will be returned with detail status from traceTx
		return nil, err
	}

	resultData, err := json.Marshal(result)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &txs.QueryTraceTxResponse{
		Data: resultData,
	}, nil
}

// TraceBlock configures a new tracer according to the provided configuration, and
// executes the given message in the provided environment for all the transactions in the queried block.
// The return value will be tracer dependent.
func (k Keeper) TraceBlock(c context.Context, req *txs.QueryTraceBlockRequest) (*txs.QueryTraceBlockResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.TraceConfig != nil && req.TraceConfig.Limit < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "output limit cannot be negative, got %d", req.TraceConfig.Limit)
	}

	// minus one to get the context of block beginning
	contextHeight := req.BlockNumber - 1
	if contextHeight < 1 {
		// 0 is a special value in `ContextWithHeight`
		contextHeight = 1
	}

	ctx := cosmos.UnwrapSDKContext(c)
	ctx = ctx.WithBlockHeight(contextHeight)
	ctx = ctx.WithBlockTime(req.BlockTime)
	ctx = ctx.WithHeaderHash(common.Hex2Bytes(req.BlockHash))
	chainID, err := getChainID(ctx, req.ChainId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	cfg, err := k.EVMConfig(ctx, GetProposerAddress(ctx, req.ProposerAddress), chainID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to load evm config")
	}
	signer := ethereum.MakeSigner(cfg.ChainConfig, big.NewInt(ctx.BlockHeight()), uint64(ctx.BlockTime().Unix()))
	txsLength := len(req.Txs)
	results := make([]*txs.TxTraceResult, 0, txsLength)

	txConfig := states.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash().Bytes()))
	for i, tx := range req.Txs {
		result := txs.TxTraceResult{}
		ethTx := tx.AsTransaction()
		txConfig.TxHash = ethTx.Hash()
		txConfig.TxIndex = uint(i)
		traceResult, logIndex, err := k.traceTx(ctx, cfg, txConfig, signer, ethTx, req.TraceConfig, true, nil)
		if err != nil {
			result.Error = err.Error()
		} else {
			txConfig.LogIndex = logIndex
			result.Result = traceResult
		}
		results = append(results, &result)
	}

	resultData, err := json.Marshal(results)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &txs.QueryTraceBlockResponse{
		Data: resultData,
	}, nil
}

// traceTx do trace on one txs, it returns a tuple: (traceResult, nextLogIndex, error).
func (k *Keeper) traceTx(
	ctx cosmos.Context,
	cfg *states.EVMConfig,
	txConfig states.TxConfig,
	signer ethereum.Signer,
	tx *ethereum.Transaction,
	traceConfig *support.TraceConfig,
	commitMessage bool,
	tracerJSONConfig json.RawMessage,
) (*interface{}, uint, error) {
	// Assemble the structured logger or the JavaScript tracer
	var (
		tracer    tracers.Tracer
		overrides *ethparams.ChainConfig
		err       error
		timeout   = defaultTraceTimeout
	)

	// Aspect Runtime Context Lifecycle: create aspect context.
	// This marks the beginning of running an aspect of TraceBlock or TraceTx, creating the aspect context,
	// and establishing the link with the SDK context.
	cacheCtx, commit := ctx.CacheContext()
	ctx, aspectCtx := k.WithAspectContext(cacheCtx, tx, cfg,
		artelatypes.NewEthBlockContextFromQuery(ctx, k.clientContext))
	defer func() {
		if commitMessage {
			commit()
		}
		aspectCtx.Destroy()
	}()

	msg, err := txs.ToMessage(tx, signer, cfg.BaseFee)
	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	if traceConfig == nil {
		traceConfig = &support.TraceConfig{}
	}

	if traceConfig.Overrides != nil {
		overrides = traceConfig.Overrides.EthereumConfig(ctx.BlockHeight(), cfg.ChainConfig.ChainID)
	}

	logConfig := logger.Config{
		EnableMemory:     traceConfig.EnableMemory,
		DisableStorage:   traceConfig.DisableStorage,
		DisableStack:     traceConfig.DisableStack,
		EnableReturnData: traceConfig.EnableReturnData,
		Debug:            traceConfig.Debug,
		Limit:            int(traceConfig.Limit),
		Overrides:        overrides,
	}

	tracer = logger.NewStructLogger(&logConfig)

	tCtx := &tracers.Context{
		BlockHash: txConfig.BlockHash,
		TxIndex:   int(txConfig.TxIndex),
		TxHash:    txConfig.TxHash,
	}

	if traceConfig.Tracer != "" {
		if tracer, err = tracers.DefaultDirectory.New(traceConfig.Tracer, tCtx, tracerJSONConfig); err != nil {
			return nil, 0, status.Error(codes.Internal, err.Error())
		}
	}

	// Define a meaningful timeout of a single txs trace
	if traceConfig.Timeout != "" {
		if timeout, err = time.ParseDuration(traceConfig.Timeout); err != nil {
			return nil, 0, status.Errorf(codes.InvalidArgument, "timeout value: %s", err.Error())
		}
	}

	// Handle timeouts and RPC cancellations
	deadlineCtx, cancel := context.WithTimeout(ctx.Context(), timeout)
	defer cancel()

	go func() {
		<-deadlineCtx.Done()
		if errors.Is(deadlineCtx.Err(), context.DeadlineExceeded) {
			tracer.Stop(errors.New("execution timeout"))
		}
	}()

	isCustomVerification := k.isCustomizedVerification(tx)
	res, err := k.ApplyMessageWithConfig(ctx, aspectCtx, msg, tracer, commitMessage, cfg, txConfig, isCustomVerification)
	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}
	var result interface{}
	result, err = tracer.GetResult()
	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	return &result, txConfig.LogIndex + uint(len(res.Logs)), nil
}

// BaseFee implements the Query/BaseFee gRPC method
func (k Keeper) BaseFee(c context.Context, _ *txs.QueryBaseFeeRequest) (*txs.QueryBaseFeeResponse, error) {
	ctx := cosmos.UnwrapSDKContext(c)

	params := k.GetParams(ctx)
	ethCfg := params.ChainConfig.EthereumConfig(ctx.BlockHeight(), k.eip155ChainID)
	baseFee := k.GetBaseFee(ctx, ethCfg)

	res := &txs.QueryBaseFeeResponse{}
	if baseFee != nil {
		aux := math.NewIntFromBigInt(baseFee)
		res.BaseFee = &aux
	}

	return res, nil
}

func (k Keeper) GetSender(c context.Context, in *txs.MsgEthereumTx) (*txs.GetSenderResponse, error) {
	ctx := cosmos.UnwrapSDKContext(c)

	evmConfig, err := k.EVMConfigFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	// Aspect Runtime Context Lifecycle: create aspect context.
	// This marks the beginning of running an aspect of EthCall, creating the aspect context,
	// and establishing the link with the SDK context.
	ctx, aspectCtx := k.WithAspectContext(ctx, in.AsEthCallTransaction(),
		evmConfig, artelatypes.NewEthBlockContextFromHeight(ctx.BlockHeight()))
	defer aspectCtx.Destroy()

	tx := in.AsTransaction()
	sender, _, err := k.tryAspectVerifier(ctx, tx)
	if err != nil {
		return nil, err
	}
	return &txs.GetSenderResponse{Sender: sender.String()}, nil
}

// WithAspectContext creates the Aspect Context and establishes the link to the SDK context.
func (k Keeper) WithAspectContext(ctx cosmos.Context, tx *ethereum.Transaction,
	evmConf *states.EVMConfig, block *artelatypes.EthBlockContext) (cosmos.Context, *artelatypes.AspectRuntimeContext) {
	ethTxContext := artelatypes.NewEthTxContext(tx)
	ethTxContext.WithEVMConfig(evmConf)

	aspectCtx := artelatypes.NewAspectRuntimeContext()
	protocol := provider.NewAspectProtocolProvider(aspectCtx.EthTxContext)
	jitManager := inherent.NewManager(protocol)

	aspectCtx.SetEthTxContext(ethTxContext, jitManager)
	aspectCtx.WithCosmosContext(ctx)
	aspectCtx.SetEthBlockContext(block)
	aspectCtx.CreateStateObject()
	return ctx.WithValue(artelatypes.AspectContextKey, aspectCtx), aspectCtx
}

// getChainID parse chainID from current context if not provided
func getChainID(ctx cosmos.Context, chainID int64) (*big.Int, error) {
	if chainID == 0 {
		return artela.ParseChainID(ctx.ChainID())
	}
	return big.NewInt(chainID), nil
}
