package types

import (
	"fmt"
	evmtypes "github.com/artela-network/artela/x/evm/types"
	artelatypes "github.com/artela-network/aspect-core/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"sync"

	"github.com/artela-network/artela-evm/vm"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	statedb "github.com/artela-network/artela/x/evm/states"
)

type HistoryStoreBuilder func(height int64, keyPrefix string) (prefix.Store, error)
type ContextBuilder func(height int64, prove bool) (cosmos.Context, error)

type GetLastBlockHeight func() int64

type AspectRuntimeContext struct {
	ethTxContext    *EthTxContext
	aspectContext   *AspectContext
	extBlockContext *ExtBlockContext
	aspectState     *AspectState
	cosmosCtx       cosmos.Context
	storeKey        storetypes.StoreKey
}

func NewAspectRuntimeContext(storeKey storetypes.StoreKey) *AspectRuntimeContext {
	return &AspectRuntimeContext{
		ethTxContext:    nil,
		cosmosCtx:       cosmos.Context{},
		aspectContext:   NewAspectContext(),
		extBlockContext: NewExtBlockContext(),
		aspectState:     NewAspectState(),
		storeKey:        storeKey,
	}
}

func (c *AspectRuntimeContext) WithCosmosContext(newTxCtx cosmos.Context) {
	c.cosmosCtx = newTxCtx
}
func (c *AspectRuntimeContext) CosmosContext() cosmos.Context {
	return c.cosmosCtx
}
func (c *AspectRuntimeContext) StoreKey() storetypes.StoreKey {
	return c.storeKey
}

func (c *AspectRuntimeContext) SetEthTxContext(newTxCtx *EthTxContext) {
	c.ethTxContext = newTxCtx
	c.aspectContext = NewAspectContext()
}

func (c *AspectRuntimeContext) ExtBlockContext() *ExtBlockContext {
	return c.extBlockContext
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

func (c *AspectRuntimeContext) StateDb() vm.StateDB {
	if c.EthTxContext() == nil {
		return nil
	}
	return c.EthTxContext().stateDb
}

func (c *AspectRuntimeContext) ClearBlockContext() {
	if c.extBlockContext != nil {
		c.extBlockContext = NewExtBlockContext()
	}
}

func (c *AspectRuntimeContext) ClearContext() {
	if c.EthTxContext().TxTo() == "" {
		c.ethTxContext = nil
		return
	}
	contractAddress := c.EthTxContext().TxTo()
	c.AspectContext().Clear(contractAddress)
	c.ethTxContext = nil
}

type EthTxContext struct {
	// eth Transaction,it is set in
	txContent     *ethtypes.Transaction
	msg           *core.Message
	vmTracer      *vm.Tracer
	receipt       *ethtypes.Receipt
	stateDb       vm.StateDB
	evmCfg        *statedb.EVMConfig
	lastEvm       *vm.EVM
	extProperties map[string]interface{}
	from          common.Address
}

func NewEthTxContext(ethTx *ethtypes.Transaction) *EthTxContext {
	return &EthTxContext{
		txContent:     ethTx,
		vmTracer:      nil,
		receipt:       nil,
		stateDb:       nil,
		extProperties: make(map[string]interface{}),
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

func (c *EthTxContext) EvmCfg() *statedb.EVMConfig            { return c.evmCfg }
func (c *EthTxContext) TxContent() *ethtypes.Transaction      { return c.txContent }
func (c *EthTxContext) VmTracer() *vm.Tracer                  { return c.vmTracer }
func (c *EthTxContext) Receipt() *ethtypes.Receipt            { return c.receipt }
func (c *EthTxContext) VmStateDB() vm.StateDB                 { return c.stateDb }
func (c *EthTxContext) ExtProperties() map[string]interface{} { return c.extProperties }
func (c *EthTxContext) LastEvm() *vm.EVM                      { return c.lastEvm }
func (c *EthTxContext) Message() *core.Message                { return c.msg }

func (c *EthTxContext) WithFrom(from common.Address) *EthTxContext {
	c.from = from
	return c
}

func (c *EthTxContext) WithMessage(msg *core.Message) *EthTxContext {
	c.msg = msg
	return c
}

func (c *EthTxContext) WithLastEvm(lastEvm *vm.EVM) *EthTxContext {
	c.lastEvm = lastEvm
	return c
}

func (c *EthTxContext) WithEvmCfg(cfg *statedb.EVMConfig) *EthTxContext {
	c.evmCfg = cfg
	return c
}

func (c *EthTxContext) WithVmMonitor(monitor *vm.Tracer) *EthTxContext {
	c.vmTracer = monitor
	return c
}

func (c *EthTxContext) WithStateDB(db vm.StateDB) *EthTxContext {
	c.stateDb = db
	return c
}

func (c *EthTxContext) WithReceipt(receipt *ethtypes.Receipt) *EthTxContext {
	c.receipt = receipt
	return c
}

func (c *EthTxContext) AddExtProperties(key string, value interface{}) *EthTxContext {
	if value == nil || key == "" {
		return c
	}
	c.extProperties[key] = value
	return c
}

func (c *EthTxContext) ClearEvmObject() *EthTxContext {
	c.stateDb = nil
	c.vmTracer = nil
	c.lastEvm = nil
	c.evmCfg = nil

	return c
}

type AspectContext struct {
	// 1.string= namespace Default
	// 2.string=key
	// 3.string=value
	context map[string]map[string]string
	mutex   sync.RWMutex
}

func NewAspectContext() *AspectContext {
	return &AspectContext{
		context: make(map[string]map[string]string),
	}
}

func (c *AspectContext) Add(aspectId string, key string, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.context[aspectId] == nil {
		c.context[aspectId] = make(map[string]string, 1)
	}
	c.context[aspectId][key] = value
}

func (c *AspectContext) Get(aspectId string, key string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.context[aspectId] == nil {
		return ""
	}
	return c.context[aspectId][key]
}

func (c *AspectContext) Remove(aspectId string, key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if aspectId == "" && key == "" {
		return false
	}
	if aspectId != "" && key == "" {
		delete(c.context, aspectId)
	}
	if aspectId != "" && key != "" {
		delete(c.context[aspectId], key)
	}
	return true
}

func (c *AspectContext) Clear(contractAddr string) {
	delete(c.context, contractAddr)
}

type ExtBlockContext struct {
	evidenceList   []abci.Misbehavior
	lastCommitInfo abci.CommitInfo
	getRpcClient   client.Context
}

func NewExtBlockContext() *ExtBlockContext {
	return &ExtBlockContext{}
}

func (c *ExtBlockContext) WithEvidenceList(cfg []abci.Misbehavior) *ExtBlockContext {
	c.evidenceList = cfg
	return c
}

func (c *ExtBlockContext) WithLastCommit(cfg abci.CommitInfo) *ExtBlockContext {
	c.lastCommitInfo = cfg
	return c
}

func (c *ExtBlockContext) WithRpcClient(cfg client.Context) *ExtBlockContext {
	c.getRpcClient = cfg
	return c
}

func (c *ExtBlockContext) EvidenceList() []abci.Misbehavior {
	return c.evidenceList
}

func (c *ExtBlockContext) LastCommitInfo() abci.CommitInfo {
	return c.lastCommitInfo
}

func (c *ExtBlockContext) RpcClient() client.Context {
	return c.getRpcClient
}

type AspectState struct {
	//int64 block height
	//string  group
	//  AspectStateObject
	stateCache map[int64]map[string]*AspectStateObject
}

func NewAspectState() *AspectState {
	return &AspectState{
		stateCache: make(map[int64]map[string]*AspectStateObject),
	}
}

func (k *AspectRuntimeContext) CreateStateObject(temporary bool, blockHeight int64, lockKey string) {
	object := NewAspectStateObject(k.cosmosCtx, k.storeKey, AspectStateKeyPrefix, temporary)
	m := k.aspectState.stateCache[blockHeight]
	if m == nil {
		k.aspectState.stateCache[blockHeight] = make(map[string]*AspectStateObject)
	}
	k.aspectState.stateCache[blockHeight][lockKey] = object
}

func (k *AspectRuntimeContext) RefreshState(needCommit bool, blockHeight int64, lockKey string) {
	if blockHeight < 0 {
		return
	}
	if len(lockKey) == 0 {
		if mapResult, ok := k.aspectState.stateCache[blockHeight]; ok {
			if needCommit {
				for _, object := range mapResult {
					object.commit()
				}
			}
			delete(k.aspectState.stateCache, blockHeight)
		}
		return
	}
	if stateObject, exist := k.aspectState.stateCache[blockHeight][lockKey]; exist {
		if needCommit {
			stateObject.commit()
		}
		delete(k.aspectState.stateCache[blockHeight], lockKey)
	}
}

func (k *AspectRuntimeContext) GetAspectState(ctx *artelatypes.RunnerContext, key string) string {
	aspectPropertyKey := AspectArrayKey(
		ctx.AspectId.Bytes(),
		[]byte(key),
	)
	point := GetAspectStatePoint(ctx.Point)
	var get []byte
	if len(point) == 0 {
		// If it's not a cut-point access to the DeliverTx stage, the data is accessed through the store at the previous block height
		store := prefix.NewStore(k.cosmosCtx.KVStore(k.storeKey), evmtypes.KeyPrefix(AspectStateKeyPrefix))
		get = store.Get(aspectPropertyKey)
	}
	if object, exist := k.aspectState.stateCache[ctx.BlockNumber][point]; exist {
		get = object.Get(aspectPropertyKey)
	}
	return artelatypes.Ternary(get != nil, func() string {
		return string(get)
	}, "")
}

func (k *AspectRuntimeContext) SetAspectState(ctx *artelatypes.RunnerContext, key, value string) bool {
	point := GetAspectStatePoint(ctx.Point)
	if len(point) == 0 {
		return false
	}
	aspectPropertyKey := AspectArrayKey(
		ctx.AspectId.Bytes(),
		[]byte(key),
	)
	if object, exist := k.aspectState.stateCache[ctx.BlockNumber][point]; exist {
		object.Set(aspectPropertyKey, []byte(value))
		return true
	}
	return false
}

// RemoveAspectState RemoveAspectState( key string) bool
func (k *AspectRuntimeContext) RemoveAspectState(ctx *artelatypes.RunnerContext, key string) bool {
	point := GetAspectStatePoint(ctx.Point)
	if len(point) == 0 {
		return false
	}
	aspectPropertyKey := AspectArrayKey(
		ctx.AspectId.Bytes(),
		[]byte(key),
	)
	if object, exist := k.aspectState.stateCache[ctx.BlockNumber][point]; exist {
		object.Set(aspectPropertyKey, nil)
		return true
	}
	return false
}

type AspectStateObject struct {
	preStore prefix.Store
	commit   func()
	storeKey storetypes.StoreKey

	log log.Logger
}

func NewAspectStateObject(ctx cosmos.Context, storeKey storetypes.StoreKey, fixKey string, temporary bool) *AspectStateObject {
	store := prefix.NewStore(ctx.KVStore(storeKey), evmtypes.KeyPrefix(fixKey))
	stateObj := &AspectStateObject{
		preStore: store,
		commit:   nil,
		storeKey: storeKey,
		log:      ctx.Logger(),
	}
	if temporary {
		cc, writeEvent := ctx.CacheContext()
		cacheStore := prefix.NewStore(cc.KVStore(storeKey), evmtypes.KeyPrefix(fixKey))
		stateObj.commit = writeEvent
		stateObj.preStore = cacheStore
	}
	return stateObj
}

func (k *AspectStateObject) Set(key, value []byte) {
	action := "updated"
	if len(value) == 0 {
		k.preStore.Delete(key)
		action = "deleted"
	} else {
		k.preStore.Set(key, value)
	}
	k.log.Info(
		fmt.Sprintf("states %s", action),
		"key", common.Bytes2Hex(key),
		"value", common.Bytes2Hex(key),
	)
}

func (k *AspectStateObject) Get(key []byte) []byte {
	return k.preStore.Get(key)
}

func (k *AspectStateObject) Commit() {
	if k.commit != nil {
		k.commit()
	}
}
