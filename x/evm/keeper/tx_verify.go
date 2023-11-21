package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/artela-network/artela/x/evm/states"
	"github.com/artela-network/aspect-core/djpm"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	ethereum "github.com/ethereum/go-ethereum/core/types"
)

func (k *Keeper) VerifySig(ctx cosmos.Context, tx *ethereum.Transaction) (common.Address, []byte, error) {
	// tx without signature
	txConfig := k.TxConfig(ctx, tx.Hash(), tx.Type())
	stateDB := states.New(ctx, k, txConfig)

	v, r, s := tx.RawSignatureValues()

	// calling a contract without signature, verify with aspect
	// this verification method is only allowed in call contract,
	// transactions that transfer value or creating contract must be signed
	if v == nil || r == nil || s == nil &&
		(tx.To() != nil && *tx.To() != common.Address{}) &&
		(stateDB.GetCodeHash(*tx.To()) != common.Hash{}) {
		return k.tryAspectVerifier(ctx, tx)
	}

	// tx with valid ec sig
	chainID := k.ChainID()
	evmParams := k.GetParams(ctx)
	chainCfg := evmParams.GetChainConfig()
	ethCfg := chainCfg.EthereumConfig(chainID)
	blockNum := big.NewInt(ctx.BlockHeight())
	signer := ethereum.MakeSigner(ethCfg, blockNum, uint64(ctx.BlockTime().Unix()))

	allowUnprotectedTxs := evmParams.GetAllowUnprotectedTxs()
	if !allowUnprotectedTxs && !tx.Protected() {
		return common.Address{}, nil, errorsmod.Wrapf(
			errortypes.ErrNotSupported,
			"rejected unprotected Ethereum transaction. Please EIP155 sign your transaction to protect it against replay-attacks")
	}
	sender, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, nil, errorsmod.Wrapf(
			errortypes.ErrorInvalidSigner,
			"couldn't retrieve sender address from the ethereum transaction: %s",
			err.Error(),
		)
	}

	return sender, nil, nil
}

func (k *Keeper) tryAspectVerifier(ctx cosmos.Context, tx *ethereum.Transaction) (common.Address, []byte, error) {
	return djpm.AspectInstance().GetSenderAndCallData(ctx.BlockHeight(), tx)
}
