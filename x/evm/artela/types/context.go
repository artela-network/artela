package types

import (
	"context"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"math/big"
	"sync"
	"time"

	evmtypes "github.com/artela-network/artela/x/evm/types"
	artelatypes "github.com/artela-network/aspect-core/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"

	"github.com/artela-network/artela-evm/vm"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	statedb "github.com/artela-network/artela/x/evm/states"
	inherent "github.com/artela-network/aspect-core/chaincoreext/jit_inherent"
)

const (
	AspectContextKey   cosmos.ContextKey = "aspect-ctx"
	ExtBlockContextKey cosmos.ContextKey = "block-ctx"

	AspectModuleName = "aspect"
)

var cachedStoreKey storetypes.StoreKey

type (
	HistoryStoreBuilder func(height int64, keyPrefix string) (prefix.Store, error)
	ContextBuilder      func(height int64, prove bool) (cosmos.Context, error)

	GetLastBlockHeight func() int64
)

// AspectRuntimeContext is the contextual object required for Aspect execution,
// containing information related to transactions (tx) and blocks. Aspects at different
// join points can access this context, and consequently, the context dynamically
// adjusts its content based on the actual execution of blocks and transactions.

// Here is the execution scenario of this context in the lifecycle of a tx process,
// listed in the order of tx execution:

// 1. initialization: Before each transaction execution, create the AspectRuntimeContext
// and establish a bidirectional connection with the sdk context.

// 2. withBlockConfig: Write information before the start of each block and destroy it
// at the end of each block. Transfer it to the AspectRuntimeContext before the execution
// of tx in the deliver state through WithExtBlock.

// 3. withEVM: Before Pre-tx-execute, incorporate the EVM context, including evm, stateDB,
// evm tracer, message, message from, etc., and pass it to the AspectRuntimeContext through
// WithTxContext.

// 4. withReceipt: After the execution of the EVM, store the result in TxContext, enabling
// subsequent JoinPoints to access the execution details of the tx.

// 5. commit: Decide whether to commit at the end of each transaction. If committing is
// necessary, write the result to the sdk context.

// 6. destory: After each transaction execution, destroy the AspectRuntimeContext.
type AspectRuntimeContext struct {
	baseCtx context.Context

	ethTxContext    *EthTxContext
	aspectContext   *AspectContext
	ethBlockContext *EthBlockContext
	aspectState     *AspectState
	cosmosCtx       cosmos.Context
	storeKey        storetypes.StoreKey

	logger     log.Logger
	jitManager *inherent.Manager
}

func NewAspectRuntimeContext() *AspectRuntimeContext {
	return &AspectRuntimeContext{
		cosmosCtx:     cosmos.Context{},
		aspectContext: NewAspectContext(),
		storeKey:      cachedStoreKey,
		logger:        log.NewNopLogger(),
	}
}

func (c *AspectRuntimeContext) Init(storeKey storetypes.StoreKey) {
	cachedStoreKey = storeKey
	c.storeKey = storeKey
}

func (c *AspectRuntimeContext) WithCosmosContext(newTxCtx cosmos.Context) {
	c.cosmosCtx = newTxCtx
	c.logger = newTxCtx.Logger().With("module", fmt.Sprintf("x/%s", AspectModuleName))
}

func (c *AspectRuntimeContext) Debug(msg string, keyvals ...interface{}) {
	if c.ethTxContext != nil {
		keyvals = append(keyvals, "tx-from", fmt.Sprintf("%s", c.ethTxContext.TxFrom().Hex()))
		if c.ethTxContext.TxContent() != nil {
			keyvals = append(keyvals, "tx-hash", fmt.Sprintf("%s", c.ethTxContext.TxContent().Hash().Hex()))
		}
	}
	c.logger.Debug(msg, keyvals...)
}

func (c *AspectRuntimeContext) Logger() log.Logger {
	return c.logger
}

func (c *AspectRuntimeContext) CosmosContext() cosmos.Context {
	return c.cosmosCtx
}

func (c *AspectRuntimeContext) StoreKey() storetypes.StoreKey {
	return c.storeKey
}

func (c *AspectRuntimeContext) SetEthTxContext(newTxCtx *EthTxContext, jitManager *inherent.Manager) {
	c.ethTxContext = newTxCtx
	c.aspectContext = NewAspectContext()
	c.jitManager = jitManager
}

func (c *AspectRuntimeContext) SetEthBlockContext(newBlockCtx *EthBlockContext) {
	c.ethBlockContext = newBlockCtx
}

func (c *AspectRuntimeContext) EthBlockContext() *EthBlockContext {
	return c.ethBlockContext
}

func (c *AspectRuntimeContext) EthTxContext() *EthTxContext {
	return c.ethTxContext
}

func (c *AspectRuntimeContext) AspectContext() *AspectContext {
	return c.aspectContext
}

func (c *AspectRuntimeContext) AspectState() *AspectState {
	return c.aspectState
}

func (c *AspectRuntimeContext) JITManager() *inherent.Manager {
	return c.jitManager
}

func (c *AspectRuntimeContext) StateDb() vm.StateDB {
	if c.EthTxContext() == nil {
		return nil
	}
	return c.EthTxContext().stateDb
}

func (c *AspectRuntimeContext) ClearBlockContext() {
	if c.ethBlockContext != nil {
		c.ethBlockContext = nil
	}
}

func (c *AspectRuntimeContext) ClearContext() {
	if c.EthTxContext().TxTo() == "" {
		c.ethTxContext = nil
		return
	}
	contractAddress := c.EthTxContext().TxTo()
	c.AspectContext().Clear(common.HexToAddress(contractAddress))
	c.ethTxContext = nil
}

func (c *AspectRuntimeContext) Destory() {
	if c.EthTxContext() != nil {
		c.EthTxContext().ClearEvmObject()
	}
	c.ClearContext()

	c.jitManager = nil
	c.aspectContext = nil
	c.cosmosCtx = cosmos.Context{}
	c.aspectState = nil
}

func (c *AspectRuntimeContext) Deadline() (deadline time.Time, ok bool) {
	return c.baseCtx.Deadline()
}

func (c *AspectRuntimeContext) Done() <-chan struct{} {
	return c.baseCtx.Done()
}

func (c *AspectRuntimeContext) Err() error {
	return c.baseCtx.Err()
}

func (c *AspectRuntimeContext) Value(key interface{}) interface{} {
	return c.baseCtx.Value(key)
}

type EthTxContext struct {
	// eth Transaction,it is set in
	txContent *ethtypes.Transaction
	msg       *core.Message
	vmTracer  *vm.Tracer
	receipt   *ethtypes.Receipt
	stateDb   vm.StateDB
	evmCfg    *statedb.EVMConfig
	lastEvm   *vm.EVM
	from      common.Address
	index     uint64
	commit    bool
}

func NewEthTxContext(ethTx *ethtypes.Transaction) *EthTxContext {
	return &EthTxContext{
		txContent: ethTx,
		vmTracer:  nil,
		receipt:   nil,
		stateDb:   nil,
	}
}

func (c *EthTxContext) TxTo() string {
	if c.txContent == nil {
		return ""
	}
	if c.txContent.To() == nil {
		return ""
	}
	return c.txContent.To().String()
}

func (c *EthTxContext) TxFrom() common.Address {
	return c.from
}
func (c *EthTxContext) TxIndex() uint64 {
	return c.index
}
func (c *EthTxContext) EvmCfg() *statedb.EVMConfig       { return c.evmCfg }
func (c *EthTxContext) TxContent() *ethtypes.Transaction { return c.txContent }
func (c *EthTxContext) VmTracer() *vm.Tracer             { return c.vmTracer }
func (c *EthTxContext) Receipt() *ethtypes.Receipt       { return c.receipt }
func (c *EthTxContext) VmStateDB() vm.StateDB            { return c.stateDb }
func (c *EthTxContext) LastEvm() *vm.EVM                 { return c.lastEvm }
func (c *EthTxContext) Message() *core.Message           { return c.msg }
func (c *EthTxContext) Commit() bool                     { return c.commit }

func (c *EthTxContext) WithEVM(
	from common.Address,
	msg *core.Message,
	lastEvm *vm.EVM,
	monitor *vm.Tracer,
	db vm.StateDB,
) *EthTxContext {
	c.from = from
	c.msg = msg
	c.lastEvm = lastEvm
	c.vmTracer = monitor
	c.stateDb = db
	return c
}

func (c *EthTxContext) WithEVMConfig(cfg *statedb.EVMConfig) *EthTxContext {
	c.evmCfg = cfg
	return c
}

func (c *EthTxContext) WithTxIndex(index uint64) *EthTxContext {
	c.index = index
	return c
}

func (c *EthTxContext) WithReceipt(receipt *ethtypes.Receipt) *EthTxContext {
	c.receipt = receipt
	return c
}

func (c *EthTxContext) WithCommit(commit bool) *EthTxContext {
	c.commit = commit
	return c
}

func (c *EthTxContext) WithStateDB(stateDb vm.StateDB) *EthTxContext {
	c.stateDb = stateDb
	return c
}

func (c *EthTxContext) ClearEvmObject() *EthTxContext {
	c.stateDb = nil
	c.vmTracer = nil
	c.lastEvm = nil
	c.evmCfg = nil

	return c
}

type EthBlockContext struct {
	blockHeader *ethtypes.Header
}

func NewEthBlockContextFromHeight(height int64) *EthBlockContext {
	return &EthBlockContext{&ethtypes.Header{Number: big.NewInt(height)}}
}

func NewEthBlockContextFromABCIBeginBlockReq(req abci.RequestBeginBlock) *EthBlockContext {
	txHash := ethtypes.EmptyTxsHash
	if len(req.Header.DataHash) != 0 {
		txHash = common.BytesToHash(req.Header.DataHash)
	}

	blockHeader := &ethtypes.Header{
		ParentHash: common.BytesToHash(req.Header.LastBlockId.Hash),
		Coinbase:   common.BytesToAddress(req.Header.ProposerAddress),
		TxHash:     txHash,
		Number:     big.NewInt(req.Header.Height),
		Time:       uint64(req.Header.Time.UTC().Unix()),
	}

	return &EthBlockContext{
		blockHeader: blockHeader,
	}
}

func NewEthBlockContextFromQuery(sdkCtx cosmos.Context, queryCtx client.Context) *EthBlockContext {
	blockHeight := sdkCtx.BlockHeight()
	resBlock, err := queryCtx.Client.Block(sdkCtx, &blockHeight)
	if err != nil || resBlock == nil || resBlock.Block == nil {
		return nil
	}

	resBlockHeader := resBlock.Block.Header

	txHash := ethtypes.EmptyTxsHash
	if len(resBlockHeader.DataHash) != 0 {
		txHash = common.BytesToHash(resBlockHeader.DataHash)
	}

	blockHeader := &ethtypes.Header{
		ParentHash: common.BytesToHash(resBlockHeader.LastBlockID.Hash),
		Coinbase:   common.BytesToAddress(resBlockHeader.ProposerAddress),
		TxHash:     txHash,
		Number:     big.NewInt(resBlockHeader.Height),
		Time:       uint64(resBlockHeader.Time.UTC().Unix()),
	}

	return &EthBlockContext{
		blockHeader: blockHeader,
	}
}

func (c *EthBlockContext) BlockHeader() *ethtypes.Header {
	return c.blockHeader
}

type AspectContext struct {
	// 1.string=namespace Default
	// 2.string=key
	// 3.string=value
	context map[common.Address]map[string][]byte
	mutex   sync.RWMutex
}

func NewAspectContext() *AspectContext {
	return &AspectContext{
		context: make(map[common.Address]map[string][]byte),
	}
}

func (c *AspectContext) Add(address common.Address, key string, value []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.context[address] == nil {
		c.context[address] = make(map[string][]byte, 1)
	}
	c.context[address][key] = value
}

func (c *AspectContext) Get(address common.Address, key string) []byte {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.context[address] == nil {
		return []byte{}
	}
	return c.context[address][key]
}

func (c *AspectContext) Clear(address common.Address) {
	delete(c.context, address)
}

func (k *AspectRuntimeContext) CreateStateObject() {
	k.aspectState = NewAspectState(k.cosmosCtx, k.storeKey, AspectStateKeyPrefix, k.logger)
}

func (k *AspectRuntimeContext) GetAspectState(ctx *artelatypes.RunnerContext, key string) []byte {
	stateKey := AspectArrayKey(
		ctx.AspectId.Bytes(),
		[]byte(key),
	)
	return k.aspectState.Get(stateKey)
}

func (k *AspectRuntimeContext) SetAspectState(ctx *artelatypes.RunnerContext, key string, value []byte) {
	stateKey := AspectArrayKey(
		ctx.AspectId.Bytes(),
		[]byte(key),
	)

	k.aspectState.Set(stateKey, value)
	return
}

type AspectState struct {
	state    prefix.Store
	storeKey storetypes.StoreKey

	logger log.Logger
}

func NewAspectState(ctx cosmos.Context, storeKey storetypes.StoreKey, fixKey string, logger log.Logger) *AspectState {
	cacheStore := prefix.NewStore(ctx.KVStore(storeKey), evmtypes.KeyPrefix(fixKey))
	stateObj := &AspectState{
		state:    cacheStore,
		storeKey: storeKey,
		logger:   logger,
	}
	return stateObj
}

func (k *AspectState) Set(key, value []byte) {
	action := "updated"
	if len(value) == 0 {
		k.state.Delete(key)
		action = "deleted"
	} else {
		k.state.Set(key, value)
	}

	if value == nil {
		k.logger.Debug("setState:", "action", action, "key", string(key), "value", "nil")
	} else {
		k.logger.Debug("setState:", "action", action, "key", string(key), "value", string(value))
	}

}

func (k *AspectState) Get(key []byte) []byte {
	val := k.state.Get(key)
	if val == nil {
		k.logger.Debug("getState:", "key", string(key), "value", "nil")
	} else {
		k.logger.Debug("getState:", "key", string(key), "value", string(val))
	}
	return val
}
