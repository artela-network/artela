package cosmos

import (
	"context"

	cosmos "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the exposed interface for using functionality of the bank keeper
// in the context of the cosmos AnteHandler package.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr cosmos.AccAddress, denom string) cosmos.Coin
	SendCoins(ctx context.Context, from, to cosmos.AccAddress, amt cosmos.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr cosmos.AccAddress, recipientModule string, amt cosmos.Coins) error
}
