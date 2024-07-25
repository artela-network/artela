package ethapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artela-network/artela/x/evm/txs"
)

// TransactionArgs represents the arguments to construct a new transaction
// or a message call.
type TransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	// Issue detail: https://github.com/ethereum/go-ethereum/issues/15628
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`
}

// from retrieves the transaction sender address.
func (args *TransactionArgs) from() common.Address {
	if args.From == nil {
		return common.Address{}
	}
	return *args.From
}

// data retrieves the transaction calldata. Input field is preferred.
func (args *TransactionArgs) data() []byte {
	if args.Input != nil {
		return *args.Input
	}
	if args.Data != nil {
		return *args.Data
	}
	return nil
}

// setDefaults fills in default values for unspecified tx fields.
func (args *TransactionArgs) setDefaults(ctx context.Context, b Backend) error {
	if err := args.setFeeDefaults(ctx, b); err != nil {
		return err
	}
	if args.Value == nil {
		args.Value = new(hexutil.Big)
	}
	if args.Nonce == nil {
		nonce, err := b.GetTransactionCount(args.from(), rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))
		if err != nil {
			return err
		}
		args.Nonce = nonce
		// TODO set nonce
	}
	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
		return errors.New(`both "data" and "input" are set and not equal. Please use "input" to pass transaction call data`)
	}
	if args.To == nil && len(args.data()) == 0 {
		return errors.New(`contract creation without any data provided`)
	}
	// Estimate the gas usage if necessary.
	if args.Gas == nil {
		// These fields are immutable during the estimation, safe to
		// pass the pointer directly.
		data := args.data()
		callArgs := TransactionArgs{
			From:                 args.From,
			To:                   args.To,
			GasPrice:             args.GasPrice,
			MaxFeePerGas:         args.MaxFeePerGas,
			MaxPriorityFeePerGas: args.MaxPriorityFeePerGas,
			Value:                args.Value,
			Data:                 (*hexutil.Bytes)(&data),
			AccessList:           args.AccessList,
		}
		pendingBlockNr := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
		// estimated, err := DoEstimateGas(ctx, b, callArgs, pendingBlockNr, b.RPCGasCap())
		// if err != nil {
		// 	return err
		// }
		// args.Gas = &estimated
		// TODO set gas
		_ = callArgs
		_ = pendingBlockNr

		log.Trace("Estimate gas usage automatically", "gas", args.Gas)
	}
	// If chain id is provided, ensure it matches the local chain id. Otherwise, set the local
	// chain id as the default.
	want := b.ChainConfig().ChainID
	if args.ChainID != nil {
		if have := (*big.Int)(args.ChainID); have.Cmp(want) != 0 {
			return fmt.Errorf("chainId does not match node's (have=%v, want=%v)", have, want)
		}
	} else {
		args.ChainID = (*hexutil.Big)(want)
	}
	return nil
}

// setFeeDefaults fills in default fee values for unspecified tx fields.
func (args *TransactionArgs) setFeeDefaults(ctx context.Context, b Backend) error {
	// If both gasPrice and at least one of the EIP-1559 fee parameters are specified, error.
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}
	// If the tx has completely specified a fee mechanism, no default is needed. This allows users
	// who are not yet synced past London to get defaults for other tx values. See
	// https://github.com/ethereum/go-ethereum/pull/23274 for more information.
	eip1559ParamsSet := args.MaxFeePerGas != nil && args.MaxPriorityFeePerGas != nil
	if (args.GasPrice != nil && !eip1559ParamsSet) || (args.GasPrice == nil && eip1559ParamsSet) {
		// Sanity check the EIP-1559 fee parameters if present.
		if args.GasPrice == nil && args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
			return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
		}
		return nil
	}
	// Now attempt to fill in default value depending on whether London is active or not.
	head, err := b.CurrentHeader()
	if err != nil {
		return err
	}

	if b.ChainConfig().IsLondon(head.Number) {
		// London is active, set maxPriorityFeePerGas and maxFeePerGas.
		if err := args.setLondonFeeDefaults(ctx, head, b); err != nil {
			return err
		}
	} else {
		if args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil {
			return errors.New("maxFeePerGas and maxPriorityFeePerGas are not valid before London is active")
		}
		// London not active, set gas price.
		price, err := b.SuggestGasTipCap(head.BaseFee)
		if err != nil {
			return err
		}
		args.GasPrice = (*hexutil.Big)(price)
	}
	return nil
}

// setLondonFeeDefaults fills in reasonable default fee values for unspecified fields.
func (args *TransactionArgs) setLondonFeeDefaults(_ context.Context, head *types.Header, b Backend) error {
	// Set maxPriorityFeePerGas if it is missing.
	if args.MaxPriorityFeePerGas == nil {
		tip, err := b.SuggestGasTipCap(head.BaseFee)
		if err != nil {
			return err
		}
		args.MaxPriorityFeePerGas = (*hexutil.Big)(tip)
	}
	// Set maxFeePerGas if it is missing.
	if args.MaxFeePerGas == nil {
		// Set the max fee to be 2 times larger than the previous block's base fee.
		// The additional slack allows the tx to not become invalidated if the base
		// fee is rising.
		val := new(big.Int).Add(
			args.MaxPriorityFeePerGas.ToInt(),
			new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
		)
		args.MaxFeePerGas = (*hexutil.Big)(val)
	}
	// Both EIP-1559 fee parameters are now set; sanity check them.
	if args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
		return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
	}
	return nil
}

// ToMessage converts the transaction arguments to the Message type used by the
// core evm. This method is used in calls and traces that do not require a real
// live transaction.
func (args *TransactionArgs) ToMessage(globalGasCap uint64, baseFee *big.Int) (*core.Message, error) {
	// Reject invalid combinations of pre- and post-1559 fee styles
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return nil, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}
	// Set sender address or use zero address if none specified.
	addr := args.from()

	// Set default gas & gas price if none were set
	gas := globalGasCap
	if gas == 0 {
		gas = uint64(math.MaxUint64 / 2)
	}
	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}
	if globalGasCap != 0 && globalGasCap < gas {
		log.Warn("Caller gas above allowance, capping", "requested", gas, "cap", globalGasCap)
		gas = globalGasCap
	}
	var (
		gasPrice  *big.Int
		gasFeeCap *big.Int
		gasTipCap *big.Int
	)
	if baseFee == nil {
		// If there's no basefee, then it must be a non-1559 execution
		gasPrice = new(big.Int)
		if args.GasPrice != nil {
			gasPrice = args.GasPrice.ToInt()
		}
		gasFeeCap, gasTipCap = gasPrice, gasPrice
	} else {
		// A basefee is provided, necessitating 1559-type execution
		if args.GasPrice != nil {
			// User specified the legacy gas field, convert to 1559 gas typing
			gasPrice = args.GasPrice.ToInt()
			gasFeeCap, gasTipCap = gasPrice, gasPrice
		} else {
			// User specified 1559 gas fields (or none), use those
			gasFeeCap = new(big.Int)
			if args.MaxFeePerGas != nil {
				gasFeeCap = args.MaxFeePerGas.ToInt()
			}
			gasTipCap = new(big.Int)
			if args.MaxPriorityFeePerGas != nil {
				gasTipCap = args.MaxPriorityFeePerGas.ToInt()
			}
			// Backfill the legacy gasPrice for EVM execution, unless we're all zeroes
			gasPrice = new(big.Int)
			if gasFeeCap.BitLen() > 0 || gasTipCap.BitLen() > 0 {
				gasPrice = math.BigMin(new(big.Int).Add(gasTipCap, baseFee), gasFeeCap)
			}
		}
	}
	value := new(big.Int)
	if args.Value != nil {
		value = args.Value.ToInt()
	}
	data := args.data()
	var accessList types.AccessList
	if args.AccessList != nil {
		accessList = *args.AccessList
	}
	msg := &core.Message{
		From:              addr,
		To:                args.To,
		Value:             value,
		GasLimit:          gas,
		GasPrice:          gasPrice,
		GasFeeCap:         gasFeeCap,
		GasTipCap:         gasTipCap,
		Data:              data,
		AccessList:        accessList,
		SkipAccountChecks: true,
	}
	return msg, nil
}

// toTransaction converts the arguments to a transaction.
// This assumes that setDefaults has been called.
func (args *TransactionArgs) toTransaction() *types.Transaction {
	var data types.TxData
	switch {
	case args.MaxFeePerGas != nil:
		al := types.AccessList{}
		if args.AccessList != nil {
			al = *args.AccessList
		}
		data = &types.DynamicFeeTx{
			To:         args.To,
			ChainID:    (*big.Int)(args.ChainID),
			Nonce:      uint64(*args.Nonce),
			Gas:        uint64(*args.Gas),
			GasFeeCap:  (*big.Int)(args.MaxFeePerGas),
			GasTipCap:  (*big.Int)(args.MaxPriorityFeePerGas),
			Value:      (*big.Int)(args.Value),
			Data:       args.data(),
			AccessList: al,
		}
	case args.AccessList != nil:
		data = &types.AccessListTx{
			To:         args.To,
			ChainID:    (*big.Int)(args.ChainID),
			Nonce:      uint64(*args.Nonce),
			Gas:        uint64(*args.Gas),
			GasPrice:   (*big.Int)(args.GasPrice),
			Value:      (*big.Int)(args.Value),
			Data:       args.data(),
			AccessList: *args.AccessList,
		}
	default:
		data = &types.LegacyTx{
			To:       args.To,
			Nonce:    uint64(*args.Nonce),
			Gas:      uint64(*args.Gas),
			GasPrice: (*big.Int)(args.GasPrice),
			Value:    (*big.Int)(args.Value),
			Data:     args.data(),
		}
	}
	return types.NewTx(data)
}

// ToTransaction converts the arguments to a transaction.
// This assumes that setDefaults has been called.
func (args *TransactionArgs) ToTransaction() *types.Transaction {
	return args.toTransaction()
}

func (args *TransactionArgs) toEVMTransaction() *txs.MsgEthereumTx {
	var (
		chainID, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas sdkmath.Int
		gas, nonce                                                   uint64
		from, to                                                     string
	)

	// Set sender address or use zero address if none specified.
	if args.ChainID != nil {
		chainID = sdkmath.NewIntFromBigInt(args.ChainID.ToInt())
	}

	if args.Nonce != nil {
		nonce = uint64(*args.Nonce)
	}

	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}

	if args.GasPrice != nil {
		gasPrice = sdkmath.NewIntFromBigInt(args.GasPrice.ToInt())
	}

	if args.MaxFeePerGas != nil {
		maxFeePerGas = sdkmath.NewIntFromBigInt(args.MaxFeePerGas.ToInt())
	}

	if args.MaxPriorityFeePerGas != nil {
		maxPriorityFeePerGas = sdkmath.NewIntFromBigInt(args.MaxPriorityFeePerGas.ToInt())
	}

	if args.Value != nil {
		value = sdkmath.NewIntFromBigInt(args.Value.ToInt())
	}

	if args.To != nil {
		to = args.To.Hex()
	}

	var data txs.TxData
	switch {
	case args.MaxFeePerGas != nil:
		al := txs.AccessList{}
		if args.AccessList != nil {
			al = txs.NewAccessList(args.AccessList)
		}

		data = &txs.DynamicFeeTx{
			To:        to,
			ChainID:   &chainID,
			Nonce:     nonce,
			GasLimit:  gas,
			GasFeeCap: &maxFeePerGas,
			GasTipCap: &maxPriorityFeePerGas,
			Amount:    &value,
			Data:      args.GetData(),
			Accesses:  al,
		}
	case args.AccessList != nil:
		data = &txs.AccessListTx{
			To:       to,
			ChainID:  &chainID,
			Nonce:    nonce,
			GasLimit: gas,
			GasPrice: &gasPrice,
			Amount:   &value,
			Data:     args.GetData(),
			Accesses: txs.NewAccessList(args.AccessList),
		}
	default:
		data = &txs.LegacyTx{
			To:       to,
			Nonce:    nonce,
			GasLimit: gas,
			GasPrice: &gasPrice,
			Amount:   &value,
			Data:     args.GetData(),
		}
	}

	anyData, err := txs.PackTxData(data)
	if err != nil {
		return nil
	}

	if args.From != nil {
		from = args.From.Hex()
	}

	msg := txs.MsgEthereumTx{
		Data: anyData,
		From: from,
	}
	msg.Hash = msg.AsTransaction().Hash().Hex()
	return &msg
}

func (args *TransactionArgs) ToEVMTransaction() *txs.MsgEthereumTx {
	return args.toEVMTransaction()
}

func (args *TransactionArgs) GetData() []byte {
	if args.Input != nil {
		return *args.Input
	}
	if args.Data != nil {
		return *args.Data
	}
	return nil
}
