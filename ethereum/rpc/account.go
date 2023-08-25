package rpc

import (
	"context"
	"fmt"
	ethapi2 "github.com/artela-network/artela/ethereum/rpc/ethapi"
	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	types2 "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/txs"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/artela-network/artela/ethereum/crypto/ethsecp256k1"
	"github.com/artela-network/artela/ethereum/crypto/hd"
)

var _ ethapi2.AccountBackend = (*AccountBackend)(nil)

type AccountBackend struct {
	ctx         context.Context
	clientCtx   client.Context
	chainID     *big.Int
	queryClient *rpctypes.QueryClient
}

func NewAccountBackend(ctx context.Context, clientCtx client.Context, queryClient *rpctypes.QueryClient) *AccountBackend {
	chainID, err := types2.ParseChainID(clientCtx.ChainID)
	if err != nil {
		panic(err)
	}

	return &AccountBackend{
		ctx:         ctx,
		clientCtx:   clientCtx,
		chainID:     chainID,
		queryClient: queryClient,
	}
}

func (ab *AccountBackend) Accounts() []common.Address {
	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	infos, err := ab.clientCtx.Keyring.List()
	if err != nil {
		return nil
	}

	for _, info := range infos {
		pubKey, err := info.GetPubKey()
		if err != nil {
			return nil
		}
		addressBytes := pubKey.Address().Bytes()
		addresses = append(addresses, common.BytesToAddress(addressBytes))
	}

	return addresses
}

func (ab *AccountBackend) NewAccount(password string) (common.AddressEIP55, error) {
	name := "key_" + time.Now().UTC().Format(time.RFC3339)

	cfg := sdktypes.GetConfig()
	basePath := cfg.GetFullBIP44Path()

	hdPathIter, err := types2.NewHDPathIterator(basePath, true)
	if err != nil {
		panic(err)
	}
	// create the mnemonic and save the account
	hdPath := hdPathIter()

	info, _, err := ab.clientCtx.Keyring.NewMnemonic(name, keyring.English, hdPath.String(), password, hd.EthSecp256k1)
	if err != nil {
		return common.AddressEIP55{}, err
	}

	pubKey, err := info.GetPubKey()
	if err != nil {
		return common.AddressEIP55{}, err
	}
	addr := common.BytesToAddress(pubKey.Address().Bytes())
	return common.AddressEIP55(addr), nil
}

func (ab *AccountBackend) ImportRawKey(privkey, password string) (common.Address, error) {
	priv, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return common.Address{}, err
	}

	privKey := &ethsecp256k1.PrivKey{Key: crypto.FromECDSA(priv)}

	addr := sdktypes.AccAddress(privKey.PubKey().Address().Bytes())
	ethereumAddr := common.BytesToAddress(addr)

	// return if the key has already been imported
	if _, err := ab.clientCtx.Keyring.KeyByAddress(addr); err == nil {
		return ethereumAddr, nil
	}

	// ignore error as we only care about the length of the list
	list, _ := ab.clientCtx.Keyring.List() // #nosec G703
	privKeyName := fmt.Sprintf("personal_%d", len(list))

	armor := sdkcrypto.EncryptArmorPrivKey(privKey, password, ethsecp256k1.KeyType)

	if err := ab.clientCtx.Keyring.ImportPrivKey(privKeyName, armor, password); err != nil {
		return common.Address{}, err
	}

	return ethereumAddr, nil
}

func (ab *AccountBackend) SignTransaction(args *ethapi2.TransactionArgs) (*ethtypes.Transaction, error) {
	_, err := ab.clientCtx.Keyring.KeyByAddress(sdktypes.AccAddress(args.From.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("failed to find key in the node's keyring; %s; %s", keystore.ErrNoMatch, err.Error())
	}

	if args.ChainID != nil && (ab.chainID).Cmp((*big.Int)(args.ChainID)) != 0 {
		return nil, fmt.Errorf("chainId does not match node's (have=%v, want=%v)", args.ChainID, (*hexutil.Big)(ab.chainID))
	}

	// TODO, set defaults
	// args, err = b.SetTxDefaults(args)
	// if err != nil {
	// 	return common.Hash{}, err
	// }

	bn, err := ab.BlockNumber()
	if err != nil {
		return nil, err
	}

	bt, err := ab.BlockTimeByNumber(int64(bn))
	if err != nil {
		return nil, err
	}

	signer := ethtypes.MakeSigner(ab.ChainConfig(), new(big.Int).SetUint64(uint64(bn)), bt)

	// LegacyTx derives chainID from the signature. To make sure the msg.ValidateBasic makes
	// the corresponding chainID validation, we need to sign the transaction before calling it

	// Sign transaction
	msg := args.ToEVMTransaction()
	return msg.SignEthereumTx(signer, ab.clientCtx.Keyring)
}

// Sign signs the provided data using the private key of address via Geth's signature standard.
func (ab *AccountBackend) Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	from := sdktypes.AccAddress(address.Bytes())

	_, err := ab.clientCtx.Keyring.KeyByAddress(from)
	if err != nil {
		return nil, fmt.Errorf("%s; %s", keystore.ErrNoMatch, err.Error())
	}

	// Sign the requested hash with the wallet
	signature, _, err := ab.clientCtx.Keyring.SignByAddress(from, data)
	if err != nil {
		return nil, err
	}

	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}

func (ab *AccountBackend) BlockNumber() (hexutil.Uint64, error) {
	// do any grpc query, ignore the response and use the returned block height
	var header metadata.MD
	_, err := ab.queryClient.Params(ab.ctx, &txs.QueryParamsRequest{}, grpc.Header(&header))
	if err != nil {
		return hexutil.Uint64(0), err
	}

	blockHeightHeader := header.Get(grpctypes.GRPCBlockHeightHeader)
	if headerLen := len(blockHeightHeader); headerLen != 1 {
		return 0, fmt.Errorf("unexpected '%s' gRPC header length; got %d, expected: %d", grpctypes.GRPCBlockHeightHeader, headerLen, 1)
	}

	height, err := strconv.ParseUint(blockHeightHeader[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block height: %w", err)
	}

	if height > math.MaxInt64 {
		return 0, fmt.Errorf("block height %d is greater than max uint64", height)
	}

	return hexutil.Uint64(height), nil
}

func (ab *AccountBackend) BlockTimeByNumber(blockNum int64) (uint64, error) {
	resBlock, err := ab.clientCtx.Client.Block(ab.ctx, &blockNum)
	if err != nil {
		return 0, err
	}
	return uint64(resBlock.Block.Time.Unix()), nil
}

// ChainConfig returns the latest ethereum chain configuration
func (ab *AccountBackend) ChainConfig() *params.ChainConfig {
	params, err := ab.queryClient.Params(ab.ctx, &txs.QueryParamsRequest{})
	if err != nil {
		return nil
	}

	return params.Params.ChainConfig.EthereumConfig(ab.chainID)
}
