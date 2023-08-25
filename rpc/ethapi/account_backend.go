package ethapi

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type AccountBackend interface {
	Accounts() []common.Address
	NewAccount(password string) (common.AddressEIP55, error)
	ImportRawKey(privkey, password string) (common.Address, error)
	SignTransaction(args *TransactionArgs) (*ethtypes.Transaction, error)
}
