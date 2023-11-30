package contract

import (
	"github.com/artela-network/artela/x/evm/artela/types"
	evmtypes "github.com/artela-network/artela/x/evm/txs"
	xtype "github.com/artela-network/artela/x/evm/types"

	"github.com/artela-network/aspect-core/djpm/contract"
	"github.com/artela-network/aspect-core/djpm/run"
	artelasdkType "github.com/artela-network/aspect-core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
)

func (anc *AspectNativeContract) entrypoint(ctx sdk.Context, tx *ethtypes.Transaction, method *abi.Method, aspectId common.Address, data []byte, commit bool) (*evmtypes.MsgEthereumTxResponse, error) {
	lastHeight := ctx.BlockHeight()
	code, _ := anc.aspectService.GetAspectCode(ctx, aspectId, nil, commit)
	runner, newErr := run.NewRunner(aspectId.String(), code)
	if newErr != nil {
		return nil, newErr
	}
	transaction := &artelasdkType.EthTransaction{
		BlockNumber: ctx.BlockHeight(),
		From:        artelasdkType.ARTELA_ADDR,
		Input:       data,
		To:          aspectId.String(),
	}
	ret, runErr := runner.JoinPoint(artelasdkType.OPERATION_METHOD, tx.Gas(), lastHeight, tx.To(), transaction)
	vmError := ""
	retByte := make([]byte, 0)
	if runErr != nil {
		vmError = runErr.Error()
	} else {
		byteData := &artelasdkType.BytesData{}
		dataUnpackErr := ret.Data.UnmarshalTo(byteData)
		if dataUnpackErr != nil {
			vmError = dataUnpackErr.Error()
		}
		retByte, dataUnpackErr = method.Outputs.Pack(byteData.Data)
		if dataUnpackErr != nil {
			vmError = dataUnpackErr.Error()
		}
	}
	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: ctx.GasMeter().GasConsumed(),
		VmError: vmError,
		Ret:     retByte,
		Logs:    nil,
		Hash:    tx.Hash().Hex(),
	}, nil
}

func (anc *AspectNativeContract) contractsOf(ctx sdk.Context, tx *ethtypes.Transaction, method *abi.Method, aspectId common.Address, commit bool) (*evmtypes.MsgEthereumTxResponse, error) {
	value, err := anc.aspectService.GetAspectOf(ctx, aspectId, commit)
	if err != nil {
		return nil, err
	}
	addressAry := make([]common.Address, 0)
	iterator := value.Iterator()
	if iterator.Next() {
		ct := iterator.Value()
		if ct != nil {
			contractAddr := common.HexToAddress(ct.(string))
			addressAry = append(addressAry, contractAddr)
		}
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
		Hash:    tx.Hash().Hex(),
	}, nil
}

func (k *AspectNativeContract) bind(ctx sdk.Context, tx *ethtypes.Transaction, aspectId common.Address, contract common.Address, aspectVersion *uint256.Int, priority int8, commit bool) (*evmtypes.MsgEthereumTxResponse, error) {
	aspectCode, _ := k.aspectService.GetAspectCode(ctx, aspectId, aspectVersion, commit)
	level, err := k.checkTransactionLevel(aspectId, aspectCode)
	if err != nil || !level {
		return nil, errors.Wrapf(xtype.ErrCallContract, "aspect not implement `IAspectTransaction`, aspectId: %s , err: %s", aspectId, err.Error())
	}
	if err := k.aspectService.aspectStore.BindContractAspects(ctx, k.storeKey, contract, aspectId, aspectVersion, priority); err != nil {
		return nil, err
	}
	if err := k.aspectService.aspectStore.StoreAspectRefValue(ctx, k.storeKey, contract, aspectId); err != nil {
		return nil, err
	}
	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: ctx.GasMeter().GasConsumed(),
		VmError: "",
		Ret:     nil,
		Logs:    nil,
		Hash:    tx.Hash().Hex(),
	}, nil
}

func (k *AspectNativeContract) unbind(ctx sdk.Context, tx *ethtypes.Transaction, aspectId common.Address, contract common.Address) (*evmtypes.MsgEthereumTxResponse, error) {
	if err := k.aspectService.aspectStore.UnBindContractAspects(ctx, k.storeKey, contract, aspectId); err != nil {
		return nil, err
	}
	if err := k.aspectService.aspectStore.UnbindAspectRefValue(ctx, k.storeKey, contract, aspectId); err != nil {
		return nil, err
	}

	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: ctx.GasMeter().GasConsumed(),
		VmError: "",
		Ret:     nil,
		Logs:    nil,
		Hash:    tx.Hash().Hex(),
	}, nil
}

func (k *AspectNativeContract) changeVersion(ctx sdk.Context, tx *ethtypes.Transaction, aspectId common.Address, contract common.Address, version uint64) (*evmtypes.MsgEthereumTxResponse, error) {
	err := k.aspectService.aspectStore.ChangeBoundAspectVersion(ctx, k.storeKey, contract, aspectId, version)
	if err != nil {
		return nil, err
	}

	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: 100,
		VmError: "",
		Ret:     nil,
		Logs:    nil,
		Hash:    tx.Hash().Hex(),
	}, nil
}

func (k *AspectNativeContract) version(ctx sdk.Context, tx *ethtypes.Transaction, method *abi.Method, aspectId common.Address) (*evmtypes.MsgEthereumTxResponse, error) {
	version := k.aspectService.aspectStore.GetAspectLastVersion(ctx, k.storeKey, aspectId)

	ret, err := method.Outputs.Pack(version)
	if err != nil {
		return nil, err
	}

	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: 100,
		VmError: "",
		Ret:     ret,
		Logs:    nil,
		Hash:    tx.Hash().Hex(),
	}, nil
}

func (k *AspectNativeContract) aspectsOf(ctx sdk.Context, tx *ethtypes.Transaction, method *abi.Method, contract common.Address, commit bool) (*evmtypes.MsgEthereumTxResponse, error) {
	aspects, err := k.aspectService.GetAspectForAddr(ctx, contract, commit)
	if err != nil {
		return nil, err
	}
	ret, pkErr := method.Outputs.Pack(aspects)
	if pkErr != nil {
		return nil, pkErr
	}
	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: ctx.GasMeter().GasConsumed(),
		VmError: "",
		Ret:     ret,
		Logs:    nil,
		Hash:    tx.Hash().Hex(),
	}, nil
}

func (k *AspectNativeContract) checkContractOwner(ctx sdk.Context, to *common.Address, nonce uint64, sender common.Address) bool {
	msg, err := contract.ArtelaOwnerMsg(to, nonce, sender)
	if err != nil {
		return false
	}
	message, err := k.applyMessageFunc(ctx, msg, nil, false)
	if err != nil {
		return false
	}
	ret := message.Ret
	result, err := contract.UnpackIsOwnerResult(ret)
	if err != nil {
		return false
	}
	return result
}

func (k *AspectNativeContract) checkContractBinding(ctx sdk.Context, aspectId common.Address, contract common.Address, commit bool) (bool, error) {
	bHeight := ctx.BlockHeight()
	code, _ := k.aspectService.GetAspectCode(ctx, aspectId, nil, commit)
	runner, newErr := run.NewRunner(aspectId.String(), code)
	if newErr != nil {
		return false, newErr
	}
	binding, runErr := runner.OnContractBinding(bHeight, 0, &contract, contract.String())
	return binding, runErr
}

func (k *AspectNativeContract) checkAspectOwner(ctx sdk.Context, aspectId common.Address, sender common.Address, commit bool) (bool, error) {
	bHeight := ctx.BlockHeight()
	code, _ := k.aspectService.GetAspectCode(ctx, aspectId, nil, commit)
	if code == nil {
		return true, nil
	}
	runner, newErr := run.NewRunner(aspectId.String(), code)
	if newErr != nil {
		return false, newErr
	}
	binding, runErr := runner.IsOwner(bHeight, 0, &sender, sender.String())
	return binding, runErr
}

func (k *AspectNativeContract) CheckIsAspectOwnerByCode(ctx sdk.Context, aspectId common.Address, code []byte, sender common.Address) (bool, error) {
	runner, newErr := run.NewRunner(aspectId.String(), code)
	if newErr != nil {
		return false, newErr
	}
	binding, runErr := runner.IsOwner(ctx.BlockHeight(), 0, &sender, sender.String())
	return binding, runErr
}

func (k *AspectNativeContract) deploy(ctx sdk.Context, tx *ethtypes.Transaction, aspectId common.Address, code []byte, properties []types.Property) (*evmtypes.MsgEthereumTxResponse, error) {
	k.aspectService.aspectStore.StoreAspectCode(ctx, k.storeKey, aspectId, code)
	k.aspectService.aspectStore.StoreAspectProperty(ctx, k.storeKey, aspectId, properties)

	level, err := k.checkBlockLevel(aspectId, code)
	if err != nil {
		return nil, err
	}
	if level {
		storeErr := k.aspectService.aspectStore.StoreBlockLevelAspect(ctx, k.storeKey, aspectId)
		if storeErr != nil {
			return nil, storeErr
		}
	} else {
		// for k.update,one aspect trans to tx level
		removeErr := k.aspectService.aspectStore.RemoveBlockLevelAspect(ctx, k.storeKey, aspectId)
		if removeErr != nil {
			return nil, removeErr
		}
	}

	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: ctx.GasMeter().GasConsumed(),
		VmError: "",
		Ret:     nil,
		Logs:    nil,
		Hash:    tx.Hash().Hex(),
	}, nil
}

func (k *AspectNativeContract) checkBlockLevel(aspectId common.Address, code []byte) (bool, error) {
	runner, newErr := run.NewRunner(aspectId.String(), code)
	if newErr != nil {
		return false, newErr
	}
	binding, runErr := runner.IsBlockLevel()
	return binding, runErr
}

func (k *AspectNativeContract) checkTransactionLevel(aspectId common.Address, code []byte) (bool, error) {
	if len(code) == 0 {
		return false, errors.New("The Aspect of aspectId could not be found.")
	}
	runner, newErr := run.NewRunner(aspectId.String(), code)
	if newErr != nil {
		return false, newErr
	}
	binding, runErr := runner.IsTransactionLevel()
	return binding, runErr
}
