package types

import (
	statedb "github.com/artela-network/artela/x/evm/states"
	artelavm "github.com/artela-network/evm/vm"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"sync"
)

type AspectRuntimeContext struct {
	ethTxContext    *EthTxContext
	aspectContext   *AspectContext
	extBlockContext *ExtBlockContext
}

func NewAspectRuntimeContext() *AspectRuntimeContext {
	return &AspectRuntimeContext{
		ethTxContext:    nil,
		aspectContext:   NewAspectContext(),
		extBlockContext: NewExtBlockContext(),
	}
}

func (c *AspectRuntimeContext) SetEthTxContext(newTxCtx *EthTxContext) {
	c.ethTxContext = newTxCtx
}
func (c *AspectRuntimeContext) NewAspectContext() {
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
	vmTracer      *artelavm.Tracer
	receipt       *ethtypes.Receipt
	stateDb       vm.StateDB
	evmCfg        *statedb.EVMConfig
	lastEvm       *vm.EVM
	extProperties map[string]interface{}
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
func (c *EthTxContext) EvmCfg() *statedb.EVMConfig            { return c.evmCfg }
func (c *EthTxContext) TxContent() *ethtypes.Transaction      { return c.txContent }
func (c *EthTxContext) VmTracer() *artelavm.Tracer            { return c.vmTracer }
func (c *EthTxContext) Receipt() *ethtypes.Receipt            { return c.receipt }
func (c *EthTxContext) VmStateDB() artelavm.StateDB           { return c.stateDb }
func (c *EthTxContext) ExtProperties() map[string]interface{} { return c.extProperties }
func (c *EthTxContext) LastEvm() *vm.EVM                      { return c.lastEvm }

func (c *EthTxContext) WithLastEvm(lastEvm *vm.EVM) *EthTxContext {
	c.lastEvm = lastEvm
	return c
}

func (c *EthTxContext) WithEvmCfg(cfg *statedb.EVMConfig) *EthTxContext {
	c.evmCfg = cfg
	return c
}
func (c *EthTxContext) WithVmMonitor(monitor *artelavm.Tracer) *EthTxContext {
	c.vmTracer = monitor
	return c
}
func (c *EthTxContext) WithStateDB(db artelavm.StateDB) *EthTxContext {
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
