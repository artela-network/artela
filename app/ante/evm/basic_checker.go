package evm

import (
	"math"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	anteutils "github.com/artela-network/artela/app/ante/utils"
	"github.com/artela-network/artela/app/interfaces"
	artela "github.com/artela-network/artela/ethereum/types"
	"github.com/artela-network/artela/x/evm/keeper"
	"github.com/artela-network/artela/x/evm/states"
	"github.com/artela-network/artela/x/evm/txs"
	"github.com/artela-network/artela/x/evm/txs/support"
	evmmodule "github.com/artela-network/artela/x/evm/types"
)

// EthAccountVerificationDecorator validates an account balance checks
type EthAccountVerificationDecorator struct {
	ak        evmmodule.AccountKeeper
	evmKeeper interfaces.EVMKeeper
}

// NewEthAccountVerificationDecorator creates a new EthAccountVerificationDecorator
func NewEthAccountVerificationDecorator(ak evmmodule.AccountKeeper, ek interfaces.EVMKeeper) EthAccountVerificationDecorator {
	return EthAccountVerificationDecorator{
		ak:        ak,
		evmKeeper: ek,
	}
}

// AnteHandle validates checks that the sender balance is greater than the total transaction cost.
// The account will be set to store if it doesn't exist, i.e. cannot be found on store.
// This AnteHandler decorator will fail if:
// - any of the msgs is not a MsgEthereumTx
// - from address is empty
// - account balance is lower than the transaction cost
func (avd EthAccountVerificationDecorator) AnteHandle(
	ctx cosmos.Context,
	tx cosmos.Tx,
	simulate bool,
	next cosmos.AnteHandler,
) (newCtx cosmos.Context, err error) {
	if !ctx.IsCheckTx() {
		return next(ctx, tx, simulate)
	}

	for i, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*txs.MsgEthereumTx)
		if !ok {
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*txs.MsgEthereumTx)(nil))
		}

		txData, err := txs.UnpackTxData(msgEthTx.Data)
		if err != nil {
			return ctx, errorsmod.Wrapf(err, "failed to unpack tx data any for tx %d", i)
		}

		// sender address should be in the tx cache from the previous AnteHandle call
		from := msgEthTx.GetFrom()
		if from.Empty() {
			return ctx, errorsmod.Wrap(errortypes.ErrInvalidAddress, "from address cannot be empty")
		}

		// check whether the sender address is EOA
		fromAddr := common.BytesToAddress(from)
		acct := avd.evmKeeper.GetAccount(ctx, fromAddr)

		if acct == nil {
			acc := avd.ak.NewAccountWithAddress(ctx, from)
			avd.ak.SetAccount(ctx, acc)
			acct = states.NewEmptyAccount()
		} else if acct.IsContract() {
			return ctx, errorsmod.Wrapf(errortypes.ErrInvalidType,
				"the sender is not EOA: address %s, codeHash <%s>", fromAddr, acct.CodeHash)
		}

		if err := keeper.CheckSenderBalance(sdkmath.NewIntFromBigInt(acct.Balance), txData); err != nil {
			return ctx, errorsmod.Wrap(err, "failed to check sender balance")
		}
	}
	return next(ctx, tx, simulate)
}

// EthGasConsumeDecorator validates enough intrinsic gas for the transaction and
// gas consumption.
type EthGasConsumeDecorator struct {
	bankKeeper         anteutils.BankKeeper
	distributionKeeper anteutils.DistributionKeeper
	evmKeeper          interfaces.EVMKeeper
	stakingKeeper      anteutils.StakingKeeper
	maxGasWanted       uint64
}

// NewEthGasConsumeDecorator creates a new EthGasConsumeDecorator
func NewEthGasConsumeDecorator(
	bankKeeper anteutils.BankKeeper,
	distributionKeeper anteutils.DistributionKeeper,
	evmKeeper interfaces.EVMKeeper,
	stakingKeeper anteutils.StakingKeeper,
	maxGasWanted uint64,
) EthGasConsumeDecorator {
	return EthGasConsumeDecorator{
		bankKeeper,
		distributionKeeper,
		evmKeeper,
		stakingKeeper,
		maxGasWanted,
	}
}

// AnteHandle validates that the Ethereum tx message has enough to cover intrinsic gas
// (during CheckTx only) and that the sender has enough balance to pay for the gas cost.
// If the balance is not sufficient, it will be attempted to withdraw enough staking rewards
// for the payment.
//
// Intrinsic gas for a transaction is the amount of gas that the transaction uses before the
// transaction is executed. The gas is a constant value plus any cost incurred by additional bytes
// of data supplied with the transaction.
//
// This AnteHandler decorator will fail if:
// - the message is not a MsgEthereumTx
// - sender account cannot be found
// - transaction's gas limit is lower than the intrinsic gas
// - user has neither enough balance nor staking rewards to deduct the transaction fees (gas_limit * gas_price)
// - transaction or block gas meter runs out of gas
// - sets the gas meter limit
// - gas limit is greater than the block gas meter limit
func (egcd EthGasConsumeDecorator) AnteHandle(ctx cosmos.Context, tx cosmos.Tx, simulate bool, next cosmos.AnteHandler) (cosmos.Context, error) {
	gasWanted := uint64(0)
	// gas consumption limit already checked during CheckTx so there's no need to
	// verify it again during ReCheckTx
	if ctx.IsReCheckTx() {
		// Use new context with gasWanted = 0
		// Otherwise, there's an error on txmempool.postCheck (tendermint)
		// that is not bubbled up. Thus, the Tx never runs on DeliverMode
		// Error: "gas wanted -1 is negative"
		// For more information, see issue #1554
		newCtx := ctx.WithGasMeter(artela.NewInfiniteGasMeterWithLimit(gasWanted))
		return next(newCtx, tx, simulate)
	}

	evmParams := egcd.evmKeeper.GetParams(ctx)
	evmDenom := evmParams.GetEvmDenom()
	chainCfg := evmParams.GetChainConfig()
	ethCfg := chainCfg.EthereumConfig(ctx.BlockHeight(), egcd.evmKeeper.ChainID())

	blockHeight := big.NewInt(ctx.BlockHeight())
	homestead := ethCfg.IsHomestead(blockHeight)
	istanbul := ethCfg.IsIstanbul(blockHeight)
	var events cosmos.Events

	// Use the lowest priority of all the messages as the final one.
	minPriority := int64(math.MaxInt64)
	baseFee := egcd.evmKeeper.GetBaseFee(ctx, ethCfg)

	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*txs.MsgEthereumTx)
		if !ok {
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*txs.MsgEthereumTx)(nil))
		}
		from := msgEthTx.GetFrom()

		txData, err := txs.UnpackTxData(msgEthTx.Data)
		if err != nil {
			return ctx, errorsmod.Wrap(err, "failed to unpack tx data")
		}

		if ctx.IsCheckTx() && egcd.maxGasWanted != 0 {
			// We can't trust the tx gas limit, because we'll refund the unused gas.
			if txData.GetGas() > egcd.maxGasWanted {
				gasWanted += egcd.maxGasWanted
			} else {
				gasWanted += txData.GetGas()
			}
		} else {
			gasWanted += txData.GetGas()
		}

		fees, err := keeper.VerifyFee(txData, evmDenom, baseFee, homestead, istanbul, ctx.IsCheckTx())
		if err != nil {
			return ctx, errorsmod.Wrapf(err, "failed to verify the fees")
		}

		// If the account balance is not sufficient, try to withdraw enough staking rewards
		// TODO artela staking.
		_ = from // TODO remove this after fix staking.
		// err = anteutils.ClaimStakingRewardsIfNecessary(ctx, egcd.bankKeeper, egcd.distributionKeeper, egcd.stakingKeeper, from, fees)
		// if err != nil {
		// 	return ctx, err
		// }
		fromAddr := msgEthTx.From
		err = egcd.evmKeeper.DeductTxCostsFromUserBalance(ctx, fees, common.HexToAddress(fromAddr))
		if err != nil {
			return ctx, errorsmod.Wrapf(err, "failed to deduct transaction costs from user balance")
		}

		events = append(events,
			cosmos.NewEvent(
				cosmos.EventTypeTx,
				cosmos.NewAttribute(cosmos.AttributeKeyFee, fees.String()),
			),
		)

		priority := txs.GetTxPriority(txData, baseFee)

		if priority < minPriority {
			minPriority = priority
		}
	}

	ctx.EventManager().EmitEvents(events)

	blockGasLimit := artela.BlockGasLimit(ctx)

	// return error if the tx gas is greater than the block limit (max gas)

	// NOTE: it's important here to use the gas wanted instead of the gas consumed
	// from the tx gas pool. The latter only has the value so far since the
	// EthSetupContextDecorator, so it will never exceed the block gas limit.
	if gasWanted > blockGasLimit {
		return ctx, errorsmod.Wrapf(
			errortypes.ErrOutOfGas,
			"tx gas (%d) exceeds block gas limit (%d)",
			gasWanted,
			blockGasLimit,
		)
	}

	// Set tx GasMeter with a limit of GasWanted (i.e. gas limit from the Ethereum tx).
	// The gas consumed will be then reset to the gas used by the states transition
	// in the EVM.

	// FIXME: use a custom gas configuration that doesn't add any additional gas and only
	// takes into account the gas consumed at the end of the EVM transaction.
	newCtx := ctx.
		WithGasMeter(artela.NewInfiniteGasMeterWithLimit(gasWanted)).
		WithPriority(minPriority)

	// we know that we have enough gas on the pool to cover the intrinsic gas
	return next(newCtx, tx, simulate)
}

// CanTransferDecorator checks if the sender is allowed to transfer funds according to the EVM block
// context rules.
type CanTransferDecorator struct {
	evmKeeper interfaces.EVMKeeper
}

// NewCanTransferDecorator creates a new CanTransferDecorator instance.
// NOTE: If this decorator is enabled, place it after the AspectRuntimeContextDecorator.
func NewCanTransferDecorator(evmKeeper interfaces.EVMKeeper) CanTransferDecorator {
	return CanTransferDecorator{
		evmKeeper: evmKeeper,
	}
}

// AnteHandle creates an EVM from the message and calls the BlockContext CanTransfer function to
// see if the address can execute the transaction.
func (ctd CanTransferDecorator) AnteHandle(ctx cosmos.Context, tx cosmos.Tx, simulate bool, next cosmos.AnteHandler) (cosmos.Context, error) {
	params := ctd.evmKeeper.GetParams(ctx)
	ethCfg := params.ChainConfig.EthereumConfig(ctx.BlockHeight(), ctd.evmKeeper.ChainID())

	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*txs.MsgEthereumTx)
		if !ok {
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*txs.MsgEthereumTx)(nil))
		}

		baseFee := ctd.evmKeeper.GetBaseFee(ctx, ethCfg)

		signer := ctd.evmKeeper.MakeSigner(ctx, msgEthTx.AsTransaction(),
			ethCfg, big.NewInt(ctx.BlockHeight()), uint64(ctx.BlockTime().Unix()))
		coreMsg, err := msgEthTx.AsMessage(signer, baseFee)
		if err != nil {
			return ctx, errorsmod.Wrapf(
				err,
				"failed to create an ethereum core.Message from signer %T", signer,
			)
		}

		if support.IsLondon(ethCfg, ctx.BlockHeight()) {
			if baseFee == nil {
				return ctx, errorsmod.Wrap(
					evmmodule.ErrInvalidBaseFee,
					"base fee is supported but evm block context value is nil",
				)
			}
			if coreMsg.GasFeeCap.Cmp(baseFee) < 0 {
				return ctx, errorsmod.Wrapf(
					errortypes.ErrInsufficientFee,
					"max fee per gas less than block base fee (%s < %s)",
					coreMsg.GasFeeCap, baseFee,
				)
			}
		}

		// NOTE: pass in an empty coinbase address and nil tracer as we don't need them for the check below
		cfg := &states.EVMConfig{
			ChainConfig: ethCfg,
			Params:      params,
			CoinBase:    common.Address{},
			BaseFee:     baseFee,
		}

		stateDB := states.New(ctx, ctd.evmKeeper, states.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash().Bytes())))
		evm := ctd.evmKeeper.NewEVM(ctx, coreMsg, cfg, txs.NewNoOpTracer(), stateDB)

		// check that caller has enough balance to cover asset transfer for **topmost** call
		// NOTE: here the gas consumed is from the context with the infinite gas meter
		if coreMsg.Value.Sign() > 0 && !evm.Context.CanTransfer(stateDB, coreMsg.From, coreMsg.Value) {
			return ctx, errorsmod.Wrapf(
				errortypes.ErrInsufficientFunds,
				"failed to transfer %s from address %s using the EVM block context transfer function",
				coreMsg.Value,
				coreMsg.From,
			)
		}
	}

	return next(ctx, tx, simulate)
}

// EthIncrementSenderSequenceDecorator increments the sequence of the signers.
type EthIncrementSenderSequenceDecorator struct {
	ak evmmodule.AccountKeeper
}

// NewEthIncrementSenderSequenceDecorator creates a new EthIncrementSenderSequenceDecorator.
func NewEthIncrementSenderSequenceDecorator(ak evmmodule.AccountKeeper) EthIncrementSenderSequenceDecorator {
	return EthIncrementSenderSequenceDecorator{
		ak: ak,
	}
}

// AnteHandle handles incrementing the sequence of the signer (i.e. sender). If the transaction is a
// contract creation, the nonce will be incremented during the transaction execution and not within
// this AnteHandler decorator.
func (issd EthIncrementSenderSequenceDecorator) AnteHandle(ctx cosmos.Context, tx cosmos.Tx, simulate bool, next cosmos.AnteHandler) (cosmos.Context, error) {
	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*txs.MsgEthereumTx)
		if !ok {
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*txs.MsgEthereumTx)(nil))
		}

		txData, err := txs.UnpackTxData(msgEthTx.Data)
		if err != nil {
			return ctx, errorsmod.Wrap(err, "failed to unpack tx data")
		}
		from := msgEthTx.GetFrom()
		// increase sequence of sender
		acc := issd.ak.GetAccount(ctx, from)
		if acc == nil {
			return ctx, errorsmod.Wrapf(
				errortypes.ErrUnknownAddress,
				"account %s is nil", common.BytesToAddress(msgEthTx.GetFrom().Bytes()),
			)
		}
		nonce := acc.GetSequence()
		// we merged the nonce verification to nonce increment, so when tx includes multiple messages
		// with same sender, they'll be accepted.
		if txData.GetNonce() != nonce {
			return ctx, errorsmod.Wrapf(
				errortypes.ErrInvalidSequence,
				"invalid nonce; got %d, expected %d", txData.GetNonce(), nonce,
			)
		}

		if err := acc.SetSequence(nonce + 1); err != nil {
			return ctx, errorsmod.Wrapf(err, "failed to set sequence to %d", acc.GetSequence()+1)
		}

		issd.ak.SetAccount(ctx, acc)
	}

	return next(ctx, tx, simulate)
}
