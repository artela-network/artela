package ethapi

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/davecgh/go-spew/spew"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/artela-evm/vm"
	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	ethtypes "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/txs"
)

// EthereumAPI provides an API to access Ethereum related information.
type EthereumAPI struct {
	b      Backend
	logger log.Logger
}

// NewEthereumAPI creates a new Ethereum protocol API.
func NewEthereumAPI(b Backend, logger log.Logger) *EthereumAPI {
	return &EthereumAPI{b, logger}
}

// ProtocolVersion returns the supported Ethereum protocol version.
func (s *EthereumAPI) ProtocolVersion() hexutil.Uint {
	s.logger.Debug("eth_protocolVersion")
	return hexutil.Uint(ethtypes.ProtocolVersion)
}

// GasPrice returns a suggestion for a gas price for legacy transactions.
func (s *EthereumAPI) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	s.logger.Debug("eth_gasPrice")
	return s.b.GasPrice(ctx)
}

// MaxPriorityFeePerGas returns a suggestion for a gas tip cap for dynamic fee transactions.
func (s *EthereumAPI) MaxPriorityFeePerGas(_ context.Context) (*hexutil.Big, error) {
	return nil, errors.New("MaxPriorityFeePerGas is not implemented")
}

// FeeHistory returns the fee market history.
func (s *EthereumAPI) FeeHistory(_ context.Context, blockCount math.HexOrDecimal64, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*rpctypes.FeeHistoryResult, error) {
	return s.b.FeeHistory(uint64(blockCount), lastBlock, rewardPercentiles)
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up-to-date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronize from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of states entries processed until now
// - knownStates:   number of known states entries that still need to be pulled
func (s *EthereumAPI) Syncing() (interface{}, error) {
	s.logger.Debug("eth_syncing")
	return s.b.Syncing()
}

// TxPoolAPI offers and API for the transaction pool. It only operates on data that is non-confidential.
type TxPoolAPI struct {
	b Backend
}

// NewTxPoolAPI creates a new tx pool service that gives information about the transaction pool.
func NewTxPoolAPI(b Backend) *TxPoolAPI {
	return &TxPoolAPI{b}
}

// Content returns the transactions contained within the transaction pool.
func (s *TxPoolAPI) Content() map[string]map[string]map[string]*RPCTransaction {
	// not implemented
	return nil
}

// ContentFrom returns the transactions contained within the transaction pool.
func (s *TxPoolAPI) ContentFrom(_ common.Address) map[string]map[string]*RPCTransaction {
	// not implemented
	return nil
}

// Status returns the number of pending and queued transaction in the pool.
func (s *TxPoolAPI) Status() map[string]hexutil.Uint {
	// not implemented
	return nil
}

// Inspect retrieves the content of the transaction pool and flattens it into an
// easily inspectable list.
func (s *TxPoolAPI) Inspect() map[string]map[string]map[string]string {
	// not implemented
	return nil
}

// EthereumAccountAPI provides an API to access accounts managed by this node.
// It offers only methods that can retrieve accounts.
type EthereumAccountAPI struct {
	b Backend
}

// NewEthereumAccountAPI creates a new EthereumAccountAPI.
func NewEthereumAccountAPI(b Backend) *EthereumAccountAPI {
	return &EthereumAccountAPI{b}
}

// Accounts returns the collection of accounts this node manages.
func (s *EthereumAccountAPI) Accounts() []common.Address {
	return s.b.Accounts()
}

// PersonalAccountAPI provides an API to access accounts managed by this node.
// It offers methods to create, (un)lock en list accounts. Some methods accept
// passwords and are therefore considered private by default.
type PersonalAccountAPI struct {
	nonceLock *AddrLocker
	logger    log.Logger
	b         Backend
}

// NewPersonalAccountAPI create a new PersonalAccountAPI.
func NewPersonalAccountAPI(b Backend, logger log.Logger, nonceLock *AddrLocker) *PersonalAccountAPI {
	return &PersonalAccountAPI{
		nonceLock: nonceLock,
		logger:    logger,
		b:         b,
	}
}

// ListAccounts will return a list of addresses for accounts this node manages.
func (s *PersonalAccountAPI) ListAccounts() []common.Address {
	return s.b.Accounts()
}

// rawWallet is a JSON representation of an accounts.Wallet interface, with its
// data contents extracted into plain fields.
type rawWallet struct {
	URL      string             `json:"url"`
	Status   string             `json:"status"`
	Failure  string             `json:"failure,omitempty"`
	Accounts []accounts.Account `json:"accounts,omitempty"`
}

// ListWallets will return a list of wallets this node manages.
func (s *PersonalAccountAPI) ListWallets() []rawWallet {
	// not implemented
	wallets := make([]rawWallet, 0) // return [] instead of nil if empty
	return wallets
}

// OpenWallet initiates a hardware wallet opening procedure, establishing a USB
// connection and attempting to authenticate via the provided passphrase. Note,
// the method may return an extra challenge requiring a second open (e.g. the
// Trezor PIN matrix challenge).
func (s *PersonalAccountAPI) OpenWallet(_ string, _ *string) error {
	return errors.New("OpenWallet is not implemented")
}

// DeriveAccount requests an HD wallet to derive a new account, optionally pinning
// it for later reuse.
func (s *PersonalAccountAPI) DeriveAccount(url string, path string, pin *bool) (accounts.Account, error) {
	return accounts.Account{}, errors.New("DeriveAccount is not implemented")
}

// NewAccount will create a new account and returns the address for the new account.
func (s *PersonalAccountAPI) NewAccount(password string) (common.AddressEIP55, error) {
	return s.b.NewAccount(password)
}

// ImportRawKey stores the given hex encoded ECDSA key into the key directory,
// encrypting it with the passphrase.
func (s *PersonalAccountAPI) ImportRawKey(privkey string, password string) (common.Address, error) {
	return s.b.ImportRawKey(privkey, password)
}

// UnlockAccount will unlock the account associated with the given address with
// the given password for duration seconds. If duration is nil it will use a
// default of 300 seconds. It returns an indication if the account was unlocked.
func (s *PersonalAccountAPI) UnlockAccount(_ context.Context, _ common.Address, _ string, duration *uint64) (bool, error) {
	// not implemented
	return false, errors.New("UnlockAccount is not implemented")
}

// LockAccount will lock the account associated with the given address when it's unlocked.
func (s *PersonalAccountAPI) LockAccount(_ common.Address) bool {
	// not implemented"
	return false
}

// signTransaction sets defaults and signs the given transaction
// NOTE: the caller needs to ensure that the nonceLock is held, if applicable,
// and release it after the transaction has been submitted to the tx pool
func (s *PersonalAccountAPI) signTransaction(_ context.Context, args *TransactionArgs, passwd string) (*types.Transaction, error) {
	// return s.b.SignTransaction(args, passwd)
	// TODO
	return nil, fmt.Errorf("signTransaction is not implemented, args: %v, passwd: %s", args, passwd)
}

// SendTransaction will create a transaction from the given arguments and
// tries to sign it with the key associated with args.From. If the given
// passwd isn't able to decrypt the key it fails.
func (s *PersonalAccountAPI) SendTransaction(ctx context.Context, args TransactionArgs, passwd string) (common.Hash, error) {
	if args.Nonce == nil {
		// Hold the mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		s.nonceLock.LockAddr(args.from())
		defer s.nonceLock.UnlockAddr(args.from())
	}
	signed, err := s.signTransaction(ctx, &args, passwd)
	if err != nil {
		log.Warn("Failed transaction send attempt", "from", args.from(), "to", args.To, "value", args.Value.ToInt(), "err", err)
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.logger, s.b, signed)
}

// SignTransaction will create a transaction from the given arguments and
// tries to sign it with the key associated with args.From. If the given passwd isn't
// able to decrypt the key it fails. The transaction is returned in RLP-form, not broadcast
// to other nodes
func (s *PersonalAccountAPI) SignTransaction(_ context.Context, args TransactionArgs, passwd string) (*SignTransactionResult, error) {
	// TODO
	return nil, fmt.Errorf("SignTransaction is not implemented, args: %v, passwd: %s", args, passwd)
}

// Sign calculates an Ethereum ECDSA signature for:
// keccak256("\x19Ethereum Signed Message:\n" + len(message) + message))
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The key used to calculate the signature is decrypted with the given password.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_sign
func (s *PersonalAccountAPI) Sign(_ context.Context, _ hexutil.Bytes, _ common.Address, _ string) (hexutil.Bytes, error) {
	// TODO
	return nil, errors.New("Sign is not implemented")
}

// EcRecover returns the address for the account that was used to create the signature.
// Note, this function is compatible with eth_sign and personal_sign. As such it recovers
// the address of:
// hash = keccak256("\x19Ethereum Signed Message:\n"${message length}${message})
// addr = ecrecover(hash, signature)
//
// Note, the signature must conform to the secp256k1 curve R, S and V values, where
// the V value must be 27 or 28 for legacy reasons.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_ecRecover
func (s *PersonalAccountAPI) EcRecover(_ context.Context, data, sig hexutil.Bytes) (common.Address, error) {
	if len(sig) != crypto.SignatureLength {
		return common.Address{}, fmt.Errorf("signature must be %d bytes long", crypto.SignatureLength)
	}
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return common.Address{}, errors.New("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	rpk, err := crypto.SigToPub(accounts.TextHash(data), sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*rpk), nil
}

// InitializeWallet initializes a new wallet at the provided URL, by generating and returning a new private key.
func (s *PersonalAccountAPI) InitializeWallet(_ context.Context, _ string) (string, error) {
	return "", errors.New("InitializeWallet is not implemented")
}

// Unpair deletes a pairing between wallet and geth.
func (s *PersonalAccountAPI) Unpair(_ context.Context, _ string, _ string) error {
	return errors.New("unpair is not implemented")
}

// BlockChainAPI provides an API to access Ethereum blockchain data.
type BlockChainAPI struct {
	logger log.Logger
	b      Backend
}

// NewBlockChainAPI creates a new Ethereum blockchain API.
func NewBlockChainAPI(b Backend, logger log.Logger) *BlockChainAPI {
	return &BlockChainAPI{logger, b}
}

// ChainId is the EIP-155 replay-protection chain id for the current Ethereum chain config.
//
// Note, this method does not conform to EIP-695 because the configured chain ID is always
// returned, regardless of the current head block. We used to return an error when the chain
// wasn't synced up to a block where EIP-155 is enabled, but this behavior caused issues
// in CL clients.
func (s *BlockChainAPI) ChainId() *hexutil.Big {
	return (*hexutil.Big)(s.b.ChainConfig().ChainID)
}

func (s *BlockChainAPI) Coinbase() (sdktypes.AccAddress, error) {
	return s.b.GetCoinbase()
}

// BlockNumber returns the block number of the chain head.
func (s *BlockChainAPI) BlockNumber() hexutil.Uint64 {
	header, err := s.b.HeaderByNumber(context.Background(), rpc.LatestBlockNumber) // latest header should always be available
	if err != nil || header == nil {
		s.logger.Debug("get BlockNumber failed", err)
		return 0
	}
	return hexutil.Uint64(header.Number.Uint64())
}

// GetBalance returns the amount of wei for the given address in the states of the
// given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers are also allowed.
func (s *BlockChainAPI) GetBalance(_ context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	return s.b.GetBalance(address, blockNrOrHash)
}

// Result structs for GetProof
type AccountResult struct {
	Address      common.Address  `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      *hexutil.Big    `json:"balance"`
	CodeHash     common.Hash     `json:"codeHash"`
	Nonce        hexutil.Uint64  `json:"nonce"`
	StorageHash  common.Hash     `json:"storageHash"`
	StorageProof []StorageResult `json:"storageProof"`
}

type StorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

// GetProof returns the Merkle-proof for a given account and optionally some storage keys.
func (s *BlockChainAPI) GetProof(_ context.Context, address common.Address, storageKeys []string, blockNrOrHash rpctypes.BlockNumberOrHash) (*rpctypes.AccountResult, error) {
	s.logger.Debug("eth_getProof", "address", address.Hex(), "keys", storageKeys, "block number or hash", blockNrOrHash)
	return s.b.GetProof(address, storageKeys, blockNrOrHash)
}

// GetHeaderByNumber returns the requested canonical block header.
// * When blockNr is -1 the chain head is returned.
// * When blockNr is -2 the pending chain head is returned.
func (s *BlockChainAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	block, err := s.b.ArtBlockByNumber(ctx, number)
	if block != nil && block.Header() != nil && err == nil {
		response := s.rpcMarshalHeader(ctx, block.Header(), block.Hash())
		if number == rpc.PendingBlockNumber {
			// Pending header need to nil out a few fields
			for _, field := range []string{"hash", "nonce", "miner"} {
				response[field] = nil
			}
		}
		return response, err
	}
	return nil, err
}

// GetHeaderByHash returns the requested header by hash.
func (s *BlockChainAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) map[string]interface{} {
	block, _ := s.b.BlockByHash(ctx, hash)
	if block != nil && block.Header() != nil {
		return s.rpcMarshalHeader(ctx, block.Header(), block.Hash())
	}
	return nil
}

// GetBlockByNumber returns the requested canonical block.
//   - When blockNr is -1 the chain head is returned.
//   - When blockNr is -2 the pending chain head is returned.
//   - When fullTx is true all transactions in the block are returned, otherwise
//     only the transaction hash is returned.
func (s *BlockChainAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	block, err := s.b.ArtBlockByNumber(ctx, number)
	if block != nil && err == nil {
		response, err := s.rpcMarshalBlock(ctx, block, true, fullTx)
		if err == nil && number == rpc.PendingBlockNumber {
			// Pending blocks need to nil out a few fields
			for _, field := range []string{"hash", "nonce", "miner"} {
				response[field] = nil
			}
		}
		return response, err
	}
	return nil, err
}

// GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.
func (s *BlockChainAPI) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	block, err := s.b.BlockByHash(ctx, hash)
	if block != nil {
		return s.rpcMarshalBlock(ctx, block, true, fullTx)
	}
	return nil, err
}

// GetUncleByBlockNumberAndIndex returns the uncle block for the given block hash and index.
func (s *BlockChainAPI) GetUncleByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) (map[string]interface{}, error) {
	block, err := s.b.ArtBlockByNumber(ctx, blockNr)
	if block != nil {
		uncles := block.Uncles()
		if index >= hexutil.Uint(len(uncles)) {
			log.Debug("Requested uncle not found", "number", blockNr, "hash", block.Hash(), "index", index)
			return nil, nil
		}
		ethblock := types.NewBlockWithHeader(uncles[index])
		hash := block.Hash()
		block = rpctypes.EthBlockToBlock(ethblock)
		block.SetHash(hash)
		return s.rpcMarshalBlock(ctx, block, false, false)
	}
	return nil, err
}

// GetUncleByBlockHashAndIndex returns the uncle block for the given block hash and index.
func (s *BlockChainAPI) GetUncleByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (map[string]interface{}, error) {
	block, err := s.b.BlockByHash(ctx, blockHash)
	if block != nil {
		uncles := block.Uncles()
		if index >= hexutil.Uint(len(uncles)) {
			log.Debug("Requested uncle not found", "number", block.Number(), "hash", blockHash, "index", index)
			return nil, nil
		}
		ethblock := types.NewBlockWithHeader(uncles[index])
		hash := block.Hash()
		block = rpctypes.EthBlockToBlock(ethblock)
		block.SetHash(hash)
		return s.rpcMarshalBlock(ctx, block, false, false)
	}
	return nil, err
}

// GetUncleCountByBlockNumber returns number of uncles in the block for the given block number
func (s *BlockChainAPI) GetUncleCountByBlockNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := s.b.ArtBlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}

// GetUncleCountByBlockHash returns number of uncles in the block for the given block hash
func (s *BlockChainAPI) GetUncleCountByBlockHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}

// GetCode returns the code stored at the given address in the states for the given block number.
func (s *BlockChainAPI) GetCode(_ context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	return s.b.GetCode(address, blockNrOrHash)
}

// GetStorageAt returns the storage from the states at the given address, key and
// block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
// numbers are also allowed.
func (s *BlockChainAPI) GetStorageAt(_ context.Context, address common.Address, hexKey string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	return s.b.GetStorageAt(address, hexKey, blockNrOrHash)
}

// OverrideAccount indicates the overriding fields of account during the execution
// of a message call.
// Note, states and stateDiff can't be specified at the same time. If states is
// set, message execution will only use the data in the given states. Otherwise
// if statDiff is set, all diff will be applied first and then execute the call
// message.
type OverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"states"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

// StateOverride is the collection of overridden accounts.
type StateOverride map[common.Address]OverrideAccount

// Apply overrides the fields of specified accounts into the given states.
func (diff *StateOverride) Apply(state *state.StateDB) error {
	if diff == nil {
		return nil
	}
	for addr, account := range *diff {
		// Override account nonce.
		if account.Nonce != nil {
			state.SetNonce(addr, uint64(*account.Nonce))
		}
		// Override account(contract) code.
		if account.Code != nil {
			state.SetCode(addr, *account.Code)
		}
		// Override account balance.
		if account.Balance != nil {
			state.SetBalance(addr, (*big.Int)(*account.Balance))
		}
		if account.State != nil && account.StateDiff != nil {
			return fmt.Errorf("account %s has both 'states' and 'stateDiff'", addr.Hex())
		}
		// Replace entire states if caller requires.
		if account.State != nil {
			state.SetStorage(addr, *account.State)
		}
		// Apply states diff into specified accounts.
		if account.StateDiff != nil {
			for key, value := range *account.StateDiff {
				state.SetState(addr, key, value)
			}
		}
	}
	// Now finalize the changes. Finalize is normally performed between transactions.
	// By using finalize, the overrides are semantically behaving as
	// if they were created in a transaction just before the tracing occur.
	state.Finalise(false)
	return nil
}

// BlockOverrides is a set of header fields to override.
type BlockOverrides struct {
	Number     *hexutil.Big
	Difficulty *hexutil.Big
	Time       *hexutil.Uint64
	GasLimit   *hexutil.Uint64
	Coinbase   *common.Address
	Random     *common.Hash
	BaseFee    *hexutil.Big
}

// Apply overrides the given header fields into the given block context.
func (diff *BlockOverrides) Apply(blockCtx *vm.BlockContext) {
	if diff == nil {
		return
	}
	if diff.Number != nil {
		blockCtx.BlockNumber = diff.Number.ToInt()
	}
	if diff.Difficulty != nil {
		blockCtx.Difficulty = diff.Difficulty.ToInt()
	}
	if diff.Time != nil {
		blockCtx.Time = uint64(*diff.Time)
	}
	if diff.GasLimit != nil {
		blockCtx.GasLimit = uint64(*diff.GasLimit)
	}
	if diff.Coinbase != nil {
		blockCtx.Coinbase = *diff.Coinbase
	}
	if diff.Random != nil {
		blockCtx.Random = diff.Random
	}
	if diff.BaseFee != nil {
		blockCtx.BaseFee = diff.BaseFee.ToInt()
	}
}

// ChainContextBackend provides methods required to implement ChainContext.
type ChainContextBackend interface {
	Engine() consensus.Engine
	HeaderByNumber(context.Context, rpc.BlockNumber) (*types.Header, error)
}

// ChainContext is an implementation of core.ChainContext. It's main use-case
// is instantiating a vm.BlockContext without having access to the BlockChain object.
type ChainContext struct {
	b   ChainContextBackend
	ctx context.Context
}

// NewChainContext creates a new ChainContext object.
func NewChainContext(ctx context.Context, backend ChainContextBackend) *ChainContext {
	return &ChainContext{ctx: ctx, b: backend}
}

func (context *ChainContext) Engine() consensus.Engine {
	return context.b.Engine()
}

func (context *ChainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	// This method is called to get the hash for a block number when executing the BLOCKHASH
	// opcode. Hence no need to search for non-canonical blocks.
	header, err := context.b.HeaderByNumber(context.ctx, rpc.BlockNumber(number))
	if err != nil || header.Hash() != hash {
		return nil
	}
	return header
}

// Call executes the given transaction on the states for the given block number.
//
// Additionally, the caller can specify a batch of contract for fields overriding.
//
// Note, this function doesn't make and changes in the states/blockchain and is
// useful to execute and retrieve values.
func (s *BlockChainAPI) Call(_ context.Context, args TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, _ *StateOverride, _ *BlockOverrides) (hexutil.Bytes, error) {
	data, err := s.b.DoCall(args, blockNrOrHash)
	if err != nil {
		return hexutil.Bytes{}, err
	}

	return (hexutil.Bytes)(data.Ret), nil
}

// EstimateGas returns an estimate of the amount of gas needed to execute the
// given transaction against the current pending block.
func (s *BlockChainAPI) EstimateGas(ctx context.Context, args TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	return s.b.EstimateGas(ctx, args, blockNrOrHash)
}

// RPCMarshalHeader converts the given header to the RPC output .
func RPCMarshalHeader(head *types.Header, hash common.Hash) map[string]interface{} {
	result := map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             hash,
		"parentHash":       head.ParentHash,
		"nonce":            head.Nonce,
		"mixHash":          head.MixDigest,
		"sha3Uncles":       head.UncleHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"miner":            head.Coinbase,
		"difficulty":       (*hexutil.Big)(head.Difficulty),
		"extraData":        hexutil.Bytes(head.Extra),
		"size":             hexutil.Uint64(head.Size()),
		"gasLimit":         hexutil.Uint64(head.GasLimit),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        hexutil.Uint64(head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}

	if head.BaseFee != nil {
		result["baseFeePerGas"] = (*hexutil.Big)(head.BaseFee)
	}

	if head.WithdrawalsHash != nil {
		result["withdrawalsRoot"] = head.WithdrawalsHash
	}

	return result
}

// RPCMarshalBlock converts the given block to the RPC output which depends on fullTx. If inclTx is true transactions are
// returned. When fullTx is true the returned block contains full transaction details, otherwise it will only contain
// transaction hashes.
func RPCMarshalBlock(block *rpctypes.Block, inclTx bool, fullTx bool, config *params.ChainConfig) (map[string]interface{}, error) {
	fields := RPCMarshalHeader(block.Header(), block.Hash())
	fields["size"] = hexutil.Uint64(block.Size())

	if inclTx {
		formatTx := func(idx int, tx *types.Transaction) interface{} {
			return tx.Hash()
		}
		if fullTx {
			formatTx = func(idx int, tx *types.Transaction) interface{} {
				return newRPCTransactionFromBlockIndex(block.EthBlock(), block.Hash(), uint64(idx), config)
			}
		}
		txs := block.Transactions()
		transactions := make([]interface{}, len(txs))
		for i, tx := range txs {
			transactions[i] = formatTx(i, tx)
		}
		fields["transactions"] = transactions
	}
	uncles := block.Uncles()
	uncleHashes := make([]common.Hash, len(uncles))
	for i, uncle := range uncles {
		uncleHashes[i] = uncle.Hash()
	}
	fields["uncles"] = uncleHashes
	if block.Header().WithdrawalsHash != nil {
		fields["withdrawals"] = block.Withdrawals()
	}
	return fields, nil
}

// rpcMarshalHeader uses the generalized output filler, then adds the total difficulty field, which requires
// a `BlockchainAPI`.
func (s *BlockChainAPI) rpcMarshalHeader(_ context.Context, header *types.Header, hash common.Hash) map[string]interface{} {
	fields := RPCMarshalHeader(header, hash)
	// fields["totalDifficulty"] = (*hexutil.Big)(s.b.GetTd(ctx, header.Hash()))
	return fields
}

// rpcMarshalBlock uses the generalized output filler, then adds the total difficulty field, which requires
// a `BlockchainAPI`.
func (s *BlockChainAPI) rpcMarshalBlock(_ context.Context, b *rpctypes.Block, inclTx bool, fullTx bool) (map[string]interface{}, error) {
	return RPCMarshalBlock(b, inclTx, fullTx, s.b.ChainConfig())
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash        *common.Hash      `json:"blockHash"`
	BlockNumber      *hexutil.Big      `json:"blockNumber"`
	From             common.Address    `json:"from"`
	Gas              hexutil.Uint64    `json:"gas"`
	GasPrice         *hexutil.Big      `json:"gasPrice"`
	GasFeeCap        *hexutil.Big      `json:"maxFeePerGas,omitempty"`
	GasTipCap        *hexutil.Big      `json:"maxPriorityFeePerGas,omitempty"`
	Hash             common.Hash       `json:"hash"`
	Input            hexutil.Bytes     `json:"input"`
	Nonce            hexutil.Uint64    `json:"nonce"`
	To               *common.Address   `json:"to"`
	TransactionIndex *hexutil.Uint64   `json:"transactionIndex"`
	Value            *hexutil.Big      `json:"value"`
	Type             hexutil.Uint64    `json:"type"`
	Accesses         *types.AccessList `json:"accessList,omitempty"`
	ChainID          *hexutil.Big      `json:"chainId,omitempty"`
	V                *hexutil.Big      `json:"v"`
	R                *hexutil.Big      `json:"r"`
	S                *hexutil.Big      `json:"s"`
}

// newRPCTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func newRPCTransaction(tx *types.Transaction, blockHash common.Hash, blockNumber uint64, blockTime uint64, index uint64, baseFee *big.Int, config *params.ChainConfig) *RPCTransaction {
	signer := types.MakeSigner(config, new(big.Int).SetUint64(blockNumber), blockTime)
	from, _ := types.Sender(signer, tx)

	return newRPCTransactionWithFrom(tx, blockHash, blockNumber, index, baseFee, from)
}

func newRPCTransactionWithFrom(tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index uint64, baseFee *big.Int, from common.Address) *RPCTransaction {
	v, r, s := tx.RawSignatureValues()
	result := &RPCTransaction{
		Type:     hexutil.Uint64(tx.Type()),
		From:     from,
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    hexutil.Bytes(tx.Data()),
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(v),
		R:        (*hexutil.Big)(r),
		S:        (*hexutil.Big)(s),
	}
	if blockHash != (common.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}
	switch tx.Type() {
	case types.LegacyTxType:
		// if a legacy transaction has an EIP-155 chain id, include it explicitly
		if id := tx.ChainId(); id.Sign() != 0 {
			result.ChainID = (*hexutil.Big)(id)
		}
	case types.AccessListTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
	case types.DynamicFeeTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		// if the transaction has been mined, compute the effective gas price
		if baseFee != nil && blockHash != (common.Hash{}) {
			// price = min(tip, gasFeeCap - baseFee) + baseFee
			price := math.BigMin(new(big.Int).Add(tx.GasTipCap(), baseFee), tx.GasFeeCap())
			result.GasPrice = (*hexutil.Big)(price)
		} else {
			result.GasPrice = (*hexutil.Big)(tx.GasFeeCap())
		}
	}
	return result
}

// NewTransactionFromMsg returns a txs that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewTransactionFromMsg(
	msg *txs.MsgEthereumTx,
	blockHash common.Hash,
	blockNumber, index uint64,
	baseFee *big.Int,
	cfg *params.ChainConfig,
) *RPCTransaction {
	tx := msg.AsTransaction()
	// use latest singer, so use time.now as block time.
	if msg.From != "" {
		return newRPCTransactionWithFrom(tx, blockHash, blockNumber, index, baseFee, common.HexToAddress(msg.From))
	}
	return newRPCTransaction(tx, blockHash, blockNumber, uint64(time.Now().Unix()), index, baseFee, cfg)
}

// NewRPCPendingTransaction returns a pending transaction that will serialize to the RPC representation
func NewRPCPendingTransaction(tx *types.Transaction, current *types.Header, config *params.ChainConfig) *RPCTransaction {
	var (
		baseFee     *big.Int
		blockNumber = uint64(0)
		blockTime   = uint64(0)
	)
	if current != nil {
		baseFee = misc.CalcBaseFee(config, current)
		blockNumber = current.Number.Uint64()
		blockTime = current.Time
	}
	return newRPCTransaction(tx, common.Hash{}, blockNumber, blockTime, 0, baseFee, config)
}

// newRPCTransactionFromBlockIndex returns a transaction that will serialize to the RPC representation.
func newRPCTransactionFromBlockIndex(b *types.Block, blockHash common.Hash, index uint64, config *params.ChainConfig) *RPCTransaction {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	return newRPCTransaction(txs[index], blockHash, b.NumberU64(), b.Time(), index, b.BaseFee(), config)
}

// newRPCRawTransactionFromBlockIndex returns the bytes of a transaction given a block and a transaction index.
func newRPCRawTransactionFromBlockIndex(b *types.Block, index uint64) hexutil.Bytes {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	blob, _ := txs[index].MarshalBinary()
	return blob
}

// AccessListResult returns an optional access list
// It's the result of the `debug_createAccessList` RPC call.
// It contains an error if the transaction itself failed.
type AccessListResult struct {
	Accesslist *types.AccessList `json:"accessList"`
	Error      string            `json:"error,omitempty"`
	GasUsed    hexutil.Uint64    `json:"gasUsed"`
}

// CreateAccessList creates an EIP-2930 type AccessList for the given transaction.
// Reexec and BlockNrOrHash can be specified to create the accessList on top of a certain states.
func (s *BlockChainAPI) CreateAccessList(_ context.Context, _ TransactionArgs, _ *rpc.BlockNumberOrHash) (*AccessListResult, error) {
	return nil, errors.New("CreateAccessList is not implemented")
}

// TransactionAPI exposes methods for reading and creating transaction data.
type TransactionAPI struct {
	b         Backend
	logger    log.Logger
	nonceLock *AddrLocker
}

// NewTransactionAPI creates a new RPC service with methods for interacting with transactions.
func NewTransactionAPI(b Backend, logger log.Logger, nonceLock *AddrLocker) *TransactionAPI {
	// The signer used by the API should always be the 'latest' known one because we expect
	// signers to be backwards-compatible with old transactions.
	return &TransactionAPI{b, logger, nonceLock}
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.
func (s *TransactionAPI) GetBlockTransactionCountByNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := s.b.ArtBlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}

// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
func (s *TransactionAPI) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}

// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.
func (s *TransactionAPI) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) *RPCTransaction {
	if block, _ := s.b.ArtBlockByNumber(ctx, blockNr); block != nil {
		return newRPCTransactionFromBlockIndex(block.EthBlock(), block.Hash(), uint64(index), s.b.ChainConfig())
	}
	return nil
}

// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
func (s *TransactionAPI) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) *RPCTransaction {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		return newRPCTransactionFromBlockIndex(block.EthBlock(), block.Hash(), uint64(index), s.b.ChainConfig())
	}
	return nil
}

// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.
func (s *TransactionAPI) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) hexutil.Bytes {
	if block, _ := s.b.ArtBlockByNumber(ctx, blockNr); block != nil {
		return newRPCRawTransactionFromBlockIndex(block.EthBlock(), uint64(index))
	}
	return nil
}

// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.
func (s *TransactionAPI) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) hexutil.Bytes {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		return newRPCRawTransactionFromBlockIndex(block.EthBlock(), uint64(index))
	}
	return nil
}

// GetTransactionCount returns the number of transactions the given address has sent for the given block number
func (s *TransactionAPI) GetTransactionCount(_ context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	return s.b.GetTransactionCount(address, blockNrOrHash)
}

// GetTransactionByHash returns the transaction for the given hash
func (s *TransactionAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) (*RPCTransaction, error) {
	return s.b.GetTransaction(ctx, hash)
}

// GetRawTransactionByHash returns the bytes of the transaction for the given hash.
func (s *TransactionAPI) GetRawTransactionByHash(_ context.Context, _ common.Hash) (hexutil.Bytes, error) {
	// TODO
	return nil, errors.New("GetRawTransactionByHash is not implemented")
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
func (s *TransactionAPI) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	return s.b.GetTransactionReceipt(ctx, hash)
}

// SubmitTransaction is a helper function that submits tx to txPool and logs a message.
func SubmitTransaction(ctx context.Context, logger log.Logger, b Backend, tx *types.Transaction) (common.Hash, error) {
	// If the transaction fee cap is already specified, ensure the
	// fee of the given transaction is _reasonable_.
	if err := checkTxFee(tx.GasPrice(), tx.Gas(), b.RPCTxFeeCap()); err != nil {
		return common.Hash{}, err
	}

	customizedVerification := isCustomizedVerificationRequired(tx)
	if !customizedVerification && !b.UnprotectedAllowed() && !tx.Protected() {
		// Ensure only eip155 signed transactions are submitted if EIP155Required is set.
		return common.Hash{}, errors.New("only replay-protected (EIP-155) transactions allowed over RPC")
	}
	if err := b.SendTx(ctx, tx); err != nil {
		return common.Hash{}, err
	}
	// Print a log with full tx details for manual investigations and interventions
	if customizedVerification {
		// no need to check customized verification tx
		return tx.Hash(), nil
	}

	head, err := b.CurrentHeader()
	if err != nil {
		return common.Hash{}, err
	}

	signer := types.MakeSigner(b.ChainConfig(), head.Number, head.Time)
	from, err := types.Sender(signer, tx)
	if err != nil {
		return common.Hash{}, err
	}

	if tx.To() == nil {
		addr := crypto.CreateAddress(from, tx.Nonce())
		logger.Debug("Submitted contract creation", "hash", tx.Hash().Hex(), "from", from, "nonce", tx.Nonce(), "contract", addr.Hex(), "value", tx.Value())
	} else {
		logger.Debug("Submitted transaction", "hash", tx.Hash().Hex(), "from", from, "nonce", tx.Nonce(), "recipient", tx.To(), "value", tx.Value())
	}
	return tx.Hash(), nil
}

// SendTransaction creates a transaction for the given argument, sign it and submit it to the
// transaction pool.
func (s *TransactionAPI) SendTransaction(ctx context.Context, args TransactionArgs) (common.Hash, error) {
	if err := args.setDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}
	signed, err := s.b.SignTransaction(&args)
	if err != nil {
		log.Warn("Failed transaction send attempt", "from", args.from(), "to", args.To, "value", args.Value.ToInt(), "err", err)
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.logger, s.b, signed)
}

// FillTransaction fills the defaults (nonce, gas, gasPrice or 1559 fields)
// on a given unsigned transaction, and returns it to the caller for further
// processing (signing + broadcast).
func (s *TransactionAPI) FillTransaction(ctx context.Context, args TransactionArgs) (*SignTransactionResult, error) {
	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	// Assemble the transaction and obtain rlp
	tx := args.toTransaction()
	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, tx}, nil
}

// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *TransactionAPI) SendRawTransaction(ctx context.Context, input hexutil.Bytes) (common.Hash, error) {
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(input); err != nil {
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.logger, s.b, tx)
}

// Sign calculates an ECDSA signature for:
// keccak256("\x19Ethereum Signed Message:\n" + len(message) + message).
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The account associated with addr must be unlocked.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func (s *TransactionAPI) Sign(addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	s.logger.Debug("eth_sign", "address", addr.Hex(), "data", common.Bytes2Hex(data))
	return s.b.Sign(addr, data)
}

// SignTransactionResult represents a RLP encoded signed transaction.
type SignTransactionResult struct {
	Raw hexutil.Bytes      `json:"raw"`
	Tx  *types.Transaction `json:"tx"`
}

// SignTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.
func (s *TransactionAPI) SignTransaction(ctx context.Context, args TransactionArgs) (*SignTransactionResult, error) {
	// gas, gas limit, nonce checking are made in SignTransaction
	signed, err := s.b.SignTransaction(&args)
	if err != nil {
		return nil, err
	}

	data, err := signed.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, signed}, nil
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (s *TransactionAPI) PendingTransactions() ([]*RPCTransaction, error) {
	// TODO
	return nil, errors.New("PendingTransactions is not implemented")
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (s *TransactionAPI) Resend(_ context.Context, _ TransactionArgs, _ *hexutil.Big, _ *hexutil.Uint64) (common.Hash, error) {
	// TODO
	return common.Hash{}, errors.New("Resend is not implemented")
}

// DebugAPI is the collection of Ethereum APIs exposed over the debugging
// namespace.
type DebugAPI struct {
	b Backend
}

// NewDebugAPI creates a new instance of DebugAPI.
func NewDebugAPI(b Backend) *DebugAPI {
	return &DebugAPI{b: b}
}

// GetRawHeader retrieves the RLP encoding for a single header.
func (api *DebugAPI) GetRawHeader(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	var hash common.Hash
	if h, ok := blockNrOrHash.Hash(); ok {
		hash = h
	} else {
		block, err := api.b.BlockByNumberOrHash(ctx, blockNrOrHash)
		if err != nil {
			return nil, err
		}
		hash = block.Hash()
	}
	header, _ := api.b.HeaderByHash(ctx, hash)
	if header == nil {
		return nil, fmt.Errorf("header #%d not found", hash)
	}
	return rlp.EncodeToBytes(header)
}

// GetRawBlock retrieves the RLP encoded for a single block.
func (api *DebugAPI) GetRawBlock(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	var hash common.Hash
	if h, ok := blockNrOrHash.Hash(); ok {
		hash = h
	} else {
		block, err := api.b.BlockByNumberOrHash(ctx, blockNrOrHash)
		if err != nil {
			return nil, err
		}
		hash = block.Hash()
	}
	block, _ := api.b.BlockByHash(ctx, hash)
	if block == nil {
		return nil, fmt.Errorf("block #%d not found", hash)
	}
	return rlp.EncodeToBytes(block)
}

// GetRawReceipts retrieves the binary-encoded receipts of a single block.
func (api *DebugAPI) GetRawReceipts(_ context.Context, _ rpc.BlockNumberOrHash) ([]hexutil.Bytes, error) {
	return nil, errors.New("GetRawReceipts is not implemented")
}

// GetRawTransaction returns the bytes of the transaction for the given hash.
func (api *DebugAPI) GetRawTransaction(_ context.Context, _ common.Hash) (hexutil.Bytes, error) {
	// TODO
	return hexutil.Bytes{}, errors.New("GetRawTransaction is not implemented")
}

// PrintBlock retrieves a block and returns its pretty printed form.
func (api *DebugAPI) PrintBlock(ctx context.Context, number uint64) (string, error) {
	block, _ := api.b.ArtBlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	return spew.Sdump(block), nil
}

// ChaindbProperty returns leveldb properties of the key-value database.
func (api *DebugAPI) ChaindbProperty(_ string) (string, error) {
	return "", errors.New("ChaindbProperty is not implemented")
}

// ChaindbCompact flattens the entire key-value database into a single level,
// removing all unused slots and merging all keys.
func (api *DebugAPI) ChaindbCompact() error {
	return errors.New("ChaindbCompact is not implemented")
}

// SetHead rewinds the head of the blockchain to a previous block.
func (api *DebugAPI) SetHead(_ hexutil.Uint64) {
	// TODO
}

// NetAPI offers network related RPC methods
type NetAPI struct {
	net            *p2p.Server
	networkVersion uint64
}

// NewNetAPI creates a new net API instance.
func NewNetAPI(net *p2p.Server, networkVersion uint64) *NetAPI {
	return &NetAPI{net, networkVersion}
}

// Listening returns an indication if the node is listening for network connections.
func (s *NetAPI) Listening() bool {
	return true // always listening
}

// PeerCount returns the number of connected peers
func (s *NetAPI) PeerCount() hexutil.Uint {
	if s.net == nil {
		return 0
	}
	return hexutil.Uint(s.net.PeerCount())
}

// Version returns the current ethereum protocol version.
func (s *NetAPI) Version() string {
	return fmt.Sprintf("%d", s.networkVersion)
}

// checkTxFee is an internal function used to check whether the fee of
// the given transaction is _reasonable_(under the cap).
func checkTxFee(gasPrice *big.Int, gas uint64, cap float64) error {
	// Short circuit if there is no cap for transaction fee at all.
	if cap == 0 {
		return nil
	}
	feeEth := new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gas))), new(big.Float).SetInt(big.NewInt(params.Ether)))
	feeFloat, _ := feeEth.Float64()
	if feeFloat > cap {
		return fmt.Errorf("tx fee (%.2f art) exceeds the configured cap (%.2f art)", feeFloat, cap)
	}
	return nil
}

func isCustomizedVerificationRequired(tx *types.Transaction) bool {
	zero := big.NewInt(0)
	v, r, s := tx.RawSignatureValues()
	return (v == nil || r == nil || s == nil) || (v.Cmp(zero) == 0 && r.Cmp(zero) == 0 && s.Cmp(zero) == 0)
}
