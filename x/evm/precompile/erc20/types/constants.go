package types

import "github.com/ethereum/go-ethereum/common"

const (
	Method_Transfer  = "transfer"
	Method_BalanceOf = "balanceOf"
	Method_Register  = "register"
)

var (
	True32Byte  = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	False32Byte = make([]byte, 32)

	PrecompiledAddress = common.HexToAddress("0x0000000000000000000000000000000000000101")
)
