package contract

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/artela-network/artela-evm/vm"
	common2 "github.com/artela-network/artela/common"
	"github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/artela/x/evm/states"
	"github.com/artela-network/aspect-core/djpm/contract"
	"github.com/artela-network/aspect-core/djpm/run"
	artelasdkType "github.com/artela-network/aspect-core/types"
	runtime "github.com/artela-network/aspect-runtime"
	runtimeTypes "github.com/artela-network/aspect-runtime/types"
)

var (
	zero      = uint256.NewInt(0)
	one       = uint256.NewInt(1)
	emptyAddr = common.Address{}
)

type HandlerContext struct {
	cosmosCtx  sdk.Context
	from       common.Address
	parameters map[string]interface{}
	commit     bool
	service    *AspectService
	logger     runtimeTypes.Logger
	evmState   *states.StateDB
	evm        *vm.EVM
	abi        *abi.Method

	rawInput  []byte
	nonce     uint64
	gasLimit  uint64
	gasPrice  *big.Int
	gasTipCap *big.Int
	gasFeeCap *big.Int
}

type Handler interface {
	Handle(ctx *HandlerContext, gas uint64) (ret []byte, remainingGas uint64, err error)
	Method() string
}

type DeployHandler struct{}

func (h DeployHandler) Handle(ctx *HandlerContext, gas uint64) ([]byte, uint64, error) {
	aspectId, code, initData, properties, joinPoint, err := h.decodeAndValidate(ctx)
	if err != nil {
		ctx.logger.Error("deploy aspect failed", "error", err, "from", ctx.from, "gasLimit", ctx.gasLimit)
		return nil, 0, err
	}

	store := ctx.service.aspectStore
	newVersion, gas, err := store.BumpAspectVersion(ctx.cosmosCtx, aspectId, gas)
	if err != nil {
		ctx.logger.Error("bump aspect version failed", "error", err)
		return nil, gas, err
	}

	gas, err = store.StoreAspectCode(ctx.cosmosCtx, aspectId, code, newVersion, gas)
	if err != nil {
		ctx.logger.Error("store aspect code failed", "error", err)
		return nil, gas, err
	}

	// join point might be nil, since there are some operation only Aspects
	if joinPoint != nil {
		store.StoreAspectJP(ctx.cosmosCtx, aspectId, *newVersion, joinPoint)
	}

	if len(properties) > 0 {
		gas, err = store.StoreAspectProperty(ctx.cosmosCtx, aspectId, properties, gas)
		if err != nil {
			ctx.logger.Error("store aspect property failed", "error", err)
		}
	}

	// initialize aspect
	aspectCtx := mustGetAspectContext(ctx.cosmosCtx)
	runner, err := run.NewRunner(aspectCtx, ctx.logger, aspectId.String(), newVersion.Uint64(), code, ctx.commit)
	if err != nil {
		ctx.logger.Error("failed to create aspect runner", "error", err)
		return nil, 0, err
	}
	defer runner.Return()

	var txHash []byte
	ethTxCtx := aspectCtx.EthTxContext()
	if ethTxCtx != nil && ethTxCtx.TxContent() != nil {
		txHash = ethTxCtx.TxContent().Hash().Bytes()
	}

	height := ctx.cosmosCtx.BlockHeight()
	heightU64 := uint64(height)

	return runner.JoinPoint(artelasdkType.INIT_METHOD, gas, height, aspectId, &artelasdkType.InitInput{
		Tx: &artelasdkType.WithFromTxInput{
			Hash: txHash,
			To:   aspectId.Bytes(),
			From: ctx.from.Bytes(),
		},
		Block:    &artelasdkType.BlockInput{Number: &heightU64},
		CallData: initData,
	})
}

func (h DeployHandler) Method() string {
	return "deploy"
}

func (h DeployHandler) decodeAndValidate(ctx *HandlerContext) (aspectId common.Address, code, initData []byte, properties []types.Property, joinPoint *big.Int, err error) {
	// input validations
	code = ctx.parameters["code"].([]byte)
	if len(code) == 0 {
		err = errors.New("code is empty")
		return
	}

	initData = ctx.parameters["initdata"].([]byte)

	propertiesArr := ctx.parameters["properties"].([]struct {
		Key   string `json:"key"`
		Value []byte `json:"value"`
	})

	for i := range propertiesArr {
		s := propertiesArr[i]
		if types.AspectProofKey == s.Key || types.AspectAccountKey == s.Key {
			// Block query of account and Proof
			err = errors.New("using reserved aspect property key")
			return
		}

		properties = append(properties, types.Property{
			Key:   s.Key,
			Value: s.Value,
		})
	}

	account := ctx.parameters["account"].(common.Address)
	if bytes.Equal(account.Bytes(), ctx.from.Bytes()) {
		accountProperty := types.Property{
			Key:   types.AspectAccountKey,
			Value: account.Bytes(),
		}
		properties = append(properties, accountProperty)
	} else {
		err = errors.New("account verification fail")
		return
	}

	proof := ctx.parameters["proof"].([]byte)
	proofProperty := types.Property{
		Key:   types.AspectProofKey,
		Value: proof,
	}
	properties = append(properties, proofProperty)

	joinPoint = ctx.parameters["joinPoints"].(*big.Int)
	if joinPoint == nil {
		err = errors.New("unable to decode join point")
		return
	}

	aspectId = crypto.CreateAddress(ctx.from, ctx.nonce)

	// check duplicate deployment
	if isAspectDeployed(ctx.cosmosCtx, ctx.service.aspectStore, aspectId) {
		err = errors.New("aspect already deployed")
		return
	}

	// validate aspect code
	code, err = validateCode(ctx.cosmosCtx, code)
	return
}

type UpgradeHandler struct{}

func (h UpgradeHandler) Handle(ctx *HandlerContext, gas uint64) ([]byte, uint64, error) {
	aspectId, code, properties, joinPoint, gas, err := h.decodeAndValidate(ctx, gas)
	if err != nil {
		return nil, gas, err
	}

	store := ctx.service.aspectStore
	newVersion, gas, err := store.BumpAspectVersion(ctx.cosmosCtx, aspectId, gas)
	if err != nil {
		ctx.logger.Error("bump aspect version failed", "error", err)
		return nil, gas, err
	}

	if gas, err = store.StoreAspectCode(ctx.cosmosCtx, aspectId, code, newVersion, gas); err != nil {
		ctx.logger.Error("store aspect code failed", "error", err)
		return nil, gas, err
	}

	// join point might be nil, since there are some operation only Aspects
	if joinPoint != nil {
		store.StoreAspectJP(ctx.cosmosCtx, aspectId, *newVersion, joinPoint)
	}

	if len(properties) > 0 {
		gas, err = store.StoreAspectProperty(ctx.cosmosCtx, aspectId, properties, gas)
	}

	return nil, gas, err
}

func (h UpgradeHandler) Method() string {
	return "upgrade"
}

func (h UpgradeHandler) decodeAndValidate(ctx *HandlerContext, gas uint64) (aspectId common.Address,
	code []byte,
	properties []types.Property,
	joinPoint *big.Int, leftover uint64, err error) {
	aspectId = ctx.parameters["aspectId"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), aspectId.Bytes()) {
		err = errors.New("aspectId not specified")
		return
	}

	// input validations
	code = ctx.parameters["code"].([]byte)
	if len(code) == 0 {
		err = errors.New("code is empty")
		return
	}

	propertiesArr := ctx.parameters["properties"].([]struct {
		Key   string `json:"key"`
		Value []byte `json:"value"`
	})

	for i := range propertiesArr {
		s := propertiesArr[i]
		if types.AspectProofKey == s.Key || types.AspectAccountKey == s.Key {
			// Block query of account and Proof
			err = errors.New("using reserved aspect property key")
			return
		}

		properties = append(properties, types.Property{
			Key:   s.Key,
			Value: s.Value,
		})
	}

	joinPoint = ctx.parameters["joinPoints"].(*big.Int)
	if joinPoint == nil {
		joinPoint = big.NewInt(0)
	}

	// check deployment
	store := ctx.service.aspectStore
	if !isAspectDeployed(ctx.cosmosCtx, store, aspectId) {
		err = errors.New("aspect not deployed")
		return
	}

	// check aspect owner
	currentCode, version := store.GetAspectCode(ctx.cosmosCtx, aspectId, nil)

	var ok bool
	ok, leftover, err = checkAspectOwner(ctx.cosmosCtx, aspectId, ctx.from, gas, currentCode, version, ctx.commit)
	if err != nil || !ok {
		err = errors.New("aspect ownership validation failed")
		return
	}

	// validate aspect code
	code, err = validateCode(ctx.cosmosCtx, code)
	return
}

type BindHandler struct{}

func (b BindHandler) Handle(ctx *HandlerContext, gas uint64) (ret []byte, remainingGas uint64, err error) {
	aspectId, account, aspectVersion, priority, isContract, leftover, err := b.decodeAndValidate(ctx, gas)
	if err != nil {
		return nil, leftover, err
	}

	// check aspect types
	store := ctx.service.aspectStore
	aspectJP, err := store.GetAspectJP(ctx.cosmosCtx, aspectId, aspectVersion)
	if err != nil {
		return nil, leftover, err
	}

	txAspect := artelasdkType.CheckIsTransactionLevel(aspectJP.Int64())
	txVerifier := artelasdkType.CheckIsTxVerifier(aspectJP.Int64())

	if !txAspect && !txVerifier {
		return nil, 0, errors.New("aspect is either for tx or verifier")
	}

	// EoA can only bind with tx verifier
	if !txVerifier && !isContract {
		return nil, 0, errors.New("only verifier aspect can be bound with eoa")
	}

	// bind tx processing aspect if account is a contract
	if txAspect && isContract {
		if err := store.BindTxAspect(ctx.cosmosCtx, account, aspectId, aspectVersion, priority); err != nil {
			ctx.logger.Error("bind tx aspect failed", "aspect", aspectId.Hex(), "version", aspectVersion.Uint64(), "contract", account.Hex(), "error", err)
			return nil, 0, err
		}
	}

	// bind tx verifier aspect
	if txVerifier {
		if err := store.BindVerificationAspect(ctx.cosmosCtx, account, aspectId, aspectVersion, priority, isContract); err != nil {
			ctx.logger.Error("bind verifier aspect failed", "aspect", aspectId.Hex(), "version", aspectVersion.Uint64(), "account", account.Hex(), "error", err)
			return nil, 0, err
		}
	}

	// save reverse index
	if err := store.StoreAspectRefValue(ctx.cosmosCtx, account, aspectId); err != nil {
		return nil, 0, err
	}

	return nil, leftover, nil
}

func (b BindHandler) Method() string {
	return "bind"
}

func (b BindHandler) decodeAndValidate(ctx *HandlerContext, gas uint64) (
	aspectId common.Address,
	account common.Address,
	aspectVersion *uint256.Int,
	priority int8,
	isContract bool,
	leftover uint64,
	err error) {
	aspectId = ctx.parameters["aspectId"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), aspectId.Bytes()) {
		err = errors.New("aspectId not specified")
		return
	}

	version := ctx.parameters["aspectVersion"].(*big.Int)
	if version != nil && version.Sign() < 0 {
		err = errors.New("aspectVersion is negative")
		return
	}

	account = ctx.parameters["contract"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), account.Bytes()) {
		err = errors.New("binding account not specified")
		return
	}

	priority = ctx.parameters["priority"].(int8)

	store := ctx.service.aspectStore
	if !isAspectDeployed(ctx.cosmosCtx, store, aspectId) {
		err = errors.New("aspect not deployed")
		return
	}

	isContract = len(ctx.evmState.GetCode(account)) > 0
	if isContract {
		var isOwner bool
		isOwner, leftover = checkContractOwner(ctx, account, gas)
		if !isOwner {
			err = errors.New("contract ownership validation failed")
			return
		}
	} else if !bytes.Equal(account.Bytes(), ctx.from.Bytes()) {
		err = errors.New("unauthorized EoA account aspect binding")
		return
	} else {
		leftover = gas
	}

	aspectVersion, _ = uint256.FromBig(version)

	// overwrite aspect version, just in case if aspect version is 0 which means we will need to overwrite
	// it to latest
	if aspectVersion == nil || aspectVersion.Cmp(zero) <= 0 {
		aspectVersion = ctx.service.aspectStore.GetAspectLastVersion(ctx.cosmosCtx, aspectId)
	}

	return
}

type UnbindHandler struct{}

func (u UnbindHandler) Handle(ctx *HandlerContext, gas uint64) (ret []byte, remainingGas uint64, err error) {
	aspectId, account, isContract, leftover, err := u.decodeAndValidate(ctx, gas)
	if err != nil {
		return nil, leftover, err
	}

	store := ctx.service.aspectStore

	if err := store.UnBindVerificationAspect(ctx.cosmosCtx, account, aspectId); err != nil {
		return nil, leftover, err
	}
	if isContract {
		if err := store.UnBindContractAspects(ctx.cosmosCtx, account, aspectId); err != nil {
			return nil, leftover, err
		}
	}

	if err := store.UnbindAspectRefValue(ctx.cosmosCtx, account, aspectId); err != nil {
		return nil, leftover, err
	}

	return nil, leftover, nil
}

func (u UnbindHandler) decodeAndValidate(ctx *HandlerContext, gas uint64) (
	aspectId common.Address,
	account common.Address,
	isContract bool,
	leftover uint64,
	err error) {
	aspectId = ctx.parameters["aspectId"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), aspectId.Bytes()) {
		err = errors.New("aspectId not specified")
		return
	}

	account = ctx.parameters["contract"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), account.Bytes()) {
		err = errors.New("binding account not specified")
		return
	}

	isContract = len(ctx.evmState.GetCode(account)) > 0
	if isContract {
		var isOwner bool
		isOwner, leftover = checkContractOwner(ctx, account, gas)
		if !isOwner {
			err = errors.New("contract ownership validation failed")
		}
	} else if !bytes.Equal(account.Bytes(), ctx.from.Bytes()) {
		err = errors.New("unauthorized EoA account aspect binding")
	} else {
		leftover = gas
	}

	return
}

func (u UnbindHandler) Method() string {
	return "unbind"
}

type ChangeVersionHandler struct{}

func (c ChangeVersionHandler) Handle(ctx *HandlerContext, gas uint64) (ret []byte, remainingGas uint64, err error) {
	aspectId, account, version, isContract, leftover, err := c.decodeAndValidate(ctx, gas)
	if err != nil {
		return nil, leftover, err
	}

	aspectJP, err := ctx.service.aspectStore.GetAspectJP(ctx.cosmosCtx, aspectId, uint256.NewInt(version))
	if err != nil {
		return nil, leftover, err
	}

	txAspect := artelasdkType.CheckIsTransactionLevel(aspectJP.Int64())
	verifierAspect := artelasdkType.CheckIsTxVerifier(aspectJP.Int64())

	if !txAspect && !verifierAspect {
		return nil, leftover, errors.New("aspect is either for tx or verifier")
	}

	if !verifierAspect && !isContract {
		return nil, 0, errors.New("only verifier aspect can be bound with eoa")
	}

	err = ctx.service.aspectStore.ChangeBoundAspectVersion(ctx.cosmosCtx, account, aspectId, version, isContract, verifierAspect, txAspect)
	remainingGas = leftover
	return
}

func (c ChangeVersionHandler) Method() string {
	return "changeversion"
}

func (c ChangeVersionHandler) decodeAndValidate(ctx *HandlerContext, gas uint64) (
	aspectId common.Address,
	account common.Address,
	version uint64,
	isContract bool,
	leftover uint64,
	err error,
) {
	aspectId = ctx.parameters["aspectId"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), aspectId.Bytes()) {
		err = errors.New("aspectId not specified")
		return
	}

	account = ctx.parameters["contract"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), account.Bytes()) {
		err = errors.New("binding account not specified")
		return
	}

	// should check whether expected version is greater than
	// the latest version we have, if so, the designated aspect
	// does not exist yet
	version = ctx.parameters["version"].(uint64)
	store := ctx.service.aspectStore
	latestVersion := store.GetAspectLastVersion(ctx.cosmosCtx, aspectId)
	if latestVersion == nil || latestVersion.Cmp(zero) == 0 || latestVersion.Uint64() < version {
		err = errors.New("given version of aspect does not exist")
		return
	}
	if version == 0 {
		version = latestVersion.Uint64()
	}

	if isContract = len(ctx.evmState.GetCode(account)) > 0; isContract {
		var isOwner bool
		if isOwner, leftover = checkContractOwner(ctx, account, gas); !isOwner {
			err = errors.New("unauthorized operation")
		}
	} else if !bytes.Equal(account.Bytes(), ctx.from.Bytes()) {
		err = errors.New("unauthorized operation")
	} else {
		leftover = gas
	}

	return
}

type GetVersionHandler struct{}

func (g GetVersionHandler) Handle(ctx *HandlerContext, gas uint64) (ret []byte, remainingGas uint64, err error) {
	aspectId, err := g.decodeAndValidate(ctx)
	if err != nil {
		return nil, 0, err
	}

	version := ctx.service.aspectStore.GetAspectLastVersion(ctx.cosmosCtx, aspectId)

	ret, err = ctx.abi.Outputs.Pack(version.Uint64())
	if err != nil {
		return nil, gas, err
	}

	return ret, gas, nil
}

func (g GetVersionHandler) Method() string {
	return "versionof"
}

func (g GetVersionHandler) decodeAndValidate(ctx *HandlerContext) (aspectId common.Address, err error) {
	aspectId = ctx.parameters["aspectId"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), aspectId.Bytes()) {
		err = errors.New("aspectId not specified")
		return
	}

	return
}

type GetBindingHandler struct{}

func (g GetBindingHandler) Handle(ctx *HandlerContext, gas uint64) (ret []byte, remainingGas uint64, err error) {
	account, isContract, err := g.decodeAndValidate(ctx)
	if err != nil {
		return nil, 0, err
	}

	aspectInfo := make([]types.AspectInfo, 0)
	deduplicationMap := make(map[common.Address]struct{})

	accountVerifiers, err := ctx.service.aspectStore.GetVerificationAspects(ctx.cosmosCtx, account)
	if err != nil {
		return nil, 0, err
	}

	for _, aspect := range accountVerifiers {
		if _, exist := deduplicationMap[aspect.Id]; exist {
			continue
		}
		deduplicationMap[aspect.Id] = struct{}{}
		info := types.AspectInfo{
			AspectId: aspect.Id,
			Version:  aspect.Version.Uint64(),
			Priority: int8(aspect.Priority),
		}
		aspectInfo = append(aspectInfo, info)
	}

	if isContract {
		txLevelAspects, err := ctx.service.aspectStore.GetTxLevelAspects(ctx.cosmosCtx, account)
		if err != nil {
			return nil, 0, err
		}

		for _, aspect := range txLevelAspects {
			if _, exist := deduplicationMap[aspect.Id]; exist {
				continue
			}
			deduplicationMap[aspect.Id] = struct{}{}
			info := types.AspectInfo{
				AspectId: aspect.Id,
				Version:  aspect.Version.Uint64(),
				Priority: int8(aspect.Priority),
			}
			aspectInfo = append(aspectInfo, info)
		}
	}

	ret, err = ctx.abi.Outputs.Pack(aspectInfo)
	return ret, gas, err
}

func (g GetBindingHandler) Method() string {
	return "aspectsof"
}

func (g GetBindingHandler) decodeAndValidate(ctx *HandlerContext) (account common.Address, isContract bool, err error) {
	account = ctx.parameters["contract"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), account.Bytes()) {
		err = errors.New("binding account not specified")
		return
	}

	isContract = len(ctx.evmState.GetCode(account)) > 0
	return
}

type GetBoundAddressHandler struct{}

func (g GetBoundAddressHandler) Handle(ctx *HandlerContext, gas uint64) (ret []byte, remainingGas uint64, err error) {
	aspectId, err := g.decodeAndValidate(ctx)
	if err != nil {
		return nil, 0, err
	}

	value, err := ctx.service.GetAspectOf(ctx.cosmosCtx, aspectId)
	if err != nil {
		return nil, 0, err
	}
	addressArr := make([]common.Address, 0)
	if value != nil {
		for _, data := range value.Values() {
			contractAddr := common.HexToAddress(data.(string))
			addressArr = append(addressArr, contractAddr)
		}
	}

	ret, err = ctx.abi.Outputs.Pack(addressArr)
	return ret, gas, err
}

func (g GetBoundAddressHandler) Method() string {
	return "boundaddressesof"
}

func (g GetBoundAddressHandler) decodeAndValidate(ctx *HandlerContext) (aspectId common.Address, err error) {
	aspectId = ctx.parameters["aspectId"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), aspectId.Bytes()) {
		err = errors.New("aspect id not specified")
		return
	}

	if !isAspectDeployed(ctx.cosmosCtx, ctx.service.aspectStore, aspectId) {
		err = errors.New("aspect not deployed")
	}

	return
}

type OperationHandler struct{}

func (o OperationHandler) Handle(ctx *HandlerContext, gas uint64) (ret []byte, remainingGas uint64, err error) {
	aspectId, args, err := o.decodeAndValidate(ctx)
	if err != nil {
		return nil, 0, err
	}

	lastHeight := ctx.cosmosCtx.BlockHeight()
	code, version := ctx.service.GetAspectCode(ctx.cosmosCtx, aspectId, nil)

	aspectCtx := mustGetAspectContext(ctx.cosmosCtx)
	runner, err := run.NewRunner(aspectCtx, ctx.logger, aspectId.String(), version.Uint64(), code, ctx.commit)
	if err != nil {
		ctx.logger.Error("failed to create aspect runner", "error", err)
		return nil, 0, err
	}
	defer runner.Return()

	var txHash []byte
	ethTxCtx := aspectCtx.EthTxContext()
	if ethTxCtx != nil && ethTxCtx.TxContent() != nil {
		txHash = ethTxCtx.TxContent().Hash().Bytes()
	}
	height := uint64(lastHeight)
	ret, remainingGas, err = runner.JoinPoint(artelasdkType.OPERATION_METHOD, gas, lastHeight, aspectId, &artelasdkType.OperationInput{
		Tx: &artelasdkType.WithFromTxInput{
			Hash: txHash,
			To:   aspectId.Bytes(),
			From: ctx.from.Bytes(),
		},
		Block:    &artelasdkType.BlockInput{Number: &height},
		CallData: args,
	})
	if err == nil {
		ret, err = ctx.abi.Outputs.Pack(ret)
		if err != nil {
			ctx.logger.Error("failed to pack operation output", "error", err)
		}
	}

	return
}

func (o OperationHandler) Method() string {
	return "entrypoint"
}

func (o OperationHandler) decodeAndValidate(ctx *HandlerContext) (aspectId common.Address, args []byte, err error) {
	aspectId = ctx.parameters["aspectId"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), aspectId.Bytes()) {
		err = errors.New("aspect id not specified")
		return
	}

	if !isAspectDeployed(ctx.cosmosCtx, ctx.service.aspectStore, aspectId) {
		err = errors.New("aspect not deployed")
		return
	}

	args = ctx.parameters["optArgs"].([]byte)
	return
}

func isAspectDeployed(ctx sdk.Context, store *AspectStore, aspectId common.Address) bool {
	return store.GetAspectLastVersion(ctx, aspectId).Cmp(zero) > 0
}

func validateCode(ctx sdk.Context, aspectCode []byte) ([]byte, error) {
	startTime := time.Now()
	validator, err := runtime.NewValidator(ctx, common2.WrapLogger(ctx.Logger()), runtime.WASM)
	if err != nil {
		return nil, err
	}
	ctx.Logger().Info("validated aspect bytecode", "duration", time.Since(startTime).String())

	startTime = time.Now()
	parsed, err := ParseByteCode(aspectCode)
	if err != nil {
		return nil, err
	}
	ctx.Logger().Info("parsed aspect bytecode", "duration", time.Since(startTime).String())

	return parsed, validator.Validate(parsed)
}

func checkContractOwner(ctx *HandlerContext, contractAddr common.Address, gas uint64) (bool, uint64) {
	msg, err := contract.ArtelaOwnerMsg(&contractAddr, ctx.nonce, ctx.from, gas, ctx.gasPrice, ctx.gasFeeCap, ctx.gasTipCap)
	if err != nil {
		return false, 0
	}
	fromAccount := vm.AccountRef(msg.From)
	ctx.evm.CloseAspectCall()
	defer ctx.evm.AspectCall()

	aspectCtx := mustGetAspectContext(ctx.cosmosCtx)
	ret, leftover, err := ctx.evm.Call(aspectCtx, fromAccount, *msg.To, msg.Data, gas, msg.Value)
	if err != nil {
		// if fail, fallback to openzeppelin ownable
		msg, err = contract.OpenZeppelinOwnableMsg(&contractAddr, ctx.nonce, ctx.from, gas, ctx.gasPrice, ctx.gasFeeCap, ctx.gasTipCap)
		if err != nil {
			return false, 0
		}

		ret, leftover, err = ctx.evm.Call(aspectCtx, fromAccount, *msg.To, msg.Data, gas, msg.Value)
		if err != nil {
			return false, leftover
		}

		result, err := contract.UnpackOwnableOwnerResult(ret)
		if err != nil {
			return false, leftover
		}

		return bytes.Equal(result.Bytes(), ctx.from.Bytes()), leftover
	}

	result, err := contract.UnpackIsOwnerResult(ret)
	if err != nil {
		return false, leftover
	}
	return result, leftover
}

func checkAspectOwner(ctx sdk.Context, aspectId common.Address, sender common.Address, gas uint64, code []byte, version *uint256.Int, commit bool) (bool, uint64, error) {
	aspectCtx := mustGetAspectContext(ctx)
	runner, err := run.NewRunner(aspectCtx, common2.WrapLogger(ctx.Logger()), aspectId.String(), version.Uint64(), code, commit)
	if err != nil {
		panic(fmt.Sprintf("failed to create runner: %v", err))
	}
	defer runner.Return()

	return runner.IsOwner(ctx.BlockHeight(), gas, sender, sender.Bytes())
}

// retrieving aspect context from sdk.Context must not fail, so we panic if it does
func mustGetAspectContext(ctx sdk.Context) *types.AspectRuntimeContext {
	aspectCtx, ok := ctx.Value(types.AspectContextKey).(*types.AspectRuntimeContext)
	if !ok {
		panic("unable to get aspect context, this should not happen")
	}

	return aspectCtx
}
