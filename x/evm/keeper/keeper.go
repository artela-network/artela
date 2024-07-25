package keeper

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"sync"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	paramsmodule "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethereum "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"

	"github.com/artela-network/artela-evm/vm"
	common2 "github.com/artela-network/artela/common"
	artela "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/artela/api"
	"github.com/artela-network/artela/x/evm/artela/provider"
	artvmtype "github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/artela/x/evm/states"
	"github.com/artela-network/artela/x/evm/txs"
	"github.com/artela-network/artela/x/evm/txs/support"
	"github.com/artela-network/artela/x/evm/types"
	"github.com/artela-network/aspect-core/djpm"
	artelaType "github.com/artela-network/aspect-core/types"
)

// Keeper grants access to the EVM module states and implements the go-ethereum StateDB interface.
type Keeper struct {
	// logger saves the logger instance of evm module
	logger log.Logger

	// protobuf codec
	cdc codec.BinaryCodec

	// store key required for the EVM Prefix KVStore. It is required by:
	// 		- storing account's Storage State
	// 		- storing account's Code
	// 		- storing txs Logs
	// 		- storing Bloom filters by block height. Needed for the Web3 API.
	storeKey storetypes.StoreKey

	// key to access the transient store, which is reset on every block during Commit
	transientKey storetypes.StoreKey

	// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
	authority cosmos.AccAddress
	// access to account states
	accountKeeper types.AccountKeeper
	// update balance and accounting operations with coins
	bankKeeper types.BankKeeper
	// access historical headers for EVM states transition execution
	stakingKeeper types.StakingKeeper
	// fetch EIP1559 base fee and parameters
	feeKeeper types.FeeKeeper

	// chain ID number obtained from the context's chain id
	eip155ChainID *big.Int

	// tracer used to collect execution traces from the EVM txs execution
	tracer string

	// legacy subspace
	ss paramsmodule.Subspace

	// keep the evm and matched stateDB instance just finished running
	aspectRuntimeContext *artvmtype.AspectRuntimeContext

	aspect *provider.ArtelaProvider

	clientContext client.Context

	// store the block context, this will be fresh every block.
	BlockContext *artvmtype.EthBlockContext

	// cache of aspect sig
	VerifySigCache *sync.Map
}

// NewKeeper generates new evm module keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey, transientKey storetypes.StoreKey,
	authority cosmos.AccAddress,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	feeKeeper types.FeeKeeper,
	tracer string,
	subSpace paramsmodule.Subspace,
	app *baseapp.BaseApp,
	logger log.Logger,
) *Keeper {
	// ensure evm module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the EVM module account has not been set")
	}

	// ensure the authority account is correct
	if err := cosmos.VerifyAddressFormat(authority); err != nil {
		panic(err)
	}

	// init aspect
	aspect := provider.NewArtelaProvider(storeKey, app.LastBlockHeight, logger)
	// new Aspect Runtime Context
	aspectRuntimeContext := artvmtype.NewAspectRuntimeContext()
	aspectRuntimeContext.Init(storeKey)

	// pass in the parameter space to the CommitStateDB in order to use custom denominations for the EVM operations
	k := &Keeper{
		logger:               logger.With("module", fmt.Sprintf("x/%s", types.ModuleName)),
		cdc:                  cdc,
		authority:            authority,
		accountKeeper:        accountKeeper,
		bankKeeper:           bankKeeper,
		stakingKeeper:        stakingKeeper,
		feeKeeper:            feeKeeper,
		storeKey:             storeKey,
		transientKey:         transientKey,
		tracer:               tracer,
		ss:                   subSpace,
		aspectRuntimeContext: aspectRuntimeContext,
		aspect:               aspect,
		VerifySigCache:       &sync.Map{},
	}
	k.WithChainID(app.ChainId())

	djpm.NewAspect(aspect, common2.WrapLogger(k.logger.With("module", "aspect")))
	api.InitAspectGlobals(k)

	// init aspect host api factory
	artelaType.GetEvmHostHook = api.GetEvmHostInstance
	artelaType.GetStateDbHook = api.GetStateDBHostInstance
	artelaType.GetAspectRuntimeContextHostHook = api.GetAspectRuntimeContextHostInstance
	artelaType.GetAspectStateHostHook = api.GetAspectStateHostInstance
	artelaType.GetAspectPropertyHostHook = api.GetAspectPropertyHostInstance
	artelaType.GetAspectTransientStorageHostHook = api.GetAspectTransientStorageHostInstance
	artelaType.GetAspectTraceHostHook = api.GetAspectTraceHostInstance

	artelaType.GetAspectContext = k.GetAspectContext
	artelaType.SetAspectContext = k.SetAspectContext

	artelaType.JITSenderAspectByContext = k.JITSenderAspectByContext
	artelaType.IsCommit = k.IsCommit
	return k
}

func (k *Keeper) SetClientContext(ctx client.Context) {
	k.clientContext = ctx
	// k.aspectRuntimeContext.ExtBlockContext().WithBlockConfig(nil, nil, &ctx)
}

func (k Keeper) GetClientContext() client.Context {
	return k.clientContext
}

// ----------------------------------------------------------------------------
// 								   Config
// ----------------------------------------------------------------------------

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx cosmos.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// WithChainID sets the chain id to the local variable in the keeper
func (k *Keeper) WithChainID(chainId string) {
	if k.eip155ChainID != nil {
		return
	}

	chainID, err := artela.ParseChainID(chainId)
	if err != nil {
		panic(err)
	}

	if k.eip155ChainID != nil && k.eip155ChainID.Cmp(chainID) != 0 {
		panic("chain id already set")
	}

	k.eip155ChainID = chainID
}

// ChainID returns the EIP155 chain ID for the EVM context
func (k Keeper) ChainID() *big.Int {
	return k.eip155ChainID
}

// GetAuthority returns the x/evm module authority address
func (k Keeper) GetAuthority() cosmos.AccAddress {
	return k.authority
}

// ----------------------------------------------------------------------------
// 								Block Bloom
// 							Required by Web3 API
// ----------------------------------------------------------------------------

// EmitBlockBloomEvent emit block bloom events
func (k Keeper) EmitBlockBloomEvent(ctx cosmos.Context, bloom ethereum.Bloom) {
	encodedBloom := base64.StdEncoding.EncodeToString(bloom.Bytes())

	sprintf := fmt.Sprintf("emit block event %d bloom %s header %d, ", len(bloom.Bytes()), encodedBloom, ctx.BlockHeight())
	k.Logger(ctx).Info(sprintf)

	ctx.EventManager().EmitEvent(
		cosmos.NewEvent(
			types.EventTypeBlockBloom,
			cosmos.NewAttribute(types.AttributeKeyEthereumBloom, encodedBloom),
		),
	)
}

// GetBlockBloomTransient returns bloom bytes for the current block height
func (k Keeper) GetBlockBloomTransient(ctx cosmos.Context) *big.Int {
	store := prefix.NewStore(ctx.TransientStore(k.transientKey), types.KeyPrefixTransientBloom)
	heightBz := cosmos.Uint64ToBigEndian(uint64(ctx.BlockHeight()))
	bz := store.Get(heightBz)
	if len(bz) == 0 {
		return big.NewInt(0)
	}

	return new(big.Int).SetBytes(bz)
}

// SetBlockBloomTransient sets the given bloom bytes to the transient store. This value is reset on
// every block.
func (k Keeper) SetBlockBloomTransient(ctx cosmos.Context, bloom *big.Int) {
	store := prefix.NewStore(ctx.TransientStore(k.transientKey), types.KeyPrefixTransientBloom)
	heightBz := cosmos.Uint64ToBigEndian(uint64(ctx.BlockHeight()))
	store.Set(heightBz, bloom.Bytes())

	k.Logger(ctx).Debug(
		"setState: SetBlockBloomTransient",
		"block-height", ctx.BlockHeight(),
		"bloom", bloom.String(),
	)
}

// ----------------------------------------------------------------------------
// 								  Tx Index
// ----------------------------------------------------------------------------

// SetTxIndexTransient set the index of processing txs
func (k Keeper) SetTxIndexTransient(ctx cosmos.Context, index uint64) {
	store := ctx.TransientStore(k.transientKey)
	store.Set(types.KeyPrefixTransientTxIndex, cosmos.Uint64ToBigEndian(index))

	k.Logger(ctx).Debug(
		"setState: SetTxIndexTransient",
		"key", "KeyPrefixTransientTxIndex",
		"index", index,
	)
}

// GetTxIndexTransient returns EVM txs index on the current block.
func (k Keeper) GetTxIndexTransient(ctx cosmos.Context) uint64 {
	store := ctx.TransientStore(k.transientKey)
	bz := store.Get(types.KeyPrefixTransientTxIndex)
	if len(bz) == 0 {
		return 0
	}

	return cosmos.BigEndianToUint64(bz)
}

// ----------------------------------------------------------------------------
// 									Log
// ----------------------------------------------------------------------------

// GetLogSizeTransient returns EVM log index on the current block.
func (k Keeper) GetLogSizeTransient(ctx cosmos.Context) uint64 {
	store := ctx.TransientStore(k.transientKey)
	bz := store.Get(types.KeyPrefixTransientLogSize)
	if len(bz) == 0 {
		return 0
	}

	return cosmos.BigEndianToUint64(bz)
}

// SetLogSizeTransient fetches the current EVM log index from the transient store, increases its
// value by one and then sets the new index back to the transient store.
func (k Keeper) SetLogSizeTransient(ctx cosmos.Context, logSize uint64) {
	store := ctx.TransientStore(k.transientKey)
	store.Set(types.KeyPrefixTransientLogSize, cosmos.Uint64ToBigEndian(logSize))

	k.Logger(ctx).Debug(
		"setState: SetLogSizeTransient",
		"key", "KeyPrefixTransientLogSize",
		"logSize", logSize,
	)
}

// ----------------------------------------------------------------------------
// 									Storage
// ----------------------------------------------------------------------------

// GetAccountStorage return states storage associated with an account
func (k Keeper) GetAccountStorage(ctx cosmos.Context, address common.Address) support.Storage {
	storage := support.Storage{}

	k.ForEachStorage(ctx, address, func(key, value common.Hash) bool {
		storage = append(storage, support.NewState(key, value))
		return true
	})

	return storage
}

// ----------------------------------------------------------------------------
//									Account
// ----------------------------------------------------------------------------

// GetAccountWithoutBalance load nonce and codeHash without balance,
// more efficient in cases where balance is not needed.
func (k *Keeper) GetAccountWithoutBalance(ctx cosmos.Context, addr common.Address) *states.StateAccount {
	cosmosAddr := cosmos.AccAddress(addr.Bytes())
	acct := k.accountKeeper.GetAccount(ctx, cosmosAddr)
	if acct == nil {
		return nil
	}

	codeHash := txs.EmptyCodeHash
	ethAcct, ok := acct.(artela.EthAccountI)
	if ok {
		codeHash = ethAcct.GetCodeHash().Bytes()
	}

	return &states.StateAccount{
		Nonce:    acct.GetSequence(),
		CodeHash: codeHash,
	}
}

// GetAccountOrEmpty returns empty account if not exist, returns error if it's not `EthAccount`
func (k *Keeper) GetAccountOrEmpty(ctx cosmos.Context, addr common.Address) states.StateAccount {
	acct := k.GetAccount(ctx, addr)
	if acct != nil {
		return *acct
	}

	// empty account
	return states.StateAccount{
		Balance:  new(big.Int),
		CodeHash: txs.EmptyCodeHash,
	}
}

// GetNonce returns the sequence number of an account, returns 0 if not exists.
func (k *Keeper) GetNonce(ctx cosmos.Context, addr common.Address) uint64 {
	cosmosAddr := cosmos.AccAddress(addr.Bytes())
	acct := k.accountKeeper.GetAccount(ctx, cosmosAddr)
	if acct == nil {
		return 0
	}

	return acct.GetSequence()
}

// GetBalance load account's balance of gas token
func (k *Keeper) GetBalance(ctx cosmos.Context, addr common.Address) *big.Int {
	cosmosAddr := cosmos.AccAddress(addr.Bytes())
	evmParams := k.GetParams(ctx)
	evmDenom := evmParams.GetEvmDenom()
	// if node is pruned, params is empty. Return invalid value
	if evmDenom == "" {
		return big.NewInt(-1)
	}
	coin := k.bankKeeper.GetBalance(ctx, cosmosAddr, evmDenom)
	return coin.Amount.BigInt()
}

// ----------------------------------------------------------------------------
// 								Gas and Fee
// ----------------------------------------------------------------------------

// Tracer return a default vm.Tracer based on current keeper states
func (k Keeper) Tracer(ctx cosmos.Context, msg *core.Message, ethCfg *params.ChainConfig) vm.EVMLogger {
	return txs.NewTracer(k.tracer, msg, ethCfg, ctx.BlockHeight())
}

// GetBaseFee returns current base fee, return values:
// - `nil`: london hardfork not enabled.
// - `0`: london hardfork enabled but fee is not enabled.
// - `n`: both london hardfork and fee are enabled.
func (k Keeper) GetBaseFee(ctx cosmos.Context, ethCfg *params.ChainConfig) *big.Int {
	return k.getBaseFee(ctx, support.IsLondon(ethCfg, ctx.BlockHeight()))
}

func (k Keeper) getBaseFee(ctx cosmos.Context, london bool) *big.Int {
	if !london {
		return nil
	}
	baseFee := k.feeKeeper.GetBaseFee(ctx)
	if baseFee == nil {
		// return 0 if fee not enabled.
		baseFee = big.NewInt(0)
	}
	return baseFee
}

// GetMinGasMultiplier returns the MinGasMultiplier param from the fee market module
func (k Keeper) GetMinGasMultiplier(ctx cosmos.Context) cosmos.Dec {
	feeParams := k.feeKeeper.GetParams(ctx)
	if feeParams.MinGasMultiplier.IsNil() {
		// in case we are executing eth_call on a legacy block, returns a zero value.
		return cosmos.ZeroDec()
	}
	return feeParams.MinGasMultiplier
}

// ResetTransientGasUsed reset gas used to prepare for execution of current cosmos txs, called in ante handler.
func (k Keeper) ResetTransientGasUsed(ctx cosmos.Context) {
	store := ctx.TransientStore(k.transientKey)
	store.Delete(types.KeyPrefixTransientGasUsed)

	k.Logger(ctx).Debug("setState: ResetTransientGasUsed, delete", "key", "KeyPrefixTransientGasUsed")
}

// GetTransientGasUsed returns the gas used by current cosmos txs.
func (k Keeper) GetTransientGasUsed(ctx cosmos.Context) uint64 {
	store := ctx.TransientStore(k.transientKey)
	bz := store.Get(types.KeyPrefixTransientGasUsed)
	if len(bz) == 0 {
		return 0
	}
	return cosmos.BigEndianToUint64(bz)
}

// SetTransientGasUsed sets the gas used by current cosmos txs.
func (k Keeper) SetTransientGasUsed(ctx cosmos.Context, gasUsed uint64) {
	store := ctx.TransientStore(k.transientKey)
	bz := cosmos.Uint64ToBigEndian(gasUsed)
	store.Set(types.KeyPrefixTransientGasUsed, bz)

	k.Logger(ctx).Debug(
		"setState: SetTransientGasUsed, set",
		"key", "KeyPrefixTransientGasUsed",
		"gasUsed", fmt.Sprintf("%d", gasUsed),
	)
}

// AddTransientGasUsed accumulate gas used by each eth msg included in current cosmos txs.
func (k Keeper) AddTransientGasUsed(ctx cosmos.Context, gasUsed uint64) (uint64, error) {
	result := k.GetTransientGasUsed(ctx) + gasUsed
	if result < gasUsed {
		return 0, errorsmod.Wrap(types.ErrGasOverflow, "transient gas used")
	}

	k.SetTransientGasUsed(ctx, result)
	return result, nil
}
