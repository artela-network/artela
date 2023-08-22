package rpc

import (
	"fmt"
	"math/big"
	"time"

	"github.com/artela-network/artela/ethereum/crypto/ethsecp256k1"
	"github.com/artela-network/artela/ethereum/crypto/hd"
	"github.com/artela-network/artela/rpc/ethapi"
	"github.com/artela-network/artela/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var _ ethapi.AccountBackend = (*AccountBackend)(nil)

type AccountBackend struct {
	clientCtx client.Context
	chainID   *big.Int
}

func NewAccountBackend(clientCtx client.Context) *AccountBackend {
	chainID, err := types.ParseChainID(clientCtx.ChainID)
	if err != nil {
		panic(err)
	}

	return &AccountBackend{
		clientCtx: clientCtx,
		chainID:   chainID,
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

	hdPathIter, err := types.NewHDPathIterator(basePath, true)
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

func (ab *AccountBackend) SignTransaction(args *ethapi.TransactionArgs, passwd string) (*ethtypes.Transaction, error) {
	/*_, err := ab.clientCtx.Keyring.KeyByAddress(sdktypes.AccAddress(args.From.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("failed to find key in the node's keyring; %s; %s", keystore.ErrNoMatch, err.Error())
	}

	if args.ChainID != nil && (ab.chainID).Cmp((*big.Int)(args.ChainID)) != 0 {
		return nil, fmt.Errorf("chainId does not match node's (have=%v, want=%v)", args.ChainID, (*hexutil.Big)(b.chainID))
	}

	args, err = b.SetTxDefaults(args)
	if err != nil {
		return common.Hash{}, err
	}

	bn, err := b.BlockNumber()
	if err != nil {
		b.logger.Debug("failed to fetch latest block number", "error", err.Error())
		return common.Hash{}, err
	}

	signer := ethtypes.MakeSigner(b.ChainConfig(), new(big.Int).SetUint64(uint64(bn)))

	// LegacyTx derives chainID from the signature. To make sure the msg.ValidateBasic makes
	// the corresponding chainID validation, we need to sign the transaction before calling it

	// Sign transaction
	msg := args.ToTransaction()
	if err := msg.Sign(signer, b.clientCtx.Keyring); err != nil {
		b.logger.Debug("failed to sign tx", "error", err.Error())
		return common.Hash{}, err
	}

	if err := msg.ValidateBasic(); err != nil {
		b.logger.Debug("tx failed basic validation", "error", err.Error())
		return common.Hash{}, err
	}

	// Query params to use the EVM denomination
	res, err := b.queryClient.QueryClient.Params(b.ctx, &evmtypes.QueryParamsRequest{})
	if err != nil {
		b.logger.Error("failed to query evm params", "error", err.Error())
		return common.Hash{}, err
	}

	// Assemble transaction from fields
	tx, err := msg.BuildTx(b.clientCtx.TxConfig.NewTxBuilder(), res.Params.EvmDenom)
	if err != nil {
		b.logger.Error("build cosmos tx failed", "error", err.Error())
		return common.Hash{}, err
	}

	// Encode transaction by default Tx encoder
	txEncoder := b.clientCtx.TxConfig.TxEncoder()
	txBytes, err := txEncoder(tx)
	if err != nil {
		b.logger.Error("failed to encode eth tx using default encoder", "error", err.Error())
		return common.Hash{}, err
	}

	ethTx := msg.AsTransaction()
	return ethTx, nil*/
	return nil, nil
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
