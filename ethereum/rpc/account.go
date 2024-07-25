package rpc

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkmath "cosmossdk.io/math"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/artela/ethereum/crypto/ethsecp256k1"
	"github.com/artela-network/artela/ethereum/crypto/hd"
	ethapi2 "github.com/artela-network/artela/ethereum/rpc/ethapi"
	"github.com/artela-network/artela/ethereum/rpc/types"
	types2 "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/ethereum/utils"
	"github.com/artela-network/artela/x/evm/txs"
)

func (b *BackendImpl) Accounts() []common.Address {
	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	infos, err := b.clientCtx.Keyring.List()
	if err != nil {
		b.logger.Info("keying list failed", "error", err)
		return nil
	}

	for _, info := range infos {
		pubKey, err := info.GetPubKey()
		if err != nil {
			b.logger.Info("getPubKey failed", "info", info, "error", err)
			return nil
		}
		addressBytes := pubKey.Address().Bytes()
		addresses = append(addresses, common.BytesToAddress(addressBytes))
	}

	return addresses
}

func (b *BackendImpl) NewAccount(password string) (common.AddressEIP55, error) {
	name := "key_" + time.Now().UTC().Format(time.RFC3339)

	cfg := sdktypes.GetConfig()
	basePath := cfg.GetFullBIP44Path()

	hdPathIter, err := types2.NewHDPathIterator(basePath, true)
	if err != nil {
		b.logger.Info("NewHDPathIterator failed", "error", err)
		return common.AddressEIP55{}, err
	}
	// create the mnemonic and save the account
	hdPath := hdPathIter()

	info, _, err := b.clientCtx.Keyring.NewMnemonic(name, keyring.English, hdPath.String(), password, hd.EthSecp256k1)
	if err != nil {
		b.logger.Info("NewMnemonic failed", "error", err)
		return common.AddressEIP55{}, err
	}

	pubKey, err := info.GetPubKey()
	if err != nil {
		b.logger.Info("GetPubKey failed", "error", err)
		return common.AddressEIP55{}, err
	}
	addr := common.BytesToAddress(pubKey.Address().Bytes())
	return common.AddressEIP55(addr), nil
}

func (b *BackendImpl) ImportRawKey(privkey, password string) (common.Address, error) {
	priv, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return common.Address{}, err
	}

	privKey := &ethsecp256k1.PrivKey{Key: crypto.FromECDSA(priv)}

	addr := sdktypes.AccAddress(privKey.PubKey().Address().Bytes())
	ethereumAddr := common.BytesToAddress(addr)

	// return if the key has already been imported
	if _, err := b.clientCtx.Keyring.KeyByAddress(addr); err == nil {
		return ethereumAddr, nil
	}

	// ignore error as we only care about the length of the list
	list, _ := b.clientCtx.Keyring.List() // #nosec G703
	privKeyName := fmt.Sprintf("personal_%d", len(list))

	armor := sdkcrypto.EncryptArmorPrivKey(privKey, password, ethsecp256k1.KeyType)

	if err := b.clientCtx.Keyring.ImportPrivKey(privKeyName, armor, password); err != nil {
		return common.Address{}, err
	}

	return ethereumAddr, nil
}

func (b *BackendImpl) SignTransaction(args *ethapi2.TransactionArgs) (*ethtypes.Transaction, error) {
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

func (b *BackendImpl) getAccountNonce(accAddr common.Address, pending bool, height int64) (uint64, error) {
	queryClient := authtypes.NewQueryClient(b.clientCtx)
	adr := sdktypes.AccAddress(accAddr.Bytes()).String()
	ctx := types.ContextWithHeight(height)
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

func (b *BackendImpl) GetBalance(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	blockNum, err := b.blockNumberFromCosmos(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	req := &txs.QueryBalanceRequest{
		Address: address.String(),
	}

	_, err = b.CosmosBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	res, err := b.queryClient.Balance(types.ContextWithHeight(blockNum.Int64()), req)
	if err != nil {
		return nil, err
	}

	val, ok := sdkmath.NewIntFromString(res.Balance)
	if !ok {
		return nil, errors.New("invalid balance")
	}

	if val.IsNegative() {
		return nil, errors.New("couldn't fetch balance. Node state is pruned")
	}

	return (*hexutil.Big)(val.BigInt()), nil
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
		ctx := types.ContextWithHeight(int64(bn))

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
