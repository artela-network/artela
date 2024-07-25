package handle

import (
	"errors"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
)

type (
	// ProposalTxVerifier defines the interface that is implemented by BaseApp,
	// that any custom ABCI PrepareProposal and ProcessProposal handler can use
	// to verify a transaction.
	ProposalTxVerifier interface {
		PrepareProposalVerifyTx(tx sdk.Tx) ([]byte, error)
		ProcessProposalVerifyTx(txBz []byte) (sdk.Tx, error)
	}

	// DefaultProposalHandler defines the default ABCI PrepareProposal and
	// ProcessProposal handlers.
	ArtelaProposalHandler struct {
		mempool    mempool.Mempool
		txVerifier ProposalTxVerifier
	}
)

func NewArtelaProposalHandler(mp mempool.Mempool, txVerifier ProposalTxVerifier) ArtelaProposalHandler {
	return ArtelaProposalHandler{
		mempool:    mp,
		txVerifier: txVerifier,
	}
}

// PrepareProposalHandler returns the default implementation for processing an
// ABCI proposal. The application's mempool is enumerated and all valid
// transactions are added to the proposal. Transactions are valid if they:
//
// 1) Successfully encode to bytes.
// 2) Are valid (i.e. pass runTx, AnteHandler only).
//
// Enumeration is halted once RequestPrepareProposal.MaxBytes of transactions is
// reached or the mempool is exhausted.
//
// Note:
//
// - Step (2) is identical to the validation step performed in
// DefaultProcessProposal. It is very important that the same validation logic
// is used in both steps, and applications must ensure that this is the case in
// non-default handlers.
//
// - If no mempool is set or if the mempool is a no-op mempool, the transactions
// requested from Tendermint will simply be returned, which, by default, are in
// FIFO order.
func (h ArtelaProposalHandler) PrepareProposalHandler() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req abci.RequestPrepareProposal) abci.ResponsePrepareProposal {
		// If the mempool is nil or NoOp we simply return the transactions
		// requested from CometBFT, which, by default, should be in FIFO order.
		_, isNoOp := h.mempool.(mempool.NoOpMempool)
		if h.mempool == nil || isNoOp {
			return abci.ResponsePrepareProposal{Txs: req.Txs}
		}

		var (
			selectedTxs  [][]byte
			totalTxBytes int64
		)

		iterator := h.mempool.Select(ctx, req.Txs)

		for iterator != nil {
			memTx := iterator.Tx()

			// NOTE: Since transaction verification was already executed in CheckTx,
			// which calls mempool.Insert, in theory everything in the pool should be
			// valid. But some mempool implementations may insert invalid txs, so we
			// check again.
			bz, err := h.txVerifier.PrepareProposalVerifyTx(memTx)
			if err != nil {
				err := h.mempool.Remove(memTx)
				if err != nil && !errors.Is(err, mempool.ErrTxNotFound) {
					panic(err)
				}
			} else {
				txSize := int64(len(bz))
				if totalTxBytes += txSize; totalTxBytes <= req.MaxTxBytes {
					selectedTxs = append(selectedTxs, bz)
				} else {
					// We've reached capacity per req.MaxTxBytes so we cannot select any
					// more transactions.
					break
				}
			}

			iterator = iterator.Next()
		}

		return abci.ResponsePrepareProposal{Txs: selectedTxs}
	}
}

// ProcessProposalHandler returns the default implementation for processing an
// ABCI proposal. Every transaction in the proposal must pass 2 conditions:
//
// 1. The transaction bytes must decode to a valid transaction.
// 2. The transaction must be valid (i.e. pass runTx, AnteHandler only)
//
// If any transaction fails to pass either condition, the proposal is rejected.
// Note that step (2) is identical to the validation step performed in
// DefaultPrepareProposal. It is very important that the same validation logic
// is used in both steps, and applications must ensure that this is the case in
// non-default handlers.
func (h ArtelaProposalHandler) ProcessProposalHandler() sdk.ProcessProposalHandler {
	// If the mempool is nil or NoOp we simply return ACCEPT,
	// because PrepareProposal may have included txs that could fail verification.
	_, isNoOp := h.mempool.(mempool.NoOpMempool)
	if h.mempool == nil || isNoOp {
		return NoOpProcessProposal()
	}

	return func(_ sdk.Context, req abci.RequestProcessProposal) abci.ResponseProcessProposal {
		for _, txBytes := range req.Txs {
			_, err := h.txVerifier.ProcessProposalVerifyTx(txBytes)
			if err != nil {
				return abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}
			}
		}

		return abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}
	}
}

// NoOpPrepareProposal defines a no-op PrepareProposal handler. It will always
// return the transactions sent by the client's request.
func NoOpPrepareProposal() sdk.PrepareProposalHandler {
	return func(_ sdk.Context, req abci.RequestPrepareProposal) abci.ResponsePrepareProposal {
		return abci.ResponsePrepareProposal{Txs: req.Txs}
	}
}

// NoOpProcessProposal defines a no-op ProcessProposal Handler. It will always
// return ACCEPT.
func NoOpProcessProposal() sdk.ProcessProposalHandler {
	return func(_ sdk.Context, _ abci.RequestProcessProposal) abci.ResponseProcessProposal {
		return abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}
	}
}
