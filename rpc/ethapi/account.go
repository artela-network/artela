package ethapi

import "github.com/ethereum/go-ethereum/common"

type AcctBackend interface {
	Accounts() []common.Address
}
