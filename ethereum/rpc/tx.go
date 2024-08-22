package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/artela-network/artela/common/aspect"
	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	rpcutils "github.com/artela-network/artela/ethereum/rpc/utils"
	"github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/ethereum/utils"
	"github.com/artela-network/artela/x/evm/txs"
	"github.com/artela-network/artela/x/evm/txs/support"
	evmtypes "github.com/artela-network/artela/x/evm/types"
)

// Transaction API

func (b *BackendImpl) SendTx(ctx context.Context, signedTx *ethtypes.Transaction) error {
	// verify the ethereum tx
	ethereumTx := &txs.MsgEthereumTx{}
	if err := ethereumTx.FromEthereumTx(signedTx); err != nil {
		b.logger.Error("transaction converting failed", "error", err.Error())
		return err
	}

	if err := ethereumTx.ValidateBasic(); err != nil {
		b.logger.Debug("tx failed basic validation", "error", err.Error())
		return err
	}

	// Query params to use the EVM denomination
	res, err := b.queryClient.QueryClient.Params(b.ctx, &txs.QueryParamsRequest{})
	if err != nil {
		b.logger.Error("failed to query evm params", "error", err.Error())
		return err
	}

	cosmosTx, err := ethereumTx.BuildTx(b.clientCtx.TxConfig.NewTxBuilder(), res.Params.EvmDenom)
	if err != nil {
		b.logger.Error("failed to build cosmos tx", "error", err.Error())
		return err
	}

	// Encode transaction by default Tx encoder
	txBytes, err := b.clientCtx.TxConfig.TxEncoder()(cosmosTx)
	if err != nil {
		b.logger.Error("failed to encode eth tx using default encoder", "error", err.Error())
		return err
	}

	// txHash := ethereumTx.AsTransaction().Hash()

	syncCtx := b.clientCtx.WithBroadcastMode(flags.BroadcastSync)
	rsp, err := syncCtx.BroadcastTx(txBytes)
	if rsp != nil && rsp.Code != 0 {
		err = errorsmod.ABCIError(rsp.Codespace, rsp.Code, rsp.RawLog)
	}
	if err != nil {
		b.logger.Error("failed to broadcast tx", "error", err.Error())
		return err
	}

	return nil
}

func (b *BackendImpl) GetTransaction(ctx context.Context, txHash common.Hash) (*rpctypes.RPCTransaction, error) {
	_, tx, err := b.getTransaction(ctx, txHash)
	return tx, err
}

func (b *BackendImpl) GetTransactionCount(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	n := hexutil.Uint64(0)
	height, err := b.blockNumberFromCosmos(blockNrOrHash)
	if err != nil {
		return &n, err
	}
	header, err := b.CurrentHeader()
	if err != nil {
		return &n, err
	}
	if height.Int64() > header.Number.Int64() {
		return &n, fmt.Errorf(
			"cannot query with height in the future (current: %d, queried: %d); please provide a valid height",
			header.Number, height)
	}
	// Get nonce (sequence) from account
	from := sdktypes.AccAddress(address.Bytes())
	accRet := b.clientCtx.AccountRetriever

	if err = accRet.EnsureExists(b.clientCtx, from); err != nil {
		// account doesn't exist yet, return 0
		b.logger.Info("GetTransactionCount faild, return 0. Account doesn't exist yet", "account", address.Hex(), "error", err)
		return &n, nil
	}

	includePending := height == rpc.PendingBlockNumber
	nonce, err := b.getAccountNonce(address, includePending, height.Int64())
	if err != nil {
		return nil, err
	}

	n = hexutil.Uint64(nonce)
	return &n, nil
}

func (b *BackendImpl) GetTxMsg(ctx context.Context, txHash common.Hash) (*txs.MsgEthereumTx, error) {
	msg, _, err := b.getTransaction(ctx, txHash)
	return msg, err
}

func (b *BackendImpl) SignTransaction(args *rpctypes.TransactionArgs) (*ethtypes.Transaction, error) {
	_, err := b.clientCtx.Keyring.KeyByAddress(sdktypes.AccAddress(args.From.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("failed to find key in the node's keyring; %s; %s", keystore.ErrNoMatch, err.Error())
	}

	if args.ChainID != nil && (b.chainID).Cmp((*big.Int)(args.ChainID)) != 0 {
		return nil, fmt.Errorf("chainId does not match node's (have=%v, want=%v)", args.ChainID, (*hexutil.Big)(b.chainID))
	}

	bn, err := b.BlockNumber()
	if err != nil {
		return nil, err
	}

	bt, err := b.BlockTimeByNumber(int64(bn))
	if err != nil {
		return nil, err
	}

	cfg, err := b.chainConfig()
	if err != nil {
		return nil, err
	}
	signer := ethtypes.MakeSigner(cfg, new(big.Int).SetUint64(uint64(bn)), bt)

	// LegacyTx derives chainID from the signature. To make sure the msg.ValidateBasic makes
	// the corresponding chainID validation, we need to sign the transaction before calling it

	// Sign transaction
	msg := args.ToEVMTransaction()
	return msg.SignEthereumTx(signer, b.clientCtx.Keyring)
}

// GetTransactionReceipt get receipt by transaction hash
func (b *BackendImpl) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	res, err := b.GetTxByEthHash(hash)
	if err != nil {
		b.logger.Debug("GetTransactionReceipt failed", "error", err)
		return nil, nil
	}
	resBlock, err := b.CosmosBlockByNumber(rpc.BlockNumber(res.Height))
	if err != nil {
		b.logger.Debug("GetTransactionReceipt failed", "error", err)
		return nil, nil
	}
	tx, err := b.clientCtx.TxConfig.TxDecoder()(resBlock.Block.Txs[res.TxIndex])
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx: %w", err)
	}
	ethMsg := tx.GetMsgs()[res.MsgIndex].(*txs.MsgEthereumTx)

	txData, err := txs.UnpackTxData(ethMsg.Data)
	if err != nil {
		return nil, err
	}

	cumulativeGasUsed := uint64(0)
	blockRes, err := b.CosmosBlockResultByNumber(&res.Height)
	if err != nil {
		b.logger.Debug("GetTransactionReceipt failed", "error", err)
		return nil, nil
	}
	for _, txResult := range blockRes.TxsResults[0:res.TxIndex] {
		cumulativeGasUsed += uint64(txResult.GasUsed)
	}
	cumulativeGasUsed += res.CumulativeGasUsed

	var status hexutil.Uint
	if res.Failed {
		status = hexutil.Uint(ethtypes.ReceiptStatusFailed)
	} else {
		status = hexutil.Uint(ethtypes.ReceiptStatusSuccessful)
	}

	// parse tx logs from events
	msgIndex := int(res.MsgIndex)
	logs, _ := rpcutils.TxLogsFromEvents(blockRes.TxsResults[res.TxIndex].Events, msgIndex)

	if res.EthTxIndex == -1 {
		// Fallback to find tx index by iterating all valid eth transactions
		msgs := b.EthMsgsFromCosmosBlock(resBlock, blockRes)
		for i := range msgs {
			if msgs[i].Hash == hash.Hex() {
				res.EthTxIndex = int32(i) // #nosec G701
				break
			}
		}
	}
	// return error if still unable to find the eth tx index
	if res.EthTxIndex == -1 {
		return nil, errors.New("can't find index of ethereum tx")
	}

	receipt := map[string]interface{}{
		// Consensus fields: These fields are defined by the Yellow Paper
		"status":            status,
		"cumulativeGasUsed": hexutil.Uint64(cumulativeGasUsed),
		"logsBloom":         ethtypes.BytesToBloom(ethtypes.LogsBloom(logs)),
		"logs":              logs,

		// Implementation fields: These fields are added by geth when processing a transaction.
		// They are stored in the chain database.
		"transactionHash": hash,
		"contractAddress": nil,
		"gasUsed":         hexutil.Uint64(res.GasUsed),

		// Inclusion information: These fields provide information about the inclusion of the
		// transaction corresponding to this receipt.
		"blockHash":        common.BytesToHash(resBlock.Block.Header.Hash()).Hex(),
		"blockNumber":      hexutil.Uint64(res.Height),
		"transactionIndex": hexutil.Uint64(res.EthTxIndex),

		// sender and receiver (contract or EOA) addreses
		"from": res.Sender,
		"to":   txData.GetTo(),
		"type": hexutil.Uint(ethMsg.AsTransaction().Type()),
	}

	if logs == nil {
		receipt["logs"] = []*ethtypes.Log{}
	}

	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if txData.GetTo() == nil || aspect.IsAspectDeploy(txData.GetTo(), txData.GetData()) {
		receipt["contractAddress"] = crypto.CreateAddress(common.HexToAddress(res.Sender), txData.GetNonce())
	}

	if dynamicTx, ok := txData.(*txs.DynamicFeeTx); ok {
		baseFee, err := b.BaseFee(blockRes)
		if err == nil {
			receipt["effectiveGasPrice"] = hexutil.Big(*dynamicTx.EffectiveGasPrice(baseFee))
		}
	}

	return receipt, nil
}

func (b *BackendImpl) RPCTxFeeCap() float64 {
	return b.cfg.RPCTxFeeCap
}

func (b *BackendImpl) UnprotectedAllowed() bool {
	if b.cfg.AppCfg == nil {
		return false
	}
	return b.cfg.AppCfg.JSONRPC.AllowUnprotectedTxs
}

func (b *BackendImpl) PendingTransactions() ([]*sdktypes.Tx, error) {
	res, err := b.clientCtx.Client.UnconfirmedTxs(b.ctx, nil)
	if err != nil {
		return nil, err
	}

	result := make([]*sdktypes.Tx, 0, len(res.Txs))
	for _, txBz := range res.Txs {
		tx, err := b.clientCtx.TxConfig.TxDecoder()(txBz)
		if err != nil {
			return nil, err
		}
		result = append(result, &tx)
	}

	return result, nil
}

func (b *BackendImpl) GetResendArgs(args rpctypes.TransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (rpctypes.TransactionArgs, error) {
	chainID, err := types.ParseChainID(b.clientCtx.ChainID)
	if err != nil {
		return rpctypes.TransactionArgs{}, err
	}

	cfg := b.ChainConfig()
	if cfg == nil {
		header, err := b.CurrentHeader()
		if err != nil {
			return rpctypes.TransactionArgs{}, err
		}
		cfg = support.DefaultChainConfig().EthereumConfig(header.Number.Int64(), chainID)
	}

	// use the latest signer for the new tx
	signer := ethtypes.LatestSigner(cfg)

	matchTx := args.ToTransaction()

	// Before replacing the old transaction, ensure the _new_ transaction fee is reasonable.
	price := matchTx.GasPrice()
	if gasPrice != nil {
		price = gasPrice.ToInt()
	}
	gas := matchTx.Gas()
	if gasLimit != nil {
		gas = uint64(*gasLimit)
	}
	if err := rpctypes.CheckTxFee(price, gas, b.RPCTxFeeCap()); err != nil {
		return rpctypes.TransactionArgs{}, err
	}

	pending, err := b.PendingTransactions()
	if err != nil {
		return rpctypes.TransactionArgs{}, err
	}

	for _, tx := range pending {
		wantSigHash := signer.Hash(matchTx)

		// TODO, wantSigHash?
		msg, err := txs.UnwrapEthereumMsg(tx, wantSigHash)
		if err != nil {
			// not ethereum tx
			continue
		}

		pendingTx := msg.AsTransaction()
		pFrom, err := ethtypes.Sender(signer, pendingTx)
		if err != nil {
			continue
		}

		if pFrom == *args.From && signer.Hash(pendingTx) == wantSigHash {
			// Match. Re-sign and send the transaction.
			if gasPrice != nil && (*big.Int)(gasPrice).Sign() != 0 {
				args.GasPrice = gasPrice
			}
			if gasLimit != nil && *gasLimit != 0 {
				args.Gas = gasLimit
			}

			return args, nil
		}
	}

	return rpctypes.TransactionArgs{}, fmt.Errorf("transaction %s not found", matchTx.Hash().String())
}

// Sign signs the provided data using the private key of address via Geth's signature standard.
func (b *BackendImpl) Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	from := sdktypes.AccAddress(address.Bytes())

	_, err := b.clientCtx.Keyring.KeyByAddress(from)
	if err != nil {
		return nil, fmt.Errorf("%s; %s", keystore.ErrNoMatch, err.Error())
	}

	// Sign the requested hash with the wallet
	signature, _, err := b.clientCtx.Keyring.SignByAddress(from, data)
	if err != nil {
		return nil, err
	}

	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}

// GetSender extracts the sender address from the signature values using the latest signer for the given chainID.
func (b *BackendImpl) GetSender(msg *txs.MsgEthereumTx, chainID *big.Int) (from common.Address, err error) {
	if msg.From != "" {
		return common.HexToAddress(msg.From), nil
	}

	tx := msg.AsTransaction()
	// retrieve sender info from aspect if tx is not signed
	if utils.IsCustomizedVerification(tx) {
		bn, err := b.BlockNumber()
		if err != nil {
			return common.Address{}, err
		}
		ctx := rpctypes.ContextWithHeight(int64(bn))

		res, err := b.queryClient.GetSender(ctx, msg)
		if err != nil {
			return common.Address{}, err
		}

		from = common.HexToAddress(res.Sender)
	} else {
		signer := ethtypes.LatestSignerForChainID(chainID)
		from, err = signer.Sender(tx)
		if err != nil {
			return common.Address{}, err
		}
	}

	msg.From = from.Hex()
	return from, nil
}

func (b *BackendImpl) getTransaction(_ context.Context, txHash common.Hash) (*txs.MsgEthereumTx, *rpctypes.RPCTransaction, error) {
	res, err := b.GetTxByEthHash(txHash)
	hexTx := txHash.Hex()

	if err != nil {
		b.logger.Debug("GetTxByEthHash failed, try to getTransactionByHashPending", "error", err)
		return b.getTransactionByHashPending(txHash)
	}

	block, err := b.CosmosBlockByNumber(rpc.BlockNumber(res.Height))
	if err != nil {
		return nil, nil, err
	}

	tx, err := b.clientCtx.TxConfig.TxDecoder()(block.Block.Txs[res.TxIndex])
	if err != nil {
		return nil, nil, err
	}

	// the `res.MsgIndex` is inferred from tx index, should be within the bound.
	msg, ok := tx.GetMsgs()[res.MsgIndex].(*txs.MsgEthereumTx)
	if !ok {
		return nil, nil, errors.New("invalid ethereum tx")
	}
	msg.From = res.Sender

	blockRes, err := b.CosmosBlockResultByNumber(&block.Block.Height)
	if err != nil {
		b.logger.Debug("block result not found", "height", block.Block.Height, "error", err.Error())
		return nil, nil, nil
	}

	if res.EthTxIndex == -1 {
		// Fallback to find tx index by iterating all valid eth transactions
		msgs := b.EthMsgsFromCosmosBlock(block, blockRes)
		for i := range msgs {
			if msgs[i].Hash == hexTx {
				res.EthTxIndex = int32(i)
				break
			}
		}
	}
	// if we still unable to find the eth tx index, return error, shouldn't happen.
	if res.EthTxIndex == -1 {
		return msg, nil, errors.New("can't find index of ethereum tx")
	}

	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle the error for pruned node.
		b.logger.Error("failed to fetch Base Fee from prunned block. Check node prunning configuration", "height", blockRes.Height, "error", err)
	}

	cfg, err := b.chainConfig()
	if err != nil {
		return msg, nil, err
	}

	return msg, rpctypes.NewTransactionFromMsg(
		msg,
		common.BytesToHash(block.BlockID.Hash.Bytes()),
		uint64(res.Height),
		uint64(res.EthTxIndex),
		baseFee,
		cfg,
	), nil
}

func (b *BackendImpl) GetPoolTransactions() (ethtypes.Transactions, error) {
	b.logger.Debug("called eth.rpc.rpctypes.GetPoolTransactions")
	return nil, errors.New("GetPoolTransactions is not implemented")
}

func (b *BackendImpl) GetPoolTransaction(txHash common.Hash) *ethtypes.Transaction {
	b.logger.Error("GetPoolTransaction is not implemented")
	return nil
}

func (b *BackendImpl) GetPoolNonce(_ context.Context, addr common.Address) (uint64, error) {
	return 0, errors.New("GetPoolNonce is not implemented")
}

func (b *BackendImpl) Stats() (int, int) {
	b.logger.Error("Stats is not implemented")
	return 0, 0
}

func (b *BackendImpl) TxPoolContent() (
	map[common.Address]ethtypes.Transactions, map[common.Address]ethtypes.Transactions,
) {
	b.logger.Error("TxPoolContent is not implemented")
	return nil, nil
}

func (b *BackendImpl) TxPoolContentFrom(addr common.Address) (
	ethtypes.Transactions, ethtypes.Transactions,
) {
	return nil, nil
}

func (b *BackendImpl) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.scope.Track(b.newTxsFeed.Subscribe(ch))
}

func (b *BackendImpl) version() (string, error) {
	cfg, err := b.chainConfig()
	if err != nil {
		return "", err
	}

	if cfg.ChainID == nil {
		b.logger.Error("eth.rpc.rpctypes.Version", "ChainID is nil")
		return "", errors.New("chain id is not valid")
	}
	return cfg.ChainID.String(), nil
}

func (b *BackendImpl) GetTxByEthHash(hash common.Hash) (*types.TxResult, error) {
	// if b.indexer != nil {
	// 	return b.indexer.GetByTxHash(hash)
	// }

	// fallback to tendermint tx indexer
	query := fmt.Sprintf("%s.%s='%s'", evmtypes.TypeMsgEthereumTx, evmtypes.AttributeKeyEthereumTxHash, hash.Hex())
	txResult, err := b.queryCosmosTxIndexer(query, func(txs *rpctypes.ParsedTxs) *rpctypes.ParsedTx {
		return txs.GetTxByHash(hash)
	})
	if err != nil {
		return nil, fmt.Errorf("GetTxByEthHash %s, %w", hash.Hex(), err)
	}
	return txResult, nil
}

func (b *BackendImpl) queryCosmosTxIndexer(query string, txGetter func(*rpctypes.ParsedTxs) *rpctypes.ParsedTx) (*types.TxResult, error) {
	resTxs, err := b.clientCtx.Client.TxSearch(b.ctx, query, false, nil, nil, "")
	if err != nil {
		return nil, err
	}
	if len(resTxs.Txs) == 0 {
		return nil, errors.New("ethereum tx not found")
	}
	txResult := resTxs.Txs[0]
	if !rpctypes.TxSuccessOrExceedsBlockGasLimit(&txResult.TxResult) {
		return nil, errors.New("invalid ethereum tx")
	}

	var tx sdktypes.Tx
	if txResult.TxResult.Code != 0 {
		// it's only needed when the tx exceeds block gas limit
		tx, err = b.clientCtx.TxConfig.TxDecoder()(txResult.Tx)
		if err != nil {
			return nil, fmt.Errorf("invalid ethereum tx, %w", err)
		}
	}

	return rpctypes.ParseTxIndexerResult(txResult, tx, txGetter)
}

// getTransactionByHashPending find pending tx from mempool
func (b *BackendImpl) getTransactionByHashPending(txHash common.Hash) (*txs.MsgEthereumTx, *rpctypes.RPCTransaction, error) {
	hexTx := txHash.Hex()
	// try to find tx in mempool
	ptxs, err := b.PendingTransactions()
	if err != nil {
		b.logger.Debug("pending tx not found", "hash", hexTx, "error", err.Error())
		return nil, nil, nil
	}

	for _, tx := range ptxs {
		msg, err := txs.UnwrapEthereumMsg(tx, txHash)
		if err != nil {
			// not ethereum tx
			continue
		}

		cfg, err := b.chainConfig()
		if err != nil {
			return msg, nil, err
		}
		if msg.Hash == hexTx {
			// use zero block values since it's not included in a block yet
			rpctx := rpctypes.NewTransactionFromMsg(
				msg,
				common.Hash{},
				uint64(0),
				uint64(0),
				nil,
				cfg,
			)
			return msg, rpctx, nil
		}
	}

	b.logger.Debug("tx not found", "hash", hexTx)
	return nil, nil, nil
}

func (b *BackendImpl) getAccountNonce(accAddr common.Address, pending bool, height int64) (uint64, error) {
	queryClient := authtypes.NewQueryClient(b.clientCtx)
	adr := sdktypes.AccAddress(accAddr.Bytes()).String()
	ctx := rpctypes.ContextWithHeight(height)
	res, err := queryClient.Account(ctx, &authtypes.QueryAccountRequest{Address: adr})
	if err != nil {
		st, ok := status.FromError(err)
		// treat as account doesn't exist yet
		if ok && st.Code() == codes.NotFound {
			b.logger.Info("getAccountNonce faild, account not found", "error", err)
			return 0, nil
		}
		return 0, err
	}
	var acc authtypes.AccountI
	if err := b.clientCtx.InterfaceRegistry.UnpackAny(res.Account, &acc); err != nil {
		return 0, err
	}

	nonce := acc.GetSequence()

	if !pending {
		return nonce, nil
	}

	// the account retriever doesn't include the uncommitted transactions on the nonce so we need to
	// to manually add them.
	pendingTxs, err := b.PendingTransactions()
	if err != nil {
		return nonce, nil
	}

	// add the uncommitted txs to the nonce counter
	// only supports `MsgEthereumTx` style tx
	for _, tx := range pendingTxs {
		for _, msg := range (*tx).GetMsgs() {
			ethMsg, ok := msg.(*txs.MsgEthereumTx)
			if !ok {
				// not ethereum tx
				break
			}

			sender, err := b.GetSender(ethMsg, b.chainID)
			if err != nil {
				continue
			}
			if sender == accAddr {
				nonce++
			}
		}
	}

	return nonce, nil
}
