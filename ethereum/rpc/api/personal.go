package api

import (
	"context"
	"errors"
	"fmt"

	rpctypes "github.com/artela-network/artela/ethereum/rpc/types"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

// PersonalAccountAPI provides an API to access accounts managed by this node.
// It offers methods to create, (un)lock en list accounts. Some methods accept
// passwords and are therefore considered private by default.
type PersonalAccountAPI struct {
	nonceLock *AddrLocker
	logger    log.Logger
	b         rpctypes.PersonalBackend
}

// NewPersonalAccountAPI create a new PersonalAccountAPI.
func NewPersonalAccountAPI(b rpctypes.PersonalBackend, logger log.Logger, nonceLock *AddrLocker) *PersonalAccountAPI {
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
	return errors.New("OpenWallet is not valid")
}

// DeriveAccount requests an HD wallet to derive a new account, optionally pinning
// it for later reuse.
func (s *PersonalAccountAPI) DeriveAccount(url string, path string, pin *bool) (accounts.Account, error) {
	return accounts.Account{}, errors.New("DeriveAccount is not valid")
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
	return false, errors.New("UnlockAccount is not valid")
}

// LockAccount will lock the account associated with the given address when it's unlocked.
func (s *PersonalAccountAPI) LockAccount(_ common.Address) bool {
	// not implemented"
	return false
}

// signTransaction sets defaults and signs the given transaction
// NOTE: the caller needs to ensure that the nonceLock is held, if applicable,
// and release it after the transaction has been submitted to the tx pool
func (s *PersonalAccountAPI) signTransaction(ctx context.Context, args *rpctypes.TransactionArgs, _ string) (*types.Transaction, error) {
	if err := args.SetDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	signed, err := s.b.SignTransaction(args)
	if err != nil {
		log.Warn("Failed transaction send attempt", "from", args.FromAddr(), "to", args.To, "value", args.Value.ToInt(), "err", err)
		return nil, err
	}
	return signed, err
}

// SendTransaction will create a transaction from the given arguments and
// tries to sign it with the key associated with args.From. If the given
// passwd isn't able to decrypt the key it fails.
func (s *PersonalAccountAPI) SendTransaction(ctx context.Context, args rpctypes.TransactionArgs, passwd string) (common.Hash, error) {
	if args.Nonce == nil {
		// Hold the mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		s.nonceLock.LockAddr(args.FromAddr())
		defer s.nonceLock.UnlockAddr(args.FromAddr())
	}

	signed, err := s.signTransaction(ctx, &args, passwd)
	if err != nil {
		log.Warn("Failed transaction send attempt", "from", args.FromAddr(), "to", args.To, "value", args.Value.ToInt(), "err", err)
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.logger, s.b, signed)
}

// SignTransaction will create a transaction from the given arguments and
// tries to sign it with the key associated with args.From. If the given passwd isn't
// able to decrypt the key it fails. The transaction is returned in RLP-form, not broadcast
// to other nodes
func (s *PersonalAccountAPI) SignTransaction(ctx context.Context, args rpctypes.TransactionArgs, passwd string) (*SignTransactionResult, error) {
	signed, err := s.signTransaction(ctx, &args, passwd)
	if err != nil {
		return nil, err
	}

	data, err := signed.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, signed}, nil
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
func (s *PersonalAccountAPI) Sign(_ context.Context, data hexutil.Bytes, addr common.Address, _ string) (hexutil.Bytes, error) {
	return s.b.Sign(addr, data)
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
	return "", errors.New("InitializeWallet is not valid")
}

// Unpair deletes a pairing between wallet and geth.
func (s *PersonalAccountAPI) Unpair(_ context.Context, _ string, _ string) error {
	return errors.New("unpair is not valid")
}
