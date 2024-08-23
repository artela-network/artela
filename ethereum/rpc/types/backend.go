package types

import (
	"context"
	"math/big"

	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
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
	"github.com/artela-network/artela/x/evm/txs"
	evmtxs "github.com/artela-network/artela/x/evm/txs"
	evmsupport "github.com/artela-network/artela/x/evm/txs/support"
)

type (
	// Backend defines the common interfaces
	Backend interface {
		CurrentHeader() (*types.Header, error)

		Accounts() []common.Address
		GetBalance(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error)
		ArtBlockByNumber(ctx context.Context, number rpc.BlockNumber) (*Block, error)
		BlockByHash(ctx context.Context, hash common.Hash) (*Block, error)
		ChainConfig() *params.ChainConfig
	}

	// EthereumBackend defines the chain related interfaces
	EthereumBackend interface {
		Backend

		SuggestGasTipCap(baseFee *big.Int) (*big.Int, error)
		GasPrice(ctx context.Context) (*hexutil.Big, error)
		FeeHistory(blockCount uint64, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*FeeHistoryResult, error)

		Engine() consensus.Engine
		Syncing() (interface{}, error)
	}

	// BlockChainBackend defines the block chain interfaces
	BlockChainBackend interface {
		Backend

		GetProof(address common.Address, storageKeys []string, blockNrOrHash BlockNumberOrHash) (*AccountResult, error)
		DoCall(args TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash) (*txs.MsgEthereumTxResponse, error)
		EstimateGas(ctx context.Context, args TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (hexutil.Uint64, error)

		HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error)
		HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
		HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error)
		CurrentBlock() *Block
		ArtBlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*Block, error)
		CosmosBlockByHash(blockHash common.Hash) (*tmrpctypes.ResultBlock, error)
		CosmosBlockByNumber(blockNum rpc.BlockNumber) (*tmrpctypes.ResultBlock, error)
		StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *types.Header, error)
		StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error)
		GetEVM(ctx context.Context, msg *core.Message, state *state.StateDB, header *types.Header, vmConfig *vm.Config, blockCtx *vm.BlockContext) (*vm.EVM, func() error)
		GetCode(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error)
		GetStorageAt(address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error)
		GetCoinbase() (sdk.AccAddress, error)
	}

	// TrancsactionBackend defines the block chain interfaces
	TrancsactionBackend interface {
		BlockChainBackend
		EthereumBackend

		SendTx(ctx context.Context, signedTx *types.Transaction) error
		GetTransaction(ctx context.Context, txHash common.Hash) (*RPCTransaction, error)
		GetTransactionCount(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error)
		GetTxMsg(ctx context.Context, txHash common.Hash) (*txs.MsgEthereumTx, error)
		SignTransaction(args *TransactionArgs) (*types.Transaction, error)
		GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error)
		RPCTxFeeCap() float64
		UnprotectedAllowed() bool

		PendingTransactions() ([]*sdk.Tx, error)
		GetResendArgs(args TransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (TransactionArgs, error)
		Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error)
		GetSender(msg *txs.MsgEthereumTx, chainID *big.Int) (from common.Address, err error)
	}

	DebugBackend interface {
		BlockChainBackend
		TrancsactionBackend

		TraceTransaction(hash common.Hash, config *evmsupport.TraceConfig) (interface{}, error)
		TraceBlock(height rpc.BlockNumber,
			config *evmsupport.TraceConfig,
			block *tmrpctypes.ResultBlock,
		) ([]*evmtxs.TxTraceResult, error)
		GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error)

		DBProperty(property string) (string, error)
		DBCompact(start []byte, limit []byte) error
	}

	PersonalBackend interface {
		TrancsactionBackend

		NewAccount(password string) (common.AddressEIP55, error)
		ImportRawKey(privkey, password string) (common.Address, error)
	}

	TxPoolBackend interface {
		TrancsactionBackend

		PendingTransactionsCount() (int, error)
	}

	// NetBackend is the collection of methods required to satisfy the net
	// RPC DebugAPI.
	NetBackend interface {
		PeerCount() hexutil.Uint
		Listening() bool
		Version() string
	}

	// Web3Backend is the collection of methods required to satisfy the net
	// RPC DebugAPI.
	Web3Backend interface {
		ClientVersion() string
	}
)
