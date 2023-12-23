package contract

import (
	"math/big"

	"github.com/artela-network/artela-evm/vm"
	"github.com/artela-network/artela/x/evm/artela/types"
	evmtypes "github.com/artela-network/artela/x/evm/txs"
	"github.com/artela-network/aspect-core/djpm/contract"
	"github.com/artela-network/aspect-core/djpm/run"
	artelasdkType "github.com/artela-network/aspect-core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
)

func (anc *AspectNativeContract) entrypoint(ctx sdk.Context, msg *core.Message, method *abi.Method, aspectId common.Address, data []byte, commit bool) (*evmtypes.MsgEthereumTxResponse, error) {
	lastHeight := ctx.BlockHeight()
	code, version := anc.aspectService.GetAspectCode(ctx, aspectId, nil, commit)

	aspectCtx, ok := ctx.Value(types.AspectContextKey).(*types.AspectRuntimeContext)
	if !ok {
		return nil, errors.New("AspectNativeContract.entrypoint: unwrap AspectRuntimeContext failed")
	}
	runner, newErr := run.NewRunner(aspectCtx, aspectId.String(), version.Uint64(), code)
	if newErr != nil {
		return nil, newErr
	}
	defer runner.Return()

	var txHash []byte
	ethTxCtx := aspectCtx.EthTxContext()
	if ethTxCtx != nil {
		ethTxContent := ethTxCtx.TxContent()
		if ethTxContent != nil {
			txHash = ethTxContent.Hash().Bytes()
		}
	}

	// ignore gas output for now, since we haven't implemented gas metering for aspect for now
	ret, _, err := runner.JoinPoint(artelasdkType.OPERATION_METHOD, msg.GasLimit, lastHeight, msg.To, &artelasdkType.OperationInput{
		Tx: &artelasdkType.WithFromTxInput{
			Hash: txHash,
			To:   msg.To.Bytes(),
			From: msg.From.Bytes(),
		},
		Block:    &artelasdkType.BlockInput{Number: uint64(lastHeight)},
		CallData: data,
	})
	var vmError string
	var retByte []byte
	if err != nil {
		vmError = err.Error()
	} else {
		retByte, err = method.Outputs.Pack(ret)
		if err != nil {
			vmError = err.Error()
		}
	}
	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: ctx.GasMeter().GasConsumed(),
		VmError: vmError,
		Ret:     retByte,
		Logs:    nil,
	}, nil
}

func (anc *AspectNativeContract) contractsOf(ctx sdk.Context, method *abi.Method, aspectId common.Address, commit bool) (*evmtypes.MsgEthereumTxResponse, error) {
	value, err := anc.aspectService.GetAspectOf(ctx, aspectId, commit)
	if err != nil {
		return nil, err
	}
	addressAry := make([]common.Address, 0)
	for _, data := range value.Values() {
		contractAddr := common.HexToAddress(data.(string))
		addressAry = append(addressAry, contractAddr)
	}

	ret, err := method.Outputs.Pack(addressAry)
	if err != nil {
		return nil, err
	}
	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: ctx.GasMeter().GasConsumed(),
		VmError: "",
		Ret:     ret,
		Logs:    nil,
	}, nil
}

func (k *AspectNativeContract) bind(ctx sdk.Context, aspectId common.Address, account common.Address, aspectVersion *uint256.Int, priority int8, isContract bool, commit bool) (*evmtypes.MsgEthereumTxResponse, error) {
	// check aspect types
	aspectJP := k.aspectService.aspectStore.GetAspectJP(ctx, aspectId, aspectVersion)
	txAspect := artelasdkType.CheckIsTransactionLevel(aspectJP.Int64())
	txVerifier := artelasdkType.CheckIsTxVerifier(aspectJP.Int64())

	if !(txAspect || txVerifier) {
		return nil, errors.New("check aspect join point fail, An aspect which can be bound, must be a transactional or Verifier")
	}
	// EoA can only bind with tx verifier
	if !txVerifier && !isContract {
		return nil, errors.New("unable to bind non-tx-verifier Aspect to an EoA account")
	}

	// bind tx processing aspect if account is a contract
	if txAspect && isContract {
		if err := k.aspectService.aspectStore.BindTxAspect(ctx, account, aspectId, aspectVersion, priority); err != nil {
			return nil, err
		}
	}

	// bind tx verifier aspect
	if txVerifier {
		if err := k.aspectService.aspectStore.BindVerificationAspect(ctx, account, aspectId, aspectVersion, priority, isContract); err != nil {
			return nil, err
		}
	}

	// save reverse index
	if err := k.aspectService.aspectStore.StoreAspectRefValue(ctx, account, aspectId); err != nil {
		return nil, err
	}

	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: ctx.GasMeter().GasConsumed(),
		VmError: "",
		Ret:     nil,
		Logs:    nil,
	}, nil
}

func (k *AspectNativeContract) unbind(ctx sdk.Context, aspectId common.Address, contract common.Address, isContract bool) (*evmtypes.MsgEthereumTxResponse, error) {
	// contract=>aspect object
	// aspectId= [contract,contract]
	if !isContract {
		if err := k.aspectService.aspectStore.UnBindVerificationAspect(ctx, contract, aspectId); err != nil {
			return nil, err
		}
	}
	if err := k.aspectService.aspectStore.UnBindContractAspects(ctx, contract, aspectId); err != nil {
		return nil, err
	}
	if err := k.aspectService.aspectStore.UnbindAspectRefValue(ctx, contract, aspectId); err != nil {
		return nil, err
	}
	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: ctx.GasMeter().GasConsumed(),
		VmError: "",
		Ret:     nil,
		Logs:    nil,
	}, nil
}

func (k *AspectNativeContract) changeVersion(ctx sdk.Context, aspectId common.Address, contract common.Address, version uint64) (*evmtypes.MsgEthereumTxResponse, error) {
	err := k.aspectService.aspectStore.ChangeBoundAspectVersion(ctx, contract, aspectId, version)
	if err != nil {
		return nil, err
	}

	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: 100,
		VmError: "",
		Ret:     nil,
		Logs:    nil,
	}, nil
}

func (k *AspectNativeContract) version(ctx sdk.Context, method *abi.Method, aspectId common.Address) (*evmtypes.MsgEthereumTxResponse, error) {
	version := k.aspectService.aspectStore.GetAspectLastVersion(ctx, aspectId)

	ret, err := method.Outputs.Pack(version.Uint64())
	if err != nil {
		return nil, err
	}

	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: 100,
		VmError: "",
		Ret:     ret,
		Logs:    nil,
	}, nil
}

func (k *AspectNativeContract) aspectsOf(ctx sdk.Context, method *abi.Method, contract common.Address, isContract bool, commit bool) (*evmtypes.MsgEthereumTxResponse, error) {

	aspectInfo := make([]types.AspectInfo, 0)
	deduplicationMap := make(map[string]uint64, 0)
	if isContract {
		aspects, err := k.aspectService.GetBoundAspectForAddr(ctx.BlockHeight(), contract)
		if err != nil {
			return nil, err
		}
		aspectAccounts, verErr := k.aspectService.GetAccountVerifiers(ctx, contract, commit)
		if verErr != nil {
			return nil, verErr
		}

		for _, aspect := range aspects {
			if _, exist := deduplicationMap[aspect.AspectId]; exist {
				continue
			}
			deduplicationMap[aspect.AspectId] = aspect.Version
			info := types.AspectInfo{
				AspectId: common.HexToAddress(aspect.AspectId),
				Version:  aspect.Version,
				Priority: int8(aspect.Priority),
			}
			aspectInfo = append(aspectInfo, info)
		}

		for _, aspect := range aspectAccounts {
			if _, exist := deduplicationMap[aspect.AspectId]; exist {
				continue
			}
			deduplicationMap[aspect.AspectId] = aspect.Version
			info := types.AspectInfo{
				AspectId: common.HexToAddress(aspect.AspectId),
				Version:  aspect.Version,
				Priority: int8(aspect.Priority),
			}
			aspectInfo = append(aspectInfo, info)
		}

	} else {
		aspectAccounts, verErr := k.aspectService.GetAccountVerifiers(ctx, contract, commit)
		if verErr != nil {
			return nil, verErr
		}
		for _, aspect := range aspectAccounts {
			if _, exist := deduplicationMap[aspect.AspectId]; exist {
				continue
			}
			deduplicationMap[aspect.AspectId] = aspect.Version
			info := types.AspectInfo{
				AspectId: common.HexToAddress(aspect.AspectId),
				Version:  aspect.Version,
				Priority: int8(aspect.Priority),
			}
			aspectInfo = append(aspectInfo, info)
		}
	}

	ret, pkErr := method.Outputs.Pack(aspectInfo)
	if pkErr != nil {
		return nil, pkErr
	}
	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: ctx.GasMeter().GasConsumed(),
		VmError: "",
		Ret:     ret,
		Logs:    nil,
	}, nil
}

func (k *AspectNativeContract) checkContractOwner(ctx sdk.Context, to *common.Address, nonce uint64, sender common.Address) bool {
	msg, err := contract.ArtelaOwnerMsg(to, nonce, sender)
	if err != nil {
		return false
	}
	fromAccount := vm.AccountRef(msg.From)
	k.evm.CloseAspectCall()
	defer k.evm.AspectCall()

	aspectCtx, ok := ctx.Value(types.AspectContextKey).(*types.AspectRuntimeContext)
	if !ok {
		return false
	}
	ret, _, err := k.evm.Call(aspectCtx, fromAccount, *msg.To, msg.Data, msg.GasLimit, msg.Value)
	if err != nil {
		return false
	}
	result, err := contract.UnpackIsOwnerResult(ret)
	if err != nil {
		return false
	}
	return result
}

func (k *AspectNativeContract) checkAspectOwner(ctx sdk.Context, aspectId common.Address, sender common.Address, commit bool) (bool, error) {
	bHeight := ctx.BlockHeight()
	code, ver := k.aspectService.GetAspectCode(ctx, aspectId, nil, commit)
	if code == nil {
		return false, nil
	}

	aspectCtx, ok := ctx.Value(types.AspectContextKey).(*types.AspectRuntimeContext)
	if !ok {
		return false, errors.New("checkAspectOwner: unwrap AspectRuntimeContext failed")
	}
	runner, newErr := run.NewRunner(aspectCtx, aspectId.String(), ver.Uint64(), code)
	if newErr != nil {
		return false, newErr
	}
	defer runner.Return()

	binding, runErr := runner.IsOwner(bHeight, 0, &sender, sender.Bytes())
	return binding, runErr
}

func (k *AspectNativeContract) deploy(ctx sdk.Context, aspectId common.Address, code []byte, properties []types.Property, joinPoint *big.Int) (*evmtypes.MsgEthereumTxResponse, error) {
	if len(code) == 0 && len(properties) == 0 && joinPoint == nil {
		return &evmtypes.MsgEthereumTxResponse{
			GasUsed: ctx.GasMeter().GasConsumed(),
			VmError: "",
			Ret:     nil,
			Logs:    nil,
		}, nil
	}

	aspectVersion := k.aspectService.aspectStore.StoreAspectCode(ctx, aspectId, code)

	err := k.aspectService.aspectStore.StoreAspectJP(ctx, aspectId, aspectVersion, joinPoint)
	if err != nil {
		return nil, err
	}
	err = k.aspectService.aspectStore.StoreAspectProperty(ctx, aspectId, properties)
	if err != nil {
		return nil, err
	}

	level := artelasdkType.CheckIsBlockLevel(joinPoint.Int64())
	if level {
		storeErr := k.aspectService.aspectStore.StoreBlockLevelAspect(ctx, aspectId)
		if storeErr != nil {
			return nil, storeErr
		}
	} else {
		// for k.update,one aspect trans to tx level
		removeErr := k.aspectService.aspectStore.RemoveBlockLevelAspect(ctx, aspectId)
		if removeErr != nil {
			return nil, removeErr
		}
	}

	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: ctx.GasMeter().GasConsumed(),
		VmError: "",
		Ret:     nil,
		Logs:    nil,
	}, nil
}
