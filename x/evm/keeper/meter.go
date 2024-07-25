package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authmodule "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethereum "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"

	"github.com/artela-network/artela/x/evm/txs"
	"github.com/artela-network/artela/x/evm/types"
	"github.com/artela-network/aspect-core/djpm"
)

// GetEthIntrinsicGas returns the intrinsic gas cost for the transaction.
func (k *Keeper) GetEthIntrinsicGas(ctx cosmos.Context, msg *core.Message, cfg *params.ChainConfig, isContractCreation bool, isCustomVerification bool) (uint64, error) {
	blockHeight := big.NewInt(ctx.BlockHeight())

	homestead := cfg.IsHomestead(blockHeight)
	istanbul := cfg.IsIstanbul(blockHeight)

	// EIP3860(limit and meter initcode): https://eips.ethereum.org/EIPS/eip-3860
	intrinsic, err := core.IntrinsicGas(msg.Data, msg.AccessList, isContractCreation, homestead, istanbul, false)
	if err != nil {
		return 0, err
	}

	// for custom verification transaction we add an extra tx verification gas cost as intrinsic gas
	if isCustomVerification {
		intrinsic += djpm.MaxTxVerificationGas
	}

	return intrinsic, nil
}

// RefundGas transfers the leftover gas to the sender of the message, caped to half of the total gas
// consumed in the transaction. Additionally, the function sets the total gas consumed to the value
// returned by the EVM execution, thus ignoring the previous intrinsic gas consumed during in the
// AnteHandler.
func (k *Keeper) RefundGas(ctx cosmos.Context, msg *core.Message, leftoverGas uint64, denom string) error {
	// return EVM tokens for remaining gas, exchanged at the original rate.
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(leftoverGas), msg.GasPrice)

	switch remaining.Sign() {
	case -1:
		// negative refund errors
		return errorsmod.Wrapf(types.ErrInvalidRefund, "refunded amount value cannot be negative %d", remaining.Int64())
	case 1:
		// positive amount refund
		refundedCoins := cosmos.Coins{cosmos.NewCoin(denom, sdkmath.NewIntFromBigInt(remaining))}

		// refund to sender from the fee collector module account, which is the escrow account in charge of collecting tx fees

		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, authmodule.FeeCollectorName, msg.From.Bytes(), refundedCoins)
		if err != nil {
			err = errorsmod.Wrapf(errortypes.ErrInsufficientFunds, "fee collector account failed to refund fees: %s", err.Error())
			return errorsmod.Wrapf(err, "failed to refund %d leftover gas (%s)", leftoverGas, refundedCoins.String())
		}
	default:
		// no refund, consume gas and update the tx gas meter
	}

	return nil
}

// ResetGasMeterAndConsumeGas reset first the gas meter consumed value to zero and set it back to the new value
// 'gasUsed'
func (k *Keeper) ResetGasMeterAndConsumeGas(ctx cosmos.Context, gasUsed uint64) {
	// reset the gas count
	ctx.GasMeter().RefundGas(ctx.GasMeter().GasConsumed(), "reset the gas count")
	ctx.GasMeter().ConsumeGas(gasUsed, "apply evm transaction")
}

// DeductTxCostsFromUserBalance deducts the fees from the user balance. Returns an
// error if the specified sender address does not exist or the account balance is not sufficient.
func (k *Keeper) DeductTxCostsFromUserBalance(
	ctx cosmos.Context,
	fees cosmos.Coins,
	from common.Address,
) error {
	// fetch sender account
	signerAcc, err := authante.GetSignerAcc(ctx, k.accountKeeper, from.Bytes())
	if err != nil {
		return errorsmod.Wrapf(err, "account not found for sender %s", from)
	}

	// deduct the full gas cost from the user balance
	if err := authante.DeductFees(k.bankKeeper, ctx, signerAcc, fees); err != nil {
		return errorsmod.Wrapf(err, "failed to deduct full gas cost %s from the user %s balance", fees, from)
	}

	return nil
}

// ----------------------------------------------------------------------------
// 							        utils
// ----------------------------------------------------------------------------

// GasToRefund calculates the amount of gas the states machine should refund to the sender.
// It is capped by the refund quotient value(do not pass 0 to refundQuotient).
func GasToRefund(availableRefund, gasConsumed, refundQuotient uint64) uint64 {
	refund := gasConsumed / refundQuotient
	if refund > availableRefund {
		return availableRefund
	}

	return refund
}

// CheckSenderBalance validates that the tx cost value is positive and that the
// sender has enough funds to pay for the fees and value of the transaction.
func CheckSenderBalance(
	balance sdkmath.Int,
	txData txs.TxData,
) error {
	cost := txData.Cost()
	if cost.Sign() < 0 {
		return errorsmod.Wrapf(
			errortypes.ErrInvalidCoins,
			"tx cost (%s) is negative and invalid", cost,
		)
	}

	if balance.IsNegative() || balance.BigInt().Cmp(cost) < 0 {
		return errorsmod.Wrapf(
			errortypes.ErrInsufficientFunds,
			"sender balance < tx cost (%s < %s)", balance, txData.Cost(),
		)
	}

	return nil
}

// VerifyFee is used to return the fee for the given transaction data in cosmos.Coins. It checks that the
// gas limit is not reached, the gas limit is higher than the intrinsic gas and that the
// base fee is higher than the gas fee cap.
func VerifyFee(
	txData txs.TxData,
	denom string,
	baseFee *big.Int,
	homestead, istanbul, isCheckTx bool,
) (cosmos.Coins, error) {
	gasLimit := txData.GetGas()
	isContractCreation := txData.GetTo() == nil

	var accessList ethereum.AccessList
	if txData.GetAccessList() != nil {
		accessList = txData.GetAccessList()
	}

	intrinsicGas, err := core.IntrinsicGas(txData.GetData(), accessList, isContractCreation, homestead, istanbul, false)
	if err != nil {
		return nil, errorsmod.Wrapf(
			err,
			"failed to retrieve intrinsic gas, contract creation = %t; homestead = %t, istanbul = %t",
			isContractCreation, homestead, istanbul,
		)
	}

	// intrinsic gas verification during CheckTx
	if isCheckTx && gasLimit < intrinsicGas {
		return nil, errorsmod.Wrapf(
			errortypes.ErrOutOfGas,
			"gas limit too low: %d (gas limit) < %d (intrinsic gas)", gasLimit, intrinsicGas,
		)
	}

	if baseFee != nil && txData.GetGasFeeCap().Cmp(baseFee) < 0 {
		return nil, errorsmod.Wrapf(errortypes.ErrInsufficientFee,
			"the tx gasfeecap is lower than the tx baseFee: %s (gasfeecap), %s (basefee) ",
			txData.GetGasFeeCap(),
			baseFee)
	}

	feeAmt := txData.EffectiveFee(baseFee)
	if feeAmt.Sign() == 0 {
		// zero fee, no need to deduct
		return cosmos.Coins{}, nil
	}

	return cosmos.Coins{{Denom: denom, Amount: sdkmath.NewIntFromBigInt(feeAmt)}}, nil
}
