package keeper

import (
	"errors"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	cometbft "github.com/cometbft/cometbft/types"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethereum "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"

	artcore "github.com/artela-network/artela-evm/core"
	"github.com/artela-network/artela-evm/vm"
	"github.com/artela-network/artela/common/aspect"
	artela "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/artela/contract"
	artelatypes "github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/artela/x/evm/states"
	"github.com/artela-network/artela/x/evm/txs"
	"github.com/artela-network/artela/x/evm/txs/support"
	"github.com/artela-network/artela/x/evm/types"
	"github.com/artela-network/aspect-core/djpm"
	asptypes "github.com/artela-network/aspect-core/types"
)

// NewEVM generates a go-ethereum VM from the provided Message fields and the chain parameters
// (ChainConfig and module Params). It additionally sets the validator operator address as the
// coinbase address to make it available for the COINBASE opcode, even though there is no
// beneficiary of the coinbase txs (since we're not mining).
func (k *Keeper) NewEVM(
	ctx cosmos.Context,
	msg *core.Message,
	cfg *states.EVMConfig,
	tracer vm.EVMLogger,
	stateDB vm.StateDB,
) *vm.EVM {
	blockCtx := vm.BlockContext{
		CanTransfer: artcore.CanTransfer,
		Transfer:    artcore.Transfer,
		GetHash:     k.GetHashFn(ctx),
		Coinbase:    cfg.CoinBase,
		GasLimit:    artela.BlockGasLimit(ctx),
		BlockNumber: big.NewInt(ctx.BlockHeight()),
		Time:        uint64(ctx.BlockHeader().Time.Unix()),
		Difficulty:  big.NewInt(0), // unused. Only required in PoW context
		BaseFee:     cfg.BaseFee,
		Random:      nil, // not supported
	}

	txCtx := artcore.NewEVMTxContext(msg)
	if tracer == nil {
		tracer = k.Tracer(ctx, msg, cfg.ChainConfig)
	}
	vmConfig := k.VMConfig(ctx, msg, cfg, tracer)
	return vm.NewEVM(blockCtx, txCtx, stateDB, cfg.ChainConfig, vmConfig)
}

// GetHashFn implements vm.GetHashFunc for Artela.
// It handles 3 cases:
//  1. The requested height matches the current height from context (and thus same epoch number)
//  2. The requested height is from a previous height from the same chain epoch
//  3. The requested height is from a height greater than the latest one
func (k Keeper) GetHashFn(ctx cosmos.Context) vm.GetHashFunc {
	return func(height uint64) common.Hash {
		h, err := artela.SafeInt64(height)
		if err != nil {
			k.Logger(ctx).Error("failed to cast height to int64", "error", err)
			return common.Hash{}
		}

		switch {
		case ctx.BlockHeight() == h:
			// Case 1: The requested height matches the one from the context so we can retrieve the header
			// hash directly from the context.
			// Note: The headerHash is only set at begin block, it will be nil in case of a query context
			headerHash := ctx.HeaderHash()
			if len(headerHash) != 0 {
				return common.BytesToHash(headerHash)
			}

			// only recompute the hash if not set (eg: checkTxState)
			contextBlockHeader := ctx.BlockHeader()
			header, err := cometbft.HeaderFromProto(&contextBlockHeader)
			if err != nil {
				k.Logger(ctx).Error("failed to cast tendermint header from proto", "error", err)
				return common.Hash{}
			}

			headerHash = header.Hash()
			return common.BytesToHash(headerHash)

		case ctx.BlockHeight() > h:
			// Case 2: if the chain is not the current height we need to retrieve the hash from the store for the
			// current chain epoch. This only applies if the current height is greater than the requested height.
			histInfo, found := k.stakingKeeper.GetHistoricalInfo(ctx, h)
			if !found {
				k.Logger(ctx).Debug("historical info not found", "height", h)
				return common.Hash{}
			}

			header, err := cometbft.HeaderFromProto(&histInfo.Header)
			if err != nil {
				k.Logger(ctx).Error("failed to cast tendermint header from proto", "error", err)
				return common.Hash{}
			}

			return common.BytesToHash(header.Hash())
		default:
			// Case 3: heights greater than the current one returns an empty hash.
			return common.Hash{}
		}
	}
}

// ApplyTransaction runs and attempts to perform a states transition with the given txs (i.e Message), that will
// only be persisted (committed) to the underlying KVStore if the txs does not fail.
//
// # Gas tracking
//
// Ethereum consumes gas according to the EVM opcodes instead of general reads and writes to store. Because of this, the
// states transition needs to ignore the SDK gas consumption mechanism defined by the GasKVStore and instead consume the
// amount of gas used by the VM execution. The amount of gas used is tracked by the EVM and returned in the execution
// result.
//
// Prior to the execution, the starting txs gas meter is saved and replaced with an infinite gas meter in a new context
// in order to ignore the SDK gas consumption config values (read, write, has, delete).
// After the execution, the gas used from the message execution will be added to the starting gas consumed, taking into
// consideration the amount of gas returned. Finally, the context is updated with the EVM gas consumed value prior to
// returning.
//
// For relevant discussion see: https://github.com/cosmos/cosmos-sdk/discussions/9072
func (k *Keeper) ApplyTransaction(ctx cosmos.Context, tx *ethereum.Transaction) (*txs.MsgEthereumTxResponse, error) {
	var (
		bloom        *big.Int
		bloomReceipt ethereum.Bloom
	)

	// build evm config and txs config
	evmConfig, err := k.EVMConfig(ctx, ctx.BlockHeader().ProposerAddress, k.eip155ChainID)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to load evm config")
	}
	txConfig := k.TxConfig(ctx, tx.Hash(), tx.Type())

	// retrieve aspectCtx from sdk.Context
	aspectCtx, ok := ctx.Value(artelatypes.AspectContextKey).(*artelatypes.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("ApplyMessageWithConfig: unwrap AspectRuntimeContext failed")
	}

	// snapshot to contain the txs processing and post processing in same scope
	// var commit func()
	// tmpCtx := ctx
	tmpCtx, commit := ctx.CacheContext()

	// use the temp ctx for later tx processing
	aspectCtx.WithCosmosContext(tmpCtx)

	// get the signer according to the chain rules from the config and block height
	signer := k.MakeSigner(ctx, tx, evmConfig.ChainConfig, big.NewInt(ctx.BlockHeight()), uint64(ctx.BlockTime().Unix()))

	msg, err := txs.ToMessage(tx, signer, evmConfig.BaseFee)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to return ethereum txs as core message")
	}

	msg.Data, err = k.processMsgData(tx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "unable to process msg data")
	}

	// pass true to commit the StateDB
	res, err := k.ApplyMessageWithConfig(tmpCtx, aspectCtx, msg, nil, true, evmConfig, txConfig, k.isCustomizedVerification(tx))
	if err != nil {
		ctx.Logger().Error("ApplyMessageWithConfig with error", "txhash", tx.Hash().String(), "error", err, "response", res)
		return nil, errorsmod.Wrap(err, "failed to apply ethereum core message")
	}
	ctx.Logger().Debug("ApplyMessageWithConfig", "txhash", tx.Hash().String(), "response", res)

	logs := support.LogsToEthereum(res.Logs)

	// Compute block bloom filter
	if len(logs) > 0 {
		bloom = k.GetBlockBloomTransient(ctx)
		bloom.Or(bloom, big.NewInt(0).SetBytes(ethereum.LogsBloom(logs)))
		bloomReceipt = ethereum.BytesToBloom(bloom.Bytes())
	}

	cumulativeGasUsed := res.GasUsed
	if ctx.BlockGasMeter() != nil {
		limit := ctx.BlockGasMeter().Limit()
		cumulativeGasUsed += ctx.BlockGasMeter().GasConsumed()
		if cumulativeGasUsed > limit {
			cumulativeGasUsed = limit
		}
	}
	res.CumulativeGasUsed = cumulativeGasUsed

	var contractAddr common.Address
	if msg.To == nil || aspect.IsAspectDeploy(msg.To, msg.Data) {
		contractAddr = crypto.CreateAddress(msg.From, msg.Nonce)
	}

	receipt := &ethereum.Receipt{
		Type:              tx.Type(),
		PostState:         nil, // TODO: intermediate states root
		CumulativeGasUsed: cumulativeGasUsed,
		Bloom:             bloomReceipt,
		Logs:              logs,
		TxHash:            txConfig.TxHash,
		ContractAddress:   contractAddr,
		GasUsed:           res.GasUsed,
		BlockHash:         txConfig.BlockHash,
		BlockNumber:       big.NewInt(ctx.BlockHeight()),
		TransactionIndex:  txConfig.TxIndex,
	}

	if !res.Failed() {
		receipt.Status = ethereum.ReceiptStatusSuccessful

		if commit != nil {
			commit()
			res.Logs = support.NewLogsFromEth(receipt.Logs)
			ctx.EventManager().EmitEvents(tmpCtx.EventManager().Events())
		}
	} else {
		ctx.Logger().Debug("apply transaction failed", "tx hash", txConfig.TxHash.String(), "error", res.VmError)
	}

	// refund gas in order to match the Ethereum gas consumption instead of the default SDK one.
	if err = k.RefundGas(ctx, msg, msg.GasLimit-res.GasUsed, evmConfig.Params.EvmDenom); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to refund gas leftover gas to sender %s", msg.From)
	}

	if len(receipt.Logs) > 0 {
		// Update transient block bloom filter
		k.SetBlockBloomTransient(ctx, receipt.Bloom.Big())
		k.SetLogSizeTransient(ctx, uint64(txConfig.LogIndex)+uint64(len(receipt.Logs)))
	}

	k.SetTxIndexTransient(ctx, uint64(txConfig.TxIndex)+1)

	totalGasUsed, err := k.AddTransientGasUsed(ctx, res.GasUsed)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to add transient gas used")
	}

	// reset the gas meter for current cosmos txs
	k.ResetGasMeterAndConsumeGas(ctx, totalGasUsed)

	return res, nil
}

// ApplyMessage calls ApplyMessageWithConfig with an empty TxConfig.
func (k *Keeper) ApplyMessage(ctx cosmos.Context, msg *core.Message, tracer vm.EVMLogger, commit bool) (*txs.MsgEthereumTxResponse, error) {
	evmConfig, err := k.EVMConfig(ctx, ctx.BlockHeader().ProposerAddress, k.eip155ChainID)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to load evm config")
	}

	txConfig := states.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash()))

	// retrieve aspectCtx from sdk.Context
	aspectCtx, ok := ctx.Value(artelatypes.AspectContextKey).(*artelatypes.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("ApplyMessageWithConfig: unwrap AspectRuntimeContext failed")
	}

	return k.ApplyMessageWithConfig(ctx, aspectCtx, msg, tracer, commit, evmConfig, txConfig, false)
}

// ApplyMessageWithConfig computes the new states by applying the given message against the existing states.
// If the message fails, the VM execution error with the reason will be returned to the client
// and the txs won't be committed to the store.
//
// # Reverted states
//
// The snapshot and rollback are supported by the `states.StateDB`.
//
// # Different Callers
//
// It's called in three scenarios:
// 1. `ApplyTransaction`, in the txs processing flow.
// 2. `EthCall/EthEstimateGas` grpc query handler.
// 3. Called by other native modules directly.
//
// # PreChecks and Preprocessing
//
// All relevant states transition preChecks for the MsgEthereumTx are performed on the AnteHandler,
// prior to running the txs against the states. The PreChecks run are the following:
//
// 1. the nonce of the message caller is correct
// 2. caller has enough balance to cover txs fee(gasLimit * gasPrice)
// 3. the amount of gas required is available in the block
// 4. the purchased gas is enough to cover intrinsic usage
// 5. there is no overflow when calculating intrinsic gas
// 6. caller has enough balance to cover asset transfer for **topmost** call
//
// The preprocessing steps performed by the AnteHandler are:
//
// 1. set up the initial access list (if fork > Berlin)
//
// # Tracer parameter
//
// It should be a `vm.Tracer` object or nil, if pass `nil`, it'll create a default one based on keeper options.
//
// # Commit parameter
//
// If commit is true, the `StateDB` will be committed, otherwise discarded.
func (k *Keeper) ApplyMessageWithConfig(ctx cosmos.Context,
	aspectCtx *artelatypes.AspectRuntimeContext,
	msg *core.Message,
	tracer vm.EVMLogger,
	commit bool,
	cfg *states.EVMConfig,
	txConfig states.TxConfig,
	isCustomVerification bool,
) (*txs.MsgEthereumTxResponse, error) {
	var (
		ret   []byte // return bytes from evm execution
		vmErr error  // vm errors do not effect consensus and are therefore not assigned to err
	)

	// return error if contract creation or call are disabled through governance
	if !cfg.Params.EnableCreate && msg.To == nil {
		return nil, errorsmod.Wrap(types.ErrCreateDisabled, "failed to create new contract")
	} else if !cfg.Params.EnableCall && msg.To != nil {
		return nil, errorsmod.Wrap(types.ErrCallDisabled, "failed to call contract")
	}

	stateDB := states.New(ctx, k, txConfig)
	evm := k.NewEVM(ctx, msg, cfg, tracer, stateDB)

	// Aspect Runtime Context Lifecycle: set EVM params.
	// Before the pre-transaction execution, establish the EVM context, encompassing details such as
	// the sender of the transaction (from), transaction message, the EVM responsible for executing
	// the transaction, EVM tracer, and the stateDB associated with running the EVM transaction.
	aspectCtx.EthTxContext().WithEVM(msg.From, msg, evm, evm.Tracer(), stateDB)
	aspectCtx.EthTxContext().WithTxIndex(uint64(txConfig.TxIndex))
	aspectCtx.EthTxContext().WithCommit(commit)

	// Aspect Runtime Context Lifecycle: create state object.
	aspectCtx.CreateStateObject()

	leftoverGas := msg.GasLimit

	// Allow the tracer captures the txs level events, mainly the gas consumption.
	var aspectLogger asptypes.AspectLogger
	if tracer != nil {
		tracer.CaptureTxStart(leftoverGas)
		defer func() {
			tracer.CaptureTxEnd(leftoverGas)
		}()

		aspectLogger, _ = tracer.(asptypes.AspectLogger)
	}

	sender := vm.AccountRef(msg.From)
	contractCreation := msg.To == nil
	isLondon := cfg.ChainConfig.IsLondon(evm.Context.BlockNumber)

	intrinsicGas, err := k.GetEthIntrinsicGas(ctx, msg, cfg.ChainConfig, contractCreation, isCustomVerification)
	if err != nil {
		// should have already been checked on Ante Handler
		return nil, errorsmod.Wrap(err, "intrinsic gas failed")
	}

	// Should check again even if it is checked on Ante Handler, because eth_call don't go through Ante Handler.
	if leftoverGas < intrinsicGas {
		// eth_estimateGas will check for this exact error
		return nil, errorsmod.Wrap(core.ErrIntrinsicGas, "apply message")
	}
	leftoverGas -= intrinsicGas

	// access list preparation is moved from ante handler to here, because it's needed when `ApplyMessage` is called
	// under contexts where ante handlers are not run, for example `eth_call` and `eth_estimateGas`.
	if rules := cfg.ChainConfig.Rules(big.NewInt(ctx.BlockHeight()), cfg.ChainConfig.MergeNetsplitBlock != nil, uint64(ctx.BlockTime().Unix())); rules.IsBerlin {
		stateDB.PrepareAccessList(msg.From, msg.To, vm.ActivePrecompiles(rules), msg.AccessList)
	}
	lastHeight := uint64(ctx.BlockHeight())
	// if transaction is Aspect operational, short the circuit and skip the processes
	if isAspectOpTx := asptypes.IsAspectContractAddr(msg.To); isAspectOpTx {
		nativeContract := contract.NewAspectNativeContract(k.storeKey, evm,
			ctx.BlockHeight, stateDB, k.logger)
		nativeContract.Init()
		ret, leftoverGas, vmErr = nativeContract.ApplyMessage(ctx, msg, leftoverGas, commit)
	} else if contractCreation {
		// take over the nonce management from evm:
		// - reset sender's nonce to msg.Nonce() before calling evm.
		// - increase sender's nonce by one no matter the result.
		stateDB.SetNonce(sender.Address(), msg.Nonce)
		ret, _, leftoverGas, vmErr = evm.Create(aspectCtx, sender, msg.Data, leftoverGas, msg.Value)
		stateDB.SetNonce(sender.Address(), msg.Nonce+1)
	} else {
		// begin pre tx aspect execution
		preTxResult := djpm.AspectInstance().PreTxExecute(aspectCtx, msg.From, *msg.To, msg.Data, ctx.BlockHeight(), leftoverGas, msg.Value, &asptypes.PreTxExecuteInput{
			Tx: &asptypes.WithFromTxInput{
				Hash: aspectCtx.EthTxContext().TxContent().Hash().Bytes(),
				To:   msg.To.Bytes(),
				From: msg.From.Bytes(),
			},
			Block: &asptypes.BlockInput{Number: &lastHeight},
		}, aspectLogger)

		leftoverGas = preTxResult.Gas
		if preTxResult.Err != nil {
			// short circuit if pre tx failed
			if preTxResult.Err.Error() == vm.ErrOutOfGas.Error() {
				vmErr = vm.ErrOutOfGas
			} else {
				vmErr = preTxResult.Err
			}
		} else {
			// execute evm call
			ret, leftoverGas, vmErr = evm.Call(aspectCtx, sender, *msg.To, msg.Data, leftoverGas, msg.Value)
			status := ethereum.ReceiptStatusSuccessful
			if vmErr != nil {
				status = ethereum.ReceiptStatusFailed
			}

			logs := stateDB.Logs()

			// compute block bloom filter
			var bloomReceipt ethereum.Bloom
			if len(logs) > 0 {
				bloom := k.GetBlockBloomTransient(ctx)
				bloom.Or(bloom, big.NewInt(0).SetBytes(ethereum.LogsBloom(logs)))
				bloomReceipt = ethereum.BytesToBloom(bloom.Bytes())
			}

			// compute gas
			gasUsed := msg.GasLimit - leftoverGas
			cumulativeGasUsed := gasUsed
			if ctx.BlockGasMeter() != nil {
				limit := ctx.BlockGasMeter().Limit()
				cumulativeGasUsed += ctx.BlockGasMeter().GasConsumed()
				if cumulativeGasUsed > limit {
					cumulativeGasUsed = limit
				}
			}

			// only trigger post tx execute if tx not reverted
			if vmErr == nil {
				// set receipt
				aspectCtx.EthTxContext().WithReceipt(&ethereum.Receipt{
					Status:            status,
					Bloom:             bloomReceipt,
					Logs:              logs,
					GasUsed:           gasUsed,
					CumulativeGasUsed: cumulativeGasUsed,
				})

				// begin post tx aspect execution
				postTxResult := djpm.AspectInstance().PostTxExecute(aspectCtx, msg.From, *msg.To, msg.Data, ctx.BlockHeight(), leftoverGas, msg.Value,
					&asptypes.PostTxExecuteInput{
						Tx: &asptypes.WithFromTxInput{
							Hash: aspectCtx.EthTxContext().TxContent().Hash().Bytes(),
							To:   msg.To.Bytes(),
							From: msg.From.Bytes(),
						},
						Block:   &asptypes.BlockInput{Number: &lastHeight},
						Receipt: &asptypes.ReceiptInput{Status: &status},
					}, aspectLogger)
				if postTxResult.Err != nil {
					// overwrite vmErr if post tx reverted
					if postTxResult.Err.Error() == vm.ErrOutOfGas.Error() {
						vmErr = vm.ErrOutOfGas
					} else {
						vmErr = postTxResult.Err
					}
					ret = postTxResult.Ret
				}
				leftoverGas = postTxResult.Gas
			}
		}
	}

	refundQuotient := params.RefundQuotient

	// After EIP-3529: refunds are capped to gasUsed / 5
	if isLondon {
		refundQuotient = params.RefundQuotientEIP3529
	}

	// calculate gas refund
	if msg.GasLimit < leftoverGas {
		return nil, errorsmod.Wrap(types.ErrGasOverflow, "apply message")
	}
	// refund gas
	temporaryGasUsed := msg.GasLimit - leftoverGas
	refund := GasToRefund(stateDB.GetRefund(), temporaryGasUsed, refundQuotient)

	// update leftoverGas and temporaryGasUsed with refund amount
	leftoverGas += refund
	temporaryGasUsed -= refund

	// EVM execution error needs to be available for the JSON-RPC client
	var vmError string
	if vmErr != nil {
		vmError = vmErr.Error()
	}

	// The dirty states in `StateDB` is either committed or discarded after return
	if commit {
		if err := stateDB.Commit(); err != nil {
			return nil, errorsmod.Wrap(err, "failed to commit stateDB")
		}
	}

	// calculate a minimum amount of gas to be charged to sender if GasLimit
	// is considerably higher than GasUsed to stay more aligned with Tendermint gas mechanics

	gasLimit := cosmos.NewDec(int64(msg.GasLimit))
	minGasMultiplier := k.GetMinGasMultiplier(ctx)
	minimumGasUsed := gasLimit.Mul(minGasMultiplier)

	if msg.GasLimit < leftoverGas {
		return nil, errorsmod.Wrapf(types.ErrGasOverflow, "message gas limit < leftover gas (%d < %d)", msg.GasLimit, leftoverGas)
	}

	gasUsed := cosmos.MaxDec(minimumGasUsed, cosmos.NewDec(int64(temporaryGasUsed))).TruncateInt().Uint64()
	// reset leftoverGas, to be used by the tracer
	leftoverGas = msg.GasLimit - gasUsed

	return &txs.MsgEthereumTxResponse{
		GasUsed: gasUsed,
		VmError: vmError,
		Ret:     ret,
		Logs:    support.NewLogsFromEth(stateDB.Logs()),
		Hash:    txConfig.TxHash.Hex(),
	}, nil
}
