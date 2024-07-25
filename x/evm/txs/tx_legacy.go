package txs

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common"
	ethereum "github.com/ethereum/go-ethereum/core/types"

	artela "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/ethereum/utils"
	"github.com/artela-network/artela/x/evm/types"
	"github.com/artela-network/aspect-core/djpm"
)

func newLegacyTx(tx *ethereum.Transaction) (*LegacyTx, error) {
	txData := &LegacyTx{
		Nonce:    tx.Nonce(),
		Data:     tx.Data(),
		GasLimit: tx.Gas(),
	}

	v, r, s := tx.RawSignatureValues()
	if to := tx.To(); to != nil {
		txData.To = to.Hex()
	}

	if tx.Value() != nil {
		amountInt, err := artela.SafeNewIntFromBigInt(tx.Value())
		if err != nil {
			return nil, err
		}
		txData.Amount = &amountInt
	}

	if tx.GasPrice() != nil {
		gasPriceInt, err := artela.SafeNewIntFromBigInt(tx.GasPrice())
		if err != nil {
			return nil, err
		}
		txData.GasPrice = &gasPriceInt
	}

	txData.SetSignatureValues(v, r, s)
	return txData, nil
}

// TxType returns the txs type
func (tx *LegacyTx) TxType() uint8 {
	return ethereum.LegacyTxType
}

// Copy returns an instance with the same field values
func (tx *LegacyTx) Copy() TxData {
	return &LegacyTx{
		Nonce:    tx.Nonce,
		GasPrice: tx.GasPrice,
		GasLimit: tx.GasLimit,
		To:       tx.To,
		Amount:   tx.Amount,
		Data:     common.CopyBytes(tx.Data),
		V:        common.CopyBytes(tx.V),
		R:        common.CopyBytes(tx.R),
		S:        common.CopyBytes(tx.S),
	}
}

// Validate performs a stateless validation of the txs fields.
func (tx LegacyTx) Validate() error {
	gasPrice := tx.GetGasPrice()
	if gasPrice == nil {
		return errorsmod.Wrap(types.ErrInvalidGasPrice, "gas price cannot be nil")
	}

	if gasPrice.Sign() == -1 {
		return errorsmod.Wrapf(types.ErrInvalidGasPrice, "gas price cannot be negative %s", gasPrice)
	}
	if !artela.IsValidInt256(gasPrice) {
		return errorsmod.Wrap(types.ErrInvalidGasPrice, "out of bound")
	}
	if !artela.IsValidInt256(tx.Fee()) {
		return errorsmod.Wrap(types.ErrInvalidGasFee, "out of bound")
	}

	amount := tx.GetValue()
	// Amount can be 0
	if amount != nil && amount.Sign() == -1 {
		return errorsmod.Wrapf(types.ErrInvalidAmount, "amount cannot be negative %s", amount)
	}
	if !artela.IsValidInt256(amount) {
		return errorsmod.Wrap(types.ErrInvalidAmount, "out of bound")
	}

	if tx.To != "" {
		if err := artela.ValidateAddress(tx.To); err != nil {
			return errorsmod.Wrap(err, "invalid to address")
		}
	}

	return nil
}

// GetChainID returns the chain id field from the derived signature values
func (tx *LegacyTx) GetChainID() *big.Int {
	v, _, _ := tx.GetRawSignatureValues()
	return DeriveChainID(v)
}

// GetAccessList returns nil
func (tx *LegacyTx) GetAccessList() ethereum.AccessList {
	return nil
}

// GetData returns a copy of the input data bytes.
func (tx *LegacyTx) GetData() []byte {
	return common.CopyBytes(tx.Data)
}

// GetGas returns the gas limit.
func (tx *LegacyTx) GetGas() uint64 {
	return tx.GasLimit
}

// GetGasPrice returns the gas price field.
func (tx *LegacyTx) GetGasPrice() *big.Int {
	if tx.GasPrice == nil {
		return nil
	}
	return tx.GasPrice.BigInt()
}

// GetGasTipCap returns the gas price field.
func (tx *LegacyTx) GetGasTipCap() *big.Int {
	return tx.GetGasPrice()
}

// GetGasFeeCap returns the gas price field.
func (tx *LegacyTx) GetGasFeeCap() *big.Int {
	return tx.GetGasPrice()
}

// Fee returns gasprice * gaslimit.
func (tx LegacyTx) Fee() *big.Int {
	return fee(tx.GetGasPrice(), tx.GetGas())
}

// Cost returns amount + gasprice * gaslimit.
func (tx LegacyTx) Cost() *big.Int {
	return cost(tx.Fee(), tx.GetValue())
}

// EffectiveGasPrice is the same as GasPrice for LegacyTx
func (tx LegacyTx) EffectiveGasPrice(_ *big.Int) *big.Int {
	return tx.GetGasPrice()
}

// EffectiveFee is the same as Fee for LegacyTx
func (tx LegacyTx) EffectiveFee(_ *big.Int) *big.Int {
	return tx.Fee()
}

// EffectiveCost is the same as Cost for LegacyTx
func (tx LegacyTx) EffectiveCost(_ *big.Int) *big.Int {
	return tx.Cost()
}

// GetValue returns the txs amount.
func (tx *LegacyTx) GetValue() *big.Int {
	if tx.Amount == nil {
		return nil
	}
	return tx.Amount.BigInt()
}

// GetNonce returns the account sequence for the txs.
func (tx *LegacyTx) GetNonce() uint64 { return tx.Nonce }

// GetTo returns the pointer to the recipient address.
func (tx *LegacyTx) GetTo() *common.Address {
	if tx.To == "" {
		return nil
	}
	to := common.HexToAddress(tx.To)
	return &to
}

// AsEthereumData returns an AccessListTx txs txs from the proto-formatted
// TxData defined on the Cosmos EVM.
func (tx *LegacyTx) AsEthereumData(stripCallData bool) ethereum.TxData {
	v, r, s := tx.GetRawSignatureValues()
	txData := &ethereum.LegacyTx{
		Nonce:    tx.GetNonce(),
		GasPrice: tx.GetGasPrice(),
		Gas:      tx.GetGas(),
		To:       tx.GetTo(),
		Value:    tx.GetValue(),
		Data:     tx.GetData(),
		V:        v,
		R:        r,
		S:        s,
	}

	if stripCallData && utils.IsCustomizedVerification(ethereum.NewTx(txData)) {
		_, txData.Data, _ = djpm.DecodeValidationAndCallData(tx.Data)
	}

	return txData
}

// GetRawSignatureValues returns the V, R, S signature values of the txs.
// The return values should not be modified by the caller.
func (tx *LegacyTx) GetRawSignatureValues() (v, r, s *big.Int) {
	return rawSignatureValues(tx.V, tx.R, tx.S)
}

// SetSignatureValues sets the signature values to the txs.
func (tx *LegacyTx) SetSignatureValues(v, r, s *big.Int) {
	if v != nil {
		tx.V = v.Bytes()
	}
	if r != nil {
		tx.R = r.Bytes()
	}
	if s != nil {
		tx.S = s.Bytes()
	}
}

func (tx *LegacyTx) SetChainId(_ *big.Int) {
}
