package evm

import (
	"errors"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	ethereum "github.com/ethereum/go-ethereum/core/types"

	"github.com/artela-network/artela/app/interfaces"
	"github.com/artela-network/artela/x/evm/txs"
	evmmodule "github.com/artela-network/artela/x/evm/types"
)

// EthSetupContextDecorator is adapted from SetUpContextDecorator from cosmos-cosmos, it ignores gas consumption
// by setting the gas meter to infinite
type EthSetupContextDecorator struct {
	evmKeeper interfaces.EVMKeeper
}

func NewEthSetUpContextDecorator(evmKeeper interfaces.EVMKeeper) EthSetupContextDecorator {
	return EthSetupContextDecorator{
		evmKeeper: evmKeeper,
	}
}

func (esc EthSetupContextDecorator) AnteHandle(ctx cosmos.Context, tx cosmos.Tx, simulate bool, next cosmos.AnteHandler) (newCtx cosmos.Context, err error) {
	// all transactions must implement GasTx
	_, ok := tx.(authante.GasTx)
	if !ok {
		return ctx, errorsmod.Wrapf(errortypes.ErrInvalidType, "invalid transaction type %T, expected GasTx", tx)
	}

	// We need to setup an empty gas config so that the gas is consistent with Ethereum.
	newCtx = ctx.WithGasMeter(cosmos.NewInfiniteGasMeter()).
		WithKVGasConfig(storetypes.GasConfig{}).
		WithTransientKVGasConfig(storetypes.GasConfig{})

	// Reset transient gas used to prepare the execution of current cosmos tx.
	// Transient gas-used is necessary to sum the gas-used of cosmos tx, when it contains multiple eth msgs.
	esc.evmKeeper.ResetTransientGasUsed(ctx)
	return next(newCtx, tx, simulate)
}

// EthEmitEventDecorator emit events in ante handler in case of tx execution failed (out of block gas limit).
type EthEmitEventDecorator struct {
	evmKeeper interfaces.EVMKeeper
}

// NewEthEmitEventDecorator creates a new EthEmitEventDecorator
func NewEthEmitEventDecorator(evmKeeper interfaces.EVMKeeper) EthEmitEventDecorator {
	return EthEmitEventDecorator{evmKeeper}
}

// AnteHandle emits some basic events for the eth messages
func (eeed EthEmitEventDecorator) AnteHandle(ctx cosmos.Context, tx cosmos.Tx, simulate bool, next cosmos.AnteHandler) (newCtx cosmos.Context, err error) {
	// After eth tx passed ante handler, the fee is deducted and nonce increased, it shouldn't be ignored by json-rpc,
	// we need to emit some basic events at the very end of ante handler to be indexed by tendermint.
	txIndex := eeed.evmKeeper.GetTxIndexTransient(ctx)

	for i, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*txs.MsgEthereumTx)
		if !ok {
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*txs.MsgEthereumTx)(nil))
		}

		// emit ethereum tx hash as an event so that it can be indexed by Tendermint for query purposes
		// it's emitted in ante handler, so we can query failed transaction (out of block gas limit).
		ctx.EventManager().EmitEvent(cosmos.NewEvent(
			evmmodule.EventTypeEthereumTx,
			cosmos.NewAttribute(evmmodule.AttributeKeyEthereumTxHash, msgEthTx.Hash),
			cosmos.NewAttribute(evmmodule.AttributeKeyTxIndex, strconv.FormatUint(txIndex+uint64(i), 10)), // #nosec G701
		))
	}

	return next(ctx, tx, simulate)
}

// EthValidateBasicDecorator is adapted from ValidateBasicDecorator from cosmos-cosmos, it ignores ErrNoSignatures
type EthValidateBasicDecorator struct {
	evmKeeper interfaces.EVMKeeper
}

// NewEthValidateBasicDecorator creates a new EthValidateBasicDecorator
func NewEthValidateBasicDecorator(ek interfaces.EVMKeeper) EthValidateBasicDecorator {
	return EthValidateBasicDecorator{
		evmKeeper: ek,
	}
}

// AnteHandle handles basic validation of tx
func (vbd EthValidateBasicDecorator) AnteHandle(ctx cosmos.Context, tx cosmos.Tx, simulate bool, next cosmos.AnteHandler) (cosmos.Context, error) {
	// no need to validate basic on recheck tx, call next antehandler
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	err := tx.ValidateBasic()
	// ErrNoSignatures is fine with eth tx
	if err != nil && !errors.Is(err, errortypes.ErrNoSignatures) {
		return ctx, errorsmod.Wrap(err, "tx basic validation failed")
	}

	// For eth type cosmos tx, some fields should be verified as zero values,
	// since we will only verify the signature against the hash of the MsgEthereumTx.Data
	wrapperTx, ok := tx.(interfaces.ProtoTxProvider)
	if !ok {
		return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid tx type %T, didn't implement interface protoTxProvider", tx)
	}

	protoTx := wrapperTx.GetProtoTx()
	body := protoTx.Body
	if body.Memo != "" || body.TimeoutHeight != uint64(0) || len(body.NonCriticalExtensionOptions) > 0 {
		return ctx, errorsmod.Wrap(errortypes.ErrInvalidRequest,
			"for eth tx body Memo TimeoutHeight NonCriticalExtensionOptions should be empty")
	}

	if len(body.ExtensionOptions) != 1 {
		return ctx, errorsmod.Wrap(errortypes.ErrInvalidRequest, "for eth tx length of ExtensionOptions should be 1")
	}

	authInfo := protoTx.AuthInfo
	if len(authInfo.SignerInfos) > 0 {
		return ctx, errorsmod.Wrap(errortypes.ErrInvalidRequest, "for eth tx AuthInfo SignerInfos should be empty")
	}

	if authInfo.Fee.Payer != "" || authInfo.Fee.Granter != "" {
		return ctx, errorsmod.Wrap(errortypes.ErrInvalidRequest, "for eth tx AuthInfo Fee payer and granter should be empty")
	}

	sigs := protoTx.Signatures
	if len(sigs) > 0 {
		return ctx, errorsmod.Wrap(errortypes.ErrInvalidRequest, "for eth tx Signatures should be empty")
	}

	txFee := cosmos.Coins{}
	txGasLimit := uint64(0)

	evmParams := vbd.evmKeeper.GetParams(ctx)
	chainCfg := evmParams.GetChainConfig()
	chainID := vbd.evmKeeper.ChainID()
	ethCfg := chainCfg.EthereumConfig(ctx.BlockHeight(), chainID)
	baseFee := vbd.evmKeeper.GetBaseFee(ctx, ethCfg)
	enableCreate := evmParams.GetEnableCreate()
	enableCall := evmParams.GetEnableCall()
	evmDenom := evmParams.GetEvmDenom()

	for _, msg := range protoTx.GetMsgs() {
		msgEthTx, ok := msg.(*txs.MsgEthereumTx)
		if !ok {
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*txs.MsgEthereumTx)(nil))
		}

		// Validate `From` field
		if msgEthTx.From != "" {
			return ctx, errorsmod.Wrapf(errortypes.ErrInvalidRequest, "invalid From %s, expect empty string", msgEthTx.From)
		}

		txGasLimit += msgEthTx.GetGas()

		txData, err := txs.UnpackTxData(msgEthTx.Data)
		if err != nil {
			return ctx, errorsmod.Wrap(err, "failed to unpack MsgEthereumTx Data")
		}

		// return error if contract creation or call are disabled through governance
		if !enableCreate && txData.GetTo() == nil {
			return ctx, errorsmod.Wrap(evmmodule.ErrCreateDisabled, "failed to create new contract")
		} else if !enableCall && txData.GetTo() != nil {
			return ctx, errorsmod.Wrap(evmmodule.ErrCallDisabled, "failed to call contract")
		}

		if baseFee == nil && txData.TxType() == ethereum.DynamicFeeTxType {
			return ctx, errorsmod.Wrap(ethereum.ErrTxTypeNotSupported, "dynamic fee tx not supported")
		}

		txFee = txFee.Add(cosmos.Coin{Denom: evmDenom, Amount: sdkmath.NewIntFromBigInt(txData.Fee())})
	}

	if !authInfo.Fee.Amount.IsEqual(txFee) {
		return ctx, errorsmod.Wrapf(errortypes.ErrInvalidRequest, "invalid AuthInfo Fee Amount (%s != %s)", authInfo.Fee.Amount, txFee)
	}

	if authInfo.Fee.GasLimit != txGasLimit {
		return ctx, errorsmod.Wrapf(errortypes.ErrInvalidRequest, "invalid AuthInfo Fee GasLimit (%d != %d)", authInfo.Fee.GasLimit, txGasLimit)
	}

	return next(ctx, tx, simulate)
}
