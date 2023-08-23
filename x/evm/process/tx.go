package process

import (
	sdkmath "cosmossdk.io/math"
	"math"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// NewTx returns a reference to a new Ethereum process message.
func NewTx(
	tx *EvmTxArgs,
) *MsgEthereumTx {
	return newMsgEthereumTx(tx)
}

func newMsgEthereumTx(
	tx *EvmTxArgs,
) *MsgEthereumTx {
	var (
		cid, amt, gp *sdkmath.Int
		toAddr       string
		txData       TxData
	)

	if tx.To != nil {
		toAddr = tx.To.Hex()
	}

	if tx.Amount != nil {
		amountInt := sdkmath.NewIntFromBigInt(tx.Amount)
		amt = &amountInt
	}

	if tx.ChainID != nil {
		chainIDInt := sdkmath.NewIntFromBigInt(tx.ChainID)
		cid = &chainIDInt
	}

	if tx.GasPrice != nil {
		gasPriceInt := sdkmath.NewIntFromBigInt(tx.GasPrice)
		gp = &gasPriceInt
	}

	switch {
	case tx.Accesses == nil:
		txData = &LegacyTx{
			To:       toAddr,
			Amount:   amt,
			GasPrice: gp,
			Nonce:    tx.Nonce,
			GasLimit: tx.GasLimit,
			Data:     tx.Input,
		}
	case tx.Accesses != nil && tx.GasFeeCap != nil && tx.GasTipCap != nil:
		gtc := sdkmath.NewIntFromBigInt(tx.GasTipCap)
		gfc := sdkmath.NewIntFromBigInt(tx.GasFeeCap)

		txData = &DynamicFeeTx{
			ChainID:   cid,
			Amount:    amt,
			To:        toAddr,
			GasTipCap: &gtc,
			GasFeeCap: &gfc,
			Nonce:     tx.Nonce,
			GasLimit:  tx.GasLimit,
			Data:      tx.Input,
			Accesses:  NewAccessList(tx.Accesses),
		}
	case tx.Accesses != nil:
		txData = &AccessListTx{
			ChainID:  cid,
			Nonce:    tx.Nonce,
			To:       toAddr,
			Amount:   amt,
			GasLimit: tx.GasLimit,
			GasPrice: gp,
			Data:     tx.Input,
			Accesses: NewAccessList(tx.Accesses),
		}
	default:
	}

	dataAny, err := PackTxData(txData)
	if err != nil {
		panic(err)
	}

	msg := MsgEthereumTx{Data: dataAny}
	msg.Hash = msg.AsTransaction().Hash().Hex()
	return &msg
}

// EvmTxArgs encapsulates all possible params to create all EVM txs types.
// This includes LegacyTx, DynamicFeeTx and AccessListTx
type EvmTxArgs struct {
	Nonce     uint64
	GasLimit  uint64
	Input     []byte
	GasFeeCap *big.Int
	GasPrice  *big.Int
	ChainID   *big.Int
	Amount    *big.Int
	GasTipCap *big.Int
	To        *common.Address
	Accesses  *ethereum.AccessList
}

// GetTxPriority returns the priority of a given Ethereum process. It relies of the
// priority reduction global variable to calculate the process priority given the process
// tip price:
//
//	tx_priority = tip_price / priority_reduction
func GetTxPriority(txData TxData, baseFee *big.Int) (priority int64) {
	// calculate priority based on effective gas price
	tipPrice := txData.EffectiveGasPrice(baseFee)
	// if london hard fork is not enabled, tipPrice is the gasPrice
	if baseFee != nil {
		tipPrice = new(big.Int).Sub(tipPrice, baseFee)
	}

	priority = math.MaxInt64
	priorityBig := new(big.Int).Quo(tipPrice, DefaultPriorityReduction.BigInt())

	// safety check
	if priorityBig.IsInt64() {
		priority = priorityBig.Int64()
	}

	return priority
}

// Failed returns if the contract execution failed in vm errors
func (m *MsgEthereumTxResponse) Failed() bool {
	return len(m.VmError) > 0
}

// Return is a helper function to help caller distinguish between revert reason
// and function return. Return returns the data after execution if no error occurs.
func (m *MsgEthereumTxResponse) Return() []byte {
	if m.Failed() {
		return nil
	}
	return common.CopyBytes(m.Ret)
}

// Revert returns the concrete revert reason if the execution is aborted by `REVERT`
// opcode. Note the reason can be nil if no data supplied with revert opcode.
func (m *MsgEthereumTxResponse) Revert() []byte {
	if m.VmError != vm.ErrExecutionReverted.Error() {
		return nil
	}
	return common.CopyBytes(m.Ret)
}
