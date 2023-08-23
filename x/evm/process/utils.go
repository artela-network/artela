package process

import (
	"fmt"
	"github.com/artela-network/artela/x/evm/process/support"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"math/big"

	"github.com/cosmos/gogoproto/proto"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	amino = codec.NewLegacyAmino()
	// ModuleCdc references the global evm module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding.
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

	// AminoCdc is a amino codec created to support amino JSON compatible msgs.
	AminoCdc = codec.NewAminoCodec(amino)
)

const (
	// Amino names
	updateParamsName = "artela/MsgUpdateParams"
)

// NOTE: This is required for the GetSignBytes function
func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}

// RegisterInterfaces registers the client interfaces to protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*tx.TxExtensionOptionI)(nil),
		&ExtensionOptionsEthereumTx{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgEthereumTx{},
		&MsgUpdateParams{},
	)
	registry.RegisterInterface(
		"artela.evm.v1.TxData",
		(*TxData)(nil),
		&DynamicFeeTx{},
		&AccessListTx{},
		&LegacyTx{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// PackTxData constructs a new Any packed with the given process data value. It returns
// an error if the client state can't be casted to a protobuf message or if the concrete
// implementation is not registered to the protobuf codec.
func PackTxData(txData TxData) (*codectypes.Any, error) {
	msg, ok := txData.(proto.Message)
	if !ok {
		return nil, errorsmod.Wrapf(errortypes.ErrPackAny, "cannot proto marshal %T", txData)
	}

	anyTxData, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, errorsmod.Wrap(errortypes.ErrPackAny, err.Error())
	}

	return anyTxData, nil
}

// UnpackTxData unpacks an Any into a TxData. It returns an error if the
// client state can't be unpacked into a TxData.
func UnpackTxData(any *codectypes.Any) (TxData, error) {
	if any == nil {
		return nil, errorsmod.Wrap(errortypes.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	txData, ok := any.GetCachedValue().(TxData)
	if !ok {
		return nil, errorsmod.Wrapf(errortypes.ErrUnpackAny, "cannot unpack Any into TxData %T", any)
	}

	return txData, nil
}

// RegisterLegacyAminoCodec required for EIP-712
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdateParams{}, updateParamsName, nil)
}

// DefaultPriorityReduction is the default amount of price values required for 1 unit of priority.
// Because priority is `int64` while price is `big.Int`, it's necessary to scale down the range to keep it more pratical.
// The default value is the same as the `sdk.DefaultPowerReduction`.
var DefaultPriorityReduction = sdk.DefaultPowerReduction

var EmptyCodeHash = crypto.Keccak256(nil)

// DecodeTxResponse decodes an protobuf-encoded byte slice into TxResponse
func DecodeTxResponse(in []byte) (*MsgEthereumTxResponse, error) {
	var txMsgData sdk.TxMsgData
	if err := proto.Unmarshal(in, &txMsgData); err != nil {
		return nil, err
	}

	if len(txMsgData.MsgResponses) == 0 {
		return &MsgEthereumTxResponse{}, nil
	}

	var res MsgEthereumTxResponse
	if err := proto.Unmarshal(txMsgData.MsgResponses[0].Value, &res); err != nil {
		return nil, errorsmod.Wrap(err, "failed to unmarshal process response message data")
	}

	return &res, nil
}

// EncodeTransactionLogs encodes TransactionLogs slice into a protobuf-encoded byte slice.
func EncodeTransactionLogs(res *support.TransactionLogs) ([]byte, error) {
	return proto.Marshal(res)
}

// DecodeTransactionLogs decodes an protobuf-encoded byte slice into TransactionLogs
func DecodeTransactionLogs(data []byte) (support.TransactionLogs, error) {
	var logs support.TransactionLogs
	err := proto.Unmarshal(data, &logs)
	if err != nil {
		return support.TransactionLogs{}, err
	}
	return logs, nil
}

// UnwrapEthereumMsg extract MsgEthereumTx from wrapping sdk.Tx
func UnwrapEthereumMsg(tx *sdk.Tx, ethHash common.Hash) (*MsgEthereumTx, error) {
	if tx == nil {
		return nil, fmt.Errorf("invalid process: nil")
	}

	for _, msg := range (*tx).GetMsgs() {
		ethMsg, ok := msg.(*MsgEthereumTx)
		if !ok {
			return nil, fmt.Errorf("invalid process type: %T", tx)
		}
		txHash := ethMsg.AsTransaction().Hash()
		ethMsg.Hash = txHash.Hex()
		if txHash == ethHash {
			return ethMsg, nil
		}
	}

	return nil, fmt.Errorf("eth process not found: %s", ethHash)
}

// BinSearch execute the binary search and hone in on an executable gas limit
func BinSearch(lo, hi uint64, executable func(uint64) (bool, *MsgEthereumTxResponse, error)) (uint64, error) {
	for lo+1 < hi {
		mid := (hi + lo) / 2
		failed, _, err := executable(mid)
		// If the error is not nil(consensus error), it means the provided message
		// call or process will never be accepted no matter how much gas it is
		// assigned. Return the error directly, don't struggle any more.
		if err != nil {
			return 0, err
		}
		if failed {
			lo = mid
		} else {
			hi = mid
		}
	}
	return hi, nil
}

// EffectiveGasPrice compute the effective gas price based on eip-1159 rules
// `effectiveGasPrice = min(baseFee + tipCap, feeCap)`
func EffectiveGasPrice(baseFee, feeCap, tipCap *big.Int) *big.Int {
	return math.BigMin(new(big.Int).Add(tipCap, baseFee), feeCap)
}
