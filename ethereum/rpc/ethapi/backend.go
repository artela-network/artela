package ethapi

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/artela-evm/vm"
	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	"github.com/artela-network/artela/x/evm/txs"
)

// Backend interface provides the common API services (that are provided by
// both full and light clients) with access to necessary functions.
type Backend interface {
	// General Ethereum API
	SuggestGasTipCap(baseFee *big.Int) (*big.Int, error)
	GasPrice(ctx context.Context) (*hexutil.Big, error)

	// Account releated
	Accounts() []common.Address
	NewAccount(password string) (common.AddressEIP55, error)
	ImportRawKey(privkey, password string) (common.Address, error)
	GetTransactionCount(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error)
	GetBalance(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error)

	// Blockchain API
	HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error)
	CurrentHeader() (*types.Header, error)
	CurrentBlock() *rpctypes.Block
	BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error)
	ArtBlockByNumber(ctx context.Context, number rpc.BlockNumber) (*rpctypes.Block, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*rpctypes.Block, error)
	BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*rpctypes.Block, error)
	StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *types.Header, error)
	StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error)
	GetEVM(ctx context.Context, msg *core.Message, state *state.StateDB, header *types.Header, vmConfig *vm.Config, blockCtx *vm.BlockContext) (*vm.EVM, func() error)
	GetCode(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error)
	GetStorageAt(address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error)
	FeeHistory(blockCount uint64, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*rpctypes.FeeHistoryResult, error)

	// Transaction pool API
	SendTx(ctx context.Context, signedTx *types.Transaction) error
	GetTransaction(ctx context.Context, txHash common.Hash) (*RPCTransaction, error)
	SignTransaction(args *TransactionArgs) (*types.Transaction, error)
	GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error)
	RPCTxFeeCap() float64
	UnprotectedAllowed() bool
	EstimateGas(ctx context.Context, args TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (hexutil.Uint64, error)
	DoCall(args TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash) (*txs.MsgEthereumTxResponse, error)

	ChainConfig() *params.ChainConfig
	Engine() consensus.Engine

	Syncing() (interface{}, error)
	// This is copied from filters.Backend
	Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error)

	GetCoinbase() (sdk.AccAddress, error)

	GetProof(address common.Address, storageKeys []string, blockNrOrHash rpctypes.BlockNumberOrHash) (*rpctypes.AccountResult, error)
}
