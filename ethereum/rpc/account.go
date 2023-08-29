package rpc

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	ethapi2 "github.com/artela-network/artela/ethereum/rpc/ethapi"
	types2 "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/txs"

	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/artela-network/artela/ethereum/crypto/ethsecp256k1"
	"github.com/artela-network/artela/ethereum/crypto/hd"
)

func (b *backend) Accounts() []common.Address {
	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	infos, err := b.clientCtx.Keyring.List()
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

func (b *backend) NewAccount(password string) (common.AddressEIP55, error) {
	name := "key_" + time.Now().UTC().Format(time.RFC3339)

	cfg := sdktypes.GetConfig()
	basePath := cfg.GetFullBIP44Path()

	hdPathIter, err := types2.NewHDPathIterator(basePath, true)
	if err != nil {
		panic(err)
	}
	// create the mnemonic and save the account
	hdPath := hdPathIter()

	info, _, err := b.clientCtx.Keyring.NewMnemonic(name, keyring.English, hdPath.String(), password, hd.EthSecp256k1)
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

func (b *backend) ImportRawKey(privkey, password string) (common.Address, error) {
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

func (b *backend) SignTransaction(args *ethapi2.TransactionArgs) (*ethtypes.Transaction, error) {
	_, err := b.clientCtx.Keyring.KeyByAddress(sdktypes.AccAddress(args.From.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("failed to find key in the node's keyring; %s; %s", keystore.ErrNoMatch, err.Error())
	}

	if args.ChainID != nil && (b.chainID).Cmp((*big.Int)(args.ChainID)) != 0 {
		return nil, fmt.Errorf("chainId does not match node's (have=%v, want=%v)", args.ChainID, (*hexutil.Big)(b.chainID))
	}

	// TODO, set defaults
	// args, err = b.SetTxDefaults(args)
	// if err != nil {
	// 	return common.Hash{}, err
	// }

	bn, err := b.BlockNumber()
	if err != nil {
		return nil, err
	}

	bt, err := b.BlockTimeByNumber(int64(bn))
	if err != nil {
		return nil, err
	}

	signer := ethtypes.MakeSigner(b.ChainConfig(), new(big.Int).SetUint64(uint64(bn)), bt)

	// LegacyTx derives chainID from the signature. To make sure the msg.ValidateBasic makes
	// the corresponding chainID validation, we need to sign the transaction before calling it

	// Sign transaction
	msg := args.ToEVMTransaction()
	return msg.SignEthereumTx(signer, b.clientCtx.Keyring)
}

// Sign signs the provided data using the private key of address via Geth's signature standard.
func (b *backend) Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
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

func (b *backend) BlockNumber() (hexutil.Uint64, error) {
	// do any grpc query, ignore the response and use the returned block height
	var header metadata.MD
	_, err := b.queryClient.Params(b.ctx, &txs.QueryParamsRequest{}, grpc.Header(&header))
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

func (b *backend) BlockTimeByNumber(blockNum int64) (uint64, error) {
	resBlock, err := b.clientCtx.Client.Block(b.ctx, &blockNum)
	if err != nil {
		return 0, err
	}
	return uint64(resBlock.Block.Time.Unix()), nil
}

func (b *backend) GetTransactionCount(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	return nil, nil
}
