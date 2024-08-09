package contract

import (
	"bytes"
	"errors"
	"fmt"
	asptool "github.com/artela-network/artela/x/aspect/common"
	"github.com/artela-network/artela/x/aspect/store"
	aspectmoduletypes "github.com/artela-network/artela/x/aspect/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/artela-network/artela-evm/vm"
	arttool "github.com/artela-network/artela/common"
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
	logger     runtimeTypes.Logger
	evmState   *states.StateDB
	evm        *vm.EVM
	abi        *abi.Method

	evmStoreKey    storetypes.StoreKey
	aspectStoreKey storetypes.StoreKey

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
	aspectID, code, initData, properties, joinPoint, paymaster, proof, err := h.decodeAndValidate(ctx)
	if err != nil {
		ctx.logger.Error("deploy aspect failed", "error", err, "from", ctx.from, "gasLimit", ctx.gasLimit)
		return nil, 0, err
	}

	// we can ignore the new store here, since new deployed aspect should not have that
	metaStore, _, err := store.GetAspectMetaStore(buildAspectStoreCtx(ctx, aspectID, gas))
	if err != nil {
		return nil, 0, err
	}

	// check duplicate deployment
	var latestVersion uint64
	if latestVersion, err = metaStore.GetLatestVersion(); err != nil {
		return nil, 0, err
	} else if latestVersion > 0 {
		return nil, 0, errors.New("aspect already deployed")
	}

	if err := metaStore.Init(); err != nil {
		ctx.logger.Error("init aspect meta failed", "error", err)
		return nil, 0, err
	}

	newVersion, err := metaStore.BumpVersion()
	if err != nil {
		ctx.logger.Error("bump aspect version failed", "error", err)
		return nil, 0, err
	}

	if err = metaStore.StoreCode(newVersion, code); err != nil {
		ctx.logger.Error("store aspect code failed", "error", err)
		return nil, 0, err
	}

	// join point might be nil, since there are some operation only Aspects
	if err = metaStore.StoreVersionMeta(newVersion, &aspectmoduletypes.VersionMeta{
		JoinPoint: joinPoint.Uint64(),
		CodeHash:  crypto.Keccak256Hash(code),
	}); err != nil {
		ctx.logger.Error("store aspect meta failed", "error", err)
		return nil, 0, err
	}

	if err = metaStore.StoreMeta(&aspectmoduletypes.AspectMeta{
		Proof:     proof,
		PayMaster: paymaster,
	}); err != nil {
		ctx.logger.Error("store aspect meta failed", "error", err)
		return nil, 0, err
	}

	if err = metaStore.StoreProperties(newVersion, properties); err != nil {
		ctx.logger.Error("store aspect property failed", "error", err)
		return nil, 0, err
	}

	// get remaining gas after updating store
	gas = metaStore.Gas()

	// initialize aspect
	aspectCtx := mustGetAspectContext(ctx.cosmosCtx)
	runner, err := run.NewRunner(aspectCtx, ctx.logger, aspectID.String(), newVersion, code, ctx.commit)
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

	return runner.JoinPoint(artelasdkType.INIT_METHOD, gas, height, aspectID, &artelasdkType.InitInput{
		Tx: &artelasdkType.WithFromTxInput{
			Hash: txHash,
			To:   aspectID.Bytes(),
			From: ctx.from.Bytes(),
		},
		Block:    &artelasdkType.BlockInput{Number: &heightU64},
		CallData: initData,
	})
}

func (h DeployHandler) Method() string {
	return "deploy"
}

func (h DeployHandler) decodeAndValidate(ctx *HandlerContext) (
	aspectId common.Address,
	code,
	initData []byte,
	properties []aspectmoduletypes.Property,
	joinPoint *big.Int,
	paymaster common.Address,
	proof []byte,
	err error) {
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
		properties = append(properties, aspectmoduletypes.Property{
			Key:   s.Key,
			Value: s.Value,
		})
	}

	paymaster = ctx.parameters["account"].(common.Address)
	if !bytes.Equal(paymaster.Bytes(), ctx.from.Bytes()) {
		err = errors.New("account verification fail")
		return
	}

	proof = ctx.parameters["proof"].([]byte)

	joinPoint = ctx.parameters["joinPoints"].(*big.Int)
	if joinPoint == nil {
		joinPoint = big.NewInt(0)
	}

	aspectId = crypto.CreateAddress(ctx.from, ctx.nonce)

	// validate aspect code
	code, err = validateCode(ctx.cosmosCtx, code)
	return
}

type UpgradeHandler struct{}

func (h UpgradeHandler) Handle(ctx *HandlerContext, gas uint64) ([]byte, uint64, error) {
	aspectID, code, properties, joinPoint, err := h.decodeAndValidate(ctx)
	if err != nil {
		return nil, 0, err
	}

	// check deployment
	storeCtx := buildAspectStoreCtx(ctx, aspectID, gas)
	currentStore, _, err := store.GetAspectMetaStore(storeCtx)
	if err != nil {
		return nil, 0, err
	}

	// check deployment
	var latestVersion uint64
	if latestVersion, err = currentStore.GetLatestVersion(); err != nil {
		return nil, 0, err
	} else if latestVersion == 0 {
		return nil, 0, errors.New("aspect not deployed")
	}

	// check aspect owner
	var currentCode []byte
	currentCode, err = currentStore.GetCode(latestVersion)
	if err != nil {
		return nil, 0, err
	}

	var ok bool
	ok, gas, err = checkAspectOwner(ctx.cosmosCtx, aspectID, ctx.from, storeCtx.Gas(), currentCode, latestVersion, ctx.commit)
	if err != nil || !ok {
		err = errors.New("aspect ownership validation failed")
		return nil, 0, err
	}
	storeCtx.UpdateGas(gas)

	// bump version
	newVersion, err := currentStore.BumpVersion()
	if err != nil {
		ctx.logger.Error("bump aspect version failed", "error", err)
		return nil, gas, err
	}

	if err = currentStore.StoreCode(newVersion, code); err != nil {
		ctx.logger.Error("store aspect code failed", "error", err)
		return nil, 0, err
	}

	// join point might be nil, since there are some operation only Aspects
	var jpU64 uint64
	if joinPoint != nil {
		jpU64 = joinPoint.Uint64()
	}

	if err = currentStore.StoreVersionMeta(newVersion, &aspectmoduletypes.VersionMeta{
		JoinPoint: jpU64,
		CodeHash:  common.BytesToHash(crypto.Keccak256(code)),
	}); err != nil {
		ctx.logger.Error("store aspect meta failed", "error", err)
		return nil, 0, err
	}

	// save properties if any
	if err = currentStore.StoreProperties(newVersion, properties); err != nil {
		ctx.logger.Error("store aspect property failed", "error", err)
		return nil, 0, err
	}

	return nil, storeCtx.Gas(), err
}

func (h UpgradeHandler) Method() string {
	return "upgrade"
}

func (h UpgradeHandler) decodeAndValidate(ctx *HandlerContext) (
	aspectID common.Address,
	code []byte,
	properties []aspectmoduletypes.Property,
	joinPoint *big.Int, err error) {
	aspectID = ctx.parameters["aspectId"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), aspectID.Bytes()) {
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

	for _, prop := range propertiesArr {
		properties = append(properties, aspectmoduletypes.Property{
			Key:   prop.Key,
			Value: prop.Value,
		})
	}

	joinPoint = ctx.parameters["joinPoints"].(*big.Int)
	if joinPoint == nil {
		joinPoint = big.NewInt(0)
	}

	// validate aspect code
	code, err = validateCode(ctx.cosmosCtx, code)
	return
}

type BindHandler struct{}

func (b BindHandler) Handle(ctx *HandlerContext, gas uint64) (ret []byte, remainingGas uint64, err error) {
	aspectID, account, aspectVersion, priority, isContract, leftover, err := b.decodeAndValidate(ctx, gas)
	if err != nil {
		return nil, 0, err
	}

	// no need to get new version store, since we only do auto migration during aspect upgrade
	metaStore, _, err := store.GetAspectMetaStore(buildAspectStoreCtx(ctx, aspectID, leftover))
	if err != nil {
		return nil, 0, err
	}

	// get latest version if aspect version is empty
	latestVersion, err := metaStore.GetLatestVersion()
	if err != nil {
		return nil, 0, err
	}
	if latestVersion == 0 {
		return nil, 0, errors.New("aspect not deployed")
	}

	if aspectVersion == 0 {
		// use latest if not specified
		aspectVersion = latestVersion
	} else if aspectVersion > latestVersion {
		return nil, 0, errors.New("aspect version not deployed")
	}

	meta, err := metaStore.GetVersionMeta(aspectVersion)
	if err != nil {
		return nil, 0, err
	}

	i64JP := int64(meta.JoinPoint)
	txAspect := artelasdkType.CheckIsTransactionLevel(i64JP)
	txVerifier := artelasdkType.CheckIsTxVerifier(i64JP)

	if !txAspect && !txVerifier {
		return nil, 0, errors.New("aspect is either for tx or verifier")
	}

	// EoA can only bind with tx verifier
	if !txVerifier && !isContract {
		return nil, 0, errors.New("only verifier aspect can be bound with eoa")
	}

	// save aspect -> contract bindings
	if err := metaStore.StoreBinding(account, aspectVersion, meta.JoinPoint, priority); err != nil {
		return nil, 0, err
	}

	// init account store
	accountStore, _, err := store.GetAccountStore(buildAccountStoreCtx(ctx, account, metaStore.Gas()))
	if err != nil {
		return nil, 0, err
	}

	// check if used
	if used, err := accountStore.Used(); err != nil {
		return nil, 0, err
	} else if !used {
		// init if not used
		if err := accountStore.Init(); err != nil {
			return nil, 0, err
		}
	}

	// save account -> contract bindings
	if err := accountStore.StoreBinding(aspectID, aspectVersion, meta.JoinPoint, priority, isContract); err != nil {
		ctx.logger.Error("bind tx aspect failed", "aspect", aspectID.Hex(), "version", aspectVersion, "contract", account.Hex(), "error", err)
		return nil, 0, err
	}

	return nil, accountStore.Gas(), nil
}

func (b BindHandler) Method() string {
	return "bind"
}

func (b BindHandler) decodeAndValidate(ctx *HandlerContext, gas uint64) (
	aspectId common.Address,
	account common.Address,
	aspectVersion uint64,
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
	} else if version == nil {
		aspectVersion = 0
	} else {
		aspectVersion = version.Uint64()
	}

	account = ctx.parameters["contract"].(common.Address)
	if bytes.Equal(emptyAddr.Bytes(), account.Bytes()) {
		err = errors.New("binding account not specified")
		return
	}

	priority = ctx.parameters["priority"].(int8)

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

	return
}

type UnbindHandler struct{}

func (u UnbindHandler) Handle(ctx *HandlerContext, gas uint64) (ret []byte, remainingGas uint64, err error) {
	aspectID, account, isContract, leftover, err := u.decodeAndValidate(ctx, gas)
	if err != nil {
		return nil, leftover, err
	}

	// init account store
	accountStore, _, err := store.GetAccountStore(buildAccountStoreCtx(ctx, account, leftover))
	if err != nil {
		return nil, 0, err
	}

	bindings, err := accountStore.LoadAccountBoundAspects(aspectmoduletypes.NewDefaultFilter(isContract))
	if err != nil {
		return nil, 0, err
	}

	// looking for binding info
	version := uint64(0)
	for _, binding := range bindings {
		if binding.Account == aspectID {
			version = binding.Version
			break
		}
	}

	if version == 0 {
		// not bound
		return nil, accountStore.Gas(), nil
	}

	// init aspect meta store
	metaStore, _, err := store.GetAspectMetaStore(buildAspectStoreCtx(ctx, aspectID, accountStore.Gas()))
	if err != nil {
		return nil, 0, err
	}

	// remove account from aspect bound list
	if err := metaStore.RemoveBinding(account); err != nil {
		return nil, 0, err
	}
	// load aspect join point
	meta, err := metaStore.GetVersionMeta(version)
	if err != nil {
		return nil, 0, err
	}

	// remove aspect from account bound list
	accountStore.TransferGasFrom(metaStore)
	if err := accountStore.RemoveBinding(aspectID, meta.JoinPoint, isContract); err != nil {
		return nil, 0, err
	}

	return nil, accountStore.Gas(), nil
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
	aspectID, account, version, isContract, leftover, err := c.decodeAndValidate(ctx, gas)
	if err != nil {
		return nil, leftover, err
	}

	// init account store
	accountStore, _, err := store.GetAccountStore(buildAccountStoreCtx(ctx, account, leftover))
	if err != nil {
		return nil, 0, err
	}

	bindings, err := accountStore.LoadAccountBoundAspects(aspectmoduletypes.NewDefaultFilter(isContract))
	if err != nil {
		return nil, 0, err
	}

	// looking for binding info
	var bindingInfo *aspectmoduletypes.Binding
	for _, binding := range bindings {
		if binding.Account == aspectID {
			bindingInfo = &binding
			break
		}
	}

	if bindingInfo == nil {
		// not bound
		return nil, 0, errors.New("aspect not bound")
	}

	// init aspect meta store
	metaStoreCtx := buildAspectStoreCtx(ctx, aspectID, accountStore.Gas())
	metaStore, _, err := store.GetAspectMetaStore(metaStoreCtx)
	if err != nil {
		return nil, 0, err
	}

	// load current bound version aspect meta
	currentVersionMeta, err := metaStore.GetVersionMeta(bindingInfo.Version)
	if err != nil {
		return nil, 0, err
	}

	latestVersion, err := metaStore.GetLatestVersion()
	if latestVersion == 0 {
		return nil, 0, errors.New("aspect not deployed")
	}

	if version > latestVersion {
		return nil, 0, errors.New("given version of aspect does not exist")
	}

	if version == 0 {
		// use latest if not specified
		version = latestVersion
	}

	// load new aspect meta
	newVersionMeta, err := metaStore.GetVersionMeta(version)
	if err != nil {
		return nil, 0, err
	}

	i64JP := int64(newVersionMeta.JoinPoint)
	txAspect := artelasdkType.CheckIsTransactionLevel(i64JP)
	txVerifier := artelasdkType.CheckIsTxVerifier(i64JP)

	if !txAspect && !txVerifier {
		return nil, 0, errors.New("aspect is either for tx or verifier")
	}

	// EoA can only bind with tx verifier
	if !txVerifier && !isContract {
		return nil, 0, errors.New("only verifier aspect can be bound with eoa")
	}

	// remove old version aspect from account bound list
	accountStore.TransferGasFrom(metaStore)
	if err := accountStore.RemoveBinding(aspectID, currentVersionMeta.JoinPoint, isContract); err != nil {
		return nil, 0, err
	}

	// update new binding in account store
	if err := accountStore.StoreBinding(aspectID, version, newVersionMeta.JoinPoint, bindingInfo.Priority, isContract); err != nil {
		return nil, 0, err
	}

	return nil, accountStore.Gas(), nil
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

	version = ctx.parameters["version"].(uint64)

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
	aspectID, err := g.decodeAndValidate(ctx)
	if err != nil {
		return nil, 0, err
	}

	// no need to get new version store, since we only do auto migration during aspect upgrade
	metaStore, _, err := store.GetAspectMetaStore(buildAspectStoreCtx(ctx, aspectID, gas))
	if err != nil {
		return nil, 0, err
	}

	version, err := metaStore.GetLatestVersion()
	if err != nil {
		return nil, 0, err
	}

	ret, err = ctx.abi.Outputs.Pack(version)
	if err != nil {
		return nil, gas, err
	}

	return ret, metaStore.Gas(), nil
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

	// init account store
	accountStore, _, err := store.GetAccountStore(buildAccountStoreCtx(ctx, account, gas))
	if err != nil {
		return nil, 0, err
	}

	bindings, err := accountStore.LoadAccountBoundAspects(aspectmoduletypes.NewDefaultFilter(isContract))
	if err != nil {
		return nil, 0, err
	}

	aspectInfo := make([]types.AspectInfo, 0)
	for _, binding := range bindings {
		info := types.AspectInfo{
			AspectId: binding.Account,
			Version:  binding.Version,
			Priority: binding.Priority,
		}
		aspectInfo = append(aspectInfo, info)
	}

	ret, err = ctx.abi.Outputs.Pack(aspectInfo)
	return ret, accountStore.Gas(), err
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
	aspectID, err := g.decodeAndValidate(ctx)
	if err != nil {
		return nil, 0, err
	}

	// init account store
	metaStore, _, err := store.GetAspectMetaStore(buildAspectStoreCtx(ctx, aspectID, gas))
	if err != nil {
		return nil, 0, err
	}

	// check deployment
	if latestVersion, err := metaStore.GetLatestVersion(); err != nil {
		return nil, 0, err
	} else if latestVersion == 0 {
		return nil, 0, errors.New("aspect not deployed")
	}

	bindings, err := metaStore.LoadAspectBoundAccounts()
	if err != nil {
		return nil, 0, err
	}

	addressArr := make([]common.Address, 0)
	for _, binding := range bindings {
		addressArr = append(addressArr, binding.Account)
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

	return
}

type OperationHandler struct{}

func (o OperationHandler) Handle(ctx *HandlerContext, gas uint64) (ret []byte, remainingGas uint64, err error) {
	aspectID, args, err := o.decodeAndValidate(ctx)
	if err != nil {
		return nil, 0, err
	}

	lastHeight := ctx.cosmosCtx.BlockHeight()
	metaStore, _, err := store.GetAspectMetaStore(buildAspectStoreCtx(ctx, aspectID, gas))
	if err != nil {
		return nil, 0, err
	}

	// check deployment
	latestVersion, err := metaStore.GetLatestVersion()
	if err != nil {
		return nil, 0, err
	} else if latestVersion == 0 {
		return nil, 0, errors.New("aspect not deployed")
	}

	code, err := metaStore.GetCode(latestVersion)
	if err != nil {
		return nil, 0, err
	}

	aspectCtx := mustGetAspectContext(ctx.cosmosCtx)
	runner, err := run.NewRunner(aspectCtx, ctx.logger, aspectID.String(), latestVersion, code, ctx.commit)
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
	ret, remainingGas, err = runner.JoinPoint(artelasdkType.OPERATION_METHOD, gas, lastHeight, aspectID, &artelasdkType.OperationInput{
		Tx: &artelasdkType.WithFromTxInput{
			Hash: txHash,
			To:   aspectID.Bytes(),
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

	args = ctx.parameters["optArgs"].([]byte)
	return
}

func validateCode(ctx sdk.Context, aspectCode []byte) ([]byte, error) {
	startTime := time.Now()
	validator, err := runtime.NewValidator(ctx, arttool.WrapLogger(ctx.Logger()), runtime.WASM)
	if err != nil {
		return nil, err
	}
	ctx.Logger().Info("validated aspect bytecode", "duration", time.Since(startTime).String())

	startTime = time.Now()
	parsed, err := asptool.ParseByteCode(aspectCode)
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

func checkAspectOwner(ctx sdk.Context, aspectId common.Address, sender common.Address, gas uint64, code []byte, version uint64, commit bool) (bool, uint64, error) {
	aspectCtx := mustGetAspectContext(ctx)
	runner, err := run.NewRunner(aspectCtx, arttool.WrapLogger(ctx.Logger()), aspectId.String(), version, code, commit)
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

func buildAspectStoreCtx(ctx *HandlerContext, aspectID common.Address, gas uint64) *aspectmoduletypes.AspectStoreContext {
	return &aspectmoduletypes.AspectStoreContext{
		StoreContext: aspectmoduletypes.NewStoreContext(ctx.cosmosCtx, ctx.aspectStoreKey, ctx.evmStoreKey, gas),
		AspectID:     aspectID,
	}
}

func buildAccountStoreCtx(ctx *HandlerContext, account common.Address, gas uint64) *aspectmoduletypes.AccountStoreContext {
	return &aspectmoduletypes.AccountStoreContext{
		StoreContext: aspectmoduletypes.NewStoreContext(ctx.cosmosCtx, ctx.aspectStoreKey, ctx.evmStoreKey, gas),
		Account:      account,
	}
}
