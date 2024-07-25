package txs

import (
	"errors"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	codec "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/common"
	cmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	ethereum "github.com/ethereum/go-ethereum/core/types"

	artela "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/ethereum/utils"
	"github.com/artela-network/artela/x/evm/types"
)

var (
	_ cosmos.Msg = &MsgEthereumTx{}
	_ cosmos.Tx  = &MsgEthereumTx{}
	_ ante.GasTx = &MsgEthereumTx{}
	_ cosmos.Msg = &MsgUpdateParams{}

	_ codec.UnpackInterfacesMessage = MsgEthereumTx{}
)

// ===============================================================
//          		      MsgUpdateParams
// ===============================================================

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m MsgUpdateParams) GetSigners() []cosmos.AccAddress {
	// #nosec G703 -- gosec raises a warning about a non-handled error which we deliberately ignore here
	addr, _ := cosmos.AccAddressFromBech32(m.Authority)
	return []cosmos.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := cosmos.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	return m.Params.Validate()
}

// GetSignBytes implements the LegacyMsg interface.
func (m MsgUpdateParams) GetSignBytes() []byte {
	return cosmos.MustSortJSON(AminoCdc.MustMarshalJSON(&m))
}

// ===============================================================
//          		      MsgEthereumTx
// ===============================================================

// AsTransaction creates an Ethereum Transaction type from the msg fields
func (msg MsgEthereumTx) AsTransaction() *ethereum.Transaction {
	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return nil
	}

	return ethereum.NewTx(txData.AsEthereumData(false))
}

func (msg MsgEthereumTx) AsEthCallTransaction() *ethereum.Transaction {
	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return nil
	}

	return ethereum.NewTx(txData.AsEthereumData(true))
}

func ToMessage(tx *ethereum.Transaction, signer ethereum.Signer, baseFee *big.Int) (*core.Message, error) {
	message := &core.Message{
		Nonce:             tx.Nonce(),
		GasLimit:          tx.Gas(),
		GasPrice:          new(big.Int).Set(tx.GasPrice()),
		GasFeeCap:         new(big.Int).Set(tx.GasFeeCap()),
		GasTipCap:         new(big.Int).Set(tx.GasTipCap()),
		To:                tx.To(),
		Value:             tx.Value(),
		Data:              tx.Data(),
		AccessList:        tx.AccessList(),
		SkipAccountChecks: false,
	}
	// If baseFee provided, set gasPrice to effectiveGasPrice.
	if baseFee != nil {
		message.GasPrice = cmath.BigMin(message.GasPrice.Add(message.GasTipCap, baseFee), message.GasFeeCap)
	}
	var err error
	message.From, err = ethereum.Sender(signer, tx)
	return message, err
}

// AsMessage creates an Ethereum core.Message from the msg fields
func (msg MsgEthereumTx) AsMessage(signer ethereum.Signer, baseFee *big.Int) (*core.Message, error) {
	tx := msg.AsTransaction()
	return ToMessage(tx, signer, baseFee)
}

// UnpackInterfaces implements UnpackInterfacesMesssage.UnPackInterfaces
func (msg MsgEthereumTx) UnpackInterfaces(unpacker codec.AnyUnpacker) error {
	return unpacker.UnpackAny(msg.Data, new(TxData))
}

// UnmarshalBinary decodes the canonical encoding of transactions.
func (msg *MsgEthereumTx) UnmarshalBinary(b []byte) error {
	tx := &ethereum.Transaction{}
	if err := tx.UnmarshalBinary(b); err != nil {
		return err
	}
	return msg.FromEthereumTx(tx)
}

// BuildTx builds the canonical cosmos tx from ethereum msg
func (msg *MsgEthereumTx) BuildTx(b client.TxBuilder, evmDenom string) (signing.Tx, error) {
	builder, ok := b.(authtx.ExtensionOptionsTxBuilder)
	if !ok {
		return nil, errors.New("unsupported builder")
	}

	option, err := codec.NewAnyWithValue(&ExtensionOptionsEthereumTx{})
	if err != nil {
		return nil, err
	}

	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return nil, err
	}
	fees := make(cosmos.Coins, 0)
	feeAmt := sdkmath.NewIntFromBigInt(txData.Fee())
	if feeAmt.Sign() > 0 {
		fees = append(fees, cosmos.NewCoin(evmDenom, feeAmt))
	}

	builder.SetExtensionOptions(option)

	// A valid msg should have empty `From`
	msg.From = ""

	err = builder.SetMsgs(msg)
	if err != nil {
		return nil, err
	}
	builder.SetFeeAmount(fees)
	builder.SetGasLimit(msg.GetGas())
	tx := builder.GetTx()
	return tx, nil
}

// Route returns the route value of an MsgEthereumTx.
func (msg MsgEthereumTx) Route() string { return types.RouterKey }

// Type returns the type value of an MsgEthereumTx.
func (msg MsgEthereumTx) Type() string { return types.TypeMsgEthereumTx }

// ValidateBasic implements the cosmos.Msg interface. It performs basic validation
// checks of a Transaction. If returns an error if validation fails.
func (msg MsgEthereumTx) ValidateBasic() error {
	if msg.From != "" {
		if err := artela.ValidateAddress(msg.From); err != nil {
			return errorsmod.Wrap(err, "invalid from address")
		}
	}

	// Validate Size_ field, should be kept empty
	if msg.Size_ != 0 {
		return errorsmod.Wrapf(errortypes.ErrInvalidRequest, "txs size is deprecated")
	}

	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return errorsmod.Wrap(err, "failed to unpack txs data")
	}

	gas := txData.GetGas()

	// prevent txs with 0 gas to fill up the mempool
	if gas == 0 {
		return errorsmod.Wrap(types.ErrInvalidGasLimit, "gas limit must not be zero")
	}

	// prevent gas limit from overflow
	if g := new(big.Int).SetUint64(gas); !g.IsInt64() {
		return errorsmod.Wrap(types.ErrGasOverflow, "gas limit must be less than math.MaxInt64")
	}

	if err := txData.Validate(); err != nil {
		return err
	}

	// Validate Hash field after validated txData to avoid panic
	txHash := msg.AsTransaction().Hash().Hex()
	if msg.Hash != txHash {
		return errorsmod.Wrapf(errortypes.ErrInvalidRequest, "invalid txs hash %s, expected: %s", msg.Hash, txHash)
	}

	return nil
}

// Sign calculates a secp256k1 ECDSA signature and signs the  It
// takes a keyring signer and the chainID to sign an Ethereum txs according to
// EIP155 standard.
// This method mutates the txs as it populates the V, R, S
// fields of the Transaction's Signature.
// The function will fail if the sender address is not defined for the msg or if
// the sender is not registered on the keyring
func (msg *MsgEthereumTx) Sign(ethSigner ethereum.Signer, keyringSigner keyring.Signer) error {
	tx, err := msg.SignEthereumTx(ethSigner, keyringSigner)
	if err != nil {
		return err
	}

	return msg.FromEthereumTx(tx)
}

func (msg *MsgEthereumTx) SignEthereumTx(ethSigner ethereum.Signer, keyringSigner keyring.Signer) (*ethereum.Transaction, error) {
	from := msg.GetFrom()
	if from.Empty() {
		return nil, fmt.Errorf("sender address not defined for message")
	}

	tx := msg.AsTransaction()
	txHash := ethSigner.Hash(tx)

	sig, _, err := keyringSigner.SignByAddress(from, txHash.Bytes())
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(ethSigner, sig)
}

// GetMsgs returns a single MsgEthereumTx as an cosmos.Msg.
func (msg *MsgEthereumTx) GetMsgs() []cosmos.Msg {
	return []cosmos.Msg{msg}
}

// GetSigners returns the expected signers for an Ethereum txs message.
// For such a message, there should exist only a single 'signer'.
//
// NOTE: This method panics if 'Sign' hasn't been called first.
func (msg *MsgEthereumTx) GetSigners() []cosmos.AccAddress {
	data, err := UnpackTxData(msg.Data)
	if err != nil {
		panic(err)
	}

	// ctx := rpctypes.ContextWithHeight(blockNum.Int64())
	sender, err := msg.GetSender(data.GetChainID())
	if err != nil {
		panic(err)
	}

	signer := cosmos.AccAddress(sender.Bytes())
	return []cosmos.AccAddress{signer}
}

// GetSignBytes returns the Amino bytes of an Ethereum txs message used
// for signing.
//
// NOTE: This method cannot be used as a chain ID is needed to create valid bytes
// to sign over. Use 'RLPSignBytes' instead.
func (msg MsgEthereumTx) GetSignBytes() []byte {
	panic("must use 'RLPSignBytes' with a chain ID to get the valid bytes to sign")
}

// GetGas implements the GasTx interface. It returns the GasLimit of the
func (msg MsgEthereumTx) GetGas() uint64 {
	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return 0
	}
	return txData.GetGas()
}

// GetFee returns the fee for non dynamic fee txs
func (msg MsgEthereumTx) GetFee() *big.Int {
	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return nil
	}
	return txData.Fee()
}

// GetEffectiveFee returns the fee for dynamic fee txs
func (msg MsgEthereumTx) GetEffectiveFee(baseFee *big.Int) *big.Int {
	txData, err := UnpackTxData(msg.Data)
	if err != nil {
		return nil
	}
	return txData.EffectiveFee(baseFee)
}

// GetFrom loads the ethereum sender address from the sigcache and returns an
// cosmos.AccAddress from its bytes
func (msg *MsgEthereumTx) GetFrom() cosmos.AccAddress {
	if msg.From == "" {
		return nil
	}

	return common.HexToAddress(msg.From).Bytes()
}

// GetSender extracts the sender address from the signature values using the latest signer for the given chainID.
func (msg *MsgEthereumTx) GetSender(chainID *big.Int) (from common.Address, err error) {
	if msg.From != "" {
		return common.HexToAddress(msg.From), nil
	}

	tx := msg.AsTransaction()
	// retrieve sender info from aspect if tx is not signed
	if utils.IsCustomizedVerification(tx) {
		// TODO, more checkings, should never reach here
		return common.Address{}, errors.New("failed to get sender of customized tx")
	}

	signer := ethereum.LatestSignerForChainID(chainID)
	from, err = signer.Sender(tx)
	if err != nil {
		return common.Address{}, err
	}

	msg.From = from.Hex()
	return from, nil
}

// FromEthereumTx populates the message fields from the given ethereum transaction
func (msg *MsgEthereumTx) FromEthereumTx(tx *ethereum.Transaction) error {
	txData, err := NewTxDataFromTx(tx)
	if err != nil {
		return err
	}

	anyTxData, err := PackTxData(txData)
	if err != nil {
		return err
	}

	msg.Data = anyTxData
	msg.Hash = tx.Hash().Hex()
	return nil
}

func NewTxDataFromTx(tx *ethereum.Transaction) (TxData, error) {
	var txData TxData
	var err error
	switch tx.Type() {
	case ethereum.DynamicFeeTxType:
		txData, err = newDynamicFeeTx(tx)
	case ethereum.AccessListTxType:
		txData, err = newAccessListTx(tx)
	default:
		txData, err = newLegacyTx(tx)
	}
	if err != nil {
		return nil, err
	}

	return txData, nil
}
