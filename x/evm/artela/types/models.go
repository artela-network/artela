package types

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

type AspectInfo struct {
	AspectId common.Address `json:"AspectId"`
	Version  uint64         `json:"Version"`
	Priority int8           `json:"Priority"`
}

type AspectMeta struct {
	Id       common.Address `json:"id"`
	Version  *uint256.Int   `json:"version"`
	Priority int64          `json:"priority"`
}
type Property struct {
	Key   string `json:"Key"`
	Value []byte `json:"Value"`
}

type BoundAspectCode struct {
	AspectId common.Address `json:"aspectId"`
	Version  *uint256.Int   `json:"version"`
	Priority int64          `json:"priority"`
	Code     []byte         `json:"code"`
}

func NewBindingPriorityComparator(x []*AspectMeta) func(i, j int) bool {
	return func(i, j int) bool {
		return x[i].Priority < x[j].Priority && (bytes.Compare(x[i].Id.Bytes(), x[j].Id.Bytes()) < 0)
	}
}
